# Pass spec (YAML)

`muster pass run --file pass.yaml` executes units and optional exec gates. No LLM nodes in v0.1.

```yaml
id: my-pass
units:
  - name: unit1
    brief: /path/to/brief.json   # or rendered brief.txt
    harness: claude              # or omit harness and set tier for dispatch
    tier: implement
gates:
  - name: smoke
    command: test -s unit1-report.json
    timeout_seconds: 30
```

State is written under `$MUSTER_HOME/passes/<id>/` (override with `--pass-dir`). Routing lines append to `ledger.jsonl`. Gates run via `sh -c` with a `context` timeout (no sleep polling).

Exit non-zero when a gate fails or a unit returns a non-ok routing verdict.
