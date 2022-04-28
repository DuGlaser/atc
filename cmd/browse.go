package cmd

import (
	"github.com/DuGlaser/atc/internal/handler"
	"github.com/spf13/cobra"
)

var browseCmd = &cobra.Command{
	Use:     "browse [problem id]",
	Short:   "open problem in browse",
	Aliases: []string{"b"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		displayID := args[0]
		handler.OpenTask(displayID)
	},
}

func init() {
	rootCmd.AddCommand(browseCmd)
}
