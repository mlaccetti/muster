package cli

import (
	"github.com/spf13/cobra"
)

func NewRoot() *cobra.Command {
	root := &cobra.Command{
		Use:   "muster",
		Short: "Consistent orchestrator and worker behavior across coding harnesses",
		Long:  "muster installs a versioned policy kit and links it into Claude Code, Codex, Cursor Agent, and Antigravity (agy).",
		Example: `  muster install
  muster link
  muster doctor`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(
		newInstallCmd(),
		newLinkCmd(),
		newDoctorCmd(),
		newKitCmd(),
	)
	setupCompletions(root)
	attachReleaseCheck(root)
	return root
}
