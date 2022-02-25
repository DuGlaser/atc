package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"strconv"
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

	res, err := fetcher.FetchProblemPage(contest.Name, task.ID)
	cobra.CheckErr(err)
	defer res.Body.Close()

	pp, err := scraper.NewProblemPage(res.Body)
	cobra.CheckErr(err)

	tests, err := pp.GetProblemSamples()
	cobra.CheckErr(err)

	if len(tests) == 0 {
		cobra.CheckErr(fmt.Errorf("%s_%s has no tests.",
			strings.ToLower(contest.Name), strings.ToLower(task.ID)))
	}

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

		got := strings.TrimSpace(string(out))
		got = strings.TrimLeft(got, "\n")

		outLines := strings.Split(got, "\n")
		testLines := strings.Split(test.Out, "\n")

		pass := true
		if len(outLines) != len(testLines) {
			pass = false
		} else {
			for i := range outLines {
				if strings.TrimSpace(outLines[i]) != strings.TrimSpace(testLines[i]) {
					pass = false
					break
				}
			}
		}

		f := fmt.Sprintf("sample test case %d", i+1)

		if pass {
			fmt.Printf("%s ... %s\n", f, color.GreenString("success"))
		} else {
			fmt.Printf("%s ... %s\n", f, color.RedString("failed"))
			failures = append(failures, result{id: i + 1, in: test.In, expected: test.Out, got: got})
		}
	}

	for _, result := range failures {
		fmt.Println()
		color.New(color.Bold).Printf("=== sample test case %d ===\n", result.id)

		logs := []struct {
			label   string
			content string
		}{
			{label: color.CyanString("input: "), content: result.in},
			{label: color.GreenString("expected: "), content: result.expected},
			{label: color.RedString("your output: "), content: result.got},
		}

		for _, log := range logs {
			fmt.Println()
			fmt.Println(log.label)
			printResultValue(log.content)
		}
	}

	if len(failures) > 0 {
		os.Exit(1)
	}
}

func printResultValue(v string) {
	ss := strings.Split(v, "\n")

	maxL := len(strconv.Itoa(len(ss))) + 1

	for i, s := range ss {
		fmt.Printf("%*d | %s\n", maxL, i+1, s)
	}
}

func init() {
	rootCmd.AddCommand(testCmd)
}
