# Agent prompt: build **muster** (standalone policy + link)

Use this document as the **single source of truth** for **muster**: a public Go CLI and installable **kit** (orchestrator + worker policy, harness link entrypoints) that gives the same orchestration behavior across coding agents **without** tying sessions to one vendor.

**muster works standalone.** Operators need only `muster install` and `muster link`. **rckt-factory is optional:** when `rckt` is on `PATH`, orchestrator policy may use `rckt run …` for durable delivery; muster must not require factory to install or run.

**This repo is not [harness-bridge](https://github.com/mlaccetti/harness-bridge).** harness-bridge moves conversation context between CLIs. Link it in README as an optional related project.

**This repo is not rckt-factory.** Do not implement factory backlog, charter, pools, harness subprocess dispatch, brief rendering, or pass runners here.

**Upstream:** `github.com/mlaccetti/muster` (public).  
**Homebrew:** `github.com/mlaccetti/homebrew-tap` (formula `muster`).  
**License:** MIT.

---

## Name and positioning

**Project name:** `muster`  
**CLI binary:** `muster`

**Tagline (README):** *Same orchestrator and worker behavior in every coding harness.*

---

## Mission

Deliver in one installable unit:

1. **Kit (data):** Versioned `kit/policy/orchestrator.md`, `kit/policy/worker.md`, and `kit/link/*` harness entrypoints (embedded, versioned).
2. **CLI (code):** `muster install`, `muster link`, `muster doctor`, `muster kit path`. No worker, dispatch, or pass subcommands in this repo.
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

## Non-goals

- Do **not** require **rckt-factory**, **rckt-skills**, or harness-bridge.
- Do **not** ship `muster worker`, `muster dispatch`, `muster pass`, adapters, `pools.json`, or JSON Schema brief/routing pipelines in this repo (defer to factory or future work).
- Do **not** call real `claude` / `codex` / `cursor-agent` / `agy` in CI or default tests.
- Do **not** use `time.Sleep` in production code; gate with `scripts/check-no-sleep.sh`.

---

## Architecture

```text
┌────────────────────────────────────────────────────────────┐
│  Optional: rckt-factory (rckt run) when rckt on PATH      │
├────────────────────────────────────────────────────────────┤
│  muster CLI: install, link, doctor, kit                     │
├────────────────────────────────────────────────────────────┤
│  Kit: policy + kit/link harness entrypoints (embedded)      │
└────────────────────────────────────────────────────────────┘
```

Orchestrator policy: plan units, delegate to **worker** subagents in the harness, verify evidence. Optionally mention `rckt run` when factory is installed. Never treat factory as a hard dependency.

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

Root `Long` help must state standalone operation and optional factory.

---

## Linking

- Auto-detect harnesses on `PATH`; flags `--harness`, `--all`, `--global`, `--project`, `--unlink`.
- Record state in `$MUSTER_HOME/link-manifest.json`.
- Codex: managed TOML block + profiles; Claude/Agy: settings JSON agent selector; Cursor: always-applied rules.
- Optional shims: `codex-muster`, `agy-muster` under `$MUSTER_HOME/bin`.

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

1. Standalone mission; optional composition with factory, skills, harness-bridge.
2. Install: Homebrew, `go install`, `muster install`, `muster link`.
3. Projections section.
4. Environment table, `just check`, MIT license.

---

## Acceptance criteria

- [ ] `just check` green without real harness binaries.
- [ ] `muster install` + `muster link` + `muster doctor` work in tests.
- [ ] Docs and kit policy never imply factory is required.
- [ ] Orchestrator link prompts delegate to workers first; `rckt run` only when installed.

---

## Agent invocation

```text
Implement or maintain **muster** per PROMPT-build.md.

Module: github.com/mlaccetti/muster. CLI: Cobra (+ Fang/Carapace as in repo).
Hard constraints: standalone muster; optional rckt-factory; no worker/dispatch/pass in this repo; no real harness binaries in CI.
```
