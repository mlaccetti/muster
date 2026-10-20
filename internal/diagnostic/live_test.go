package diagnostic

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type evidenceFixture struct {
	home, rootPath, rootID, childID, rootTurn, childTurn, nonce string
	events                                                      []byte
}

func TestValidateLiveEvidenceAcceptsActualTraceShape(t *testing.T) {
	f := validFixture(t)
	got, err := ValidateLiveEvidence(f.events, f.rootPath, f.home, f.rootID, f.nonce)
	if err != nil {
		t.Fatal(err)
	}
	if got.Agent != "muster-worker" || got.RootModel != "gpt-6.1-sol" || got.RootEffort != "high" || got.Model != "gpt-6-luna" || got.Effort != "medium" {
		t.Fatalf("unexpected model or role evidence: %+v", got)
	}
	if got.RootID != f.rootID || got.ChildID != f.childID || got.RootTurnID != f.rootTurn || got.ChildTurnID != f.childTurn {
		t.Fatalf("unexpected thread/turn linkage: %+v", got)
	}
}

func TestValidateLiveEvidenceFailsClosedForSpecificDefects(t *testing.T) {
	tests := []struct {
		name, expected string
		mutate         func(*testing.T, *evidenceFixture)
	}{
		{"failed root turn", "root turn did not complete", func(t *testing.T, f *evidenceFixture) {
			f.events = []byte("{\"type\":\"thread.started\",\"thread_id\":\"" + f.rootID + "\"}\n{\"type\":\"turn.completed\"}\n{\"type\":\"turn.failed\"}\n")
		}},
		{"ambiguous root thread", "root turn did not complete", func(t *testing.T, f *evidenceFixture) {
			f.events = []byte("{\"type\":\"thread.started\",\"thread_id\":\"" + f.rootID + "\"}\n{\"type\":\"thread.started\",\"thread_id\":\"00000000-0000-4000-8000-000000000099\"}\n{\"type\":\"turn.completed\"}\n")
		}},
		{"malformed event stream", "root turn did not complete", func(t *testing.T, f *evidenceFixture) { f.events = []byte("{bad json\n") }},
		{"wrong root model", "root model evidence", func(t *testing.T, f *evidenceFixture) {
			replaceOnce(t, f.rootPath, "gpt-6.1-sol", "gpt-6-astra")
		}},
		{"wrong child model", "child model or reasoning effort was missing or unexpected", func(t *testing.T, f *evidenceFixture) {
			replaceOnce(t, childPath(*f), "gpt-6-luna", "gpt-6-sol")
		}},
		{"wrong effort", "child model or reasoning effort was missing or unexpected", func(t *testing.T, f *evidenceFixture) {
			replaceOnce(t, childPath(*f), `"effort":"medium"`, `"effort":"high"`)
		}},
		{"child parent mismatch", "parent linkage is invalid", func(t *testing.T, f *evidenceFixture) {
			replaceOnce(t, childPath(*f), `"parent_thread_id":"`+f.rootID+`"`, `"parent_thread_id":"00000000-0000-4000-8000-000000000099"`)
		}},
		{"child session ID mismatch", "parent linkage is invalid", func(t *testing.T, f *evidenceFixture) {
			replaceOnce(t, childPath(*f), `"id":"`+f.childID+`"`, `"id":"00000000-0000-4000-8000-000000000099"`)
		}},
		{"nonce only in user message", "did not complete with the check nonce", func(t *testing.T, f *evidenceFixture) {
			replaceOnce(t, childPath(*f), `"role":"assistant"`, `"role":"user"`)
		}},
		{"missing nonce from assistant result", "did not complete with the check nonce", func(t *testing.T, f *evidenceFixture) {
			replaceOnce(t, childPath(*f), f.nonce, "different-value")
		}},
		{"mismatched result turn", "did not complete with the check nonce", func(t *testing.T, f *evidenceFixture) {
			replaceOnce(t, childPath(*f), `"turn_id":"`+f.childTurn+`"}}`, `"turn_id":"00000000-0000-4000-8000-000000000099"}}`)
		}},
		{"missing same-turn effort", "child model or reasoning effort was missing or unexpected", func(t *testing.T, f *evidenceFixture) {
			replaceOnce(t, childPath(*f), `,"effort":"medium"`, "")
		}},
		{"mismatched parent activity turn", "native child delegation evidence", func(t *testing.T, f *evidenceFixture) {
			replaceOnce(t, f.rootPath, `"type":"item_completed","turn_id":"`+f.rootTurn+`","item":{"type":"SubAgentActivity","kind":"completed"`, `"type":"item_completed","turn_id":"00000000-0000-4000-8000-000000000099","item":{"type":"SubAgentActivity","kind":"completed"`)
		}},
		{"ambiguous children", "ambiguous native child delegation evidence", func(t *testing.T, f *evidenceFixture) {
			data, err := os.ReadFile(f.rootPath)
			if err != nil {
				t.Fatal(err)
			}
			lines := strings.Split(strings.TrimSpace(string(data)), "\n")
			lines = append(lines, `{"type":"event_msg","payload":{"type":"item_completed","turn_id":"`+f.rootTurn+`","item":{"type":"SubAgentActivity","kind":"started","agent_thread_id":"00000000-0000-4000-8000-000000000099"}}}`)
			writeRecords(t, f.rootPath, lines)
		}},
		{"completion before delegation start", "ambiguous native child delegation evidence", func(t *testing.T, f *evidenceFixture) {
			data, err := os.ReadFile(f.rootPath)
			if err != nil {
				t.Fatal(err)
			}
			lines := strings.Split(strings.TrimSpace(string(data)), "\n")
			lines[2], lines[3] = lines[3], lines[2]
			writeRecords(t, f.rootPath, lines)
		}},
		{"task completion before final output", "did not complete with the check nonce", func(t *testing.T, f *evidenceFixture) {
			data, err := os.ReadFile(childPath(*f))
			if err != nil {
				t.Fatal(err)
			}
			lines := strings.Split(strings.TrimSpace(string(data)), "\n")
			if len(lines) != 4 {
				t.Fatalf("unexpected fixture record count: %d", len(lines))
			}
			lines[2], lines[3] = lines[3], lines[2]
			writeRecords(t, childPath(*f), lines)
		}},
		{"malformed child rollout", "parent linkage is invalid", func(t *testing.T, f *evidenceFixture) {
			writeRecords(t, childPath(*f), []string{"{malformed"})
		}},
		{"oversized child rollout", "parent linkage is invalid", func(t *testing.T, f *evidenceFixture) {
			if err := os.WriteFile(childPath(*f), []byte(strings.Repeat("x", 8<<20+1)), 0o600); err != nil {
				t.Fatal(err)
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := validFixture(t)
			tt.mutate(t, &f)
			_, err := ValidateLiveEvidence(f.events, f.rootPath, f.home, f.rootID, f.nonce)
			if err == nil || !strings.Contains(err.Error(), tt.expected) {
				t.Fatalf("expected error containing %q, got %v", tt.expected, err)
			}
		})
	}
}

func TestThreadIDMustBeStrictUUID(t *testing.T) {
	for _, id := range []string{"", "not-a-uuid", "00000000-0000-4000-8000-00000000000*", "00000000000040008000000000000000"} {
		if IsThreadID(id) {
			t.Fatalf("accepted invalid thread ID %q", id)
		}
	}
}

func validFixture(t *testing.T) evidenceFixture {
	t.Helper()
	f := evidenceFixture{
		home:      t.TempDir(),
		rootID:    "00000000-0000-4000-8000-000000000001",
		childID:   "00000000-0000-4000-8000-000000000002",
		rootTurn:  "00000000-0000-4000-8000-000000000011",
		childTurn: "00000000-0000-4000-8000-000000000022",
		nonce:     "nonce-789",
	}
	f.rootPath = writeRollout(t, f.home, f.rootID,
		`{"type":"session_meta","payload":{"id":"`+f.rootID+`","session_id":"`+f.rootID+`"}}`,
		`{"type":"turn_context","payload":{"turn_id":"`+f.rootTurn+`","model":"gpt-6.1-sol","effort":"high"}}`,
		`{"type":"event_msg","payload":{"type":"item_completed","turn_id":"`+f.rootTurn+`","item":{"type":"SubAgentActivity","kind":"started","agent_thread_id":"`+f.childID+`"}}}`,
		`{"type":"event_msg","payload":{"type":"item_completed","turn_id":"`+f.rootTurn+`","item":{"type":"SubAgentActivity","kind":"completed","agent_thread_id":"`+f.childID+`"}}}`)
	writeRollout(t, f.home, f.childID,
		`{"type":"session_meta","payload":{"id":"`+f.childID+`","session_id":"`+f.rootID+`","parent_thread_id":"`+f.rootID+`","agent_role":"muster-worker"}}`,
		`{"type":"turn_context","payload":{"turn_id":"`+f.childTurn+`","model":"gpt-6-luna","effort":"medium"}}`,
		`{"type":"response_item","payload":{"type":"message","role":"assistant","phase":"final_answer","content":[{"text":"done `+f.nonce+`"}],"internal_chat_message_metadata_passthrough":{"turn_id":"`+f.childTurn+`"}}}`,
		`{"type":"event_msg","payload":{"type":"task_complete","turn_id":"`+f.childTurn+`","root_turn_id":"`+f.rootTurn+`"}}`)
	f.events = []byte("{\"type\":\"thread.started\",\"thread_id\":\"" + f.rootID + "\"}\n{\"type\":\"turn.completed\",\"thread_id\":\"" + f.rootID + "\"}\n")
	return f
}

func childPath(f evidenceFixture) string {
	matches, _ := filepath.Glob(filepath.Join(f.home, "sessions", "*", "*", "*", "rollout-*-"+f.childID+".jsonl"))
	if len(matches) != 1 {
		panic("test fixture child path missing")
	}
	return matches[0]
}

func replaceOnce(t *testing.T, path, old, replacement string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), old) {
		t.Fatalf("fixture does not contain %q", old)
	}
	writeRecords(t, path, []string{strings.Replace(string(data), old, replacement, 1)})
}

func writeRecords(t *testing.T, path string, records []string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(strings.TrimSpace(strings.Join(records, "\n"))+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeRollout(t *testing.T, home, id string, records ...string) string {
	t.Helper()
	dir := filepath.Join(home, "sessions", "2026", "10", "19")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "rollout-2026-10-19T00-00-00-"+id+".jsonl")
	writeRecords(t, path, records)
	return path
}
