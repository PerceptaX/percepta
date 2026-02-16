package errors

import (
	"strings"
	"testing"
)

func TestUserError_Error(t *testing.T) {
	err := &UserError{
		Message:    "Something went wrong",
		Suggestion: "Try this fix",
		DocsURL:    "https://example.com/docs",
	}

	errMsg := err.Error()

	// Verify message is included
	if !strings.Contains(errMsg, "Something went wrong") {
		t.Errorf("Expected error message to contain 'Something went wrong', got: %s", errMsg)
	}

	// Verify suggestion is included
	if !strings.Contains(errMsg, "Suggestion: Try this fix") {
		t.Errorf("Expected error to contain suggestion, got: %s", errMsg)
	}

	// Verify docs URL is included
	if !strings.Contains(errMsg, "Docs: https://example.com/docs") {
		t.Errorf("Expected error to contain docs URL, got: %s", errMsg)
	}
}

func TestUserError_Error_NoSuggestion(t *testing.T) {
	err := &UserError{
		Message: "Something went wrong",
		DocsURL: "https://example.com/docs",
	}

	errMsg := err.Error()

	// Verify suggestion section is not included
	if strings.Contains(errMsg, "Suggestion:") {
		t.Errorf("Expected no suggestion section, got: %s", errMsg)
	}

	// Verify docs URL is still included
	if !strings.Contains(errMsg, "Docs:") {
		t.Errorf("Expected docs URL to be included, got: %s", errMsg)
	}
}

func TestUserError_Error_NoDocsURL(t *testing.T) {
	err := &UserError{
		Message:    "Something went wrong",
		Suggestion: "Try this fix",
	}

	errMsg := err.Error()

	// Verify docs section is not included
	if strings.Contains(errMsg, "Docs:") {
		t.Errorf("Expected no docs section, got: %s", errMsg)
	}

	// Verify suggestion is still included
	if !strings.Contains(errMsg, "Suggestion:") {
		t.Errorf("Expected suggestion to be included, got: %s", errMsg)
	}
}

func TestMissingAPIKey(t *testing.T) {
	tests := []struct {
		service     string
		expectedVar string
		expectedMsg string
	}{
		{"ANTHROPIC", "ANTHROPIC_API_KEY", "ANTHROPIC API key not set"},
		{"OPENAI", "OPENAI_API_KEY", "OPENAI API key not set"},
		{"claude", "CLAUDE_API_KEY", "claude API key not set"},
	}

	for _, tt := range tests {
		t.Run(tt.service, func(t *testing.T) {
			err := MissingAPIKey(tt.service)
			errMsg := err.Error()

			// Verify message contains service name
			if !strings.Contains(errMsg, tt.expectedMsg) {
				t.Errorf("Expected message to contain '%s', got: %s", tt.expectedMsg, errMsg)
			}

			// Verify suggestion mentions the environment variable
			if !strings.Contains(errMsg, tt.expectedVar) {
				t.Errorf("Expected suggestion to mention '%s', got: %s", tt.expectedVar, errMsg)
			}

			// Verify docs URL is present
			if !strings.Contains(errMsg, "Docs:") {
				t.Errorf("Expected docs URL, got: %s", errMsg)
			}

			// Verify suggestion mentions config file
			if !strings.Contains(errMsg, "config.yaml") {
				t.Errorf("Expected suggestion to mention config.yaml, got: %s", errMsg)
			}
		})
	}
}

func TestDeviceNotFound(t *testing.T) {
	err := DeviceNotFound("my-esp32")
	errMsg := err.Error()

	// Verify message contains device ID
	if !strings.Contains(errMsg, "my-esp32") {
		t.Errorf("Expected message to contain 'my-esp32', got: %s", errMsg)
	}

	// Verify suggestion mentions device list command
	if !strings.Contains(errMsg, "percepta device list") {
		t.Errorf("Expected suggestion to mention 'percepta device list', got: %s", errMsg)
	}

	// Verify suggestion mentions device add command
	if !strings.Contains(errMsg, "percepta device add") {
		t.Errorf("Expected suggestion to mention 'percepta device add', got: %s", errMsg)
	}

	// Verify docs URL is present
	if !strings.Contains(errMsg, "Docs:") {
		t.Errorf("Expected docs URL, got: %s", errMsg)
	}
}

