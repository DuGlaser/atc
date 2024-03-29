package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/DuGlaser/atc/internal/util"
)

type Task struct {
	RunCmd       string
	BuildCmd     string
	alreadyBuild bool
}

type result struct {
	Out    string
	TimeMs int64
}

func (t *Task) ExecHandleCode() (result, error) {
	var result result
	if !t.alreadyBuild {
		err := t.BuildCode()
		if err != nil {
			return result, err
		}

		fmt.Println("Ready to go!!")
	}

	util.InfoLog("Run code")

	start := time.Now()
	cmd := exec.Command("sh", "-c", t.RunCmd)
	cmd.Stdin = os.Stdin

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		return result, errors.New(stderr.String())
	}

	got := strings.TrimSpace(string(out))
	got = strings.TrimLeft(got, "\n")

	result.Out = got
	result.TimeMs = time.Since(start).Milliseconds()

	return result, nil
}

func (t *Task) ExecCode(input string) (result, error) {
	var result result
	err := t.BuildCode()
	if err != nil {
		return result, err
	}

	util.InfoLog("Run code")

	start := time.Now()
	cmd := exec.Command("sh", "-c", t.RunCmd)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr

	pipe, err := cmd.StdinPipe()
	if err != nil {
		return result, err
	}

	io.WriteString(pipe, input)
	pipe.Close()

	err = cmd.Run()

	got := strings.TrimSpace(out.String())
	got = strings.TrimLeft(got, "\n")

	result.Out = got
	result.TimeMs = time.Since(start).Milliseconds()

	if err != nil {
		return result, errors.New(stderr.String())
	}

	return result, nil
}

func (t *Task) BuildCode() error {
	if t.BuildCmd == "" {
		t.alreadyBuild = true
		return nil
	}

	if t.alreadyBuild {
		return nil
	}

	util.InfoLog("Build code")

	cmd := exec.Command("sh", "-c", t.BuildCmd)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	_, err := cmd.Output()
	if err != nil {
		return errors.New(stderr.String())
	}

	t.alreadyBuild = true

	return nil
}
