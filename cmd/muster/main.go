package main

import (
	"context"
	"os"

	"charm.land/fang/v2"
	"github.com/spf13/cobra"

	"github.com/mlaccetti/muster/internal/buildversion"
	"github.com/mlaccetti/muster/internal/cli"
)

func main() {
	root := cli.NewRoot()
	v := buildversion.String()
	root.Version = v
	root.SetVersionTemplate(buildversion.Line())

	if err := fang.Execute(context.Background(), root, fang.WithVersion(v)); err != nil {
		code := cli.ExitCode(err)
		if root.SilenceErrors {
			if code == 1 {
				cobra.CheckErr(err)
			} else {
				os.Exit(code)
			}
		}
		os.Exit(code)
	}
}
