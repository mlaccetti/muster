package cli

import (
	"os"
	"path/filepath"

	"github.com/mlaccetti/muster/internal/kit"
	"github.com/mlaccetti/muster/internal/render"
	"github.com/spf13/cobra"
)

func newRenderCmd() *cobra.Command {
	briefCmd := &cobra.Command{
		Use:   "brief",
		Short: "Render a brief from JSON input",
		RunE: func(cmd *cobra.Command, args []string) error {
			tmplName, _ := cmd.Flags().GetString("template")
			input, _ := cmd.Flags().GetString("input")
			out, _ := cmd.Flags().GetString("out")
			kitDir, err := kit.ActiveDir()
			if err != nil {
				kitDir, err = kit.Install()
				if err != nil {
					return err
				}
			}
			tmplPath := filepath.Join(kitDir, "templates", tmplName+".tmpl")
			text, err := render.BriefFromFile(tmplPath, input)
			if err != nil {
				return err
			}
			if out == "" || out == "-" {
				_, err = cmd.OutOrStdout().Write([]byte(text))
				return err
			}
			return os.WriteFile(out, []byte(text), 0o644)
		},
	}
	briefCmd.Flags().String("template", "implement", "Template name")
	briefCmd.Flags().String("input", "", "Brief JSON input path")
	briefCmd.Flags().String("out", "-", "Output path")
	_ = briefCmd.MarkFlagRequired("input")

	renderCmd := &cobra.Command{Use: "render", Short: "Render kit templates"}
	renderCmd.AddCommand(briefCmd)
	return renderCmd
}
