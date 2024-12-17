package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	atcgModules "atcg/internal/atcg/modules"
	"atcg/pkg/utils"
)

// ParseAndGenerateTask parses module documentation and generates task YAML.
var generateTaskFunc = GenerateTask

func ParseAndGenerateTask(module string, executor atcgModules.CommandExecutor) (string, error) {
	module = strings.TrimSpace(module)

	doc, err := atcgModules.ParseModuleDoc(executor, module)
	if err != nil {
		return "", fmt.Errorf("error fetching documentation for module %s: %w", module, err)
	}

	task, err := generateTaskFunc(module, doc)
	if err != nil {
		return "", fmt.Errorf("error generating task for module %s: %w", module, err)
	}

	return task, nil
}

// WriteTaskToFile writes the task YAML to a file.
var WriteFile = os.WriteFile // Global variable for dependency injection

func WriteTaskToFile(task string, module string, outputDir string) (string, error) {
	basename := utils.Basename(module)
	outputFile := filepath.Join(outputDir, basename+".yml")

	if err := WriteFile(outputFile, []byte(task), 0644); err != nil {
		return "", fmt.Errorf("error writing task to file %s: %w", outputFile, err)
	}

	return outputFile, nil
}

// ProcessModule processes a single module by parsing documentation, generating tasks, and writing to a file.
func ProcessModule(module string, outputDir string, executor atcgModules.CommandExecutor) (*Module, error) {
	task, err := ParseAndGenerateTask(module, executor)
	if err != nil {
		return nil, err
	}

	outputFile, err := WriteTaskToFile(task, module, outputDir)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Generated task for %s: %s\n", module, outputFile)

	return &Module{
		Name:     module,
		Basename: utils.Basename(module),
	}, nil
}
