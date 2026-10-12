package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mlaccetti/muster/internal/testutil"
)

func TestRunWorkerStubGolden(t *testing.T) {
	cleanup := testutil.SetupStubs(t)
	defer cleanup()
	dir := t.TempDir()
	report := filepath.Join(dir, "report.json")
	brief := testutil.Path(t, "testdata/briefs/minimal.json")
	res, err := RunWorker(context.Background(), WorkerRequest{
		Harness: "claude", BriefPath: brief, ReportPath: report, WorkDir: dir, Sync: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 0 || res.Line.Verdict != "ok" {
		t.Fatalf("res %+v", res)
	}
	if _, err := os.Stat(report); err != nil {
		t.Fatal(err)
	}
}
