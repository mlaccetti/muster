# Claude Code projection

Prefer **`muster link`** (see [linking.md](../linking.md)) so `~/.claude/settings.json` selects `muster-orchestrator` and agent definitions symlink into the installed kit.

Manual check:

```bash
export KIT="$(muster kit path)"
jq -r '.agent' ~/.claude/settings.json
readlink ~/.claude/agents/muster-orchestrator.md
# Policy: $KIT/policy/orchestrator.md
```

The default orchestrator uses Claude Code's native `Agent` subagents for bounded parallel work and follows `$KIT/policy/worker.md` for their scope and handoff. It uses `EnterPlanMode`/`ExitPlanMode` and `TodoWrite` for substantial multi-step work, and `AskUserQuestion` only when a material decision blocks progress. Claude Code's `/goal` is for an explicitly continuing task with a verifiable finish condition; `/loop` is time-based and must not be used to wait for workers or CI. Background task completion arrives as a notification; `/tasks` can inspect details without creating a polling loop. Keep assignments and handoffs conversational; no fixed format is required.
