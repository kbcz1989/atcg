package utils

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

// ValidateInputs checks if the modules slice is empty.
func ValidateInputs(modules []string) {
	if len(modules) == 0 {
		fmt.Println("Error: No modules specified. Use -m or --module to specify modules.")
		pflag.Usage()
		os.Exit(1)
	}
}

// EnsureOutputDirectory creates the output directory if it doesn't exist.
func EnsureOutputDirectory(outputDir string) {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Printf("Error creating output directory '%s': %v\n", outputDir, err)
		os.Exit(1)
	}
}
