# Architecture

muster is a **policy distribution and link** tool for interactive coding-agent harnesses.

## Components

| Piece | Role |
|-------|------|
| Embedded kit (`kit/policy/`, `kit/link/`) | Orchestrator and worker behavior for harness sessions |
| `muster install` | Copy kit to `$MUSTER_HOME/kit/<version>/`, symlink `current` |
| `muster link` | Symlink harness entrypoints and set default orchestrator selectors |
| `muster doctor` | Kit path, link manifest, optional harness CLI probe |

## Related

- [linking.md](linking.md)
