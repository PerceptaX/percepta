package filter

import (
	"testing"
	"time"

	"github.com/perceptumx/percepta/internal/core"
	"github.com/perceptumx/percepta/internal/storage"
)

func TestTemporalSmoother_NoHistory(t *testing.T) {
	// Setup: empty storage
	store := storage.NewMemoryStorage()
	smoother := NewTemporalSmoother(store)

	// Create new observation
	obs := &core.Observation{
		ID:        "test-1",
		DeviceID:  "test-device",
		Timestamp: time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, Confidence: 0.95},
		},
	}

	// Smooth (should return unchanged - no history)
	smoothed, err := smoother.Smooth(obs)
	if err != nil {
		t.Fatalf("Smooth failed: %v", err)
	}

	// Verify unchanged
	if len(smoothed.Signals) != 1 {
		t.Errorf("Expected 1 signal, got %d", len(smoothed.Signals))
	}

	led := smoothed.Signals[0].(core.LEDSignal)
	if led.Name != "LED1" || !led.On {
		t.Errorf("LED state changed unexpectedly: %+v", led)
	}
}

func TestTemporalSmoother_StableSignal(t *testing.T) {
	// Setup: storage with stable LED history
	store := storage.NewMemoryStorage()
	smoother := NewTemporalSmoother(store)

	// Add 3 historical observations with stable LED state
	baseTime := time.Now().Add(-3 * time.Second)
	for i := 0; i < 3; i++ {
		obs := core.Observation{
			ID:        core.GenerateID(),
			DeviceID:  "test-device",
			Timestamp: baseTime.Add(time.Duration(i) * time.Second),
			Signals: []core.Signal{
				core.LEDSignal{Name: "LED1", On: true, BlinkHz: 2.0, Confidence: 0.95},
			},
		}
		store.Save(obs)
	}

	// New observation with same stable state
	newObs := &core.Observation{
		ID:        "new-obs",
		DeviceID:  "test-device",
		Timestamp: time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, BlinkHz: 2.0, Confidence: 0.95},
		},
	}

	// Smooth (should preserve stable state)
	smoothed, err := smoother.Smooth(newObs)
	if err != nil {
		t.Fatalf("Smooth failed: %v", err)
	}

	if len(smoothed.Signals) != 1 {
		t.Fatalf("Expected 1 signal, got %d", len(smoothed.Signals))
	}

	led := smoothed.Signals[0].(core.LEDSignal)
	if led.Name != "LED1" || !led.On || led.BlinkHz != 2.0 {
		t.Errorf("Stable LED state not preserved: %+v", led)
	}
}

func TestTemporalSmoother_SingleFrameGlitch(t *testing.T) {
	// Setup: storage with stable LED history
	store := storage.NewMemoryStorage()
	smoother := NewTemporalSmoother(store)

	// Add 3 historical observations with LED ON
	baseTime := time.Now().Add(-3 * time.Second)
	for i := 0; i < 3; i++ {
		obs := core.Observation{
			ID:        core.GenerateID(),
			DeviceID:  "test-device",
			Timestamp: baseTime.Add(time.Duration(i) * time.Second),
			Signals: []core.Signal{
				core.LEDSignal{Name: "LED1", On: true, BlinkHz: 2.0, Confidence: 0.95},
			},
		}
		store.Save(obs)
	}

	// New observation with glitch (LED OFF - anomaly)
	newObs := &core.Observation{
		ID:        "glitch-obs",
		DeviceID:  "test-device",
		Timestamp: time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: false, BlinkHz: 0.0, Confidence: 0.60}, // Glitch!
		},
	}

	// Smooth (should filter out glitch, use historical state)
	smoothed, err := smoother.Smooth(newObs)
	if err != nil {
		t.Fatalf("Smooth failed: %v", err)
	}

	if len(smoothed.Signals) != 1 {
		t.Fatalf("Expected 1 signal, got %d", len(smoothed.Signals))
	}

	led := smoothed.Signals[0].(core.LEDSignal)
	if !led.On {
		t.Errorf("Glitch not filtered: LED should be ON (historical state), but got OFF")
	}
	if led.BlinkHz != 2.0 {
		t.Errorf("Expected BlinkHz 2.0 from history, got %f", led.BlinkHz)
	}
}

