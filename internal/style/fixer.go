package style

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

// Fixer is the interface for auto-fixing violations
type Fixer interface {
	// Fix returns fixed source if applicable, or original if no fix
	Fix(violation Violation, source []byte) ([]byte, bool)
}

// NamingFixer fixes function and variable naming violations
type NamingFixer struct{}

// Fix applies auto-fixes for naming violations
func (n *NamingFixer) Fix(v Violation, source []byte) ([]byte, bool) {
	// Only fix function names, not variables/constants (too risky without full context)
	if v.Rule.ID != RuleFunctionNaming {
		return source, false
	}

	// Extract the function name from the violation message
	// Message format: "Function '<name>' should use Module_Function format (got: <name>)"
	re := regexp.MustCompile(`Function '([^']+)' should use Module_Function format`)
	matches := re.FindStringSubmatch(v.Message)
	if len(matches) < 2 {
		return source, false
	}
	oldName := matches[1]

	// The suggestion is just the new name directly (e.g., "Module_InitDevice")
	if v.Suggestion == "" {
		return source, false
	}
	newName := v.Suggestion

	// Replace the function name in the source
	// Use word boundary to avoid replacing partial matches
	pattern := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(oldName))
	re = regexp.MustCompile(pattern)
	fixed := re.ReplaceAll(source, []byte(newName))

	return fixed, !bytes.Equal(source, fixed)
}

// TypesFixer fixes type-related violations (stdint.h types)
type TypesFixer struct{}

// Fix applies auto-fixes for type violations
func (t *TypesFixer) Fix(v Violation, source []byte) ([]byte, bool) {
	if v.Rule.ID != RuleStdintTypes {
		return source, false
	}

	// Suggestion format: "Replace 'unsigned char status' with 'uint8_t status'"
	// Extract old and new parts
	re := regexp.MustCompile(`Replace '([^']+)' with '([^']+)'`)
	matches := re.FindStringSubmatch(v.Suggestion)
	if len(matches) < 3 {
		return source, false
	}

	oldDecl := matches[1]
	newDecl := matches[2]

	// Split source into lines
	lines := bytes.Split(source, []byte("\n"))

	// Fix the specific line (v.Line is 1-indexed)
	if v.Line > 0 && v.Line <= len(lines) {
		lineIdx := v.Line - 1
		line := lines[lineIdx]

		// Replace the old declaration with new one on this specific line
		fixed := bytes.Replace(line, []byte(oldDecl), []byte(newDecl), 1)
		lines[lineIdx] = fixed
	}

	// Rejoin the lines
	result := bytes.Join(lines, []byte("\n"))

	return result, !bytes.Equal(source, result)
}

// StyleFixer orchestrates multiple fixers
type StyleFixer struct {
	fixers map[string]Fixer
}

// NewStyleFixer creates a new style fixer with all fixers enabled
func NewStyleFixer() *StyleFixer {
	return &StyleFixer{
		fixers: map[string]Fixer{
			"naming": &NamingFixer{},
			"types":  &TypesFixer{},
		},
	}
}

// ApplyFixes applies all applicable fixes to the source code
func (s *StyleFixer) ApplyFixes(violations []Violation, source []byte) ([]byte, []string) {
	current := source
	fixed := make([]string, 0)

	// Apply fixes in order: types first (safer), then naming
	// This ensures we don't mess up naming fixes with type changes
	categories := []string{"types", "naming"}

	for _, category := range categories {
		fixer, ok := s.fixers[category]
		if !ok {
			continue
		}

		for _, v := range violations {
			if v.Rule.Category == category {
				if newSource, applied := fixer.Fix(v, current); applied {
					current = newSource
					fixed = append(fixed, fmt.Sprintf("%s:%d:%d - Fixed: %s", v.File, v.Line, v.Column, v.Rule.Name))
				}
			}
		}
	}

	return current, fixed
}

// EnsureStdintHeader adds #include <stdint.h> if types were fixed and header is missing
func (s *StyleFixer) EnsureStdintHeader(source []byte, fixedRules []string) []byte {
	// Check if any type fixes were applied
	hasTypeFix := false
	for _, fix := range fixedRules {
		if strings.Contains(fix, "Stdint Type Usage") {
			hasTypeFix = true
			break
		}
	}

	if !hasTypeFix {
		return source
	}

	// Check if <stdint.h> is already included
	if bytes.Contains(source, []byte("#include <stdint.h>")) {
		return source
	}

	// Find the first include statement or insert at the top
	lines := bytes.Split(source, []byte("\n"))
	insertIdx := -1
	lastIncludeIdx := -1

	// Look for existing include statements
	for i, line := range lines {
		if bytes.Contains(line, []byte("#include")) {
			lastIncludeIdx = i
		}
	}

	// If there are existing includes, insert after the last one
	if lastIncludeIdx >= 0 {
		insertIdx = lastIncludeIdx + 1
	} else {
		// No existing includes - find first non-comment, non-blank line
		for i, line := range lines {
			trimmed := bytes.TrimSpace(line)
			if len(trimmed) > 0 && !bytes.HasPrefix(trimmed, []byte("//")) && !bytes.HasPrefix(trimmed, []byte("/*")) {
				insertIdx = i
				break
			}
		}
		// If we didn't find any code, insert at beginning
		if insertIdx == -1 {
			insertIdx = 0
		}
	}

	// Insert the header
	newLine := []byte("#include <stdint.h>")

	// Insert at the determined position
	result := make([][]byte, 0, len(lines)+1)
	result = append(result, lines[:insertIdx]...)
	result = append(result, newLine)
	result = append(result, lines[insertIdx:]...)

	return bytes.Join(result, []byte("\n"))
}
