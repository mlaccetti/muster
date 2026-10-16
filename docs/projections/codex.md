# Codex projection

Use **`muster link`** to compose the shared orchestrator policy with Codex's native-control adapter and add the rendered prompt to Codex's global `AGENTS.md`. It also installs a native worker agent.

Codex has no separate root-agent selector in this setup: the active Codex session is the orchestrator because native instruction discovery loads the managed global guidance. The `muster-worker` custom agent is available for delegation.

```bash
rg 'MUSTER_LINK_START|default orchestrator' ~/.codex/AGENTS.md
readlink ~/.codex/agents/muster-worker.toml
codex debug prompt-input "Summarize the active Muster instructions"
```

The `muster-worker` agent uses `gpt-6-luna` with `medium` reasoning effort. The active Codex session is the orchestrator; the installed instructions combine `kit/policy/orchestrator.md` with `kit/policy/adapters/codex.md`. This keeps shared orchestration behavior in one source while limiting the adapter to Codex-specific controls and availability caveats.

