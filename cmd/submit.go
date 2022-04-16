package cmd

import (
	"github.com/DuGlaser/atc/internal/handler"
	"github.com/spf13/cobra"
)

var skipTests = false

var submitCmd = &cobra.Command{
	Use:     "submit [problem id[",
	Short:   "Submit answer",
	Aliases: []string{"s"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		displayID := args[0]
		if !skipTests {
			handler.TestCode(displayID, verbose)
		}
		handler.SubmitCode(displayID, verbose)
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().BoolVar(&skipTests, "skip-tests", false, "Whether to skip the test and submit the answers")
}
