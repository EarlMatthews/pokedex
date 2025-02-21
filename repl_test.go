package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{input: "  hello  world  ", expected: []string{"hello", "world"}},
		// Add more test cases as needed
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("For input '%s', expected %v, but got %v", c.input, c.expected, actual)
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("For input '%s', at index %d, expected word '%s' but got '%s'", c.input, i, expectedWord, word)
			}
		}
	}
}
