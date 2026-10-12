package cli

import (
	"context"
	"os"

	"github.com/mlaccetti/muster/internal/release"
	"github.com/spf13/cobra"
)

var buildVersion = "dev"

// SetBuildVersion sets the CLI version from GoReleaser ldflags (call from main).
func SetBuildVersion(v string) {
	buildVersion = v
}

func attachReleaseCheck(root *cobra.Command) {
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "completion" {
			return nil
		}
		release.MaybeNotify(context.Background(), buildVersion, os.Stderr)
		return nil
	}
}
