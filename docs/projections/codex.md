# Codex projection

Use **`muster link`** to compose the shared orchestrator policy with Codex's native-control adapter and add the rendered prompt to Codex's global `AGENTS.md`. Linking preserves the root session's selected model and reasoning effort; it does not switch the main UI session's model. Choose a **Sol** model for the root when you want Sol orchestration. Muster registers its `muster-worker` native agent with the **Luna Medium** defaults, so delegated workers receive narrow assignments while the root keeps coordination.

If Codex cannot load the registered custom role, the instructions also describe a built-in worker fallback. A fresh Codex session or app restart may be needed for newly linked or repaired agent definitions to replace roles cached by an existing session.

During `muster link`, Muster narrowly repairs the legacy Claude-compatibility mapping that could otherwise interfere with Codex worker selection. It does not change unrelated Claude configuration.

```bash
rg 'MUSTER_LINK_START|default orchestrator' ~/.codex/AGENTS.md
readlink ~/.codex/agents/muster-worker.toml
codex debug prompt-input "Summarize the active Muster instructions"
```

The regular `muster doctor` validates the installed kit and recorded links. `muster doctor --probe` checks harness binary versions. The opt-in **`muster doctor --live`** starts a read-only task in an isolated context and checks actual session and event evidence that the Sol root delegated to a completed Luna Medium worker. For this live check, validation fails closed if the root or worker uses the wrong model or effort, or if delegation is absent. The check consumes usage from the signed-in Codex account.

```bash
muster doctor --live
muster doctor --live --live-timeout 2m
muster doctor --live --codex-bin "$(command -v codex)"
```

Use `--live-timeout` to set the live check's duration. `--codex-bin` selects the executable to test; inspect the Codex desktop app bundle on your system if you need to point it at that exact binary.
