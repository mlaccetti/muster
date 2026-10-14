# Linking harnesses

`muster install` lays down the versioned **kit** (policy, schemas, `link/` harness entrypoints). **`muster link`** symlinks those entrypoints into each harness and sets the **default orchestrator** selector. No prose is copied into `AGENTS.md` or harness markdown.

## Behavior

| Mode | Default | What it does |
|------|---------|----------------|
| Auto-detect | yes (`muster link`) | Links only harnesses whose CLIs are on `PATH` (`claude`, `codex`, `cursor-agent`, `agy`) |
| `--harness` | optional | Link specific IDs even if not on `PATH` |
| `--all` | off | Link all four supported harnesses |
| `--global` | on | User-level harness config |
| `--project` | on | Repo-level harness config + `.muster/` pointers |
| `--project-dir` | cwd | Repo root to link |
| `--unlink` | off | Remove symlinks, config patches, and shims |

State is recorded in `$MUSTER_HOME/link-manifest.json`. `muster doctor` reports link status and runs **`link.VerifyReports`** on each linked harness.

## Per harness (global)

| ID | Default orchestrator | Definitions (symlinks into `$KIT/link/...`) |
|----|----------------------|---------------------------------------------|
| `claude` | `~/.claude/settings.json` → `"agent": "muster-orchestrator"` | `~/.claude/agents/muster-orchestrator.md`, `muster-worker.md` |
| `codex` | `~/.codex/config.toml` managed block: `model_instructions_file` + `developer_instructions` (every `codex` run). Profile: `codex -p muster-orchestrator` | `~/.codex/muster-orchestrator.config.toml`, `muster-worker.config.toml`; `[agents.muster-worker]` in `config.toml` |
| `cursor` | `~/.cursor/rules/muster-orchestrator.mdc` with `alwaysApply: true` | `muster-orchestrator.mdc`, `muster-worker.mdc` |
| `antigravity` (`agy`) | `~/.gemini/antigravity-cli/settings.json` → `"agent": "muster-orchestrator"` | `~/.gemini/config/agents/muster-orchestrator/agent.md`, `muster-worker/agent.md` |

Project linking adds `.claude/`, `.codex/`, `.cursor/rules/`, or `.agents/agents/` equivalents under the repo.

## Optional shims

After link, `$MUSTER_HOME/bin/` contains:

- **`codex-muster`**: `codex -p muster-orchestrator "$@"`
- **`agy-muster`**: `agy --agent muster-orchestrator "$@"`

Add `$MUSTER_HOME/bin` to `PATH` if you want those names instead of the stock CLIs.

## Validation

```bash
muster install
muster link --all
muster doctor          # link.VerifyReports on manifest harnesses
muster doctor --probe  # harness binaries on PATH
```

**Claude**

```bash
jq -r '.agent' ~/.claude/settings.json
readlink ~/.claude/agents/muster-orchestrator.md
```

**Codex**

```bash
rg MUSTER_LINK_BEGIN ~/.codex/config.toml
readlink ~/.codex/muster-orchestrator.config.toml
codex -p muster-orchestrator -h   # profile recognized
```

**cursor-agent**

```bash
test -f ~/.cursor/rules/muster-orchestrator.mdc
cursor-agent -p --print "List always-applied rules that mention muster or muster kit path"
```

**agy**

```bash
agy agents | rg muster
jq -r '.agent' ~/.gemini/antigravity-cli/settings.json
agy -p "Reply with your agent name only" --agent muster-orchestrator
```

Re-run **`muster link`** after `muster install` upgrades the kit so symlinks target the active kit tree.

## Optional durable runs (rckt-factory)

muster does not require **rckt-factory**. Linked orchestrator policy tells the main session to delegate units to **worker** subagents in the harness. If the operator has installed factory and `rckt` is on `PATH`, the same policy allows `rckt run …` for durable delivery instead of only in-chat delegation.

## Testing without real harnesses

CI uses `muster link --harness claude --project <dir> --global=false`. Package tests cover codex/cursor/agy global link with temp `HOME` and `CODEX_HOME`.
