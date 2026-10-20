package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mlaccetti/muster/internal/cli"
	"github.com/mlaccetti/muster/internal/testutil"
)

func applyStubEnv(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	bin := t.TempDir()
	for _, name := range []string{"claude", "codex", "cursor-agent", "agy"} {
		p := filepath.Join(bin, name)
		_ = os.WriteFile(p, []byte("#!/bin/sh\nexit 0\n"), 0o755)
	}
	for _, e := range testutil.StubEnv(t, home, bin, t.TempDir()) {
		k, v, ok := splitEnv(e)
		if ok {
			t.Setenv(k, v)
		}
	}
}

func splitEnv(e string) (key, val string, ok bool) {
	i := strings.IndexByte(e, '=')
	if i < 0 {
		return "", "", false
	}
	return e[:i], e[i+1:], true
}

func TestCLI_installAndKitPath(t *testing.T) {
	applyStubEnv(t)
	var out bytes.Buffer
	if cli.RunForTest([]string{"install"}, &out, &out) != 0 {
		t.Fatalf("install failed: %s", out.String())
	}
	out.Reset()
	if cli.RunForTest([]string{"kit", "path"}, &out, &out) != 0 {
		t.Fatalf("kit path: %s", out.String())
	}
	kitPath := strings.TrimSpace(out.String())
	if kitPath == "" {
		t.Fatal("empty kit path")
	}
	if _, err := os.Stat(filepath.Join(kitPath, "policy", "orchestrator.md")); err != nil {
		t.Fatal(err)
	}
}

func TestCLI_doctorSkipsProbe(t *testing.T) {
	applyStubEnv(t)
	var out bytes.Buffer
	if cli.RunForTest([]string{"doctor"}, &out, &out) != 0 {
		t.Fatalf("doctor: %s", out.String())
	}
	if !strings.Contains(out.String(), "probe: skipped") {
		t.Fatalf("expected skipped probe: %s", out.String())
	}
}

func TestCLI_doctorLiveRejectsSelfReportOnly(t *testing.T) {
	applyStubEnv(t)
	codex := filepath.Join(t.TempDir(), "codex")
	if err := os.WriteFile(codex, []byte("#!/bin/sh\nprintf '{\\\"type\\\":\\\"message\\\",\\\"text\\\":\\\"ok\\\"}\\n'\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if cli.RunForTest([]string{"doctor", "--live", "--codex-bin", codex}, &out, &out) == 0 {
		t.Fatalf("self-reporting stub unexpectedly passed: %s", out.String())
	}
	if strings.Contains(out.String(), "live: ok") {
		t.Fatalf("self-reporting stub reported success: %s", out.String())
	}
}

func TestCLI_doctorLiveHonorsTimeout(t *testing.T) {
	applyStubEnv(t)
	codex := filepath.Join(t.TempDir(), "codex")
	if err := os.WriteFile(codex, []byte("#!/bin/sh\nsleep 5\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if cli.RunForTest([]string{"doctor", "--live", "--live-timeout", "30ms", "--codex-bin", codex}, &out, &out) == 0 {
		t.Fatalf("timed out stub unexpectedly passed: %s", out.String())
	}
	if !strings.Contains(out.String(), "timed out") {
		t.Fatalf("expected timeout diagnostic: %s", out.String())
	}
}
