package vision

import (
	"testing"

	"github.com/perceptumx/percepta/internal/core"
)

func TestStructuredParser_ParseLEDSignals(t *testing.T) {
	// Mock test - actual API calls require ANTHROPIC_API_KEY
	// Real integration tests should be run manually with API key

	t.Skip("Skipping integration test - requires ANTHROPIC_API_KEY and API call")
}

func TestStructuredParser_ParseDisplaySignals(t *testing.T) {
	// Mock test - actual API calls require ANTHROPIC_API_KEY

	t.Skip("Skipping integration test - requires ANTHROPIC_API_KEY and API call")
}

func TestParseLEDToolResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{
			name: "single LED",
			input: map[string]interface{}{
				"leds": []interface{}{
					map[string]interface{}{
						"name":       "LED1",
						"on":         true,
						"color":      "red",
						"blink_hz":   float64(0),
						"confidence": float64(0.95),
					},
				},
			},
			expected: 1,
		},
		{
			name: "multiple LEDs",
			input: map[string]interface{}{
				"leds": []interface{}{
					map[string]interface{}{
						"name":       "LED1",
						"on":         true,
						"confidence": float64(0.90),
					},
					map[string]interface{}{
						"name":       "LED2",
						"on":         false,
						"confidence": float64(0.85),
					},
				},
			},
			expected: 2,
		},
		{
			name:     "empty input",
			input:    map[string]interface{}{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := parseLEDToolResponse(tt.input)
			if len(signals) != tt.expected {
				t.Errorf("expected %d signals, got %d", tt.expected, len(signals))
			}

			// Verify signal types
			for _, sig := range signals {
				if _, ok := sig.(core.LEDSignal); !ok {
					t.Errorf("expected LEDSignal, got %T", sig)
				}
			}
		})
	}
}

func TestParseDisplayToolResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{
			name: "single display",
			input: map[string]interface{}{
				"displays": []interface{}{
					map[string]interface{}{
						"name":       "LCD",
						"text":       "Hello World",
						"confidence": float64(0.92),
					},
				},
			},
			expected: 1,
		},
		{
			name: "multiple displays",
			input: map[string]interface{}{
				"displays": []interface{}{
					map[string]interface{}{
						"name":       "OLED",
						"text":       "Status: OK",
						"confidence": float64(0.88),
					},
					map[string]interface{}{
						"name":       "LCD",
						"text":       "Temp: 25C",
						"confidence": float64(0.93),
					},
				},
			},
			expected: 2,
		},
		{
			name:     "empty input",
			input:    map[string]interface{}{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := parseDisplayToolResponse(tt.input)
			if len(signals) != tt.expected {
				t.Errorf("expected %d signals, got %d", tt.expected, len(signals))
			}

			// Verify signal types
			for _, sig := range signals {
				if _, ok := sig.(core.DisplaySignal); !ok {
					t.Errorf("expected DisplaySignal, got %T", sig)
				}
			}
		})
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		colorStr string
		expected core.RGB
	}{
		{"red", core.RGB{R: 255, G: 0, B: 0}},
		{"green", core.RGB{R: 0, G: 255, B: 0}},
		{"blue", core.RGB{R: 0, G: 0, B: 255}},
		{"yellow", core.RGB{R: 255, G: 255, B: 0}},
		{"white", core.RGB{R: 255, G: 255, B: 255}},
		{"orange", core.RGB{R: 255, G: 165, B: 0}},
		{"unknown", core.RGB{R: 0, G: 0, B: 0}},
	}

	for _, tt := range tests {
		t.Run(tt.colorStr, func(t *testing.T) {
			result := parseColor(tt.colorStr)
			if result != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}
