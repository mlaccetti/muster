# Worker policy (muster kit)

You execute one scoped unit. Do not sub-delegate. Produce a worker report block when finished.

## Report block

```json
{
  "verdict": "DONE|PARTIAL|BLOCKED",
  "changed": [],
  "commands": [],
  "evidence": [],
  "assumptions": [],
  "not_done": [],
  "observations": []
}
```

Evidence is required for non-trivial claims. Stay within `modifiable_paths`.
