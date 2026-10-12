package cli

import (
	"context"
	"path/filepath"

	"github.com/mlaccetti/muster/internal/pass"
	"github.com/mlaccetti/muster/internal/paths"
	"github.com/spf13/cobra"
)

func newPassCmd() *cobra.Command {
	passDir := ""
	file := ""
	poolsPath := ""
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run a declarative pass file",
		RunE: func(cmd *cobra.Command, args []string) error {
			spec, err := pass.LoadSpec(file)
			if err != nil {
				return err
			}
			if passDir == "" {
				passDir, err = paths.PassDir(spec.ID)
				if err != nil {
					return err
				}
			}
			if poolsPath == "" {
				poolsPath, _ = paths.PoolsFile()
			}
			code, err := pass.Run(context.Background(), spec, passDir, poolsPath, true)
			if err != nil {
				return err
			}
			if code != 0 {
				return exitErr(code)
			}
			return nil
		},
	}
	runCmd.Flags().StringVar(&file, "file", "", "Pass spec YAML")
	runCmd.Flags().StringVar(&passDir, "pass-dir", "", "Pass state directory")
	runCmd.Flags().StringVar(&poolsPath, "pools", "", "Pools JSON path")
	_ = runCmd.MarkFlagRequired("file")

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Scaffold a pass directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := passDir
			if dir == "" {
				dir = filepath.Join(".", "pass-example")
			}
			if err := pass.Init(dir); err != nil {
				return err
			}
			writef(cmd.OutOrStdout(), "Created %s/pass.yaml\n", dir)
			return nil
		},
	}
	initCmd.Flags().StringVar(&passDir, "pass-dir", "", "Directory to create")

	root := &cobra.Command{Use: "pass", Short: "Deterministic pass driver"}
	root.AddCommand(runCmd, initCmd)
	return root
}
