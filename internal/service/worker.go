package service

import (
	"context"
	"os"
	"path/filepath"

	"github.com/mlaccetti/muster/internal/adapter"
	"github.com/mlaccetti/muster/internal/pools"
	"github.com/mlaccetti/muster/internal/routing"
)

type WorkerRequest struct {
	Harness    string
	BriefPath  string
	ReportPath string
	WorkDir    string
	Model      string
	Write      bool
	Sync       bool
}

type WorkerResult struct {
	Line     routing.Line
	ExitCode int
}

func RunWorker(ctx context.Context, req WorkerRequest) (WorkerResult, error) {
	a, err := adapter.Get(req.Harness)
	if err != nil {
		return WorkerResult{}, err
	}
	inv, err := a.Plan(adapter.PlanRequest{
		BriefPath:  req.BriefPath,
		ReportPath: req.ReportPath,
		WorkDir:    req.WorkDir,
		Model:      req.Model,
		Write:      req.Write,
	})
	if err != nil {
		return WorkerResult{}, err
	}
	res, err := a.Run(ctx, inv, adapter.RunOptions{ReportPath: req.ReportPath, Sync: req.Sync})
	if err != nil {
		return WorkerResult{}, err
	}
	if m, ok := adapter.LastRoutingLine(res.Stdout); ok && os.Getenv("MUSTER_ADAPTER_ROOT") != "" {
		line, err := routingFromMap(m, req.ReportPath)
		if err != nil {
			return WorkerResult{}, err
		}
		if line.Verdict == "ok" {
			abs, _ := filepath.Abs(req.ReportPath)
			if err := adapter.ReportOK(abs); err != nil {
				return WorkerResult{ExitCode: 1, Line: routing.Line{
					Verdict: "error", Harness: routing.Ptr(a.ID()), ReasonCode: "report_missing", Report: "", Reset: 0,
				}}, nil
			}
			line.Report = abs
		}
		return WorkerResult{Line: line, ExitCode: exitForVerdict(line.Verdict)}, nil
	}
	verdict := a.ClassifyEnd(res.Events, res.ExitCode)
	line := routing.Line{
		Report:     "",
		ReasonCode: "none",
		Reset:      0,
		Harness:    routing.Ptr(a.ID()),
	}
	switch verdict.Kind {
	case "ok":
		abs, _ := filepath.Abs(req.ReportPath)
		if err := adapter.ReportOK(abs); err != nil {
			line.Verdict = "error"
			line.ReasonCode = "report_missing"
			return WorkerResult{Line: line, ExitCode: 1}, nil
		}
		line.Verdict = "ok"
		line.Report = abs
	case "limit":
		line.Verdict = "exhausted"
		line.ReasonCode = "limit"
		line.Reset = verdict.Reset
	default:
		line.Verdict = "error"
		line.ReasonCode = "harness_error"
	}
	return WorkerResult{Line: line, ExitCode: exitForVerdict(line.Verdict)}, nil
}

func routingFromMap(m map[string]any, reportPath string) (routing.Line, error) {
	line := routing.Line{ReasonCode: "none", Reset: 0}
	if v, ok := m["verdict"].(string); ok {
		line.Verdict = v
	}
	if v, ok := m["harness"].(string); ok {
		line.Harness = routing.Ptr(v)
	}
	if v, ok := m["family"].(string); ok {
		line.Family = routing.Ptr(v)
	}
	if v, ok := m["model"].(string); ok {
		line.Model = routing.Ptr(v)
	}
	if v, ok := m["reason_code"].(string); ok {
		line.ReasonCode = v
	}
	if v, ok := m["reset"].(float64); ok {
		line.Reset = int(v)
	}
	if line.Verdict == "ok" {
		abs, _ := filepath.Abs(reportPath)
		line.Report = abs
	}
	return line, nil
}

func exitForVerdict(v string) int {
	switch v {
	case "ok":
		return 0
	case "exhausted", "no_candidate":
		return 3
	default:
		return 1
	}
}

type DispatchRequest struct {
	Tier          string
	BriefPath     string
	ReportPath    string
	WorkDir       string
	PoolsPath     string
	AvoidFamilies []string
	Sync          bool
}

func RunDispatch(ctx context.Context, req DispatchRequest) (WorkerResult, error) {
	pf, err := pools.Load(req.PoolsPath)
	if err != nil {
		return WorkerResult{}, err
	}
	requireWrite := req.Tier == "implement"
	cands, err := pools.Walk(req.Tier, pf, req.AvoidFamilies, requireWrite)
	if err != nil {
		return WorkerResult{}, err
	}
	if len(cands) == 0 {
		reason := "all_avoided"
		if len(req.AvoidFamilies) == 0 {
			reason = "none_usable"
		}
		return WorkerResult{
			ExitCode: 3,
			Line: routing.Line{
				Verdict: "no_candidate", ReasonCode: reason, Report: "", Reset: 0,
			},
		}, nil
	}
	earliestReset := 0
	for _, c := range cands {
		wr := WorkerRequest{
			Harness: c.Harness, BriefPath: req.BriefPath, ReportPath: req.ReportPath,
			WorkDir: req.WorkDir, Model: c.Model, Write: c.Write, Sync: req.Sync,
		}
		res, err := RunWorker(ctx, wr)
		if err != nil {
			return WorkerResult{}, err
		}
		res.Line.Family = routing.Ptr(c.Family)
		res.Line.Model = routing.Ptr(c.Model)
		if res.Line.Harness == nil {
			res.Line.Harness = routing.Ptr(c.Harness)
		}
		if res.Line.Verdict == "ok" || res.Line.Verdict == "error" {
			return res, nil
		}
		if res.Line.Reset > 0 && (earliestReset == 0 || res.Line.Reset < earliestReset) {
			earliestReset = res.Line.Reset
		}
	}
	ended := true
	return WorkerResult{
		ExitCode: 3,
		Line: routing.Line{
			Verdict: "exhausted", ReasonCode: "limit", Report: "", Reset: earliestReset, Ended: &ended,
		},
	}, nil
}
