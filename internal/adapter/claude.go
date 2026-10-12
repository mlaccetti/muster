package adapter

import (
	"context"
	"encoding/json"
	"fmt"
)

type claudeAdapter struct{}

func (c *claudeAdapter) ID() string { return "claude" }

func (c *claudeAdapter) ResolveCLI() (string, string, error) {
	bin, err := ResolveBinary("claude")
	if err != nil {
		return "", "", err
	}
	ver, err := probeVersion(bin)
	if err != nil {
		return bin, "", err
	}
	return bin, ver, nil
}

func (c *claudeAdapter) Plan(req PlanRequest) (Invocation, error) {
	bin, err := ResolveBinary("claude")
	if err != nil {
		return Invocation{}, err
	}
	prompt := fmt.Sprintf("Follow the brief at %s. Write your final report to %s when you are done.", req.BriefPath, req.ReportPath)
	args := []string{"-p", prompt}
	if req.Model != "" {
		args = append(args, "--model", req.Model)
	}
	return Invocation{Command: bin, Args: args, WorkDir: req.WorkDir}, nil
}

func (c *claudeAdapter) Run(ctx context.Context, inv Invocation, opts RunOptions) (RunResult, error) {
	return RunInvocation(ctx, c, inv, opts)
}

func (c *claudeAdapter) ParseStream(record []byte) (NormalizedEvent, error) {
	var ev NormalizedEvent
	if err := json.Unmarshal(record, &ev); err != nil {
		return nil, err
	}
	return ev, nil
}

func (c *claudeAdapter) ClassifyEnd(events []NormalizedEvent, exitCode int) EndVerdict {
	return classifyClaude(events, exitCode)
}

func (c *claudeAdapter) Capabilities() Capabilities {
	return Capabilities{Write: true, ReadOnly: true, UsageLimitSignal: true, Resume: true}
}
