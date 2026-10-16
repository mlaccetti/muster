package link

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	markerStart = "<!-- MUSTER_LINK_START -->"
	markerEnd   = "<!-- MUSTER_LINK_END -->"
)

func orchestratorBlock(kitDir string) (string, error) {
	template, err := os.ReadFile(filepath.Join(kitDir, "link", "codex", "muster-orchestrator.md"))
	if err != nil {
		return "", err
	}
	rendered, err := renderOrchestrator(kitDir, "codex", template)
	if err != nil {
		return "", err
	}
	return string(rendered), nil
}

func upsertManagedBlock(existing string, block string) string {
	if strings.Contains(existing, markerStart) && strings.Contains(existing, markerEnd) {
		before, _, _ := strings.Cut(existing, markerStart)
		_, after, _ := strings.Cut(existing, markerEnd)
		after = strings.TrimPrefix(after, "\n")
		out := strings.TrimRight(before, "\n") + "\n\n" + block + "\n"
		if strings.TrimSpace(after) != "" {
			out += after
		}
		return out
	}
	trim := strings.TrimRight(existing, "\n")
	if trim == "" {
		return block + "\n"
	}
	return trim + "\n\n" + block + "\n"
}

func removeManagedBlock(existing string) string {
	if !strings.Contains(existing, markerStart) {
		return existing
	}
	before, _, _ := strings.Cut(existing, markerStart)
	_, after, ok := strings.Cut(existing, markerEnd)
	if !ok {
		return strings.TrimRight(before, "\n") + "\n"
	}
	after = strings.TrimPrefix(after, "\n")
	return strings.TrimRight(before, "\n") + "\n" + after
}
