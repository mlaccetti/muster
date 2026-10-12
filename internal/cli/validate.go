package cli

import (
	"github.com/mlaccetti/muster/internal/report"
	"github.com/mlaccetti/muster/internal/schema"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate JSON documents against kit schemas",
	}
	for _, kind := range []string{"brief", "report", "routing"} {
		k := kind
		sub := &cobra.Command{
			Use:   k,
			Short: "Validate " + k,
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				path := args[0]
				if k == "report" {
					norm, err := report.NormalizeFile(path)
					if err != nil {
						return err
					}
					return schema.ValidateBytes("report", norm)
				}
				return schema.ValidateFile(k, path)
			},
		}
		cmd.AddCommand(sub)
	}
	return cmd
}
