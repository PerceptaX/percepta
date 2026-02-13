package filter

import (
	"time"

	"github.com/perceptumx/percepta/internal/core"
)

// TemporalSmoother filters out single-frame anomalies by comparing consecutive observations
type TemporalSmoother struct {
	window       time.Duration // Time window for smoothing (e.g., 5 seconds)
	minAgreement int           // Minimum observations that must agree (e.g., 2 out of 3)
	storage      core.StorageDriver
}

func NewTemporalSmoother(storage core.StorageDriver) *TemporalSmoother {
	return &TemporalSmoother{
		window:       5 * time.Second,
		minAgreement: 2, // Require 2/3 agreement
		storage:      storage,
	}
}

// Smooth filters the new observation against recent history
func (t *TemporalSmoother) Smooth(newObs *core.Observation) (*core.Observation, error) {
	// Get recent observations for same device
	recent, err := t.storage.Query(newObs.DeviceID, 10) // Last 10 observations
	if err != nil {
		return newObs, nil // If query fails, return unsmoothed (graceful degradation)
	}

	// Filter to observations within time window
	cutoff := time.Now().Add(-t.window)
	var windowObs []core.Observation
	for _, obs := range recent {
		if obs.Timestamp.After(cutoff) {
			windowObs = append(windowObs, obs)
		}
	}

	// If no recent observations, return as-is (nothing to smooth against)
	if len(windowObs) == 0 {
		return newObs, nil
	}

	// Smooth LED signals
	smoothedLEDs := t.smoothLEDs(newObs, windowObs)

	// Smooth display signals
	smoothedDisplays := t.smoothDisplays(newObs, windowObs)

	// Combine smoothed signals
	var smoothedSignals []core.Signal
	smoothedSignals = append(smoothedSignals, smoothedLEDs...)
	smoothedSignals = append(smoothedSignals, smoothedDisplays...)

	return &core.Observation{
		ID:           newObs.ID,
		DeviceID:     newObs.DeviceID,
		FirmwareHash: newObs.FirmwareHash,
		Timestamp:    newObs.Timestamp,
		Signals:      smoothedSignals,
	}, nil
}

func (t *TemporalSmoother) smoothLEDs(newObs *core.Observation, history []core.Observation) []core.Signal {
	var smoothed []core.Signal

	// Extract LED signals from new observation
	newLEDs := extractLEDs(newObs.Signals)

	for _, newLED := range newLEDs {
		// Find this LED in recent history
		historicalStates := t.findLEDHistory(newLED.Name, history)

		// Check if new state agrees with historical trend
		if len(historicalStates) < t.minAgreement {
			// Not enough history → trust new observation
			smoothed = append(smoothed, newLED)
			continue
		}

		// Count agreement with historical states
		agreementCount := 0
		for _, histLED := range historicalStates {
			if t.ledsMatch(newLED, histLED) {
				agreementCount++
			}
		}

		// If majority agrees, keep new observation
		// If majority disagrees, it's likely noise → use historical state
		if agreementCount >= t.minAgreement {
			smoothed = append(smoothed, newLED)
		} else {
			// Use most recent historical state (likely more stable)
			if len(historicalStates) > 0 {
				smoothed = append(smoothed, historicalStates[0])
			}
		}
	}

	return smoothed
}

func (t *TemporalSmoother) smoothDisplays(newObs *core.Observation, history []core.Observation) []core.Signal {
	var smoothed []core.Signal

	// Extract display signals from new observation
	newDisplays := extractDisplays(newObs.Signals)

	for _, newDisplay := range newDisplays {
		// Skip smoothing for displays with detected state changes
		if newDisplay.Changed {
			smoothed = append(smoothed, newDisplay)
			continue
		}

		// Find this display in recent history
		historicalTexts := t.findDisplayHistory(newDisplay.Name, history)

		// Check if new text agrees with historical trend
		if len(historicalTexts) < t.minAgreement {
			// Not enough history → trust new observation
			smoothed = append(smoothed, newDisplay)
			continue
		}

		// Count agreement (exact text match)
		agreementCount := 0
		for _, histText := range historicalTexts {
			if newDisplay.Text == histText {
				agreementCount++
			}
		}

		// If majority agrees, keep new observation
		// If majority disagrees, it's likely OCR glitch → use historical text
		if agreementCount >= t.minAgreement {
			smoothed = append(smoothed, newDisplay)
		} else {
			// Use most recent historical text
			if len(historicalTexts) > 0 {
				smoothed = append(smoothed, core.DisplaySignal{
					Name:       newDisplay.Name,
					Text:       historicalTexts[0],
					Confidence: newDisplay.Confidence,
				})
			}
		}
	}

	return smoothed
}

func (t *TemporalSmoother) findLEDHistory(name string, history []core.Observation) []core.LEDSignal {
	var leds []core.LEDSignal

	for _, obs := range history {
		for _, signal := range obs.Signals {
			if led, ok := signal.(core.LEDSignal); ok && led.Name == name {
				leds = append(leds, led)
			}
		}
	}

	return leds
}

func (t *TemporalSmoother) findDisplayHistory(name string, history []core.Observation) []string {
	var texts []string

	for _, obs := range history {
		for _, signal := range obs.Signals {
			if display, ok := signal.(core.DisplaySignal); ok && display.Name == name {
				texts = append(texts, display.Text)
			}
		}
	}

	return texts
}

func (t *TemporalSmoother) ledsMatch(a, b core.LEDSignal) bool {
	// LEDs match if on/off state and blink frequency are similar
	if a.On != b.On {
		return false
	}

	// Allow 10% tolerance on blink frequency
	if a.BlinkHz > 0 || b.BlinkHz > 0 {
		diff := abs(a.BlinkHz - b.BlinkHz)
		avg := (a.BlinkHz + b.BlinkHz) / 2.0
		if avg > 0 && diff/avg > 0.1 {
			return false
		}
	}

	return true
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func extractLEDs(signals []core.Signal) []core.LEDSignal {
	var leds []core.LEDSignal
	for _, signal := range signals {
		if led, ok := signal.(core.LEDSignal); ok {
			leds = append(leds, led)
		}
	}
	return leds
}

func extractDisplays(signals []core.Signal) []core.DisplaySignal {
	var displays []core.DisplaySignal
	for _, signal := range signals {
		if display, ok := signal.(core.DisplaySignal); ok {
			displays = append(displays, display)
		}
	}
	return displays
}
