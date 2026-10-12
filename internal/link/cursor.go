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

func cursorOrchestratorMDC(kitDir string) string {
	return fmt.Sprintf(`---
description: Muster orchestrator (default agent policy)
alwaysApply: true
---

%s
`, orchestratorBlock(kitDir))
}

func cursorWorkerMDC(kitDir string) string {
	return fmt.Sprintf(`---
description: Muster worker policy reference
alwaysApply: false
---

%s
`, workerBlock(kitDir))
}

func linkCursorGlobal(kitDir string, dryRun bool) ([]string, error) {
	dir, err := cursorRulesDir()
	if err != nil {
		return nil, err
	}
	var targets []string
	orchRule := filepath.Join(dir, "muster-orchestrator.mdc")
	if err := writeFile(orchRule, cursorOrchestratorMDC(kitDir), dryRun); err != nil {
		return nil, err
	}
	targets = append(targets, orchRule)
	workerRule := filepath.Join(dir, "muster-worker.mdc")
	if err := writeFile(workerRule, cursorWorkerMDC(kitDir), dryRun); err != nil {
		return nil, err
	}
	targets = append(targets, workerRule)
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	agents := filepath.Join(home, ".cursor", "AGENTS.md")
	if err := patchOrchestratorAgents(agents, kitDir, dryRun); err != nil {
		return nil, err
	}
	targets = append(targets, agents)
	return targets, nil
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
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return unlinkOrchestratorAgents(filepath.Join(home, ".cursor", "AGENTS.md"), dryRun)
}

func linkCursorProject(projectDir, kitDir string, dryRun bool) ([]string, error) {
	dir := filepath.Join(projectDir, ".cursor", "rules")
	target := filepath.Join(dir, "muster-orchestrator.mdc")
	if err := writeFile(target, cursorOrchestratorMDC(kitDir), dryRun); err != nil {
		return nil, err
	}
	worker := filepath.Join(projectDir, ".muster", "cursor-worker.md")
	if err := writeFile(worker, workerBlock(kitDir), dryRun); err != nil {
		return nil, err
	}
	return []string{target, worker}, nil
}

func unlinkCursorProject(projectDir string, dryRun bool) error {
	if !dryRun {
		_ = os.Remove(filepath.Join(projectDir, ".cursor", "rules", "muster-orchestrator.mdc"))
		_ = os.Remove(filepath.Join(projectDir, ".muster", "cursor-worker.md"))
	}
	return nil
}
