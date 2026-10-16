# Antigravity projection

Prefer **`muster link`**: `agy --agent muster-orchestrator` via `~/.gemini/antigravity-cli/settings.json`.

```bash
export KIT="$(muster kit path)"
agy agents | rg muster
jq -r '.agent' ~/.gemini/antigravity-cli/settings.json
# Orchestrator policy: $KIT/policy/orchestrator.md
```

CLI binary `agy` (harness id `antigravity`). Use `/plan` for substantial or ambiguous work and `/goal` for explicitly continuing work with a verifiable finish line. Delegate bounded work to native worker subagents per kit policy. Let lifecycle notifications deliver results; use `/tasks` only to inspect a task when needed, not as a polling loop. Assignments and handoffs are ordinary conversation.

