package cmd

import (
	"fmt"

	"github.com/DuGlaser/atc/internal"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of atc.",
	Run: func(_ *cobra.Command, args []string) {
		fmt.Print(internal.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
