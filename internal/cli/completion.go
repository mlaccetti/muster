package cli

import (
	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
)

func setupCompletions(root *cobra.Command) {
	carapace.Gen(root).Standalone()
	harness := carapace.ActionValues("claude", "codex", "cursor", "antigravity", "agy").Tag("harness")
	tier := carapace.ActionValues("implement", "review", "plan", "cheap").Tag("tier")

	if w, _, _ := root.Find([]string{"worker"}); w != nil {
		carapace.Gen(w).FlagCompletion(carapace.ActionMap{
			"harness": harness,
			"tier":    tier,
		})
	}
	if d, _, _ := root.Find([]string{"dispatch"}); d != nil {
		carapace.Gen(d).FlagCompletion(carapace.ActionMap{
			"tier": tier,
		})
	}
}
