package adapter

import (
	"context"
	"encoding/json"
	"fmt"
)

type codexAdapter struct{}

func (c *codexAdapter) ID() string { return "codex" }

func (c *codexAdapter) ResolveCLI() (string, string, error) {
	bin, err := ResolveBinary("codex")
	if err != nil {
		return "", "", err
	}
	ver, err := probeVersion(bin)
	if err != nil {
		return bin, "", err
	}
	return bin, ver, nil
}

func (c *codexAdapter) Plan(req PlanRequest) (Invocation, error) {
	bin, err := ResolveBinary("codex")
	if err != nil {
		return Invocation{}, err
	}
	prompt := fmt.Sprintf("Follow the brief at %s. Write your final report to %s when you are done.", req.BriefPath, req.ReportPath)
	args := []string{"exec", prompt}
	if req.Model != "" {
		args = append(args, "--model", req.Model)
	}
	return Invocation{Command: bin, Args: args, WorkDir: req.WorkDir}, nil
}

func (c *codexAdapter) Run(ctx context.Context, inv Invocation, opts RunOptions) (RunResult, error) {
	return RunInvocation(ctx, c, inv, opts)
}

func (c *codexAdapter) ParseStream(record []byte) (NormalizedEvent, error) {
	var ev NormalizedEvent
	if err := json.Unmarshal(record, &ev); err != nil {
		return nil, err
	}
	return ev, nil
}

func (c *codexAdapter) ClassifyEnd(events []NormalizedEvent, exitCode int) EndVerdict {
	return classifyCodex(events, exitCode)
}

func (c *codexAdapter) Capabilities() Capabilities {
	return Capabilities{Write: true, ReadOnly: true, NativeSandbox: true, UsageLimitSignal: true}
}
