package link

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mlaccetti/muster/internal/kit"
)

func TestLinkProjectClaudeOnly(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("MUSTER_HOME", filepath.Join(home, "muster"))
	project := t.TempDir()

	kitDir, err := kit.Install()
	if err != nil {
		t.Fatal(err)
	}

	res, err := Run(Options{
		ProjectDir: project,
		Harnesses:  []string{"claude"},
		Global:     false,
		Project:    true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Harnesses) != 1 || res.Harnesses[0].ID != "claude" {
		t.Fatalf("res %+v", res)
	}
	settingsPath := filepath.Join(project, ".claude", "settings.json")
	raw, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatal(err)
	}
	var settings map[string]interface{}
	if err := json.Unmarshal(raw, &settings); err != nil {
		t.Fatal(err)
	}
	if settings["agent"] != orchestratorAgentName {
		t.Fatalf("settings agent: %+v", settings)
	}
	orchAgent := filepath.Join(project, ".claude", "agents", orchestratorAgentName+".md")
	target, err := filepath.EvalSymlinks(orchAgent)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(kitDir, "link", "claude", orchestratorAgentName+".md")
	if target != want {
		t.Fatalf("agent link %s != %s", target, want)
	}
	if _, err := os.Stat(filepath.Join(project, ".muster", "orchestrator.path")); err != nil {
		t.Fatal(err)
	}
	m, err := LoadManifest()
	if err != nil || len(m.Harnesses) != 1 {
		t.Fatalf("manifest %+v err %v", m, err)
	}
}

func TestLinkGlobalCodexCursorAgy(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("MUSTER_HOME", filepath.Join(home, "muster"))
	t.Setenv("CODEX_HOME", filepath.Join(home, ".codex"))
	_ = os.MkdirAll(filepath.Join(home, ".gemini", "antigravity-cli"), 0o755)

	kitDir, err := kit.Install()
	if err != nil {
		t.Fatal(err)
	}
	_, err = Run(Options{
		Harnesses: []string{"codex", "cursor", "antigravity"},
		Global:    true,
		Project:   false,
	})
	if err != nil {
		t.Fatal(err)
	}
	codexProfile := filepath.Join(home, ".codex", codexProfileName+".config.toml")
	if target, err := filepath.EvalSymlinks(codexProfile); err != nil || target != filepath.Join(kitDir, "link", "codex", codexProfileName+".config.toml") {
		t.Fatalf("codex profile: %v %s", err, target)
	}
	cfg, err := os.ReadFile(filepath.Join(home, ".codex", "config.toml"))
	if err != nil || !strings.Contains(string(cfg), tomlMarkerBegin) {
		t.Fatalf("codex config: %v %s", err, cfg)
	}
	cursorRule := filepath.Join(home, ".cursor", "rules", "muster-orchestrator.mdc")
	cursorTarget, err := filepath.EvalSymlinks(cursorRule)
	if err != nil {
		t.Fatal(err)
	}
	if cursorTarget != filepath.Join(kitDir, "link", "cursor", "muster-orchestrator.mdc") {
		t.Fatalf("cursor rule: %s", cursorTarget)
	}
	agyAgent := filepath.Join(home, ".gemini", "config", "agents", orchestratorAgentName, "agent.md")
	if _, err := os.Stat(agyAgent); err != nil {
		t.Fatal(err)
	}
}

func TestLinkAllWithoutGlobalPaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("MUSTER_HOME", filepath.Join(home, "muster"))
	project := t.TempDir()
	_, err := kit.Install()
	if err != nil {
		t.Fatal(err)
	}
	res, err := Run(Options{
		ProjectDir: project,
		All:        true,
		Global:     false,
		Project:    true,
		DryRun:     true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Harnesses) != 4 {
		t.Fatalf("expected 4 harnesses, got %d", len(res.Harnesses))
	}
}

func TestClaudeSettingsRestoreOnUnlink(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	claudeDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	settings := `{"agent":"orchestrator","model":"opus"}` + "\n"
	if err := os.WriteFile(filepath.Join(claudeDir, "settings.json"), []byte(settings), 0o644); err != nil {
		t.Fatal(err)
	}
	kitDir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(kitDir, "link", "claude"), 0o755)
	_ = os.WriteFile(filepath.Join(kitDir, "link", "claude", orchestratorAgentName+".md"), []byte("---\nname: x\n---\n"), 0o644)
	_ = os.WriteFile(filepath.Join(kitDir, "link", "claude", workerAgentName+".md"), []byte("---\nname: y\n---\n"), 0o644)
	if _, err := linkClaudeAgentWiring(claudeDir, kitDir, false); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(filepath.Join(claudeDir, "settings.json"))
	if !strings.Contains(string(raw), orchestratorAgentName) {
		t.Fatalf("expected linked agent: %s", raw)
	}
	if err := unlinkClaudeAgentWiring(claudeDir, false); err != nil {
		t.Fatal(err)
	}
	raw, _ = os.ReadFile(filepath.Join(claudeDir, "settings.json"))
	if !strings.Contains(string(raw), `"agent": "orchestrator"`) {
		t.Fatalf("expected restored agent: %s", raw)
	}
}
