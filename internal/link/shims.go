package link

import (
	"os"
	"path/filepath"

	"github.com/mlaccetti/muster/internal/paths"
)

func installLinkShims(dryRun bool) error {
	home, err := paths.Home()
	if err != nil {
		return err
	}
	bin := filepath.Join(home, "bin")
	if dryRun {
		return nil
	}
	if err := os.MkdirAll(bin, 0o755); err != nil {
		return err
	}
	agyShim := "#!/bin/sh\nexec agy --agent " + orchestratorAgentName + " \"$@\"\n"
	return writeExecutable(filepath.Join(bin, "agy-muster"), agyShim)
}

func removeLinkShims(dryRun bool) error {
	if dryRun {
		return nil
	}
	home, err := paths.Home()
	if err != nil {
		return err
	}
	_ = os.Remove(filepath.Join(home, "bin", "codex-muster"))
	_ = os.Remove(filepath.Join(home, "bin", "agy-muster"))
	return nil
}

func writeExecutable(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		return err
	}
	return nil
}
