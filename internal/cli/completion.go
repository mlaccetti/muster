package cli

import (
	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
)

func setupCompletions(root *cobra.Command) {
	carapace.Gen(root).Standalone()
}
