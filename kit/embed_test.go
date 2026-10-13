package kit

import (
	"testing"
)

func TestEmbedPolicy(t *testing.T) {
	for _, name := range []string{
		"policy/orchestrator.md",
		"policy/worker.md",
	} {
		if _, err := Source.ReadFile(name); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
}
