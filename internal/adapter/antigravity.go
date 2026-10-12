package adapter

import (
	"context"
	"encoding/json"
	"fmt"
)

type antigravityAdapter struct{}

func (a *antigravityAdapter) ID() string { return "antigravity" }

func (a *antigravityAdapter) ResolveCLI() (string, string, error) {
	bin, err := ResolveBinary("antigravity")
	if err != nil {
		return "", "", err
	}
	ver, err := probeVersion(bin)
	if err != nil {
		return bin, "", err
	}
	return bin, ver, nil
}

func (a *antigravityAdapter) Plan(req PlanRequest) (Invocation, error) {
	bin, err := ResolveBinary("antigravity")
	if err != nil {
		return Invocation{}, err
	}
	prompt := fmt.Sprintf("Follow the brief at %s. Write your final report to %s when you are done.", req.BriefPath, req.ReportPath)
	args := []string{"chat", prompt}
	return Invocation{Command: bin, Args: args, WorkDir: req.WorkDir}, nil
}

func (a *antigravityAdapter) Run(ctx context.Context, inv Invocation, opts RunOptions) (RunResult, error) {
	return RunInvocation(ctx, a, inv, opts)
}

func (a *antigravityAdapter) ParseStream(record []byte) (NormalizedEvent, error) {
	var ev NormalizedEvent
	if err := json.Unmarshal(record, &ev); err != nil {
		return nil, err
	}
	return ev, nil
}

func (a *antigravityAdapter) ClassifyEnd(events []NormalizedEvent, exitCode int) EndVerdict {
	return classifyClaude(events, exitCode)
}

func (a *antigravityAdapter) Capabilities() Capabilities {
	return Capabilities{Write: true, ReadOnly: true, UsageLimitSignal: true}
}
