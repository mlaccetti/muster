package link

import (
	"fmt"
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

func cursorOrchestratorRule(kitDir string) string {
	orch := kitDir + "/policy/orchestrator.md"
	work := kitDir + "/policy/worker.md"
	return fmt.Sprintf(`---
description: Muster orchestrator and worker policies
alwaysApply: true
---

# Muster (managed by muster link)

Orchestrator: read and follow %s

Workers: follow %s and return the structured handback defined there.

Harness dispatch and durable runs: use rckt-factory (`+"`rckt run`"+`) when installed.
`, orch, work)
}

func linkCursorGlobal(kitDir string, dryRun bool) ([]string, error) {
	dir, err := cursorRulesDir()
	if err != nil {
		return nil, err
	}
	target := filepath.Join(dir, "muster-orchestrator.mdc")
	if err := writeFile(target, cursorOrchestratorRule(kitDir), dryRun); err != nil {
		return nil, err
	}
	worker := filepath.Join(dir, "muster-worker.mdc")
	body := fmt.Sprintf(`---
description: Muster worker policy reference
alwaysApply: false
---

%s
`, workerBlock(kitDir))
	if err := writeFile(worker, body, dryRun); err != nil {
		return nil, err
	}
	return []string{target, worker}, nil
}

func unlinkCursorGlobal(dryRun bool) error {
	dir, err := cursorRulesDir()
	if err != nil {
		return err
	}
	if !dryRun {
		_ = os.Remove(filepath.Join(dir, "muster-orchestrator.mdc"))
		_ = os.Remove(filepath.Join(dir, "muster-worker.mdc"))
	}
	return nil
}

func linkCursorProject(projectDir, kitDir string, dryRun bool) ([]string, error) {
	dir := filepath.Join(projectDir, ".cursor", "rules")
	target := filepath.Join(dir, "muster-orchestrator.mdc")
	if err := writeFile(target, cursorOrchestratorRule(kitDir), dryRun); err != nil {
		return nil, err
	}
	return []string{target}, nil
}

func unlinkCursorProject(projectDir string, dryRun bool) error {
	if !dryRun {
		_ = os.Remove(filepath.Join(projectDir, ".cursor", "rules", "muster-orchestrator.mdc"))
	}
	return nil
}
