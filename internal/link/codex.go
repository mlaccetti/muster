package link

import (
	"os"
	"path/filepath"
)

func codexHome() string {
	if v := os.Getenv("CODEX_HOME"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".codex")
}

func linkCodexGlobal(kitDir string, dryRun bool) ([]string, error) {
	dir := codexHome()
	if dir == "" {
		return nil, os.ErrNotExist
	}
	if !dryRun {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
		_ = removeLink(filepath.Join(dir, "muster-orchestrator.config.toml"))
		_ = removeLink(filepath.Join(dir, workerAgentName+".config.toml"))
	}
	var touched []string
	if err := unlinkTomlFile(filepath.Join(dir, "config.toml"), dryRun); err != nil {
		return nil, err
	}
	workerAgent := filepath.Join(dir, "agents", workerAgentName+".toml")
	if err := linkKitFileDryRun(workerAgent, dryRun, kitDir, "codex", workerAgentName+".config.toml"); err != nil {
		return nil, err
	}
	touched = append(touched, workerAgent)

	if err := unlinkCodexAgentsRole(dir, dryRun); err != nil {
		return nil, err
	}

	globalAgents := filepath.Join(dir, "AGENTS.md")
	if err := patchFile(globalAgents, func(s string) string { return upsertManagedBlock(s, orchestratorBlock(kitDir)) }, dryRun); err != nil {
		return nil, err
	}
	touched = append(touched, globalAgents)

	// Legacy prose cleanup.
	if !dryRun {
		_ = os.Remove(filepath.Join(dir, "muster-worker.md"))
	}
	return touched, nil
}

func unlinkCodexGlobal(dryRun bool) error {
	dir := codexHome()
	if err := unlinkTomlFile(filepath.Join(dir, "config.toml"), dryRun); err != nil {
		return err
	}
	if err := unlinkCodexAgentsRole(dir, dryRun); err != nil {
		return err
	}
	if !dryRun {
		_ = removeLink(filepath.Join(dir, "agents", workerAgentName+".toml"))
		_ = removeLink(filepath.Join(dir, "muster-orchestrator.config.toml"))
		_ = removeLink(filepath.Join(dir, workerAgentName+".config.toml"))
		_ = patchFile(filepath.Join(dir, "AGENTS.md"), removeManagedBlock, false)
		_ = os.Remove(filepath.Join(dir, "muster-worker.md"))
	}
	return nil
}
