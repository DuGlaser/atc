package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"
	"time"

	"github.com/DuGlaser/atc/assets"
	"github.com/DuGlaser/atc/internal/core"
	"github.com/DuGlaser/atc/internal/repository/fetcher"
	"github.com/DuGlaser/atc/internal/repository/scraper"
	"github.com/spf13/viper"
)

type Contest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Task struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

type ContestConfig struct {
	v       *viper.Viper
	config  *core.Config
	contest *Contest
	task    map[string]*Task
}

func NewContestConfig() (*ContestConfig, error) {
	cc := &ContestConfig{
		v:    viper.New(),
		task: map[string]*Task{},
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

func (cc *ContestConfig) ReadContestSetting() (*Contest, error) {
	if cc.contest != nil {
		return cc.contest, nil
	}

	var contest Contest
	err := cc.v.UnmarshalKey("contest", &contest)
	if err != nil {
		return nil, err
	}

	cc.contest = &contest

	return &contest, nil
}

func (cc *ContestConfig) ReadTaskSetting(displayID string) (*Task, error) {
	if task, ok := cc.task[displayID]; ok {
		return task, nil
	}

	var task Task
	err := cc.v.UnmarshalKey(fmt.Sprintf("tasks.%s", displayID), &task)
	if task.Path == "" && task.ID == "" {
		return nil, fmt.Errorf("Not found task %s.", displayID)
	}

	if err != nil {
		return nil, err
	}

	cc.task[displayID] = &task
	return &task, nil
}

// FIXME: configの中にfetchがあるのは違和感がある
func (cc *ContestConfig) SetTaskID(displayID string) error {
	task, err := cc.ReadTaskSetting(displayID)
	if err != nil {
		return err
	}

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

	task.ID = p.ID

	return nil
}

func GetCurrentVersionLangList(targetTime time.Time) ([]core.Language, error) {
	validPaths := make([]string, 0)

	err := fs.WalkDir(assets.AtCoderAssets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		timestamp, err := time.Parse("200601021504", d.Name()[:12])
		if err != nil {
			return err
		}

		if timestamp.Before(targetTime) || timestamp.Equal(targetTime) {
			validPaths = append(validPaths, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(validPaths) == 0 {
		return nil, err
	}

	sort.Strings(validPaths)

	last := validPaths[len(validPaths)-1]

	data, err := fs.ReadFile(assets.AtCoderAssets, last)
	if err != nil {
		return nil, err
	}

	var l []core.Language
	err = json.Unmarshal(data, &l)
	if err != nil {
		return nil, err
	}

	return l, nil
}
