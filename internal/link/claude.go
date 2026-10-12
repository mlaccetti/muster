package link

import (
	"os"
	"path/filepath"
)

func claudeGlobalDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude"), nil
}

func linkClaudeGlobal(kitDir string, dryRun bool) ([]string, error) {
	dir, err := claudeGlobalDir()
	if err != nil {
		return nil, err
	}
	target := filepath.Join(dir, "CLAUDE.md")
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

func unlinkClaudeGlobal(dryRun bool) error {
	dir, err := claudeGlobalDir()
	if err != nil {
		return err
	}
	target := filepath.Join(dir, "CLAUDE.md")
	if err := patchFile(target, removeManagedBlock, dryRun); err != nil {
		return err
	}
	if !dryRun {
		_ = os.Remove(filepath.Join(dir, "muster-worker.md"))
	}
	return nil
}

func linkClaudeProject(projectDir, kitDir string, dryRun bool) ([]string, error) {
	target := filepath.Join(projectDir, "CLAUDE.md")
	block := orchestratorBlock(kitDir)
	if err := patchFile(target, func(s string) string { return upsertManagedBlock(s, block) }, dryRun); err != nil {
		return nil, err
	}
	worker := filepath.Join(projectDir, ".muster", "claude-worker.md")
	if err := writeFile(worker, workerBlock(kitDir), dryRun); err != nil {
		return nil, err
	}
	return []string{target, worker}, nil
}

func unlinkClaudeProject(projectDir string, dryRun bool) error {
	target := filepath.Join(projectDir, "CLAUDE.md")
	if err := patchFile(target, removeManagedBlock, dryRun); err != nil {
		return err
	}
	if !dryRun {
		_ = os.Remove(filepath.Join(projectDir, ".muster", "claude-worker.md"))
	}
	return nil
}
