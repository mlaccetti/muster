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

The orchestrator should delegate bounded tasks, parallelize independent work, use `gh-stacking` for dependent pull requests, and wait on lifecycle events or callbacks rather than polling. Delivery depends on repository ownership and configuration, after GitHub Actions and required reviews pass on the reviewed head:

- Personal repositories use normal protected GitHub merging, pinned to the reviewed head with `--match-head-commit`. An absent queue is expected and does not block this flow.
- `rckt-corp` repositories use GitHub's native `enqueuePullRequest` with `expectedHeadOid` matching the reviewed PR head, and `dequeuePullRequest` for removal. A required queue that is missing or rejects entry is a blocker; direct merge is not its fallback. `rckt-corp/apple-signing` retains its existing signing-specific workflow and queue exemption.
- Other organizations follow the repository's configured delivery flow, including its queue when configured.

Every path preserves branch protections and required reviews, rejects forged checks/statuses, and verifies the actual merge and post-merge CI on the exact merge commit. Queue labels, statuses, comments, and accepted merge or enqueue requests are not evidence that a merge occurred.

## Related

- [linking.md](linking.md)
