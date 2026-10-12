package cli

import (
	"os"

	"github.com/mlaccetti/muster/internal/adapter"
	"github.com/mlaccetti/muster/internal/kit"
	"github.com/mlaccetti/muster/internal/link"
	"github.com/mlaccetti/muster/internal/schema"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	probe := false
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check local muster installation",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := kit.ActiveDir()
			if err != nil {
				writef(cmd.ErrOrStderr(), "kit: %v\n", err)
			} else {
				writef(cmd.OutOrStdout(), "kit: ok (%s)\n", dir)
			}
			for _, kind := range []string{"brief", "report", "routing", "pools"} {
				if _, err := schema.SchemaOK(kind); err != nil {
					writef(cmd.ErrOrStderr(), "schema %s: %v\n", kind, err)
				}
			}
			if m, err := link.LoadManifest(); err == nil {
				writef(cmd.OutOrStdout(), "link: ok (%d harnesses, updated %s)\n", len(m.Harnesses), m.UpdatedAt)
				if staleKitLink(m.KitDir, dir) {
					writef(cmd.ErrOrStderr(), "link: stale kit paths in harness snippets (manifest %s, active %s); run muster link\n", m.KitDir, dir)
				}
			} else {
				writef(cmd.OutOrStderr(), "link: not configured (run muster link)\n")
			}
			if os.Getenv("MUSTER_SKIP_PROBE") == "1" || !probe {
				writelnf(cmd.OutOrStdout(), "probe: skipped")
				return nil
			}
			for _, id := range adapter.AllIDs() {
				a, err := adapter.Get(id)
				if err != nil {
					continue
				}
				path, ver, err := a.ResolveCLI()
				if err != nil {
					writef(cmd.ErrOrStderr(), "%s: %v\n", id, err)
					continue
				}
				writef(cmd.OutOrStdout(), "%s: %s (%s)\n", id, path, ver)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&probe, "probe", false, "Probe harness binaries with --version")
	return cmd
}
