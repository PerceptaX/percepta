package diff

import (
	"testing"
	"time"

	"github.com/perceptumx/percepta/internal/core"
)

func TestCompare_NoChanges(t *testing.T) {
	from := &core.Observation{
		ID:           "obs1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, Color: core.RGB{R: 255, G: 0, B: 0}},
		},
	}

	to := &core.Observation{
		ID:           "obs2",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, Color: core.RGB{R: 255, G: 0, B: 0}},
		},
	}

	result := Compare(from, to)

	if result.HasChanges() {
		t.Errorf("Expected no changes, but got %d changes", len(result.Changes))
	}
}

func TestCompare_LEDAdded(t *testing.T) {
	from := &core.Observation{
		ID:           "obs1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true},
		},
	}

	to := &core.Observation{
		ID:           "obs2",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true},
			core.LEDSignal{Name: "LED2", On: true, Color: core.RGB{R: 128, G: 0, B: 128}, BlinkHz: 0.8},
		},
	}

	result := Compare(from, to)

	if !result.HasChanges() {
		t.Fatal("Expected changes, but got none")
	}

	if len(result.Changes) != 1 {
		t.Fatalf("Expected 1 change, got %d", len(result.Changes))
	}

	change := result.Changes[0]
	if change.Type != ChangeAdded {
		t.Errorf("Expected ChangeAdded, got %v", change.Type)
	}
	if change.Name != "LED2" {
		t.Errorf("Expected LED2, got %s", change.Name)
	}
}

func TestCompare_LEDRemoved(t *testing.T) {
	from := &core.Observation{
		ID:           "obs1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true},
			core.LEDSignal{Name: "LED2", On: false},
		},
	}

	to := &core.Observation{
		ID:           "obs2",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true},
		},
	}

	result := Compare(from, to)

	if !result.HasChanges() {
		t.Fatal("Expected changes, but got none")
	}

	if len(result.Changes) != 1 {
		t.Fatalf("Expected 1 change, got %d", len(result.Changes))
	}

	change := result.Changes[0]
	if change.Type != ChangeRemoved {
		t.Errorf("Expected ChangeRemoved, got %v", change.Type)
	}
	if change.Name != "LED2" {
		t.Errorf("Expected LED2, got %s", change.Name)
	}
}

func TestCompare_LEDStateChange(t *testing.T) {
	from := &core.Observation{
		ID:           "obs1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true},
		},
	}

	to := &core.Observation{
		ID:           "obs2",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: false},
		},
	}

	result := Compare(from, to)

	if !result.HasChanges() {
		t.Fatal("Expected changes, but got none")
	}

	change := result.Changes[0]
	if change.Type != ChangeModified {
		t.Errorf("Expected ChangeModified, got %v", change.Type)
	}
}

func TestCompare_LEDColorChange(t *testing.T) {
	from := &core.Observation{
		ID:           "obs1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, Color: core.RGB{R: 255, G: 0, B: 0}},
		},
	}

	to := &core.Observation{
		ID:           "obs2",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, Color: core.RGB{R: 0, G: 0, B: 255}},
		},
	}

	result := Compare(from, to)

	if !result.HasChanges() {
		t.Fatal("Expected changes, but got none")
	}

	change := result.Changes[0]
	if change.Type != ChangeModified {
		t.Errorf("Expected ChangeModified, got %v", change.Type)
	}
}

func TestCompare_LEDBlinkRateChange(t *testing.T) {
	from := &core.Observation{
		ID:           "obs1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, BlinkHz: 2.0},
		},
	}

	to := &core.Observation{
		ID:           "obs2",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, BlinkHz: 2.5},
		},
	}

	result := Compare(from, to)

	if !result.HasChanges() {
		t.Fatal("Expected changes, but got none")
	}

	change := result.Changes[0]
	if change.Type != ChangeModified {
		t.Errorf("Expected ChangeModified, got %v", change.Type)
	}
}

func TestCompare_BlinkHzNormalization(t *testing.T) {
	// Test that slight variations in BlinkHz are normalized (rounded to 1 decimal)
	from := &core.Observation{
		ID:           "obs1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, BlinkHz: 2.04}, // Should normalize to 2.0
		},
	}

	to := &core.Observation{
		ID:           "obs2",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, BlinkHz: 2.01}, // Should normalize to 2.0
		},
	}

	result := Compare(from, to)

	if result.HasChanges() {
		t.Errorf("Expected no changes due to normalization, but got %d changes", len(result.Changes))
	}
}

func TestCompare_ExactColorComparison(t *testing.T) {
	// Test that colors are compared exactly (no tolerance)
	from := &core.Observation{
		ID:           "obs1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, Color: core.RGB{R: 255, G: 0, B: 0}},
		},
	}

	to := &core.Observation{
		ID:           "obs2",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, Color: core.RGB{R: 254, G: 0, B: 0}}, // Off by 1
		},
	}

	result := Compare(from, to)

	if !result.HasChanges() {
		t.Fatal("Expected changes for exact color comparison, but got none")
	}

	change := result.Changes[0]
	if change.Type != ChangeModified {
		t.Errorf("Expected ChangeModified, got %v", change.Type)
	}
}

