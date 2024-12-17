package tasks

import (
	"strings"
	"testing"

	atcgModules "atcg/internal/atcg/modules"
)

func TestGenerateTask_Success(t *testing.T) {
	module := "ansible.builtin.debug"
	doc := &atcgModules.ModuleDoc{
		Options: map[string]atcgModules.ModuleOption{
			"msg": {
				Default:     "Hello, World!",
				Description: []string{"Message to display."},
				Required:    false,
				Type:        "str",
			},
		},
	}

	output, err := GenerateTask(module, doc)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedOutput := `---
- name: Configure debug
  ansible.builtin.debug:
    msg: "{{ item.msg | default('Hello, World!') }}"
  tags: [debug]
`
	if strings.TrimSpace(output) != strings.TrimSpace(expectedOutput) {
		t.Errorf("unexpected output: got %q, expected %q", output, expectedOutput)
	}
}

func TestGenerateTask_NilOptions(t *testing.T) {
	module := "ansible.builtin.debug"
	doc := &atcgModules.ModuleDoc{
		Options: nil, // Invalid case
	}

	_, err := GenerateTask(module, doc)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	// Update the expected error message
	if !strings.Contains(err.Error(), "doc.Options cannot be empty") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGenerateTask_EmptyOptions(t *testing.T) {
	module := "ansible.builtin.debug"
	doc := &atcgModules.ModuleDoc{
		Options: map[string]atcgModules.ModuleOption{}, // Empty options
	}

	_, err := GenerateTask(module, doc)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "doc.Options cannot be empty") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGenerateTask_TemplateParsingError(t *testing.T) {
	originalTaskTemplate := TaskTemplate
	defer func() { TaskTemplate = originalTaskTemplate }() // Restore original after test

	// Simulate a parsing error
	TaskTemplate = `{{ define invalid-template {{ end }}`

	module := "ansible.builtin.debug"
	doc := &atcgModules.ModuleDoc{
		Options: map[string]atcgModules.ModuleOption{
			"msg": {
				Default:     "Hello, World!",
				Description: []string{"Message to display."},
				Required:    false,
				Type:        "str",
			},
		},
	}

	_, err := GenerateTask(module, doc)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "parsing task template") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGenerateTask_TemplateExecutionError(t *testing.T) {
	module := "ansible.builtin.debug"
	doc := &atcgModules.ModuleDoc{
		// Inject an invalid data structure that the template can't handle
		Options: map[string]atcgModules.ModuleOption{
			"msg": {
				Default: func() {}, // Invalid type that causes template execution to fail
			},
		},
	}

	// Call GenerateTask and expect it to fail
	_, err := GenerateTask(module, doc)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	// Validate that the error message contains the expected failure
	if !strings.Contains(err.Error(), "executing task template") {
		t.Errorf("unexpected error message: %v", err)
	}
}
