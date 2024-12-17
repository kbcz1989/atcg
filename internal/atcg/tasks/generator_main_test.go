package tasks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateMain_Success(t *testing.T) {
	modules := []Module{
		{Name: "ansible.builtin.debug", Basename: "debug"},
		{Name: "ansible.builtin.ping", Basename: "ping"},
	}
	outputDir := t.TempDir()

	err := GenerateMain(modules, outputDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify the main.yml file
	mainFile := filepath.Join(outputDir, "main.yml")
	content, err := os.ReadFile(mainFile)
	if err != nil {
		t.Fatalf("failed to read main.yml: %v", err)
	}

	expectedContent := `---
- name: Configure debug
  ansible.builtin.include_tasks:
    file: debug.yml
    apply:
      tags: debug
  when: debug is defined
  loop: "{{ debug }}"
  tags: debug

- name: Configure ping
  ansible.builtin.include_tasks:
    file: ping.yml
    apply:
      tags: ping
  when: ping is defined
  loop: "{{ ping }}"
  tags: ping
`
	if strings.TrimSpace(string(content)) != strings.TrimSpace(expectedContent) {
		t.Errorf("unexpected content: got %q, expected %q", string(content), expectedContent)
	}
}

func TestGenerateMain_TemplateParsingError(t *testing.T) {
	originalMainTemplate := MainTemplate
	defer func() { MainTemplate = originalMainTemplate }() // Restore original after test

	// Simulate a parsing error
	MainTemplate = `{{ define invalid-template {{ end }}`

	modules := []Module{
		{Name: "ansible.builtin.debug", Basename: "debug"},
	}
	outputDir := t.TempDir()

	err := GenerateMain(modules, outputDir)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !strings.Contains(err.Error(), "parsing main template") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestGenerateMain_ErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		modules     []Module
		outputDir   string
		overrideTpl string
		wantErrMsg  string
	}{
		{
			name: "Empty Basename Error",
			modules: []Module{
				{Name: "ansible.builtin.debug", Basename: ""}, // Invalid Basename
			},
			outputDir:  t.TempDir(),
			wantErrMsg: "module.Basename cannot be empty",
		},
		{
			name: "Template Execution Error",
			modules: []Module{
				{Name: "ansible.builtin.debug", Basename: "debug"},
			},
			outputDir:   t.TempDir(),
			overrideTpl: "{{ range .Modules }}{{ .InvalidField }}{{ end }}", // Invalid template
			wantErrMsg:  "executing main template",
		},
		{
			name: "File Writing Error",
			modules: []Module{
				{Name: "ansible.builtin.debug", Basename: "debug"},
			},
			outputDir:  filepath.Join(t.TempDir(), "invalid/file_as_dir"),
			wantErrMsg: "writing main.yml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Override template if needed
			originalMainTemplate := MainTemplate
			if tt.overrideTpl != "" {
				MainTemplate = tt.overrideTpl
				defer func() { MainTemplate = originalMainTemplate }()
			}

			// Run GenerateMain
			err := GenerateMain(tt.modules, tt.outputDir)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}

			if !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("unexpected error message: got %q, want to contain %q", err.Error(), tt.wantErrMsg)
			}
		})
	}
}
