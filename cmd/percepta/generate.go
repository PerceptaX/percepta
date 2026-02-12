package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/perceptumx/percepta/internal/codegen"
	"github.com/perceptumx/percepta/internal/knowledge"
	"github.com/spf13/cobra"
)

var (
	boardType  string
	outputFile string
)

var generateCmd = &cobra.Command{
	Use:   "generate <spec> --board <type> [--output <file>]",
	Short: "Generate BARR-C compliant firmware from specification",
	Long: `Uses AI and validated patterns to generate firmware code.

The generate command uses Claude AI with knowledge from validated patterns
to create professional, BARR-C compliant embedded C code. Generated code
follows established working patterns and includes proper error handling,
non-blocking architecture, and static allocation.

Example:
  percepta generate "Blink LED at 1Hz" --board esp32 --output led_blink.c
  percepta generate "Read temperature sensor every 2 seconds" --board stm32
  percepta generate "Toggle LED on button press" --board arduino --output button_led.c

Requires:
  ANTHROPIC_API_KEY environment variable set with your API key from:
  https://console.anthropic.com/

Optional:
  OPENAI_API_KEY for semantic pattern search (graceful degradation without it)`,
	Args: cobra.MinimumNArgs(1),
	RunE: runGenerate,
}

func init() {
	generateCmd.Flags().StringVarP(&boardType, "board", "b", "", "Board type (required): esp32, stm32, arduino, etc.")
	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file (optional, prints to stdout if not set)")
	generateCmd.MarkFlagRequired("board")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	spec := strings.Join(args, " ")

	// Check for API key
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return fmt.Errorf(`ANTHROPIC_API_KEY not set

Get your API key from: https://console.anthropic.com/
Set it: export ANTHROPIC_API_KEY=your-key-here`)
	}

	fmt.Printf("Generating firmware...\n")
	fmt.Printf("Spec: %s\n", spec)
	fmt.Printf("Board: %s\n\n", boardType)

	// 1. Initialize pattern store
	patternStore, err := knowledge.NewPatternStore()
	if err != nil {
		return fmt.Errorf("failed to load patterns: %w", err)
	}
	defer patternStore.Close()

	// 2. Initialize vector store for semantic search (optional)
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey != "" {
		if err := patternStore.InitializeVectorStore(); err != nil {
			// Non-fatal: can still generate without semantic search
			fmt.Printf("Warning: semantic search unavailable: %v\n", err)
			fmt.Printf("Continuing with basic pattern matching...\n\n")
		} else {
			fmt.Printf("✓ Semantic search enabled\n")
		}
	} else {
		fmt.Printf("Note: OPENAI_API_KEY not set, semantic search disabled\n")
		fmt.Printf("Continuing with basic BARR-C requirements...\n\n")
	}

	// 3. Build prompt with patterns
	promptBuilder := codegen.NewPromptBuilder(patternStore)
	systemPrompt, err := promptBuilder.BuildSystemPrompt(spec, boardType)
	if err != nil {
		return fmt.Errorf("failed to build prompt: %w", err)
	}

	// 4. Generate code
	fmt.Printf("Querying Claude API...\n")
	client := codegen.NewClaudeClient(apiKey)
	code, err := client.GenerateCode(spec, boardType, systemPrompt, 4096)
	if err != nil {
		return fmt.Errorf("code generation failed: %w", err)
	}

	lineCount := len(strings.Split(code, "\n"))
	fmt.Printf("✓ Code generated (%d lines)\n\n", lineCount)

	// 5. Output
	if outputFile != "" {
		err = os.WriteFile(outputFile, []byte(code), 0644)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		fmt.Printf("Saved to: %s\n", outputFile)
	} else {
		fmt.Println("--- Generated Code ---")
		fmt.Println(code)
		fmt.Println("--- End Generated Code ---")
	}

	// 6. Suggest next steps
	fmt.Println("\n--- Next Steps ---")
	if outputFile != "" {
		fmt.Printf("1. Validate style: percepta style-check %s\n", outputFile)
		fmt.Printf("2. Review code for correctness\n")
		fmt.Printf("3. Flash to hardware and test\n")
		fmt.Printf("4. Observe behavior: percepta observe <device>\n")
		fmt.Printf("5. Store validated pattern: percepta knowledge store ...\n")
	} else {
		fmt.Printf("1. Save to file with --output flag\n")
		fmt.Printf("2. Validate style: percepta style-check <file>\n")
		fmt.Printf("3. Review, flash, and test on hardware\n")
	}

	return nil
}
