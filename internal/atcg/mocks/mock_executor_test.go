package mocks

import (
	"errors"
	"testing"
)

func TestMockExecutor_WithMockExecute(t *testing.T) {
	mockOutput := []byte("mock output")
	mockError := errors.New("mock error")

	executor := &MockExecutor{
		MockExecute: func(command string, args ...string) ([]byte, error) {
			// Validate inputs
			if command != "ansible-doc" {
				t.Errorf("unexpected command: got %q, want %q", command, "ansible-doc")
			}
			if len(args) != 2 || args[0] != "-j" || args[1] != "ping" {
				t.Errorf("unexpected arguments: got %v, want %v", args, []string{"-j", "ping"})
			}
			return mockOutput, mockError
		},
	}

	// Call Execute and validate
	output, err := executor.Execute("ansible-doc", "-j", "ping")
	if string(output) != string(mockOutput) {
		t.Errorf("unexpected output: got %q, want %q", output, mockOutput)
	}
	if err != mockError {
		t.Errorf("unexpected error: got %v, want %v", err, mockError)
	}
}

func TestMockExecutor_WithoutMockExecute(t *testing.T) {
	executor := &MockExecutor{
		MockExecute: nil, // No mock behavior defined
	}

	// Call Execute and expect default nil behavior
	output, err := executor.Execute("ansible-doc", "-j", "ping")
	if output != nil {
		t.Errorf("expected output to be nil, got %q", output)
	}
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
}
