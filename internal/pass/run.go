package pass

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/mlaccetti/muster/internal/routing"
	"github.com/mlaccetti/muster/internal/service"
)

func LoadSpec(path string) (*Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var spec Spec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, err
	}
	return &spec, nil
}

func Run(ctx context.Context, spec *Spec, passDir string, poolsPath string, sync bool) (int, error) {
	if err := os.MkdirAll(passDir, 0o755); err != nil {
		return 1, err
	}
	ledgerPath := filepath.Join(passDir, "ledger.jsonl")
	ledger, err := os.OpenFile(ledgerPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return 1, err
	}
	defer func() { _ = ledger.Close() }()

	exitCode := 0
	for _, unit := range spec.Units {
		reportPath := filepath.Join(passDir, unit.Name+"-report.json")
		briefPath := unit.Brief
		var res service.WorkerResult
		if unit.Harness != "" {
			res, err = service.RunWorker(ctx, service.WorkerRequest{
				Harness: unit.Harness, BriefPath: briefPath, ReportPath: reportPath,
				WorkDir: passDir, Sync: sync,
			})
		} else {
			tier := unit.Tier
			if tier == "" {
				tier = "implement"
			}
			res, err = service.RunDispatch(ctx, service.DispatchRequest{
				Tier: tier, BriefPath: briefPath, ReportPath: reportPath, WorkDir: passDir,
				PoolsPath: poolsPath, Sync: sync,
			})
		}
		if err != nil {
			return 1, err
		}
		if err := appendLedger(ledger, res.Line); err != nil {
			return 1, err
		}
		if res.ExitCode != 0 {
			exitCode = res.ExitCode
		}
	}

	for _, gate := range spec.Gates {
		code, err := runGate(ctx, gate)
		if err != nil {
			return 1, err
		}
		if code != 0 {
			return code, nil
		}
	}
	return exitCode, nil
}

func appendLedger(f *os.File, line routing.Line) error {
	b, err := json.Marshal(line)
	if err != nil {
		return err
	}
	_, err = f.Write(append(b, '\n'))
	return err
}

func runGate(ctx context.Context, gate Gate) (int, error) {
	timeout := 60
	if gate.Timeout > 0 {
		timeout = gate.Timeout
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", "-c", gate.Command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return 1, fmt.Errorf("gate %s timed out", gate.Name)
		}
		return 1, nil
	}
	return 0, nil
}

func Init(passDir string) error {
	if err := os.MkdirAll(passDir, 0o755); err != nil {
		return err
	}
	example := filepath.Join(passDir, "pass.yaml")
	if _, err := os.Stat(example); err == nil {
		return nil
	}
	spec := `id: example-pass
units:
  - name: unit1
    brief: ../testdata/briefs/minimal.json
    harness: claude
gates:
  - name: verify-report
    command: test -s unit1-report.json
    timeout_seconds: 30
`
	return os.WriteFile(example, []byte(spec), 0o644)
}
