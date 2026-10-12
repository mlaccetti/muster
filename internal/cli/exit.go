package cli

import "fmt"

type exitError struct {
	code int
}

func (e *exitError) Error() string {
	return fmt.Sprintf("exit %d", e.code)
}

func exitErr(code int) error {
	return &exitError{code: code}
}

func ExitCode(err error) int {
	if err == nil {
		return 0
	}
	if ee, ok := err.(*exitError); ok {
		return ee.code
	}
	return 1
}
