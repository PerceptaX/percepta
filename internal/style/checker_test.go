//go:build linux

package style

import (
	"testing"
)

func TestStyleChecker_NamingViolations(t *testing.T) {
	checker := NewStyleChecker()

	// Code with naming violations
	source := []byte(`
// Bad function name (should be Module_Function)
void foo() {
	return;
}

// Good function name
void LED_Init() {
	return;
}

// Bad variable name (should be snake_case)
int myVariable = 5;

// Good variable name
int status_flag = 0;

// Bad constant name (should be UPPER_SNAKE)
const int maxSize = 100;

// Good constant name
const int MAX_BUFFER_SIZE = 256;
`)

	violations, err := checker.CheckSource(source, "test.c")
	if err != nil {
		t.Fatalf("CheckSource failed: %v", err)
	}

	// We expect violations for: foo, myVariable, maxSize
	expectedViolations := 3
	if len(violations) < expectedViolations {
		t.Errorf("Expected at least %d violations, got %d", expectedViolations, len(violations))
		for i, v := range violations {
			t.Logf("Violation %d: %s at line %d: %s", i+1, v.Rule.Name, v.Line, v.Message)
		}
	}

	// Check that violations have correct information
	for _, v := range violations {
		if v.File != "test.c" {
			t.Errorf("Expected file 'test.c', got '%s'", v.File)
		}
		if v.Line == 0 {
			t.Error("Line number should not be 0")
		}
		if v.Message == "" {
			t.Error("Message should not be empty")
		}
		if v.Suggestion == "" {
			t.Error("Suggestion should not be empty")
		}
	}
}

func TestStyleChecker_TypeViolations(t *testing.T) {
	checker := NewStyleChecker()

	// Code with type violations
	source := []byte(`
// Bad: using unsigned char instead of uint8_t
unsigned char status;

// Bad: using unsigned short instead of uint16_t
unsigned short counter;

// Bad: using unsigned int instead of uint32_t
unsigned int value;

// Good: using stdint.h types
uint8_t flags;
uint32_t timestamp;
`)

	violations, err := checker.CheckSource(source, "test.c")
	if err != nil {
		t.Fatalf("CheckSource failed: %v", err)
	}

	// We expect violations for: unsigned char, unsigned short, unsigned int
	typeViolationCount := 0
	for _, v := range violations {
		if v.Rule.ID == RuleStdintTypes {
			typeViolationCount++
		}
	}

	expectedTypeViolations := 3
	if typeViolationCount != expectedTypeViolations {
		t.Errorf("Expected %d type violations, got %d", expectedTypeViolations, typeViolationCount)
	}
}

func TestStyleChecker_PointerConstness(t *testing.T) {
	checker := NewStyleChecker()

	// Code with pointer const violations
	source := []byte(`
// Bad: pointer without const
void process(uint8_t* data) {
	return;
}

// Good: const pointer
void read(const uint8_t* data) {
	return;
}
`)

	violations, err := checker.CheckSource(source, "test.c")
	if err != nil {
		t.Fatalf("CheckSource failed: %v", err)
	}

	// We expect at least one const pointer warning
	constViolationCount := 0
	for _, v := range violations {
		if v.Rule.ID == RuleConstPointers {
			constViolationCount++
		}
	}

	if constViolationCount == 0 {
		t.Error("Expected at least one const pointer warning")
	}
}

func TestStyleChecker_CompleteExample(t *testing.T) {
	checker := NewStyleChecker()

	// Complete code with multiple violation types
	source := []byte(`
#include <stdint.h>

// VIOLATIONS:
// 1. Function name not Module_Function
void initLED() {
	// 2. Variable not snake_case
	int ledPin = 13;

	// 3. Using unsigned char instead of uint8_t
	unsigned char brightness = 255;
}

// GOOD CODE:
void LED_Init() {
	uint8_t led_pin = 13;
	const uint8_t max_brightness = 255;
}
`)

	violations, err := checker.CheckSource(source, "test.c")
	if err != nil {
		t.Fatalf("CheckSource failed: %v", err)
	}

	if len(violations) == 0 {
		t.Error("Expected violations but got none")
	}

	// Log all violations for inspection
	t.Logf("Found %d violations:", len(violations))
	for i, v := range violations {
		t.Logf("  %d. [%s] Line %d: %s", i+1, v.Rule.Category, v.Line, v.Message)
		if v.Suggestion != "" {
			t.Logf("     Suggestion: %s", v.Suggestion)
		}
	}
}

func TestStyleChecker_CleanCode(t *testing.T) {
	checker := NewStyleChecker()

	// Code that follows BARR-C standards
	source := []byte(`
#include <stdint.h>

void LED_Init(void) {
	const uint8_t led_pin = 13;
	uint8_t brightness = 0;
}

void LED_SetBrightness(uint8_t value) {
	// Implementation
}

int main(void) {
	LED_Init();
	LED_SetBrightness(128);
	return 0;
}
`)

	violations, err := checker.CheckSource(source, "test.c")
	if err != nil {
		t.Fatalf("CheckSource failed: %v", err)
	}

	// Clean code might still have some warnings (like const pointers)
	// But should not have any naming or type errors
	errorCount := 0
	for _, v := range violations {
		if v.Rule.Severity == "error" {
			errorCount++
			t.Logf("Unexpected error: %s at line %d: %s", v.Rule.Name, v.Line, v.Message)
		}
	}

	if errorCount > 0 {
		t.Errorf("Expected no errors in clean code, got %d", errorCount)
	}
}
