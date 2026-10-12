package cli

import (
	"strings"

	"github.com/mlaccetti/muster/internal/link"
	"github.com/spf13/cobra"
)

func newLinkCmd() *cobra.Command {
	var (
		projectDir string
		harnesses  []string
		all        bool
		global     bool
		project    bool
		dryRun     bool
		unlink     bool
	)
	cmd := &cobra.Command{
		Use:   "link",
		Short: "Wire muster orchestrator and worker policies into harness configs",
		Long: `Link installs managed snippets that point every harness at the same kit policies.

By default, muster auto-detects harnesses whose CLIs are on PATH and links both global
(user) and project (current repo) entrypoints. Use --harness to pin specific IDs, or
--all to link every supported harness even when not installed.`,
		Example: `  muster install
  muster link
  muster link --harness claude --harness codex --project-dir .
  muster link --all --dry-run
  muster link --unlink`,
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := link.Run(link.Options{
				ProjectDir: projectDir,
				Harnesses:  harnesses,
				All:        all,
				Global:     global,
				Project:    project,
				DryRun:     dryRun,
				Unlink:     unlink,
			})
			if err != nil {
				return err
			}
			if unlink {
				writelnf(cmd.OutOrStdout(), "muster: unlinked")
				return nil
			}
			writef(cmd.OutOrStdout(), "kit: %s\n", res.KitDir)
			writef(cmd.OutOrStdout(), "project: %s\n", res.Project)
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
	cmd.Flags().StringVar(&projectDir, "project-dir", "", "Project directory (default: cwd)")
	cmd.Flags().StringArrayVar(&harnesses, "harness", nil, "Harness id (repeatable); default: auto-detect on PATH")
	cmd.Flags().BoolVar(&all, "all", false, "Link all supported harnesses (do not require CLI on PATH)")
	cmd.Flags().BoolVar(&global, "global", true, "Link user-level harness config")
	cmd.Flags().BoolVar(&project, "project", true, "Link project-level config (.muster, CLAUDE.md, AGENTS.md, .cursor/rules)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print actions without writing files")
	cmd.Flags().BoolVar(&unlink, "unlink", false, "Remove managed muster snippets")
	return cmd
}
