package cmd

import (
	"os"

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
		trs := handler.TestCode(displayID, verbose)
		if !trs.Pass {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
