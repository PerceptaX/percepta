//go:build linux

package codegen

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/perceptumx/percepta/internal/style"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrintGenerationReport_StyleCompliant(t *testing.T) {
	result := &GenerationResult{
		Code:           "void LED_Init(void) {\n    // code\n}\n",
		StyleCompliant: true,
		Violations:     []style.Violation{},
		AutoFixed:      false,
		PatternStored:  true,
		IterationsUsed: 1,
	}

	output := captureOutput(func() {
		PrintGenerationReport(result)
	})

	// Verify output contains expected elements
	if !strings.Contains(output, "GENERATION REPORT") {
		t.Error("Expected report header")
	}

	if !strings.Contains(output, "✓ Style: BARR-C compliant") {
		t.Error("Expected style compliance message")
	}

	if !strings.Contains(output, "✓ Pattern: Stored in knowledge graph") {
		t.Error("Expected pattern storage message")
	}

	if !strings.Contains(output, "4 lines generated in 1 iteration(s)") {
		t.Error("Expected code stats")
	}
}

func TestPrintGenerationReport_WithViolations(t *testing.T) {
	result := &GenerationResult{
		Code:           "int main() { return 0; }",
		StyleCompliant: false,
		Violations: []style.Violation{
			{
				File:    "test.c",
				Line:    1,
				Column:  5,
				Message: "Function name should use PascalCase",
				Rule: style.Rule{
					Name:     "function-naming",
					Severity: "error",
					Category: "naming",
				},
			},
			{
				File:    "test.c",
				Line:    1,
				Column:  1,
				Message: "Should use uint8_t instead of int",
				Rule: style.Rule{
					Name:     "fixed-width-types",
					Severity: "error",
					Category: "types",
				},
			},
		},
		AutoFixed:      true,
		PatternStored:  false,
		IterationsUsed: 1,
	}

	output := captureOutput(func() {
		PrintGenerationReport(result)
	})

	// Verify violations are shown
	if !strings.Contains(output, "✗ Style: 2 violation(s) remaining") {
		t.Error("Expected violation count")
	}

	if !strings.Contains(output, "Function name should use PascalCase") {
		t.Error("Expected first violation message")
	}

	if !strings.Contains(output, "Should use uint8_t instead of int") {
		t.Error("Expected second violation message")
	}

	if !strings.Contains(output, "[function-naming]") {
		t.Error("Expected rule name for first violation")
	}

	if !strings.Contains(output, "[fixed-width-types]") {
		t.Error("Expected rule name for second violation")
	}

	// Verify auto-fix message
	if !strings.Contains(output, "✓ Auto-fix: Applied deterministic corrections") {
		t.Error("Expected auto-fix message")
	}
}

func TestPrintGenerationReport_NoPatternStorage(t *testing.T) {
	result := &GenerationResult{
		Code:           "void Test(void) {}",
		StyleCompliant: true,
		Violations:     []style.Violation{},
		AutoFixed:      false,
		PatternStored:  false,
		IterationsUsed: 1,
	}

	output := captureOutput(func() {
		PrintGenerationReport(result)
	})

	// Verify pattern not stored message
	if !strings.Contains(output, "✗ Pattern: Not stored (storage unavailable)") {
		t.Error("Expected pattern not stored message")
	}
}

func TestPrintGenerationReport_MultipleIterations(t *testing.T) {
	result := &GenerationResult{
		Code:           "void LED_Init(void) {\n}\n",
		StyleCompliant: true,
		Violations:     []style.Violation{},
		AutoFixed:      false,
		PatternStored:  true,
		IterationsUsed: 3,
	}

	output := captureOutput(func() {
		PrintGenerationReport(result)
	})

	// Verify iteration count
	if !strings.Contains(output, "3 lines generated in 3 iteration(s)") {
		t.Error("Expected correct iteration count")
	}
}

func TestPrintGenerationReport_NoAutoFix(t *testing.T) {
	result := &GenerationResult{
		Code:           "void LED_Init(void) {}",
		StyleCompliant: true,
		Violations:     []style.Violation{},
		AutoFixed:      false,
		PatternStored:  true,
		IterationsUsed: 1,
	}

	output := captureOutput(func() {
		PrintGenerationReport(result)
	})

	// Verify no auto-fix message when AutoFixed is false
	if strings.Contains(output, "Auto-fix:") {
		t.Error("Should not show auto-fix message when AutoFixed is false")
	}
}

