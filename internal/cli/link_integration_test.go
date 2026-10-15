package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mlaccetti/muster/internal/cli"
	"github.com/mlaccetti/muster/internal/kit"
)

func TestCLI_linkClaudeHarnessWide(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("MUSTER_HOME", filepath.Join(home, "muster"))
	if _, err := kit.Install(); err != nil {
		t.Fatal(err)
	}
	code := cli.RunForTest([]string{
		"link", "--harness", "claude",
	}, nil, nil)
	if code != 0 {
		t.Fatalf("link exit %d", code)
	}
	raw, err := os.ReadFile(filepath.Join(home, ".claude", "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "muster-orchestrator") {
		t.Fatalf("harness-wide settings not linked: %s", raw)
	}
	if _, err := os.Stat(filepath.Join(home, ".codex")); !os.IsNotExist(err) {
		t.Fatalf("explicit harness selection linked an unrelated harness, err=%v", err)
	}
}
