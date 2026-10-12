package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

func Path(t *testing.T, rel string) string {
	return filepath.Join(RepoRoot(t), rel)
}

func RepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for dir := wd; ; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		if dir == filepath.Dir(dir) {
			t.Fatal("go.mod not found")
		}
	}
}
