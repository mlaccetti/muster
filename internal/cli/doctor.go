package cli

import (
	"os"
	"path/filepath"
	"time"

	"github.com/mlaccetti/muster/internal/harnesscli"
	"github.com/mlaccetti/muster/internal/kit"
	"github.com/mlaccetti/muster/internal/link"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	probe := false
	live := false
	liveTimeout := 3 * time.Minute
	codexBin := "codex"
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check local muster installation",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := kit.ActiveDir()
			if err != nil {
				writef(cmd.ErrOrStderr(), "kit: %v\n", err)
			} else {
				writef(cmd.OutOrStdout(), "kit: ok (%s)\n", dir)
				for _, name := range []string{"orchestrator.md", "worker.md"} {
					p := filepath.Join(dir, "policy", name)
					if _, err := os.Stat(p); err != nil {
						writef(cmd.ErrOrStderr(), "policy %s: missing\n", name)
					}
				}
			}
			if m, err := link.LoadManifest(); err == nil {
				writef(cmd.OutOrStdout(), "link: ok (%d harnesses, updated %s)\n", len(m.Harnesses), m.UpdatedAt)
				if staleKitLink(m.KitDir, dir) {
					writef(cmd.ErrOrStderr(), "link: stale kit paths in harness snippets (manifest %s, active %s); run muster link\n", m.KitDir, dir)
				}
				for _, issue := range link.VerifyReports(m) {
					writef(cmd.ErrOrStderr(), "link: %s\n", issue)
				}
			} else {
				writef(cmd.ErrOrStderr(), "link: not configured (run muster link)\n")
			}
			if os.Getenv("MUSTER_SKIP_PROBE") == "1" || !probe {
				writelnf(cmd.OutOrStdout(), "probe: skipped")
			} else {
				for _, id := range harnesscli.AllIDs() {
					path, ver, err := harnesscli.ResolveCLI(id)
					if err != nil {
						writef(cmd.ErrOrStderr(), "%s: %v\n", id, err)
						continue
					}
					writef(cmd.OutOrStdout(), "%s: %s (%s)\n", id, path, ver)
				}
			}
			if live {
				result, err := runLiveDoctor(cmd, codexBin, liveTimeout)
				if err != nil {
					writef(cmd.ErrOrStderr(), "live: failed (%v)\n", err)
					return err
				}
				if staticErr := validateCodexLiveRouting(); staticErr != nil {
					if result.Agent != "muster-worker" {
						writef(cmd.ErrOrStderr(), "live: fallback verified (root %s/%s -> child role %q %s/%s; root %s, child %s)\n", result.RootModel, result.RootEffort, result.Agent, result.Model, result.Effort, result.RootID, result.ChildID)
					}
					writef(cmd.ErrOrStderr(), "live: failed (Codex muster-worker routing is missing or invalid)\n")
					return staticErr
				}
				if result.Agent != "muster-worker" {
					writef(cmd.OutOrStdout(), "live: fallback verified (root %s/%s -> child role %q %s/%s; root %s, child %s)\n", result.RootModel, result.RootEffort, result.Agent, result.Model, result.Effort, result.RootID, result.ChildID)
				} else {
					writef(cmd.OutOrStdout(), "live: ok (root %s/%s -> %s %s/%s; root %s, child %s; turns %s/%s)\n", result.RootModel, result.RootEffort, result.Agent, result.Model, result.Effort, result.RootID, result.ChildID, result.RootTurnID, result.ChildTurnID)
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&probe, "probe", false, "Probe harness binaries with --version")
	cmd.Flags().BoolVar(&live, "live", false, "Run an opt-in Codex delegation check")
	cmd.Flags().DurationVar(&liveTimeout, "live-timeout", liveTimeout, "Timeout for the live Codex delegation check")
	cmd.Flags().StringVar(&codexBin, "codex-bin", codexBin, "Codex executable to use for --live")
	return cmd
}
