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

func claudeProjectDir(projectDir string) string {
	return filepath.Join(projectDir, ".claude")
}

func linkClaudeGlobal(kitDir string, dryRun bool) ([]string, error) {
	dir, err := claudeGlobalDir()
	if err != nil {
		return nil, err
	}
	touched, err := linkClaudeAgentWiring(dir, kitDir, dryRun)
	if err != nil {
		return nil, err
	}
	if !dryRun {
		_ = os.Remove(filepath.Join(dir, "muster-worker.md"))
		legacy := filepath.Join(dir, "CLAUDE.md")
		_ = patchFile(legacy, removeManagedBlock, false)
	}
	return touched, nil
}

func unlinkClaudeGlobal(dryRun bool) error {
	dir, err := claudeGlobalDir()
	if err != nil {
		return err
	}
	if err := unlinkClaudeAgentWiring(dir, dryRun); err != nil {
		return err
	}
	if !dryRun {
		_ = os.Remove(filepath.Join(dir, "muster-worker.md"))
		legacy := filepath.Join(dir, "CLAUDE.md")
		_ = patchFile(legacy, removeManagedBlock, false)
	}
	return nil
}

func linkClaudeProject(projectDir, kitDir string, dryRun bool) ([]string, error) {
	return linkClaudeAgentWiring(claudeProjectDir(projectDir), kitDir, dryRun)
}

func unlinkClaudeProject(projectDir string, dryRun bool) error {
	target := filepath.Join(projectDir, "CLAUDE.md")
	if err := patchFile(target, removeManagedBlock, dryRun); err != nil {
		return err
	}
	if !dryRun {
		_ = os.Remove(filepath.Join(projectDir, ".muster", "claude-worker.md"))
	}
	return unlinkClaudeAgentWiring(claudeProjectDir(projectDir), dryRun)
}
