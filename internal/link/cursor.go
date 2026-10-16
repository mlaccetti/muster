package link

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const cursorManagedMarker = "muster-managed-plugin-v1\n"

func cursorPluginDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cursor", "plugins", "local", "muster-orchestrator"), nil
}

func cursorPluginFiles(kitDir string) (manifest, rule []byte, err error) {
	manifest, err = os.ReadFile(filepath.Join(kitDir, "link", "cursor", "plugin.json"))
	if err != nil {
		return nil, nil, err
	}
	template, err := os.ReadFile(filepath.Join(kitDir, "link", "cursor", "rules", "muster-orchestrator.mdc"))
	if err != nil {
		return nil, nil, err
	}
	policy, err := os.ReadFile(filepath.Join(kitDir, "policy", "orchestrator.md"))
	if err != nil {
		return nil, nil, err
	}
	if !bytes.Contains(template, []byte("<!-- MUSTER_POLICY -->")) {
		return nil, nil, fmt.Errorf("cursor rule template is missing the policy insertion marker")
	}
	rule = []byte(strings.Replace(string(template), "<!-- MUSTER_POLICY -->", string(policy), 1))
	return manifest, rule, nil
}

func linkCursorGlobal(kitDir string, dryRun bool) ([]string, error) {
	dir, err := cursorPluginDir()
	if err != nil {
		return nil, err
	}
	manifestPath := filepath.Join(dir, ".cursor-plugin", "plugin.json")
	rulePath := filepath.Join(dir, "rules", "muster-orchestrator.mdc")
	markerPath := filepath.Join(dir, ".muster-managed")
	manifest, rule, err := cursorPluginFiles(kitDir)
	if err != nil {
		return nil, err
	}
	if info, err := os.Lstat(dir); err == nil {
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("cursor plugin path exists and is not a Muster-owned directory: %s", dir)
		}
		marker, err := os.ReadFile(markerPath)
		if err != nil || string(marker) != cursorManagedMarker {
			return nil, fmt.Errorf("cursor plugin path already exists and is not managed by Muster: %s", dir)
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	if !dryRun {
		for _, path := range []string{filepath.Join(dir, ".cursor-plugin"), filepath.Join(dir, "rules")} {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return nil, err
			}
		}
		if err := os.WriteFile(markerPath, []byte(cursorManagedMarker), 0o644); err != nil {
			return nil, err
		}
		if err := os.WriteFile(manifestPath, manifest, 0o644); err != nil {
			return nil, err
		}
		if err := os.WriteFile(rulePath, rule, 0o644); err != nil {
			return nil, err
		}
	}
	return []string{manifestPath, rulePath}, nil
}

func unlinkCursorGlobal(dryRun bool) error {
	dir, err := cursorPluginDir()
	if err != nil {
		return err
	}
	markerPath := filepath.Join(dir, ".muster-managed")
	marker, err := os.ReadFile(markerPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if string(marker) != cursorManagedMarker || dryRun {
		return nil
	}
	for _, path := range []string{
		filepath.Join(dir, "rules", "muster-orchestrator.mdc"),
		filepath.Join(dir, ".cursor-plugin", "plugin.json"),
		markerPath,
	} {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	for _, path := range []string{filepath.Join(dir, "rules"), filepath.Join(dir, ".cursor-plugin"), dir} {
		entries, err := os.ReadDir(path)
		if os.IsNotExist(err) || (err == nil && len(entries) > 0) {
			continue
		}
		if err != nil {
			return err
		}
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func verifyCursor(kitDir string, _ []string) []string {
	dir, err := cursorPluginDir()
	if err != nil {
		return []string{err.Error()}
	}
	marker, err := os.ReadFile(filepath.Join(dir, ".muster-managed"))
	if err != nil || string(marker) != cursorManagedMarker {
		return []string{dir + ": Muster Cursor plugin is missing or unmanaged"}
	}
	manifest, rule, err := cursorPluginFiles(kitDir)
	if err != nil {
		return []string{err.Error()}
	}
	actualManifest, err := os.ReadFile(filepath.Join(dir, ".cursor-plugin", "plugin.json"))
	if err != nil || !json.Valid(actualManifest) || !bytes.Equal(actualManifest, manifest) {
		return []string{filepath.Join(dir, ".cursor-plugin", "plugin.json") + ": missing or stale"}
	}
	actualRule, err := os.ReadFile(filepath.Join(dir, "rules", "muster-orchestrator.mdc"))
	if err != nil || !bytes.Equal(actualRule, rule) {
		return []string{filepath.Join(dir, "rules", "muster-orchestrator.mdc") + ": missing or stale"}
	}
	return nil
}
