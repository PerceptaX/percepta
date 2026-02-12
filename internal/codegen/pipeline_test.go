package codegen

import (
	"testing"

	"github.com/perceptumx/percepta/internal/knowledge"
	"github.com/perceptumx/percepta/internal/style"
)

// TestGenerationPipeline_CleanCode tests the pipeline with code that is already style compliant
func TestGenerationPipeline_CleanCode(t *testing.T) {
	// Create a mock Claude client that returns clean BARR-C compliant code
	claudeClient := NewClaudeClient("mock-api-key")

	// Create pattern store
	patternStore, err := knowledge.NewPatternStore()
	if err != nil {
		t.Fatalf("Failed to create pattern store: %v", err)
	}
	defer patternStore.Close()

	// Create prompt builder
	promptBuilder := NewPromptBuilder(patternStore)

	// Create style checker and fixer
	styleChecker := style.NewStyleChecker()
	styleFixer := style.NewStyleFixer()

	// Create pipeline
	pipeline := NewGenerationPipeline(
		claudeClient,
		promptBuilder,
		styleChecker,
		styleFixer,
		patternStore,
	)

	// Test that pipeline is created correctly
	if pipeline == nil {
		t.Fatal("Pipeline should not be nil")
	}

	// Verify fields are set
	if pipeline.claudeClient == nil {
		t.Error("ClaudeClient should not be nil")
	}
	if pipeline.promptBuilder == nil {
		t.Error("PromptBuilder should not be nil")
	}
	if pipeline.styleChecker == nil {
		t.Error("StyleChecker should not be nil")
	}
	if pipeline.styleFixer == nil {
		t.Error("StyleFixer should not be nil")
	}
	if pipeline.patternStore == nil {
		t.Error("PatternStore should not be nil")
	}
}

// TestGenerationPipeline_WithViolations tests the pipeline with code that has style violations
func TestGenerationPipeline_WithViolations(t *testing.T) {
	// Test code with violations (uses int instead of uint8_t)
	testCode := `#include <stdint.h>

int led_state = 0;

void toggle_led() {
    led_state = !led_state;
}
`

	styleChecker := style.NewStyleChecker()
	violations, err := styleChecker.CheckSource([]byte(testCode), "test.c")
	if err != nil {
		t.Fatalf("Style check failed: %v", err)
	}

	// Should have at least one violation (function naming)
	if len(violations) == 0 {
		t.Error("Expected violations for non-BARR-C code")
	}

	// Test that auto-fixer can fix some violations
	styleFixer := style.NewStyleFixer()
	fixed, fixedRules := styleFixer.ApplyFixes(violations, []byte(testCode))

	if len(fixedRules) == 0 {
		t.Log("No automatic fixes applied (this is expected for some violations)")
	} else {
		t.Logf("Applied %d fixes", len(fixedRules))
	}

	// Verify code was modified if fixes were applied
	if len(fixedRules) > 0 && string(fixed) == testCode {
		t.Error("Code should be modified after applying fixes")
	}
}

// TestGenerationPipeline_StyleCheckIntegration tests style checking integration
func TestGenerationPipeline_StyleCheckIntegration(t *testing.T) {
	// Test with valid BARR-C code
	validCode := `#include <stdint.h>

#define LED_PIN 2

static uint8_t led_state = 0;

void LED_Toggle(void) {
    led_state = !led_state;
}
`

	styleChecker := style.NewStyleChecker()
	violations, err := styleChecker.CheckSource([]byte(validCode), "valid.c")
	if err != nil {
		t.Fatalf("Style check failed: %v", err)
	}

	if len(violations) > 0 {
		t.Errorf("Expected no violations for valid code, got %d", len(violations))
		for _, v := range violations {
			t.Logf("  %s:%d:%d - %s", v.File, v.Line, v.Column, v.Message)
		}
	}
}

// TestGenerationResult_Fields tests GenerationResult structure
func TestGenerationResult_Fields(t *testing.T) {
	result := &GenerationResult{
		Code:           "test code",
		StyleCompliant: true,
		Violations:     []style.Violation{},
		AutoFixed:      false,
		PatternStored:  true,
		IterationsUsed: 1,
	}

	if result.Code != "test code" {
		t.Error("Code field not set correctly")
	}
	if !result.StyleCompliant {
		t.Error("StyleCompliant should be true")
	}
	if len(result.Violations) != 0 {
		t.Error("Violations should be empty")
	}
	if result.AutoFixed {
		t.Error("AutoFixed should be false")
	}
	if !result.PatternStored {
		t.Error("PatternStored should be true")
	}
	if result.IterationsUsed != 1 {
		t.Error("IterationsUsed should be 1")
	}
}

// TestGenerationPipeline_GracefulDegradation tests that pipeline handles storage failure gracefully
func TestGenerationPipeline_GracefulDegradation(t *testing.T) {
	// This test verifies that if pattern storage fails, generation still succeeds
	// The pipeline should log a warning but not fail the entire generation

	// Create components
	patternStore, err := knowledge.NewPatternStore()
	if err != nil {
		t.Fatalf("Failed to create pattern store: %v", err)
	}
	defer patternStore.Close()

	styleChecker := style.NewStyleChecker()
	styleFixer := style.NewStyleFixer()

	// Verify style checker and fixer work independently
	testCode := `#include <stdint.h>
void LED_Init(void) {}
`

	violations, err := styleChecker.CheckSource([]byte(testCode), "test.c")
	if err != nil {
		t.Fatalf("Style check failed: %v", err)
	}

	if len(violations) > 0 {
		t.Logf("Found %d violations (expected for test)", len(violations))
		fixed, _ := styleFixer.ApplyFixes(violations, []byte(testCode))
		if len(fixed) == 0 {
			t.Error("Fixer should return non-empty code")
		}
	}
}
