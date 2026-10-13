package link

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const claudeLinkStateName = ".muster-link.json"

type claudeLinkState struct {
	PreviousAgent string `json:"previous_agent,omitempty"`
}

func claudeSettingsPath(dir string) string {
	return filepath.Join(dir, "settings.json")
}

func claudeLinkStatePath(dir string) string {
	return filepath.Join(dir, claudeLinkStateName)
}

func linkClaudeAgentWiring(claudeDir, kitDir string, dryRun bool) ([]string, error) {
	agentsDir := filepath.Join(claudeDir, "agents")
	var touched []string

	orchPath := filepath.Join(agentsDir, orchestratorAgentName+".md")
	if err := linkKitFileDryRun(orchPath, dryRun, kitDir, "claude", orchestratorAgentName+".md"); err != nil {
		return nil, err
	}
	touched = append(touched, orchPath)

	workerPath := filepath.Join(agentsDir, workerAgentName+".md")
	if err := linkKitFileDryRun(workerPath, dryRun, kitDir, "claude", workerAgentName+".md"); err != nil {
		return nil, err
	}
	touched = append(touched, workerPath)

	settingsPath := claudeSettingsPath(claudeDir)
	if dryRun {
		touched = append(touched, settingsPath)
		return touched, nil
	}
	if err := setClaudeDefaultAgent(settingsPath, claudeLinkStatePath(claudeDir), orchestratorAgentName); err != nil {
		return nil, err
	}
	touched = append(touched, settingsPath)
	return touched, nil
}

func unlinkClaudeAgentWiring(claudeDir string, dryRun bool) error {
	if !dryRun {
		_ = removeLink(filepath.Join(claudeDir, "agents", orchestratorAgentName+".md"))
		_ = removeLink(filepath.Join(claudeDir, "agents", workerAgentName+".md"))
	}
	if err := restoreClaudeDefaultAgent(claudeSettingsPath(claudeDir), claudeLinkStatePath(claudeDir), dryRun); err != nil {
		return err
	}
	if !dryRun {
		_ = os.Remove(claudeLinkStatePath(claudeDir))
	}
	return nil
}

func setClaudeDefaultAgent(settingsPath, statePath, agentName string) error {
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
	if v, ok := doc["agent"]; ok {
		if s, ok := v.(string); ok {
			previous = s
		}
	}
	if previous == agentName {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0o755); err != nil {
		return err
	}
	state := claudeLinkState{PreviousAgent: previous}
	stateBytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(statePath, stateBytes, 0o644); err != nil {
		return err
	}
	doc["agent"] = agentName
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(settingsPath, out, 0o644)
}

func restoreClaudeDefaultAgent(settingsPath, statePath string, dryRun bool) error {
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
	var state claudeLinkState
	if err := json.Unmarshal(stateBytes, &state); err != nil {
		return err
	}
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
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
	return os.WriteFile(settingsPath, out, 0o644)
}
