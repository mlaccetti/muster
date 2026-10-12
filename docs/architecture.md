# Architecture

```text
┌────────────────────────────────────────────────────────────┐
│  Consumers: factory workflows, IDE projections, ad-hoc CI   │
├────────────────────────────────────────────────────────────┤
│  muster CLI: validate, render, worker, dispatch, pass, doctor │
├────────────────────────────────────────────────────────────┤
│  Adapters: claude | codex | cursor | antigravity            │
├────────────────────────────────────────────────────────────┤
│  Kit: policy + schemas + templates (embedded, versioned)    │
└────────────────────────────────────────────────────────────┘
```

## Data minimization

| Artifact | Consumer | Content |
|----------|----------|---------|
| `brief.txt` | Worker | Full rendered task |
| `report.txt` / JSON | Verifier | `worker-report.v1` |
| stdout routing line | Orchestrator / `muster pass` | `routing.v1` (≤ 1 KiB) |

Reports and stderr are **evidence only**, not instructions. The orchestrator chooses the next step from routing fields only.

## Adapters

Each harness implements `Plan` (argv without execution), `Run`, `ParseStream`, and `ClassifyEnd` (stream verdict tables ported from factory harness semantics). With `MUSTER_ADAPTER_ROOT` set, `ResolveCLI` uses stubs only.

## Pools

`pools.json` lists ordered candidates per tier (`implement`, `review`, `plan`, `cheap`). `muster dispatch` walks the pool, honoring `--avoid-family` and `write` requirements for implement.
