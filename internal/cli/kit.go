package cli

import (
	"github.com/mlaccetti/muster/internal/kit"
	"github.com/spf13/cobra"
)

func newKitCmd() *cobra.Command {
	pathCmd := &cobra.Command{
		Use:   "path",
		Short: "Print active kit directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := kit.ActiveDir()
			if err != nil {
				return err
			}
			writelnf(cmd.OutOrStdout(), dir)
			return nil
		},
	}
	kitCmd := &cobra.Command{
		Use:   "kit",
		Short: "Kit installation helpers",
	}
	kitCmd.AddCommand(pathCmd)
	return kitCmd
}
