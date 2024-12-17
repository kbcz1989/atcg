package modules

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// CommandExecutor is an interface for executing commands.
type CommandExecutor interface {
	Execute(command string, args ...string) ([]byte, error)
}

// RealExecutor is the default implementation of CommandExecutor.
type RealExecutor struct{}

// Execute runs the given command and returns its output.
func (r *RealExecutor) Execute(command string, args ...string) ([]byte, error) {
	cmd := exec.Command(command, args...)
	return cmd.Output()
}

// ModuleOption represents a module option from ansible-doc output.
type ModuleOption struct {
	Choices     []string    `json:"choices,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Description []string    `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Type        string      `json:"type,omitempty"`
}

// ModuleDoc represents the structure of ansible-doc JSON output.
type ModuleDoc struct {
	Options map[string]ModuleOption `json:"options"`
}

// ParseModuleDoc runs ansible-doc and parses the JSON output for a module.
func ParseModuleDoc(exec CommandExecutor, module string) (*ModuleDoc, error) {
	output, err := exec.Execute("ansible-doc", "-j", module)
	if err != nil {
		return nil, err
	}

	var docs map[string]struct {
		Doc ModuleDoc `json:"doc"`
	}
	if err := json.Unmarshal(output, &docs); err != nil {
		return nil, err
	}

	docStruct, found := docs[module]
	if !found {
		return nil, fmt.Errorf("module %s not found in ansible-doc output", module)
	}

	return &docStruct.Doc, nil
}
