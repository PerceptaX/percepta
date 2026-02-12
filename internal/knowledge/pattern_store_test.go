package knowledge

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/perceptumx/percepta/internal/core"
)

func setupTestPatternStore(t *testing.T) (*PatternStore, func()) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Override HOME for testing
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	// Create .local/share/percepta directory
	dbDir := filepath.Join(tmpDir, ".local", "share", "percepta")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("failed to create test db dir: %v", err)
	}

	store, err := NewPatternStore()
	if err != nil {
		t.Fatalf("failed to create test pattern store: %v", err)
	}

	cleanup := func() {
		store.Close()
		os.Setenv("HOME", oldHome)
		os.RemoveAll(tmpDir)
	}

	return store, cleanup
}

func TestPatternStore_StoreValidatedPattern_Success(t *testing.T) {
	store, cleanup := setupTestPatternStore(t)
	defer cleanup()

	// First, store an observation that the pattern can reference
	obs := core.Observation{
		ID:           "obs-test-123",
		DeviceID:     "esp32-devkit-v1",
		FirmwareHash: "v1.0.0",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{
				Name:       "LED1",
				On:         true,
				BlinkHz:    1.0,
				Confidence: 0.95,
			},
		},
	}

	err := store.storage.Save(obs)
	if err != nil {
		t.Fatalf("failed to save observation: %v", err)
	}

	// Now store a validated pattern
	// This code is BARR-C compliant
	code := `#include <stdint.h>

void LED_Blink(void) {
    uint8_t status = 1;
}
`

	spec := "Blink LED at 1Hz"
	patternID, err := store.StoreValidatedPattern(spec, code, "esp32-devkit-v1", "v1.0.0")

	if err != nil {
		t.Fatalf("StoreValidatedPattern failed: %v", err)
	}

	if patternID == "" {
		t.Fatal("StoreValidatedPattern returned empty ID")
	}

	// Verify pattern was stored
	pattern, err := store.graph.GetPattern(patternID)
	if err != nil {
		t.Fatalf("failed to retrieve pattern: %v", err)
	}

	if pattern.Code != code {
		t.Errorf("expected code to match")
	}

	if pattern.Spec != spec {
		t.Errorf("expected spec %s, got %s", spec, pattern.Spec)
	}

	if pattern.BoardType != "esp32" {
		t.Errorf("expected board type esp32, got %s", pattern.BoardType)
	}

	if !pattern.StyleCompliant {
		t.Error("expected pattern to be style compliant")
	}

	// Verify relationships were created
	edges, err := store.graph.GetEdgesFrom(patternID, PRODUCES)
	if err != nil {
		t.Fatalf("failed to get edges: %v", err)
	}

	if len(edges) != 1 {
		t.Errorf("expected 1 PRODUCES edge, got %d", len(edges))
	}
}

func TestPatternStore_StoreValidatedPattern_StyleViolations(t *testing.T) {
	store, cleanup := setupTestPatternStore(t)
	defer cleanup()

	// Store observation first
	obs := core.Observation{
		ID:           "obs-test-456",
		DeviceID:     "esp32-devkit-v1",
		FirmwareHash: "v1.0.0",
		Timestamp:    time.Now(),
		Signals:      []core.Signal{},
	}

	err := store.storage.Save(obs)
	if err != nil {
		t.Fatalf("failed to save observation: %v", err)
	}

	// Code with BARR-C violations (wrong naming convention)
	code := `void blinkLED() {
    int x = 5;
}
`

	spec := "Blink LED"
	_, err = store.StoreValidatedPattern(spec, code, "esp32-devkit-v1", "v1.0.0")

	if err == nil {
		t.Fatal("expected error for non-compliant code, got nil")
	}

	// Verify error message mentions violations
	expectedMsg := "not BARR-C compliant"
	if err != nil && len(err.Error()) > 0 {
		if !contains(err.Error(), expectedMsg) {
			t.Errorf("expected error to contain %q, got %q", expectedMsg, err.Error())
		}
	}
}