func TestTemporalSmoother_DisplayOCRError(t *testing.T) {
	// Setup: storage with stable display text
	store := storage.NewMemoryStorage()
	smoother := NewTemporalSmoother(store)

	// Add 3 historical observations with stable display text
	baseTime := time.Now().Add(-3 * time.Second)
	for i := 0; i < 3; i++ {
		obs := core.Observation{
			ID:        core.GenerateID(),
			DeviceID:  "test-device",
			Timestamp: baseTime.Add(time.Duration(i) * time.Second),
			Signals: []core.Signal{
				core.DisplaySignal{Name: "DISPLAY1", Text: "READY", Confidence: 0.95},
			},
		}
		store.Save(obs)
	}

	// New observation with OCR error (garbled text)
	newObs := &core.Observation{
		ID:        "ocr-error-obs",
		DeviceID:  "test-device",
		Timestamp: time.Now(),
		Signals: []core.Signal{
			core.DisplaySignal{Name: "DISPLAY1", Text: "R3ADY", Confidence: 0.50}, // OCR error!
		},
	}

	// Smooth (should correct OCR error with historical text)
	smoothed, err := smoother.Smooth(newObs)
	if err != nil {
		t.Fatalf("Smooth failed: %v", err)
	}

	if len(smoothed.Signals) != 1 {
		t.Fatalf("Expected 1 signal, got %d", len(smoothed.Signals))
	}

	display := smoothed.Signals[0].(core.DisplaySignal)
	if display.Text != "READY" {
		t.Errorf("OCR error not corrected: expected 'READY', got '%s'", display.Text)
	}
}

func TestTemporalSmoother_GracefulDegradation(t *testing.T) {
	// Setup: mock storage that fails Query
	store := &failingStorage{}
	smoother := NewTemporalSmoother(store)

	// Create new observation
	obs := &core.Observation{
		ID:        "test-1",
		DeviceID:  "test-device",
		Timestamp: time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, Confidence: 0.95},
		},
	}

	// Smooth (should gracefully degrade - return unsmoothed)
	smoothed, err := smoother.Smooth(obs)
	if err != nil {
		t.Fatalf("Smooth should not return error on storage failure: %v", err)
	}

	// Verify observation returned unchanged
	if len(smoothed.Signals) != 1 {
		t.Errorf("Expected 1 signal, got %d", len(smoothed.Signals))
	}

	led := smoothed.Signals[0].(core.LEDSignal)
	if led.Name != "LED1" || !led.On {
		t.Errorf("LED state changed unexpectedly: %+v", led)
	}
}

func TestTemporalSmoother_BlinkFrequencyTolerance(t *testing.T) {
	// Setup: storage with LED at 2.0 Hz
	store := storage.NewMemoryStorage()
	smoother := NewTemporalSmoother(store)

	// Add historical observations with 2.0 Hz blink
	baseTime := time.Now().Add(-3 * time.Second)
	for i := 0; i < 3; i++ {
		obs := core.Observation{
			ID:        core.GenerateID(),
			DeviceID:  "test-device",
			Timestamp: baseTime.Add(time.Duration(i) * time.Second),
			Signals: []core.Signal{
				core.LEDSignal{Name: "LED1", On: true, BlinkHz: 2.0, Confidence: 0.95},
			},
		}
		store.Save(obs)
	}

	// Test within 10% tolerance (2.0 * 0.1 = 0.2, so 1.8-2.2 should match)
	newObs := &core.Observation{
		ID:        "test-tolerance",
		DeviceID:  "test-device",
		Timestamp: time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, BlinkHz: 2.15, Confidence: 0.95}, // 7.5% deviation - should match
		},
	}

	smoothed, err := smoother.Smooth(newObs)
	if err != nil {
		t.Fatalf("Smooth failed: %v", err)
	}

	led := smoothed.Signals[0].(core.LEDSignal)
	if led.BlinkHz != 2.15 {
		t.Errorf("Expected new observation preserved (within tolerance), got BlinkHz=%f", led.BlinkHz)
	}
}

