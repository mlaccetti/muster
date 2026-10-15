# Cursor Agent projection

Prefer **`muster link`**: `~/.cursor/rules/muster-orchestrator.mdc` with `alwaysApply: true`.

```bash
export KIT="$(muster kit path)"
test -f ~/.cursor/rules/muster-orchestrator.mdc
# Orchestrator policy: $KIT/policy/orchestrator.md
```

Binary resolved as `cursor-agent`. Delegate bounded work to native worker subagents per kit policy. Assignments and handoffs are ordinary conversation.

