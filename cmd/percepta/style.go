package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/perceptumx/percepta/internal/style"
	"github.com/spf13/cobra"
)

var styleCmd = &cobra.Command{
	Use:   "style-check <file-or-directory>",
	Short: "Check C code for BARR-C compliance",
	Long: `Validates generated firmware against BARR-C Embedded C Coding Standard.

The BARR-C standard defines professional coding practices for embedded systems:
- Module_Function naming for all functions
- snake_case for variables
- UPPER_SNAKE for global constants
- stdint.h types (uint8_t) instead of primitives (unsigned char)
- const correctness for pointers

Use --fix to auto-correct deterministic violations (naming, types).
Manual review required for magic numbers and const correctness.

Examples:
  # Check single file
  percepta style-check led_blink.c

  # Check entire directory
  percepta style-check ./src

  # Auto-fix violations
  percepta style-check led_blink.c --fix

Exit codes:
  0 - No violations found (BARR-C compliant)
  1 - Violations found`,
	Args: cobra.ExactArgs(1),
	RunE: runStyleCheck,
}

var fixFlag bool

func init() {
	styleCmd.Flags().BoolVar(&fixFlag, "fix", false, "Auto-fix violations where possible")
}

func runStyleCheck(cmd *cobra.Command, args []string) error {
	path := args[0]

	// Detect if file or directory
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path not found: %w", err)
	}

	var files []string
	if info.IsDir() {
		// Walk directory, find *.c and *.h files
		err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && (strings.HasSuffix(p, ".c") || strings.HasSuffix(p, ".h")) {
				files = append(files, p)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to walk directory: %w", err)
		}
	} else {
		files = []string{path}
	}

	if len(files) == 0 {
		fmt.Println("No C files found.")
		return nil
	}

	checker := style.NewStyleChecker()
	fixer := style.NewStyleFixer()

	totalViolations := 0
	totalFixed := 0
	filesWithViolations := 0

	for _, file := range files {
		violations, err := checker.CheckFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking %s: %v\n", file, err)
			continue
		}

		if fixFlag && len(violations) > 0 {
			// Read the source
			source, err := os.ReadFile(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", file, err)
				continue
			}

			// Apply fixes
			fixed, fixedList := fixer.ApplyFixes(violations, source)

			// Add stdint header if types were fixed
			fixed = fixer.EnsureStdintHeader(fixed, fixedList)

			// Write fixed source back
			if len(fixedList) > 0 {
				if err := os.WriteFile(file, fixed, 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", file, err)
					continue
				}
				totalFixed += len(fixedList)

				// Print fixed violations
				fmt.Printf("\nFixed in %s:\n", file)
				for _, f := range fixedList {
					fmt.Printf("  ✓ %s\n", f)
				}
			}

			// Re-check to see remaining violations
			violations, err = checker.CheckFile(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error re-checking %s: %v\n", file, err)
				continue
			}
		}

		if len(violations) > 0 {
			filesWithViolations++
			fmt.Printf("\n%s:\n", file)
			for _, v := range violations {
				// Format: file:line:col severity [rule] message
				fmt.Printf("  %d:%d %s [%s] %s\n",
					v.Line, v.Column, v.Rule.Severity, v.Rule.Name, v.Message)
				if v.Suggestion != "" {
					fmt.Printf("    → %s\n", v.Suggestion)
				}
			}
			totalViolations += len(violations)
		}
	}

	// Summary
	fmt.Println()
	if totalViolations == 0 && totalFixed == 0 {
		fmt.Println("✅ No style violations found. Code is BARR-C compliant.")
		return nil
	}

	if fixFlag && totalFixed > 0 {
		fmt.Printf("✅ Fixed %d violation(s) automatically.\n", totalFixed)
	}

	if totalViolations > 0 {
		fmt.Printf("⚠️  %d violation(s) remain in %d file(s).\n", totalViolations, filesWithViolations)
		if !fixFlag {
			fmt.Println("\nRun with --fix to auto-correct fixable violations.")
		} else {
			fmt.Println("\nRemaining violations require manual review.")
		}
		return fmt.Errorf("style violations detected")
	}

	return nil
}
