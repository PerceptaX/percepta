package main

import (
	"fmt"

	"github.com/perceptumx/percepta/internal/config"
	"github.com/spf13/cobra"
)

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Manage device configurations",
	Long:  "Add, list, and configure devices for observation.",
}

var deviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured devices",
	Long:  "Displays all devices configured in ~/.config/percepta/config.yaml",
	RunE:  runDeviceList,
}

func init() {
	deviceCmd.AddCommand(deviceListCmd)
}

func runDeviceList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config load failed: %w", err)
	}

	if len(cfg.Devices) == 0 {
		fmt.Println("No devices configured. Add one with: percepta device add <name>")
		return nil
	}

	fmt.Println("Configured devices:")
	fmt.Println()

	for name, dev := range cfg.Devices {
		fmt.Println(name)
		if dev.Type != "" {
			fmt.Printf("  Type: %s\n", dev.Type)
		}
		if dev.CameraID != "" {
			fmt.Printf("  Camera: %s\n", dev.CameraID)
		}
		if dev.Firmware != "" {
			fmt.Printf("  Firmware: %s\n", dev.Firmware)
		}
		fmt.Println()
	}

	return nil
}
