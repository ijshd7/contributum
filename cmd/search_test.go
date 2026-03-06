package cmd

import (
	"reflect"
	"testing"
)

func TestSplitCSV(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "single value",
			input:    "go",
			expected: []string{"go"},
		},
		{
			name:     "multiple values",
			input:    "python,javascript,rust",
			expected: []string{"python", "javascript", "rust"},
		},
		{
			name:     "values with spaces",
			input:    "go , python , rust",
			expected: []string{"go", "python", "rust"},
		},
		{
			name:     "values with leading spaces",
			input:    "  go,  python,  rust",
			expected: []string{"go", "python", "rust"},
		},
		{
			name:     "values with trailing spaces",
			input:    "go  ,python  ,rust  ",
			expected: []string{"go", "python", "rust"},
		},
		{
			name:     "mixed spacing",
			input:    "  go , python  ,  rust  ",
			expected: []string{"go", "python", "rust"},
		},
		{
			name:     "single value with spaces",
			input:    "  golang  ",
			expected: []string{"golang"},
		},
		{
			name:     "empty parts are skipped",
			input:    "go,,python",
			expected: []string{"go", "python"},
		},
		{
			name:     "only commas",
			input:    ",,",
			expected: []string{},
		},
		{
			name:     "comma at start",
			input:    ",go,python",
			expected: []string{"go", "python"},
		},
		{
			name:     "comma at end",
			input:    "go,python,",
			expected: []string{"go", "python"},
		},
		{
			name:     "spaces only",
			input:    "   ,   ,   ",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitCSV(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("splitCSV(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
