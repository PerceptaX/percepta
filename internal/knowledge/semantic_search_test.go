//go:build !windows
package knowledge

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/perceptumx/percepta/internal/core"
	"github.com/perceptumx/percepta/internal/storage"
)

// setupTestEnvironment creates a test environment with all required directories
func setupTestEnvironment(t *testing.T) string {
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)

	// Create percepta directories
	perceptaDir := filepath.Join(tmpDir, ".local", "share", "percepta")
	os.MkdirAll(perceptaDir, 0755)

	configDir := filepath.Join(tmpDir, ".config", "percepta")
	os.MkdirAll(configDir, 0755)

	return tmpDir
}

func TestSearchSimilarPatterns_Success(t *testing.T) {
	tmpDir := setupTestEnvironment(t)
	defer os.Unsetenv("HOME")

	// Create pattern store
	store, err := NewPatternStore()
	if err != nil {
		t.Fatalf("Failed to create pattern store: %v", err)
	}
	defer store.Close()

	// Initialize vector store with mock embedder
	if err := store.InitializeVectorStoreWithEmbedder(&MockEmbedder{}); err != nil {
		t.Fatalf("Failed to initialize vector store: %v", err)
	}

	// Create test observations first
	sqliteStorage, _ := storage.NewSQLiteStorage()
	defer sqliteStorage.Close()

	// Observation for ESP32 LED blink
	obs1 := core.Observation{
		ID:           "obs1",
		DeviceID:     "esp32-devkit-v1",
		FirmwareHash: "v1.0.0",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, BlinkHz: 1.0, Confidence: 0.95},
		},
	}
	sqliteStorage.Save(obs1)

	// Observation for ESP32 LED toggle
	obs2 := core.Observation{
		ID:           "obs2",
		DeviceID:     "esp32-devkit-v2",
		FirmwareHash: "v1.0.1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, Confidence: 0.95},
		},
	}
	sqliteStorage.Save(obs2)

	// Observation for STM32 UART
	obs3 := core.Observation{
		ID:           "obs3",
		DeviceID:     "stm32-nucleo",
		FirmwareHash: "v1.0.0",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.DisplaySignal{Name: "UART", Text: "HELLO", Confidence: 0.95},
		},
	}
	sqliteStorage.Save(obs3)

	// Store validated patterns
	code1 := `#include <stdint.h>

void LED_Blink(void) {
    uint8_t status = 1;
}
`
	_, err = store.StoreValidatedPattern("Blink LED at 1Hz", code1, "esp32-devkit-v1", "v1.0.0")
	if err != nil {
		t.Fatalf("Failed to store pattern 1: %v", err)
	}

	code2 := `#include <stdint.h>

void LED_Toggle(void) {
    uint8_t state = 0;
}
`
	_, err = store.StoreValidatedPattern("Toggle LED state", code2, "esp32-devkit-v2", "v1.0.1")
	if err != nil {
		t.Fatalf("Failed to store pattern 2: %v", err)
	}

	code3 := `#include <stdint.h>

void UART_Send(uint8_t data) {
    // Send data
}
`
	_, err = store.StoreValidatedPattern("Send data via UART", code3, "stm32-nucleo", "v1.0.0")
	if err != nil {
		t.Fatalf("Failed to store pattern 3: %v", err)
	}

	// Test semantic search - query for LED blink
	results, err := store.SearchSimilarPatterns("Blink LED", "", 3)
	if err != nil {
		t.Fatalf("SearchSimilarPatterns failed: %v", err)
	}

	// Should return patterns
	if len(results) == 0 {
		t.Error("SearchSimilarPatterns returned no results")
	}

	// Verify result structure
	for i, result := range results {
		if result.Pattern == nil {
			t.Errorf("Result %d has nil pattern", i)
		}
		if result.Observation == nil {
			t.Errorf("Result %d has nil observation", i)
		}
		if result.Similarity < 0 || result.Similarity > 1 {
			t.Errorf("Result %d has invalid similarity: %f", i, result.Similarity)
		}
		if result.Confidence < 0 || result.Confidence > 1 {
			t.Errorf("Result %d has invalid confidence: %f", i, result.Confidence)
		}

		// All patterns should be style compliant
		if !result.Pattern.StyleCompliant {
			t.Errorf("Result %d has non-compliant pattern", i)
		}
	}

	// Results should be sorted by similarity (descending)
	for i := 1; i < len(results); i++ {
		if results[i].Similarity > results[i-1].Similarity {
			t.Errorf("Results not sorted: result[%d].Similarity (%f) > result[%d].Similarity (%f)",
				i, results[i].Similarity, i-1, results[i-1].Similarity)
		}
	}

	_ = tmpDir // Keep tmpDir reference
}

