package core

import "strings"

type TestCase struct {
	In       string
	Expected string
}

func (tc *TestCase) Compare(result string) bool {
	results := strings.Split(result, "\n")
	expecteds := strings.Split(tc.Expected, "\n")

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
