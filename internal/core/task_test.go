package core

import (
	"testing"
)

func TestExecCode(t *testing.T) {
	tests := []struct {
		task     Task
		hasError bool
	}{
		{
			task:     Task{RunCmd: "go run ../../test/valid.go", BuildCmd: ""},
			hasError: false,
		},
		{
			task:     Task{RunCmd: "go run ../../test/syntax_error.go", BuildCmd: ""},
			hasError: true,
		},
	}

	for _, test := range tests {
		_, err := test.task.ExecCode("", false)

		if test.hasError && err == nil {
			t.Fatalf("No error occurd. commad=`%s`", test.task.RunCmd)
		}

		if !test.hasError && err != nil {
			t.Fatalf("Error occurd. commad=`%s`", test.task.RunCmd)
		}
	}
}

func TestBuildCode(t *testing.T) {
	tests := []struct {
		task     Task
		hasError bool
	}{
		{
			task:     Task{RunCmd: "", BuildCmd: "go build ../../test/valid.go && rm valid"},
			hasError: false,
		},
		{
			task:     Task{RunCmd: "", BuildCmd: "go build ../../test/syntax_error.go"},
			hasError: true,
		},
	}

	for _, test := range tests {
		err := test.task.BuildCode(false)

		if test.hasError && err == nil {
			t.Fatalf("No error occurd. commad=`%s`", test.task.RunCmd)
		}

		if !test.hasError && err != nil {
			t.Fatalf("Error occurd. commad=`%s`", test.task.RunCmd)
		}
	}
}