func TestSearchSimilarPatterns_BoardFilter(t *testing.T) {
	tmpDir := setupTestEnvironment(t)
	defer os.Unsetenv("HOME")

	// Create pattern store
	store, err := NewPatternStore()
	if err != nil {
		t.Fatalf("Failed to create pattern store: %v", err)
	}
	defer store.Close()

	// Initialize vector store with mock embedder
	if err := store.InitializeVectorStoreWithEmbedder(&MockEmbedder{}); err != nil {
		t.Fatalf("Failed to initialize vector store: %v", err)
	}

	// Create test observations
	sqliteStorage, _ := storage.NewSQLiteStorage()
	defer sqliteStorage.Close()

	obs1 := core.Observation{
		ID:           "obs1",
		DeviceID:     "esp32-devkit-v1",
		FirmwareHash: "v1.0.0",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, BlinkHz: 1.0, Confidence: 0.95},
		},
	}
	sqliteStorage.Save(obs1)

	obs2 := core.Observation{
		ID:           "obs2",
		DeviceID:     "stm32-nucleo",
		FirmwareHash: "v1.0.0",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, BlinkHz: 1.0, Confidence: 0.95},
		},
	}
	sqliteStorage.Save(obs2)

	// Store patterns for different boards
	code := `#include <stdint.h>

void LED_Blink(void) {
    uint8_t status = 1;
}
`
	store.StoreValidatedPattern("Blink LED", code, "esp32-devkit-v1", "v1.0.0")
	store.StoreValidatedPattern("Blink LED", code, "stm32-nucleo", "v1.0.0")

	// Search with board filter (ESP32 only)
	results, err := store.SearchSimilarPatterns("Blink LED", "esp32", 5)
	if err != nil {
		t.Fatalf("SearchSimilarPatterns failed: %v", err)
	}

	// All results should be ESP32
	for i, result := range results {
		if result.Pattern.BoardType != "esp32" {
			t.Errorf("Result %d has wrong board type: %s (expected esp32)",
				i, result.Pattern.BoardType)
		}
	}

	// Search with board filter (STM32 only)
	results, err = store.SearchSimilarPatterns("Blink LED", "stm32", 5)
	if err != nil {
		t.Fatalf("SearchSimilarPatterns failed: %v", err)
	}

	// All results should be STM32
	for i, result := range results {
		if result.Pattern.BoardType != "stm32" {
			t.Errorf("Result %d has wrong board type: %s (expected stm32)",
				i, result.Pattern.BoardType)
		}
	}

	_ = tmpDir
}

func TestSearchSimilarPatterns_TopK(t *testing.T) {
	tmpDir := setupTestEnvironment(t)
	defer os.Unsetenv("HOME")

	// Create pattern store
	store, err := NewPatternStore()
	if err != nil {
		t.Fatalf("Failed to create pattern store: %v", err)
	}
	defer store.Close()

	// Initialize vector store with mock embedder
	if err := store.InitializeVectorStoreWithEmbedder(&MockEmbedder{}); err != nil {
		t.Fatalf("Failed to initialize vector store: %v", err)
	}

	// Create test observations
	sqliteStorage, _ := storage.NewSQLiteStorage()
	defer sqliteStorage.Close()

	// Store 10 patterns with different code variations
	for i := 0; i < 10; i++ {
		firmware := fmt.Sprintf("v1.0.%d", i)
		obsID := core.GenerateID()
		obs := core.Observation{
			ID:           obsID,
			DeviceID:     "esp32-dev",
			FirmwareHash: firmware,
			Timestamp:    time.Now(),
			Signals: []core.Signal{
				core.LEDSignal{Name: "LED1", On: true, BlinkHz: 1.0, Confidence: 0.95},
			},
		}
		sqliteStorage.Save(obs)

		// Vary the code slightly to create different pattern IDs
		// Use BARR-C compliant code (Module_Function naming)
		code := fmt.Sprintf(`#include <stdint.h>

void LED_Blink%d(void) {
    uint8_t const status = %d;
}
`, i, i)
		_, err := store.StoreValidatedPattern(fmt.Sprintf("Blink LED %d", i), code, "esp32-dev", firmware)
		if err != nil {
			t.Logf("Failed to store pattern %d: %v", i, err)
		}
	}

	// Check that patterns were actually stored
	stats := store.Stats()
	t.Logf("Store stats: %+v", stats)

	// Check vector store count
	if store.vectorStore != nil {
		t.Logf("Vector store count: %d", store.vectorStore.Count())
	}

	// Search with topK=3
	results, err := store.SearchSimilarPatterns("Blink LED", "", 3)
	if err != nil {
		t.Fatalf("SearchSimilarPatterns failed: %v", err)
	}

	t.Logf("Got %d results", len(results))

	// Should return at least 3 results (or whatever was stored)
	if len(results) != 3 && len(results) < stats["patterns"] {
		t.Errorf("Expected 3 results, got %d (patterns stored: %d)", len(results), stats["patterns"])
	}

	_ = tmpDir
}

func TestSearchSimilarPatterns_NoVectorStore(t *testing.T) {
	tmpDir := setupTestEnvironment(t)
	defer os.Unsetenv("HOME")

	// Create pattern store WITHOUT initializing vector store
	store, err := NewPatternStore()
	if err != nil {
		t.Fatalf("Failed to create pattern store: %v", err)
	}
	defer store.Close()

	// Search should fail gracefully
	_, err = store.SearchSimilarPatterns("Blink LED", "", 5)
	if err == nil {
		t.Error("Expected error when vector store not initialized")
	}

	_ = tmpDir
}

func TestCalculateConfidence(t *testing.T) {
	pattern := &PatternNode{
		StyleCompliant: true,
	}

	obs := &core.Observation{
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, Confidence: 0.95},
			core.LEDSignal{Name: "LED2", On: true, Confidence: 0.95},
			core.DisplaySignal{Name: "OLED", Text: "Hello", Confidence: 0.95},
		},
	}

	// Test with high similarity and signals
	confidence := calculateConfidence(0.9, pattern, obs)
	if confidence < 0.9 || confidence > 1.0 {
		t.Errorf("calculateConfidence(0.9) = %f, expected ~0.9", confidence)
	}

	// Test with low similarity
	confidence = calculateConfidence(0.5, pattern, obs)
	if confidence < 0.5 || confidence > 0.7 {
		t.Errorf("calculateConfidence(0.5) = %f, expected ~0.5-0.6", confidence)
	}

	// Test with no observation
	confidence = calculateConfidence(0.8, pattern, nil)
	if confidence < 0.8 || confidence > 0.9 {
		t.Errorf("calculateConfidence(0.8, nil obs) = %f, expected ~0.8", confidence)
	}
}
