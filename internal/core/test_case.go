package core

import "strings"

type TestCase struct {
	ID       int    `toml:"id"`
	In       string `toml:"in"`
	Expected string `toml:"expected"`
}

func (tc *TestCase) Compare(result string) bool {
	result = strings.TrimRight(result, "\n")
	expected := strings.TrimRight(tc.Expected, "\n")

	results := strings.Split(result, "\n")
	expecteds := strings.Split(expected, "\n")

	pass := true
	if len(results) != len(expecteds) {
		pass = false
	} else {
		for i := range results {
			if strings.TrimSpace(results[i]) != strings.TrimSpace(expecteds[i]) {
				pass = false
				break
			}
		}
	}

	return pass
}
