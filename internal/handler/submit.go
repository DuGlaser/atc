package handler

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"

	"github.com/DuGlaser/atc/internal/repository/config"
	"github.com/DuGlaser/atc/internal/repository/fetcher"
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

	err = openbrowser(fetcher.GetAtcoderUrl("contests", contest.Name, "submissions", "me"))
	cobra.CheckErr(err)
}

// FYI: https://gist.github.com/hyg/9c4afcd91fe24316cbf0
// Thank you!
func openbrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
