# Linking harnesses

`muster install` lays down the versioned **kit** (policy and harness entrypoints). **`muster link`** connects those entrypoints to harness-wide user configuration and selects Muster as the default orchestrator where the harness supports that setting. Link and unlink never write into the current repository.

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
| `cursor` | Add Muster policy to User Rules in Cursor Settings | Manual setup; the documented interface is Settings, with no supported CLI or file interface for configuring global User Rules |
| `antigravity` (`agy`) | `~/.gemini/antigravity-cli/settings.json` → `"agent": "muster-orchestrator"` | `~/.gemini/config/agents/muster-orchestrator/agent.md`, `muster-worker/agent.md` |

Codex's worker agent uses `gpt-6-luna` with `medium` reasoning effort. Muster reports Cursor as skipped because it cannot configure Cursor's global User Rules through a supported interface. To use Muster in Cursor, copy the contents of `muster kit path` + `/policy/orchestrator.md` into Cursor Settings → Rules → User Rules. Cursor's worker and event behavior depends on the capabilities exposed by the active session.

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

Confirm Muster's orchestrator policy appears under Cursor Settings → Rules → User Rules.

**agy**

```bash
agy agents | rg muster
jq -r '.agent' ~/.gemini/antigravity-cli/settings.json
agy -p "Reply with your agent name only" --agent muster-orchestrator
```

Re-run **`muster link`** after `muster install` upgrades the kit so symlinks target the active kit tree.

## Testing without real harnesses

Package and CLI tests use temporary homes to validate harness-wide configuration, worker definitions, preservation of existing settings, unlink behavior, and the Cursor Settings limitation. Tests never write into the caller's project directory.
