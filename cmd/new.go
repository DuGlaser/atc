package cmd

import (
	"github.com/DuGlaser/atc/internal/handler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newCmd = &cobra.Command{
	Use:     "new [CONTEST]",
	Short:   "Create contest project",
	Aliases: []string{"n"},
	Args:    cobra.MinimumNArgs(1),
	PreRun:  func(cmd *cobra.Command, args []string) { handler.CheckLogin() },
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(viper.ReadInConfig())
		contest := args[0]
		handler.CreateProject(contest, true)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
