//go:build linux

package style

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestTypesFixer_BasicFix(t *testing.T) {
	source := []byte(`void test() {
    unsigned char x = 255;
}`)

	// Create violation with actual format from checker
	v := Violation{
		Rule:       StdintTypesRule,
		Line:       2,
		Column:     5,
		Message:    "Type 'unsigned char' should be replaced with stdint.h type (use uint8_t instead)",
		Suggestion: "Replace 'unsigned char x' with 'uint8_t x'",
	}

	fixer := &TypesFixer{}
	fixed, applied := fixer.Fix(v, source)

	if !applied {
		t.Fatal("Expected fix to be applied")
	}

	expected := `void test() {
    uint8_t x = 255;
}`

	if string(fixed) != expected {
		t.Errorf("Expected:\n%s\n\nGot:\n%s", expected, string(fixed))
	}
}

func TestNamingFixer_FunctionRename(t *testing.T) {
	source := []byte(`void initLED() {
    return;
}

void initLED() {
    initLED();
}`)

	// Create violation with actual format from checker
	v := Violation{
		Rule:       FunctionNamingRule,
		Line:       1,
		Column:     6,
		Message:    "Function 'initLED' should use Module_Function format (got: initLED)",
		Suggestion: "Module_InitLED",
	}

	fixer := &NamingFixer{}
	fixed, applied := fixer.Fix(v, source)

	if !applied {
		t.Fatal("Expected fix to be applied")
	}

	// Should replace ALL occurrences of the function name
	if !bytes.Contains(fixed, []byte("Module_InitLED")) {
		t.Error("Expected fixed source to contain 'Module_InitLED'")
	}

	// Count occurrences - should be 3 (2 definitions + 1 call)
	count := bytes.Count(fixed, []byte("Module_InitLED"))
	if count != 3 {
		t.Errorf("Expected 3 occurrences of Module_InitLED, got %d", count)
	}

	// Should not contain old name
	if bytes.Contains(fixed, []byte("initLED")) {
		t.Error("Fixed source should not contain old name 'initLED'")
	}
}

func TestStyleFixer_ApplyFixes(t *testing.T) {
	source := []byte(`void initLED() {
    unsigned char brightness = 255;
}`)

	// Create violations with actual format from checkers
	violations := []Violation{
		{
			Rule:       FunctionNamingRule,
			File:       "test.c",
			Line:       1,
			Column:     6,
			Message:    "Function 'initLED' should use Module_Function format (got: initLED)",
			Suggestion: "LED_Init",
		},
		{
			Rule:       StdintTypesRule,
			File:       "test.c",
			Line:       2,
			Column:     5,
			Message:    "Type 'unsigned char' should be replaced with stdint.h type (use uint8_t instead)",
			Suggestion: "Replace 'unsigned char brightness' with 'uint8_t brightness'",
		},
	}

	fixer := NewStyleFixer()
	fixed, fixedList := fixer.ApplyFixes(violations, source)

	// Should have applied 2 fixes
	if len(fixedList) != 2 {
		t.Errorf("Expected 2 fixes, got %d", len(fixedList))
	}

	// Check fixed source contains both changes
	if !bytes.Contains(fixed, []byte("LED_Init")) {
		t.Error("Expected fixed source to contain 'LED_Init'")
	}

	if !bytes.Contains(fixed, []byte("uint8_t")) {
		t.Error("Expected fixed source to contain 'uint8_t'")
	}
}

func TestStyleFixer_EnsureStdintHeader(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		fixedRules  []string
		shouldAdd   bool
		description string
	}{
		{
			name:        "no header, types fixed",
			source:      "void test() {\n    uint8_t x;\n}",
			fixedRules:  []string{"test.c:1:1 - Fixed: Stdint Type Usage"},
			shouldAdd:   true,
			description: "Should add header when types are fixed and header is missing",
		},
		{
			name:        "header exists, types fixed",
			source:      "#include <stdint.h>\nvoid test() {\n    uint8_t x;\n}",
			fixedRules:  []string{"test.c:1:1 - Fixed: Stdint Type Usage"},
			shouldAdd:   false,
			description: "Should not add duplicate header",
		},
		{
			name:        "no types fixed",
			source:      "void test() {\n    int x;\n}",
			fixedRules:  []string{"test.c:1:1 - Fixed: Function Naming Convention"},
			shouldAdd:   false,
			description: "Should not add header when no type fixes applied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixer := NewStyleFixer()
			result := fixer.EnsureStdintHeader([]byte(tt.source), tt.fixedRules)

			hasHeader := bytes.Contains(result, []byte("#include <stdint.h>"))

			if tt.shouldAdd && !hasHeader {
				t.Errorf("%s: expected header to be added", tt.description)
			}

			if !tt.shouldAdd && bytes.Contains([]byte(tt.source), []byte("#include <stdint.h>")) {
				// Original had header, result should still have it
				if !hasHeader {
					t.Errorf("%s: header was removed", tt.description)
				}
			}
		})
	}
}

func TestFixer_IntegrationWithChecker(t *testing.T) {
	// Create a test file with violations
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.c")

	source := []byte(`void initDevice() {
    unsigned char status = 1;
    unsigned short counter = 0;
}`)

	if err := os.WriteFile(testFile, source, 0644); err != nil {
		t.Fatal(err)
	}

	// Check for violations
	checker := NewStyleChecker()
	violations, err := checker.CheckFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(violations) == 0 {
		t.Fatal("Expected violations, got none")
	}

	// Debug: print violations
	t.Logf("Found %d violations:", len(violations))
	for i, v := range violations {
		t.Logf("  [%d] %s: %s (suggestion: %s)", i, v.Rule.Name, v.Message, v.Suggestion)
	}

	// Apply fixes
	fixer := NewStyleFixer()
	fixed, fixedList := fixer.ApplyFixes(violations, source)

	// Debug: print fixes
	t.Logf("Applied %d fixes:", len(fixedList))
	for i, f := range fixedList {
		t.Logf("  [%d] %s", i, f)
	}

	// Should have fixed some violations
	if len(fixedList) == 0 {
		t.Error("Expected some fixes to be applied")
	}

	// Add stdint header if needed
	fixed = fixer.EnsureStdintHeader(fixed, fixedList)

	// Write fixed source back
	if err := os.WriteFile(testFile, fixed, 0644); err != nil {
		t.Fatal(err)
	}

	// Re-check - should have fewer violations
	newViolations, err := checker.CheckFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	// Should have fewer violations (at least the type violations should be fixed)
	if len(newViolations) >= len(violations) {
		t.Errorf("Expected fewer violations after fix. Before: %d, After: %d", len(violations), len(newViolations))
	}

	// Check that stdint header was added
	fixedSource, _ := os.ReadFile(testFile)
	if !bytes.Contains(fixedSource, []byte("#include <stdint.h>")) {
		t.Error("Expected stdint.h header to be added")
	}
}
