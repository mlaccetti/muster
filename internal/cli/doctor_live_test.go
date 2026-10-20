package cli

import (
	"testing"

	"github.com/mlaccetti/muster/internal/link"
)

func TestValidateCodexLiveRoutingIgnoresOtherHarnessesButRequiresCodex(t *testing.T) {
	t.Setenv("MUSTER_HOME", t.TempDir())
	t.Setenv("CODEX_HOME", t.TempDir())
	if err := link.SaveManifest(&link.Manifest{Harnesses: []link.HarnessLink{{ID: "cursor"}}}); err != nil {
		t.Fatal(err)
	}
	if err := validateCodexLiveRouting(); err == nil {
		t.Fatal("manifest without Codex routing unexpectedly passed")
	}

	if err := link.SaveManifest(&link.Manifest{Harnesses: []link.HarnessLink{{ID: "codex"}}}); err != nil {
		t.Fatal(err)
	}
	if err := validateCodexLiveRouting(); err == nil {
		t.Fatal("Codex manifest without the worker target unexpectedly passed")
	}
}
