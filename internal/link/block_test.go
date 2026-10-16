package link

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpsertManagedBlock_idempotent(t *testing.T) {
	kitDir := t.TempDir()
	for path, contents := range map[string]string{
		"policy/orchestrator.md":            "# Shared policy\nShared delivery guidance.",
		"policy/adapters/codex.md":          "## Codex controls\nUse native controls.",
		"link/codex/muster-orchestrator.md": markerStart + "\n<!-- MUSTER_POLICY -->\n<!-- MUSTER_ADAPTER -->\n" + markerEnd,
	} {
		fullPath := filepath.Join(kitDir, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	block, err := orchestratorBlock(kitDir)
	if err != nil {
		t.Fatal(err)
	}
	once := upsertManagedBlock("", block)
	twice := upsertManagedBlock(once, block)
	if strings.TrimSpace(once) != strings.TrimSpace(twice) {
		t.Fatalf("expected idempotent upsert\nonce=%q\ntwice=%q", once, twice)
	}
	if !contains(twice, markerStart) || !contains(twice, "Shared delivery guidance.") || !contains(twice, "## Codex controls") || contains(twice, "MUSTER_POLICY") {
		t.Fatalf("block missing markers: %s", twice)
	}
}

func TestRemoveManagedBlock(t *testing.T) {
	block := markerStart + "\nmanaged\n" + markerEnd
	body := upsertManagedBlock("hello\n", block)
	clean := removeManagedBlock(body)
	if contains(clean, markerStart) || !contains(clean, "hello") {
		t.Fatalf("clean=%q", clean)
	}
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
