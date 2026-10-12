package adapter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunInvocation(ctx context.Context, a Adapter, inv Invocation, opts RunOptions) (RunResult, error) {
	cmd := exec.CommandContext(ctx, inv.Command, inv.Args...)
	cmd.Dir = inv.WorkDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return buildResult(a, stdout.Bytes(), ee.ExitCode())
		}
		return RunResult{}, err
	}
	return buildResult(a, stdout.Bytes(), 0)
}

func buildResult(a Adapter, out []byte, code int) (RunResult, error) {
	events := make([]NormalizedEvent, 0)
	sc := bufio.NewScanner(bytes.NewReader(out))
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		ev, err := a.ParseStream(line)
		if err == nil {
			events = append(events, ev)
		}
	}
	return RunResult{Events: events, ExitCode: code, Stdout: out}, nil
}

func LastRoutingLine(out []byte) (map[string]any, bool) {
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(line, "{") {
			continue
		}
		var m map[string]any
		if json.Unmarshal([]byte(line), &m) == nil {
			if v, ok := m["verdict"].(string); ok && v != "" {
				return m, true
			}
		}
	}
	return nil, false
}

func ReportOK(path string) error {
	st, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("report missing: %w", err)
	}
	if st.Size() == 0 {
		return fmt.Errorf("report empty")
	}
	return nil
}
