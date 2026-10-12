# Pools (`pools.v1`)

Default pools ship in the kit as `pools.json` and copy to `$MUSTER_HOME/pools.json` on `muster install`.

```json
{
  "version": 1,
  "pass_order": ["claude", "codex", "cursor"],
  "tiers": {
    "implement": [
      {"harness": "claude", "model": "", "family": "anthropic", "write": true}
    ]
  }
}
```

Fields:

- `harness`: `claude`, `codex`, `cursor`, `antigravity`
- `family`: used for `--avoid-family` filtering
- `write`: implement tier skips `write: false` candidates
- `boundary`: optional `hook` or `sandbox` (codex)

Validate against schema `kit/schemas/pools.v1.json` (loaded on dispatch).
