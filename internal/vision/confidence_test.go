package vision

import (
	"testing"

	"github.com/perceptumx/percepta/internal/core"
)

func TestConfidenceCalibrator_CalibrateLED_HighDetectionRate(t *testing.T) {
	cal := NewConfidenceCalibrator()

	led := core.LEDSignal{
		Name:       "LED1",
		On:         true,
		Color:      core.RGB{R: 255, G: 0, B: 0},
		Confidence: 0.8,
	}

	// Detected in all frames (100% detection rate)
	calibrated := cal.CalibrateLED(led, 1.0)

	// Base: 0.8
	// Agreement boost: (1.0 - 0.5) * 0.2 = 0.1
	// Color boost: 0.05
	// Blink boost: 0 (not blinking, so 0.05)
	// Total: 0.8 + 0.1 + 0.05 + 0.05 = 1.0
	if calibrated.Confidence != 1.0 {
		t.Errorf("expected confidence 1.0, got %f", calibrated.Confidence)
	}
}

func TestConfidenceCalibrator_CalibrateLED_LowDetectionRate(t *testing.T) {
	cal := NewConfidenceCalibrator()

	led := core.LEDSignal{
		Name:       "LED2",
		On:         true,
		Confidence: 0.7,
	}

	// Detected in 60% of frames
	calibrated := cal.CalibrateLED(led, 0.6)

	// Base: 0.7
	// Agreement boost: (0.6 - 0.5) * 0.2 = 0.02
	// Color boost: 0 (no color)
	// Blink boost: 0.05 (not blinking)
	// Total: 0.7 + 0.02 + 0 + 0.05 = 0.77
	if calibrated.Confidence != 0.77 {
		t.Errorf("expected confidence 0.77, got %f", calibrated.Confidence)
	}
}

func TestConfidenceCalibrator_CalibrateLED_VeryLowDetectionRate(t *testing.T) {
	cal := NewConfidenceCalibrator()

	led := core.LEDSignal{
		Name:       "LED3",
		On:         true,
		Confidence: 0.6,
	}

	// Detected in only 30% of frames (below threshold)
	calibrated := cal.CalibrateLED(led, 0.3)

	// Base: 0.6
	// Agreement boost: 0 (0.3 - 0.5 is negative, clamped to 0)
	// Color boost: 0
	// Blink boost: 0.05
	// Total: 0.6 + 0 + 0 + 0.05 = 0.65
	if calibrated.Confidence != 0.65 {
		t.Errorf("expected confidence 0.65, got %f", calibrated.Confidence)
	}
}

func TestConfidenceCalibrator_CalibrateDisplay_TypicalText(t *testing.T) {
	cal := NewConfidenceCalibrator()

	display := core.DisplaySignal{
		Name:       "LCD",
		Text:       "Hello World",
		Confidence: 0.85,
	}

	calibrated := cal.CalibrateDisplay(display)

	// Base: 0.85
	// Length factor: 0.05 (11 chars, in 5-50 range)
	// Special penalty: 0 (only space, no excessive special chars)
	// Total: 0.85 + 0.05 = 0.90
	if calibrated.Confidence != 0.90 {
		t.Errorf("expected confidence 0.90, got %f", calibrated.Confidence)
	}
}

func TestConfidenceCalibrator_CalibrateDisplay_ShortText(t *testing.T) {
	cal := NewConfidenceCalibrator()

	display := core.DisplaySignal{
		Name:       "OLED",
		Text:       "OK",
		Confidence: 0.80,
	}

	calibrated := cal.CalibrateDisplay(display)

	// Base: 0.80
	// Length factor: -0.1 (2 chars, < 5)
	// Special penalty: 0
	// Total: 0.80 - 0.1 = 0.70
	expected := 0.70
	if calibrated.Confidence < expected-0.001 || calibrated.Confidence > expected+0.001 {
		t.Errorf("expected confidence ~%f, got %f", expected, calibrated.Confidence)
	}
}

func TestConfidenceCalibrator_CalibrateDisplay_LongText(t *testing.T) {
	cal := NewConfidenceCalibrator()

	display := core.DisplaySignal{
		Name:       "LCD",
		Text:       "This is a very long text that appears on the display with more than 50 characters total",
		Confidence: 0.75,
	}

	calibrated := cal.CalibrateDisplay(display)

	// Base: 0.75
	// Length factor: 0.1 (> 50 chars)
	// Special penalty: 0
	// Total: 0.75 + 0.1 = 0.85
	if calibrated.Confidence != 0.85 {
		t.Errorf("expected confidence 0.85, got %f", calibrated.Confidence)
	}
}

func TestConfidenceCalibrator_CalibrateDisplay_NoisyText(t *testing.T) {
	cal := NewConfidenceCalibrator()

	display := core.DisplaySignal{
		Name:       "LCD",
		Text:       "A@#$%^&*()",
		Confidence: 0.90,
	}

	calibrated := cal.CalibrateDisplay(display)

	// Base: 0.90
	// Length factor: 0.05 (10 chars, in 5-50 range)
	// Special penalty: -0.15 (9/10 = 90% special chars, > 30%)
	// Total: 0.90 + 0.05 - 0.15 = 0.80
	if calibrated.Confidence != 0.80 {
		t.Errorf("expected confidence 0.80, got %f", calibrated.Confidence)
	}
}

func TestConfidenceCalibrator_CalibrateDisplay_Floor(t *testing.T) {
	cal := NewConfidenceCalibrator()

	display := core.DisplaySignal{
		Name:       "LCD",
		Text:       "@#$",
		Confidence: 0.50,
	}

	calibrated := cal.CalibrateDisplay(display)

	// Base: 0.50
	// Length factor: -0.1 (< 5 chars)
	// Special penalty: -0.15 (100% special chars)
	// Total: 0.50 - 0.1 - 0.15 = 0.25, floored to 0.50
	if calibrated.Confidence != 0.50 {
		t.Errorf("expected confidence floored to 0.50, got %f", calibrated.Confidence)
	}
}
