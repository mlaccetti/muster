package kit

import (
	"strings"
	"testing"
)

func TestEmbedPolicy(t *testing.T) {
	for _, name := range []string{
		"policy/orchestrator.md",
		"policy/worker.md",
		"policy/adapters/claude.md",
		"policy/adapters/codex.md",
		"policy/adapters/antigravity.md",
		"policy/adapters/cursor.md",
		"link/codex/muster-orchestrator.md",
		"link/claude/muster-orchestrator.md",
		"link/agy/muster-orchestrator/agent.md",
		"link/cursor/plugin.json",
		"link/cursor/rules/muster-orchestrator.mdc",
	} {
		if _, err := Source.ReadFile(name); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
}

func TestOrchestrationPolicyRequiresNativeParallelEventDrivenWorkflow(t *testing.T) {
	orchestrator, err := Source.ReadFile("policy/orchestrator.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"Keep coordination small",
		"Start independent tasks concurrently",
		"gh-stacking",
		"lifecycle events or callbacks",
		"Do not poll, sleep",
		"repository's configured merge queue",
		"enqueuePullRequest",
		"expectedHeadOid",
		"exact reviewed PR head",
		"dequeuePullRequest",
		"Do not fall back to direct merge",
		"actual merge",
		"that exact commit",
		"post-merge CI",
	} {
		if !strings.Contains(string(orchestrator), required) {
			t.Errorf("orchestrator policy missing %q", required)
		}
	}
	if strings.Contains(string(orchestrator), "@mergify queue") {
		t.Error("orchestrator policy must not mandate a queue provider through a hardcoded bot comment")
	}
	worker, err := Source.ReadFile("policy/worker.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"ordinary harness conversation",
		"There is no required field list, status enum, or response schema",
	} {
		if !strings.Contains(string(worker), required) {
			t.Errorf("worker policy missing %q", required)
		}
	}
}
