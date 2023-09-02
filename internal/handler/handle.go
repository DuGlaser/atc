package handler

import (
	"fmt"

	"github.com/DuGlaser/atc/internal/core"
	"github.com/DuGlaser/atc/internal/repository/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func HaldleExec(displayID string) {
	cc, err := config.NewContestConfig()
	cobra.CheckErr(err)

	task, err := cc.ReadTaskSetting(displayID)
	cobra.CheckErr(err)

	config, err := cc.ReadConfig()
	cobra.CheckErr(err)

	cobra.CheckErr(config.GenerateCmd(task.Path, config.FileName))

	t := &core.Task{
		RunCmd:   config.RunCmd,
		BuildCmd: config.BuildCmd,
	}

	for {
		result, err := t.ExecHandleCode()
		cobra.CheckErr(err)
		fmt.Println()

		logs := []struct {
			label   string
			content string
		}{
			{label: color.CyanString("your output: "), content: result.Out},
		}

		for _, log := range logs {
			fmt.Println(log.label)
			printResultValue(log.content)
			fmt.Println()
		}
	}
}
