# Architecture

muster is a **standalone policy distribution and link** tool. It is not a dispatcher and does not require **rckt-factory**.

## Components

| Piece | Role |
|-------|------|
| Embedded kit (`kit/policy/`, `kit/link/`) | Orchestrator and worker behavior for interactive harness sessions |
| `muster install` | Copy kit to `$MUSTER_HOME/kit/<version>/`, symlink `current` |
| `muster link` | Symlink harness entrypoints and set default orchestrator selectors |
| `muster doctor` | Kit path, link manifest, optional harness CLI probe |

## Optional composition (rckt-factory)

muster works without factory. When an operator installs **[rckt-factory](https://github.com/mlaccetti/rckt-factory)** and `rckt` is on `PATH`, orchestrator policy may delegate durable delivery to `rckt run …` instead of only harness worker subagents.

Factory owns harness adapters, tier pools, quota/capacity admission, brief rendering, worker subprocess dispatch, pass runners, and delivery graphs (`rckt`, `rcktd`). None of that ships in muster.

## Related

- [linking.md](linking.md)
- [harness-bridge](https://github.com/mlaccetti/harness-bridge) for session handoff between CLIs (optional)
