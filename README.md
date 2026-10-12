# muster

*Harness-neutral briefs, workers, and routing for coding agents.*

**muster** ships a versioned **kit** (policy, JSON Schemas, templates) and a Go **CLI** to validate, render briefs, run workers through harness **adapters**, walk optional **pools**, and drive deterministic **passes** without tying orchestration to one vendor.

## Related projects

- **[harness-bridge](https://github.com/mlaccetti/harness-bridge)**: moves conversation context between CLIs (sessions). muster does not parse foreign session JSONL; use harness-bridge for that.
- **[rckt-factory](https://github.com/mlaccetti/rckt-factory)**: factory workflows, charter, and PR mechanics that can **consume** muster contracts. Factory-specific logic stays out of this repo.
- **muster**: how work is briefed, executed, and reported back with minimal cross-boundary data (brief, report, routing line).

## Install

```bash
brew tap mlaccetti/tap
brew install muster

# or
go install github.com/mlaccetti/muster/cmd/muster@latest

muster install   # lay down kit under MUSTER_HOME
muster link      # wire orchestrator + worker into detected harnesses (see docs/linking.md)
```

After `brew upgrade muster`, run `muster install` and `muster link` again. The CLI checks GitHub once per day and nudges you when a newer release exists (see [docs/upgrades.md](docs/upgrades.md)).

Developer setup:

```bash
mise trust && mise install
just check
```

## Quickstart (hermetic stubs)

```bash
export MUSTER_ADAPTER_ROOT="$(pwd)/testdata/stubs/bin"
export MUSTER_STUB_DATA="$(pwd)/testdata/stubs/data"
# copy testdata/stubs/cli to each name under bin/ (see internal/testutil)

muster render brief --template implement --input testdata/briefs/minimal.json --out /tmp/brief.txt
muster worker --harness claude --brief /tmp/brief.txt --report /tmp/report.json --sync
```

The worker prints a single **routing.v1** JSON line on stdout; the report is written to `--report`.

## Projections

IDE-specific one-liners pointing at the installed kit:

- [Claude Code](docs/projections/claude.md)
- [Codex](docs/projections/codex.md)
- [Cursor Agent](docs/projections/cursor.md)
- [Antigravity (`agy`)](docs/projections/antigravity.md)

## Environment

| Variable | Purpose |
|----------|---------|
| `MUSTER_HOME` | Config root (default `$XDG_CONFIG_HOME/muster` or `~/.config/muster`) |
| `MUSTER_ADAPTER_ROOT` | **Tests:** directory of stub harness binaries |
| `MUSTER_STUB_DATA` | **Tests:** per-stub control files (routing, report fixtures) |
| `MUSTER_TEST` | `1` enables synchronous worker defaults in tests |
| `MUSTER_SKIP_PROBE` | `1` makes `muster doctor` skip real binary probe (CI) |

## Validation

CI and local gate:

```bash
just check   # gofumpt, golangci-lint, check-no-sleep, tests (race)
```

No real `claude`, `codex`, `cursor-agent`, or `agy` binaries are invoked in tests or CI. Adapters use hermetic stubs under `testdata/stubs` when `MUSTER_ADAPTER_ROOT` is set.

Test layers:

- **Unit:** schema, render, pools, adapter stream classification, report normalize
- **Service:** worker/dispatch/pass with stubs (`internal/service`, `internal/pass`)
- **CLI in-process:** `internal/cli` via `RunForTest` (Cobra, no Fang)
- **CLI E2E:** `cmd/muster/e2e_test.go` builds the binary and runs subprocesses with `StubEnv` (`PATH` trimmed, stub harness root only)

## Release

Tag `v*` on `main` runs GoReleaser and updates the `muster` formula in [homebrew-tap](https://github.com/mlaccetti/homebrew-tap) (requires `HOMEBREW_TAP_GITHUB_TOKEN`).

Preflight: `just release-dry`

## License

MIT. See [LICENSE](LICENSE). Copyright Michael Laccetti.
