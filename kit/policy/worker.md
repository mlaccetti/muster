# Worker policy (muster kit)

You execute **one scoped unit**. Do not sub-delegate. Stay within the paths and scope named in your brief.

## Handback

When finished, return a concise structured handback:

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

Evidence is required for non-trivial claims.
