package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/perceptumx/percepta/internal/knowledge"
	"github.com/spf13/cobra"
)

var knowledgeCmd = &cobra.Command{
	Use:   "knowledge",
	Short: "Manage validated pattern knowledge graph",
	Long: `Store and query validated firmware patterns with behavioral and style metadata.

The knowledge graph stores only hardware-validated, BARR-C compliant patterns
that have been tested on real devices. Semantic search finds similar patterns
by code meaning, not just exact text matches.

Commands:
  store   - Store a validated pattern from a device observation
  search  - Search for similar patterns semantically
  list    - List all validated patterns`,
}

var storePatternCmd = &cobra.Command{
	Use:   "store <spec> <file.c> --device <device-id> --firmware <tag>",
	Short: "Store validated pattern in knowledge graph",
	Long: `Validates code against BARR-C, links to observation, stores in graph.

The pattern will only be stored if:
1. Code is BARR-C compliant (no style violations)
2. An observation exists for this device+firmware combination
3. All relationships can be created in the graph

Example:
  percepta knowledge store "Blink LED at 1Hz" led.c --device esp32-dev --firmware v1.0.0`,
	Args: cobra.ExactArgs(2),
	RunE: runStorePattern,
}

var searchPatternsCmd = &cobra.Command{
	Use:   "search <query> [--board <type>] [--limit <n>]",
	Short: "Search for similar validated patterns",
	Long: `Semantic search for patterns solving similar problems.

Uses vector embeddings to find patterns by code similarity, not exact matches.
Results are ranked by similarity and filtered by board type if specified.

Example:
  percepta knowledge search "blink LED" --board esp32 --limit 5`,
	Args: cobra.ExactArgs(1),
	RunE: runSearchPatterns,
}

var listPatternsCmd = &cobra.Command{
	Use:   "list [--board <type>]",
	Short: "List all validated patterns",
	Long: `List all validated patterns in the knowledge graph.

Optionally filter by board type to see only patterns for specific hardware.

Example:
  percepta knowledge list --board esp32`,
	Args: cobra.NoArgs,
	RunE: runListPatterns,
}

var (
	deviceFlag   string
	firmwareFlag string
	boardFlag    string
	limitFlag    int
)

func init() {
	// Add subcommands
	knowledgeCmd.AddCommand(storePatternCmd)
	knowledgeCmd.AddCommand(searchPatternsCmd)
	knowledgeCmd.AddCommand(listPatternsCmd)

	// Store flags
	storePatternCmd.Flags().StringVarP(&deviceFlag, "device", "d", "", "Device ID (required)")
	storePatternCmd.Flags().StringVarP(&firmwareFlag, "firmware", "f", "", "Firmware tag (required)")
	storePatternCmd.MarkFlagRequired("device")
	storePatternCmd.MarkFlagRequired("firmware")

	// Search flags
	searchPatternsCmd.Flags().StringVarP(&boardFlag, "board", "b", "", "Filter by board type")
	searchPatternsCmd.Flags().IntVarP(&limitFlag, "limit", "l", 5, "Number of results")

	// List flags
	listPatternsCmd.Flags().StringVarP(&boardFlag, "board", "b", "", "Filter by board type")
}

func runStorePattern(cmd *cobra.Command, args []string) error {
	spec := args[0]
	codeFile := args[1]

	// Read code file
	code, err := os.ReadFile(codeFile)
	if err != nil {
		return fmt.Errorf("failed to read code file: %w", err)
	}

	// Initialize pattern store
	store, err := knowledge.NewPatternStore()
	if err != nil {
		return fmt.Errorf("failed to initialize pattern store: %w", err)
	}
	defer store.Close()

	// Initialize vector store for semantic search
	if err := store.InitializeVectorStore(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: semantic search unavailable (vector store init failed: %v)\n", err)
		// Continue anyway - pattern can still be stored in graph
	} else {
		defer store.CloseVectorStore()
	}

	// Store validated pattern
	patternID, err := store.StoreValidatedPattern(spec, string(code), deviceFlag, firmwareFlag)
	if err != nil {
		return fmt.Errorf("failed to store pattern: %w", err)
	}

	// Success output
	fmt.Printf("✓ Pattern stored successfully\n")
	fmt.Printf("  ID:       %s\n", patternID[:16]+"...")
	fmt.Printf("  Spec:     %s\n", spec)
	fmt.Printf("  Device:   %s\n", deviceFlag)
	fmt.Printf("  Firmware: %s\n", firmwareFlag)
	fmt.Printf("  File:     %s\n", codeFile)

	// Show stats
	stats := store.Stats()
	fmt.Printf("\nKnowledge graph stats:\n")
	fmt.Printf("  Patterns:      %d\n", stats["patterns"])
	fmt.Printf("  Observations:  %d\n", stats["observations"])

	return nil
}

