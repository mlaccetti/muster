package testutil

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// PrepareStubHarnessDirs builds temp stub binaries and data (no process env mutation).
func PrepareStubHarnessDirs(t *testing.T) (binDir, dataDir string) {
	t.Helper()
	root := RepoRoot(t)
	stubRoot := filepath.Join(root, "testdata", "stubs")
	cli := filepath.Join(stubRoot, "cli")
	_ = os.Chmod(cli, 0o755)
	bin := t.TempDir()
	names := []string{"claude", "codex", "cursor-agent", "agy"}
	for _, name := range names {
		dest := filepath.Join(bin, name)
		copyFile(t, cli, dest)
		_ = os.Chmod(dest, 0o755)
	}
	data := t.TempDir()
	copyDir(t, filepath.Join(stubRoot, "data"), data)
	return bin, data
}

// StubEnv returns environment variables for hermetic harness runs (subprocess or in-process).
func StubEnv(t *testing.T, home string, binDir, dataDir string) []string {
	t.Helper()
	musterHome := filepath.Join(home, "muster-home")
	env := dropKeys(os.Environ(),
		"MUSTER_ADAPTER_ROOT", "MUSTER_STUB_DATA", "MUSTER_TEST", "MUSTER_SKIP_PROBE",
		"MUSTER_HOME", "XDG_CONFIG_HOME", "HOME",
	)
	env = append(env,
		"HOME="+home,
		"MUSTER_HOME="+musterHome,
		"XDG_CONFIG_HOME="+filepath.Join(home, "xdg"),
		"MUSTER_ADAPTER_ROOT="+binDir,
		"MUSTER_STUB_DATA="+dataDir,
		"MUSTER_TEST=1",
		"MUSTER_SKIP_PROBE=1",
		"PATH=/usr/bin:/bin:/usr/sbin:/sbin",
	)
	return env
}

// SetupStubs prepares stub dirs and sets env for in-process tests in this package.
func SetupStubs(t *testing.T) func() {
	t.Helper()
	home := t.TempDir()
	bin, data := PrepareStubHarnessDirs(t)
	for _, e := range StubEnv(t, home, bin, data) {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			t.Setenv(parts[0], parts[1])
		}
	}
	return func() {}
}

func StubDataDir(t *testing.T) string {
	return os.Getenv("MUSTER_STUB_DATA")
}

func dropKeys(env []string, keys ...string) []string {
	skip := map[string]bool{}
	for _, k := range keys {
		skip[k] = true
	}
	var out []string
	for _, e := range env {
		k := e
		if i := strings.IndexByte(e, '='); i >= 0 {
			k = e[:i]
		}
		if skip[k] {
			continue
		}
		out = append(out, e)
	}
	return out
}

func copyDir(t *testing.T, src, dest string) {
	t.Helper()
	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		in := filepath.Join(src, e.Name())
		out := filepath.Join(dest, e.Name())
		if e.IsDir() {
			_ = os.MkdirAll(out, 0o755)
			copyDir(t, in, out)
			continue
		}
		copyFile(t, in, out)
	}
}

func copyFile(t *testing.T, src, dest string) {
	t.Helper()
	in, err := os.Open(src)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = in.Close() }()
	out, err := os.Create(dest)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = out.Close() }()
	if _, err := io.Copy(out, in); err != nil {
		t.Fatal(err)
	}
}
