package link

import (
	"os"
	"path/filepath"
)

func cursorRulesDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cursor", "rules"), nil
}

func linkCursorGlobal(_ string, _ bool) ([]string, error) {
	// Cursor User Rules are configured in Cursor Settings, not in ~/.cursor/rules.
	// Global rule installation is not supported by the filesystem linker.
	return nil, nil
}

func unlinkCursorGlobal(dryRun bool) error {
	dir, err := cursorRulesDir()
	if err != nil {
		return err
	}
	if !dryRun {
		_ = removeLink(filepath.Join(dir, "muster-orchestrator.mdc"))
		_ = removeLink(filepath.Join(dir, "muster-worker.mdc"))
	}
	return nil
}

func linkCursorProject(projectDir, kitDir string, dryRun bool) ([]string, error) {
	dir := filepath.Join(projectDir, ".cursor", "rules")
	target := filepath.Join(dir, "muster-orchestrator.mdc")
	if err := linkKitFileDryRun(target, dryRun, kitDir, "cursor", "muster-orchestrator.mdc"); err != nil {
		return nil, err
	}
	return []string{target}, nil
}

func unlinkCursorProject(projectDir string, dryRun bool) error {
	if !dryRun {
		_ = removeLink(filepath.Join(projectDir, ".cursor", "rules", "muster-orchestrator.mdc"))
	}
	return nil
}
