package vision

import (
	"testing"
	"time"

	"github.com/perceptumx/percepta/internal/core"
)

func TestLEDAggregator_SteadyOn(t *testing.T) {
	agg := &ledAggregator{
		name: "LED1",
		observations: []core.LEDSignal{
			{Name: "LED1", On: true, Confidence: 0.9},
			{Name: "LED1", On: true, Confidence: 0.92},
			{Name: "LED1", On: true, Confidence: 0.88},
			{Name: "LED1", On: true, Confidence: 0.91},
			{Name: "LED1", On: true, Confidence: 0.89},
		},
	}

	result := agg.aggregate()

	if result.Name != "LED1" {
		t.Errorf("expected LED1, got %s", result.Name)
	}
	if !result.On {
		t.Errorf("expected LED to be on")
	}
	if result.BlinkHz != 0 {
		t.Errorf("expected BlinkHz 0, got %f", result.BlinkHz)
	}

	// Confidence should be average
	expectedConf := (0.9 + 0.92 + 0.88 + 0.91 + 0.89) / 5.0
	if result.Confidence != expectedConf {
		t.Errorf("expected confidence %f, got %f", expectedConf, result.Confidence)
	}
}

func TestLEDAggregator_SteadyOff(t *testing.T) {
	agg := &ledAggregator{
		name: "LED2",
		observations: []core.LEDSignal{
			{Name: "LED2", On: false, Confidence: 0.85},
			{Name: "LED2", On: false, Confidence: 0.87},
			{Name: "LED2", On: false, Confidence: 0.86},
		},
	}

	result := agg.aggregate()

	if result.On {
		t.Errorf("expected LED to be off")
	}
	if result.BlinkHz != 0 {
		t.Errorf("expected BlinkHz 0, got %f", result.BlinkHz)
	}
}

func TestLEDAggregator_Blinking(t *testing.T) {
	agg := &ledAggregator{
		name: "LED3",
		observations: []core.LEDSignal{
			{Name: "LED3", On: true, Confidence: 0.9},
			{Name: "LED3", On: false, Confidence: 0.88},
			{Name: "LED3", On: true, Confidence: 0.91},
			{Name: "LED3", On: false, Confidence: 0.89},
			{Name: "LED3", On: true, Confidence: 0.90},
		},
	}

	result := agg.aggregate()

	if !result.On {
		t.Errorf("expected blinking LED to be logically on")
	}

	// With 4 transitions (on->off, off->on, on->off, off->on), frequency = 4/2 = 2 Hz
	if result.BlinkHz != 2.0 {
		t.Errorf("expected BlinkHz 2.0, got %f", result.BlinkHz)
	}
}

func TestAggregateLEDs(t *testing.T) {
	frames := []FrameResult{
		{
			Signals: []core.Signal{
				core.LEDSignal{Name: "LED1", On: true, Confidence: 0.9},
				core.LEDSignal{Name: "LED2", On: false, Confidence: 0.85},
			},
			CapturedAt: time.Now(),
		},
		{
			Signals: []core.Signal{
				core.LEDSignal{Name: "LED1", On: true, Confidence: 0.92},
				core.LEDSignal{Name: "LED2", On: true, Confidence: 0.87},
			},
			CapturedAt: time.Now(),
		},
		{
			Signals: []core.Signal{
				core.LEDSignal{Name: "LED1", On: true, Confidence: 0.88},
				core.LEDSignal{Name: "LED2", On: false, Confidence: 0.86},
			},
			CapturedAt: time.Now(),
		},
	}

	leds := AggregateLEDs(frames)

	if len(leds) != 2 {
		t.Fatalf("expected 2 LEDs, got %d", len(leds))
	}

	// Find LED1 and LED2
	var led1, led2 *core.LEDSignal
	for i := range leds {
		if leds[i].Name == "LED1" {
			led1 = &leds[i]
		} else if leds[i].Name == "LED2" {
			led2 = &leds[i]
		}
	}

	if led1 == nil || led2 == nil {
		t.Fatalf("expected to find LED1 and LED2")
	}

	// LED1 should be steady on
	if !led1.On {
		t.Errorf("LED1 should be on")
	}
	if led1.BlinkHz != 0 {
		t.Errorf("LED1 should have BlinkHz 0, got %f", led1.BlinkHz)
	}

	// LED2 should be blinking (on in 1 frame, off in 2)
	if !led2.On {
		t.Errorf("LED2 should be logically on (blinking)")
	}
	// 2 transitions: off->on, on->off
	if led2.BlinkHz != 1.0 {
		t.Errorf("LED2 should have BlinkHz 1.0, got %f", led2.BlinkHz)
	}
}
