package diff

import (
	"fmt"
	"math"

	"github.com/perceptumx/percepta/internal/core"
)

// Compare compares two observations and returns detected changes
func Compare(from, to *core.Observation) *DiffResult {
	result := &DiffResult{
		DeviceID:      from.DeviceID,
		FromFirmware:  from.FirmwareHash,
		ToFirmware:    to.FirmwareHash,
		FromTimestamp: from.Timestamp.Format("2006-01-02 15:04:05"),
		ToTimestamp:   to.Timestamp.Format("2006-01-02 15:04:05"),
		Changes:       make([]SignalChange, 0),
	}

	// Normalize signals for comparison
	fromSignals := normalizeSignals(from.Signals)
	toSignals := normalizeSignals(to.Signals)

	// Build maps by signal name for comparison
	fromMap := make(map[string]NormalizedSignal)
	toMap := make(map[string]NormalizedSignal)

	for _, sig := range fromSignals {
		fromMap[sig.Name] = sig
	}
	for _, sig := range toSignals {
		toMap[sig.Name] = sig
	}

	// Detect removed signals (in 'from' but not in 'to')
	for name, fromSig := range fromMap {
		if _, exists := toMap[name]; !exists {
			result.Changes = append(result.Changes, SignalChange{
				Type:      ChangeRemoved,
				Name:      name,
				FromState: formatSignalState(fromSig),
				ToState:   "",
				Details:   "",
			})
		}
	}

	// Detect added and modified signals
	for name, toSig := range toMap {
		fromSig, existsInFrom := fromMap[name]

		if !existsInFrom {
			// Signal added
			result.Changes = append(result.Changes, SignalChange{
				Type:      ChangeAdded,
				Name:      name,
				FromState: "",
				ToState:   formatSignalState(toSig),
				Details:   "",
			})
		} else {
			// Signal exists in both - check for modifications
			if !signalsEqual(fromSig, toSig) {
				details := describeChange(fromSig, toSig)
				result.Changes = append(result.Changes, SignalChange{
					Type:      ChangeModified,
					Name:      name,
					FromState: formatSignalState(fromSig),
					ToState:   formatSignalState(toSig),
					Details:   details,
				})
			}
		}
	}

	return result
}

// normalizeSignals converts signals to normalized form for comparison
func normalizeSignals(signals []core.Signal) []NormalizedSignal {
	normalized := make([]NormalizedSignal, 0, len(signals))

	for _, sig := range signals {
		switch s := sig.(type) {
		case core.LEDSignal:
			normalized = append(normalized, NormalizedSignal{
				Name:    s.Name,
				Type:    "led",
				Signal:  s,
				OnState: s.On,
				Blink:   s.BlinkHz > 0,
				BlinkHz: normalizeBlinkHz(s.BlinkHz),
				Color:   s.Color,
			})
		case core.DisplaySignal:
			normalized = append(normalized, NormalizedSignal{
				Name:   s.Name,
				Type:   "display",
				Signal: s,
				Text:   s.Text,
			})
		case core.BootTimingSignal:
			normalized = append(normalized, NormalizedSignal{
				Name:   "boot",
				Type:   "boot_timing",
				Signal: s,
			})
		}
	}

	return normalized
}

// normalizeBlinkHz rounds blink rate to 1 decimal place to handle Claude Vision fluctuations
func normalizeBlinkHz(hz float64) float64 {
	if hz == 0 {
		return 0
	}
	return math.Round(hz*10) / 10
}

// signalsEqual compares two normalized signals for equality
func signalsEqual(a, b NormalizedSignal) bool {
	if a.Type != b.Type {
		return false
	}

	switch a.Type {
	case "led":
		aLED := a.Signal.(core.LEDSignal)
		bLED := b.Signal.(core.LEDSignal)

		// Compare state
		if aLED.On != bLED.On {
			return false
		}

		// Compare blinking (normalized)
		if normalizeBlinkHz(aLED.BlinkHz) != normalizeBlinkHz(bLED.BlinkHz) {
			return false
		}

		// Compare color exactly
		if aLED.Color != bLED.Color {
			return false
		}

		return true

	case "display":
		aDisplay := a.Signal.(core.DisplaySignal)
		bDisplay := b.Signal.(core.DisplaySignal)

		// Compare text exactly
		return aDisplay.Text == bDisplay.Text

	case "boot_timing":
		aBoot := a.Signal.(core.BootTimingSignal)
		bBoot := b.Signal.(core.BootTimingSignal)

		// Compare duration exactly
		return aBoot.DurationMs == bBoot.DurationMs
	}

	return false
}

