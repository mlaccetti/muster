package main_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mlaccetti/muster/internal/cli"
	"github.com/mlaccetti/muster/internal/testutil"
)

func TestE2E_linkProjectNoRealHarness(t *testing.T) {
	home := t.TempDir()
	project := t.TempDir()
	bin := t.TempDir()
	for _, name := range []string{"claude", "codex", "cursor-agent", "agy"} {
		_ = os.WriteFile(filepath.Join(bin, name), []byte("#!/bin/sh\nexit 0\n"), 0o755)
	}
	for _, e := range testutil.StubEnv(t, home, bin, t.TempDir()) {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			t.Setenv(parts[0], parts[1])
		}
	}
	t.Setenv("MUSTER_HOME", filepath.Join(home, "muster"))

	var out bytes.Buffer
	if cli.RunForTest([]string{"install"}, &out, &out) != 0 {
		t.Fatalf("install: %s", out.String())
	}
	if cli.RunForTest([]string{
		"link", "--harness", "claude", "--project-dir", project, "--global=false",
	}, &out, &out) != 0 {
		t.Fatalf("link: %s", out.String())
	}
	raw, err := os.ReadFile(filepath.Join(project, ".claude", "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "muster-orchestrator") {
		t.Fatal(".claude/settings.json not linked")
	}
	if cli.RunForTest([]string{"doctor"}, &out, &out) != 0 {
		t.Fatalf("doctor: %s", out.String())
	}
}
