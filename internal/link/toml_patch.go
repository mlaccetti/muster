package link

import (
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

func unlinkTomlFile(path string, dryRun bool) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return patchFile(path, func(s string) string { return removeMarkedBlock(s, tomlMarkerBegin, tomlMarkerEnd) }, dryRun)
}

func unlinkCodexAgentsRole(codexDir string, dryRun bool) error {
	path := filepath.Join(codexDir, "config.toml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return patchFile(path, func(s string) string { return removeMarkedBlock(s, tomlAgentsBegin, tomlAgentsEnd) }, dryRun)
}
