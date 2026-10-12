package buildversion

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestString_ldflags(t *testing.T) {
	t.Setenv("GOCACHE", t.TempDir())
	dir := t.TempDir()
	out := filepath.Join(dir, "muster")
	cmd := exec.Command("go", "build", "-o", out,
		"-ldflags", "-X github.com/mlaccetti/muster/internal/buildversion.version=v1.2.3 "+
			"-X github.com/mlaccetti/muster/internal/buildversion.commit=abc1234 "+
			"-X github.com/mlaccetti/muster/internal/buildversion.date=2026-01-01",
		"./cmd/muster")
	cmd.Dir = moduleRoot(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	got, err := exec.Command(out, "--version").Output()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(got), "v1.2.3") {
		t.Fatalf("version line %q", got)
	}
}

func TestNormalizeTag(t *testing.T) {
	if normalizeTag("0.1.0") != "v0.1.0" {
		t.Fatal()
	}
	if normalizeTag("v0.1.0") != "v0.1.0" {
		t.Fatal()
	}
}

func moduleRoot(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			t.Fatal("go.mod not found")
		}
		wd = parent
	}
}
