# muster

*Same orchestrator and worker behavior in every coding harness.*

**muster** installs a small, versioned **policy kit** and **links** it into Claude Code, Codex, Cursor Agent, and Antigravity (`agy`). Orchestrators plan work and delegate to **worker** subagents inside the harness.

## Commands

- **`muster install`**: copy the kit to `MUSTER_HOME` (orchestrator + worker policy).
- **`muster link`**: wire harness entrypoints so the main session acts as **orchestrator** (plan, delegate, verify) and **worker** subagents stay scoped.
- **`muster doctor`**: kit on disk, link manifest, optional harness CLI probe.

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
