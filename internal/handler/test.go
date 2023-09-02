package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DuGlaser/atc/internal"
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
	Err      error
}

type TestResults struct {
	Results   []result
	MaxTimeMs int64
	Pass      bool
}

type TestOption struct {
	DisplayID     string
	EnableCaseIDs []int
}

func TestCode(option TestOption) TestResults {
	cc, err := config.NewContestConfig()
	cobra.CheckErr(err)

	task, err := cc.ReadTaskSetting(option.DisplayID)
	cobra.CheckErr(err)

	if task.ID == "" {
		err := cc.SetTaskID(option.DisplayID)
		cobra.CheckErr(err)
	}

	contest, err := cc.ReadContestSetting()
	cobra.CheckErr(err)

	if internal.Verbose {
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

	cobra.CheckErr(config.GenerateCmd(task.Path, config.FileName))

	t := &core.Task{
		RunCmd:   config.RunCmd,
		BuildCmd: config.BuildCmd,
	}

	tests, err = filterTestCase(tests, option)
	cobra.CheckErr(err)

	results := execTestCase(t, tests)
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

func filterTestCase(tests []core.TestCase, option TestOption) ([]core.TestCase, error) {
	if len(option.EnableCaseIDs) == 0 {
		return tests, nil
	}

	filtered := []core.TestCase{}
	for _, i := range option.EnableCaseIDs {
		if i > len(tests) || 0 >= i {
			return nil, fmt.Errorf("%d-th test does not exist.", i)
		}

		filtered = append(filtered, tests[i-1])
	}

	return filtered, nil
}

func execTestCase(t *core.Task, tests []core.TestCase) []result {
	results := []result{}

	for _, test := range tests {
		r, err := t.ExecCode(test.In)

		pass := test.Compare(r.Out)

		f := fmt.Sprintf("sample test case %d", test.ID)

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
				ID:       test.ID,
				In:       test.In,
				Expected: test.Expected,
				Got:      r.Out,
				Pass:     pass,
				TimeMs:   r.TimeMs,
				Err:      err,
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

		if result.Err != nil {
			fmt.Println()
			fmt.Printf("%s %s\n", color.HiRedString("Error:"), result.Err)
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
