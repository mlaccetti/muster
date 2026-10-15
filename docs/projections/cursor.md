# Cursor Agent projection

Prefer **`muster link --project`**: it installs `.cursor/rules/muster-orchestrator.mdc` with `alwaysApply: true`. Cursor CLI loads project rules. Global User Rules are managed in Cursor Settings, so Muster cannot set them through a supported file path.

```bash
test -f .cursor/rules/muster-orchestrator.mdc
cat .muster/orchestrator.path
```

Binary resolved as `cursor-agent`. The linked always-applied project rule makes the current session follow the Muster orchestrator policy. Cursor CLI documentation currently describes a single CLI agent and does not document a native subagent dispatch or lifecycle callback interface. The rule therefore tells the agent to use only worker and event tools actually exposed by the session; Cursor Background Agents are a separate capability.

