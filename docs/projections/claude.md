# Claude Code projection

Prefer **`muster link`** (see [linking.md](../linking.md)) so `~/.claude/settings.json` selects `muster-orchestrator` and its rendered definition points to the active kit.

Manual check:

```bash
jq -r '.agent' ~/.claude/settings.json
readlink ~/.claude/agents/muster-orchestrator.md
# Rendered from policy/orchestrator.md and policy/adapters/claude.md
```

The rendered agent combines the shared policy with Claude Code's native `Agent`, planning, task-list, question, and completion-notification controls. Claude workers use the harness's native model settings.
