package handler

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/DuGlaser/atc/internal/core"
	"github.com/DuGlaser/atc/internal/repository/config"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func HaldleExec(displayID string, verbose bool) {
	cc, err := config.NewContestConfig()
	cobra.CheckErr(err)

	task, err := cc.ReadTaskSetting(displayID)
	cobra.CheckErr(err)

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

	for {
		result, err := t.ExecHandleCode(verbose)
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