func TestCompare_ConfidenceIgnored(t *testing.T) {
	// Test that confidence values are ignored in comparison
	from := &core.Observation{
		ID:           "obs1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, Confidence: 0.95},
		},
	}

	to := &core.Observation{
		ID:           "obs2",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, Confidence: 0.85}, // Different confidence
		},
	}

	result := Compare(from, to)

	if result.HasChanges() {
		t.Errorf("Expected no changes (confidence should be ignored), but got %d changes", len(result.Changes))
	}
}

func TestCompare_DisplayTextChange(t *testing.T) {
	from := &core.Observation{
		ID:           "obs1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.DisplaySignal{Name: "LCD1", Text: "Hello"},
		},
	}

	to := &core.Observation{
		ID:           "obs2",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.DisplaySignal{Name: "LCD1", Text: "World"},
		},
	}

	result := Compare(from, to)

	if !result.HasChanges() {
		t.Fatal("Expected changes, but got none")
	}

	change := result.Changes[0]
	if change.Type != ChangeModified {
		t.Errorf("Expected ChangeModified, got %v", change.Type)
	}
}

func TestCompare_CountByType(t *testing.T) {
	from := &core.Observation{
		ID:           "obs1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true},
			core.LEDSignal{Name: "LED2", On: false},
			core.LEDSignal{Name: "LED3", On: true},
		},
	}

	to := &core.Observation{
		ID:           "obs2",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: false}, // Modified
			core.LEDSignal{Name: "LED2", On: false}, // Unchanged
			core.LEDSignal{Name: "LED4", On: true},  // Added (LED3 removed)
		},
	}

	result := Compare(from, to)

	added, removed, modified := result.CountByType()

	if added != 1 {
		t.Errorf("Expected 1 added, got %d", added)
	}
	if removed != 1 {
		t.Errorf("Expected 1 removed, got %d", removed)
	}
	if modified != 1 {
		t.Errorf("Expected 1 modified, got %d", modified)
	}
}

func TestCompare_DisplayHistoryChange(t *testing.T) {
	from := &core.Observation{
		ID:           "obs1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.DisplaySignal{
				Name:    "LCD",
				Text:    "Ready",
				Changed: true,
				History: []core.DisplayTextEntry{
					{OffsetMs: 0, Text: "Boot"},
					{OffsetMs: 400, Text: "Ready"},
				},
			},
		},
	}

	to := &core.Observation{
		ID:           "obs2",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.DisplaySignal{
				Name:    "LCD",
				Text:    "Done",
				Changed: true,
				History: []core.DisplayTextEntry{
					{OffsetMs: 0, Text: "Boot"},
					{OffsetMs: 400, Text: "Init"},
					{OffsetMs: 800, Text: "Done"},
				},
			},
		},
	}

	result := Compare(from, to)

	if !result.HasChanges() {
		t.Fatal("Expected changes for different display histories")
	}

	change := result.Changes[0]
	if change.Type != ChangeModified {
		t.Errorf("Expected ChangeModified, got %v", change.Type)
	}
}

func TestCompare_DisplayHistoryNoChange(t *testing.T) {
	from := &core.Observation{
		ID:           "obs1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.DisplaySignal{
				Name:    "LCD",
				Text:    "Ready",
				Changed: true,
				History: []core.DisplayTextEntry{
					{OffsetMs: 0, Text: "Boot"},
					{OffsetMs: 400, Text: "Ready"},
				},
			},
		},
	}

	to := &core.Observation{
		ID:           "obs2",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.DisplaySignal{
				Name:    "LCD",
				Text:    "Ready",
				Changed: true,
				History: []core.DisplayTextEntry{
					{OffsetMs: 0, Text: "Boot"},
					{OffsetMs: 500, Text: "Ready"}, // Different timing, same text sequence
				},
			},
		},
	}

	result := Compare(from, to)

	if result.HasChanges() {
		t.Errorf("Expected no changes for same text sequence with different timing, got %d", len(result.Changes))
	}
}

func TestCompare_StaticToChangingDisplay(t *testing.T) {
	from := &core.Observation{
		ID:           "obs1",
		DeviceID:     "fpga",
		FirmwareHash: "v1",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.DisplaySignal{Name: "LCD", Text: "Ready"},
		},
	}

	to := &core.Observation{
		ID:           "obs2",
		DeviceID:     "fpga",
		FirmwareHash: "v2",
		Timestamp:    time.Now(),
		Signals: []core.Signal{
			core.DisplaySignal{
				Name:    "LCD",
				Text:    "Done",
				Changed: true,
				History: []core.DisplayTextEntry{
					{OffsetMs: 0, Text: "Boot"},
					{OffsetMs: 400, Text: "Done"},
				},
			},
		},
	}

	result := Compare(from, to)

	if !result.HasChanges() {
		t.Fatal("Expected changes for static→changing transition")
	}

	change := result.Changes[0]
	if change.Type != ChangeModified {
		t.Errorf("Expected ChangeModified, got %v", change.Type)
	}
}

func TestNormalizeBlinkHz(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0.0, 0.0},
		{2.0, 2.0},
		{2.04, 2.0},
		{2.05, 2.1},
		{2.14, 2.1},
		{2.15, 2.2},
		{0.84, 0.8},
		{0.85, 0.9},
	}

	for _, tt := range tests {
		result := normalizeBlinkHz(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeBlinkHz(%v) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}
