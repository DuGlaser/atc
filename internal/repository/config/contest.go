package config

import (
	"fmt"

	"github.com/DuGlaser/atc/internal"
	"github.com/DuGlaser/atc/internal/core"
	"github.com/DuGlaser/atc/internal/repository/fetcher"
	"github.com/DuGlaser/atc/internal/repository/scraper"
	"github.com/spf13/viper"
)

type ContestConfig struct {
	v       *viper.Viper
	config  *core.Config
	contest *internal.Contest
	task    *internal.Task
}

func NewContestConfig() (*ContestConfig, error) {
	cc := &ContestConfig{
		v: viper.New(),
	}
	err := cc.load()
	return cc, err
}

func (cc *ContestConfig) load() error {
	cc.v.AddConfigPath(".")
	cc.v.AddConfigPath("../")
	cc.v.SetConfigType("toml")
	cc.v.SetConfigName("contest")

	return cc.v.ReadInConfig()
}

func (cc *ContestConfig) ReadConfig() (*core.Config, error) {
	if cc.config != nil {
		return cc.config, nil
	}

	var config core.Config
	err := cc.v.UnmarshalKey("config", &config)
	if err != nil {
		return nil, err
	}

	cc.config = &config
	return &config, nil
}

func (cc *ContestConfig) ReadContestSetting() (*internal.Contest, error) {
	if cc.contest != nil {
		return cc.contest, nil
	}

	var contest internal.Contest
	err := cc.v.UnmarshalKey("contest", &contest)
	if err != nil {
		return nil, err
	}

	cc.contest = &contest

	return &contest, nil
}

func (cc *ContestConfig) ReadTaskSetting(displayID string) (*internal.Task, error) {
	if cc.task != nil {
		return cc.task, nil
	}

	var task internal.Task
	err := cc.v.UnmarshalKey(fmt.Sprintf("tasks.%s", displayID), &task)
	if task.Path == "" && task.ID == "" {
		return nil, fmt.Errorf("Not found task %s.", displayID)
	}

	if err != nil {
		return nil, err
	}

	cc.task = &task

	return &task, nil
}

// FIXME: configの中にfetchがあるのは違和感がある
func (cc *ContestConfig) SetTaskID(displayID string) error {
	contest, err := cc.ReadContestSetting()
	if err != nil {
		return err
	}

	res, err := fetcher.FetchProblems(contest.Name)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	tp, err := scraper.NewTasksPage(res.Body)
	if err != nil {
		return err
	}

	p := tp.GetProblemId(displayID)

	key := fmt.Sprintf("tasks.%s", p.DisplayID)
	cc.v.Set(fmt.Sprintf("%s.id", key), p.ID)
	err = cc.v.WriteConfig()
	if err != nil {
		return err
	}

	cc.task.ID = p.ID

	return nil
}
