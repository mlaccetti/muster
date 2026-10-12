package kit

import (
	"io/fs"
	"testing"
)

func TestEmbeddedKitLayout(t *testing.T) {
	for _, path := range []string{
		"policy/orchestrator.md",
		"policy/worker.md",
		"schemas/brief.v1.json",
		"schemas/routing.v1.json",
		"templates/implement.tmpl",
		"pools.json",
	} {
		if _, err := fs.Stat(Source, path); err != nil {
			t.Fatalf("%s: %v", path, err)
		}
	}
}
