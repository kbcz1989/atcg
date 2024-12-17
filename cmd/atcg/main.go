package main

import (
	"fmt"
	"os"

	atcgModules "atcg/internal/atcg/modules"
	atcgTasks "atcg/internal/atcg/tasks"
	atcgUtils "atcg/internal/atcg/utils"

	"github.com/spf13/pflag"
)

// Run encapsulates the core logic of the main function for testing.
func Run(modules []string, outputDir string, executor atcgModules.CommandExecutor) error {
	// Input validation
	atcgUtils.ValidateInputs(modules)

	// Ensure output directory exists
	atcgUtils.EnsureOutputDirectory(outputDir)

	// Process modules
	var moduleDetails []atcgTasks.Module

	for _, module := range modules {
		result, err := atcgTasks.ProcessModule(module, outputDir, executor)
		if err != nil {
			fmt.Println(err)
			continue
		}
		moduleDetails = append(moduleDetails, *result)
	}

	// Generate main.yml
	if len(moduleDetails) > 0 {
		if err := atcgTasks.GenerateMain(moduleDetails, outputDir); err != nil {
			return fmt.Errorf("error generating main.yml: %w", err)
		}
		fmt.Printf("Generated main.yml in %s\n", outputDir)
	} else {
		fmt.Println("No valid modules processed. Skipping main.yml generation.")
	}

	return nil
}

func main() {
	var modules []string
	var outputDir string
	pflag.StringSliceVarP(&modules, "module", "m", nil, "Ansible module name (can be used multiple times)")
	pflag.StringVarP(&outputDir, "output", "o", "tasks", "Output directory for generated tasks")
	pflag.Parse()

	// Execute the core logic
	executor := &atcgModules.RealExecutor{}
	if err := Run(modules, outputDir, executor); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