// formatSignalState creates a human-readable representation of signal state
func formatSignalState(sig NormalizedSignal) string {
	switch sig.Type {
	case "led":
		led := sig.Signal.(core.LEDSignal)
		state := "OFF"
		if led.On {
			state = "ON"
		}

		var parts []string
		parts = append(parts, state)

		// Add color if present
		if led.Color.R > 0 || led.Color.G > 0 || led.Color.B > 0 {
			parts = append(parts, formatColor(led.Color))
		}

		// Add blink info
		if led.BlinkHz > 0 {
			parts = append(parts, fmt.Sprintf("blinking %.1fHz", normalizeBlinkHz(led.BlinkHz)))
		} else {
			parts = append(parts, "solid")
		}

		result := ""
		for i, part := range parts {
			if i > 0 {
				result += " "
			}
			result += part
		}
		return result

	case "display":
		display := sig.Signal.(core.DisplaySignal)
		return fmt.Sprintf("\"%s\"", display.Text)

	case "boot_timing":
		boot := sig.Signal.(core.BootTimingSignal)
		return fmt.Sprintf("%dms", boot.DurationMs)
	}

	return ""
}

// formatColor converts RGB to a color name or hex if no name matches
func formatColor(rgb core.RGB) string {
	// Common color names
	colors := map[string]core.RGB{
		"red":    {R: 255, G: 0, B: 0},
		"green":  {R: 0, G: 255, B: 0},
		"blue":   {R: 0, G: 0, B: 255},
		"yellow": {R: 255, G: 255, B: 0},
		"purple": {R: 128, G: 0, B: 128},
		"cyan":   {R: 0, G: 255, B: 255},
		"white":  {R: 255, G: 255, B: 255},
	}

	// Check for exact match
	for name, color := range colors {
		if rgb.R == color.R && rgb.G == color.G && rgb.B == color.B {
			return name
		}
	}

	// Check for close match (within 30 for each channel)
	for name, color := range colors {
		rDiff := abs(int(rgb.R) - int(color.R))
		gDiff := abs(int(rgb.G) - int(color.G))
		bDiff := abs(int(rgb.B) - int(color.B))

		if rDiff <= 30 && gDiff <= 30 && bDiff <= 30 {
			return name
		}
	}

	// Return hex if no match
	return fmt.Sprintf("#%02x%02x%02x", rgb.R, rgb.G, rgb.B)
}

// describeChange creates a detailed description of what changed
func describeChange(from, to NormalizedSignal) string {
	if from.Type != to.Type {
		return fmt.Sprintf("type changed: %s → %s", from.Type, to.Type)
	}

	switch from.Type {
	case "led":
		fromLED := from.Signal.(core.LEDSignal)
		toLED := to.Signal.(core.LEDSignal)

		var changes []string

		// State change
		if fromLED.On != toLED.On {
			fromState := "OFF"
			toState := "OFF"
			if fromLED.On {
				fromState = "ON"
			}
			if toLED.On {
				toState = "ON"
			}
			changes = append(changes, fmt.Sprintf("%s→%s", fromState, toState))
		}

		// Color change
		if fromLED.Color != toLED.Color {
			changes = append(changes, fmt.Sprintf("color: %s→%s", formatColor(fromLED.Color), formatColor(toLED.Color)))
		}

		// Blink rate change
		fromHz := normalizeBlinkHz(fromLED.BlinkHz)
		toHz := normalizeBlinkHz(toLED.BlinkHz)
		if fromHz != toHz {
			if fromHz == 0 && toHz > 0 {
				changes = append(changes, fmt.Sprintf("solid→blinking %.1fHz", toHz))
			} else if fromHz > 0 && toHz == 0 {
				changes = append(changes, "blinking→solid")
			} else {
				changes = append(changes, fmt.Sprintf("blink: %.1fHz→%.1fHz", fromHz, toHz))
			}
		}

		if len(changes) == 0 {
			return ""
		}

		result := ""
		for i, change := range changes {
			if i > 0 {
				result += ", "
			}
			result += change
		}
		return result

	case "display":
		fromDisplay := from.Signal.(core.DisplaySignal)
		toDisplay := to.Signal.(core.DisplaySignal)
		return fmt.Sprintf("text: \"%s\"→\"%s\"", fromDisplay.Text, toDisplay.Text)

	case "boot_timing":
		fromBoot := from.Signal.(core.BootTimingSignal)
		toBoot := to.Signal.(core.BootTimingSignal)
		return fmt.Sprintf("duration: %dms→%dms", fromBoot.DurationMs, toBoot.DurationMs)
	}

	return ""
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
