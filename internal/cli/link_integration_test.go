package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mlaccetti/muster/internal/cli"
	"github.com/mlaccetti/muster/internal/kit"
)

func TestCLI_linkProjectClaude(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("MUSTER_HOME", filepath.Join(home, "muster"))
	project := t.TempDir()
	if _, err := kit.Install(); err != nil {
		t.Fatal(err)
	}
	code := cli.RunForTest([]string{
		"link", "--harness", "claude", "--project-dir", project, "--global=false",
	}, nil, nil)
	if code != 0 {
		t.Fatalf("link exit %d", code)
	}
	raw, err := os.ReadFile(filepath.Join(project, ".claude", "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "muster-orchestrator") {
		t.Fatalf("project settings not linked: %s", raw)
	}
}
