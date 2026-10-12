package cli

import (
	"github.com/mlaccetti/muster/internal/kit"
	"github.com/mlaccetti/muster/internal/paths"
	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install policy kit to MUSTER_HOME",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := kit.Install()
			if err != nil {
				return err
			}
			home, _ := paths.Home()
			writef(cmd.OutOrStdout(), "Installed kit to %s\n", dir)
			writef(cmd.OutOrStdout(), "Orchestrator policy: %s/policy/orchestrator.md\n", dir)
			writef(cmd.OutOrStdout(), "Worker policy: %s/policy/worker.md\n", dir)
			_ = home
			writef(cmd.OutOrStdout(), "Next: muster link   # wire orchestrator + worker into harness configs\n")
			return nil
		},
	}
	return cmd
}
