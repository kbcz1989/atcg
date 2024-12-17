package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	atcgModules "atcg/internal/atcg/modules"
)

// Module represents details for a single module.
type Module struct {
	Name     string
	Basename string
}

// Global templates for easier testing.
var TaskTemplate = `---
- name: Configure {{ .Module | basename }}
  {{ .Module }}:
{{- range $key, $option := .Options }}
    {{ $key }}: "{{ "{{ item." }}{{ $key }}{{ if $option.Required }}{{ "" }}{{ else if $option.Default }}{{ " | default('" }}{{ $option.Default }}{{ "')" }}{{ else }}{{ " | default(omit)" }}{{ end }} }}"
{{- end }}
  tags: [{{ .Module | basename }}]
`

var MainTemplate = `---
{{- range $index, $module := .Modules }}
- name: Configure {{ $module.Basename }}
  ansible.builtin.include_tasks:
    file: {{ $module.Basename }}.yml
    apply:
      tags: {{ $module.Basename }}
  when: {{ $module.Basename }} is defined
  loop: "{{ "{{ " }}{{ $module.Basename }}{{ " }}" }}"
  tags: {{ $module.Basename }}
{{ end -}}
`

// GenerateTask generates a YAML task from the module schema.
func GenerateTask(module string, doc *atcgModules.ModuleDoc) (string, error) {
	// Ensure at least one option exists
	if len(doc.Options) == 0 {
		return "", fmt.Errorf("doc.Options cannot be empty")
	}

	// Define the template function map
	funcMap := template.FuncMap{
		"replace": strings.ReplaceAll,
		"basename": func(s string) string {
			parts := strings.Split(s, ".")
			return parts[len(parts)-1]
		},
	}

	// Parse the task template
	tmpl, err := template.New("task").Funcs(funcMap).Parse(TaskTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing task template: %w", err)
	}

	// Execute the template with the provided data
	var output strings.Builder
	if err := tmpl.Execute(&output, map[string]interface{}{
		"Module":  module,
		"Options": doc.Options,
	}); err != nil {
		return "", fmt.Errorf("executing task template: %w", err)
	}

	return output.String(), nil
}

// GenerateMain generates the main.yml file with include_tasks for each module.
func GenerateMain(modules []Module, outputDir string) error {
	// Ensure modules have valid Basenames
	for _, module := range modules {
		if strings.TrimSpace(module.Basename) == "" {
			return fmt.Errorf("module.Basename cannot be empty")
		}
	}

	tmpl, err := template.New("main").Parse(MainTemplate)
	if err != nil {
		return fmt.Errorf("parsing main template: %w", err)
	}

	var output strings.Builder
	if err := tmpl.Execute(&output, map[string]interface{}{
		"Modules": modules,
	}); err != nil {
		return fmt.Errorf("executing main template: %w", err)
	}

	mainFile := filepath.Join(outputDir, "main.yml")
	if err := os.WriteFile(mainFile, []byte(output.String()), 0644); err != nil {
		return fmt.Errorf("writing main.yml: %w", err)
	}

	return nil
}
