package cmd

import (
	"github.com/spf13/cobra"

	"github.com/DuGlaser/atc/internal/handler"
)

var testCmd = &cobra.Command{
	Use:     "test [problem id]",
	Short:   "Test answer",
	Aliases: []string{"t"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		displayID := args[0]
		handler.TestCode(displayID, true)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
