# Antigravity projection

Prefer **`muster link`**: `agy --agent muster-orchestrator` via `~/.gemini/antigravity-cli/settings.json`.

```bash
export KIT="$(muster kit path)"
agy agents | rg muster
jq -r '.agent' ~/.gemini/antigravity-cli/settings.json
# Orchestrator policy: $KIT/policy/orchestrator.md
```

CLI binary `agy` (harness id `antigravity`). Delegate units to worker subagents per kit policy.

**Optional:** when `rckt` is on `PATH`, use `rckt run …` for durable delivery (rckt-factory).
