package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Task struct {
	RunCmd       string
	BuildCmd     string
	alreadyBuild bool
}

func (t *Task) ExecCode(input string, verbose bool) (string, error) {
	err := t.BuildCode(verbose)
	if err != nil {
		return "", err
	}

	if verbose {
		fmt.Println("Run code...")
	}

	cmd := exec.Command("sh", "-c", t.RunCmd)
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	io.WriteString(pipe, input)
	pipe.Close()

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		return "", errors.New(stderr.String())
	}

	got := strings.TrimSpace(string(out))
	got = strings.TrimLeft(got, "\n")

	return got, nil
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

	t.alreadyBuild = true

	return nil
}
