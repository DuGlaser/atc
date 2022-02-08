package fetcher

import "testing"

func TestGetAtcoderUrl(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{input: []string{"login"}, expected: "https://atcoder.jp/login?lang=ja"},
		{input: []string{"contests", "abc001"}, expected: "https://atcoder.jp/contests/abc001?lang=ja"},
		{input: []string{"contests", "abc001", "submit"}, expected: "https://atcoder.jp/contests/abc001/submit?lang=ja"},
	}

	for _, test := range tests {
		u := GetAtcoderUrl(test.input...)
		if test.expected != u {
			t.Errorf("GetAtcoderUrl has wrong value. got=%v want=%v.", u, test.expected)
		}
	}
}
