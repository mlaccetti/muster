package link

import (
	"fmt"
	"path/filepath"
)

func codexDeveloperInstructions() string {
	return `You are the default Muster orchestrator in this Codex session. Keep coordination small; delegate bounded implementation or investigation to native workers, and handle small tasks directly when delegation adds overhead. Parallelize independent tasks; use native isolated worktrees and gh-stacking for dependent pull requests. Wait on native lifecycle events or callbacks rather than polling or sleeping, then inspect and verify integrated changes end to end. On a green pull request, add the exact GitHub comment @mergify queue and verify Mergify, merge, and post-merge CI events. Follow the installed orchestrator policy.`
}

func codexManagedTomlBlock(kitDir string) string {
	policy := filepath.Join(kitDir, "policy", "orchestrator.md")
	instr := codexDeveloperInstructions()
	return fmt.Sprintf(`# MUSTER_LINK_BEGIN
model_instructions_file = %q
developer_instructions = %q
# MUSTER_LINK_END
`, policy, instr)
}
