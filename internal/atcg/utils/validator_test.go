package utils

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestValidateInputs_ValidModules(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()

	ValidateInputs([]string{"module1", "module2"}) // Should not exit
}

func TestEnsureOutputDirectory_Success(t *testing.T) {
	tempDir := os.TempDir() + "/test_output_dir"
	defer os.RemoveAll(tempDir) // Cleanup

	EnsureOutputDirectory(tempDir)

	// Verify directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Fatalf("expected directory %s to exist, but it does not", tempDir)
	}
}

// Helper function to run a test function in a subprocess and capture its output.
func runSubprocessTest(t *testing.T, testFuncName string) (string, error) {
	cmd := exec.Command(os.Args[0], "-test.run="+testFuncName)
	cmd.Env = append(os.Environ(), "SUBPROCESS_TEST=1")

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	err := cmd.Run()
	return output.String(), err
}

func TestValidateInputs_EmptyModules(t *testing.T) {
	if os.Getenv("SUBPROCESS_TEST") == "1" {
		ValidateInputs([]string{}) // This should call os.Exit
		return
	}

	output, err := runSubprocessTest(t, "TestValidateInputs_EmptyModules")

	if err == nil || !strings.Contains(output, "Error: No modules specified") {
		t.Fatalf("expected os.Exit with error, got: %v, output: %s", err, output)
	}
}

func TestEnsureOutputDirectory_Error(t *testing.T) {
	if os.Getenv("SUBPROCESS_TEST") == "1" {
		EnsureOutputDirectory("") // Invalid path to trigger os.Exit
		return
	}

	output, err := runSubprocessTest(t, "TestEnsureOutputDirectory_Error")

	if err == nil || !strings.Contains(output, "Error creating output directory") {
		t.Fatalf("expected os.Exit with error, got: %v, output: %s", err, output)
	}
}
