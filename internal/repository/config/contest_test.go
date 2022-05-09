package config

import (
	"bytes"
	"testing"

	"github.com/DuGlaser/atc/internal/core"
	"github.com/spf13/viper"
)

func prepareContestConfig() (*ContestConfig, error) {
	cc := &ContestConfig{
		v:    viper.New(),
		task: map[string]Task{},
	}
	cc.v.SetConfigType("toml")

	testConfig := []byte(`
[config]
  buildcmd = "test buildcmd"
  filename = "test filename"
  lang = "test lang"
  runcmd = "test runcmd"

[contest]
  name = "test001"
  url = "https://atcoder.jp/contests/test001?lang=ja"

[tasks]

  [tasks.a]
    id = "a"
    path = "test_a_path"

  [tasks.b]
    id = "b"
    path = "test_b_path"

  [tasks.c]
    id = "c"
    path = "test_c_path"
  `)

	if err := cc.v.ReadConfig(bytes.NewBuffer(testConfig)); err != nil {
		return nil, err
	}

	return cc, nil
}

func TestReadContestSetting(t *testing.T) {
	cc, err := prepareContestConfig()
	if err != nil {
		t.Fatalf("prepareContestConfig() returned error. got=%s", err.Error())
	}

	contest, err := cc.ReadContestSetting()
	if err != nil {
		t.Fatalf("cc.ReadContestSetting()() returned error. got=%s", err.Error())
	}

	test := Contest{"test001", "https://atcoder.jp/contests/test001?lang=ja"}

	if contest.Name != test.Name {
		t.Errorf("contest.Name is not %s. got=%s", test.Name, contest.Name)
	}

	if contest.Url != test.Url {
		t.Errorf("contest.Url is not %s. got=%s", test.Url, contest.Url)
	}
}

func TestReadTasksSetting(t *testing.T) {
	cc, err := prepareContestConfig()
	if err != nil {
		t.Fatalf("prepareContestConfig() returned error. got=%s", err.Error())
	}

	tests := map[string]Task{
		"a": {"a", "test_a_path"},
		"b": {"b", "test_b_path"},
		"c": {"c", "test_c_path"},
	}

	for key, test := range tests {
		task, err := cc.ReadTaskSetting(key)
		if err != nil {
			t.Errorf("cc.ReadTaskSetting()() returned error. got=%s", err.Error())
			continue
		}

		if task.ID != test.ID {
			t.Errorf("task.ID is not %s. got=%s", test.ID, task.ID)
		}

		if task.Path != test.Path {
			t.Errorf("task.Url is not %s. got=%s", test.Path, task.Path)
		}
	}
}

func TestReadConfig(t *testing.T) {
	cc, err := prepareContestConfig()
	if err != nil {
		t.Fatalf("prepareContestConfig() returned error. got=%s", err.Error())
	}

	config, err := cc.ReadConfig()
	if err != nil {
		t.Fatalf("cc.ReadConfig() returned error. got=%s", err.Error())
	}

	test := core.Config{
		BuildCmd: "test buildcmd",
		RunCmd:   "test runcmd",
		Lang:     "test lang",
		FileName: "test filename",
	}

	if config.BuildCmd != test.BuildCmd {
		t.Errorf("config.BuildCmd is not %s. got=%s", test.BuildCmd, config.BuildCmd)
	}

	if config.RunCmd != test.RunCmd {
		t.Errorf("config.RunCmd is not %s. got=%s", test.RunCmd, config.RunCmd)
	}

	if config.Lang != test.Lang {
		t.Errorf("config.Lang is not %s. got=%s", test.Lang, config.Lang)
	}

	if config.FileName != test.FileName {
		t.Errorf("config.FileName is not %s. got=%s", test.FileName, config.FileName)
	}
}
