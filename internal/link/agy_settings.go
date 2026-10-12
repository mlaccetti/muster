package link

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type agyLinkState struct {
	PreviousAgent string `json:"previous_agent,omitempty"`
}

func setAgyDefaultAgent(settingsPath, statePath string, dryRun bool) error {
	if dryRun {
		return nil
	}
	var doc map[string]interface{}
	data, err := os.ReadFile(settingsPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		if err := json.Unmarshal(data, &doc); err != nil {
			return err
		}
	}
	if doc == nil {
		doc = map[string]interface{}{}
	}
	previous := ""
	for _, key := range []string{"agent", "defaultAgent", "selectedAgent"} {
		if v, ok := doc[key]; ok {
			if s, ok := v.(string); ok && s != "" {
				previous = s
				break
			}
		}
	}
	if previous == orchestratorAgentName {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		return err
	}
	state := agyLinkState{PreviousAgent: previous}
	stateBytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(statePath, stateBytes, 0o644); err != nil {
		return err
	}
	doc["agent"] = orchestratorAgentName
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(settingsPath, out, 0o644)
}

func restoreAgyDefaultAgent(settingsPath, statePath string, dryRun bool) error {
	if dryRun {
		return nil
	}
	stateBytes, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var state agyLinkState
	if err := json.Unmarshal(stateBytes, &state); err != nil {
		return err
	}
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.Remove(statePath)
			return nil
		}
		return err
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	cur, _ := doc["agent"].(string)
	if cur != orchestratorAgentName {
		_ = os.Remove(statePath)
		return nil
	}
	if state.PreviousAgent == "" {
		delete(doc, "agent")
	} else {
		doc["agent"] = state.PreviousAgent
	}
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	if err := os.WriteFile(settingsPath, out, 0o644); err != nil {
		return err
	}
	_ = os.Remove(statePath)
	return nil
}
