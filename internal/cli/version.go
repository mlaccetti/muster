package cli

import (
	"context"
	"os"

	"github.com/mlaccetti/muster/internal/buildversion"
	"github.com/mlaccetti/muster/internal/release"
	"github.com/spf13/cobra"
)

func attachReleaseCheck(root *cobra.Command) {
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "completion" {
			return nil
		}
		release.MaybeNotify(context.Background(), buildversion.String(), os.Stderr)
		return nil
	}
}
