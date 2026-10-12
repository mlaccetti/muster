package adapter

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"
)

func loadStream(t *testing.T, name string) []NormalizedEvent {
	t.Helper()
	path := filepath.Join("..", "..", "testdata", "streams", name+".jsonl")
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	a := &claudeAdapter{}
	var events []NormalizedEvent
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 || line[0] != '{' {
			continue
		}
		ev, err := a.ParseStream(line)
		if err != nil {
			t.Fatal(err)
		}
		events = append(events, ev)
	}
	return events
}

func TestClassifyClaudeOK(t *testing.T) {
	ev := loadStream(t, "claude-ok")
	v := classifyClaude(ev, 0)
	if v.Kind != "ok" {
		t.Fatalf("got %q", v.Kind)
	}
}

func TestClassifyClaudeLimit(t *testing.T) {
	ev := loadStream(t, "claude-limit-result")
	v := classifyClaude(ev, 0)
	if v.Kind != "limit" {
		t.Fatalf("got %q", v.Kind)
	}
}

func TestClassifyCodexOK(t *testing.T) {
	ev := loadStream(t, "codex-ok")
	v := classifyCodex(ev, 0)
	if v.Kind != "ok" {
		t.Fatalf("got %q", v.Kind)
	}
}

func TestClassifyCursorOK(t *testing.T) {
	ev := loadStream(t, "cursor-ok")
	v := classifyCursor(ev, 0)
	if v.Kind != "ok" {
		t.Fatalf("got %q", v.Kind)
	}
}
