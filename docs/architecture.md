# Architecture

muster is a **policy distribution and link** tool for interactive coding-agent harnesses.

## Components

| Piece | Role |
|-------|------|
| Embedded kit (`kit/policy/`, `kit/link/`) | Shared orchestration policy, concise harness adapters, and native prompt templates |
| `muster install` | Copy kit to `$MUSTER_HOME/kit/<version>/`, symlink `current` |
| `muster link` | Compose shared policy with each adapter, render native prompts, and set default orchestrator selectors |
| `muster doctor` | Kit path, link manifest, optional harness CLI probe |

The native harness handles task assignment, worker handoffs, concurrency, worktrees, and event waits. Muster does not define a task/result schema or run worker processes. Codex workers use Luna Medium; other harnesses use their native model settings.

The orchestrator should delegate bounded tasks, parallelize independent work, use `gh-stacking` for dependent pull requests, and wait on lifecycle events or callbacks rather than polling. After required checks and reviews pass, it explicitly enters the repository's configured merge queue and verifies the queue response, actual merge, and post-merge CI on the exact merge commit. For `rckt-corp`, this means GitHub's native `enqueuePullRequest` with `expectedHeadOid` matching the reviewed PR head, and `dequeuePullRequest` for removal. A missing queue or rejected enqueue is a blocker, with no direct-merge fallback. Queue labels, statuses, comments, and accepted enqueue requests are not evidence that a merge occurred.

## Related

- [linking.md](linking.md)
