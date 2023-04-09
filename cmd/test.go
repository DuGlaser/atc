package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/DuGlaser/atc/internal/handler"
	"github.com/DuGlaser/atc/internal/util"
)

var testCmd = &cobra.Command{
	Use:     "test [problem id] [optional: test case numbers]",
	Short:   "Test answer",
	Aliases: []string{"t"},
	Args:    cobra.MinimumNArgs(1),
	PreRun:  func(cmd *cobra.Command, args []string) { handler.CheckLogin() },
	Run: func(_ *cobra.Command, args []string) {
		to := handler.TestOption{}
		to.DisplayID = args[0]

		if len(args) > 1 {
			is, err := util.ParseIntArguments(args[1:])
			cobra.CheckErr(err)

			to.EnableCaseIDs = is
		}

		trs := handler.TestCode(to, verbose)
		if !trs.Pass {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
