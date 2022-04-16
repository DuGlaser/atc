package handler

import (
	"fmt"
	"os"
	"path"

	"github.com/DuGlaser/atc/internal/core"
	"github.com/DuGlaser/atc/internal/repository/fetcher"
	"github.com/DuGlaser/atc/internal/repository/scraper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateProject(contestID string, verbose bool) {
	res, err := fetcher.FetchContestPage(contestID)
	cobra.CheckErr(err)

	defer res.Body.Close()

	cp, err := scraper.NewContestPage(res.Body)
	cobra.CheckErr(err)

	ids := cp.GetProblemIds()
	if len(ids) == 0 {
		res, err = fetcher.FetchProblems(contestID)
		cobra.CheckErr(err)
		defer res.Body.Close()

		tp, err := scraper.NewTasksPage(res.Body)
		cobra.CheckErr(err)

		ids = tp.GetProblemIds()

		if len(ids) == 0 {
			cobra.CheckErr(fmt.Sprintf("Problems doesn't exists in %s.\n", contestID))
		}
	}

	wd, err := os.Getwd()
	cobra.CheckErr(err)

	contestPath := path.Join(wd, contestID)

	if _, err := os.Stat(contestPath); os.IsExist(err) {
		cobra.CheckErr(fmt.Sprintf("%s directory is already exists.\n", contestID))
	}

	err = os.Mkdir(path.Join(wd, contestID), 0754)
	cobra.CheckErr(err)

	v := viper.New()

	var c core.Config
	cobra.CheckErr(viper.UnmarshalKey("config", &c))
	cobra.CheckErr(c.Validate())

	v.Set("config.lang", c.Lang)
	v.Set("config.runcmd", c.RunCmd)
	v.Set("config.buildcmd", c.BuildCmd)
	v.Set("config.filename", c.FileName)

	v.Set("contest.name", contestID)
	v.Set("contest.url", fetcher.GetAtcoderUrl("contests", contestID))

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
