package harnesscli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var allIDs = []string{"claude", "codex", "cursor", "antigravity"}

func AllIDs() []string {
	return append([]string{}, allIDs...)
}

func NormalizeHarnessID(s string) (string, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "agy" {
		return "antigravity", nil
	}
	for _, id := range allIDs {
		if s == id {
			return id, nil
		}
	}
	return "", fmt.Errorf("unknown harness %q", s)
}

func binaryName(id string) string {
	switch id {
	case "cursor":
		return "cursor-agent"
	case "antigravity":
		return "agy"
	default:
		return id
	}
}

func ResolveBinary(id string) (string, error) {
	n, err := NormalizeHarnessID(id)
	if err != nil {
		return "", err
	}
	if root := os.Getenv("MUSTER_ADAPTER_ROOT"); root != "" {
		name := binaryName(n)
		p := filepath.Join(root, name)
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p, nil
		}
		return "", fmt.Errorf("stub binary %s not found under MUSTER_ADAPTER_ROOT", name)
	}
	name := binaryName(n)
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("%s not found on PATH", name)
	}
	return path, nil
}

func ProbeVersion(bin string) (string, error) {
	out, err := exec.Command(bin, "--version").Output()
	if err != nil {
		out, err = exec.Command(bin, "-h").Output()
		if err != nil {
			return "", err
		}
	}
	line := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	return line, nil
}

func ResolveCLI(id string) (absPath string, versionLine string, err error) {
	path, err := ResolveBinary(id)
	if err != nil {
		return "", "", err
	}
	ver, err := ProbeVersion(path)
	if err != nil {
		ver = "unknown"
	}
	return path, ver, nil
}
