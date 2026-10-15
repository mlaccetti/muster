package cli

import (
	"strings"

	"github.com/mlaccetti/muster/internal/link"
	"github.com/spf13/cobra"
)

func newLinkCmd() *cobra.Command {
	var (
		harnesses []string
		all       bool
		dryRun    bool
		unlink    bool
	)
	cmd := &cobra.Command{
		Use:   "link",
		Short: "Wire muster orchestrator and worker policies into harness configs",
		Long: `Link installs managed snippets that point every harness at the same kit policies.

By default, muster auto-detects harnesses whose CLIs are on PATH and links their
harness-wide (user) configuration. Use --harness to pin specific IDs, or --all to
link every supported harness even when not installed.`,
		Example: `  muster install
  muster link
  muster link --harness claude --harness codex
  muster link --all --dry-run
  muster link --unlink`,
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := link.Run(link.Options{
				Harnesses: harnesses,
				All:       all,
				DryRun:    dryRun,
				Unlink:    unlink,
			})
			if err != nil {
				return err
			}
			if unlink {
				writelnf(cmd.OutOrStdout(), "muster: unlinked")
				return nil
			}
			writef(cmd.OutOrStdout(), "kit: %s\n", res.KitDir)
			for _, h := range res.Harnesses {
				writef(cmd.OutOrStdout(), "linked %s:\n", h.ID)
				for _, t := range h.Targets {
					writef(cmd.OutOrStdout(), "  %s\n", t)
				}
			}
			if len(res.Skipped) > 0 {
				writef(cmd.OutOrStdout(), "skipped: %s\n", strings.Join(res.Skipped, "; "))
			}
			if dryRun {
				writelnf(cmd.OutOrStdout(), "(dry run)")
			}
			return nil
		},
	}
	cmd.Flags().StringArrayVar(&harnesses, "harness", nil, "Harness id (repeatable); default: auto-detect on PATH")
	cmd.Flags().BoolVar(&all, "all", false, "Link all supported harnesses (do not require CLI on PATH)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print actions without writing files")
	cmd.Flags().BoolVar(&unlink, "unlink", false, "Remove managed muster snippets")
	return cmd
}
