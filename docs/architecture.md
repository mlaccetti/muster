# Architecture

muster is a **policy distribution and link** tool for interactive coding-agent harnesses.

## Components

| Piece | Role |
|-------|------|
| Embedded kit (`kit/policy/`, `kit/link/`) | Prose orchestration and worker behavior for harness sessions |
| `muster install` | Copy kit to `$MUSTER_HOME/kit/<version>/`, symlink `current` |
| `muster link` | Symlink harness entrypoints and set default orchestrator selectors |
| `muster doctor` | Kit path, link manifest, optional harness CLI probe |

The native harness handles task assignment, worker handoffs, concurrency, worktrees, and event waits. Muster does not define a task/result schema or run worker processes. Codex workers use Luna Medium; other harnesses use their native model settings.

The orchestrator should delegate bounded tasks, parallelize independent work, use `gh-stacking` for dependent pull requests, and wait on lifecycle events or callbacks rather than polling. When checks are green, it adds the exact `@mergify queue` GitHub comment and verifies the subsequent Mergify, merge, and post-merge CI events.

## Related

- [linking.md](linking.md)
