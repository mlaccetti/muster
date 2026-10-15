package link

import (
	"fmt"
	"os"
	"path/filepath"
)

type harnessWiring struct {
	id string
}

func (h harnessWiring) linkGlobal(kitDir string, dryRun bool) ([]string, error) {
	switch h.id {
	case "claude":
		return linkClaudeGlobal(kitDir, dryRun)
	case "codex":
		return linkCodexGlobal(kitDir, dryRun)
	case "cursor":
		return linkCursorGlobal(kitDir, dryRun)
	case "antigravity":
		return linkAntigravityGlobal(kitDir, dryRun)
	default:
		return nil, fmt.Errorf("unknown harness %q", h.id)
	}
}

func (h harnessWiring) unlinkGlobal(dryRun bool) error {
	switch h.id {
	case "claude":
		return unlinkClaudeGlobal(dryRun)
	case "codex":
		return unlinkCodexGlobal(dryRun)
	case "cursor":
		return unlinkCursorGlobal(dryRun)
	case "antigravity":
		return unlinkAntigravityGlobal(dryRun)
	default:
		return fmt.Errorf("unknown harness %q", h.id)
	}
}

func patchFile(path string, mutate func(string) string, dryRun bool) error {
	if dryRun {
		return nil
	}
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	body := ""
	if err == nil {
		body = string(existing)
	}
	newBody := mutate(body)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(newBody), 0o644)
}
