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

func linkCodexWorkerAgent(dest, kitDir string, dryRun bool) error {
	if dryRun {
		return nil
	}
	source := kitLinkFile(kitDir, "codex", workerAgentName+".config.toml")
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	return writeFileAtomic(dest, data, 0o644)
}

func writeFileAtomic(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".muster-link-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer func() { _ = os.Remove(tmpPath) }()
	if err := tmp.Chmod(mode); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}
