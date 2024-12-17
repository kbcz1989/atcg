package mocks

// MockExecutor is a mock implementation of the CommandExecutor interface.
type MockExecutor struct {
	// MockExecute allows you to define behavior for the Execute method in tests.
	MockExecute func(command string, args ...string) ([]byte, error)
}

// Execute calls the mock implementation provided via MockExecute.
func (m *MockExecutor) Execute(command string, args ...string) ([]byte, error) {
	if m.MockExecute != nil {
		return m.MockExecute(command, args...)
	}
	return nil, nil
}
