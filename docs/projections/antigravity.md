# Antigravity projection

Prefer **`muster link`**: `agy --agent muster-orchestrator` via `~/.gemini/antigravity-cli/settings.json`.

```bash
agy agents | rg muster
jq -r '.agent' ~/.gemini/antigravity-cli/settings.json
# Rendered from policy/orchestrator.md and policy/adapters/antigravity.md
```

The rendered agent combines the shared policy with Antigravity's native `/plan`, `/goal`, subagent, and lifecycle-notification controls. CLI binary `agy` (harness id `antigravity`).

