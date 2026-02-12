package vision

import (
	"github.com/perceptumx/percepta/internal/core"
)

// ConfidenceCalibrator adjusts confidence scores based on signal quality metrics
type ConfidenceCalibrator struct{}

func NewConfidenceCalibrator() *ConfidenceCalibrator {
	return &ConfidenceCalibrator{}
}

// CalibrateLED adjusts LED confidence based on:
// - Multi-frame agreement (if detected in all frames → higher confidence)
// - State stability (steady state → higher confidence than flickering)
// - Color detection (if color present → higher confidence)
func (c *ConfidenceCalibrator) CalibrateLED(led core.LEDSignal, detectionRate float64) core.LEDSignal {
	baseConf := led.Confidence

	// Multi-frame agreement boost
	// detectionRate = fraction of frames where LED was detected
	// 1.0 = detected in all frames → +0.1 boost
	// 0.5 = detected in half → no boost
	agreementBoost := (detectionRate - 0.5) * 0.2
	if agreementBoost < 0 {
		agreementBoost = 0
	}

	// Color detection boost
	colorBoost := 0.0
	if led.Color != (core.RGB{}) {
		colorBoost = 0.05 // Color detected → +0.05
	}

	// Blink frequency confidence
	// Steady state (0 Hz or no blink) → +0.05
	// Measured blink → confidence in frequency
	blinkBoost := 0.0
	if led.BlinkHz == 0 {
		blinkBoost = 0.05 // Steady state
	}

	// Total confidence (capped at 1.0)
	totalConf := baseConf + agreementBoost + colorBoost + blinkBoost
	if totalConf > 1.0 {
		totalConf = 1.0
	}

	led.Confidence = totalConf
	return led
}

// CalibrateDisplay adjusts display confidence based on:
// - Text length (longer text → higher confidence in OCR)
// - Special characters (if present → might be OCR noise, lower confidence)
// - Base confidence from StructuredParser (tool use confidence)
func (c *ConfidenceCalibrator) CalibrateDisplay(display core.DisplaySignal) core.DisplaySignal {
	baseConf := display.Confidence

	// Text length factor
	// Short text (< 5 chars) might be noise
	// Medium text (5-50 chars) is typical
	// Long text (> 50 chars) is confident OCR
	lengthFactor := 0.0
	textLen := len(display.Text)
	if textLen < 5 {
		lengthFactor = -0.1 // Penalize very short text
	} else if textLen >= 5 && textLen <= 50 {
		lengthFactor = 0.05 // Typical display text
	} else {
		lengthFactor = 0.1 // Long text → confident OCR
	}

	// Special character penalty
	// Excessive special chars might indicate OCR noise
	specialChars := 0
	for _, ch := range display.Text {
		if !isAlphanumeric(ch) && ch != ' ' && ch != '.' && ch != ':' {
			specialChars++
		}
	}
	specialPenalty := 0.0
	if textLen > 0 && float64(specialChars)/float64(textLen) > 0.3 {
		specialPenalty = -0.15 // >30% special chars → likely noise
	}

	// Total confidence (capped at 1.0, floored at 0.5)
	totalConf := baseConf + lengthFactor + specialPenalty
	if totalConf > 1.0 {
		totalConf = 1.0
	}
	if totalConf < 0.5 {
		totalConf = 0.5 // Don't go below 0.5 for valid displays
	}

	display.Confidence = totalConf
	return display
}

func isAlphanumeric(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')
}
