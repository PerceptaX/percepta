//go:build linux || darwin

package main

import (
	"fmt"
	"time"

	"github.com/perceptumx/percepta/internal/config"
	"github.com/perceptumx/percepta/internal/core"
	perceptaErrors "github.com/perceptumx/percepta/internal/errors"
	"github.com/perceptumx/percepta/internal/storage"
	"github.com/perceptumx/percepta/internal/ui"
	"github.com/perceptumx/percepta/pkg/percepta"
	"github.com/spf13/cobra"
)

var observeCmd = &cobra.Command{
	Use:   "observe <device>",
	Short: "Observe hardware state via computer vision",
	Long: `Capture and analyze hardware behavior using vision.

Percepta uses your camera to observe LED states, display content, and boot
behavior. Observations are stored in SQLite for diffing and assertions.

Examples:
  # Observe device with default camera
  percepta observe my-esp32

  # Observe with specific camera
  percepta observe my-esp32 --camera /dev/video0

  # Observe for 10 seconds
  percepta observe my-esp32 --duration 10s

  # Save observation to file
  percepta observe my-esp32 --output observation.json`,
	Args: cobra.ExactArgs(1),
	RunE: runObserve,
}

func runObserve(cmd *cobra.Command, args []string) error {
	deviceID := args[0]

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return perceptaErrors.ConfigNotFound()
	}

	// Check if any devices configured
	if len(cfg.Devices) == 0 {
		return perceptaErrors.NoDevicesConfigured()
	}

	// Get camera path and firmware tag for device
	cameraPath := "/dev/video0" // Default
	firmwareTag := ""
	deviceCfg, ok := cfg.Devices[deviceID]
	if !ok {
		return perceptaErrors.DeviceNotFound(deviceID)
	}

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

	// Initialize Core with storage
	perceptaCore, err := percepta.NewCore(cameraPath, sqliteStorage)
	if err != nil {
		return perceptaErrors.CameraNotFound(cameraPath)
	}

	// Capture observation with spinner
	spinner := ui.NewSpinner(fmt.Sprintf("Capturing frames from %s...", deviceID))
	obs, err := perceptaCore.Observe(deviceID)
	if err != nil {
		spinner.Stop(false)
		return perceptaErrors.ObservationFailed(err)
	}
	spinner.Stop(true)

	// Inject firmware tag
	obs.FirmwareHash = firmwareTag

	// Save observation with firmware tag
	if err := sqliteStorage.Save(*obs); err != nil {
		return fmt.Errorf("failed to save observation: %w", err)
	}

	// Format output
	printObservation(obs, perceptaCore.ObservationCount())

	return nil
}

func printObservation(obs *core.Observation, count int) {
	fmt.Printf("✅ Observation captured: %s\n", obs.ID)
	fmt.Printf("Device: %s\n", obs.DeviceID)
	fmt.Printf("Timestamp: %s\n", obs.Timestamp.Format(time.RFC3339))
	fmt.Printf("\n")

	if len(obs.Signals) == 0 {
		fmt.Println("No signals detected")
		return
	}

	fmt.Printf("Signals (%d):\n", len(obs.Signals))
	for i, signal := range obs.Signals {
		switch s := signal.(type) {
		case core.LEDSignal:
			state := "OFF"
			if s.On {
				state = "ON"
			}
			fmt.Printf("  %d. LED '%s': %s", i+1, s.Name, state)
			if s.BlinkHz > 0 {
				fmt.Printf(" (blinking at %.2f Hz)", s.BlinkHz)
			}
			if s.Color.R > 0 || s.Color.G > 0 || s.Color.B > 0 {
				fmt.Printf(" [RGB(%d,%d,%d)]", s.Color.R, s.Color.G, s.Color.B)
			}
			fmt.Printf(" [confidence: %.2f]\n", s.Confidence)

		case core.DisplaySignal:
			fmt.Printf("  %d. Display '%s': \"%s\" [confidence: %.2f]\n",
				i+1, s.Name, s.Text, s.Confidence)

		case core.BootTimingSignal:
			fmt.Printf("  %d. Boot timing: %dms [confidence: %.2f]\n",
				i+1, s.DurationMs, s.Confidence)
		}
	}

	fmt.Printf("\nStored in memory (%d total observations)\n", count)
}
