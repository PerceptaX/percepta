package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/perceptumx/percepta/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Manage device configurations",
	Long: `Add, list, and configure devices for observation.

Devices associate a hardware identifier with a camera path and optional
firmware tag. This allows Percepta to track behavior across firmware versions.

Examples:
  # Add a new device
  percepta device add my-esp32

  # List all devices
  percepta device list

  # Set firmware version
  percepta device set-firmware my-esp32 v1.0.0`,
}

var deviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured devices",
	Long: `Displays all devices configured in ~/.config/percepta/config.yaml.

Shows device name, type, camera path, and firmware tag for each device.`,
	RunE: runDeviceList,
}

var deviceAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new device configuration",
	Long: `Interactively configure a new device for observation.

You will be prompted for:
- Device type (e.g., esp32, stm32, arduino)
- Camera device path (e.g., /dev/video0)
- Firmware version (optional)

Example:
  percepta device add my-board`,
	Args: cobra.ExactArgs(1),
	RunE: runDeviceAdd,
}

var deviceSetFirmwareCmd = &cobra.Command{
	Use:   "set-firmware <device> <firmware>",
	Short: "Update firmware tag for a device",
	Long: `Sets the firmware version tag for an existing device.

Use this before running observations to associate them with a specific
firmware version. Enables firmware diffing with 'percepta diff'.

Examples:
  percepta device set-firmware my-esp32 v1.0
  percepta device set-firmware my-board baseline
  percepta device set-firmware test-board abc123`,
	Args: cobra.ExactArgs(2),
	RunE: runDeviceSetFirmware,
}

func init() {
	deviceCmd.AddCommand(deviceListCmd)
	deviceCmd.AddCommand(deviceAddCmd)
	deviceCmd.AddCommand(deviceSetFirmwareCmd)
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

func runDeviceAdd(cmd *cobra.Command, args []string) error {
	deviceName := args[0]

	// Load existing config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config load failed: %w", err)
	}

	// Check if device already exists
	if cfg.Devices == nil {
		cfg.Devices = make(map[string]config.DeviceConfig)
	}
	if _, exists := cfg.Devices[deviceName]; exists {
		return fmt.Errorf("Device '%s' already exists. Use 'percepta device list' to see all devices.", deviceName)
	}

	scanner := bufio.NewScanner(os.Stdin)

	// Prompt for device type
	fmt.Print("Device type (e.g., fpga, esp32, stm32): ")
	scanner.Scan()
	deviceType := strings.TrimSpace(scanner.Text())

	// Prompt for camera path
	fmt.Print("Camera device path (default: /dev/video0): ")
	scanner.Scan()
	cameraPath := strings.TrimSpace(scanner.Text())
	if cameraPath == "" {
		cameraPath = "/dev/video0"
	}

	// Prompt for firmware version (optional)
	fmt.Print("Firmware version (optional, press Enter to skip): ")
	scanner.Scan()
	firmware := strings.TrimSpace(scanner.Text())

	// Create device config
	cfg.Devices[deviceName] = config.DeviceConfig{
		Type:     deviceType,
		CameraID: cameraPath,
		Firmware: firmware,
	}

	// Save config
	if err := saveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("\n✓ Device '%s' added successfully\n", deviceName)
	return nil
}

func runDeviceSetFirmware(cmd *cobra.Command, args []string) error {
	deviceName := args[0]
	newFirmware := args[1]

	// Load existing config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config load failed: %w", err)
	}

	// Check if device exists
	dev, exists := cfg.Devices[deviceName]
	if !exists {
		return fmt.Errorf("Device '%s' not found. Use 'percepta device list' to see all devices.", deviceName)
	}

	oldFirmware := dev.Firmware

	// Update firmware
	dev.Firmware = newFirmware
	cfg.Devices[deviceName] = dev

	// Save config
	if err := saveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Show confirmation
	if oldFirmware == "" {
		fmt.Printf("Set firmware for '%s': %s\n", deviceName, newFirmware)
	} else {
		fmt.Printf("Updated firmware for '%s': %s → %s\n", deviceName, oldFirmware, newFirmware)
	}

	return nil
}

func saveConfig(cfg *config.Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := fmt.Sprintf("%s/.config/percepta", homeDir)
	configFile := fmt.Sprintf("%s/config.yaml", configPath)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set Viper values
	viper.Set("vision.provider", cfg.Vision.Provider)
	viper.Set("vision.api_key", cfg.Vision.APIKey)
	viper.Set("devices", cfg.Devices)

	// Write config file
	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Set file permissions
	if err := os.Chmod(configFile, 0644); err != nil {
		return fmt.Errorf("failed to set config file permissions: %w", err)
	}

	return nil
}
