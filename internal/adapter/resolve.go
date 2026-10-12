package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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
	if root := os.Getenv("MUSTER_ADAPTER_ROOT"); root != "" {
		name := binaryName(id)
		p := filepath.Join(root, name)
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p, nil
		}
		return "", fmt.Errorf("stub binary %s not found under MUSTER_ADAPTER_ROOT", name)
	}
	name := binaryName(id)
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("%s not found on PATH", name)
	}
	return path, nil
}

func probeVersion(bin string) (string, error) {
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
