package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"

	"github.com/DuGlaser/atc/internal/core"
	"github.com/DuGlaser/atc/internal/repository/config"
	"github.com/DuGlaser/atc/internal/repository/fetcher"
	"github.com/DuGlaser/atc/internal/repository/scraper"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type result struct {
	id       int
	in       string
	expected string
	got      string
}

func TestCode(displayID string, verbose bool) {
	cc, err := config.NewContestConfig()
	cobra.CheckErr(err)

	task, err := cc.ReadTaskSetting(displayID)
	cobra.CheckErr(err)

	if task.ID == "" {
		err := cc.SetTaskID(displayID)
		cobra.CheckErr(err)
	}

	contest, err := cc.ReadContestSetting()
	cobra.CheckErr(err)

	if verbose {
		fmt.Println("Fetch test cases...")
	}
	res, err := fetcher.FetchProblemPage(contest.Name, task.ID)
	cobra.CheckErr(err)
	defer res.Body.Close()

	pp, err := scraper.NewTaskPage(res.Body)
	cobra.CheckErr(err)

	tests, err := pp.GetTaskTestCases()
	cobra.CheckErr(err)

	if len(tests) == 0 {
		cobra.CheckErr(fmt.Errorf("%s_%s has no tests.",
			strings.ToLower(contest.Name), strings.ToLower(task.ID)))
	}

	config, err := cc.ReadConfig()
	cobra.CheckErr(err)

	rcTmpl, err := template.New("runCmd").Parse(config.RunCmd)
	cobra.CheckErr(err)

	var rc bytes.Buffer
	cobra.CheckErr(rcTmpl.Execute(&rc, map[string]interface{}{
		"file": task.Path,
		"dir":  task.Path[0 : len(task.Path)-len(config.FileName)-1],
	}))

	t := &core.Task{
		RunCmd: rc.String(),
	}

	if config.BuildCmd != "" {
		bcTmpl, err := template.New("buildCmd").Parse(config.BuildCmd)
		cobra.CheckErr(err)

		var bc bytes.Buffer
		cobra.CheckErr(bcTmpl.Execute(&bc, map[string]interface{}{
			"file": task.Path,
			"dir":  task.Path[0 : len(task.Path)-len(config.FileName)-1],
		}))

		t.BuildCmd = bc.String()
	}

	failures := execTestCase(t, tests, verbose)
	printFailedCase(failures)
	if len(failures) > 0 {
		fmt.Println("")
		os.Exit(1)
	}
}

func execTestCase(t *core.Task, tests []*core.TestCase, verbose bool) []result {
	failures := []result{}

	for i, test := range tests {
		got, err := t.ExecCode(test.In, verbose)
		cobra.CheckErr(err)

		pass := test.Compare(got)

		f := fmt.Sprintf("sample test case %d", i+1)
		if pass {
			fmt.Printf("%s ... %s\n", f, color.GreenString("success"))
		} else {
			fmt.Printf("%s ... %s\n", f, color.RedString("failed"))
			failures = append(failures, result{id: i + 1, in: test.In, expected: test.Expected, got: got})
		}
	}

	return failures
}

func printFailedCase(failures []result) {
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
}

func printResultValue(v string) {
	ss := strings.Split(v, "\n")

	maxL := len(strconv.Itoa(len(ss))) + 1

	for i, s := range ss {
		fmt.Printf("%*d | %s\n", maxL, i+1, s)
	}
}
