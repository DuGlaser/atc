package handler

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/DuGlaser/atc/internal/repository/config"
	"github.com/DuGlaser/atc/internal/repository/fetcher"
	"github.com/DuGlaser/atc/internal/repository/scraper"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
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

	fmt.Println()
	sm := watchSubmissionStatus(contest.Name)

	res, err := fetcher.FetchSubmissionDetail(contest.Name, sm.ID)
	cobra.CheckErr(err)

	defer res.Body.Close()

	sd, err := scraper.NewSubmissionDetailPage(res.Body)
	cobra.CheckErr(err)

	m, err := sd.GetSubmissionStatusMap()
	cobra.CheckErr(err)

	printSubmissionDetail(sm, m)
}

func printSubmissionDetail(sm *scraper.Submission, statuses map[scraper.StatusCode]int) {
	data := [][]string{
		{"Problem", sm.Task},
		{"Date", sm.Date},
		{"Memory", sm.ResultMetaData.Memory},
		{"Exec time", sm.ResultMetaData.ExecTime},
	}

	s := []string{}
	for sc, count := range statuses {
		c := color.BgYellow
		if sc == scraper.AC {
			c = color.BgGreen
		}

		status := color.New(c).Sprintf(" %s ", sc)

		s = append(s, fmt.Sprintf("%s x %d", status, count))
	}

	data = append(data, []string{"Status", strings.Join(s, " ")})

	fmt.Println()
	for _, v := range data {
		fmt.Printf("%-12v%s\n", v[0], v[1])
	}
}

func fetchSubmissionStatus(contest string) (*scraper.Submission, error) {
	res, err := fetcher.FetchSubmissionsMe(contest)
	cobra.CheckErr(err)

	defer res.Body.Close()

	sp, err := scraper.NewSubmissionsPage(res.Body)
	cobra.CheckErr(err)

	return sp.GetLatestSubmission()
}

func watchSubmissionStatus(contest string) *scraper.Submission {
	status := color.New(color.BgCyan).Sprintf(" %s ", scraper.WJ)
	total := 0

	getDescription := func(task string) string {
		return fmt.Sprintf(" [ %s ]  %s ", color.New(color.Bold).Sprint(task), status)
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.HideCursor = false
	s.Suffix = " Waiting for judge..."
	s.Start()

	// NOTE: 最新の提出が今提出したやつになるように遅延を入れる
	time.Sleep(2 * time.Second)

	var bar *progressbar.ProgressBar
	for {
		sm, err := fetchSubmissionStatus(contest)
		cobra.CheckErr(err)

		if sm.Status != scraper.WJ {
			c := color.BgYellow
			if sm.Status == scraper.AC {
				c = color.BgGreen
			}

			status = color.New(c).Sprintf(" %s ", sm.Status)
		}

		if sm.ResultMetaData != nil {
			if bar != nil {
				bar.Describe(getDescription(sm.Task))
				bar.Set(total)
			}

			fmt.Println()
			return sm
		}

		if sm.Counter != nil {
			if bar == nil {
				s.Stop()
				total = sm.Counter.Total
				bar = progressbar.NewOptions(total,
					progressbar.OptionEnableColorCodes(true),
					progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
					progressbar.OptionShowCount(),
					progressbar.OptionSetWidth(30),
					progressbar.OptionSetDescription(getDescription(sm.Task)),
					progressbar.OptionSetTheme(progressbar.Theme{
						Saucer:        "[green]=[reset]",
						SaucerHead:    "[green]>[reset]",
						SaucerPadding: " ",
						BarStart:      "[",
						BarEnd:        "]",
					}))
			}

			bar.Set(sm.Counter.Current)
			bar.Describe(getDescription(sm.Task))
		}

		time.Sleep(time.Second)
	}
}