func TestInvalidBoardType(t *testing.T) {
	err := InvalidBoardType("unknown-board")
	errMsg := err.Error()

	// Verify message contains board name
	if !strings.Contains(errMsg, "unknown-board") {
		t.Errorf("Expected message to contain 'unknown-board', got: %s", errMsg)
	}

	// Verify suggestion lists supported boards
	supportedBoards := []string{"esp32", "stm32", "arduino", "atmega", "generic"}
	for _, board := range supportedBoards {
		if !strings.Contains(errMsg, board) {
			t.Errorf("Expected suggestion to mention supported board '%s', got: %s", board, errMsg)
		}
	}

	// Verify docs URL is present
	if !strings.Contains(errMsg, "Docs:") {
		t.Errorf("Expected docs URL, got: %s", errMsg)
	}
}

func TestCameraNotFound(t *testing.T) {
	err := CameraNotFound("/dev/video0")
	errMsg := err.Error()

	// Verify message contains camera path
	if !strings.Contains(errMsg, "/dev/video0") {
		t.Errorf("Expected message to contain '/dev/video0', got: %s", errMsg)
	}

	// Verify suggestion mentions Linux command
	if !strings.Contains(errMsg, "ls /dev/video*") {
		t.Errorf("Expected suggestion to mention Linux camera listing, got: %s", errMsg)
	}

	// Verify suggestion mentions camera index for other platforms
	if !strings.Contains(errMsg, "camera index") {
		t.Errorf("Expected suggestion to mention camera index for other platforms, got: %s", errMsg)
	}

	// Verify docs URL is present
	if !strings.Contains(errMsg, "Docs:") {
		t.Errorf("Expected docs URL, got: %s", errMsg)
	}

	// Verify links to troubleshooting
	if !strings.Contains(errMsg, "troubleshooting") {
		t.Errorf("Expected docs URL to link to troubleshooting, got: %s", errMsg)
	}
}

func TestConfigNotFound(t *testing.T) {
	err := ConfigNotFound()
	errMsg := err.Error()

	// Verify message mentions config path
	if !strings.Contains(errMsg, "~/.config/percepta/config.yaml") {
		t.Errorf("Expected message to mention config path, got: %s", errMsg)
	}

	// Verify suggestion mentions device add command
	if !strings.Contains(errMsg, "percepta device add") {
		t.Errorf("Expected suggestion to mention 'percepta device add', got: %s", errMsg)
	}

	// Verify docs URL is present
	if !strings.Contains(errMsg, "Docs:") {
		t.Errorf("Expected docs URL, got: %s", errMsg)
	}
}

func TestNoDevicesConfigured(t *testing.T) {
	err := NoDevicesConfigured()
	errMsg := err.Error()

	// Verify message mentions no devices
	if !strings.Contains(errMsg, "No devices configured") {
		t.Errorf("Expected message about no devices, got: %s", errMsg)
	}

	// Verify suggestion mentions device add command
	if !strings.Contains(errMsg, "percepta device add") {
		t.Errorf("Expected suggestion to mention 'percepta device add', got: %s", errMsg)
	}

	// Verify docs URL is present
	if !strings.Contains(errMsg, "Docs:") {
		t.Errorf("Expected docs URL, got: %s", errMsg)
	}
}

