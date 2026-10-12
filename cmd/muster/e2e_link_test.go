package main_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mlaccetti/muster/internal/testutil"
)

func TestE2E_linkProjectNoRealHarness(t *testing.T) {
	home := t.TempDir()
	t.Setenv("MUSTER_HOME", filepath.Join(home, "muster"))
	project := t.TempDir()
	env := testutil.StubEnv(t, home, "/nonexistent/bin", "/nonexistent/data")
	env = append(env, "MUSTER_HOME="+filepath.Join(home, "muster"))

	_, _, code := runCLI(t, env, "install")
	if code != 0 {
		t.Fatalf("install exit %d", code)
	}
	_, _, code = runCLI(t, env, "link", "--harness", "claude", "--project-dir", project, "--global=false")
	if code != 0 {
		t.Fatalf("link exit %d", code)
	}
	raw, err := os.ReadFile(filepath.Join(project, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "MUSTER_LINK_START") {
		t.Fatal("CLAUDE.md not linked")
	}
	_, _, code = runCLI(t, env, "doctor")
	if code != 0 {
		t.Fatalf("doctor exit %d", code)
	}
}
