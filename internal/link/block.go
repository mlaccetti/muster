package link

import (
	"fmt"
	"strings"
)

const (
	markerStart = "<!-- MUSTER_LINK_START -->"
	markerEnd   = "<!-- MUSTER_LINK_END -->"
)

func orchestratorBlock(kitDir string) string {
	orch := kitDir + "/policy/orchestrator.md"
	return fmt.Sprintf(`%s
# Muster orchestrator (managed by muster link)

You are the default orchestrator for this harness. Read and follow **%s**. Use native workers and ordinary conversational assignments and handoffs. Delegate bounded work, parallelize independent tasks, use gh-stacking for dependent pull requests, and add the exact @mergify queue GitHub comment after required checks are green. Wait on lifecycle events or callbacks rather than polling or sleeping, then integrate and verify end to end. Do small tasks directly when delegation adds overhead.

Use Codex's native Goal for an explicitly multi-turn objective with a verifiable finish line, and its native plan/task tracking for substantial work when available. A Goal is not a replacement for worker decomposition or evidence checks. Ask a focused question only when a material decision blocks progress.
%s
`, markerStart, orch, markerEnd)
}

func upsertManagedBlock(existing string, block string) string {
	if strings.Contains(existing, markerStart) && strings.Contains(existing, markerEnd) {
		before, _, _ := strings.Cut(existing, markerStart)
		_, after, _ := strings.Cut(existing, markerEnd)
		after = strings.TrimPrefix(after, "\n")
		out := strings.TrimRight(before, "\n") + "\n\n" + block + "\n"
		if strings.TrimSpace(after) != "" {
			out += after
		}
		return out
	}
	trim := strings.TrimRight(existing, "\n")
	if trim == "" {
		return block + "\n"
	}
	return trim + "\n\n" + block + "\n"
}

func removeManagedBlock(existing string) string {
	if !strings.Contains(existing, markerStart) {
		return existing
	}
	before, _, _ := strings.Cut(existing, markerStart)
	_, after, ok := strings.Cut(existing, markerEnd)
	if !ok {
		return strings.TrimRight(before, "\n") + "\n"
	}
	after = strings.TrimPrefix(after, "\n")
	return strings.TrimRight(before, "\n") + "\n" + after
}