func TestStorageInitFailed(t *testing.T) {
	originalErr := &UserError{Message: "database locked"}
	err := StorageInitFailed(originalErr)
	errMsg := err.Error()

	// Verify message mentions storage init failure
	if !strings.Contains(errMsg, "Failed to initialize storage") {
		t.Errorf("Expected message about storage init failure, got: %s", errMsg)
	}

	// Verify original error is included
	if !strings.Contains(errMsg, "database locked") {
		t.Errorf("Expected original error to be included, got: %s", errMsg)
	}

	// Verify suggestion mentions storage directory
	if !strings.Contains(errMsg, ".local/share/percepta") {
		t.Errorf("Expected suggestion to mention storage directory, got: %s", errMsg)
	}

	// Verify docs URL is present
	if !strings.Contains(errMsg, "Docs:") {
		t.Errorf("Expected docs URL, got: %s", errMsg)
	}
}

func TestObservationFailed(t *testing.T) {
	originalErr := &UserError{Message: "camera timeout"}
	err := ObservationFailed(originalErr)
	errMsg := err.Error()

	// Verify message mentions observation failure
	if !strings.Contains(errMsg, "Observation failed") {
		t.Errorf("Expected message about observation failure, got: %s", errMsg)
	}

	// Verify original error is included
	if !strings.Contains(errMsg, "camera timeout") {
		t.Errorf("Expected original error to be included, got: %s", errMsg)
	}

	// Verify suggestion mentions hardware visibility
	if !strings.Contains(errMsg, "camera") {
		t.Errorf("Expected suggestion to mention camera, got: %s", errMsg)
	}

	// Verify docs URL is present
	if !strings.Contains(errMsg, "Docs:") {
		t.Errorf("Expected docs URL, got: %s", errMsg)
	}
}

func TestAssertionTimeout(t *testing.T) {
	err := AssertionTimeout("LED.POWER")
	errMsg := err.Error()

	// Verify message mentions the signal
	if !strings.Contains(errMsg, "LED.POWER") {
		t.Errorf("Expected message to contain signal 'LED.POWER', got: %s", errMsg)
	}

	// Verify message mentions timeout
	if !strings.Contains(errMsg, "timeout") {
		t.Errorf("Expected message to mention timeout, got: %s", errMsg)
	}

	// Verify suggestion mentions observe command
	if !strings.Contains(errMsg, "percepta observe") {
		t.Errorf("Expected suggestion to mention 'percepta observe', got: %s", errMsg)
	}

	// Verify docs URL is present
	if !strings.Contains(errMsg, "Docs:") {
		t.Errorf("Expected docs URL, got: %s", errMsg)
	}
}

func TestInvalidSpec(t *testing.T) {
	originalErr := &UserError{Message: "unexpected token"}
	err := InvalidSpec(originalErr)
	errMsg := err.Error()

	// Verify message mentions invalid specification
	if !strings.Contains(errMsg, "Invalid specification") {
		t.Errorf("Expected message about invalid specification, got: %s", errMsg)
	}

	// Verify original error is included
	if !strings.Contains(errMsg, "unexpected token") {
		t.Errorf("Expected original error to be included, got: %s", errMsg)
	}

	// Verify suggestion provides example
	if !strings.Contains(errMsg, "Example:") {
		t.Errorf("Expected suggestion to provide example, got: %s", errMsg)
	}

	// Verify example mentions LED blink
	if !strings.Contains(errMsg, "Blink LED") {
		t.Errorf("Expected example to mention LED blink, got: %s", errMsg)
	}

	// Verify docs URL is present
	if !strings.Contains(errMsg, "Docs:") {
		t.Errorf("Expected docs URL, got: %s", errMsg)
	}
}

