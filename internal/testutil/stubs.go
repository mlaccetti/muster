package testutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// PrepareStubHarnessDirs builds temp stub binaries (no process env mutation).
func PrepareStubHarnessDirs(t *testing.T) (binDir, dataDir string) {
	t.Helper()
	bin := t.TempDir()
	names := []string{"claude", "codex", "cursor-agent", "agy"}
	for _, name := range names {
		dest := filepath.Join(bin, name)
		_ = os.WriteFile(dest, []byte("#!/bin/sh\nexit 0\n"), 0o755)
	}
	return bin, t.TempDir()
}

// StubEnv returns environment variables for hermetic CLI runs (subprocess or in-process).
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
