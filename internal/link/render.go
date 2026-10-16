package link

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	policyMarker  = "<!-- MUSTER_POLICY -->"
	adapterMarker = "<!-- MUSTER_ADAPTER -->"
)

func renderOrchestrator(kitDir, harness string, template []byte) ([]byte, error) {
	policy, err := os.ReadFile(filepath.Join(kitDir, "policy", "orchestrator.md"))
	if err != nil {
		return nil, err
	}
	adapter, err := os.ReadFile(filepath.Join(kitDir, "policy", "adapters", harness+".md"))
	if err != nil {
		return nil, err
	}
	if strings.Count(string(template), policyMarker) != 1 || strings.Count(string(template), adapterMarker) != 1 {
		return nil, fmt.Errorf("%s orchestrator template must contain exactly one policy marker and one adapter marker", harness)
	}
	rendered := strings.Replace(string(template), policyMarker, strings.TrimSpace(string(policy)), 1)
	rendered = strings.Replace(rendered, adapterMarker, strings.TrimSpace(string(adapter)), 1)
	return []byte(rendered), nil
}

func renderOrchestratorFile(kitDir, harness, templateRel string, dryRun bool) (string, error) {
	template, err := os.ReadFile(filepath.Join(kitDir, filepath.FromSlash(templateRel)))
	if err != nil {
		return "", err
	}
	rendered, err := renderOrchestrator(kitDir, harness, template)
	if err != nil {
		return "", err
	}
	path := filepath.Join(kitDir, "rendered", harness, "muster-orchestrator.md")
	if !dryRun {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return "", err
		}
		if err := os.WriteFile(path, rendered, 0o644); err != nil {
			return "", err
		}
	}
	return path, nil
}

func verifyRenderedOrchestrator(kitDir, harness, templateRel string) string {
	path := filepath.Join(kitDir, "rendered", harness, "muster-orchestrator.md")
	template, err := os.ReadFile(filepath.Join(kitDir, filepath.FromSlash(templateRel)))
	if err != nil {
		return err.Error()
	}
	want, err := renderOrchestrator(kitDir, harness, template)
	if err != nil {
		return err.Error()
	}
	got, err := os.ReadFile(path)
	if err != nil || string(got) != string(want) {
		return path + ": missing or stale"
	}
	return ""
}

func linkRenderedOrchestrator(dest, kitDir, harness, templateRel string, dryRun bool) error {
	path, err := renderOrchestratorFile(kitDir, harness, templateRel, dryRun)
	if err != nil {
		return err
	}
	if dryRun {
		return nil
	}
	return replaceWithSymlink(dest, path)
}
