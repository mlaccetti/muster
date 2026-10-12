package cli

import (
	"io"
)

// RunForTest executes the CLI in-process (Cobra only, no Fang). For integration tests.
func RunForTest(args []string, stdout, stderr io.Writer) int {
	root := NewRoot()
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs(args)
	if err := root.Execute(); err != nil {
		return ExitCode(err)
	}
	return 0
}
