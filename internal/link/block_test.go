package link

import (
	"strings"
	"testing"
)

func TestUpsertManagedBlock_idempotent(t *testing.T) {
	block := orchestratorBlock("/kit")
	once := upsertManagedBlock("", block)
	twice := upsertManagedBlock(once, block)
	if strings.TrimSpace(once) != strings.TrimSpace(twice) {
		t.Fatalf("expected idempotent upsert\nonce=%q\ntwice=%q", once, twice)
	}
	if !contains(twice, markerStart) || !contains(twice, "/kit/policy/orchestrator.md") {
		t.Fatalf("block missing markers: %s", twice)
	}
}

func TestRemoveManagedBlock(t *testing.T) {
	block := orchestratorBlock("/kit")
	body := upsertManagedBlock("hello\n", block)
	clean := removeManagedBlock(body)
	if contains(clean, markerStart) || !contains(clean, "hello") {
		t.Fatalf("clean=%q", clean)
	}
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