func TestPatternStore_StoreValidatedPattern_NoObservation(t *testing.T) {
	store, cleanup := setupTestPatternStore(t)
	defer cleanup()

	// Try to store pattern without corresponding observation
	code := `#include <stdint.h>

void LED_Blink(void) {
    uint8_t status = 1;
}
`

	spec := "Blink LED"
	_, err := store.StoreValidatedPattern(spec, code, "esp32-devkit-v1", "v1.0.0")

	if err == nil {
		t.Fatal("expected error when no observation exists, got nil")
	}

	// Verify error message mentions observation
	expectedMsg := "no observation found"
	if !contains(err.Error(), expectedMsg) {
		t.Errorf("expected error to contain %q, got %q", expectedMsg, err.Error())
	}
}

func TestPatternStore_QueryPatternsByBoard(t *testing.T) {
	store, cleanup := setupTestPatternStore(t)
	defer cleanup()

	// Store observations for different devices
	esp32Obs := core.Observation{
		ID:           "obs-esp32",
		DeviceID:     "esp32-devkit-v1",
		FirmwareHash: "v1.0.0",
		Timestamp:    time.Now(),
		Signals:      []core.Signal{},
	}
	store.storage.Save(esp32Obs)

	stm32Obs := core.Observation{
		ID:           "obs-stm32",
		DeviceID:     "stm32-nucleo",
		FirmwareHash: "v1.0.0",
		Timestamp:    time.Now(),
		Signals:      []core.Signal{},
	}
	store.storage.Save(stm32Obs)

	// Store patterns for different boards
	esp32Code := `#include <stdint.h>

void ESP_Init(void) {
    uint8_t x = 0;
}
`

	stm32Code := `#include <stdint.h>

void STM_Init(void) {
    uint8_t y = 0;
}
`

	store.StoreValidatedPattern("ESP32 Init", esp32Code, "esp32-devkit-v1", "v1.0.0")
	store.StoreValidatedPattern("STM32 Init", stm32Code, "stm32-nucleo", "v1.0.0")

	// Query ESP32 patterns
	esp32Patterns, err := store.QueryPatternsByBoard("esp32")
	if err != nil {
		t.Fatalf("QueryPatternsByBoard failed: %v", err)
	}

	if len(esp32Patterns) != 1 {
		t.Fatalf("expected 1 ESP32 pattern, got %d", len(esp32Patterns))
	}

	if esp32Patterns[0].BoardType != "esp32" {
		t.Errorf("expected board type esp32, got %s", esp32Patterns[0].BoardType)
	}

	// Query STM32 patterns
	stm32Patterns, err := store.QueryPatternsByBoard("stm32")
	if err != nil {
		t.Fatalf("QueryPatternsByBoard failed: %v", err)
	}

	if len(stm32Patterns) != 1 {
		t.Fatalf("expected 1 STM32 pattern, got %d", len(stm32Patterns))
	}
}

func TestPatternStore_QueryPatternsBySpec(t *testing.T) {
	store, cleanup := setupTestPatternStore(t)
	defer cleanup()

	// Store observation
	obs := core.Observation{
		ID:           "obs-1",
		DeviceID:     "esp32-devkit-v1",
		FirmwareHash: "v1.0.0",
		Timestamp:    time.Now(),
		Signals:      []core.Signal{},
	}
	store.storage.Save(obs)

	// Store patterns with different specs
	code1 := `#include <stdint.h>

void LED_Blink1(void) {
    uint8_t x = 0;
}
`

	code2 := `#include <stdint.h>

void LED_Blink2(void) {
    uint8_t y = 0;
}
`

	spec1 := "Blink LED at 1Hz"
	spec2 := "Blink LED at 2Hz"

	store.StoreValidatedPattern(spec1, code1, "esp32-devkit-v1", "v1.0.0")
	store.StoreValidatedPattern(spec2, code2, "esp32-devkit-v1", "v1.0.0")

	// Query by spec
	patterns, err := store.QueryPatternsBySpec(spec1)
	if err != nil {
		t.Fatalf("QueryPatternsBySpec failed: %v", err)
	}

	if len(patterns) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(patterns))
	}

	if patterns[0].Spec != spec1 {
		t.Errorf("expected spec %s, got %s", spec1, patterns[0].Spec)
	}
}

