package handler

import (
	"io/ioutil"

	"github.com/DuGlaser/atc/internal/repository/config"
	"github.com/DuGlaser/atc/internal/repository/fetcher"
	"github.com/DuGlaser/atc/internal/util"
	"github.com/spf13/cobra"
)

func SubmitCode(displayID string, verbose bool) {
	cc, err := config.NewContestConfig()
	cobra.CheckErr(err)

	task, err := cc.ReadTaskSetting(displayID)
	cobra.CheckErr(err)

	if task.ID == "" {
		err := cc.SetTaskID(displayID)
		cobra.CheckErr(err)
	}

	bytes, err := ioutil.ReadFile(task.Path)
	cobra.CheckErr(err)

	contest, err := cc.ReadContestSetting()
	cobra.CheckErr(err)

	config, err := cc.ReadConfig()
	cobra.CheckErr(err)

	_, err = fetcher.PostProblemAnswer(contest.Name, task.ID, config.Lang, string(bytes))
	cobra.CheckErr(err)

	err = util.Openbrowser(fetcher.GetAtcoderUrl("contests", contest.Name, "submissions", "me"))
	cobra.CheckErr(err)
}
