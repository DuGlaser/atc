package cmd

import (
	"fmt"

	"github.com/DuGlaser/atc/internal/auth"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout to atcoder",
	Long:  "Logout to atcoder and save the session cookie locally.",
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(auth.ClearSession())
		fmt.Println("Success!")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
