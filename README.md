# muster

*Native orchestration, shared policy, prose handoffs.*

**muster** makes its orchestrator the default for Claude Code, Codex, Cursor Agent, and Antigravity (`agy`). Each harness uses its native workers, worktrees, and lifecycle events. Muster installs prose policy; orchestrator assignments and worker handoffs stay ordinary conversation with no required format.

## Commands

- **`muster install`**: copy the kit to `MUSTER_HOME` (orchestrator + worker policy).
- **`muster link`**: wire harness entrypoints so the main session acts as **orchestrator** (plan, delegate, verify) and **worker** subagents stay scoped.
- **`muster doctor`**: kit on disk and link manifest. `--probe` checks harness binary versions; the opt-in `--live` check runs an isolated delegation task against Codex.

The orchestrator keeps coordination small, delegates bounded work, parallelizes independent tasks, and integrates and verifies worker changes. Use `gh-stacking` for dependent pull requests. After required checks are green, add the exact GitHub comment `@mergify queue` and wait for Mergify, merge, and post-merge CI events. Do not poll or sleep while waiting.

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

`muster doctor --live` uses the signed-in Codex account's model usage to verify a real orchestrator-to-worker handoff. Use `--live-timeout` to set its time limit and `--codex-bin` to select the exact Codex executable, such as the desktop app's binary. The check is read-only and runs in an isolated task.

## Validation

```bash
just check
```

## Release

Tag `v*` on `main` runs GoReleaser and updates the Homebrew formula. Preflight: `just release-dry`.

## License

MIT. See [LICENSE](LICENSE).
