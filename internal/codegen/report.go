package codegen

import (
	"fmt"
	"strings"
)

// PrintGenerationReport prints a detailed report of the generation results
func PrintGenerationReport(result *GenerationResult) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("GENERATION REPORT")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Style compliance
	if result.StyleCompliant {
		fmt.Println("✓ Style: BARR-C compliant")
	} else {
		fmt.Printf("✗ Style: %d violation(s) remaining\n", len(result.Violations))
		for _, v := range result.Violations {
			fmt.Printf("  Line %d: %s [%s]\n", v.Line, v.Message, v.Rule.Name)
		}
	}

	// Auto-fix status
	if result.AutoFixed {
		fmt.Println("✓ Auto-fix: Applied deterministic corrections")
	}

	// Pattern storage
	if result.PatternStored {
		fmt.Println("✓ Pattern: Stored in knowledge graph")
		fmt.Println("  (Will improve future generations)")
	} else if result.StyleCompliant {
		fmt.Println("✗ Pattern: Not stored (storage unavailable)")
	}

	// Code stats
	lines := len(strings.Split(result.Code, "\n"))
	fmt.Printf("\nCode: %d lines generated in %d iteration(s)\n", lines, result.IterationsUsed)

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
