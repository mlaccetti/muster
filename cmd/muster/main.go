package main

import (
	"context"
	"os"

	"charm.land/fang/v2"
	"github.com/spf13/cobra"

	"github.com/mlaccetti/muster/internal/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.SetBuildVersion(version)
	root := cli.NewRoot()
	root.Version = version
	if commit != "none" {
		root.SetVersionTemplate(version + " (" + commit + ", " + date + ")\n")
	}

	if err := fang.Execute(context.Background(), root, fang.WithVersion(version)); err != nil {
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
