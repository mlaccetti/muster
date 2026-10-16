## Codex controls

- Use Codex's native plan or task-tracking controls when the active session exposes them. Plan mode shapes the approach before risky or broad changes; do not create TODO files to imitate a missing tracker.
- Use `/goal` only when the user explicitly wants work to continue across turns, with a verifiable outcome and constraints. A Goal does not replace worker decomposition or evidence checks.
- Delegate bounded implementation and investigation through the native subagent or worker controls exposed by the active Codex session. Keep assignments and handoffs in ordinary conversation.
- Ask the user directly only when a material decision blocks progress; Codex has no portable `AskUserQuestion` tool. Wait on native worker lifecycle events or callbacks instead of polling or sleeping. Name controls only when the active Codex surface exposes them.
