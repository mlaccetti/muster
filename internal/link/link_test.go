package link

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mlaccetti/muster/internal/kit"
)

func TestLinkProjectClaudeOnly(t *testing.T) {
	home := t.TempDir()
	t.Setenv("MUSTER_HOME", filepath.Join(home, "muster"))
	project := t.TempDir()

	_, err := kit.Install()
	if err != nil {
		t.Fatal(err)
	}

	res, err := Run(Options{
		ProjectDir: project,
		Harnesses:  []string{"claude"},
		Global:     false,
		Project:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Harnesses) != 1 || res.Harnesses[0].ID != "claude" {
		t.Fatalf("res %+v", res)
	}
	claudeMD := filepath.Join(project, "CLAUDE.md")
	raw, err := os.ReadFile(claudeMD)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), markerStart) {
		t.Fatalf("CLAUDE.md not wired: %s", raw)
	}
	agents, err := os.ReadFile(filepath.Join(project, "AGENTS.md"))
	if err != nil || !strings.Contains(string(agents), markerStart) {
		t.Fatalf("AGENTS.md not wired: %v %s", err, agents)
	}
	if !strings.Contains(string(agents), "muster worker --harness") {
		t.Fatalf("AGENTS.md missing CLI dispatch: %s", agents)
	}
	if _, err := os.Stat(filepath.Join(project, ".muster", "orchestrator.path")); err != nil {
		t.Fatal(err)
	}
	m, err := LoadManifest()
	if err != nil || len(m.Harnesses) != 1 {
		t.Fatalf("manifest %+v err %v", m, err)
	}
}

func TestLinkAllWithoutGlobalPaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("MUSTER_HOME", filepath.Join(home, "muster"))
	project := t.TempDir()
	_, err := kit.Install()
	if err != nil {
		t.Fatal(err)
	}
	res, err := Run(Options{
		ProjectDir: project,
		All:        true,
		Global:     false,
		Project:    true,
		DryRun:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Harnesses) != 4 {
		t.Fatalf("expected 4 harnesses, got %d", len(res.Harnesses))
	}
}
