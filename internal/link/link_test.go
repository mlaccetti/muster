package link

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mlaccetti/muster/internal/kit"
)

func TestLinkClaudeHarnessWide(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("MUSTER_HOME", filepath.Join(home, "muster"))

	kitDir, err := kit.Install()
	if err != nil {
		t.Fatal(err)
	}

	res, err := Run(Options{
		Harnesses: []string{"claude"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Harnesses) != 1 || res.Harnesses[0].ID != "claude" {
		t.Fatalf("res %+v", res)
	}
	settingsPath := filepath.Join(home, ".claude", "settings.json")
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
	orchAgent := filepath.Join(home, ".claude", "agents", orchestratorAgentName+".md")
	target, err := filepath.EvalSymlinks(orchAgent)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(kitDir, "rendered", "claude", "muster-orchestrator.md")
	if target != want {
		t.Fatalf("agent link %s != %s", target, want)
	}
	m, err := LoadManifest()
	if err != nil || len(m.Harnesses) != 1 {
		t.Fatalf("manifest %+v err %v", m, err)
	}
	if m.Harnesses[0].ID != "claude" || len(m.Harnesses[0].Targets) == 0 {
		t.Fatalf("manifest does not describe harness-wide targets: %+v", m.Harnesses[0])
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
	res, err := Run(Options{
		Harnesses: []string{"codex", "cursor", "antigravity"},
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
	for _, required := range []string{"# Muster orchestrator policy", "## Codex controls", "Use `/goal`"} {
		if !strings.Contains(string(globalAgents), required) {
			t.Errorf("Codex instructions missing composed section %q", required)
		}
	}
	cursorPlugin := filepath.Join(home, ".cursor", "plugins", "local", "muster-orchestrator")
	for _, path := range []string{
		filepath.Join(cursorPlugin, ".cursor-plugin", "plugin.json"),
		filepath.Join(cursorPlugin, "rules", "muster-orchestrator.mdc"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("Cursor plugin file %s: %v", path, err)
		}
	}
	cursorRule, err := os.ReadFile(filepath.Join(cursorPlugin, "rules", "muster-orchestrator.mdc"))
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"alwaysApply: true", "Use Cursor's native Plan mode", "Use a native plan or TODO list"} {
		if !strings.Contains(string(cursorRule), required) {
			t.Errorf("Cursor rule missing %q", required)
		}
	}
	if len(res.Skipped) != 0 {
		t.Fatalf("supported harnesses should not be skipped: %+v", res.Skipped)
	}
	manifest, err := LoadManifest()
	if err != nil {
		t.Fatal(err)
	}
	if len(manifest.Harnesses) != 3 {
		t.Fatalf("manifest should include only configured harnesses: %+v", manifest.Harnesses)
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

func TestRenderedOrchestratorsComposeSharedPolicyAndHarnessAdapter(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("MUSTER_HOME", filepath.Join(home, "muster"))
	kitDir, err := kit.Install()
	if err != nil {
		t.Fatal(err)
	}

	templates := map[string]string{
		"claude":      filepath.Join("link", "claude", "muster-orchestrator.md"),
		"codex":       filepath.Join("link", "codex", "muster-orchestrator.md"),
		"cursor":      filepath.Join("link", "cursor", "rules", "muster-orchestrator.mdc"),
		"antigravity": filepath.Join("link", "agy", "muster-orchestrator", "agent.md"),
	}
	for harness, template := range templates {
		t.Run(harness, func(t *testing.T) {
			source, err := os.ReadFile(filepath.Join(kitDir, template))
			if err != nil {
				t.Fatal(err)
			}
			rendered, err := renderOrchestrator(kitDir, harness, source)
			if err != nil {
				t.Fatal(err)
			}
			body := string(rendered)
			if !strings.Contains(body, "# Muster orchestrator policy") || strings.Contains(body, policyMarker) || strings.Contains(body, adapterMarker) {
				t.Fatalf("shared policy was not composed cleanly: %s", body)
			}
			if !strings.Contains(body, "## "+map[string]string{
				"claude":      "Claude Code controls",
				"codex":       "Codex controls",
				"cursor":      "Cursor controls",
				"antigravity": "Antigravity controls",
			}[harness]) {
				t.Fatalf("%s adapter missing from rendered prompt", harness)
			}
		})
	}
}

func TestLinkNeverWritesProjectFiles(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CODEX_HOME", filepath.Join(home, ".codex"))
	t.Setenv("MUSTER_HOME", filepath.Join(home, "muster"))
	project := t.TempDir()
	if err := os.WriteFile(filepath.Join(project, "AGENTS.md"), []byte("# Existing project guidance\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := kit.Install(); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(Options{Harnesses: []string{"codex", "cursor"}}); err != nil {
		t.Fatal(err)
	}
	guidance, err := os.ReadFile(filepath.Join(project, "AGENTS.md"))
	if err != nil || string(guidance) != "# Existing project guidance\n" {
		t.Fatalf("project guidance changed: %v %q", err, guidance)
	}
	for _, path := range []string{".codex", ".cursor", ".claude", ".agents", ".muster"} {
		if _, err := os.Lstat(filepath.Join(project, path)); !os.IsNotExist(err) {
			t.Fatalf("link wrote project path %s (err=%v)", path, err)
		}
	}
	manifest, err := LoadManifest()
	if err != nil {
		t.Fatal(err)
	}
	if issues := VerifyReports(manifest); len(issues) != 0 {
		t.Fatalf("linked Codex harness wiring failed verification: %v", issues)
	}
}

func TestLinkAllWithoutGlobalPaths(t *testing.T) {
	home := t.TempDir()
	t.Setenv("MUSTER_HOME", filepath.Join(home, "muster"))
	_, err := kit.Install()
	if err != nil {
		t.Fatal(err)
	}
	res, err := Run(Options{
		All:    true,
		DryRun: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Harnesses) != 4 || len(res.Skipped) != 0 {
		t.Fatalf("expected four configured harnesses with no skipped setup, got %+v", res)
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
	_ = os.MkdirAll(filepath.Join(kitDir, "policy", "adapters"), 0o755)
	_ = os.WriteFile(filepath.Join(kitDir, "link", "claude", orchestratorAgentName+".md"), []byte("---\nname: x\n---\n<!-- MUSTER_POLICY -->\n<!-- MUSTER_ADAPTER -->\n"), 0o644)
	_ = os.WriteFile(filepath.Join(kitDir, "link", "claude", workerAgentName+".md"), []byte("---\nname: y\n---\n"), 0o644)
	_ = os.WriteFile(filepath.Join(kitDir, "policy", "orchestrator.md"), []byte("# Shared policy\n"), 0o644)
	_ = os.WriteFile(filepath.Join(kitDir, "policy", "adapters", "claude.md"), []byte("# Claude adapter\n"), 0o644)
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
