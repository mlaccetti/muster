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

func linkCursorGlobal(kitDir string, dryRun bool) ([]string, error) {
	dir, err := cursorRulesDir()
	if err != nil {
		return nil, err
	}
	var targets []string
	orchRule := filepath.Join(dir, "muster-orchestrator.mdc")
	if err := linkKitFileDryRun(orchRule, dryRun, kitDir, "cursor", "muster-orchestrator.mdc"); err != nil {
		return nil, err
	}
	targets = append(targets, orchRule)
	workerRule := filepath.Join(dir, "muster-worker.mdc")
	if err := linkKitFileDryRun(workerRule, dryRun, kitDir, "cursor", "muster-worker.mdc"); err != nil {
		return nil, err
	}
	targets = append(targets, workerRule)
	return targets, nil
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
