package kit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mlaccetti/muster/internal/paths"
)

func TestInstallLayout(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("MUSTER_HOME", dir)
	got, err := Install()
	if err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{
		"policy/orchestrator.md",
		"policy/worker.md",
		"schemas/brief.v1.json",
		"templates/implement.tmpl",
	} {
		if _, err := os.Stat(filepath.Join(got, rel)); err != nil {
			t.Fatalf("%s: %v", rel, err)
		}
	}
	pools := filepath.Join(dir, "pools.json")
	if _, err := os.Stat(pools); err != nil {
		t.Fatal(err)
	}
	current, err := os.Readlink(filepath.Join(dir, "kit", "current"))
	if err != nil || current != paths.KitVersion {
		t.Fatalf("current link %q err %v", current, err)
	}
}
