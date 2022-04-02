package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/DuGlaser/atc/internal"
	"github.com/DuGlaser/atc/internal/fetcher"
	"github.com/DuGlaser/atc/internal/scraper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newCmd = &cobra.Command{
	Use:     "new [CONTEST]",
	Short:   "Create contest project",
	Aliases: []string{"n"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(viper.ReadInConfig())
		contest := args[0]
		createProject(contest)
	},
}

func createProject(contest string) {
	res, err := fetcher.FetchContestPage(contest)
	cobra.CheckErr(err)

	defer res.Body.Close()

	cp, err := scraper.NewContestPage(res.Body)
	cobra.CheckErr(err)

	ids := cp.GetProblemIds()
	if len(ids) == 0 {
		res, err = fetcher.FetchProblems(contest)
		cobra.CheckErr(err)
		defer res.Body.Close()

		tp, err := scraper.NewTasksPage(res.Body)
		cobra.CheckErr(err)

		ids = tp.GetProblemIds()

		if len(ids) == 0 {
			cobra.CheckErr(fmt.Sprintf("Problems doesn't exists in %s.\n", contest))
		}
	}

	wd, err := os.Getwd()
	cobra.CheckErr(err)

	contestPath := path.Join(wd, contest)

	if _, err := os.Stat(contestPath); os.IsExist(err) {
		cobra.CheckErr(fmt.Sprintf("%s directory is already exists.\n", contest))
	}

	err = os.Mkdir(path.Join(wd, contest), 0754)
	cobra.CheckErr(err)

	v := viper.New()

	var c internal.Config
	cobra.CheckErr(viper.UnmarshalKey("config", &c))
	cobra.CheckErr(c.Validate())

	v.Set("config.lang", c.Lang)
	v.Set("config.cmd", c.Cmd)
	v.Set("config.filename", c.FileName)

	v.Set("contest.name", contest)
	v.Set("contest.url", fetcher.GetAtcoderUrl("contests", contest))

	for _, id := range ids {
		key := fmt.Sprintf("tasks.%s", id.DisplayID)

		dir := path.Join(contestPath, id.DisplayID)
		err = os.Mkdir(dir, 0754)
		cobra.CheckErr(err)

		fp := path.Join(dir, c.FileName)
		f, err := os.Create(fp)
		cobra.CheckErr(err)

		_, err = f.WriteString(c.Template)
		cobra.CheckErr(err)

		v.Set(fmt.Sprintf("%s.id", key), id.ID)
		v.Set(fmt.Sprintf("%s.path", key), fp)
	}

	v.SafeWriteConfigAs(path.Join(contestPath, "contest.toml"))
}

func init() {
	rootCmd.AddCommand(newCmd)
}
