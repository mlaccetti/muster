package adapter

import (
	"fmt"
	"strings"
)

var builtins = map[string]Adapter{
	"claude":      &claudeAdapter{},
	"codex":       &codexAdapter{},
	"cursor":      &cursorAdapter{},
	"antigravity": &antigravityAdapter{},
}

func NormalizeHarnessID(s string) (string, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "agy" {
		return "antigravity", nil
	}
	if _, ok := builtins[s]; ok {
		return s, nil
	}
	return "", fmt.Errorf("unknown harness %q", s)
}

func Get(id string) (Adapter, error) {
	n, err := NormalizeHarnessID(id)
	if err != nil {
		return nil, err
	}
	return builtins[n], nil
}

func AllIDs() []string {
	return []string{"claude", "codex", "cursor", "antigravity"}
}
