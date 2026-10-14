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
	work := kitDir + "/policy/worker.md"
	return fmt.Sprintf(`%s
# Muster orchestrator (managed by muster link)

You orchestrate work using the **muster** kit. Before planning or delegating units:

1. Read and follow **%s** for orchestration rules.
2. Delegate scoped units to **worker** subagents. Do not implement unit work in this thread.
3. Workers follow **%s** and return the structured handback defined there.
%s
`, markerStart, orch, work, markerEnd)
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
