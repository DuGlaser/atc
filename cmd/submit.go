package cmd

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"

	"github.com/DuGlaser/atc/internal"
	"github.com/DuGlaser/atc/internal/fetcher"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var skipTests = false

var submitCmd = &cobra.Command{
	Use:     "submit [problem id[",
	Short:   "Submit answer",
	Aliases: []string{"s"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		problemId := args[0]
		if !skipTests {
			testAnswer(problemId)
		}
		submitAnswer(problemId)
	},
}

func submitAnswer(problemId string) {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("../")
	v.SetConfigType("toml")
	v.SetConfigName("contest")

	cobra.CheckErr(v.ReadInConfig())

	var config internal.Config
	cobra.CheckErr(v.UnmarshalKey("config", &config))
	cobra.CheckErr(config.Validate())

	var contest internal.Contest
	v.UnmarshalKey("contest", &contest)

	var task internal.Task
	v.UnmarshalKey(fmt.Sprintf("tasks.%s", problemId), &task)

	bytes, err := ioutil.ReadFile(task.Path)
	cobra.CheckErr(err)

	_, err = fetcher.PostProblemAnswer(contest.Name, problemId, config.Lang, string(bytes))
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

func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().BoolVar(&skipTests, "skip-tests", false, "Whether to skip the test and submit the answers")
}
