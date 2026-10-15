package kit

import (
	"strings"
	"testing"
)

func TestEmbedPolicy(t *testing.T) {
	for _, name := range []string{
		"policy/orchestrator.md",
		"policy/worker.md",
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
		"@mergify queue",
		"post-merge CI",
	} {
		if !strings.Contains(string(orchestrator), required) {
			t.Errorf("orchestrator policy missing %q", required)
		}
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
