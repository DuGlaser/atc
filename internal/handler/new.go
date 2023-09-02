package handler

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/DuGlaser/atc/internal/core"
	"github.com/DuGlaser/atc/internal/repository/config"
	"github.com/DuGlaser/atc/internal/repository/fetcher"
	"github.com/DuGlaser/atc/internal/repository/scraper"
	"github.com/DuGlaser/atc/internal/ui"
	"github.com/DuGlaser/atc/internal/util"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateProject(contestID string) {
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
	util.InfoLog(fmt.Sprintf("problem ids %s", util.JsonLog(ids)))

	wd, err := os.Getwd()
	cobra.CheckErr(err)

	contestPath := path.Join(wd, contestID)

	if _, err := os.Stat(contestPath); os.IsExist(err) {
		cobra.CheckErr(fmt.Sprintf("%s directory is already exists.\n", contestID))
	}

	v := viper.New()

	var c core.Config
	cobra.CheckErr(viper.UnmarshalKey("config", &c))
	cobra.CheckErr(c.Validate())

	startAt, err := cp.GetStartAt()
	cobra.CheckErr(err)

	ls, err := config.GetCurrentVersionLangList(startAt)
	cobra.CheckErr(err)

	match := false
	for _, l := range ls {
		if l.ID == c.Lang {
			v.Set("config.lang", c.Lang)
			match = true
			break
		}
	}

	if !match {
		fmt.Println(color.RedString("The language described in `.atc.toml` is not supported in this contest. Please choose the language you would like to use for this contest."))

		lang, err := ui.SelectLanguage(ls)
		cobra.CheckErr(err)
		v.Set("config.lang", lang.ID)
	}

	v.Set("config.runcmd", c.RunCmd)
	v.Set("config.buildcmd", c.BuildCmd)
	v.Set("config.filename", c.FileName)

	v.Set("contest.name", strings.ToLower(contestID))
	v.Set("contest.url", fetcher.GetAtcoderUrl("contests", strings.ToLower(contestID)))

	err = os.Mkdir(path.Join(wd, contestID), 0754)
	cobra.CheckErr(err)

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

	util.InfoLog(fmt.Sprintf("contest config %s", util.JsonLog(c)))
	v.SafeWriteConfigAs(path.Join(contestPath, "contest.toml"))
}
