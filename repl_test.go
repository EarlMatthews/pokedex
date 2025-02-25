package main

import (
	"testing"
	"reflect"
)

func TestCleanInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Normal input",
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			name:     "Multiple spaces between words",
			input:    "hello   world",
			expected: []string{"hello", "world"},
		},
		{
			name:     "Leading and trailing spaces",
			input:    "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			name:     "Empty input",
			input:    "",
			expected: []string{},
		},
		{
			name:     "All whitespace characters",
			input:    "\t\n\r\f ",
			expected: []string{},
		},
		{
			name:     "Single word",
			input:    "test",
			expected: []string{"test"},
		},
		{
			name:     "Mixed whitespace characters",
			input:    "hello\tworld\nanother\r\fword",
			expected: []string{"hello", "world", "another", "word"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := cleanInput(tt.input)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("cleanInput(%q) = %v, expected %v", tt.input, actual, tt.expected)
			}
		})
	}
}