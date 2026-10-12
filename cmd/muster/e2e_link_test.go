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
	project := filepath.Join(home, "project")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatal(err)
	}
	env := testutil.StubEnv(t, home, "/nonexistent/bin", "/nonexistent/data")
	env = append(env,
		"MUSTER_HOME="+filepath.Join(home, "muster"),
		"HOME="+home,
		"GOPATH="+filepath.Join(home, "gopath"),
		"GOMODCACHE="+filepath.Join(home, "gomodcache"),
		"GOCACHE="+filepath.Join(home, "gocache"),
	)

	_, _, code := runCLI(t, env, "install")
	if code != 0 {
		t.Fatalf("install exit %d", code)
	}
	_, _, code = runCLI(t, env, "link", "--harness", "claude", "--project-dir", project, "--global=false")
	if code != 0 {
		t.Fatalf("link exit %d", code)
	}
	raw, err := os.ReadFile(filepath.Join(project, ".claude", "settings.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "muster-orchestrator") {
		t.Fatal(".claude/settings.json not linked")
	}
	_, _, code = runCLI(t, env, "doctor")
	if code != 0 {
		t.Fatalf("doctor exit %d", code)
	}
}
