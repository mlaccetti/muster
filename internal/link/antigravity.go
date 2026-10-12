package link

import (
	"os"
	"path/filepath"
)

func agySettingsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gemini", "antigravity-cli"), nil
}

func agyAgentsRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gemini", "config", "agents"), nil
}

func linkAntigravityGlobal(kitDir string, dryRun bool) ([]string, error) {
	agentsRoot, err := agyAgentsRoot()
	if err != nil {
		return nil, err
	}
	var touched []string

	orchAgent := filepath.Join(agentsRoot, orchestratorAgentName, "agent.md")
	if err := linkKitFileDryRun(orchAgent, dryRun, kitDir, "agy", orchestratorAgentName, "agent.md"); err != nil {
		return nil, err
	}
	touched = append(touched, orchAgent)

	workerAgent := filepath.Join(agentsRoot, workerAgentName, "agent.md")
	if err := linkKitFileDryRun(workerAgent, dryRun, kitDir, "agy", workerAgentName, "agent.md"); err != nil {
		return nil, err
	}
	touched = append(touched, workerAgent)

	settingsDir, err := agySettingsDir()
	if err != nil {
		return nil, err
	}
	settingsPath := filepath.Join(settingsDir, "settings.json")
	if err := setAgyDefaultAgent(settingsPath, filepath.Join(settingsDir, ".muster-link.json"), dryRun); err != nil {
		return nil, err
	}
	touched = append(touched, settingsPath)

	// Legacy ~/.antigravity prose files.
	if !dryRun {
		home, _ := os.UserHomeDir()
		legacy := filepath.Join(home, ".antigravity")
		_ = os.Remove(filepath.Join(legacy, "muster-orchestrator.md"))
		_ = os.Remove(filepath.Join(legacy, "muster-worker.md"))
		_ = os.Remove(filepath.Join(legacy, "muster-README.md"))
	}
	return touched, nil
}

func unlinkAntigravityGlobal(dryRun bool) error {
	agentsRoot, err := agyAgentsRoot()
	if err != nil {
		return err
	}
	if !dryRun {
		_ = removeLink(filepath.Join(agentsRoot, orchestratorAgentName, "agent.md"))
		_ = os.Remove(filepath.Join(agentsRoot, orchestratorAgentName))
		_ = removeLink(filepath.Join(agentsRoot, workerAgentName, "agent.md"))
		_ = os.Remove(filepath.Join(agentsRoot, workerAgentName))
	}
	settingsDir, err := agySettingsDir()
	if err != nil {
		return err
	}
	return restoreAgyDefaultAgent(filepath.Join(settingsDir, "settings.json"), filepath.Join(settingsDir, ".muster-link.json"), dryRun)
}

func linkAntigravityProject(projectDir, kitDir string, dryRun bool) ([]string, error) {
	dir := filepath.Join(projectDir, ".agents", "agents", orchestratorAgentName)
	target := filepath.Join(dir, "agent.md")
	if err := linkKitFileDryRun(target, dryRun, kitDir, "agy", orchestratorAgentName, "agent.md"); err != nil {
		return nil, err
	}
	return []string{target}, nil
}

func unlinkAntigravityProject(projectDir string, dryRun bool) error {
	if !dryRun {
		_ = removeLink(filepath.Join(projectDir, ".agents", "agents", orchestratorAgentName, "agent.md"))
		_ = os.RemoveAll(filepath.Join(projectDir, ".agents", "agents", orchestratorAgentName))
	}
	return nil
}
