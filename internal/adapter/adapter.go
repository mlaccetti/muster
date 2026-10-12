package adapter

import (
	"context"
)

type Capabilities struct {
	Write            bool
	ReadOnly         bool
	NativeSandbox    bool
	UsageLimitSignal bool
	Resume           bool
}

type PlanRequest struct {
	BriefPath  string
	ReportPath string
	WorkDir    string
	Model      string
	Write      bool
}

type Invocation struct {
	Command string
	Args    []string
	WorkDir string
}

type RunOptions struct {
	ReportPath string
	Sync       bool
}

type NormalizedEvent map[string]any

type EndVerdict struct {
	Kind   string // ok, limit, error
	Reset  int
	Reason string
}

type RunResult struct {
	Events   []NormalizedEvent
	ExitCode int
	Stdout   []byte
}

type Adapter interface {
	ID() string
	ResolveCLI() (absPath string, versionLine string, err error)
	Plan(req PlanRequest) (Invocation, error)
	Run(ctx context.Context, inv Invocation, opts RunOptions) (RunResult, error)
	ParseStream(record []byte) (NormalizedEvent, error)
	ClassifyEnd(events []NormalizedEvent, exitCode int) EndVerdict
	Capabilities() Capabilities
}
