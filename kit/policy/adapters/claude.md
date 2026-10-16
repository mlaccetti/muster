## Claude Code controls

- Delegate bounded work to native `Agent` subagents. Run independent assignments concurrently, keep ownership distinct, and integrate their completion notifications. Do not create an experimental agent team for ordinary orchestration.
- Use `EnterPlanMode` and `ExitPlanMode` when substantial, ambiguous, or high-risk work benefits from a reviewable plan. Use `TodoWrite` for the orchestrator's short list of substantial steps, not as a worker contract or required artifact.
- Use `AskUserQuestion` only when a material decision blocks progress. Make low-impact assumptions and continue when safe.
- Use `/goal <condition>` only for an explicitly continuing task with concrete evidence and a bounded stop condition. Do not use `/loop` or repeated `/tasks` checks to wait for an agent, command, or CI event. Let Claude Code deliver background completion notifications; use `/tasks` to inspect a task when needed.
