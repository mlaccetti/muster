# Agent prompt: build **muster** (standalone policy + link)

Use this document as the **single source of truth** for **muster**: a public Go CLI and installable **kit** (orchestrator + worker policy, harness link entrypoints) that gives the same orchestration behavior across coding agents **without** tying sessions to one vendor.

Operators need only `muster install` and `muster link`.

**Upstream:** `github.com/mlaccetti/muster` (public).  
**Homebrew:** `github.com/mlaccetti/homebrew-tap` (formula `muster`).  
**License:** MIT.

---

## Name and positioning

**Project name:** `muster`  
**CLI binary:** `muster`

**Tagline (README):** *Native orchestration, shared policy, prose handoffs.*

---

## Mission

Deliver in one installable unit:

1. **Kit (data):** Versioned `kit/policy/orchestrator.md`, `kit/policy/worker.md`, and `kit/link/*` harness entrypoints (embedded, versioned).
2. **CLI (code):** `muster install`, `muster link`, `muster doctor`, `muster kit path`.
3. **Docs:** Architecture, linking, upgrades, and short per-harness projections under `docs/projections/`.

**Supported harness IDs** (closed set):

| ID | Real binary |
|----|-------------|
| `claude` | `claude` |
| `codex` | `codex` |
| `cursor` | `cursor-agent` |
| `antigravity` | `agy` |

**Default install location:**

| Variable | Default |
|----------|---------|
| `MUSTER_HOME` | `$XDG_CONFIG_HOME/muster` or `~/.config/muster` |
| Kit files | `$MUSTER_HOME/kit/<kit-version>/` with `current` symlink |

---

## Non-goals (this repo)

- No `muster worker`, `muster dispatch`, or `muster pass` subcommands.
- No harness subprocess dispatch, tier pools, brief rendering, or pass runners in this codebase.
- No required structured format for orchestrator-to-worker assignments or worker handoffs.
- No calls to real `claude` / `codex` / `cursor-agent` / `agy` in CI or default tests.
- No `time.Sleep` in production code; gate with `scripts/check-no-sleep.sh`.

---

## Architecture

```text
┌────────────────────────────────────────────┐
│  muster CLI: install, link, doctor, kit   │
├────────────────────────────────────────────┤
│  Kit: policy + kit/link harness entrypoints │
└────────────────────────────────────────────┘
```

Orchestrator policy: plan units, delegate to **worker** subagents in the harness, verify evidence.

The orchestrator keeps coordination small, delegates bounded work, parallelizes independent tasks, and uses native lifecycle events or callbacks to wait. Use `gh-stacking` for dependent pull requests. After required checks are green, add the exact `@mergify queue` comment and verify Mergify, merge, and post-merge CI events. Do not poll or sleep.

Worker assignments and handoffs are ordinary harness conversation. Muster must not define a required schema, status enum, or serialized task/result artifact.

See [docs/architecture.md](docs/architecture.md) and [docs/linking.md](docs/linking.md).

---

## CLI surface

| Command | Purpose |
|---------|---------|
| `muster install` | Lay down embedded kit under `MUSTER_HOME` |
| `muster link` | Symlink `kit/link/*`, patch harness defaults, optional shims |
| `muster unlink` | Remove link artifacts |
| `muster doctor` | Kit path, link manifest, `link.VerifyReports`, optional `--probe` |
| `muster kit path` | Print active kit directory |

---

## Linking

- Auto-detect harnesses on `PATH`; flags `--harness`, `--all`, `--global`, `--project`, `--unlink`.
- Record state in `$MUSTER_HOME/link-manifest.json`.
- Codex: native `AGENTS.md` instructions + custom worker agent; Claude/Agy: settings JSON agent selector; Cursor: always-applied project rule (global User Rules remain a Settings-only step).
- Optional Antigravity shim: `agy-muster` under `$MUSTER_HOME/bin`.

Legacy managed markdown blocks (`internal/link/block.go`) remain for upsert helpers/tests; primary delivery is symlinks + kit/link prompts.

---

## Testing

- Unit and integration tests with temp `MUSTER_HOME` and temp `HOME` / `CODEX_HOME`.
- `cmd/muster/e2e_link_test.go` subprocess smoke where needed.
- `just check`: fmt, lint, check-no-sleep, `go test ./... -race`.
- No network in CI after module download.

---

## CI and release

- `.github/workflows/ci.yml`: `just check` on Ubuntu.
- Tag `v*` → GoReleaser, Homebrew formula in `mlaccetti/homebrew-tap`.

---

## README requirements

1. Mission and install (`brew`, `go install`, `muster install`, `muster link`).
2. Projections section.
3. Environment table, `just check`, MIT license.

---

## Acceptance criteria

- [ ] `just check` green without real harness binaries.
- [ ] `muster install` + `muster link` + `muster doctor` work in tests.
- [ ] Orchestrator link prompts delegate to workers; kit policy matches docs.
- [ ] Policies require parallel independent work, event-driven waits, `gh-stacking`, and the exact Mergify queue comment without imposing a handoff schema.

---

## Agent invocation

```text
Implement or maintain **muster** per PROMPT-build.md.

Module: github.com/mlaccetti/muster. CLI: Cobra (+ Fang/Carapace as in repo).
Hard constraints: install/link kit only; no worker/dispatch/pass subcommands; no real harness binaries in CI.
```