func TestPrintGenerationReport_EmptyCode(t *testing.T) {
	result := &GenerationResult{
		Code:           "",
		StyleCompliant: false,
		Violations:     []style.Violation{},
		AutoFixed:      false,
		PatternStored:  false,
		IterationsUsed: 1,
	}

	output := captureOutput(func() {
		PrintGenerationReport(result)
	})

	// Verify it handles empty code gracefully
	if !strings.Contains(output, "GENERATION REPORT") {
		t.Error("Expected report header even for empty code")
	}

	// Empty string split by newline gives 1 line
	if !strings.Contains(output, "1 lines generated") {
		t.Error("Expected line count for empty code")
	}
}

func TestPrintGenerationReport_LargeCodeWithManyViolations(t *testing.T) {
	// Generate a large code sample
	code := strings.Repeat("void function() {}\n", 100)

	// Generate many violations
	violations := make([]style.Violation, 50)
	for i := 0; i < 50; i++ {
		violations[i] = style.Violation{
			File:    "large.c",
			Line:    i + 1,
			Column:  1,
			Message: "Violation " + strings.Repeat("X", i),
			Rule: style.Rule{
				Name:     "rule-" + strings.Repeat("A", i),
				Severity: "error",
				Category: "test",
			},
		}
	}

	result := &GenerationResult{
		Code:           code,
		StyleCompliant: false,
		Violations:     violations,
		AutoFixed:      true,
		PatternStored:  false,
		IterationsUsed: 5,
	}

	output := captureOutput(func() {
		PrintGenerationReport(result)
	})

	// Verify all violations are printed
	if !strings.Contains(output, "50 violation(s) remaining") {
		t.Error("Expected all 50 violations to be counted")
	}

	// Verify line count
	if !strings.Contains(output, "101 lines generated in 5 iteration(s)") {
		t.Error("Expected correct line count")
	}
}

func TestPrintGenerationReport_SingleLineCode(t *testing.T) {
	result := &GenerationResult{
		Code:           "void LED_Init(void) {}",
		StyleCompliant: true,
		Violations:     []style.Violation{},
		AutoFixed:      false,
		PatternStored:  true,
		IterationsUsed: 1,
	}

	output := captureOutput(func() {
		PrintGenerationReport(result)
	})

	// Verify line count for single line
	if !strings.Contains(output, "1 lines generated") {
		t.Error("Expected 1 line count for single line code")
	}
}

func TestPrintGenerationReport_WithComplexViolation(t *testing.T) {
	result := &GenerationResult{
		Code:           "int x = 0;",
		StyleCompliant: false,
		Violations: []style.Violation{
			{
				File:    "test.c",
				Line:    1,
				Column:  1,
				Message: "Variable 'x' should have a descriptive name (min 3 chars)",
				Rule: style.Rule{
					Name:     "variable-naming",
					Severity: "error",
					Category: "naming",
				},
			},
		},
		AutoFixed:      false,
		PatternStored:  false,
		IterationsUsed: 1,
	}

	output := captureOutput(func() {
		PrintGenerationReport(result)
	})

	// Verify complex violation message is preserved
	if !strings.Contains(output, "Variable 'x' should have a descriptive name") {
		t.Error("Expected full violation message")
	}

	if !strings.Contains(output, "[variable-naming]") {
		t.Error("Expected rule name in brackets")
	}
}

func TestPrintGenerationReport_HeaderAndFooter(t *testing.T) {
	result := &GenerationResult{
		Code:           "void Test(void) {}",
		StyleCompliant: true,
		Violations:     []style.Violation{},
		AutoFixed:      false,
		PatternStored:  true,
		IterationsUsed: 1,
	}

	output := captureOutput(func() {
		PrintGenerationReport(result)
	})

	// Verify report has proper header and footer separators
	separatorCount := strings.Count(output, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if separatorCount < 2 {
		t.Errorf("Expected at least 2 separator lines, got %d", separatorCount)
	}
}
