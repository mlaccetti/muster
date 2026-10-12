# Architecture

muster is a **policy distribution and link** tool. It is not a dispatcher.

## Components

| Piece | Role |
|-------|------|
| Embedded kit (`kit/policy/`) | Orchestrator and worker behavior for interactive harness sessions |
| `muster install` | Copy kit to `$MUSTER_HOME/kit/<version>/`, symlink `current` |
| `muster link` | Inject managed blocks into Claude, Codex, Cursor, Antigravity configs |
| `muster doctor` | Kit path, link manifest, optional harness CLI probe |

## Out of scope (rckt-factory)

Harness adapters, tier pools, quota/capacity admission, brief rendering, worker subprocess dispatch, pass runners, and delivery graphs live in **rckt-factory** (`rckt`, `rcktd`). See factory design docs `internal/harness/pool` and delivery pipeline sections.

## Related

- [linking.md](linking.md)
- [harness-bridge](https://github.com/mlaccetti/harness-bridge) for session handoff between CLIs (optional)
