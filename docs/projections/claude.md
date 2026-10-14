# Claude Code projection

Prefer **`muster link`** (see [linking.md](../linking.md)) so `~/.claude/settings.json` selects `muster-orchestrator` and agent definitions symlink into the installed kit.

Manual check:

```bash
export KIT="$(muster kit path)"
jq -r '.agent' ~/.claude/settings.json
readlink ~/.claude/agents/muster-orchestrator.md
# Policy: $KIT/policy/orchestrator.md
```

Orchestrator delegates units to **worker** subagents per `$KIT/policy/worker.md`.

**Optional:** when `rckt` is on `PATH`, orchestrator policy allows `rckt run …` for durable delivery (rckt-factory).
