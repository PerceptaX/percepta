package main

import (
	"fmt"
	"os"
	"time"

	"github.com/perceptumx/percepta/internal/config"
	"github.com/perceptumx/percepta/internal/core"
	"github.com/perceptumx/percepta/internal/storage"
	"github.com/perceptumx/percepta/pkg/percepta"
	"github.com/spf13/cobra"
)

var observeCmd = &cobra.Command{
	Use:   "observe <device>",
	Short: "Capture current hardware state",
	Long:  "Captures webcam frame, analyzes with Claude Vision, and stores observation.",
	Args:  cobra.ExactArgs(1),
	RunE:  runObserve,
}

func runObserve(cmd *cobra.Command, args []string) error {
	deviceID := args[0]

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config load failed: %w", err)
	}

	// Get camera path and firmware tag for device
	cameraPath := "/dev/video0" // Default
	firmwareTag := ""
	if deviceCfg, ok := cfg.Devices[deviceID]; ok {
		if deviceCfg.CameraID != "" {
			cameraPath = deviceCfg.CameraID
		}
		firmwareTag = deviceCfg.Firmware
	}

	// Initialize SQLite storage
	sqliteStorage, err := storage.NewSQLiteStorage()
	if err != nil {
		return fmt.Errorf("storage init failed: %w", err)
	}
	defer sqliteStorage.Close()

	// Initialize Core with storage
	perceptaCore, err := percepta.NewCore(cameraPath, sqliteStorage)
	if err != nil {
		return err
	}

	// Capture observation
	fmt.Fprintf(os.Stderr, "Observing %s via %s...\n", deviceID, cameraPath)
	obs, err := perceptaCore.Observe(deviceID)
	if err != nil {
		return err
	}

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
