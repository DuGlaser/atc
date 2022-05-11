package handler

import (
	"fmt"

	"github.com/DuGlaser/atc/internal/repository/config"
	"github.com/DuGlaser/atc/internal/repository/fetcher"
	"github.com/DuGlaser/atc/internal/util"
	"github.com/spf13/cobra"
)

func OpenTask(displayID string, verbose bool) {
	cc, err := config.NewContestConfig()
	cobra.CheckErr(err)

	task, err := cc.ReadTaskSetting(displayID)
	cobra.CheckErr(err)

	contest, err := cc.ReadContestSetting()

	if task.ID == "" {
		err = cc.SetTaskID(displayID)
		cobra.CheckErr(err)
	}

	url := fetcher.GetAtcoderUrl("contests", contest.Name, "tasks", task.ID)

	if verbose {
		fmt.Printf("Open %s", url)
	}
	util.Openbrowser(url)
}
