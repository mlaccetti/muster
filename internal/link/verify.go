package link

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// VerifyReports returns human-readable problems for linked harnesses (empty if ok).
func VerifyReports(m *Manifest) []string {
	if m == nil {
		return []string{"link manifest missing"}
	}
	var issues []string
	for _, h := range m.Harnesses {
		issues = append(issues, verifyHarness(h.ID, m.KitDir)...)
	}
	return issues
}

func verifyHarness(id, kitDir string) []string {
	switch id {
	case "claude":
		return verifyClaude(kitDir)
	case "codex":
		return verifyCodex(kitDir)
	case "cursor":
		return verifyCursor(kitDir)
	case "antigravity":
		return verifyAgy(kitDir)
	default:
		return nil
	}
}

func verifySymlinkToKit(path, kitDir string, parts ...string) string {
	want := kitLinkFile(kitDir, parts...)
	if _, err := os.Lstat(path); err != nil {
		return path + ": missing"
	}
	target, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path + ": broken symlink"
	}
	if target != want {
		return path + ": not linked to active kit"
	}
	return ""
}

func verifyClaude(kitDir string) []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return []string{err.Error()}
	}
	dir := filepath.Join(home, ".claude")
	var issues []string
	if msg := verifySymlinkToKit(filepath.Join(dir, "agents", orchestratorAgentName+".md"), kitDir, "claude", orchestratorAgentName+".md"); msg != "" {
		issues = append(issues, msg)
	}
	data, err := os.ReadFile(filepath.Join(dir, "settings.json"))
	if err != nil {
		issues = append(issues, filepath.Join(dir, "settings.json")+": missing")
	} else if !strings.Contains(string(data), orchestratorAgentName) {
		issues = append(issues, "claude settings.json: agent is not "+orchestratorAgentName)
	}
	return issues
}

func verifyCodex(kitDir string) []string {
	dir := codexHome()
	var issues []string
	if msg := verifySymlinkToKit(filepath.Join(dir, codexProfileName+".config.toml"), kitDir, "codex", codexProfileName+".config.toml"); msg != "" {
		issues = append(issues, msg)
	}
	cfg, err := os.ReadFile(filepath.Join(dir, "config.toml"))
	if err != nil || !strings.Contains(string(cfg), tomlMarkerBegin) {
		issues = append(issues, filepath.Join(dir, "config.toml")+": muster block missing")
	}
	return issues
}

func verifyCursor(kitDir string) []string {
	dir, err := cursorRulesDir()
	if err != nil {
		return []string{err.Error()}
	}
	if msg := verifySymlinkToKit(filepath.Join(dir, "muster-orchestrator.mdc"), kitDir, "cursor", "muster-orchestrator.mdc"); msg != "" {
		return []string{msg}
	}
	return nil
}

func verifyAgy(kitDir string) []string {
	root, err := agyAgentsRoot()
	if err != nil {
		return []string{err.Error()}
	}
	var issues []string
	if msg := verifySymlinkToKit(filepath.Join(root, orchestratorAgentName, "agent.md"), kitDir, "agy", orchestratorAgentName, "agent.md"); msg != "" {
		issues = append(issues, msg)
	}
	settingsDir, err := agySettingsDir()
	if err != nil {
		return issues
	}
	data, err := os.ReadFile(filepath.Join(settingsDir, "settings.json"))
	if err != nil {
		issues = append(issues, filepath.Join(settingsDir, "settings.json")+": missing")
		return issues
	}
	var doc map[string]interface{}
	if json.Unmarshal(data, &doc) == nil {
		a, _ := doc["agent"].(string)
		if a != orchestratorAgentName {
			issues = append(issues, "agy settings.json: agent is not "+orchestratorAgentName)
		}
	}
	return issues
}
