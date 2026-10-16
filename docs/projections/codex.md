# Codex projection

Use **`muster link`** to add the orchestrator policy to Codex's native global `AGENTS.md` instructions and install a native worker agent.

Codex has no separate root-agent selector in this setup: the active Codex session is the orchestrator because native instruction discovery loads the managed global guidance. The `muster-worker` custom agent is available for delegation.

```bash
rg 'MUSTER_LINK_START|default orchestrator' ~/.codex/AGENTS.md
readlink ~/.codex/agents/muster-worker.toml
codex debug prompt-input "Summarize the active Muster instructions"
```

The `muster-worker` agent uses `gpt-6-luna` with `medium` reasoning effort. The active Codex session is the orchestrator; delegate bounded implementation and investigation to its native subagent/worker controls, and keep assignments and handoffs in ordinary conversation.

Use the native plan or task-tracking controls exposed by the active Codex session for substantial work. In Codex builds that expose Plan mode, use it to shape the approach before risky or broad changes; the plan is for the path, not the completion condition. Do not create TODO files to imitate a native tracker. Use `/goal` only when the user explicitly wants work to continue across turns, with a verifiable outcome and constraints. A Goal is not permission to keep working indefinitely and does not replace worker decomposition or evidence checks.

Ask the user directly only when a material decision blocks progress; Codex does not have a portable `AskUserQuestion` tool. Wait on the session's native worker lifecycle events or callbacks rather than polling or sleeping. Exact controls vary across Codex app, CLI, IDE, and enabled tools, so instructions should name a control only when that surface exposes it.

