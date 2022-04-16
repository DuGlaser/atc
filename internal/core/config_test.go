package core

import (
	"fmt"
	"testing"
)

func createEmptyError(key string) error {
	return fmt.Errorf("Config.%s is empty.", key)
}

func TestValidate(t *testing.T) {
	tests := []struct {
		input  Config
		expect error
	}{
		{input: Config{Cmd: "", Lang: "", FileName: "", Template: ""}, expect: createEmptyError("Cmd")},
		{input: Config{Cmd: "Cmd", Lang: "", FileName: "", Template: ""}, expect: createEmptyError("Lang")},
		{input: Config{Cmd: "Cmd", Lang: "Lang", FileName: "", Template: ""}, expect: createEmptyError("FileName")},
		{input: Config{Cmd: "Cmd", Lang: "Lang", FileName: "FileName", Template: ""}, expect: nil},
	}

	for _, test := range tests {
		err := test.input.Validate()
		if err == nil && test.expect == nil {
			continue
		}

		if !(err.Error() == test.expect.Error()) {
			t.Errorf("config.Validate() return wrong value. got=%v expect=%v.", err, test.expect)
		}
	}
}
