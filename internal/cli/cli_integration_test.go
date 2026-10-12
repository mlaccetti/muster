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

func applyStubEnv(t *testing.T) (home string) {
	t.Helper()
	home = t.TempDir()
	bin, data := testutil.PrepareStubHarnessDirs(t)
	for _, e := range testutil.StubEnv(t, home, bin, data) {
		k, v, ok := splitEnv(e)
		if ok {
			t.Setenv(k, v)
		}
	}
	return home
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

func TestCLI_validateBrief(t *testing.T) {
	applyStubEnv(t)
	brief := testutil.Path(t, "testdata/briefs/minimal.json")
	if cli.RunForTest([]string{"validate", "brief", brief}, nil, nil) != 0 {
		t.Fatal("validate brief failed")
	}
}

func TestCLI_renderBrief(t *testing.T) {
	applyStubEnv(t)
	_ = cli.RunForTest([]string{"install"}, nil, nil)
	brief := testutil.Path(t, "testdata/briefs/minimal.json")
	outPath := filepath.Join(t.TempDir(), "brief.txt")
	var out bytes.Buffer
	code := cli.RunForTest([]string{
		"render", "brief", "--template", "implement", "--input", brief, "--out", outPath,
	}, &out, &out)
	if code != 0 {
		t.Fatalf("render exit %d: %s", code, out.String())
	}
	raw, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "Touch example.txt") {
		t.Fatalf("brief content: %s", raw)
	}
}

func TestCLI_workerStubRouting(t *testing.T) {
	applyStubEnv(t)
	dir := t.TempDir()
	report := filepath.Join(dir, "report.json")
	brief := testutil.Path(t, "testdata/briefs/minimal.json")
	var out bytes.Buffer
	code := cli.RunForTest([]string{
		"worker", "--harness", "claude", "--brief", brief, "--report", report, "--sync",
	}, &out, &out)
	if code != 0 {
		t.Fatalf("worker exit %d stderr/out: %s", code, out.String())
	}
	line := testutil.LastRoutingLine(t, out.String())
	if line.Verdict != "ok" || line.Harness == nil || *line.Harness != "claude" {
		t.Fatalf("routing %+v", line)
	}
	if _, err := os.Stat(report); err != nil {
		t.Fatal(err)
	}
}

func TestCLI_dispatchPoolWalk(t *testing.T) {
	applyStubEnv(t)
	data := testutil.StubDataDir(t)
	_ = os.WriteFile(filepath.Join(data, "claude.routing.json"), []byte(
		`{"verdict":"exhausted","harness":"claude","family":"anthropic","model":"","report":"","reason_code":"limit","reset":0}`,
	), 0o644)
	_ = os.WriteFile(filepath.Join(data, "codex.routing.json"), []byte(
		`{"verdict":"ok","harness":"codex","family":"openai","model":"m","report":"","reason_code":"none","reset":0}`,
	), 0o644)

	dir := t.TempDir()
	pools := filepath.Join(dir, "pools.json")
	_ = os.WriteFile(pools, []byte(testPoolsJSON), 0o644)
	report := filepath.Join(dir, "report.json")
	brief := testutil.Path(t, "testdata/briefs/minimal.json")
	var out bytes.Buffer
	code := cli.RunForTest([]string{
		"dispatch", "--tier", "implement", "--brief", brief, "--report", report, "--pools", pools,
	}, &out, &out)
	if code != 0 {
		t.Fatalf("dispatch exit %d: %s", code, out.String())
	}
	line := testutil.LastRoutingLine(t, out.String())
	if line.Verdict != "ok" || line.Harness == nil || *line.Harness != "codex" {
		t.Fatalf("routing %+v", line)
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

const testPoolsJSON = `{
  "version": 1,
  "pass_order": ["claude", "codex"],
  "tiers": {
    "implement": [
      {"harness": "claude", "model": "", "family": "anthropic", "write": true},
      {"harness": "codex", "model": "m", "family": "openai", "write": true}
    ],
    "review": [{"harness": "codex", "model": "m", "family": "openai", "write": false}],
    "plan": [{"harness": "claude", "model": "", "family": "anthropic", "write": false}],
    "cheap": [{"harness": "claude", "model": "haiku", "family": "anthropic", "write": false}]
  }
}`
