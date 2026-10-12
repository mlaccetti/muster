package link

import (
	"path/filepath"
)

func patchOrchestratorAgents(path string, kitDir string, dryRun bool) error {
	return patchFile(path, func(s string) string { return upsertManagedBlock(s, orchestratorBlock(kitDir)) }, dryRun)
}

func unlinkOrchestratorAgents(path string, dryRun bool) error {
	return patchFile(path, func(s string) string { return removeManagedBlock(s) }, dryRun)
}

func linkProjectAgentsMD(projectDir, kitDir string, dryRun bool) error {
	return patchOrchestratorAgents(filepath.Join(projectDir, "AGENTS.md"), kitDir, dryRun)
}

func unlinkProjectAgentsMD(projectDir string, dryRun bool) error {
	return unlinkOrchestratorAgents(filepath.Join(projectDir, "AGENTS.md"), dryRun)
}
