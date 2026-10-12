package cli

import (
	"github.com/spf13/cobra"
)

func NewRoot() *cobra.Command {
	root := &cobra.Command{
		Use:   "muster",
		Short: "Harness-neutral briefs, workers, and routing for coding agents",
		Long:  "muster validates kit data, renders briefs, runs workers through adapters, and walks harness pools.",
		Example: `  muster install
  muster render brief --template implement --input unit.json --out brief.txt
  muster worker --harness cursor --brief brief.txt --report report.json --sync`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(
		newInstallCmd(),
		newLinkCmd(),
		newDoctorCmd(),
		newKitCmd(),
		newValidateCmd(),
		newRenderCmd(),
		newWorkerCmd(),
		newDispatchCmd(),
		newPassCmd(),
	)
	setupCompletions(root)
	attachReleaseCheck(root)
	return root
}
