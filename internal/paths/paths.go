package paths

import (
	"os"
	"path/filepath"
)

const KitVersion = "v0.2.0"

func Home() (string, error) {
	if v := os.Getenv("MUSTER_HOME"); v != "" {
		return v, nil
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "muster"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "muster"), nil
}

// KitDir returns the active kit directory (kit/current when present, else kit/<KitVersion>).
func KitDir() (string, error) {
	home, err := Home()
	if err != nil {
		return "", err
	}
	current := filepath.Join(home, "kit", "current")
	if resolved, err := filepath.EvalSymlinks(current); err == nil {
		return resolved, nil
	}
	return filepath.Join(home, "kit", KitVersion), nil
}

func KitVersionedDir() (string, error) {
	home, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "kit", KitVersion), nil
}

func KitCurrentLink() (string, error) {
	home, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "kit", "current"), nil
}
