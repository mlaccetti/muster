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
	}
	var touched []string
	if err := ensureFile(filepath.Join(dir, "config.toml"), dryRun); err != nil {
		return nil, err
	}
	orchProfile := filepath.Join(dir, codexProfileName+".config.toml")
	if err := linkKitFileDryRun(orchProfile, dryRun, kitDir, "codex", codexProfileName+".config.toml"); err != nil {
		return nil, err
	}
	touched = append(touched, orchProfile)

	workerProfile := filepath.Join(dir, workerAgentName+".config.toml")
	if err := linkKitFileDryRun(workerProfile, dryRun, kitDir, "codex", workerAgentName+".config.toml"); err != nil {
		return nil, err
	}
	touched = append(touched, workerProfile)

	cfg := filepath.Join(dir, "config.toml")
	if err := patchTomlFile(cfg, codexManagedTomlBlock(kitDir), dryRun); err != nil {
		return nil, err
	}
	touched = append(touched, cfg)

	if err := patchCodexAgentsRole(dir, dryRun); err != nil {
		return nil, err
	}

	// Legacy prose cleanup.
	if !dryRun {
		_ = patchFile(filepath.Join(dir, "AGENTS.md"), removeManagedBlock, false)
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
		_ = removeLink(filepath.Join(dir, codexProfileName+".config.toml"))
		_ = removeLink(filepath.Join(dir, workerAgentName+".config.toml"))
		_ = patchFile(filepath.Join(dir, "AGENTS.md"), removeManagedBlock, false)
		_ = os.Remove(filepath.Join(dir, "muster-worker.md"))
	}
	return nil
}

func linkCodexProject(projectDir, kitDir string, dryRun bool) ([]string, error) {
	orchProfile := filepath.Join(projectDir, ".codex", codexProfileName+".config.toml")
	if err := linkKitFileDryRun(orchProfile, dryRun, kitDir, "codex", codexProfileName+".config.toml"); err != nil {
		return nil, err
	}
	return []string{orchProfile}, nil
}

func unlinkCodexProject(projectDir string, dryRun bool) error {
	if !dryRun {
		_ = removeLink(filepath.Join(projectDir, ".codex", codexProfileName+".config.toml"))
	}
	return nil
}
