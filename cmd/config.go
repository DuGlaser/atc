package cmd

import (
	"github.com/DuGlaser/atc/internal/repository/fetcher"
	"github.com/DuGlaser/atc/internal/repository/scraper"
	"github.com/DuGlaser/atc/internal/ui"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config [subcommand]",
	Short: "Create and edit atc config",
	Run: func(_ *cobra.Command, args []string) {
		res, err := fetcher.FetchSubmitPage("abc001")
		cobra.CheckErr(err)
		defer res.Body.Close()

		sp, err := scraper.NewSubmitPage(res.Body)
		cobra.CheckErr(err)

		ls := sp.GetLanguageIds()
		lang, err := ui.SelectLanguage(ls)
		cobra.CheckErr(err)

		cmdPrompt := promptui.Prompt{
			Label:   "Input run command",
			Default: "g++ -o {{ .dir }}/main {{ .file }} && {{ .dir }}/main",
		}

		runCmd, err := cmdPrompt.Run()
		cobra.CheckErr(err)

		fileNamePrompt := promptui.Prompt{
			Label:   "Input file name",
			Default: "main.cpp",
		}

		fileName, err := fileNamePrompt.Run()
		cobra.CheckErr(err)

		viper.Set("config.lang", lang.ID)
		viper.Set("config.runcmd", runCmd)
		viper.Set("config.buildcmd", "")
		viper.Set("config.fileName", fileName)
		viper.Set("config.template", "")

		cobra.CheckErr(viper.SafeWriteConfig())
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
