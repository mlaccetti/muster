//go:build !darwin && !freebsd && !linux && !netbsd && !openbsd

package cli

import (
	"errors"
	"os/exec"
)

func configureLiveProcess(cmd *exec.Cmd) error {
	return errors.New("live doctor is unsupported on this operating system")
}
