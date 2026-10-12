package link

import (
	"os"
	"path/filepath"
)

func antigravityConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// Antigravity (agy) user config; override with ANTIGRAVITY_HOME if set.
	if v := os.Getenv("ANTIGRAVITY_HOME"); v != "" {
		return v, nil
	}
	return filepath.Join(home, ".antigravity"), nil
}

func linkAntigravityGlobal(kitDir string, dryRun bool) ([]string, error) {
	dir, err := antigravityConfigDir()
	if err != nil {
		return nil, err
	}
	target := filepath.Join(dir, "muster-orchestrator.md")
	if err := writeFile(target, orchestratorBlock(kitDir), dryRun); err != nil {
		return nil, err
	}
	worker := filepath.Join(dir, "muster-worker.md")
	if err := writeFile(worker, workerBlock(kitDir), dryRun); err != nil {
		return nil, err
	}
	readme := filepath.Join(dir, "muster-README.md")
	instructions := "Load muster-orchestrator.md for orchestrator sessions and muster-worker.md for worker sessions (managed by muster link).\n"
	if err := writeFile(readme, instructions, dryRun); err != nil {
		return nil, err
	}
	return []string{target, worker, readme}, nil
}

func unlinkAntigravityGlobal(dryRun bool) error {
	dir, err := antigravityConfigDir()
	if err != nil {
		return err
	}
	if !dryRun {
		_ = os.Remove(filepath.Join(dir, "muster-orchestrator.md"))
		_ = os.Remove(filepath.Join(dir, "muster-worker.md"))
		_ = os.Remove(filepath.Join(dir, "muster-README.md"))
	}
	return nil
}

func linkAntigravityProject(projectDir, kitDir string, dryRun bool) ([]string, error) {
	dir := filepath.Join(projectDir, ".muster", "antigravity")
	target := filepath.Join(dir, "orchestrator.md")
	if err := writeFile(target, orchestratorBlock(kitDir), dryRun); err != nil {
		return nil, err
	}
	worker := filepath.Join(dir, "worker.md")
	if err := writeFile(worker, workerBlock(kitDir), dryRun); err != nil {
		return nil, err
	}
	return []string{target, worker}, nil
}

func unlinkAntigravityProject(projectDir string, dryRun bool) error {
	if !dryRun {
		_ = os.RemoveAll(filepath.Join(projectDir, ".muster", "antigravity"))
	}
	return nil
}
