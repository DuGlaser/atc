package util

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseIntArguments(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
	}{
		{"1 2 3 4 5", []int{1, 2, 3, 4, 5}},
		{"1, 2 3 4 5", nil},
		{"1,2 3 4 5", nil},
		{"1", []int{1}},
		{"10000", nil},
	}

	for _, test := range tests {
		in := strings.Split(test.input, " ")
		result, _ := ParseIntArguments(in)

		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("result is wrong value. got=%v, expected=%v", result, test.expected)
		}
	}
}
