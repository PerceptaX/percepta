package main

import (
	"fmt"
	"os"

	"github.com/perceptumx/percepta/internal/assertions"
	"github.com/perceptumx/percepta/internal/config"
	"github.com/perceptumx/percepta/pkg/percepta"
	"github.com/spf13/cobra"
)

var assertCmd = &cobra.Command{
	Use:   "assert <device> <assertion>",
	Short: "Validate hardware state against expected behavior",
	Long:  "Captures observation and evaluates assertion. Returns 0 if passed, 1 if failed.",
	Args:  cobra.ExactArgs(2),
	RunE:  runAssert,
}

func runAssert(cmd *cobra.Command, args []string) error {
	deviceID := args[0]
	assertionDSL := args[1]

	// Parse assertion DSL
	assertion, err := assertions.Parse(assertionDSL)
	if err != nil {
		return fmt.Errorf("invalid assertion: %w", err)
	}

	// Load config and initialize Core
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config load failed: %w", err)
	}

	cameraPath := "/dev/video0"
	if deviceCfg, ok := cfg.Devices[deviceID]; ok {
		if deviceCfg.CameraID != "" {
			cameraPath = deviceCfg.CameraID
		}
	}

	perceptaCore, err := percepta.NewCore(cameraPath)
	if err != nil {
		return err
	}

	// Capture observation
	fmt.Fprintf(os.Stderr, "Observing %s (evaluating assertion)...\n", deviceID)
	obs, err := perceptaCore.Observe(deviceID)
	if err != nil {
		return err
	}

	// Evaluate assertion
	result := assertion.Evaluate(obs)

	// Format and print result
	printAssertionResult(assertion, result)

	// Exit with appropriate code
	if !result.Passed {
		os.Exit(1)
	}

	return nil
}

func printAssertionResult(assertion assertions.Assertion, result assertions.AssertionResult) {
	// Header with pass/fail indicator
	if result.Passed {
		fmt.Printf("✅ PASS: %s\n", assertion.String())
	} else {
		fmt.Printf("❌ FAIL: %s\n", assertion.String())
	}

	// Details section
	fmt.Printf("\nExpected: %s\n", result.Expected)
	fmt.Printf("Actual:   %s\n", result.Actual)

	// Confidence indicator
	fmt.Printf("Confidence: %.2f\n", result.Confidence)

	// Additional message if present
	if result.Message != "" {
		fmt.Printf("\nDetails: %s\n", result.Message)
	}
}
