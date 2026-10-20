// Package diagnostic validates machine evidence emitted by Codex live checks.
package diagnostic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type LiveResult struct{ RootID, ChildID, RootTurnID, ChildTurnID, Agent, RootModel, RootEffort, Model, Effort string }

func RootThreadID(events []byte) (string, error) {
	for _, event := range decodeLines(events) {
		if event["type"] == "thread.started" {
			if id := stringAt(event, "thread_id"); id != "" {
				return id, nil
			}
		}
	}
	return "", errors.New("Codex did not emit a root thread id")
}

func ValidateLiveEvidence(events []byte, rootPath, home, rootID, nonce string) (LiveResult, error) {
	var out LiveResult
	stream := decodeLines(events)
	if malformedJSONL(events) || hasFailure(stream) || !hasCompletion(stream) || ambiguousThreadIDs(stream, rootID) {
		return out, errors.New("Codex root turn did not complete")
	}
	root, err := readRecords(rootPath)
	if err != nil {
		return out, errors.New("could not read root session metadata")
	}
	if sessionID(root) != rootID {
		return out, errors.New("Codex root session ID did not match the event stream")
	}
	childID, rootTurnID, ok := delegation(root, rootID)
	if !ok {
		return out, errors.New("Codex root session contains missing or ambiguous native child delegation evidence")
	}
	rootModel, rootEffort, rootTurnOK := turnSettings(root, rootTurnID)
	if !rootTurnOK || rootModel != "gpt-6.1-sol" || rootEffort == "" {
		return out, errors.New("Codex root model evidence was missing or unexpected")
	}
	childPath, err := findSession(home, childID)
	if err != nil {
		return out, errors.New("delegated child session evidence unavailable")
	}
	child, err := readRecords(childPath)
	if err != nil || sessionID(child) != childID || sessionParent(child) != rootID {
		return out, errors.New("delegated child parent linkage is invalid")
	}
	model, effort, childTurnID, childTurnOK := completedChildTurn(child, nonce, rootTurnID)
	agent := sessionAgentRole(child)
	if agent == "" {
		return out, errors.New("delegated child agent role was not recorded")
	}
	if !childTurnOK {
		return out, errors.New("delegated child did not complete with the check nonce")
	}
	if model != "gpt-6-luna" || effort != "medium" {
		return out, errors.New("delegated child model or reasoning effort was missing or unexpected")
	}
	out = LiveResult{RootID: rootID, ChildID: childID, RootTurnID: rootTurnID, ChildTurnID: childTurnID, Agent: agent, RootModel: rootModel, RootEffort: rootEffort, Model: model, Effort: effort}
	return out, nil
}

func decodeLines(data []byte) []map[string]any {
	var result []map[string]any
	s := bufio.NewScanner(bytes.NewReader(data))
	s.Buffer(make([]byte, 64*1024), 2<<20)
	for s.Scan() {
		var v map[string]any
		if json.Unmarshal(s.Bytes(), &v) == nil {
			result = append(result, v)
		}
	}
	return result
}

func malformedJSONL(data []byte) bool {
	s := bufio.NewScanner(bytes.NewReader(data))
	s.Buffer(make([]byte, 64*1024), 2<<20)
	for s.Scan() {
		var v map[string]any
		if json.Unmarshal(s.Bytes(), &v) != nil {
			return true
		}
	}
	return s.Err() != nil
}

func readRecords(path string) ([]map[string]any, error) {
	b, err := os.ReadFile(path)
	if err != nil || len(b) > 8<<20 {
		return nil, errors.New("session evidence unavailable")
	}
	if malformedJSONL(b) {
		return nil, errors.New("malformed session evidence")
	}
	records := decodeLines(b)
	if len(records) == 0 {
		return nil, errors.New("empty session evidence")
	}
	return records, nil
}

func findSession(home, id string) (string, error) {
	if !IsThreadID(id) {
		return "", errors.New("invalid child id")
	}
	paths, err := filepath.Glob(filepath.Join(home, "sessions", "*", "*", "*", "rollout-*-"+id+".jsonl"))
	if err != nil || len(paths) != 1 {
		return "", errors.New("child rollout not found")
	}
	return paths[0], nil
}

func IsThreadID(id string) bool {
	if len(id) != 36 {
		return false
	}
	for i, r := range id {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if r != '-' {
				return false
			}
			continue
		}
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}
	return true
}

func hasCompletion(events []map[string]any) bool {
	complete := false
	for _, e := range events {
		if e["type"] == "turn.completed" {
			complete = true
		}
	}
	return complete
}

func hasFailure(events []map[string]any) bool {
	for _, e := range events {
		switch e["type"] {
		case "turn.failed", "turn.error", "error":
			return true
		}
	}
	return false
}

func ambiguousThreadIDs(events []map[string]any, rootID string) bool {
	for _, e := range events {
		if id := stringAt(e, "thread_id"); id != "" && id != rootID {
			return true
		}
	}
	return false
}

