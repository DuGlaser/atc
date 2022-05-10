package handler

import (
	"bytes"
	"fmt"
	"html/template"
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
	ID       int
	In       string
	Expected string
	Got      string
	Pass     bool
	TimeMs   int64
}

type TestResults struct {
	Results   []result
	MaxTimeMs int64
	Pass      bool
}

func TestCode(displayID string, verbose bool) TestResults {
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

	results := execTestCase(t, tests, verbose)
	failures := []result{}

	trs := TestResults{
		Results:   results,
		Pass:      true,
		MaxTimeMs: 0,
	}

	for _, result := range results {
		if !result.Pass {
			failures = append(failures, result)
			trs.Pass = false
		}

		if result.TimeMs > trs.MaxTimeMs {
			trs.MaxTimeMs = result.TimeMs
		}
	}

	printFailedCase(failures)

	return trs
}

func execTestCase(t *core.Task, tests []*core.TestCase, verbose bool) []result {
	results := []result{}

	for i, test := range tests {
		r, err := t.ExecCode(test.In, verbose)
		cobra.CheckErr(err)

		pass := test.Compare(r.Out)

		f := fmt.Sprintf("sample test case %d", i+1)

		time := fmt.Sprintf(" %dms ", r.TimeMs)

		if r.TimeMs >= 1000 {
			warn := color.New(color.BgRed).SprintFunc()
			time = warn(time)
		} else {
			info := color.New(color.BgBlue).SprintFunc()
			time = info(time)
		}

		if pass {
			fmt.Printf("%s ... %s  %s\n", f, color.GreenString("success"), time)
		} else {
			fmt.Printf("%s ... %s   %s\n", f, color.RedString("failed"), time)
		}

		results = append(results,
			result{
				ID:       i + 1,
				In:       test.In,
				Expected: test.Expected,
				Got:      r.Out,
				Pass:     pass,
				TimeMs:   r.TimeMs,
			})
	}

	return results
}

func printFailedCase(failures []result) {
	for _, result := range failures {
		fmt.Println()
		color.New(color.Bold).Printf("=== sample test case %d ===\n", result.ID)

		logs := []struct {
			label   string
			content string
		}{
			{label: color.CyanString("input: "), content: result.In},
			{label: color.GreenString("expected: "), content: result.Expected},
			{label: color.RedString("your output: "), content: result.Got},
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
