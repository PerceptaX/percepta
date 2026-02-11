package assertions

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/perceptumx/percepta/internal/core"
)

// Parse converts DSL string to Assertion
func Parse(dsl string) (Assertion, error) {
	dsl = strings.TrimSpace(dsl)

	// LED patterns
	if strings.HasPrefix(dsl, "LED.") {
		return parseLED(dsl)
	}

	// Display patterns
	if strings.HasPrefix(dsl, "Display.") {
		return parseDisplay(dsl)
	}

	// Timing patterns
	if strings.HasPrefix(dsl, "BootTime") {
		return parseTiming(dsl)
	}

	return nil, fmt.Errorf("unknown assertion type: %s", dsl)
}

func parseLED(dsl string) (*LEDAssertion, error) {
	// LED.name [ON|OFF|BLINK freq|COLOR RGB(r,g,b)]
	// Match: LED.{name} {rest}
	parts := strings.SplitN(dsl, " ", 2)
	if len(parts) < 1 {
		return nil, fmt.Errorf("invalid LED assertion syntax: %s", dsl)
	}

	// Extract name from LED.name
	namePattern := regexp.MustCompile(`^LED\.([a-zA-Z0-9_-]+)`)
	nameMatch := namePattern.FindStringSubmatch(parts[0])
	if nameMatch == nil {
		return nil, fmt.Errorf("invalid LED name format: %s", dsl)
	}

	assertion := &LEDAssertion{
		Name:     nameMatch[1],
		Expected: LEDState{},
	}

	// If no state specified, just check for existence
	if len(parts) == 1 {
		return assertion, nil
	}

	state := strings.TrimSpace(parts[1])

	// Check for ON
	if strings.EqualFold(state, "ON") {
		on := true
		assertion.Expected.On = &on
		return assertion, nil
	}

	// Check for OFF
	if strings.EqualFold(state, "OFF") {
		off := false
		assertion.Expected.On = &off
		return assertion, nil
	}

	// Check for BLINK {freq}Hz
	blinkPattern := regexp.MustCompile(`(?i)^BLINK\s+([\d.]+)\s*Hz$`)
	if blinkMatch := blinkPattern.FindStringSubmatch(state); blinkMatch != nil {
		freq, err := strconv.ParseFloat(blinkMatch[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid blink frequency: %s", blinkMatch[1])
		}
		assertion.Expected.BlinkHz = &freq
		return assertion, nil
	}

	// Check for COLOR RGB(r,g,b)
	colorPattern := regexp.MustCompile(`(?i)^COLOR\s+RGB\((\d+),\s*(\d+),\s*(\d+)\)$`)
	if colorMatch := colorPattern.FindStringSubmatch(state); colorMatch != nil {
		r, err := strconv.ParseUint(colorMatch[1], 10, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid red value: %s", colorMatch[1])
		}
		g, err := strconv.ParseUint(colorMatch[2], 10, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid green value: %s", colorMatch[2])
		}
		b, err := strconv.ParseUint(colorMatch[3], 10, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid blue value: %s", colorMatch[3])
		}
		assertion.Expected.Color = &core.RGB{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
		}
		return assertion, nil
	}

	return nil, fmt.Errorf("unknown LED state format: %s", state)
}

func parseDisplay(dsl string) (*DisplayAssertion, error) {
	// Display.name "text"
	pattern := regexp.MustCompile(`^Display\.([a-zA-Z0-9_-]+)\s+"([^"]+)"$`)
	matches := pattern.FindStringSubmatch(dsl)
	if matches == nil {
		return nil, fmt.Errorf("invalid Display assertion syntax: %s (expected: Display.name \"text\")", dsl)
	}

	return &DisplayAssertion{
		Name:     matches[1],
		Expected: matches[2],
	}, nil
}

func parseTiming(dsl string) (*TimingAssertion, error) {
	// BootTime < {ms}ms
	pattern := regexp.MustCompile(`^BootTime\s+<\s+(\d+)\s*ms$`)
	matches := pattern.FindStringSubmatch(dsl)
	if matches == nil {
		return nil, fmt.Errorf("invalid BootTime assertion syntax: %s (expected: BootTime < 3000ms)", dsl)
	}

	duration, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid duration value: %s", matches[1])
	}

	return &TimingAssertion{
		MaxDurationMs: duration,
	}, nil
}
