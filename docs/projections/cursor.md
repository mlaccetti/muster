# Cursor Agent projection

```bash
export KIT="$(muster kit path)"
# Orchestrator: $KIT/policy/orchestrator.md
```

Binary resolved as `cursor-agent`.

```bash
muster worker --harness cursor --brief brief.txt --report report.txt --sync
muster dispatch --tier plan --brief brief.txt --report report.txt
```
