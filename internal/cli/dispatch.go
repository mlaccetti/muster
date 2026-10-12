package cli

import (
	"context"
	"os"

	"github.com/mlaccetti/muster/internal/paths"
	"github.com/mlaccetti/muster/internal/routing"
	"github.com/mlaccetti/muster/internal/service"
	"github.com/spf13/cobra"
)

func newDispatchCmd() *cobra.Command {
	var tier, brief, reportPath, workdir, poolsPath string
	var avoidFamilies []string
	sync := true
	cmd := &cobra.Command{
		Use:   "dispatch",
		Short: "Walk a tier pool and run the first usable worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			if workdir == "" {
				workdir, _ = os.Getwd()
			}
			if poolsPath == "" {
				poolsPath, _ = paths.PoolsFile()
			}
			if os.Getenv("MUSTER_TEST") == "1" {
				sync = true
			}
			res, err := service.RunDispatch(context.Background(), service.DispatchRequest{
				Tier: tier, BriefPath: brief, ReportPath: reportPath, WorkDir: workdir,
				PoolsPath: poolsPath, AvoidFamilies: avoidFamilies, Sync: sync,
			})
			if err != nil {
				return err
			}
			if err := routing.Emit(cmd.OutOrStdout(), res.Line); err != nil {
				return err
			}
			if res.ExitCode != 0 {
				return exitErr(res.ExitCode)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&tier, "tier", "implement", "Pool tier")
	cmd.Flags().StringVar(&brief, "brief", "", "Brief file path")
	cmd.Flags().StringVar(&reportPath, "report", "", "Report output path")
	cmd.Flags().StringVar(&workdir, "workdir", "", "Working directory")
	cmd.Flags().StringArrayVar(&avoidFamilies, "avoid-family", nil, "Skip harness families (repeatable)")
	cmd.Flags().StringVar(&poolsPath, "pools", "", "Pools JSON path")
	cmd.Flags().BoolVar(&sync, "sync", true, "Run in foreground")
	_ = cmd.MarkFlagRequired("brief")
	_ = cmd.MarkFlagRequired("report")
	return cmd
}
