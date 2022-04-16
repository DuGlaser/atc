package cmd

import (
	"fmt"
	"strings"

	"github.com/DuGlaser/atc/internal/repository/fetcher"
	"github.com/DuGlaser/atc/internal/repository/scraper"
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

		searcher := func(input string, index int) bool {
			l := ls[index]
			name := strings.Replace(strings.ToLower(l.Name), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)

			return strings.Contains(name, input)
		}

		templates := &promptui.SelectTemplates{
			Active:   fmt.Sprintf("%s {{ .Name | cyan | underline }}", promptui.IconSelect),
			Inactive: "  {{ .Name }}",
		}

		langPrompt := promptui.Select{
			Label:     "Select language",
			Items:     ls,
			Searcher:  searcher,
			Templates: templates,
		}

		i, _, err := langPrompt.Run()
		cobra.CheckErr(err)

		lang := ls[i]

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

		viper.Set("config.lang", lang.Value)
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
