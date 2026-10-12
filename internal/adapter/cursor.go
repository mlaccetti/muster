package adapter

import (
	"context"
	"encoding/json"
	"fmt"
)

type cursorAdapter struct{}

func (c *cursorAdapter) ID() string { return "cursor" }

func (c *cursorAdapter) ResolveCLI() (string, string, error) {
	bin, err := ResolveBinary("cursor")
	if err != nil {
		return "", "", err
	}
	ver, err := probeVersion(bin)
	if err != nil {
		return bin, "", err
	}
	return bin, ver, nil
}

func (c *cursorAdapter) Plan(req PlanRequest) (Invocation, error) {
	bin, err := ResolveBinary("cursor")
	if err != nil {
		return Invocation{}, err
	}
	prompt := fmt.Sprintf("Follow the brief at %s. Write your final report to %s when you are done.", req.BriefPath, req.ReportPath)
	args := []string{"-p", prompt}
	return Invocation{Command: bin, Args: args, WorkDir: req.WorkDir}, nil
}

func (c *cursorAdapter) Run(ctx context.Context, inv Invocation, opts RunOptions) (RunResult, error) {
	return RunInvocation(ctx, c, inv, opts)
}

func (c *cursorAdapter) ParseStream(record []byte) (NormalizedEvent, error) {
	var ev NormalizedEvent
	if err := json.Unmarshal(record, &ev); err != nil {
		return nil, err
	}
	return ev, nil
}

func (c *cursorAdapter) ClassifyEnd(events []NormalizedEvent, exitCode int) EndVerdict {
	return classifyCursor(events, exitCode)
}

func (c *cursorAdapter) Capabilities() Capabilities {
	return Capabilities{Write: true, ReadOnly: true, UsageLimitSignal: true}
}
