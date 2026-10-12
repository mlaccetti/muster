# Linking harnesses

`muster install` lays down the kit. **`muster link`** wires the same orchestrator and worker policies into each harness so behavior stays consistent.

## Behavior

| Mode | Default | What it does |
|------|---------|----------------|
| Auto-detect | yes (`muster link`) | Links only harnesses whose CLIs are on `PATH` (`claude`, `codex`, `cursor-agent`, `agy`) |
| `--harness` | optional | Link specific IDs even if not on `PATH` |
| `--all` | off | Link all four supported harnesses |
| `--global` | on | User config (`~/.claude`, `~/.codex`, `~/.cursor/rules`, `~/.antigravity`) |
| `--project` | on | Repo config (`CLAUDE.md`, `AGENTS.md`, `.cursor/rules`, `.muster/`) |
| `--project-dir` | cwd | Repo root to link |

Managed blocks are delimited with `<!-- MUSTER_LINK_START -->` / `<!-- MUSTER_LINK_END -->` so re-running `muster link` updates in place. `muster link --unlink` removes them.

## Per harness

| ID | Global | Project |
|----|--------|---------|
| `claude` | `~/.claude/CLAUDE.md` + `muster-worker.md` | `CLAUDE.md` + `.muster/claude-worker.md` |
| `codex` | `~/.codex/AGENTS.md` + `muster-worker.md` | `AGENTS.md` + `.muster/codex-worker.md` |
| `cursor` | `~/.cursor/rules/muster-*.mdc` | `.cursor/rules/muster-orchestrator.mdc` |
| `antigravity` | `~/.antigravity/muster-*.md` | `.muster/antigravity/` |

State is recorded in `$MUSTER_HOME/link-manifest.json`. `muster doctor` reports link status.

## Testing without real harnesses

CI uses `muster link --harness claude --project <dir> --global=false` (no `PATH` probe). E2E subprocess tests can do the same after `muster install` in a temp `MUSTER_HOME`.
