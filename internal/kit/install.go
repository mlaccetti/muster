package kit

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mlaccetti/muster/internal/paths"
)

func closeIgnore(c io.Closer) {
	_ = c.Close()
}

func Install() (string, error) {
	dest, err := paths.KitVersionedDir()
	if err != nil {
		return "", err
	}
	if err := copyTree(SourceFS(), dest); err != nil {
		return "", err
	}
	home, err := paths.Home()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Join(home, "kit"), 0o755); err != nil {
		return "", err
	}
	if err := setCurrentKitLink(home, paths.KitVersion); err != nil {
		return "", err
	}
	if err := os.MkdirAll(home, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(home, "installed-kit-version.txt"), []byte(paths.KitVersion+"\n"), 0o644); err != nil {
		return "", err
	}
	active, err := paths.KitDir()
	if err != nil {
		return dest, err
	}
	return active, nil
}

func setCurrentKitLink(home, version string) error {
	link := filepath.Join(home, "kit", "current")
	_ = os.Remove(link)
	return os.Symlink(version, link)
}

func ActiveDir() (string, error) {
	dir, err := paths.KitDir()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(dir); err != nil {
		return "", fmt.Errorf("kit not installed at %s (run muster install)", dir)
	}
	return dir, nil
}

func copyTree(src fs.FS, dest string) error {
	return fs.WalkDir(src, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		target := filepath.Join(dest, name)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFSFile(src, name, target)
	})
}

func copyFSFile(src fs.FS, name, dest string) error {
	in, err := src.Open(name)
	if err != nil {
		return err
	}
	defer closeIgnore(in)
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer closeIgnore(out)
	_, err = io.Copy(out, in)
	return err
}
