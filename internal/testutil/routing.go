package testutil

import (
	"strings"
	"testing"

	"github.com/mlaccetti/muster/internal/routing"
)

func LastRoutingLine(t *testing.T, stdout string) routing.Line {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(line, "{") {
			continue
		}
		parsed, err := routing.ParseLine(line)
		if err != nil {
			t.Fatalf("parse routing: %v", err)
		}
		return parsed
	}
	t.Fatalf("no routing line in stdout:\n%s", stdout)
	return routing.Line{}
}
