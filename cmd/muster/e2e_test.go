package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/mlaccetti/muster/internal/testutil"
)

var (
	binMu     sync.Mutex
	cachedBin string
)

func cliBin(t *testing.T) string {
	t.Helper()
	binMu.Lock()
	defer binMu.Unlock()
	if cachedBin != "" {
		return cachedBin
	}
	modRoot := testutil.RepoRoot(t)
	out := filepath.Join(os.TempDir(), "muster-e2e-test-bin")
	cmd := exec.Command("go", "build", "-o", out, "./cmd/muster")
	cmd.Dir = modRoot
	cmd.Env = append(os.Environ(),
		"GOMODCACHE="+filepath.Join(modRoot, ".cache", "gomod"),
		"GOCACHE="+filepath.Join(modRoot, ".cache", "gobuild"),
	)
	if err := cmd.Run(); err != nil {
		t.Fatalf("build muster: %v", err)
	}
	cachedBin = out
	return cachedBin
}

func runCLI(t *testing.T, env []string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(cliBin(t), args...)
	cmd.Env = env
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	code := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			code = ee.ExitCode()
		} else {
			t.Fatalf("run %v: %v", args, err)
		}
	}
	return outBuf.String(), errBuf.String(), code
}

func stubEnv(t *testing.T) []string {
	t.Helper()
	home := t.TempDir()
	bin, data := testutil.PrepareStubHarnessDirs(t)
	return testutil.StubEnv(t, home, bin, data)
}

func TestE2E_version(t *testing.T) {
	env := stubEnv(t)
	out, _, code := runCLI(t, env, "--version")
	if code != 0 || strings.TrimSpace(out) == "" {
		t.Fatalf("exit %d out=%q", code, out)
	}
}

func TestE2E_installValidateRender(t *testing.T) {
	env := stubEnv(t)
	if _, _, code := runCLI(t, env, "install"); code != 0 {
		t.Fatalf("install exit %d", code)
	}
	brief := testutil.Path(t, "testdata/briefs/minimal.json")
	if _, _, code := runCLI(t, env, "validate", "brief", brief); code != 0 {
		t.Fatalf("validate exit %d", code)
	}
	outPath := filepath.Join(t.TempDir(), "brief.txt")
	if _, _, code := runCLI(t, env, "render", "brief", "--template", "implement", "--input", brief, "--out", outPath); code != 0 {
		t.Fatalf("render exit %d", code)
	}
}

func TestE2E_workerHermeticRouting(t *testing.T) {
	env := stubEnv(t)
	dir := t.TempDir()
	report := filepath.Join(dir, "report.json")
	brief := testutil.Path(t, "testdata/briefs/minimal.json")
	out, _, code := runCLI(t, env,
		"worker", "--harness", "claude", "--brief", brief, "--report", report, "--sync",
	)
	if code != 0 {
		t.Fatalf("worker exit %d out=%q", code, out)
	}
	line := testutil.LastRoutingLine(t, out)
	if line.Verdict != "ok" {
		t.Fatalf("routing %+v", line)
	}
	if _, err := os.Stat(report); err != nil {
		t.Fatal(err)
	}
}

func TestE2E_dispatchWithoutRealHarness(t *testing.T) {
	home := t.TempDir()
	bin, data := testutil.PrepareStubHarnessDirs(t)
	_ = os.WriteFile(filepath.Join(data, "claude.routing.json"), []byte(
		`{"verdict":"ok","harness":"claude","family":"anthropic","model":"","report":"","reason_code":"none","reset":0}`,
	), 0o644)
	env := testutil.StubEnv(t, home, bin, data)
	dir := t.TempDir()
	pools := filepath.Join(dir, "pools.json")
	_ = os.WriteFile(pools, []byte(e2ePoolsJSON), 0o644)
	report := filepath.Join(dir, "report.json")
	brief := testutil.Path(t, "testdata/briefs/minimal.json")
	out, _, code := runCLI(t, env,
		"dispatch", "--tier", "implement", "--brief", brief, "--report", report, "--pools", pools,
	)
	if code != 0 {
		t.Fatalf("dispatch exit %d out=%q", code, out)
	}
	line := testutil.LastRoutingLine(t, out)
	if line.Verdict != "ok" {
		t.Fatalf("routing %+v", line)
	}
}

func TestE2E_workerRejectsMissingStubRoot(t *testing.T) {
	home := t.TempDir()
	env := testutil.StubEnv(t, home, "/nonexistent/stub/bin", "/nonexistent/stub/data")
	// No MUSTER_ADAPTER_ROOT effective stubs — also ensure PATH cannot find real CLIs.
	env = append(env, "PATH=/nonexistent")
	brief := testutil.Path(t, "testdata/briefs/minimal.json")
	report := filepath.Join(t.TempDir(), "report.json")
	_, _, code := runCLI(t, env,
		"worker", "--harness", "claude", "--brief", brief, "--report", report, "--sync",
	)
	if code == 0 {
		t.Fatal("expected worker failure without stub root")
	}
}

func TestE2E_passLedgerAndGateFailure(t *testing.T) {
	env := stubEnv(t)
	passDir := t.TempDir()
	spec := filepath.Join(passDir, "pass.yaml")
	brief := testutil.Path(t, "testdata/briefs/minimal.json")
	_ = os.WriteFile(spec, []byte(`id: e2e-pass
units:
  - name: unit1
    brief: `+brief+`
    harness: claude
gates:
  - name: fail
    command: "false"
    timeout_seconds: 5
`), 0o644)
	pools := testutil.Path(t, "kit/pools.json")
	_, _, code := runCLI(t, env, "pass", "run", "--file", spec, "--pass-dir", passDir, "--pools", pools)
	if code == 0 {
		t.Fatal("expected non-zero exit on gate failure")
	}
	ledger := filepath.Join(passDir, "ledger.jsonl")
	if _, err := os.Stat(ledger); err != nil {
		t.Fatal("ledger missing")
	}
}

const e2ePoolsJSON = `{
  "version": 1,
  "pass_order": ["claude"],
  "tiers": {
    "implement": [{"harness": "claude", "model": "", "family": "anthropic", "write": true}],
    "review": [{"harness": "claude", "model": "", "family": "anthropic", "write": false}],
    "plan": [{"harness": "claude", "model": "", "family": "anthropic", "write": false}],
    "cheap": [{"harness": "claude", "model": "haiku", "family": "anthropic", "write": false}]
  }
}`
