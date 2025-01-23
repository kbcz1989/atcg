package tasks

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"atcg/internal/atcg/mocks"
	atcgModules "atcg/internal/atcg/modules"
)

func TestParseAndGenerateTask_Success(t *testing.T) {
	mockExecutor := &mocks.MockExecutor{
		MockExecute: func(command string, args ...string) ([]byte, error) {
			return []byte(`{"ansible.builtin.debug": {"doc": {"options": {"msg": {"default": "Hello, World!", "required": false}}}}}`), nil
		},
	}

	task, err := ParseAndGenerateTask("ansible.builtin.debug", mockExecutor)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := `---
- name: Configure debug
  ansible.builtin.debug:
    msg: "{{ item.msg | default('Hello, World!') }}"
  tags: [debug]
`
	if task != expected {
		t.Errorf("unexpected task output, got:\n%s\nexpected:\n%s", task, expected)
	}
}
func TestParseAndGenerateTask_ParseError(t *testing.T) {
	// Mock executor to simulate an error during command execution
	expectedErr := fmt.Errorf("failed to execute ansible-doc")
	mockExecutor := &mocks.MockExecutor{
		MockExecute: func(command string, args ...string) ([]byte, error) {
			return nil, expectedErr
		},
	}

	// Call ParseAndGenerateTask
	_, err := ParseAndGenerateTask("ansible.builtin.debug", mockExecutor)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	// Verify that the error contains the expected underlying error
	if !errors.Is(err, expectedErr) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWriteTaskToFile_Success(t *testing.T) {
	taskContent := "---\n- name: Sample task\n"
	module := "ansible.builtin.debug"
	outputDir := t.TempDir()

	outputFile, err := WriteTaskToFile(taskContent, module, outputDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedFile := filepath.Join(outputDir, "debug.yml")
	if outputFile != expectedFile {
		t.Errorf("unexpected file path, got %q, expected %q", outputFile, expectedFile)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}

	if string(content) != taskContent {
		t.Errorf("unexpected file content, got %q, expected %q", string(content), taskContent)
	}
}

func TestProcessModule_Success(t *testing.T) {
	mockExecutor := &mocks.MockExecutor{
		MockExecute: func(command string, args ...string) ([]byte, error) {
			return []byte(`{"ansible.builtin.debug": {"doc": {"options": {"msg": {"default": "Hello, World!", "required": false}}}}}`), nil
		},
	}

	outputDir := t.TempDir()
	result, err := ProcessModule("ansible.builtin.debug", outputDir, mockExecutor)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Basename != "debug" {
		t.Errorf("unexpected basename, got %q, expected %q", result.Basename, "debug")
	}

	expectedFile := filepath.Join(outputDir, "debug.yml")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("expected file %q to exist, but it does not", expectedFile)
	}
}

func TestProcessModule_ParseError(t *testing.T) {
	mockExecutor := &mocks.MockExecutor{
		MockExecute: func(command string, args ...string) ([]byte, error) {
			return nil, fmt.Errorf("failed to execute ansible-doc")
		},
	}

	outputDir := t.TempDir()
	result, err := ProcessModule("ansible.builtin.debug", outputDir, mockExecutor)

	if result != nil {
		t.Fatalf("expected result to be nil, got %v", result)
	}

	if err == nil || !strings.Contains(err.Error(), "error fetching documentation") {
		t.Errorf("unexpected error: got %v, expected error containing 'error fetching documentation'", err)
	}
}

func TestProcessModule_WriteError(t *testing.T) {
	// Save original WriteFile and restore after the test
	originalWriteFile := WriteFile
	defer func() { WriteFile = originalWriteFile }()

	// Mock WriteFile to simulate a failure
	WriteFile = func(name string, data []byte, perm os.FileMode) error {
		return fmt.Errorf("mock write error")
	}

	mockExecutor := &mocks.MockExecutor{
		MockExecute: func(command string, args ...string) ([]byte, error) {
			return []byte(`{"ansible.builtin.debug": {"doc": {"options": {"msg": {"default": "Hello"}}}}}`), nil
		},
	}

	outputDir := t.TempDir()
	result, err := ProcessModule("ansible.builtin.debug", outputDir, mockExecutor)

	if result != nil {
		t.Fatalf("expected result to be nil, got %v", result)
	}

	if err == nil || !strings.Contains(err.Error(), "error writing task to file") {
		t.Errorf("unexpected error: got %v, expected error containing 'error writing task to file'", err)
	}
}

func TestParseAndGenerateTask_GenerateTaskError(t *testing.T) {
	// Save and restore the original GenerateTask function
	originalGenerateTask := generateTaskFunc
	defer func() { generateTaskFunc = originalGenerateTask }()

	// Mock GenerateTask to return an error
	generateTaskFunc = func(module string, doc *atcgModules.ModuleDoc) (string, error) {
		return "", fmt.Errorf("mock generate task error")
	}

	mockExecutor := &mocks.MockExecutor{
		MockExecute: func(command string, args ...string) ([]byte, error) {
			// Return valid JSON to simulate successful ParseModuleDoc
			return []byte(`{"ansible.builtin.debug": {"doc": {"options": {"msg": {"default": "Hello"}}}}}`), nil
		},
	}

	module := "ansible.builtin.debug"
	output, err := ParseAndGenerateTask(module, mockExecutor)

	// Validate that an error occurred
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "error generating task for module ansible.builtin.debug") {
		t.Errorf("unexpected error message: %v", err)
	}

	// Ensure output is empty
	if output != "" {
		t.Errorf("expected output to be empty, got %q", output)
	}
}
