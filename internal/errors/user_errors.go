package errors

import (
	"fmt"
	"strings"
)

// UserError represents an error with actionable guidance for users
type UserError struct {
	Message    string
	Suggestion string
	DocsURL    string
}

func (e *UserError) Error() string {
	var b strings.Builder
	b.WriteString(e.Message)
	if e.Suggestion != "" {
		b.WriteString("\n\nSuggestion: ")
		b.WriteString(e.Suggestion)
	}
	if e.DocsURL != "" {
		b.WriteString("\nDocs: ")
		b.WriteString(e.DocsURL)
	}
	return b.String()
}

// Common error constructors

func MissingAPIKey(service string) error {
	envVar := strings.ToUpper(service) + "_API_KEY"
	return &UserError{
		Message:    fmt.Sprintf("%s API key not set", service),
		Suggestion: fmt.Sprintf("Set %s environment variable or add to ~/.config/percepta/config.yaml", envVar),
		DocsURL:    "https://github.com/Perceptax/percepta/blob/main/docs/installation.md#environment-setup",
	}
}

func DeviceNotFound(deviceID string) error {
	return &UserError{
		Message:    fmt.Sprintf("Device '%s' not found in config", deviceID),
		Suggestion: "Run 'percepta device list' to see available devices or 'percepta device add <name>' to add a new one",
		DocsURL:    "https://github.com/Perceptax/percepta/blob/main/docs/getting-started.md#step-1-configure-your-first-device",
	}
}

func InvalidBoardType(board string) error {
	return &UserError{
		Message:    fmt.Sprintf("Unknown board type '%s'", board),
		Suggestion: "Supported boards: esp32, stm32, arduino, atmega, generic. See docs for full list.",
		DocsURL:    "https://github.com/Perceptax/percepta/blob/main/docs/commands.md#supported-boards",
	}
}

func CameraNotFound(cameraPath string) error {
	return &UserError{
		Message:    fmt.Sprintf("Camera '%s' not found", cameraPath),
		Suggestion: "Check available cameras: 'ls /dev/video*' (Linux) or try camera index 0-2 (macOS/Windows)",
		DocsURL:    "https://github.com/Perceptax/percepta/blob/main/docs/troubleshooting.md#camera-not-found",
	}
}

func ConfigNotFound() error {
	return &UserError{
		Message:    "No config file found at ~/.config/percepta/config.yaml",
		Suggestion: "Create a device first with 'percepta device add <name>'",
		DocsURL:    "https://github.com/Perceptax/percepta/blob/main/docs/getting-started.md",
	}
}

func NoDevicesConfigured() error {
	return &UserError{
		Message:    "No devices configured yet",
		Suggestion: "Add your first device with 'percepta device add <name>'",
		DocsURL:    "https://github.com/Perceptax/percepta/blob/main/docs/getting-started.md#step-1-configure-your-first-device",
	}
}

func StorageInitFailed(err error) error {
	return &UserError{
		Message:    fmt.Sprintf("Failed to initialize storage: %v", err),
		Suggestion: "Check that ~/.local/share/percepta/ directory is writable",
		DocsURL:    "https://github.com/Perceptax/percepta/blob/main/docs/troubleshooting.md#storage-errors",
	}
}

func ObservationFailed(err error) error {
	return &UserError{
		Message:    fmt.Sprintf("Observation failed: %v", err),
		Suggestion: "Ensure hardware is visible to camera and properly lit. Try adjusting camera position.",
		DocsURL:    "https://github.com/Perceptax/percepta/blob/main/docs/troubleshooting.md#no-signals-detected",
	}
}

func AssertionTimeout(signal string) error {
	return &UserError{
		Message:    fmt.Sprintf("Assertion timeout: signal '%s' not found in observation", signal),
		Suggestion: "Run 'percepta observe <device>' first to see available signals",
		DocsURL:    "https://github.com/Perceptax/percepta/blob/main/docs/getting-started.md#step-4-your-first-assertion",
	}
}

func InvalidSpec(err error) error {
	return &UserError{
		Message:    fmt.Sprintf("Invalid specification: %v", err),
		Suggestion: "Check specification format. Example: 'Blink LED at 1Hz when button pressed'",
		DocsURL:    "https://github.com/Perceptax/percepta/blob/main/docs/commands.md#percepta-generate",
	}
}

func CodeGenerationFailed(err error) error {
	return &UserError{
		Message:    fmt.Sprintf("Code generation failed: %v", err),
		Suggestion: "Ensure ANTHROPIC_API_KEY is set and valid. Check your API quota.",
		DocsURL:    "https://github.com/Perceptax/percepta/blob/main/docs/troubleshooting.md#code-generation-errors",
	}
}

func StyleCheckFailed(violations int) error {
	return &UserError{
		Message:    fmt.Sprintf("Style check failed: %d violation(s) found", violations),
		Suggestion: "Run 'percepta style-check --fix <file>' to auto-fix deterministic violations",
		DocsURL:    "https://github.com/Perceptax/percepta/blob/main/docs/commands.md#percepta-style-check",
	}
}
