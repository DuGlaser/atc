package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/DuGlaser/atc/internal"
	"github.com/DuGlaser/atc/internal/fetcher"
	"github.com/DuGlaser/atc/internal/scraper"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var testCmd = &cobra.Command{
	Use:     "test [problem id]",
	Short:   "Test answer",
	Aliases: []string{"t"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		problemId := args[0]
		testAnswer(problemId)
	},
}

type result struct {
	id       int
	in       string
	expected string
	got      string
}

func testAnswer(problemId string) {
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

	res, err := fetcher.FetchProblemPage(contest.Name, problemId)
	cobra.CheckErr(err)
	defer res.Body.Close()

	pp, err := scraper.NewProblemPage(res.Body)
	cobra.CheckErr(err)

	tests, err := pp.GetProblemSamples()
	cobra.CheckErr(err)

	tmpl, err := template.New("test").Parse(config.Cmd)
	cobra.CheckErr(err)

	var c bytes.Buffer
	cobra.CheckErr(tmpl.Execute(&c, map[string]interface{}{
		"file": task.Path,
		"dir":  task.Path[0 : len(task.Path)-len(config.FileName)-1],
	}))

	failures := []result{}

	for i, test := range tests {
		cmd := exec.Command("sh", "-c", c.String())

		pipe, err := cmd.StdinPipe()
		cobra.CheckErr(err)

		io.WriteString(pipe, test.In)
		pipe.Close()

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		out, err := cmd.Output()
		if err != nil {
			cobra.CheckErr(stderr.String())
		}

		got := strings.TrimRight(string(out), "\n")
		got = strings.TrimSpace(got)

		if got == test.Out {
			fmt.Printf("sample test case 1 ... %s\n", color.GreenString("success"))
		} else {
			fmt.Printf("sample test case 1 ... %s\n", color.RedString("failed"))
			failures = append(failures, result{id: i, in: test.In, expected: test.Out, got: got})
		}
	}

	for _, result := range failures {
		fmt.Println()
		color.New(color.Bold).Printf("=== sample test case %d ===\n", result.id)
		fmt.Println()
		color.Cyan("input:")
		printResultValue(result.in)
		fmt.Println()
		color.Green("expected:")
		printResultValue(result.expected)
		fmt.Println()
		color.Red("your output:")
		printResultValue(result.got)
	}

	if len(failures) > 0 {
		os.Exit(1)
	}
}

func printResultValue(v string) {
	ss := strings.Split(v, "\n")

	maxL := len(ss)/10 + 1

	for i, s := range ss {
		fmt.Printf("%*d | %s\n", maxL, i+1, s)
	}
}

func init() {
	rootCmd.AddCommand(testCmd)
}
