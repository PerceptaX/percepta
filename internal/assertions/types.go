package assertions

import (
	"fmt"
	"strings"

	"github.com/perceptumx/percepta/internal/core"
)

// AssertionResult represents the outcome of an assertion evaluation
type AssertionResult struct {
	Passed     bool
	Expected   string
	Actual     string
	Confidence float64
	Message    string
}

// Assertion interface - evaluates observed state
type Assertion interface {
	Evaluate(obs *core.Observation) AssertionResult
	String() string
}

// LEDAssertion validates LED state
type LEDAssertion struct {
	Name     string
	Expected LEDState
}

type LEDState struct {
	On      *bool // nil = don't care
	Color   *core.RGB
	BlinkHz *float64
}

func (a *LEDAssertion) Evaluate(obs *core.Observation) AssertionResult {
	// Find matching LED signal (case-insensitive)
	var matchedSignal *core.LEDSignal
	lowerName := strings.ToLower(a.Name)

	for _, sig := range obs.Signals {
		if ledSig, ok := sig.(*core.LEDSignal); ok {
			if strings.ToLower(ledSig.Name) == lowerName {
				matchedSignal = ledSig
				break
			}
		}
	}

	// Fallback: if no match and exactly one LED exists, use it
	if matchedSignal == nil {
		var ledSignals []*core.LEDSignal
		for _, sig := range obs.Signals {
			if ledSig, ok := sig.(*core.LEDSignal); ok {
				ledSignals = append(ledSignals, ledSig)
			}
		}
		if len(ledSignals) == 1 {
			matchedSignal = ledSignals[0]
		}
	}

	// If still no match, fail
	if matchedSignal == nil {
		return AssertionResult{
			Passed:     false,
			Expected:   a.String(),
			Actual:     "LED not found in observation",
			Confidence: 0.0,
			Message:    fmt.Sprintf("LED '%s' not found in observation", a.Name),
		}
	}

	// Check On/Off state if specified
	if a.Expected.On != nil {
		if matchedSignal.On != *a.Expected.On {
			return AssertionResult{
				Passed:     false,
				Expected:   a.String(),
				Actual:     fmt.Sprintf("LED '%s' is %s", matchedSignal.Name, onOffString(matchedSignal.On)),
				Confidence: matchedSignal.Confidence,
				Message:    fmt.Sprintf("Expected %s, but LED is %s", onOffString(*a.Expected.On), onOffString(matchedSignal.On)),
			}
		}
	}

	// Check color if specified
	if a.Expected.Color != nil {
		// Check if color is zero value (unset)
		if matchedSignal.Color.R == 0 && matchedSignal.Color.G == 0 && matchedSignal.Color.B == 0 {
			return AssertionResult{
				Passed:     false,
				Expected:   a.String(),
				Actual:     fmt.Sprintf("LED '%s' has no color information", matchedSignal.Name),
				Confidence: matchedSignal.Confidence,
				Message:    "Expected color, but LED has no color information",
			}
		}
		if !colorsMatch(*a.Expected.Color, matchedSignal.Color) {
			return AssertionResult{
				Passed:     false,
				Expected:   a.String(),
				Actual:     fmt.Sprintf("LED '%s' is RGB(%d,%d,%d)", matchedSignal.Name, matchedSignal.Color.R, matchedSignal.Color.G, matchedSignal.Color.B),
				Confidence: matchedSignal.Confidence,
				Message:    fmt.Sprintf("Expected RGB(%d,%d,%d), got RGB(%d,%d,%d)", a.Expected.Color.R, a.Expected.Color.G, a.Expected.Color.B, matchedSignal.Color.R, matchedSignal.Color.G, matchedSignal.Color.B),
			}
		}
	}

	// Check blink rate if specified
	if a.Expected.BlinkHz != nil {
		if matchedSignal.BlinkHz == 0 {
			return AssertionResult{
				Passed:     false,
				Expected:   a.String(),
				Actual:     fmt.Sprintf("LED '%s' is not blinking", matchedSignal.Name),
				Confidence: matchedSignal.Confidence,
				Message:    fmt.Sprintf("Expected blink rate %.2f Hz, but LED is not blinking", *a.Expected.BlinkHz),
			}
		}
		// Allow 10% tolerance on blink rate
		tolerance := *a.Expected.BlinkHz * 0.1
		if matchedSignal.BlinkHz < *a.Expected.BlinkHz-tolerance || matchedSignal.BlinkHz > *a.Expected.BlinkHz+tolerance {
			return AssertionResult{
				Passed:     false,
				Expected:   a.String(),
				Actual:     fmt.Sprintf("LED '%s' blinks at %.2f Hz", matchedSignal.Name, matchedSignal.BlinkHz),
				Confidence: matchedSignal.Confidence,
				Message:    fmt.Sprintf("Expected %.2f Hz, got %.2f Hz (outside tolerance)", *a.Expected.BlinkHz, matchedSignal.BlinkHz),
			}
		}
	}

	// All checks passed
	return AssertionResult{
		Passed:     true,
		Expected:   a.String(),
		Actual:     a.String(),
		Confidence: matchedSignal.Confidence,
		Message:    fmt.Sprintf("LED '%s' matches expected state", matchedSignal.Name),
	}
}

