package link

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	legacySequentialMapping = "- Task (subagent dispatch) / Subagent / Parallel: run sequentially in main thread; use multi_tool_use.parallel for tool calls"
	nativeCodexMapping      = "- Task (subagent dispatch) / Subagent / Parallel: use Codex's native subagent controls for delegated work; run independent subtasks concurrently and integrate their results."
)

func codexHome() string {
	if v := os.Getenv("CODEX_HOME"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".codex")
}

func linkCodexGlobal(kitDir string, dryRun bool) ([]string, error) {
	dir := codexHome()
	if dir == "" {
		return nil, os.ErrNotExist
	}
	if !dryRun {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	var touched []string
	configPath := filepath.Join(dir, "config.toml")
	if err := patchCodexAgentConfig(configPath, true); err != nil {
		return nil, err
	}
	if !dryRun {
		_ = removeLink(filepath.Join(dir, "muster-orchestrator.config.toml"))
		_ = removeLink(filepath.Join(dir, workerAgentName+".config.toml"))
	}
	workerAgent := filepath.Join(dir, "agents", workerAgentName+".toml")
	if err := linkCodexWorkerAgent(workerAgent, kitDir, dryRun); err != nil {
		return nil, err
	}
	touched = append(touched, workerAgent)
	if err := patchCodexAgentConfig(configPath, dryRun); err != nil {
		return nil, err
	}
	touched = append(touched, configPath)

	globalAgents := filepath.Join(dir, "AGENTS.md")
	block, err := orchestratorBlock(kitDir)
	if err != nil {
		return nil, err
	}
	if err := patchFile(globalAgents, func(s string) string {
		s = strings.ReplaceAll(s, legacySequentialMapping, nativeCodexMapping)
		return upsertManagedBlock(s, block)
	}, dryRun); err != nil {
		return nil, err
	}
	touched = append(touched, globalAgents)

	// Legacy prose cleanup.
	if !dryRun {
		_ = os.Remove(filepath.Join(dir, "muster-worker.md"))
	}
	return touched, nil
}

func unlinkCodexGlobal(dryRun bool) error {
	dir := codexHome()
	if err := unlinkCodexAgentConfig(filepath.Join(dir, "config.toml"), dryRun); err != nil {
		return err
	}
	if !dryRun {
		_ = removeLink(filepath.Join(dir, "agents", workerAgentName+".toml"))
		_ = removeLink(filepath.Join(dir, "muster-orchestrator.config.toml"))
		_ = removeLink(filepath.Join(dir, workerAgentName+".config.toml"))
		_ = patchFile(filepath.Join(dir, "AGENTS.md"), func(s string) string {
			s = removeManagedBlock(s)
			return strings.ReplaceAll(s, nativeCodexMapping, legacySequentialMapping)
		}, false)
		_ = os.Remove(filepath.Join(dir, "muster-worker.md"))
	}
	return nil
}
