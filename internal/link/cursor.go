package link

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var errManualHarnessSetup = errors.New("requires manual harness setup")

func cursorRulesDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cursor", "rules"), nil
}

func linkCursorGlobal(_ string, _ bool) ([]string, error) {
	// Cursor's documented global User Rules are edited in Cursor Settings, so
	// filesystem linking cannot configure this harness-wide setting.
	return nil, fmt.Errorf("global User Rules must be added in Cursor Settings; no supported CLI or file interface is documented: %w", errManualHarnessSetup)
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
