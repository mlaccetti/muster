package link

import (
	"io"
	"os"
	"path/filepath"
)

func kitLinkFile(kitDir string, parts ...string) string {
	return filepath.Join(kitDir, "link", filepath.Join(parts...))
}

func linkKitFile(dest, kitDir string, kitParts ...string) error {
	src := kitLinkFile(kitDir, kitParts...)
	if _, err := os.Stat(src); err != nil {
		return err
	}
	return replaceWithSymlink(dest, src)
}

func replaceWithSymlink(dest, src string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	_ = os.Remove(dest)
	if err := os.Symlink(src, dest); err != nil {
		return copyFile(dest, src)
	}
	return nil
}

func removeLink(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func copyFile(dest, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()
	_, err = io.Copy(out, in)
	return err
}

func linkKitFileDryRun(dest string, dryRun bool, kitDir string, kitParts ...string) error {
	if dryRun {
		return nil
	}
	return linkKitFile(dest, kitDir, kitParts...)
}
