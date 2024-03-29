package cmd

import (
	"github.com/spf13/cobra"

	"github.com/DuGlaser/atc/internal/handler"
)

var handleCmd = &cobra.Command{
	Use:     "handle [problem id]",
	Short:   "Manual execution code",
	Aliases: []string{"h"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		displayID := args[0]

		handler.HaldleExec(displayID)
	},
}

func init() {
	rootCmd.AddCommand(handleCmd)
}
