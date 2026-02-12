package main

import (
	"fmt"
	"os"

	"github.com/perceptumx/percepta/internal/assertions"
	"github.com/perceptumx/percepta/internal/config"
	perceptaErrors "github.com/perceptumx/percepta/internal/errors"
	"github.com/perceptumx/percepta/internal/storage"
	"github.com/perceptumx/percepta/internal/ui"
	"github.com/perceptumx/percepta/pkg/percepta"
	"github.com/spf13/cobra"
)

var assertCmd = &cobra.Command{
	Use:   "assert <device> <assertion>",
	Short: "Validate hardware state against expected behavior",
	Long: `Validate hardware behavior using assertions.

Captures an observation and evaluates the assertion expression. Returns exit
code 0 if passed, 1 if failed.

Examples:
  # LED is ON
  percepta assert my-board "led power is ON"

  # LED blinks at specific rate
  percepta assert my-board "led status blinks at 2Hz"

  # Display contains text
  percepta assert my-board "display LCD shows 'Ready'"

  # Multiple conditions
  percepta assert my-board "led power is ON" "led error is OFF"`,
	Args: cobra.MinimumNArgs(2),
	RunE: runAssert,
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
		return perceptaErrors.ConfigNotFound()
	}

	if len(cfg.Devices) == 0 {
		return perceptaErrors.NoDevicesConfigured()
	}

	deviceCfg, ok := cfg.Devices[deviceID]
	if !ok {
		return perceptaErrors.DeviceNotFound(deviceID)
	}

	cameraPath := "/dev/video0"
	firmwareTag := ""
	if deviceCfg.CameraID != "" {
		cameraPath = deviceCfg.CameraID
	}
	firmwareTag = deviceCfg.Firmware

	// Initialize SQLite storage
	sqliteStorage, err := storage.NewSQLiteStorage()
	if err != nil {
		return perceptaErrors.StorageInitFailed(err)
	}
	defer sqliteStorage.Close()

	perceptaCore, err := percepta.NewCore(cameraPath, sqliteStorage)
	if err != nil {
		return perceptaErrors.CameraNotFound(cameraPath)
	}

	// Capture observation with spinner
	spinner := ui.NewSpinner(fmt.Sprintf("Evaluating assertion on %s...", deviceID))
	obs, err := perceptaCore.Observe(deviceID)
	if err != nil {
		spinner.Stop(false)
		return perceptaErrors.ObservationFailed(err)
	}

	// Inject firmware tag and save
	obs.FirmwareHash = firmwareTag
	if err := sqliteStorage.Save(*obs); err != nil {
		spinner.Stop(false)
		return fmt.Errorf("failed to save observation: %w", err)
	}

	// Evaluate assertion
	result := assertion.Evaluate(obs)
	spinner.Stop(result.Passed)

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
