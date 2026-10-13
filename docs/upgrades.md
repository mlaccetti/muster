# Upgrades and new releases

## CLI binary

After `brew upgrade muster` (or `go install`), run:

```bash
muster install   # refresh embedded kit for this binary
muster link      # refresh harness snippets if you use link
```

The CLI checks GitHub **at most once per day** for a newer release and prints a hint to stderr:

```text
muster v0.2.0 is available (you have v0.1.0). Run: brew update && brew upgrade muster
```

Skip the check in CI with `MUSTER_SKIP_UPDATE_CHECK=1`. `dev` builds never notify.

## Kit on disk

`muster install` writes `$MUSTER_HOME/kit/<version>/` and updates the `kit/current` symlink. Old `kit/vX.Y.Z/` trees are kept for rollback. Harness pools and dispatch live in **rckt-factory**, not muster.

When the kit version changes, run **`muster link`** again so orchestrator/worker snippets pick up new absolute paths.

`muster doctor` warns if `link-manifest.json` points at an older kit directory than `muster kit path`.

## Maintainers

Bump `paths.KitVersion` when kit data changes, tag `v*`, publish via GoReleaser.
