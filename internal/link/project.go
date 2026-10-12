package link

import (
	"os"
	"path/filepath"
)

func linkProject(projectDir, kitDir string, dryRun bool) (string, error) {
	dir := filepath.Join(projectDir, ".muster")
	if dryRun {
		return dir, nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	readme := `# Muster project wiring

Orchestrator policy: ` + kitDir + `/policy/orchestrator.md
Worker policy: ` + kitDir + `/policy/worker.md

Harness-specific entrypoints are linked by ` + "`muster link`" + ` (see $MUSTER_HOME/link-manifest.json).
`
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(readme), 0o644); err != nil {
		return "", err
	}
	// Stable relative pointers for tools that read from the repo.
	if err := os.WriteFile(filepath.Join(dir, "orchestrator.path"), []byte(kitDir+"/policy/orchestrator.md\n"), 0o644); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(dir, "worker.path"), []byte(kitDir+"/policy/worker.md\n"), 0o644); err != nil {
		return "", err
	}
	return dir, nil
}

func unlinkProject(projectDir string, dryRun bool) error {
	dir := filepath.Join(projectDir, ".muster")
	if dryRun {
		return nil
	}
	return os.RemoveAll(dir)
}
