package pass

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mlaccetti/muster/internal/testutil"
)

func TestPassRunGateFailure(t *testing.T) {
	cleanup := testutil.SetupStubs(t)
	defer cleanup()
	dir := t.TempDir()
	pools := testutil.Path(t, "kit/pools.json")
	spec := &Spec{
		ID: "minimal-pass",
		Units: []Unit{{
			Name: "unit1", Harness: "claude",
			Brief: testutil.Path(t, "testdata/briefs/minimal.json"),
		}},
		Gates: []Gate{{Name: "always-fail", Command: "false", Timeout: 5}},
	}
	code, err := Run(context.Background(), spec, dir, pools, true)
	if err != nil {
		t.Fatal(err)
	}
	if code == 0 {
		t.Fatal("expected non-zero exit on gate failure")
	}
	ledger := filepath.Join(dir, "ledger.jsonl")
	if _, err := os.Stat(ledger); err != nil {
		t.Fatal("ledger not written")
	}
}
