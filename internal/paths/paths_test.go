package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHomeMUSTERHome(t *testing.T) {
	t.Setenv("MUSTER_HOME", "/custom/muster")
	t.Setenv("XDG_CONFIG_HOME", "/xdg")
	got, err := Home()
	if err != nil || got != "/custom/muster" {
		t.Fatalf("Home()=%q err=%v", got, err)
	}
}

func TestHomeXDG(t *testing.T) {
	t.Setenv("MUSTER_HOME", "")
	t.Setenv("XDG_CONFIG_HOME", "/xdg/config")
	got, err := Home()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join("/xdg/config", "muster")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestKitDirAndPoolsFile(t *testing.T) {
	t.Setenv("MUSTER_HOME", "/m")
	kit, err := KitDir()
	if err != nil {
		t.Fatal(err)
	}
	if kit != filepath.Join("/m", "kit", KitVersion) {
		t.Fatalf("KitDir %q", kit)
	}
	pools, err := PoolsFile()
	if err != nil || pools != filepath.Join("/m", "pools.json") {
		t.Fatalf("PoolsFile %q err=%v", pools, err)
	}
}

func TestKitCurrentSymlink(t *testing.T) {
	home := t.TempDir()
	t.Setenv("MUSTER_HOME", home)
	versioned := filepath.Join(home, "kit", KitVersion)
	if err := os.MkdirAll(versioned, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(home, "kit", "current")
	if err := os.Symlink(KitVersion, link); err != nil {
		t.Fatal(err)
	}
	got, err := KitDir()
	if err != nil {
		t.Fatal(err)
	}
	want, err := filepath.EvalSymlinks(versioned)
	if err != nil {
		t.Fatal(err)
	}
	got, err = filepath.EvalSymlinks(got)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("KitDir %q want %q", got, want)
	}
}

func TestPassDir(t *testing.T) {
	t.Setenv("MUSTER_HOME", "/m")
	dir, err := PassDir("pass-1")
	if err != nil {
		t.Fatal(err)
	}
	if dir != filepath.Join("/m", "passes", "pass-1") {
		t.Fatalf("PassDir %q", dir)
	}
}
