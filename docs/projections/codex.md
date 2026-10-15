# Codex projection

Prefer **`muster link`**: it adds the orchestrator policy to native global and project `AGENTS.md` instructions and installs a native worker agent.

Codex has no separate root-agent selector in this setup: the active Codex session is the orchestrator because native instruction discovery loads the managed global and project guidance. The `muster-worker` custom agent is available for delegation.

```bash
rg 'MUSTER_LINK_START|default orchestrator' ~/.codex/AGENTS.md
readlink ~/.codex/agents/muster-worker.toml
codex debug prompt-input "Summarize the active Muster instructions"
```

The `muster-worker` agent uses `gpt-6-luna` with `medium` reasoning effort. Assignments and handoffs remain ordinary conversation.

