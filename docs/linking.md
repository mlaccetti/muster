# Linking harnesses

`muster install` lays down the versioned **kit** (prose policies and `link/` harness entrypoints). **`muster link`** symlinks those entrypoints into each harness and sets the **default orchestrator** selector. Assignments and handoffs are ordinary harness conversation; Muster requires no fixed message format.

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

| ID | Native activation | Linked definition |
|----|----------------------|---------------------------------------------|
| `claude` | `~/.claude/settings.json` → `"agent": "muster-orchestrator"` | `~/.claude/agents/muster-orchestrator.md`, `muster-worker.md` |
| `codex` | `~/.codex/AGENTS.md` (every Codex run) | `~/.codex/agents/muster-worker.toml` |
| `cursor` | Project-scoped `.cursor/rules/muster-orchestrator.mdc` with `alwaysApply: true`; global User Rules must be set in Cursor Settings | `.cursor/rules/muster-orchestrator.mdc` |
| `antigravity` (`agy`) | `~/.gemini/antigravity-cli/settings.json` → `"agent": "muster-orchestrator"` | `~/.gemini/config/agents/muster-orchestrator/agent.md`, `muster-worker/agent.md` |

Codex also receives the orchestrator instructions in the linked project's root `AGENTS.md`. Its native worker agent uses `gpt-6-luna` with `medium` reasoning effort. Cursor CLI automatically loads the project rule; global Cursor User Rules must be configured in Cursor Settings.

Project linking adds `.claude/`, `.codex/`, `.cursor/rules/`, or `.agents/agents/` equivalents under the repo.

## Optional shims

After link, `$MUSTER_HOME/bin/` contains:

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
rg 'MUSTER_LINK_START|default orchestrator' ~/.codex/AGENTS.md
readlink ~/.codex/agents/muster-worker.toml
codex debug prompt-input "Summarize the active Muster instructions"  # inspect loaded prompt inputs
```

**cursor-agent**

```bash
test -f .cursor/rules/muster-orchestrator.mdc
cursor-agent -p --output-format text "Summarize the active Muster project rule"
```

**agy**

```bash
agy agents | rg muster
jq -r '.agent' ~/.gemini/antigravity-cli/settings.json
agy -p "Reply with your agent name only" --agent muster-orchestrator
```

Re-run **`muster link`** after `muster install` upgrades the kit so symlinks target the active kit tree.

## Testing without real harnesses

CI uses `muster link --harness claude --project <dir> --global=false`. Package tests cover Codex global and project `AGENTS.md` instruction wiring, native worker agent files, Cursor project rules, and Antigravity global linking with temporary homes.
