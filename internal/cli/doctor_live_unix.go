//go:build darwin || freebsd || linux || netbsd || openbsd

package cli

import (
	"os/exec"
	"syscall"
)

func configureLiveProcess(cmd *exec.Cmd) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
	return nil
}
