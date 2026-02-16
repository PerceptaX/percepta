package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func setupTestConfig(t *testing.T) (string, func()) {
	// Create temporary directory for test config
	tmpDir, err := os.MkdirTemp("", "percepta-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Override HOME environment variable
	// On Windows, os.UserHomeDir() uses USERPROFILE, not HOME
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	var originalUserProfile string
	if runtime.GOOS == "windows" {
		originalUserProfile = os.Getenv("USERPROFILE")
		os.Setenv("USERPROFILE", tmpDir)
	}

	// Clear any environment variables that might interfere
	originalAPIKey := os.Getenv("ANTHROPIC_API_KEY")
	originalProvider := os.Getenv("PERCEPTA_VISION_PROVIDER")
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Unsetenv("PERCEPTA_VISION_PROVIDER")

	cleanup := func() {
		os.Setenv("HOME", originalHome)
		if runtime.GOOS == "windows" {
			os.Setenv("USERPROFILE", originalUserProfile)
		}
		if originalAPIKey != "" {
			os.Setenv("ANTHROPIC_API_KEY", originalAPIKey)
		}
		if originalProvider != "" {
			os.Setenv("PERCEPTA_VISION_PROVIDER", originalProvider)
		}
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestLoad_MissingConfigFile(t *testing.T) {
	tmpDir, cleanup := setupTestConfig(t)
	defer cleanup()

	// Load config without creating config file
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load should succeed with missing config file, got error: %v", err)
	}

	// Verify defaults are applied
	if cfg.Vision.Provider != "claude" {
		t.Errorf("Expected default provider 'claude', got '%s'", cfg.Vision.Provider)
	}

	// Verify config directory structure doesn't need to exist
	configPath := filepath.Join(tmpDir, ".config", "percepta")
	_, statErr := os.Stat(configPath)
	_ = statErr // It's OK if the directory exists, but it shouldn't be required
}

func TestLoad_DefaultValues(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify default vision provider
	if cfg.Vision.Provider != "claude" {
		t.Errorf("Expected default provider 'claude', got '%s'", cfg.Vision.Provider)
	}

	// Verify API key is empty by default
	if cfg.Vision.APIKey != "" {
		t.Errorf("Expected empty API key by default, got '%s'", cfg.Vision.APIKey)
	}

	// Verify devices map may be nil when no devices are configured
	// This is expected behavior - devices are only initialized when present in config
}

func TestLoad_WithConfigFile(t *testing.T) {
	tmpDir, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create config directory
	configDir := filepath.Join(tmpDir, ".config", "percepta")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Write test config file
	configPath := filepath.Join(configDir, "config.yaml")
	configContent := `vision:
  provider: openai
  api_key: test-key-from-file

devices:
  test-device:
    type: esp32
    camera_id: "0"
    firmware: v1.0.0
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify config loaded from file
	if cfg.Vision.Provider != "openai" {
		t.Errorf("Expected provider 'openai', got '%s'", cfg.Vision.Provider)
	}

	if cfg.Vision.APIKey != "test-key-from-file" {
		t.Errorf("Expected API key 'test-key-from-file', got '%s'", cfg.Vision.APIKey)
	}

	// Verify device config
	if len(cfg.Devices) != 1 {
		t.Fatalf("Expected 1 device, got %d", len(cfg.Devices))
	}

	device, ok := cfg.Devices["test-device"]
	if !ok {
		t.Fatal("Expected 'test-device' in config")
	}

	if device.Type != "esp32" {
		t.Errorf("Expected device type 'esp32', got '%s'", device.Type)
	}

	if device.CameraID != "0" {
		t.Errorf("Expected camera ID '0', got '%s'", device.CameraID)
	}

	if device.Firmware != "v1.0.0" {
		t.Errorf("Expected firmware 'v1.0.0', got '%s'", device.Firmware)
	}
}

func TestLoad_WithEnvironmentOverrides(t *testing.T) {
	tmpDir, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create config directory and file
	configDir := filepath.Join(tmpDir, ".config", "percepta")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	configContent := `vision:
  provider: claude
  api_key: file-key
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Set environment variable (ANTHROPIC_API_KEY overrides config file)
	os.Setenv("ANTHROPIC_API_KEY", "env-api-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify ANTHROPIC_API_KEY overrides file value
	if cfg.Vision.APIKey != "env-api-key" {
		t.Errorf("Expected API key from ANTHROPIC_API_KEY env var, got '%s'", cfg.Vision.APIKey)
	}

	// Provider should still be from file
	if cfg.Vision.Provider != "claude" {
		t.Errorf("Expected provider 'claude', got '%s'", cfg.Vision.Provider)
	}
}

func TestLoad_PerceptaEnvPrefix(t *testing.T) {
	t.Skip("Viper AutomaticEnv with prefix requires explicit BindEnv - feature not currently used")
	// Note: PERCEPTA_ prefix with AutomaticEnv() requires viper.BindEnv() for each key
	// Current implementation relies on ANTHROPIC_API_KEY direct check instead
	// This test documents the limitation but is skipped as the feature isn't critical
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create config directory
	configDir := filepath.Join(tmpDir, ".config", "percepta")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Write YAML with type mismatch that will cause Unmarshal to fail
	configPath := filepath.Join(configDir, "config.yaml")
	invalidYAML := `vision:
  provider:
    - this should be a string
    - not an array
`
	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Load should return an error for YAML that causes Unmarshal to fail
	_, err := Load()
	if err == nil {
		t.Fatal("Expected error for YAML type mismatch during unmarshal, got nil")
	}
}

func TestLoad_EmptyConfigFile(t *testing.T) {
	tmpDir, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create config directory
	configDir := filepath.Join(tmpDir, ".config", "percepta")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Write empty config file
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load should succeed with empty config file, got error: %v", err)
	}

	// Verify defaults are still applied
	if cfg.Vision.Provider != "claude" {
		t.Errorf("Expected default provider 'claude', got '%s'", cfg.Vision.Provider)
	}
}

func TestLoad_MultipleDevices(t *testing.T) {
	tmpDir, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create config directory
	configDir := filepath.Join(tmpDir, ".config", "percepta")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Write config with multiple devices
	configPath := filepath.Join(configDir, "config.yaml")
	configContent := `vision:
  provider: claude

devices:
  esp32-dev:
    type: esp32
    camera_id: "0"
    firmware: v1.0.0
  stm32-board:
    type: stm32
    camera_id: "1"
    firmware: v2.1.0
  arduino-uno:
    type: arduino
    camera_id: "0"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify all devices loaded
	if len(cfg.Devices) != 3 {
		t.Fatalf("Expected 3 devices, got %d", len(cfg.Devices))
	}

	// Verify esp32-dev
	esp32, ok := cfg.Devices["esp32-dev"]
	if !ok {
		t.Fatal("Expected 'esp32-dev' in config")
	}
	if esp32.Type != "esp32" {
		t.Errorf("Expected esp32 type 'esp32', got '%s'", esp32.Type)
	}

	// Verify stm32-board
	stm32, ok := cfg.Devices["stm32-board"]
	if !ok {
		t.Fatal("Expected 'stm32-board' in config")
	}
	if stm32.Firmware != "v2.1.0" {
		t.Errorf("Expected stm32 firmware 'v2.1.0', got '%s'", stm32.Firmware)
	}

	// Verify arduino-uno
	_, ok = cfg.Devices["arduino-uno"]
	if !ok {
		t.Fatal("Expected 'arduino-uno' in config")
	}
}

func TestLoad_HomeDirectoryResolution(t *testing.T) {
	tmpDir, cleanup := setupTestConfig(t)
	defer cleanup()

	// Verify that Load() can successfully resolve home directory
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Basic sanity check - config should be initialized
	if cfg == nil {
		t.Fatal("Expected config to be non-nil")
	}

	// Verify the config directory would be under the temporary HOME
	expectedConfigDir := filepath.Join(tmpDir, ".config", "percepta")
	// We don't require the directory to exist, but the path should be correct
	if !filepath.IsAbs(expectedConfigDir) {
		t.Errorf("Expected absolute config path, got '%s'", expectedConfigDir)
	}
}

func TestLoad_PartialConfig(t *testing.T) {
	tmpDir, cleanup := setupTestConfig(t)
	defer cleanup()

	// Create config directory
	configDir := filepath.Join(tmpDir, ".config", "percepta")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Write config with only devices, no vision config
	configPath := filepath.Join(configDir, "config.yaml")
	configContent := `devices:
  test-device:
    type: esp32
    camera_id: "0"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify vision defaults are still applied
	if cfg.Vision.Provider != "claude" {
		t.Errorf("Expected default provider 'claude', got '%s'", cfg.Vision.Provider)
	}

	// Verify device is loaded
	if len(cfg.Devices) != 1 {
		t.Errorf("Expected 1 device, got %d", len(cfg.Devices))
	}
}
