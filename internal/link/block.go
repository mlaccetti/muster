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

For substantial work, use Codex's native plan or task-tracking controls when the active session exposes them (for example, Plan mode or its task tracker); do not simulate them with TODO files. Use /goal only for an explicitly multi-turn objective with a verifiable finish line, and only when the user has asked for continued work. A Goal does not replace worker decomposition or evidence checks. Ask the user directly only when a material decision blocks progress; Codex has no portable AskUserQuestion tool. Delegate through Codex's native subagent/worker controls when available, and use lifecycle waits or callbacks exposed by the active session instead of polling or sleeping.
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