func runSearchPatterns(cmd *cobra.Command, args []string) error {
	query := args[0]

	// Initialize pattern store
	store, err := knowledge.NewPatternStore()
	if err != nil {
		return fmt.Errorf("failed to initialize pattern store: %w", err)
	}
	defer store.Close()

	// Initialize vector store (required for semantic search)
	if err := store.InitializeVectorStore(); err != nil {
		return fmt.Errorf("semantic search unavailable: %w\n\nMake sure OPENAI_API_KEY environment variable is set", err)
	}
	defer store.CloseVectorStore()

	// Search for similar patterns
	results, err := store.SearchSimilarPatterns(query, boardFlag, limitFlag)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	// Display results
	if len(results) == 0 {
		fmt.Println("No similar patterns found.")
		if boardFlag != "" {
			fmt.Printf("(searched board type: %s)\n", boardFlag)
		}
		return nil
	}

	fmt.Printf("Found %d similar pattern(s):\n\n", len(results))

	for i, result := range results {
		fmt.Printf("%d. %s\n", i+1, result.Pattern.Spec)
		fmt.Printf("   Board:      %s\n", result.Pattern.BoardType)
		fmt.Printf("   Similarity: %.0f%%\n", result.Similarity*100)
		fmt.Printf("   Confidence: %.0f%%\n", result.Confidence*100)
		fmt.Printf("   Style:      BARR-C compliant\n")

		// Show code preview (first 3 lines)
		codeLines := strings.Split(result.Pattern.Code, "\n")
		previewLines := 3
		if len(codeLines) < previewLines {
			previewLines = len(codeLines)
		}

		fmt.Printf("   Code:\n")
		for j := 0; j < previewLines; j++ {
			fmt.Printf("     %s\n", codeLines[j])
		}
		if len(codeLines) > previewLines {
			fmt.Printf("     ... (%d more lines)\n", len(codeLines)-previewLines)
		}

		// Show observation signals
		if result.Observation != nil && len(result.Observation.Signals) > 0 {
			fmt.Printf("   Signals:    %d observed\n", len(result.Observation.Signals))
		}

		if i < len(results)-1 {
			fmt.Println()
		}
	}

	return nil
}

func runListPatterns(cmd *cobra.Command, args []string) error {
	// Initialize pattern store
	store, err := knowledge.NewPatternStore()
	if err != nil {
		return fmt.Errorf("failed to initialize pattern store: %w", err)
	}
	defer store.Close()

	// Query patterns
	var patterns []*knowledge.PatternNode
	if boardFlag != "" {
		patterns, err = store.QueryPatternsByBoard(boardFlag)
	} else {
		// Get all patterns by querying stats and then getting each board type
		stats := store.Stats()
		if stats["patterns"] == 0 {
			fmt.Println("No patterns stored in knowledge graph.")
			return nil
		}

		// For now, query common board types
		// In a real implementation, we'd have a method to list all unique board types
		boardTypes := []string{"esp32", "stm32", "arduino", "rp2040"}
		seen := make(map[string]bool)

		for _, bt := range boardTypes {
			boardPatterns, _ := store.QueryPatternsByBoard(bt)
			for _, p := range boardPatterns {
				if !seen[p.ID] {
					patterns = append(patterns, p)
					seen[p.ID] = true
				}
			}
		}
	}

	if err != nil {
		return fmt.Errorf("failed to query patterns: %w", err)
	}

	// Display results
	if len(patterns) == 0 {
		if boardFlag != "" {
			fmt.Printf("No patterns found for board type: %s\n", boardFlag)
		} else {
			fmt.Println("No patterns stored in knowledge graph.")
		}
		return nil
	}

	fmt.Printf("Validated patterns (%d total):\n\n", len(patterns))

	for i, pattern := range patterns {
		fmt.Printf("%d. %s\n", i+1, pattern.Spec)
		fmt.Printf("   Board:    %s\n", pattern.BoardType)
		fmt.Printf("   Style:    BARR-C compliant\n")
		fmt.Printf("   Created:  %s\n", pattern.CreatedAt.Format("2006-01-02 15:04:05"))

		// Show code preview (first 2 lines)
		codeLines := strings.Split(pattern.Code, "\n")
		previewLines := 2
		if len(codeLines) < previewLines {
			previewLines = len(codeLines)
		}

		fmt.Printf("   Code:\n")
		for j := 0; j < previewLines; j++ {
			fmt.Printf("     %s\n", codeLines[j])
		}
		if len(codeLines) > previewLines {
			fmt.Printf("     ... (%d more lines)\n", len(codeLines)-previewLines)
		}

		if i < len(patterns)-1 {
			fmt.Println()
		}
	}

	return nil
}
