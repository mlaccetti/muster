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
	target := filepath.Join(dir, "AGENTS.md")
	block := orchestratorBlock(kitDir)
	if err := patchFile(target, func(s string) string { return upsertManagedBlock(s, block) }, dryRun); err != nil {
		return nil, err
	}
	worker := filepath.Join(dir, "muster-worker.md")
	if err := writeFile(worker, workerBlock(kitDir), dryRun); err != nil {
		return nil, err
	}
	return []string{target, worker}, nil
}

func unlinkCodexGlobal(dryRun bool) error {
	dir := codexHome()
	target := filepath.Join(dir, "AGENTS.md")
	if err := patchFile(target, removeManagedBlock, dryRun); err != nil {
		return err
	}
	if !dryRun {
		_ = os.Remove(filepath.Join(dir, "muster-worker.md"))
	}
	return nil
}

func linkCodexProject(projectDir, kitDir string, dryRun bool) ([]string, error) {
	target := filepath.Join(projectDir, "AGENTS.md")
	block := orchestratorBlock(kitDir)
	if err := patchFile(target, func(s string) string { return upsertManagedBlock(s, block) }, dryRun); err != nil {
		return nil, err
	}
	worker := filepath.Join(projectDir, ".muster", "codex-worker.md")
	if err := writeFile(worker, workerBlock(kitDir), dryRun); err != nil {
		return nil, err
	}
	return []string{target, worker}, nil
}

func unlinkCodexProject(projectDir string, dryRun bool) error {
	target := filepath.Join(projectDir, "AGENTS.md")
	if err := patchFile(target, removeManagedBlock, dryRun); err != nil {
		return err
	}
	if !dryRun {
		_ = os.Remove(filepath.Join(projectDir, ".muster", "codex-worker.md"))
	}
	return nil
}
