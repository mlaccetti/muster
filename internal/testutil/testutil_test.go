package testutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRepoRootAndPath(t *testing.T) {
	root := RepoRoot(t)
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err != nil {
		t.Fatal(err)
	}
	brief := Path(t, "testdata/briefs/minimal.json")
	if _, err := os.Stat(brief); err != nil {
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

func TestLastRoutingLine(t *testing.T) {
	out := "noise\n{\"verdict\":\"ok\",\"harness\":\"claude\",\"family\":\"anthropic\",\"model\":\"\",\"report\":\"/r.json\",\"reason_code\":\"none\",\"reset\":0}\n"
	line := LastRoutingLine(t, out)
	if line.Verdict != "ok" || line.Harness == nil || *line.Harness != "claude" {
		t.Fatalf("line %+v", line)
	}
}
