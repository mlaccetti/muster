# Codex projection

Prefer **`muster link`**: `codex -p muster-orchestrator` and symlinks under `~/.codex/`.

```bash
export KIT="$(muster kit path)"
rg MUSTER_LINK_BEGIN ~/.codex/config.toml
readlink ~/.codex/muster-orchestrator.config.toml
# Orchestrator policy: $KIT/policy/orchestrator.md
```

Workers follow `$KIT/policy/worker.md` when invoked as the muster-worker agent profile.