func TestCodeGenerationFailed(t *testing.T) {
	originalErr := &UserError{Message: "API rate limit exceeded"}
	err := CodeGenerationFailed(originalErr)
	errMsg := err.Error()

	// Verify message mentions code generation failure
	if !strings.Contains(errMsg, "Code generation failed") {
		t.Errorf("Expected message about code generation failure, got: %s", errMsg)
	}

	// Verify original error is included
	if !strings.Contains(errMsg, "API rate limit exceeded") {
		t.Errorf("Expected original error to be included, got: %s", errMsg)
	}

	// Verify suggestion mentions API key
	if !strings.Contains(errMsg, "ANTHROPIC_API_KEY") {
		t.Errorf("Expected suggestion to mention ANTHROPIC_API_KEY, got: %s", errMsg)
	}

	// Verify suggestion mentions API quota
	if !strings.Contains(errMsg, "quota") {
		t.Errorf("Expected suggestion to mention quota, got: %s", errMsg)
	}

	// Verify docs URL is present
	if !strings.Contains(errMsg, "Docs:") {
		t.Errorf("Expected docs URL, got: %s", errMsg)
	}
}

func TestStyleCheckFailed(t *testing.T) {
	tests := []struct {
		name       string
		violations int
	}{
		{"single violation", 1},
		{"multiple violations", 5},
		{"many violations", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StyleCheckFailed(tt.violations)
			errMsg := err.Error()

			// Verify message mentions style check failure
			if !strings.Contains(errMsg, "Style check failed") {
				t.Errorf("Expected message about style check failure, got: %s", errMsg)
			}

			// Verify message includes violation count
			if !strings.Contains(errMsg, "violation") {
				t.Errorf("Expected message to mention violations, got: %s", errMsg)
			}

			// Verify suggestion mentions fix flag
			if !strings.Contains(errMsg, "--fix") {
				t.Errorf("Expected suggestion to mention --fix flag, got: %s", errMsg)
			}

			// Verify suggestion mentions style-check command
			if !strings.Contains(errMsg, "percepta style-check") {
				t.Errorf("Expected suggestion to mention style-check command, got: %s", errMsg)
			}

			// Verify docs URL is present
			if !strings.Contains(errMsg, "Docs:") {
				t.Errorf("Expected docs URL, got: %s", errMsg)
			}
		})
	}
}

func TestUserError_ImplementsError(t *testing.T) {
	var err error = &UserError{Message: "test"}
	if err.Error() != "test" {
		t.Errorf("Expected UserError to implement error interface properly")
	}
}

func TestAllErrors_HaveDocsURL(t *testing.T) {
	// Test that all error constructors return errors with docs URLs
	errors := []error{
		MissingAPIKey("test"),
		DeviceNotFound("test"),
		InvalidBoardType("test"),
		CameraNotFound("test"),
		ConfigNotFound(),
		NoDevicesConfigured(),
		StorageInitFailed(&UserError{Message: "test"}),
		ObservationFailed(&UserError{Message: "test"}),
		AssertionTimeout("test"),
		InvalidSpec(&UserError{Message: "test"}),
		CodeGenerationFailed(&UserError{Message: "test"}),
		StyleCheckFailed(1),
	}

	for i, err := range errors {
		errMsg := err.Error()
		if !strings.Contains(errMsg, "Docs:") {
			t.Errorf("Error #%d missing docs URL: %s", i, errMsg)
		}
		if !strings.Contains(errMsg, "github.com/Perceptax/percepta") {
			t.Errorf("Error #%d docs URL doesn't point to correct repository: %s", i, errMsg)
		}
	}
}

func TestAllErrors_HaveSuggestion(t *testing.T) {
	// Test that all error constructors return errors with suggestions
	errors := []error{
		MissingAPIKey("test"),
		DeviceNotFound("test"),
		InvalidBoardType("test"),
		CameraNotFound("test"),
		ConfigNotFound(),
		NoDevicesConfigured(),
		StorageInitFailed(&UserError{Message: "test"}),
		ObservationFailed(&UserError{Message: "test"}),
		AssertionTimeout("test"),
		InvalidSpec(&UserError{Message: "test"}),
		CodeGenerationFailed(&UserError{Message: "test"}),
		StyleCheckFailed(1),
	}

	for i, err := range errors {
		errMsg := err.Error()
		if !strings.Contains(errMsg, "Suggestion:") {
			t.Errorf("Error #%d missing suggestion: %s", i, errMsg)
		}
	}
}
