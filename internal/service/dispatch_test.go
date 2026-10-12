package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mlaccetti/muster/internal/testutil"
)

func writePools(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "pools.json")
	body := `{
  "version": 1,
  "pass_order": ["claude", "codex", "cursor"],
  "tiers": {
    "implement": [
      {"harness": "claude", "model": "", "family": "anthropic", "write": true},
      {"harness": "codex", "model": "m", "family": "openai", "write": true}
    ],
    "review": [
      {"harness": "codex", "model": "m", "family": "openai", "write": false}
    ],
    "plan": [
      {"harness": "claude", "model": "", "family": "anthropic", "write": false}
    ],
    "cheap": [
      {"harness": "claude", "model": "haiku", "family": "anthropic", "write": false}
    ]
  }
}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestDispatchSuccessFirstCandidate(t *testing.T) {
	cleanup := testutil.SetupStubs(t)
	defer cleanup()
	dir := t.TempDir()
	pools := writePools(t, dir)
	report := filepath.Join(dir, "report.json")
	res, err := RunDispatch(context.Background(), DispatchRequest{
		Tier: "implement", BriefPath: testutil.Path(t, "testdata/briefs/minimal.json"), ReportPath: report,
		WorkDir: dir, PoolsPath: pools, Sync: true,
	})
	if err != nil || res.Line.Verdict != "ok" || res.ExitCode != 0 {
		t.Fatalf("res %+v err %v", res, err)
	}
}

func TestDispatchAvoidFamily(t *testing.T) {
	cleanup := testutil.SetupStubs(t)
	defer cleanup()
	dir := t.TempDir()
	pools := writePools(t, dir)
	data := testutil.StubDataDir(t)
	// codex returns ok when claude skipped
	_ = os.WriteFile(filepath.Join(data, "codex.routing.json"), []byte(
		`{"verdict":"ok","harness":"codex","family":"openai","model":"m","report":"","reason_code":"none","reset":0}`,
	), 0o644)
	report := filepath.Join(dir, "report.json")
	res, err := RunDispatch(context.Background(), DispatchRequest{
		Tier: "implement", BriefPath: testutil.Path(t, "testdata/briefs/minimal.json"), ReportPath: report,
		WorkDir: dir, PoolsPath: pools, AvoidFamilies: []string{"anthropic"}, Sync: true,
	})
	if err != nil || res.Line.Verdict != "ok" || res.Line.Harness == nil || *res.Line.Harness != "codex" {
		t.Fatalf("res %+v err %v", res, err)
	}
}

func TestDispatchExhausted(t *testing.T) {
	cleanup := testutil.SetupStubs(t)
	defer cleanup()
	dir := t.TempDir()
	pools := writePools(t, dir)
	data := testutil.StubDataDir(t)
	for _, h := range []string{"claude", "codex"} {
		_ = os.WriteFile(filepath.Join(data, h+".routing.json"), []byte(
			`{"verdict":"exhausted","harness":"`+h+`","family":"x","model":"","report":"","reason_code":"limit","reset":42}`,
		), 0o644)
	}
	report := filepath.Join(dir, "report.json")
	res, err := RunDispatch(context.Background(), DispatchRequest{
		Tier: "implement", BriefPath: testutil.Path(t, "testdata/briefs/minimal.json"), ReportPath: report,
		WorkDir: dir, PoolsPath: pools, Sync: true,
	})
	if err != nil || res.ExitCode != 3 || res.Line.Verdict != "exhausted" {
		t.Fatalf("res %+v err %v", res, err)
	}
}

func TestDispatchPoolOrder(t *testing.T) {
	cleanup := testutil.SetupStubs(t)
	defer cleanup()
	dir := t.TempDir()
	pools := writePools(t, dir)
	data := testutil.StubDataDir(t)
	_ = os.Remove(filepath.Join(data, "stub.calls"))
	_ = os.WriteFile(filepath.Join(data, "claude.routing.json"), []byte(
		`{"verdict":"exhausted","harness":"claude","family":"anthropic","model":"","report":"","reason_code":"limit","reset":0}`,
	), 0o644)
	_ = os.WriteFile(filepath.Join(data, "codex.routing.json"), []byte(
		`{"verdict":"ok","harness":"codex","family":"openai","model":"m","report":"","reason_code":"none","reset":0}`,
	), 0o644)
	report := filepath.Join(dir, "report.json")
	res, err := RunDispatch(context.Background(), DispatchRequest{
		Tier: "implement", BriefPath: testutil.Path(t, "testdata/briefs/minimal.json"), ReportPath: report,
		WorkDir: dir, PoolsPath: pools, Sync: true,
	})
	if err != nil || res.Line.Verdict != "ok" || *res.Line.Harness != "codex" {
		t.Fatalf("res %+v err %v", res, err)
	}
	log, err := os.ReadFile(filepath.Join(data, "stub.calls"))
	if err != nil {
		t.Fatal(err)
	}
	if !containsSeq(string(log), "name=claude", "name=codex") {
		t.Fatalf("call order missing: %s", log)
	}
}

func containsSeq(s string, a, b string) bool {
	return len(s) > 0 && indexOf(s, a) >= 0 && indexOf(s, b) > indexOf(s, a)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
