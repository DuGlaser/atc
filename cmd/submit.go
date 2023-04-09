package cmd

import (
	"fmt"
	"os"

	"github.com/DuGlaser/atc/internal/handler"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var skipTests = false

var submitCmd = &cobra.Command{
	Use:     "submit [problem id[",
	Short:   "Submit answer",
	Aliases: []string{"s"},
	Args:    cobra.MinimumNArgs(1),
	PreRun:  func(cmd *cobra.Command, args []string) { handler.CheckLogin() },
	Run: func(cmd *cobra.Command, args []string) {
		displayID := args[0]
		if !skipTests {
			to := handler.TestOption{}
			to.DisplayID = displayID

			trs := handler.TestCode(to, verbose)
			if !trs.Pass {
				os.Exit(1)
			}

			if trs.MaxTimeMs >= 1000 {
				confirm()
			}
		}
		handler.SubmitCode(displayID, verbose)
	},
}

func confirm() {
	fmt.Println()
	prompt := promptui.Prompt{
		Label: "Execution time exceeds 1000ms. Would you like to submit? (yes / no)",
	}
	result, err := prompt.Run()
	cobra.CheckErr(err)

	if result != "yes" {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().BoolVar(&skipTests, "skip-tests", false, "Whether to skip the test and submit the answers")
}
