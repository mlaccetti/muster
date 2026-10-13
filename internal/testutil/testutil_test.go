package testutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRepoRoot(t *testing.T) {
	root := RepoRoot(t)
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatal(err)
	}
}

func TestStubEnvHermetic(t *testing.T) {
	home := t.TempDir()
	bin, data := PrepareStubHarnessDirs(t)
	env := StubEnv(t, home, bin, data)
	joined := strings.Join(env, "\n")
	for _, needle := range []string{
		"MUSTER_ADAPTER_ROOT=" + bin,
		"MUSTER_STUB_DATA=" + data,
		"MUSTER_TEST=1",
		"MUSTER_SKIP_PROBE=1",
		"PATH=/usr/bin:/bin:/usr/sbin:/sbin",
	} {
		if !strings.Contains(joined, needle) {
			t.Fatalf("missing %q in env", needle)
		}
	}
	for _, name := range []string{"claude", "codex", "cursor-agent", "agy"} {
		if _, err := os.Stat(filepath.Join(bin, name)); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
}
