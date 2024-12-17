package utils

import "testing"

func TestBasename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid module name with dots",
			input:    "ansible.builtin.debug",
			expected: "debug",
		},
		{
			name:     "Single segment module name",
			input:    "debug",
			expected: "debug",
		},
		{
			name:     "Empty string input",
			input:    "",
			expected: "",
		},
		{
			name:     "Trailing dot in module name",
			input:    "ansible.builtin.",
			expected: "",
		},
		{
			name:     "Only dots in module name",
			input:    "...",
			expected: "",
		},
		{
			name:     "Module name with leading dots",
			input:    ".ansible.builtin.debug",
			expected: "debug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Basename(tt.input)
			if result != tt.expected {
				t.Errorf("Basename(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
