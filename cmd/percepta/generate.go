package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/perceptumx/percepta/internal/codegen"
	"github.com/perceptumx/percepta/internal/config"
	perceptaErrors "github.com/perceptumx/percepta/internal/errors"
	"github.com/perceptumx/percepta/internal/knowledge"
	"github.com/perceptumx/percepta/internal/style"
	"github.com/perceptumx/percepta/internal/ui"
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
	//nolint:errcheck // Flag name is hardcoded, cannot fail
	_ = generateCmd.MarkFlagRequired("board")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	spec := strings.Join(args, " ")

	// Check for API key
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return perceptaErrors.MissingAPIKey("Anthropic")
	}

	// Validate spec
	if spec == "" {
		return perceptaErrors.InvalidSpec(fmt.Errorf("specification cannot be empty"))
	}

	// Get device ID from config (for pattern linkage)
	cfg, err := config.Load()
	deviceID := ""
	if err == nil && cfg != nil && len(cfg.Devices) > 0 {
		// Use first device from config
		for id := range cfg.Devices {
			deviceID = id
			break
		}
	}
	if deviceID == "" {
		deviceID = "unknown-device" // Fallback for testing
	}

	fmt.Printf("Generating firmware...\n")
	fmt.Printf("Spec: %s\n", spec)
	fmt.Printf("Board: %s\n", boardType)
	fmt.Printf("Device: %s\n\n", deviceID)

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

	// 3. Initialize pipeline
	styleChecker := style.NewStyleChecker()
	styleFixer := style.NewStyleFixer()
	claudeClient := codegen.NewClaudeClient(apiKey)
	promptBuilder := codegen.NewPromptBuilder(patternStore)

	pipeline := codegen.NewGenerationPipeline(
		claudeClient,
		promptBuilder,
		styleChecker,
		styleFixer,
		patternStore,
	)

	// 4. Generate with validation
	spinner := ui.NewSpinner("Generating code with Claude...")
	result, err := pipeline.Generate(spec, boardType, deviceID)
	if err != nil {
		spinner.Stop(false)
		return perceptaErrors.CodeGenerationFailed(err)
	}
	spinner.Stop(true)

	fmt.Println()

	// 5. Print detailed report
	codegen.PrintGenerationReport(result)

	fmt.Println()

	// 6. Output code
	if outputFile != "" {
		err = os.WriteFile(outputFile, []byte(result.Code), 0644)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		fmt.Printf("\nSaved to: %s\n", outputFile)
	} else {
		fmt.Println("--- Generated Code ---")
		fmt.Println(result.Code)
		fmt.Println("--- End Generated Code ---")
	}

	// 7. Suggest next steps
	fmt.Println("\n--- Next Steps ---")
	if outputFile != "" {
		if !result.StyleCompliant {
			fmt.Printf("1. Fix remaining violations manually\n")
			fmt.Printf("2. Review code for correctness\n")
			fmt.Printf("3. Flash to hardware and test\n")
		} else {
			fmt.Printf("1. Review code for correctness\n")
			fmt.Printf("2. Flash to hardware and test\n")
			fmt.Printf("3. Observe behavior: percepta observe %s\n", deviceID)
		}
	} else {
		fmt.Printf("1. Save to file with --output flag\n")
		fmt.Printf("2. Review and test on hardware\n")
	}

	return nil
}
