package vision

import (
	"encoding/json"
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

// Helper function tests

func TestGetString(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		key      string
		expected string
	}{
		{
			name:     "existing string",
			data:     map[string]interface{}{"name": "test"},
			key:      "name",
			expected: "test",
		},
		{
			name:     "missing key",
			data:     map[string]interface{}{},
			key:      "name",
			expected: "",
		},
		{
			name:     "wrong type",
			data:     map[string]interface{}{"name": 123},
			key:      "name",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getString(tt.data, tt.key)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		key      string
		expected bool
	}{
		{
			name:     "true value",
			data:     map[string]interface{}{"on": true},
			key:      "on",
			expected: true,
		},
		{
			name:     "false value",
			data:     map[string]interface{}{"on": false},
			key:      "on",
			expected: false,
		},
		{
			name:     "missing key",
			data:     map[string]interface{}{},
			key:      "on",
			expected: false,
		},
		{
			name:     "wrong type",
			data:     map[string]interface{}{"on": "true"},
			key:      "on",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBool(tt.data, tt.key)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetFloat(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		key      string
		expected float64
	}{
		{
			name:     "float64 value",
			data:     map[string]interface{}{"confidence": float64(0.95)},
			key:      "confidence",
			expected: 0.95,
		},
		{
			name:     "int value",
			data:     map[string]interface{}{"count": 5},
			key:      "count",
			expected: 5.0,
		},
		{
			name:     "int64 value",
			data:     map[string]interface{}{"count": int64(10)},
			key:      "count",
			expected: 10.0,
		},
		{
			name:     "missing key",
			data:     map[string]interface{}{},
			key:      "confidence",
			expected: 0.0,
		},
		{
			name:     "wrong type",
			data:     map[string]interface{}{"confidence": "high"},
			key:      "confidence",
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFloat(tt.data, tt.key)
			if result != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestParseLEDToolResponse_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name:  "invalid input type",
			input: "not a map",
		},
		{
			name: "invalid leds type",
			input: map[string]interface{}{
				"leds": "not an array",
			},
		},
		{
			name: "invalid LED entry type",
			input: map[string]interface{}{
				"leds": []interface{}{
					"not a map",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := parseLEDToolResponse(tt.input)
			if len(signals) != 0 {
				t.Errorf("expected 0 signals for invalid input, got %d", len(signals))
			}
		})
	}
}

func TestParseLEDToolResponse_WithColor(t *testing.T) {
	input := map[string]interface{}{
		"leds": []interface{}{
			map[string]interface{}{
				"name":       "RGB1",
				"on":         true,
				"color":      "red",
				"confidence": float64(0.90),
			},
		},
	}

	signals := parseLEDToolResponse(input)
	if len(signals) != 1 {
		t.Fatalf("expected 1 signal, got %d", len(signals))
	}

	led, ok := signals[0].(core.LEDSignal)
	if !ok {
		t.Fatalf("expected LEDSignal, got %T", signals[0])
	}

	if led.Color.R != 255 || led.Color.G != 0 || led.Color.B != 0 {
		t.Errorf("expected red color RGB(255,0,0), got RGB(%d,%d,%d)", led.Color.R, led.Color.G, led.Color.B)
	}
}

func TestParseLEDToolResponse_WithBlinkHz(t *testing.T) {
	input := map[string]interface{}{
		"leds": []interface{}{
			map[string]interface{}{
				"name":       "STATUS",
				"on":         true,
				"blink_hz":   float64(2.5),
				"confidence": float64(0.92),
			},
		},
	}

	signals := parseLEDToolResponse(input)
	if len(signals) != 1 {
		t.Fatalf("expected 1 signal, got %d", len(signals))
	}

	led, ok := signals[0].(core.LEDSignal)
	if !ok {
		t.Fatalf("expected LEDSignal, got %T", signals[0])
	}

	if led.BlinkHz != 2.5 {
		t.Errorf("expected blink rate 2.5 Hz, got %f", led.BlinkHz)
	}
}

func TestParseLEDToolResponse_ZeroBlinkHz(t *testing.T) {
	input := map[string]interface{}{
		"leds": []interface{}{
			map[string]interface{}{
				"name":       "POWER",
				"on":         true,
				"blink_hz":   float64(0),
				"confidence": float64(0.95),
			},
		},
	}

	signals := parseLEDToolResponse(input)
	led := signals[0].(core.LEDSignal)

	if led.BlinkHz != 0 {
		t.Errorf("expected blink rate 0 (steady), got %f", led.BlinkHz)
	}
}

func TestParseDisplayToolResponse_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name:  "invalid input type",
			input: "not a map",
		},
		{
			name: "invalid displays type",
			input: map[string]interface{}{
				"displays": "not an array",
			},
		},
		{
			name: "invalid display entry type",
			input: map[string]interface{}{
				"displays": []interface{}{
					"not a map",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signals := parseDisplayToolResponse(tt.input)
			if len(signals) != 0 {
				t.Errorf("expected 0 signals for invalid input, got %d", len(signals))
			}
		})
	}
}

func TestParseDisplayToolResponse_FullData(t *testing.T) {
	input := map[string]interface{}{
		"displays": []interface{}{
			map[string]interface{}{
				"name":       "LCD",
				"text":       "Temperature: 25C",
				"confidence": float64(0.93),
			},
		},
	}

	signals := parseDisplayToolResponse(input)
	if len(signals) != 1 {
		t.Fatalf("expected 1 signal, got %d", len(signals))
	}

	display, ok := signals[0].(core.DisplaySignal)
	if !ok {
		t.Fatalf("expected DisplaySignal, got %T", signals[0])
	}

	if display.Name != "LCD" {
		t.Errorf("expected name 'LCD', got '%s'", display.Name)
	}

	if display.Text != "Temperature: 25C" {
		t.Errorf("expected text 'Temperature: 25C', got '%s'", display.Text)
	}

	if display.Confidence != 0.93 {
		t.Errorf("expected confidence 0.93, got %f", display.Confidence)
	}
}

func TestLedDetectionTool_Schema(t *testing.T) {
	tool := ledDetectionTool()

	if tool.Name != "report_led_signals" {
		t.Errorf("expected tool name 'report_led_signals', got '%s'", tool.Name)
	}

	// Verify schema structure
	if tool.InputSchema.Type != "object" {
		t.Errorf("expected schema type 'object', got '%s'", tool.InputSchema.Type)
	}

	if tool.InputSchema.Properties == nil {
		t.Fatal("expected non-nil properties")
	}

	if len(tool.InputSchema.Required) == 0 {
		t.Error("expected required fields")
	}
}

func TestDisplayDetectionTool_Schema(t *testing.T) {
	tool := displayDetectionTool()

	if tool.Name != "report_display_content" {
		t.Errorf("expected tool name 'report_display_content', got '%s'", tool.Name)
	}

	// Verify schema structure
	if tool.InputSchema.Type != "object" {
		t.Errorf("expected schema type 'object', got '%s'", tool.InputSchema.Type)
	}

	if tool.InputSchema.Properties == nil {
		t.Fatal("expected non-nil properties")
	}

	if len(tool.InputSchema.Required) == 0 {
		t.Error("expected required fields")
	}
}

func TestNewStructuredParser(t *testing.T) {
	// Test with nil client
	parser := NewStructuredParser(nil)
	if parser == nil {
		t.Fatal("expected non-nil parser")
	}

	if parser.client != nil {
		t.Error("expected nil client to be stored")
	}
}

func TestParseColor_AllColors(t *testing.T) {
	colors := []struct {
		name string
		rgb  core.RGB
	}{
		{"red", core.RGB{R: 255, G: 0, B: 0}},
		{"green", core.RGB{R: 0, G: 255, B: 0}},
		{"blue", core.RGB{R: 0, G: 0, B: 255}},
		{"yellow", core.RGB{R: 255, G: 255, B: 0}},
		{"white", core.RGB{R: 255, G: 255, B: 255}},
		{"orange", core.RGB{R: 255, G: 165, B: 0}},
	}

	for _, tt := range colors {
		t.Run(tt.name, func(t *testing.T) {
			result := parseColor(tt.name)
			if result != tt.rgb {
				t.Errorf("color '%s': expected RGB(%d,%d,%d), got RGB(%d,%d,%d)",
					tt.name, tt.rgb.R, tt.rgb.G, tt.rgb.B, result.R, result.G, result.B)
			}
		})
	}
}

func TestParseLEDToolResponse_CompleteSignal(t *testing.T) {
	// Test LED with all fields populated
	input := map[string]interface{}{
		"leds": []interface{}{
			map[string]interface{}{
				"name":       "RGB_STATUS",
				"on":         true,
				"color":      "green",
				"blink_hz":   float64(1.5),
				"confidence": float64(0.88),
			},
		},
	}

	signals := parseLEDToolResponse(input)
	if len(signals) != 1 {
		t.Fatalf("expected 1 signal, got %d", len(signals))
	}

	led := signals[0].(core.LEDSignal)

	if led.Name != "RGB_STATUS" {
		t.Errorf("expected name 'RGB_STATUS', got '%s'", led.Name)
	}

	if !led.On {
		t.Error("expected LED to be on")
	}

	if led.Color.G != 255 {
		t.Errorf("expected green color, got RGB(%d,%d,%d)", led.Color.R, led.Color.G, led.Color.B)
	}

	if led.BlinkHz != 1.5 {
		t.Errorf("expected blink rate 1.5, got %f", led.BlinkHz)
	}

	if led.Confidence != 0.88 {
		t.Errorf("expected confidence 0.88, got %f", led.Confidence)
	}
}

func TestGetFloat_JSONNumber(t *testing.T) {
	// Test json.Number type handling
	tests := []struct {
		name     string
		value    interface{}
		expected float64
	}{
		{
			name:     "valid json.Number",
			value:    json.Number("3.14"),
			expected: 3.14,
		},
		{
			name:     "invalid json.Number",
			value:    json.Number("not-a-number"),
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := map[string]interface{}{"value": tt.value}
			result := getFloat(data, "value")
			if result != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}
