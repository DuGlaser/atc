package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
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

func (t *Task) ExecCode(input string, verbose bool) (result, error) {
	var result result
	err := t.BuildCode(verbose)
	if err != nil {
		return result, err
	}

	if verbose {
		fmt.Println("Run code...")
	}

	start := time.Now()
	cmd := exec.Command("sh", "-c", t.RunCmd)
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return result, err
	}

	io.WriteString(pipe, input)
	pipe.Close()

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

func (t *Task) BuildCode(verbose bool) error {
	if t.BuildCmd == "" {
		return nil
	}

	if t.alreadyBuild {
		return nil
	}

	if verbose {
		fmt.Println("Build code...")
	}

	cmd := exec.Command("sh", "-c", t.BuildCmd)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	_, err := cmd.Output()
	if err != nil {
		return errors.New(stderr.String())
	}

	if verbose {
		fmt.Println(stderr.String())
	}

	t.alreadyBuild = true

	return nil
}
