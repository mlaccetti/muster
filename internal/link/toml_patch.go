package link

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	tomlMarkerBegin = "# MUSTER_LINK_BEGIN"
	tomlMarkerEnd   = "# MUSTER_LINK_END"
	tomlAgentsBegin = "# MUSTER_LINK_AGENTS_BEGIN"
	tomlAgentsEnd   = "# MUSTER_LINK_AGENTS_END"
)

func upsertMarkedBlock(existing, begin, end, block string) string {
	if strings.Contains(existing, begin) && strings.Contains(existing, end) {
		before, _, _ := strings.Cut(existing, begin)
		_, after, _ := strings.Cut(existing, end)
		after = strings.TrimPrefix(after, "\n")
		out := strings.TrimRight(before, "\n") + "\n\n" + block + "\n"
		if strings.TrimSpace(after) != "" {
			out += after
		}
		return out
	}
	trim := strings.TrimRight(existing, "\n")
	if trim == "" {
		return block + "\n"
	}
	return trim + "\n\n" + block + "\n"
}

func removeMarkedBlock(existing, begin, end string) string {
	if !strings.Contains(existing, begin) {
		return existing
	}
	before, _, _ := strings.Cut(existing, begin)
	_, after, ok := strings.Cut(existing, end)
	if !ok {
		return strings.TrimRight(before, "\n") + "\n"
	}
	after = strings.TrimPrefix(after, "\n")
	return strings.TrimRight(before, "\n") + "\n" + after
}

func patchTomlFile(path string, block string, dryRun bool) error {
	return patchFile(path, func(s string) string { return upsertMarkedBlock(s, tomlMarkerBegin, tomlMarkerEnd, block) }, dryRun)
}

func unlinkTomlFile(path string, dryRun bool) error {
	return patchFile(path, func(s string) string { return removeMarkedBlock(s, tomlMarkerBegin, tomlMarkerEnd) }, dryRun)
}

func codexAgentsBlock(codexDir string) string {
	workerCfg := filepath.Join(codexDir, workerAgentName+".config.toml")
	return fmt.Sprintf(`%s
[agents.%s]
description = "Muster worker for one bounded task with a concise conversational handoff"
config_file = %q
%s
`, tomlAgentsBegin, workerAgentName, workerCfg, tomlAgentsEnd)
}

func patchCodexAgentsRole(codexDir string, dryRun bool) error {
	path := filepath.Join(codexDir, "config.toml")
	block := codexAgentsBlock(codexDir)
	return patchFile(path, func(s string) string { return upsertMarkedBlock(s, tomlAgentsBegin, tomlAgentsEnd, block) }, dryRun)
}

func unlinkCodexAgentsRole(codexDir string, dryRun bool) error {
	path := filepath.Join(codexDir, "config.toml")
	return patchFile(path, func(s string) string { return removeMarkedBlock(s, tomlAgentsBegin, tomlAgentsEnd) }, dryRun)
}

func ensureFile(path string, dryRun bool) error {
	if dryRun {
		return nil
	}
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte{}, 0o644)
}
