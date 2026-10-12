# Codex projection

```bash
export KIT="$(muster kit path)"
# Orchestrator: $KIT/policy/orchestrator.md
```

Worker: follow `$KIT/policy/worker.md` with `brief.txt` / `report.txt`.

```bash
muster worker --harness codex --brief brief.txt --report report.txt --sync
muster dispatch --tier review --brief brief.txt --report report.txt
```
