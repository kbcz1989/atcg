package modules

import (
	"errors"
	"os/exec"
	"strings"
	"testing"
)

// MockExecutor is a mock implementation of CommandExecutor.
type MockExecutor struct {
	OutputFunc func(command string, args ...string) ([]byte, error)
}

// Execute runs the mocked command and returns the predefined output.
func (m *MockExecutor) Execute(command string, args ...string) ([]byte, error) {
	if m.OutputFunc != nil {
		return m.OutputFunc(command, args...)
	}
	return nil, errors.New("mock executor not configured")
}

func TestParseModuleDoc_Success(t *testing.T) {
	// Mock output for `ansible-doc`
	mockOutput := `{
		"ansible.windows.win_user_right": {
			"doc": {
				"options": {
					"action": {
						"choices": ["add", "remove", "set"],
						"default": "set",
						"description": ["The action to take."],
						"required": false,
						"type": "str"
					},
					"name": {
						"description": ["The name of the user right."],
						"required": true,
						"type": "str"
					}
				}
			}
		}
	}`

	// Mock CommandExecutor
	executor := &MockExecutor{
		OutputFunc: func(command string, args ...string) ([]byte, error) {
			return []byte(mockOutput), nil
		},
	}

	module := "ansible.windows.win_user_right"
	doc, err := ParseModuleDoc(executor, module)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Assert expected options
	if len(doc.Options) != 2 {
		t.Errorf("expected 2 options, got %d", len(doc.Options))
	}

	action, ok := doc.Options["action"]
	if !ok {
		t.Fatalf("expected 'action' option, but it was missing")
	}
	if action.Default != "set" {
		t.Errorf("expected default 'set', got %v", action.Default)
	}
	if action.Type != "str" {
		t.Errorf("expected type 'str', got %v", action.Type)
	}
}

func TestParseModuleDoc_Error(t *testing.T) {
	// Mock error scenario
	expectedErr := errors.New("command not found")
	executor := &MockExecutor{
		OutputFunc: func(command string, args ...string) ([]byte, error) {
			return nil, expectedErr
		},
	}

	module := "non_existent_module"
	_, err := ParseModuleDoc(executor, module)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	// Print the error for debugging
	t.Logf("error: %v", err)

	// Verify that the error contains the expected underlying error
	if !errors.Is(err, expectedErr) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRealExecutor_Execute_Success(t *testing.T) {
	executor := &RealExecutor{}

	// Execute a simple command
	output, err := executor.Execute("echo", "hello")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "hello\n" // echo adds a newline
	if string(output) != expected {
		t.Errorf("expected output %q, got %q", expected, string(output))
	}
}

func TestRealExecutor_Execute_Error(t *testing.T) {
	executor := &RealExecutor{}

	// Execute a non-existent command
	_, err := executor.Execute("nonexistent-command")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	// Verify that the error contains "executable file not found"
	if !errors.Is(err, exec.ErrNotFound) && !strings.Contains(err.Error(), "executable file not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseModuleDoc_InvalidJSON(t *testing.T) {
	// Mock executor to return invalid JSON
	executor := &MockExecutor{
		OutputFunc: func(command string, args ...string) ([]byte, error) {
			return []byte("{invalid-json}"), nil
		},
	}

	_, err := ParseModuleDoc(executor, "some_module")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	// Verify that the error is related to JSON unmarshalling
	if !strings.Contains(err.Error(), "invalid character") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestParseModuleDoc_EmptyModuleName(t *testing.T) {
	// Mock executor to return valid output
	mockOutput := `{
		"some_module": {
			"doc": {
				"options": {}
			}
		}
	}`
	executor := &MockExecutor{
		OutputFunc: func(command string, args ...string) ([]byte, error) {
			return []byte(mockOutput), nil
		},
	}

	// Call ParseModuleDoc with an empty module name
	_, err := ParseModuleDoc(executor, "")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	// Verify that the error is about the module not being found
	if err.Error() != "module  not found in ansible-doc output" {
		t.Errorf("unexpected error: %v", err)
	}
}
