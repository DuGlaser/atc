package cmd

import (
	"fmt"

	"github.com/DuGlaser/atc/internal/repository/fetcher"
	"github.com/DuGlaser/atc/internal/repository/scraper"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Output the information of logged in users",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := fetcher.FetchHomePage()
		cobra.CheckErr(err)
		defer res.Body.Close()

		hp, err := scraper.NewHomePage(res.Body)
		cobra.CheckErr(err)

		name := hp.GetUserName()

		fmt.Printf("Hi, %s!\n", name)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
