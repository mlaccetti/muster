# Cursor Agent projection

Prefer **`muster link`**: `~/.cursor/rules/muster-orchestrator.mdc` with `alwaysApply: true`.

```bash
export KIT="$(muster kit path)"
test -f ~/.cursor/rules/muster-orchestrator.mdc
# Orchestrator policy: $KIT/policy/orchestrator.md
```

Binary resolved as `cursor-agent`. Delegate units to worker subagents per kit policy.

**Optional:** when `rckt` is on `PATH`, orchestrator policy allows `rckt run …` (rckt-factory).