func TestTemporalSmoother_TimeWindow(t *testing.T) {
	// Setup: storage with old observations (outside time window)
	store := storage.NewMemoryStorage()
	smoother := NewTemporalSmoother(store)

	// Add old observation (10 seconds ago, outside 5-second window)
	oldObs := core.Observation{
		ID:        core.GenerateID(),
		DeviceID:  "test-device",
		Timestamp: time.Now().Add(-10 * time.Second), // Outside window
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: false, BlinkHz: 0.0, Confidence: 0.95},
		},
	}
	store.Save(oldObs)

	// New observation with different state
	newObs := &core.Observation{
		ID:        "new-obs",
		DeviceID:  "test-device",
		Timestamp: time.Now(),
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, BlinkHz: 2.0, Confidence: 0.95},
		},
	}

	// Smooth (old observation outside window, should trust new observation)
	smoothed, err := smoother.Smooth(newObs)
	if err != nil {
		t.Fatalf("Smooth failed: %v", err)
	}

	led := smoothed.Signals[0].(core.LEDSignal)
	if !led.On || led.BlinkHz != 2.0 {
		t.Errorf("Expected new observation preserved (no recent history), got %+v", led)
	}
}

func TestTemporalSmoother_ChangingDisplayBypassesSmoothing(t *testing.T) {
	// Setup: storage with stable display history
	store := storage.NewMemoryStorage()
	smoother := NewTemporalSmoother(store)

	// Add historical observations with stable "READY" text
	baseTime := time.Now().Add(-3 * time.Second)
	for i := 0; i < 3; i++ {
		obs := core.Observation{
			ID:        core.GenerateID(),
			DeviceID:  "test-device",
			Timestamp: baseTime.Add(time.Duration(i) * time.Second),
			Signals: []core.Signal{
				core.DisplaySignal{Name: "LCD", Text: "READY", Confidence: 0.95},
			},
		}
		store.Save(obs)
	}

	// New observation with Changed=true (display is transitioning)
	newObs := &core.Observation{
		ID:        "changing-obs",
		DeviceID:  "test-device",
		Timestamp: time.Now(),
		Signals: []core.Signal{
			core.DisplaySignal{
				Name:       "LCD",
				Text:       "Rebooting",
				Confidence: 0.85,
				Changed:    true,
				History: []core.DisplayTextEntry{
					{OffsetMs: 0, Text: "READY", Confidence: 0.90},
					{OffsetMs: 400, Text: "Rebooting", Confidence: 0.85},
				},
			},
		},
	}

	// Smooth should NOT override changing display with historical "READY"
	smoothed, err := smoother.Smooth(newObs)
	if err != nil {
		t.Fatalf("Smooth failed: %v", err)
	}

	if len(smoothed.Signals) != 1 {
		t.Fatalf("Expected 1 signal, got %d", len(smoothed.Signals))
	}

	display := smoothed.Signals[0].(core.DisplaySignal)
	if !display.Changed {
		t.Errorf("Expected Changed=true to be preserved")
	}
	if display.Text != "Rebooting" {
		t.Errorf("Expected text 'Rebooting' (not smoothed away), got '%s'", display.Text)
	}
	if len(display.History) != 2 {
		t.Errorf("Expected history preserved, got %d entries", len(display.History))
	}
}

// failingStorage is a mock storage that fails Query
type failingStorage struct{}

func (f *failingStorage) Save(obs core.Observation) error {
	return nil
}

func (f *failingStorage) Query(deviceID string, limit int) ([]core.Observation, error) {
	return nil, core.ErrStorageUnavailable
}

func (f *failingStorage) Count() int {
	return 0
}
