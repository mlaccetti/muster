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
	codexWorker := filepath.Join(home, ".codex", "agents", workerAgentName+".toml")
	if target, err := filepath.EvalSymlinks(codexWorker); err != nil || target != filepath.Join(kitDir, "link", "codex", workerAgentName+".config.toml") {
		t.Fatalf("codex worker agent: %v %s", err, target)
	}
	workerConfig, err := os.ReadFile(filepath.Join(kitDir, "link", "codex", workerAgentName+".config.toml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{`name = "muster-worker"`, `model = "gpt-6-luna"`, `model_reasoning_effort = "medium"`, `developer_instructions = """`} {
		if !strings.Contains(string(workerConfig), required) {
			t.Errorf("Codex worker config missing %q: %s", required, workerConfig)
		}
	}
	globalAgents, err := os.ReadFile(filepath.Join(home, ".codex", "AGENTS.md"))
	if err != nil || !strings.Contains(string(globalAgents), markerStart) {
		t.Fatalf("codex global instructions: %v %s", err, globalAgents)
	}
	if _, err := os.Stat(filepath.Join(home, ".cursor", "rules", "muster-orchestrator.mdc")); !os.IsNotExist(err) {
		t.Fatalf("Cursor global rules are not a supported filesystem target; got err %v", err)
	}
	agyAgent := filepath.Join(home, ".gemini", "config", "agents", orchestratorAgentName, "agent.md")
	if _, err := os.Stat(agyAgent); err != nil {
		t.Fatal(err)
	}
	agySettings := filepath.Join(home, ".gemini", "antigravity-cli", "settings.json")
	if _, err := os.Stat(agySettings); err != nil {
		t.Fatal(err)
	}
}

func TestLinkProjectCodexCursorDefaultInstructions(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CODEX_HOME", filepath.Join(home, ".codex"))
	t.Setenv("MUSTER_HOME", filepath.Join(home, "muster"))
	project := t.TempDir()
	if err := os.WriteFile(filepath.Join(project, "AGENTS.md"), []byte("# Existing project guidance\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	kitDir, err := kit.Install()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Run(Options{ProjectDir: project, Harnesses: []string{"codex", "cursor"}, Global: false, Project: true}); err != nil {
		t.Fatal(err)
	}
	projectInstructions, err := os.ReadFile(filepath.Join(project, "AGENTS.md"))
	if err != nil || !strings.Contains(string(projectInstructions), markerStart) || !strings.Contains(string(projectInstructions), "Existing project guidance") {
		t.Fatalf("project instructions did not preserve and add managed block: %v %s", err, projectInstructions)
	}
	workerAgent := filepath.Join(project, ".codex", "agents", workerAgentName+".toml")
	if target, err := filepath.EvalSymlinks(workerAgent); err != nil || target != filepath.Join(kitDir, "link", "codex", workerAgentName+".config.toml") {
		t.Fatalf("project worker agent: %v %s", err, target)
	}
	cursorRule := filepath.Join(project, ".cursor", "rules", "muster-orchestrator.mdc")
	if target, err := filepath.EvalSymlinks(cursorRule); err != nil || target != filepath.Join(kitDir, "link", "cursor", "muster-orchestrator.mdc") {
		t.Fatalf("project Cursor rule: %v %s", err, target)
	}
	manifest, err := LoadManifest()
	if err != nil {
		t.Fatal(err)
	}
	if issues := VerifyReports(manifest); len(issues) != 0 {
		t.Fatalf("linked Codex and Cursor project wiring failed verification: %v", issues)
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
