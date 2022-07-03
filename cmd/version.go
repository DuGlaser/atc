package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of atc.",
	Run: func(_ *cobra.Command, args []string) {
		fmt.Print(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
