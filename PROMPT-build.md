# Agent prompt: build **muster** (harness-neutral orchestrator kit)

Use this document as the **single source of truth** for implementing **muster**: a public Go CLI and installable **kit** (policy, schemas, templates) that implements the **orchestrator + worker** pattern for coding agents **without** tying orchestration to any one vendor.

**This repo is not [harness-bridge](https://github.com/mlaccetti/harness-bridge).** harness-bridge moves **conversation context** between CLIs. muster defines **how work is briefed, executed, and reported back** with minimal cross-boundary data. They complement each other; do not merge codebases.

**This repo is not rckt-factory.** Factory workflows, charter, backlog, and PR mechanics are **consumers** of muster contracts. Keep factory-specific logic out of muster.

**Upstream (new):** `github.com/mlaccetti/muster` (public).  
**Homebrew:** `github.com/mlaccetti/homebrew-tap` (formula `muster`).  
**License:** MIT.

---

## Name and positioning

**Project name:** `muster`  
**CLI binary:** `muster` (no required short alias; optional `mu` only if trivial to register in Cobra without colliding with musl packagers).

**Why this name:** To *muster* is to assemble people or resources for a task. The orchestrator musters workers; each harness runs its own orchestrator loop using the same kit. The name is short, distinctive, and not overloaded in the agent-CLI space (unlike "dispatch", "harness", or "agent").

**Tagline (README):** *Harness-neutral briefs, workers, and routing for coding agents.*

**Alternatives if the GitHub org/name is taken:** `billet` (work-order), `spindle` (fan-out). Do not ship under a name that collides with an existing popular CLI on Homebrew.

---

## Mission

Deliver three things in one installable unit:

1. **Kit (data):** Versioned policy (`orchestrator.md`, `worker.md`), JSON Schemas (`brief`, `worker-report`, `routing`), and text templates to render briefs from structured input.
2. **CLI (code):** Go tool to validate, render, run workers through **adapters**, optionally walk a **pool** (multi-harness failover), and run a **deterministic pass driver** that never calls an LLM.
3. **Projections (docs):** Short, copy-paste snippets for Cursor, Claude Code, Codex, and Antigravity that point at the installed kit path. No vendor is the required orchestrator host.

**Supported harness IDs** (closed set, same IDs everywhere):

| ID | User-facing | Real binary (resolved at runtime) |
|----|-------------|-----------------------------------|
| `claude` | Claude Code | `claude` |
| `codex` | Codex CLI | `codex` |
| `cursor` | Cursor Agent | `cursor-agent` |
| `antigravity` | Antigravity | `agy` (accept `agy` as synonym in flags) |

**Default install location:**

| Variable | Default |
|----------|---------|
| `MUSTER_HOME` | `$XDG_CONFIG_HOME/muster` if set, else `~/.config/muster` |
| Kit files | `$MUSTER_HOME/kit/<kit-version>/` |
| User pools | `$MUSTER_HOME/pools.json` (optional; overrides embedded default pools) |
| Pass state | `$MUSTER_HOME/passes/<pass-id>/` (when using `muster pass`) |

`muster install` lays down kit + default `pools.json` and prints projection one-liners.

---

## Non-goals

- Do **not** implement session bridging (no parsing/writing foreign session JSONL). Link to harness-bridge in README.
- Do **not** implement factory backlog, charter, ledger, or GitHub PR automation.
- Do **not** require Claude (or any harness) as the orchestrator. Pool dispatch may run Codex from a shell script; that is not "Claude orchestrates Codex" as a product requirement.
- Do **not** call real LLM APIs or real `claude` / `codex` / `cursor-agent` / `agy` binaries in **CI or default tests**.
- Do **not** use `time.Sleep`, `time.After`, busy-wait timing loops, or "poll every N seconds" in tests or in `muster pass` waiting logic. Synchronization must be **event-driven** (process exit, closed pipe, file appeared, context cancel) or **synchronous** (`--sync` / injected `Runner` in tests).

---

## Architecture

### Layers

```text
┌────────────────────────────────────────────────────────────┐
│  Consumers: factory workflows, IDE projections, ad-hoc CI   │
├────────────────────────────────────────────────────────────┤
│  muster CLI: validate, render, worker, dispatch, pass, doctor │
├────────────────────────────────────────────────────────────┤
│  Adapters: claude | codex | cursor | antigravity            │
│  (Plan argv, parse stream, ClassifyEnd, Capabilities)       │
├────────────────────────────────────────────────────────────┤
│  Kit: policy + schemas + templates (versioned, embedded)    │
└────────────────────────────────────────────────────────────┘
```

### Data minimization (required)

Every worker attempt produces:

| Artifact | Consumer | Content |
|----------|----------|---------|
| `brief.txt` | Worker only | Full task (rendered from template + JSON) |
| `report.txt` | Verifier / gate ops | Worker report (`worker-report.v1`) |
| **stdout routing line** | Orchestrator or `muster pass` | Single JSON object (`routing.v1`), ≤ 1 KiB |

The orchestrator (human, IDE agent, or `muster pass`) **must not** treat `report.txt` or CLI stderr as instructions. Those files are untrusted evidence only. Routing fields are listed in `routing.v1.json`; nothing else may be required to choose the next step.

### Orchestrator vs worker (policy)

- **`kit/policy/orchestrator.md`:** Plan units, render briefs, dispatch workers, verify via independent check, track retries. **No** vendor-specific tool names in the normative rules; an appendix maps tools per harness.
- **`kit/policy/worker.md`:** One scoped unit; fixed report block; no sub-delegation; evidence required.

### Pools (optional multi-harness)

`pools.json` mirrors the semantics of rckt-factory tiers: `implement`, `review`, `plan`, `cheap`, ordered candidates `{harness, model, family, write, boundary?}`. `muster dispatch` walks the pool like `dispatch-worker.sh`. `muster worker --harness codex` pins a harness (default local dev).

`--avoid-family` is repeatable. `distinct_from` is a factory concern; muster only implements the flag.

### Deterministic pass driver

`muster pass` runs a **declarative pass file** (YAML or JSON, you choose one; document in `docs/pass-spec.md`) that:

- Lists units with `brief` input (JSON or path)
- Invokes `worker` or `dispatch` per step
- Appends routing lines to `ledger.jsonl`
- Runs **exec gates** (shell commands with timeout via `context`, no sleep-based polling) that read report paths

No LLM nodes in v1. This proves the kit + CLI for headless automation and factory integration later.

---

## JSON Schemas (kit/schemas/)

Author with JSON Schema draft 2020-12; validate with `github.com/santhosh-tekuri/jsonschema/v6` (same family as rckt-factory).

### `brief.v1.json`

Required fields:

- `goal` (string)
- `modifiable_paths` (array of glob strings)
- `done_when` (string)
- `verify_command` (string)
- `out_of_scope` (string or array)
- `worktree` (object: `path`, `branch` nullable for "none, leave uncommitted")
- Optional: `read_first` (array of strings), `tier` (enum), `harness_hints` (object, non-normative)

### `worker-report.v1.json`

Matches the markdown block in `worker.md` but as structured JSON for machine gates:

- `verdict`: `DONE` | `PARTIAL` | `BLOCKED`
- `changed`, `commands`, `evidence`, `assumptions`, `not_done`, `observations`

CLI accepts worker output as **either** the markdown block **or** JSON file; `muster validate report` normalizes.

### `routing.v1.json`

Required on successful CLI invocation:

- `verdict`: `ok` | `exhausted` | `error` | `no_candidate`
- `harness`, `family`, `model` (nullable strings)
- `report` (absolute path string when `verdict=ok`)
- `reason_code` (enum, when error)
- `reset` (int epoch or 0)
- Optional: `dispatch_id`, `baseline_sha256`, `ended` (bool)

Forbidden in routing: report body, stderr text, raw stream events.

### `pools.v1.json`

Document in `docs/pools.md`; validate on load.

---

## CLI framework (required)

Follow the same stack as harness-bridge ([PROMPT-go-conversion.md](../harness-bridge/PROMPT-go-conversion.md)):

1. **Cobra** command tree in `internal/cli/`
2. **Charm Fang** (`charm.land/fang/v2`) for execution UX, `--version`, completions
3. **Carapace** for flag completion (`--harness`, `--tier`, harness IDs)

`cmd/muster/main.go` only calls `fang.Execute`.

Exit codes: `0` success, `1` runtime error, `2` usage, `3` exhausted/no candidate (match factory dispatch semantics for scripting).

---

## Subcommands (v1)

| Command | Purpose |
|---------|---------|
| `muster install` | Write kit to `MUSTER_HOME`, default pools, print projections |
| `muster doctor` | Local checks: kit present, schemas compile, stub adapters in test mode; optional `--probe` tries `ResolveCLI` without invoking LLM (argv `--version` or `-h` only) |
| `muster kit path` | Print active kit directory |
| `muster validate brief\|report\|routing` | Schema validation; exit non-zero on failure |
| `muster render brief` | `--template implement` `--input unit.json` `--out brief.txt` |
| `muster worker` | Run one harness: `--harness`, `--brief`, `--report`, `--workdir`, `--tier` (informational), `--sync` (foreground; default true in tests) |
| `muster dispatch` | Pool walk: `--tier`, `--brief`, `--report`, `--workdir`, `--avoid-family` (repeatable), `--pools` (path) |
| `muster pass run` | Run pass spec; `--file`, `--pass-dir` |
| `muster pass init` | Scaffold pass dir + example units |
| `completion` | Fang/Cobra |

**Environment for tests only:**

- `MUSTER_ADAPTER_ROOT`: directory of stub binaries (see Testing)
- `MUSTER_SKIP_PROBE=1`: doctor skips real binary probe (CI)

---

## Go module layout

```text
muster/
├── go.mod
├── LICENSE
├── README.md
├── PROMPT-build.md          # this file
├── mise.toml
├── justfile
├── .golangci.yml
├── .goreleaser.yaml
├── .github/
│   ├── CODEOWNERS
│   └── workflows/
│       ├── ci.yml
│       └── release.yml
├── cmd/muster/main.go
├── docs/
│   ├── architecture.md
│   ├── pass-spec.md
│   ├── pools.md
│   └── projections/
│       ├── claude.md
│       ├── codex.md
│       ├── cursor.md
│       └── antigravity.md
├── internal/
│   ├── cli/
│   ├── kit/                 # embed FS + install
│   ├── schema/              # validate brief/report/routing/pools
│   ├── render/              # text/template execute
│   ├── pools/               # eligibility + walk order
│   ├── routing/             # parse/emit single-line JSON
│   ├── pass/                # pass runner (no LLM)
│   ├── adapter/
│   │   ├── adapter.go       # interface
│   │   ├── registry.go
│   │   ├── claude.go
│   │   ├── codex.go
│   │   ├── cursor.go
│   │   └── antigravity.go
│   └── paths/               # MUSTER_HOME, XDG
├── kit/                     # source of truth (copied/embed)
│   ├── policy/
│   ├── schemas/
│   └── templates/
└── testdata/
    ├── streams/             # jsonl per harness (hand-written shapes)
    ├── stubs/               # fake CLI scripts used by tests
    ├── briefs/
    ├── reports/
    └── passes/
```

### Adapter interface (required)

```go
type Adapter interface {
    ID() string // claude, codex, cursor, antigravity
    ResolveCLI() (absPath string, versionLine string, err error)
    Plan(req PlanRequest) (Invocation, error)
    Run(ctx context.Context, inv Invocation, opts RunOptions) (RunResult, error)
    ParseStream(record []byte) (NormalizedEvent, error)
    ClassifyEnd(events []NormalizedEvent, exitCode int) EndVerdict
    Capabilities() Capabilities
}
```

- `Plan` builds argv **without executing**. `doctor` and tests assert argv flags ⊆ declared capabilities.
- `Run` executes the process (or stub), writes worker report to `RunOptions.ReportPath`, returns events for classification.
- **No adapter may invoke an LLM in tests:** when `MUSTER_ADAPTER_ROOT` is set, `ResolveCLI` must resolve to stubs under that root only.

`Capabilities` includes: `Write`, `ReadOnly`, `NativeSandbox`, `UsageLimitSignal`, `Resume` (bool fields).

Port **classification tables** from rckt-factory `scripts/factory/testdata/harness/streams/*.jsonl` into Go table-driven tests (`internal/adapter/*_test.go`).

---

## Testing and validation (no LLM, no sleep)

### Principles

1. **Hermetic:** Tests set `MUSTER_ADAPTER_ROOT` to `testdata/stubs/` and `PATH` to `stubdir:system-path-for-git-only`.
2. **Synchronous:** `muster worker` defaults to `--sync` when `MUSTER_TEST=1` or always in adapter `Run` for stubs (immediate exit).
3. **No `time.Sleep` in module `muster`:** Enforce with a CI script:

   ```bash
   ! rg -n 'time\.Sleep|Sleep\(' --glob '*.go' --glob '!*_test.go' .
   ```

   Tests may use `context.WithTimeout` for hung-process safety only if the stub always exits immediately; prefer no timeout in unit tests.

4. **No real harness binaries in CI:** `doctor` full probe runs only in manual/`muster doctor --probe` locally.

### Stub CLI contract (`testdata/stubs/cli`)

One executable; behavior selected by `basename $0` and a control file:

- Writes brief path to `stub.calls` log
- Writes a valid minimal `report.txt` from fixture
- Prints one `routing.v1` JSON line to stdout
- Exit code from fixture name or env `MUSTER_STUB_EXIT`

Hand-written `testdata/streams/*.jsonl` feed **parser/classify** unit tests only (no subprocess).

### Test layers

| Layer | What | How |
|-------|------|-----|
| Unit | schema, render, pools walk, routing parse | `go test` |
| Adapter | ParseStream + ClassifyEnd | jsonl fixtures |
| Integration | `muster worker`, `muster dispatch`, `muster pass run` | `os/exec` stub root, temp dir, golden stdout |
| Install | `muster install` layout | golden tree compare |

Add `scripts/check-no-sleep.sh` and wire into `just check`.

Functional validation story (document in README **Validation**):

- `just check` is the only required gate; it proves dispatch order, avoid-family, exhausted handling, and pass ledger append **using stubs**.

---

## Tooling: mise + just

Same as harness-bridge:

- `go` 1.22+
- `golangci-lint`, `gofumpt`, `goreleaser`, `just`

```just
default:
    just check

fmt:
    gofumpt -w .

lint:
    golangci-lint run ./...

test:
    go test ./... -race -count=1

check-no-sleep:
    ./scripts/check-no-sleep.sh

build:
    go build -o dist/muster ./cmd/muster

install:
    go install ./cmd/muster

check: fmt lint check-no-sleep test

release-dry:
    goreleaser release --snapshot --clean
```

---

## CI (`.github/workflows/ci.yml`)

Copy harness-bridge pattern:

- Triggers: `pull_request`, `push` to `main`
- Job name: `check` (required for ruleset)
- `jdx/mise-action@v2` → `just check`
- `permissions: contents: read`
- concurrency cancel in-progress

No macOS runner required for v1.

---

## Release (`.github/workflows/release.yml` + `.goreleaser.yaml`)

- Tag `v*` → GoReleaser
- `CGO_ENABLED=0`, darwin/linux amd64+arm64
- Archives include `LICENSE`, `README.md`
- Homebrew formula `muster` in `mlaccetti/homebrew-tap`
- Secret `HOMEBREW_TAP_GITHUB_TOKEN`
- `before.hooks`: `go mod tidy`, `go test ./... -race -count=1`

---

## README requirements

1. Mission and relation to harness-bridge and factory (3 bullets).
2. Install: Homebrew, `go install`, `muster install`.
3. Quickstart: render brief → `muster worker --harness cursor --sync` with stubs example.
4. **Projections** section linking `docs/projections/*.md`.
5. Environment variables table.
6. Validation: `just check`, no LLM in CI.
7. MIT + copyright.

---

## Projections (docs only in v1)

Each projection file is ≤ 40 lines:

- Where to point the harness's orchestrator config at `$MUSTER_HOME/kit/.../policy/orchestrator.md`
- How to spawn a worker (native tool) with `brief.txt` / `report.txt`
- Optional: invoke `muster worker` for pinned harness; `muster dispatch` for pool

Do not embed factory charter or backlog references.

---

## Implementation order (for the agent)

1. Scaffold module, Fang+Cobra root, `mise.toml`, `justfile`, CI workflow.
2. Embed `kit/` via `embed.FS`; `muster install` + `muster kit path`.
3. JSON Schema validation commands.
4. `render brief` from templates.
5. Adapter interface + stubs + stream classification tests (port factory jsonl).
6. `muster worker` + `muster dispatch` with routing stdout contract.
7. `muster pass run` with example pass in `testdata/passes/minimal.yaml`.
8. `muster doctor`, README, projections, GoReleaser.
9. Final pass: `just check`, `just release-dry`.

---

## Acceptance criteria (release v0.1.0)

- [ ] `just check` green on Ubuntu CI without network after module download.
- [ ] No `time.Sleep` in non-test `.go` files; script enforces.
- [ ] No test spawns real `claude`, `codex`, `cursor-agent`, or `agy`.
- [ ] `muster install` + `muster validate` + `muster worker` works against stubs with golden routing line.
- [ ] `muster dispatch` reproduces at least: success, exhausted, avoid-family, pool order (from factory test scenarios, stubbed).
- [ ] `muster pass run` appends ledger and exits non-zero on gate failure.
- [ ] Homebrew formula published on tag.
- [ ] Architecture and pass spec committed under `docs/`.

---

## Agent invocation (paste below)

```text
Implement the public Go project **muster** per PROMPT-build.md in this repository.

Module: github.com/mlaccetti/muster. CLI: Cobra + Fang + Carapace.
Follow harness-bridge PROMPT-go-conversion.md for CI, release, mise, just, and Homebrew.

Hard constraints:
- No LLM invocation in tests or CI.
- No time.Sleep in production code; add scripts/check-no-sleep.sh to just check.
- Hermetic adapter stubs under testdata/stubs; MUSTER_ADAPTER_ROOT in tests.
- Keep factory and harness-bridge code out of this repo; link only in README.

Work in implementation order through acceptance criteria for v0.1.0.
```
