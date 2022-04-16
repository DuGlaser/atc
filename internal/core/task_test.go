package core

import (
	"testing"
)

func TestExecCode(t *testing.T) {
	tests := []struct {
		task     Task
		hasError bool
		output   string
	}{
		{
			task:     Task{RunCmd: "go run ../../test/valid.go", BuildCmd: ""},
			hasError: false,
			output:   "OK",
		},
		{
			task:     Task{RunCmd: "go run ../../test/syntax_error.go", BuildCmd: ""},
			hasError: true,
			output:   "# command-line-arguments\n../../test/syntax_error.go:6:19: syntax error: unexpected newline, expecting comma or )\n",
		},
	}

	for _, test := range tests {
		out, err := test.task.ExecCode("", false)

		if test.hasError && test.output != err.Error() {
			t.Fatalf("The error statement returned by the %s command is invalid. got=%s expect=%s",
				test.task.RunCmd, err.Error(), test.output)
		}

		if !test.hasError && out != test.output {
			t.Fatalf("The result of executing the %s command is different. got=%s expect=%s",
				test.task.RunCmd, out, test.output)
		}
	}
}

func TestBuildCode(t *testing.T) {
	tests := []struct {
		task     Task
		hasError bool
		output   string
	}{
		{
			task:     Task{RunCmd: "", BuildCmd: "go build ../../test/valid.go"},
			hasError: false,
			output:   "",
		},
		{
			task:     Task{RunCmd: "", BuildCmd: "go build ../../test/syntax_error.go"},
			hasError: true,
			output:   "# command-line-arguments\n../../test/syntax_error.go:6:19: syntax error: unexpected newline, expecting comma or )\n",
		},
	}

	for _, test := range tests {
		err := test.task.BuildCode(false)

		if test.hasError && test.output != err.Error() {
			t.Fatalf("The error statement returned by the %s command is invalid. got=%s expect=%s",
				test.task.RunCmd, err.Error(), test.output)
		}
	}
}
