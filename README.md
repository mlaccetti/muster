# muster

*Same orchestrator and worker behavior in every coding harness.*

**muster** installs a small, versioned **policy kit** and **links** it into Claude Code, Codex, Cursor Agent, and Antigravity (`agy`). It does not dispatch harnesses, run delivery graphs, or replace **rckt-factory**.

## What muster does

- **`muster install`**: copy the kit to `MUSTER_HOME` (orchestrator + worker policy).
- **`muster link`**: inject managed snippets so the main session acts as **orchestrator** (plan, delegate, verify) and **worker** subagents stay scoped.
- **`muster doctor`**: kit on disk, link manifest, optional harness CLI probe.

## What muster does not do

- No `muster worker`, pools, brief rendering, or pass runner. **rckt-factory** owns harness adapters, tier pools, quota/capacity, and durable runs (`rckt` / `rcktd`).
- No skills or planning content. **rckt-skills** owns brainstorm/plan skills; factory consumes approved plans when you run delivery.

## Related projects

- **[rckt-factory](https://github.com/mlaccetti/rckt-factory)**: execution, harness selection, delivery workflows.
- **[rckt-skills](https://github.com/rckt-corp/rckt-skills)**: planning and decomposition skills (optional).
- **[harness-bridge](https://github.com/mlaccetti/harness-bridge)**: move conversation context between CLIs (optional).

## Install

```bash
brew tap mlaccetti/tap
brew install muster

# or
go install github.com/mlaccetti/muster/cmd/muster@latest

muster install
muster link
```

After `brew upgrade muster`, run `muster install` and `muster link` again.

Developer setup:

```bash
mise trust && mise install
just check
```

## Projections

- [Claude Code](docs/projections/claude.md)
- [Codex](docs/projections/codex.md)
- [Cursor Agent](docs/projections/cursor.md)
- [Antigravity (`agy`)](docs/projections/antigravity.md)

## Environment

| Variable | Purpose |
|----------|---------|
| `MUSTER_HOME` | Config root (default `$XDG_CONFIG_HOME/muster` or `~/.config/muster`) |
| `MUSTER_SKIP_PROBE` | `1` makes `muster doctor` skip harness binary probe (CI) |

## License

MIT. See [LICENSE](LICENSE).
