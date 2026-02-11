package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/perceptumx/percepta/internal/core"
)

func setupTestDB(t *testing.T) (*SQLiteStorage, func()) {
	// Create temporary directory for test database
	tmpDir, err := os.MkdirTemp("", "percepta-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	storage, err := NewSQLiteStorage()
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create storage: %v", err)
	}

	cleanup := func() {
		storage.Close()
		os.Setenv("HOME", originalHome)
		os.RemoveAll(tmpDir)
	}

	return storage, cleanup
}

func TestSQLiteStorage_SaveAndQuery(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	obs := core.Observation{
		ID:           "test-obs-1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, Color: core.RGB{R: 255, G: 0, B: 0}},
		},
	}

	// Save observation
	err := storage.Save(obs)
	if err != nil {
		t.Fatalf("Failed to save observation: %v", err)
	}

	// Query observations
	results, err := storage.Query("fpga", 10)
	if err != nil {
		t.Fatalf("Failed to query observations: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 observation, got %d", len(results))
	}

	if results[0].ID != obs.ID {
		t.Errorf("Expected ID %s, got %s", obs.ID, results[0].ID)
	}

	if results[0].FirmwareHash != "v1" {
		t.Errorf("Expected firmware v1, got %s", results[0].FirmwareHash)
	}
}

func TestSQLiteStorage_QueryByFirmware(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Save observations with different firmware versions
	obs1 := core.Observation{
		ID:           "obs-v1-1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true},
		},
	}

	obs2 := core.Observation{
		ID:           "obs-v2-1",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now().Add(1 * time.Minute),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: false},
		},
	}

	obs3 := core.Observation{
		ID:           "obs-v1-2",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now().Add(2 * time.Minute),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true},
		},
	}

	storage.Save(obs1)
	storage.Save(obs2)
	storage.Save(obs3)

	// Query by firmware v1
	v1Results, err := storage.QueryByFirmware("fpga", "v1", 10)
	if err != nil {
		t.Fatalf("Failed to query by firmware: %v", err)
	}

	if len(v1Results) != 2 {
		t.Errorf("Expected 2 observations for v1, got %d", len(v1Results))
	}

	// Query by firmware v2
	v2Results, err := storage.QueryByFirmware("fpga", "v2", 10)
	if err != nil {
		t.Fatalf("Failed to query by firmware: %v", err)
	}

	if len(v2Results) != 1 {
		t.Errorf("Expected 1 observation for v2, got %d", len(v2Results))
	}
}

func TestSQLiteStorage_GetLatestForFirmware(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Save multiple observations for same firmware
	now := time.Now()

	obs1 := core.Observation{
		ID:           "obs-1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    now,
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true},
		},
	}

	obs2 := core.Observation{
		ID:           "obs-2",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    now.Add(10 * time.Minute),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: false},
		},
	}

	obs3 := core.Observation{
		ID:           "obs-3",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    now.Add(5 * time.Minute),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true},
		},
	}

	storage.Save(obs1)
	storage.Save(obs2)
	storage.Save(obs3)

	// Get latest observation
	latest, err := storage.GetLatestForFirmware("fpga", "v1")
	if err != nil {
		t.Fatalf("Failed to get latest observation: %v", err)
	}

	if latest.ID != "obs-2" {
		t.Errorf("Expected latest observation to be obs-2, got %s", latest.ID)
	}
}

func TestSQLiteStorage_GetLatestForFirmware_NotFound(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Try to get observation for non-existent firmware
	_, err := storage.GetLatestForFirmware("fpga", "non-existent")
	if err == nil {
		t.Fatal("Expected error for non-existent firmware, got nil")
	}
}

func TestSQLiteStorage_EmptyFirmwareTag(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Save observation with empty firmware tag
	obs := core.Observation{
		ID:           "obs-no-fw",
		DeviceID:     "fpga",
		FirmwareHash: "",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true},
		},
	}

	err := storage.Save(obs)
	if err != nil {
		t.Fatalf("Failed to save observation with empty firmware: %v", err)
	}

	// Query by empty firmware
	results, err := storage.QueryByFirmware("fpga", "", 10)
	if err != nil {
		t.Fatalf("Failed to query by empty firmware: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 observation with empty firmware, got %d", len(results))
	}
}

func TestSQLiteStorage_SignalDeserialization(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Save observation with different signal types
	obs := core.Observation{
		ID:           "obs-multi-signal",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{
				Name:       "LED1",
				On:         true,
				Color:      core.RGB{R: 255, G: 0, B: 0},
				BlinkHz:    2.5,
				Confidence: 0.95,
			},
			core.DisplaySignal{
				Name:       "LCD1",
				Text:       "Hello World",
				Confidence: 0.90,
			},
			core.BootTimingSignal{
				DurationMs: 1500,
				Confidence: 0.85,
			},
		},
	}

	err := storage.Save(obs)
	if err != nil {
		t.Fatalf("Failed to save observation: %v", err)
	}

	// Retrieve and verify
	results, err := storage.Query("fpga", 1)
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	retrieved := results[0]

	if len(retrieved.Signals) != 3 {
		t.Fatalf("Expected 3 signals, got %d", len(retrieved.Signals))
	}

	// Verify LED signal
	led, ok := retrieved.Signals[0].(core.LEDSignal)
	if !ok {
		t.Fatal("First signal should be LEDSignal")
	}
	if led.Name != "LED1" || !led.On || led.BlinkHz != 2.5 {
		t.Errorf("LED signal not deserialized correctly: %+v", led)
	}

	// Verify Display signal
	display, ok := retrieved.Signals[1].(core.DisplaySignal)
	if !ok {
		t.Fatal("Second signal should be DisplaySignal")
	}
	if display.Name != "LCD1" || display.Text != "Hello World" {
		t.Errorf("Display signal not deserialized correctly: %+v", display)
	}

	// Verify Boot timing signal
	boot, ok := retrieved.Signals[2].(core.BootTimingSignal)
	if !ok {
		t.Fatal("Third signal should be BootTimingSignal")
	}
	if boot.DurationMs != 1500 {
		t.Errorf("Boot signal not deserialized correctly: %+v", boot)
	}
}

func TestSQLiteStorage_DatabasePath(t *testing.T) {
	// Verify database is created at correct path
	tmpDir, err := os.MkdirTemp("", "percepta-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	storage, err := NewSQLiteStorage()
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	expectedPath := filepath.Join(tmpDir, ".local", "share", "percepta", "percepta.db")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Database not created at expected path: %s", expectedPath)
	}
}
