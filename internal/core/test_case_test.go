package core

import "testing"

func TestCompare(t *testing.T) {
	tests := []struct {
		got      string
		expected string
		pass     bool
	}{
		{
			got:      "test1\ntest1",
			expected: "test1\ntest1",
			pass:     true,
		},
		{
			got:      "test2\ntest2\n",
			expected: "test2\ntest2",
			pass:     true,
		},
		{
			got:      "test3\ntest3",
			expected: "test3\ntest3\n",
			pass:     true,
		},
		{
			got:      "test4  \ntest4",
			expected: "test4\ntest4",
			pass:     true,
		},
		{
			got:      "test5\t\ntest5",
			expected: "test5\ntest5",
			pass:     true,
		},
		{
			got:      "test6test6",
			expected: "test6\ntest6",
			pass:     false,
		},
		{
			got:      "\ntest7\ntest7",
			expected: "test7\ntest7",
			pass:     false,
		},
		{
			got:      "test8\ntest8\ntest8",
			expected: "test8\ntest8",
			pass:     false,
		},
	}

	for _, test := range tests {
		tc := TestCase{
			In:       "",
			Expected: test.expected,
		}

		if tc.Compare(test.got) != test.pass {
			t.Errorf("The result of the test comparison is wrong. got=%s expect=%s", test.got, test.expected)
		}
	}
}
