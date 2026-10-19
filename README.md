# muster

*Native orchestration, shared policy, prose handoffs.*

**muster** makes its orchestrator the default for Claude Code, Codex, Cursor Agent, and Antigravity (`agy`). Each harness uses its native workers, worktrees, and lifecycle events. Muster installs prose policy; orchestrator assignments and worker handoffs stay ordinary conversation with no required format.

## Commands

- **`muster install`**: copy the kit to `MUSTER_HOME` (orchestrator + worker policy).
- **`muster link`**: wire harness entrypoints so the main session acts as **orchestrator** (plan, delegate, verify) and **worker** subagents stay scoped.
- **`muster doctor`**: kit on disk, link manifest, optional harness CLI probe.

The orchestrator keeps coordination small, delegates bounded work, parallelizes independent tasks, and integrates and verifies worker changes. Use `gh-stacking` for dependent pull requests. After required checks and reviews are green, explicitly enter the repository's configured merge queue. For `rckt-corp`, use native GitHub `enqueuePullRequest` with the exact reviewed `expectedHeadOid` and `dequeuePullRequest` when removing a queued PR. Verify the actual merge and post-merge CI on its exact commit. If the queue is missing or rejects entry, report the blocker; do not fall back to direct merge. Do not poll or sleep while waiting.

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
