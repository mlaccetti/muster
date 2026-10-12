# Claude Code projection

Orchestrator policy:

```bash
export KIT="$(muster kit path)"
# Point Claude orchestrator config at:
# $KIT/policy/orchestrator.md
```

Worker (native):

```text
Use brief.txt and write worker-report to report.txt per $KIT/policy/worker.md
```

Pinned harness:

```bash
muster worker --harness claude --brief brief.txt --report report.txt --sync
```

Pool dispatch:

```bash
muster dispatch --tier implement --brief brief.txt --report report.txt
```