func (a *LEDAssertion) String() string {
	parts := []string{fmt.Sprintf("LED.%s", a.Name)}

	if a.Expected.On != nil {
		parts = append(parts, onOffString(*a.Expected.On))
	}

	if a.Expected.Color != nil {
		parts = append(parts, fmt.Sprintf("RGB(%d,%d,%d)", a.Expected.Color.R, a.Expected.Color.G, a.Expected.Color.B))
	}

	if a.Expected.BlinkHz != nil {
		parts = append(parts, fmt.Sprintf("%.2f Hz", *a.Expected.BlinkHz))
	}

	return strings.Join(parts, " ")
}

func onOffString(on bool) string {
	if on {
		return "ON"
	}
	return "OFF"
}

func colorsMatch(a, b core.RGB) bool {
	// Allow small tolerance for RGB values (±5)
	tolerance := 5
	return abs(int(a.R)-int(b.R)) <= tolerance &&
		abs(int(a.G)-int(b.G)) <= tolerance &&
		abs(int(a.B)-int(b.B)) <= tolerance
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// DisplayAssertion validates display content
type DisplayAssertion struct {
	Name     string
	Expected string
}

func (a *DisplayAssertion) Evaluate(obs *core.Observation) AssertionResult {
	// Find matching Display signal by name
	var matchedSignal *core.DisplaySignal
	for _, sig := range obs.Signals {
		if dispSig, ok := sig.(*core.DisplaySignal); ok {
			if dispSig.Name == a.Name {
				matchedSignal = dispSig
				break
			}
		}
	}

	if matchedSignal == nil {
		return AssertionResult{
			Passed:     false,
			Expected:   a.String(),
			Actual:     "Display not found in observation",
			Confidence: 0.0,
			Message:    fmt.Sprintf("Display '%s' not found in observation", a.Name),
		}
	}

	// Use contains() instead of exact match (OCR is noisy)
	if !strings.Contains(matchedSignal.Text, a.Expected) {
		return AssertionResult{
			Passed:     false,
			Expected:   a.String(),
			Actual:     fmt.Sprintf("Display '%s' shows: \"%s\"", matchedSignal.Name, matchedSignal.Text),
			Confidence: matchedSignal.Confidence,
			Message:    fmt.Sprintf("Expected text containing \"%s\", but display shows \"%s\"", a.Expected, matchedSignal.Text),
		}
	}

	return AssertionResult{
		Passed:     true,
		Expected:   a.String(),
		Actual:     fmt.Sprintf("Display '%s' contains \"%s\"", matchedSignal.Name, a.Expected),
		Confidence: matchedSignal.Confidence,
		Message:    fmt.Sprintf("Display '%s' contains expected text", matchedSignal.Name),
	}
}

func (a *DisplayAssertion) String() string {
	return fmt.Sprintf("Display.%s \"%s\"", a.Name, a.Expected)
}

// TimingAssertion validates boot timing
type TimingAssertion struct {
	MaxDurationMs int64
}

func (a *TimingAssertion) Evaluate(obs *core.Observation) AssertionResult {
	// Find BootTimingSignal
	var matchedSignal *core.BootTimingSignal
	for _, sig := range obs.Signals {
		if bootSig, ok := sig.(*core.BootTimingSignal); ok {
			matchedSignal = bootSig
			break
		}
	}

	// Graceful failure if signal missing
	if matchedSignal == nil {
		return AssertionResult{
			Passed:     false,
			Expected:   a.String(),
			Actual:     "No boot timing signal",
			Confidence: 0.0,
			Message:    "boot timing signal not present (did you capture from power-on?)",
		}
	}

	if matchedSignal.DurationMs > a.MaxDurationMs {
		return AssertionResult{
			Passed:     false,
			Expected:   a.String(),
			Actual:     fmt.Sprintf("Boot took %d ms", matchedSignal.DurationMs),
			Confidence: matchedSignal.Confidence,
			Message:    fmt.Sprintf("Expected boot <= %d ms, but took %d ms", a.MaxDurationMs, matchedSignal.DurationMs),
		}
	}

	return AssertionResult{
		Passed:     true,
		Expected:   a.String(),
		Actual:     fmt.Sprintf("Boot took %d ms", matchedSignal.DurationMs),
		Confidence: matchedSignal.Confidence,
		Message:    fmt.Sprintf("Boot completed in %d ms (within %d ms limit)", matchedSignal.DurationMs, a.MaxDurationMs),
	}
}

func (a *TimingAssertion) String() string {
	return fmt.Sprintf("BootTime < %d ms", a.MaxDurationMs)
}
