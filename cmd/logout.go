package cmd

import (
	"fmt"

	"github.com/DuGlaser/atc/internal/auth"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout to AtCoder",
	Long:  "Delete local session cookie and logout of AtCoder.",
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(auth.ClearSession())
		fmt.Println("Success!")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