func TestPatternStore_GetPatternWithObservation(t *testing.T) {
	store, cleanup := setupTestPatternStore(t)
	defer cleanup()

	// Store observation
	obs := core.Observation{
		ID:           "obs-full",
		DeviceID:     "esp32-devkit-v1",
		FirmwareHash: "v1.0.0",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{
				Name:       "LED1",
				On:         true,
				BlinkHz:    1.0,
				Confidence: 0.95,
			},
		},
	}
	store.storage.Save(obs)

	// Store pattern
	code := `#include <stdint.h>

void LED_Test(void) {
    uint8_t z = 0;
}
`

	patternID, err := store.StoreValidatedPattern("Test", code, "esp32-devkit-v1", "v1.0.0")
	if err != nil {
		t.Fatalf("StoreValidatedPattern failed: %v", err)
	}

	// Get pattern with observation
	pattern, observation, err := store.GetPatternWithObservation(patternID)
	if err != nil {
		t.Fatalf("GetPatternWithObservation failed: %v", err)
	}

	if pattern == nil {
		t.Fatal("expected pattern, got nil")
	}

	if observation == nil {
		t.Fatal("expected observation, got nil")
	}

	if pattern.Code != code {
		t.Error("pattern code mismatch")
	}

	if observation.ID != obs.ID {
		t.Errorf("expected observation ID %s, got %s", obs.ID, observation.ID)
	}

	if len(observation.Signals) != 1 {
		t.Errorf("expected 1 signal, got %d", len(observation.Signals))
	}
}

func TestPatternStore_Stats(t *testing.T) {
	store, cleanup := setupTestPatternStore(t)
	defer cleanup()

	// Store observation
	obs := core.Observation{
		ID:           "obs-stats",
		DeviceID:     "esp32-devkit-v1",
		FirmwareHash: "v1.0.0",
		Timestamp:    time.Now(),
		Signals:      []core.Signal{},
	}
	store.storage.Save(obs)

	// Store pattern
	code := `#include <stdint.h>

void Stats_Test(void) {
    uint8_t a = 0;
}
`

	store.StoreValidatedPattern("Stats Test", code, "esp32-devkit-v1", "v1.0.0")

	// Get stats
	stats := store.Stats()

	if stats["patterns"] < 1 {
		t.Errorf("expected at least 1 pattern, got %d", stats["patterns"])
	}

	if stats["observations"] < 1 {
		t.Errorf("expected at least 1 observation, got %d", stats["observations"])
	}

	if stats["edges"] < 4 {
		// Should have at least 4 edges (IMPLEMENTED_BY, RUNS_ON, PRODUCES, VALIDATED_BY)
		t.Errorf("expected at least 4 edges, got %d", stats["edges"])
	}
}

func TestGetBoardType(t *testing.T) {
	tests := []struct {
		deviceID  string
		boardType string
	}{
		{"esp32-devkit-v1", "esp32"},
		{"esp32-s3", "esp32"},
		{"stm32-nucleo", "stm32"},
		{"stm32f4", "stm32"},
		{"arduino-uno", "arduino"},
		{"rp2040", "rp2040"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		result := getBoardType(tt.deviceID)
		if result != tt.boardType {
			t.Errorf("getBoardType(%s) = %s, want %s", tt.deviceID, result, tt.boardType)
		}
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		len(s) > len(substr)*2 && s[len(s)/2-len(substr)/2:len(s)/2+len(substr)/2+1] == substr ||
		findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