func delegation(records []map[string]any, rootID string) (string, string, bool) {
	type activity struct {
		turn  string
		index int
	}
	started, completed := map[string]activity{}, map[string]activity{}
	startedCount, completedCount := 0, 0
	for index, record := range records {
		if record["type"] != "event_msg" {
			continue
		}
		p, _ := record["payload"].(map[string]any)
		item, _ := p["item"].(map[string]any)
		if item["type"] != "SubAgentActivity" {
			continue
		}
		id := stringAt(item, "agent_thread_id")
		turnID := stringAt(p, "turn_id")
		if !IsThreadID(id) || id == rootID || !IsThreadID(turnID) {
			continue
		}
		switch stringAt(item, "kind") {
		case "started":
			startedCount++
			started[id] = activity{turn: turnID, index: index}
		case "completed":
			completedCount++
			completed[id] = activity{turn: turnID, index: index}
		}
	}
	if startedCount != 1 || completedCount != 1 || len(started) != 1 || len(completed) != 1 {
		return "", "", false
	}
	for id, start := range started {
		finish, ok := completed[id]
		if ok && finish.turn == start.turn && finish.index > start.index {
			return id, start.turn, true
		}
	}
	return "", "", false
}

func sessionID(records []map[string]any) string {
	meta := sessionMeta(records)
	if id := stringAt(meta, "id"); id != "" {
		return id
	}
	return stringAt(meta, "session_id")
}

func sessionParent(records []map[string]any) string {
	meta := sessionMeta(records)
	if parent := stringAt(meta, "parent_thread_id"); parent != "" {
		return parent
	}
	source, _ := meta["source"].(map[string]any)
	subagent, _ := source["subagent"].(map[string]any)
	spawn, _ := subagent["thread_spawn"].(map[string]any)
	return stringAt(spawn, "parent_thread_id")
}

func sessionMeta(records []map[string]any) map[string]any {
	for _, r := range records {
		if r["type"] == "session_meta" {
			m, _ := r["payload"].(map[string]any)
			return m
		}
	}
	return nil
}

func turnSettings(records []map[string]any, turnID string) (string, string, bool) {
	model, effort, found := "", "", false
	for _, r := range records {
		if r["type"] == "turn_context" {
			p, _ := r["payload"].(map[string]any)
			if stringAt(p, "turn_id") == turnID {
				if found {
					return "", "", false
				}
				model, effort, found = stringAt(p, "model"), stringAt(p, "effort"), true
			}
		}
	}
	return model, effort, found
}

func sessionAgentRole(records []map[string]any) string {
	meta := sessionMeta(records)
	if role := strings.TrimSpace(stringAt(meta, "agent_role")); role != "" {
		return role
	}
	source, _ := meta["source"].(map[string]any)
	subagent, _ := source["subagent"].(map[string]any)
	spawn, _ := subagent["thread_spawn"].(map[string]any)
	return strings.TrimSpace(stringAt(spawn, "agent_role"))
}

func completedChildTurn(records []map[string]any, nonce, parentTurnID string) (string, string, string, bool) {
	contexts := map[string]int{}
	for i, r := range records {
		if r["type"] == "turn_context" {
			p, _ := r["payload"].(map[string]any)
			id := stringAt(p, "turn_id")
			if IsThreadID(id) {
				contexts[id] = i
			}
		}
	}
	var finalID, nonceID, completeID, completeRootTurnID string
	finalIndex, completeIndex, completionCount, finalCount := -1, -1, 0, 0
	for i, r := range records {
		p, _ := r["payload"].(map[string]any)
		if r["type"] == "event_msg" && p["type"] == "task_complete" {
			completionCount++
			completeID, completeIndex = stringAt(p, "turn_id"), i
			completeRootTurnID = stringAt(p, "root_turn_id")
		}
		if r["type"] == "response_item" && p["type"] == "message" {
			if p["role"] == "assistant" && p["phase"] == "final_answer" {
				finalCount++
				meta, _ := p["internal_chat_message_metadata_passthrough"].(map[string]any)
				turnID := stringAt(meta, "turn_id")
				if body, err := json.Marshal(p["content"]); err == nil && bytes.Contains(body, []byte(nonce)) {
					nonceID = turnID
				}
				finalID, finalIndex = turnID, i
			}
		}
	}
	if completionCount != 1 || finalCount != 1 || finalID == "" || nonceID != finalID || completeID != finalID || completeRootTurnID != parentTurnID || finalIndex < 0 || completeIndex <= finalIndex {
		return "", "", "", false
	}
	contextIndex, ok := contexts[finalID]
	if !ok || contextIndex >= finalIndex || len(contexts) != 1 {
		return "", "", "", false
	}
	model, effort, found := turnSettings(records, finalID)
	return model, effort, finalID, found
}

func stringAt(m map[string]any, key string) string { v, _ := m[key].(string); return v }
