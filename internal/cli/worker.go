package cli

import (
	"context"
	"os"

	"github.com/mlaccetti/muster/internal/routing"
	"github.com/mlaccetti/muster/internal/service"
	"github.com/spf13/cobra"
)

func newWorkerCmd() *cobra.Command {
	var harness, brief, reportPath, workdir, tier string
	sync := true
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Run one harness worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			if workdir == "" {
				workdir, _ = os.Getwd()
			}
			if os.Getenv("MUSTER_TEST") == "1" {
				sync = true
			}
			res, err := service.RunWorker(context.Background(), service.WorkerRequest{
				Harness: harness, BriefPath: brief, ReportPath: reportPath, WorkDir: workdir,
				Sync: sync,
			})
			if err != nil {
				return err
			}
			_ = tier
			if err := routing.Emit(cmd.OutOrStdout(), res.Line); err != nil {
				return err
			}
			if res.ExitCode != 0 {
				return exitErr(res.ExitCode)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&harness, "harness", "cursor", "Harness id")
	cmd.Flags().StringVar(&brief, "brief", "", "Brief file path")
	cmd.Flags().StringVar(&reportPath, "report", "", "Report output path")
	cmd.Flags().StringVar(&workdir, "workdir", "", "Working directory")
	cmd.Flags().StringVar(&tier, "tier", "", "Informational tier label")
	cmd.Flags().BoolVar(&sync, "sync", true, "Run in foreground")
	_ = cmd.MarkFlagRequired("brief")
	_ = cmd.MarkFlagRequired("report")
	return cmd
}
