# Linking harnesses

`muster install` lays down the versioned **kit** (shared policy, small harness adapters, and native templates). **`muster link`** composes the shared policy with each adapter, renders the native orchestrator instructions, and connects them to harness-wide user configuration. Link and unlink never write into the current repository.

## Behavior

| Command or option | Behavior |
|-------------------|----------|
| `muster link` | Auto-detect installed harness CLIs and link their harness-wide configuration |
| `--harness <id>` | Link a specific harness, whether or not its CLI is on `PATH` |
| `--all` | Attempt all supported harnesses |
| `--unlink` | Remove Muster's harness-wide configuration |

State for successfully linked harnesses is recorded in `$MUSTER_HOME/link-manifest.json`. `muster doctor` checks each recorded link.

## Harness support

| ID | Harness-wide activation | Linked definition |
|----|------------------------|-------------------|
| `claude` | `~/.claude/settings.json` → `"agent": "muster-orchestrator"` | `~/.claude/agents/muster-orchestrator.md` |
| `codex` | `~/.codex/AGENTS.md` (loaded for every Codex run) | `~/.codex/agents/muster-worker.toml` |
| `cursor` | User-scoped local plugin discovered from `~/.cursor/plugins/local/` | `~/.cursor/plugins/local/muster-orchestrator/` (always-applied global rule) |
| `antigravity` (`agy`) | `~/.gemini/antigravity-cli/settings.json` → `"agent": "muster-orchestrator"` | `~/.gemini/config/agents/muster-orchestrator/agent.md`, `muster-worker/agent.md` |

Claude Code and Antigravity point to rendered orchestrator agents under the installed kit's `rendered/` directory. Codex receives the rendered prompt in its managed global `AGENTS.md` block. Cursor receives a rendered user-scoped local plugin rule. The source of truth is `kit/policy/orchestrator.md` plus one short adapter in `kit/policy/adapters/` per harness. Codex's worker agent uses `gpt-6-luna` with `medium` reasoning effort. Cursor discovers the generated local plugin after a restart or Developer: Reload Window. Local plugin imports can be disabled by a Cursor Enterprise administrator. Cursor's plan and to-do features are unavailable in Auto mode; worker and event behavior depends on the active session.

## Optional shims

After link, `$MUSTER_HOME/bin/` contains `agy-muster`, which runs `agy --agent muster-orchestrator "$@"`. Add that directory to `PATH` to use the shim.

## Validation

```bash
muster install
muster link --all
muster doctor          # validate recorded harness-wide links
muster doctor --probe  # probe harness binaries on PATH
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
codex debug prompt-input "Summarize the active Muster instructions"
```

**Cursor**

Confirm `~/.cursor/plugins/local/muster-orchestrator/rules/muster-orchestrator.mdc` exists, then restart Cursor or run Developer: Reload Window. The rule should appear under the installed local plugin and apply to Agent chats.

**agy**

```bash
agy agents | rg muster
jq -r '.agent' ~/.gemini/antigravity-cli/settings.json
agy -p "Reply with your agent name only" --agent muster-orchestrator
```

Re-run **`muster link`** after `muster install` upgrades the kit so native instructions are rendered from the active kit and symlinks target its rendered agents.

## Testing without real harnesses

Package and CLI tests use temporary homes to validate shared-policy rendering, harness adapters, harness-wide configuration, worker definitions, preservation of existing settings, unlink behavior, and the generated Cursor local plugin. Tests never write into the caller's project directory.
