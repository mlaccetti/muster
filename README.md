# muster

*Same orchestrator and worker behavior in every coding harness.*

**muster** is **standalone**: it installs a small, versioned **policy kit** and **links** it into Claude Code, Codex, Cursor Agent, and Antigravity (`agy`). You do not need **rckt-factory** or any other rckt repo to use it.

Orchestrators plan work and delegate to **worker** subagents inside the harness. When the operator has installed **rckt-factory** and `rckt` is on `PATH`, they may optionally start durable delivery with `rckt run …` instead of re-implementing delivery graphs in chat.

## What muster does

- **`muster install`**: copy the kit to `MUSTER_HOME` (orchestrator + worker policy).
- **`muster link`**: wire harness entrypoints so the main session acts as **orchestrator** (plan, delegate, verify) and **worker** subagents stay scoped.
- **`muster doctor`**: kit on disk, link manifest, optional harness CLI probe.

## What muster does not do

- No harness subprocess dispatch, tier pools, brief rendering, or pass runner in this repo. Those concerns live in **[rckt-factory](https://github.com/mlaccetti/rckt-factory)** when you choose to install it (optional).
- No skills or planning content. **[rckt-skills](https://github.com/rckt-corp/rckt-skills)** is optional for brainstorm/plan skills.

## Related projects (optional composition)

- **[rckt-factory](https://github.com/mlaccetti/rckt-factory)** *(optional)*: durable runs, harness selection, delivery workflows (`rckt` / `rcktd`). Use `rckt run` only when this CLI is installed.
- **[rckt-skills](https://github.com/rckt-corp/rckt-skills)** *(optional)*: planning and decomposition skills.
- **[harness-bridge](https://github.com/mlaccetti/harness-bridge)** *(optional)*: move conversation context between CLIs.

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

## Validation

```bash
just check
```

## Release

Tag `v*` on `main` runs GoReleaser and updates the Homebrew formula. Preflight: `just release-dry`.

## License

MIT. See [LICENSE](LICENSE).
