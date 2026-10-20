package link

import (
	"os"
	"path/filepath"
	"regexp"
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

var codexTomlSection = regexp.MustCompile(`(?m)^\[([^]]+)\][ \t]*(?:#.*)?$`)

// patchCodexAgentConfig updates only Codex's subagent defaults and worker
// registration, preserving all unrelated TOML sections and keys.
func patchCodexAgentConfig(path, workerConfig string, dryRun bool) error {
	return patchFile(path, func(s string) string {
		return upsertCodexAgentConfig(s, workerConfig)
	}, dryRun)
}

func unlinkCodexAgentConfig(path string, dryRun bool) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return patchFile(path, removeCodexAgentConfig, dryRun)
}

func upsertCodexAgentConfig(s, workerConfig string) string {
	s = setTomlSectionValues(s, "agents", map[string]string{
		"default_subagent_model":            "gpt-6-luna",
		"default_subagent_reasoning_effort": "medium",
	})
	return setTomlSectionValues(s, "agents.muster-worker", map[string]string{"config_file": workerConfig})
}

func setTomlSectionValues(s, name string, values map[string]string) string {
	sections := codexTomlSection.FindAllStringSubmatchIndex(s, -1)
	start, end := -1, len(s)
	for i, sec := range sections {
		if s[sec[2]:sec[3]] != name {
			continue
		}
		start = sec[1]
		if i+1 < len(sections) {
			end = sections[i+1][0]
		}
		break
	}
	if start < 0 {
		if s != "" && !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		s += "\n[" + name + "]\n"
		for k, v := range values {
			s += k + " = " + tomlQuote(v) + "\n"
		}
		return s
	}
	headerEnd := strings.IndexByte(s[start:], '\n')
	if headerEnd < 0 {
		headerEnd = len(s) - start
	}
	bodyStart := start + headerEnd
	body := s[bodyStart:end]
	for k, v := range values {
		body = setTomlValue(body, k, tomlQuote(v))
	}
	return s[:bodyStart] + body + s[end:]
}

func setTomlValue(body, key, value string) string {
	line := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(key) + `[ \t]*=.*(?:\n|$)`)
	if line.MatchString(body) {
		return line.ReplaceAllString(body, key+" = "+value+"\n")
	}
	if body != "" && !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	return body + key + " = " + value + "\n"
}

func tomlQuote(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}

func removeCodexAgentConfig(s string) string {
	sections := codexTomlSection.FindAllStringSubmatchIndex(s, -1)
	for i := len(sections) - 1; i >= 0; i-- {
		sec := sections[i]
		name := s[sec[2]:sec[3]]
		keys := []string(nil)
		if name == "agents" {
			keys = []string{"default_subagent_model", "default_subagent_reasoning_effort"}
		}
		if name == "agents.muster-worker" {
			keys = []string{"config_file"}
		}
		if len(keys) == 0 {
			continue
		}
		end := len(s)
		if i+1 < len(sections) {
			end = sections[i+1][0]
		}
		bodyStart := sec[1]
		headerEnd := strings.IndexByte(s[bodyStart:], '\n')
		if headerEnd < 0 {
			headerEnd = len(s) - bodyStart
		}
		bodyStart += headerEnd
		body := s[bodyStart:end]
		for _, key := range keys {
			body = regexp.MustCompile(`(?m)^`+regexp.QuoteMeta(key)+`[ \t]*=.*(?:\n|$)`).ReplaceAllString(body, "")
		}
		if strings.TrimSpace(body) == "" {
			s = s[:sec[0]] + s[end:]
		} else {
			s = s[:bodyStart] + body + s[end:]
		}
		sections = codexTomlSection.FindAllStringSubmatchIndex(s, -1)
	}
	return strings.TrimLeft(s, "\n")
}
