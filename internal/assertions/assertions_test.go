package assertions

import (
	"strings"
	"testing"

	"github.com/perceptumx/percepta/internal/core"
)

func TestDisplayChangedAssertion_Pass(t *testing.T) {
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			core.DisplaySignal{
				Name:       "LCD",
				Text:       "Ready",
				Confidence: 0.90,
				Changed:    true,
				History: []core.DisplayTextEntry{
					{OffsetMs: 0, Text: "Boot", Confidence: 0.90},
					{OffsetMs: 400, Text: "Init...", Confidence: 0.88},
					{OffsetMs: 800, Text: "Ready", Confidence: 0.90},
				},
			},
		},
	}

	assertion := &DisplayChangedAssertion{
		Name:     "LCD",
		FromText: "Boot",
		ToText:   "Ready",
	}

	result := assertion.Evaluate(obs)
	if !result.Passed {
		t.Errorf("Expected assertion to pass, but got: %s", result.Message)
	}
}

func TestDisplayChangedAssertion_SubstringMatch(t *testing.T) {
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			core.DisplaySignal{
				Name:       "LCD",
				Text:       "System Ready v2.0",
				Confidence: 0.90,
				Changed:    true,
				History: []core.DisplayTextEntry{
					{OffsetMs: 0, Text: "Booting System...", Confidence: 0.90},
					{OffsetMs: 800, Text: "System Ready v2.0", Confidence: 0.90},
				},
			},
		},
	}

	assertion := &DisplayChangedAssertion{
		Name:     "LCD",
		FromText: "Boot",
		ToText:   "Ready",
	}

	result := assertion.Evaluate(obs)
	if !result.Passed {
		t.Errorf("Expected assertion to pass with substring match, but got: %s", result.Message)
	}
}

func TestDisplayChangedAssertion_StaticDisplay(t *testing.T) {
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			core.DisplaySignal{
				Name:       "LCD",
				Text:       "Ready",
				Confidence: 0.90,
				Changed:    false,
			},
		},
	}

	assertion := &DisplayChangedAssertion{
		Name:     "LCD",
		FromText: "Boot",
		ToText:   "Ready",
	}

	result := assertion.Evaluate(obs)
	if result.Passed {
		t.Errorf("Expected assertion to fail for static display")
	}
}

func TestDisplayChangedAssertion_FromTextNotFound(t *testing.T) {
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			core.DisplaySignal{
				Name:       "LCD",
				Text:       "Ready",
				Confidence: 0.90,
				Changed:    true,
				History: []core.DisplayTextEntry{
					{OffsetMs: 0, Text: "Init", Confidence: 0.88},
					{OffsetMs: 400, Text: "Ready", Confidence: 0.90},
				},
			},
		},
	}

	assertion := &DisplayChangedAssertion{
		Name:     "LCD",
		FromText: "Boot",
		ToText:   "Ready",
	}

	result := assertion.Evaluate(obs)
	if result.Passed {
		t.Errorf("Expected assertion to fail when FromText not found")
	}
}

func TestDisplayChangedAssertion_WrongOrder(t *testing.T) {
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			core.DisplaySignal{
				Name:       "LCD",
				Text:       "Boot",
				Confidence: 0.90,
				Changed:    true,
				History: []core.DisplayTextEntry{
					{OffsetMs: 0, Text: "Ready", Confidence: 0.90},
					{OffsetMs: 400, Text: "Boot", Confidence: 0.88},
				},
			},
		},
	}

	assertion := &DisplayChangedAssertion{
		Name:     "LCD",
		FromText: "Boot",
		ToText:   "Ready",
	}

	result := assertion.Evaluate(obs)
	if result.Passed {
		t.Errorf("Expected assertion to fail when texts are in wrong order")
	}
}

func TestDisplayChangedAssertion_DisplayNotFound(t *testing.T) {
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true},
		},
	}

	assertion := &DisplayChangedAssertion{
		Name:     "LCD",
		FromText: "Boot",
		ToText:   "Ready",
	}

	result := assertion.Evaluate(obs)
	if result.Passed {
		t.Errorf("Expected assertion to fail when display not found")
	}
}

func TestDisplayChangedAssertion_CaseInsensitiveName(t *testing.T) {
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			core.DisplaySignal{
				Name:       "lcd",
				Text:       "Ready",
				Confidence: 0.90,
				Changed:    true,
				History: []core.DisplayTextEntry{
					{OffsetMs: 0, Text: "Boot", Confidence: 0.90},
					{OffsetMs: 400, Text: "Ready", Confidence: 0.90},
				},
			},
		},
	}

	assertion := &DisplayChangedAssertion{
		Name:     "LCD",
		FromText: "Boot",
		ToText:   "Ready",
	}

	result := assertion.Evaluate(obs)
	if !result.Passed {
		t.Errorf("Expected case-insensitive name match to pass, but got: %s", result.Message)
	}
}

func TestParse_DisplayChanged(t *testing.T) {
	assertion, err := Parse(`Display.LCD CHANGED "Boot" -> "Ready"`)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	dca, ok := assertion.(*DisplayChangedAssertion)
	if !ok {
		t.Fatalf("Expected *DisplayChangedAssertion, got %T", assertion)
	}

	if dca.Name != "LCD" {
		t.Errorf("Expected name 'LCD', got '%s'", dca.Name)
	}
	if dca.FromText != "Boot" {
		t.Errorf("Expected FromText 'Boot', got '%s'", dca.FromText)
	}
	if dca.ToText != "Ready" {
		t.Errorf("Expected ToText 'Ready', got '%s'", dca.ToText)
	}
}

func TestParse_DisplayChangedWithSpaces(t *testing.T) {
	assertion, err := Parse(`Display.OLED CHANGED "System Boot" -> "System Ready"`)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	dca, ok := assertion.(*DisplayChangedAssertion)
	if !ok {
		t.Fatalf("Expected *DisplayChangedAssertion, got %T", assertion)
	}

	if dca.FromText != "System Boot" {
		t.Errorf("Expected FromText 'System Boot', got '%s'", dca.FromText)
	}
	if dca.ToText != "System Ready" {
		t.Errorf("Expected ToText 'System Ready', got '%s'", dca.ToText)
	}
}

func TestParse_DisplayStatic(t *testing.T) {
	assertion, err := Parse(`Display.LCD "Ready"`)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	_, ok := assertion.(*DisplayAssertion)
	if !ok {
		t.Fatalf("Expected *DisplayAssertion, got %T", assertion)
	}
}

func TestDisplayChangedAssertion_String(t *testing.T) {
	assertion := &DisplayChangedAssertion{
		Name:     "LCD",
		FromText: "Boot",
		ToText:   "Ready",
	}

	expected := `Display.LCD CHANGED "Boot" -> "Ready"`
	if assertion.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, assertion.String())
	}
}

// LED Assertion Tests

func TestLEDAssertion_ON(t *testing.T) {
	on := true
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.LEDSignal{
				Name:       "POWER",
				On:         true,
				Confidence: 0.95,
			},
		},
	}

	assertion := &LEDAssertion{
		Name:     "POWER",
		Expected: LEDState{On: &on},
	}

	result := assertion.Evaluate(obs)
	if !result.Passed {
		t.Errorf("Expected LED ON assertion to pass, got: %s", result.Message)
	}
}

func TestLEDAssertion_OFF(t *testing.T) {
	off := false
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.LEDSignal{
				Name:       "ERROR",
				On:         false,
				Confidence: 0.90,
			},
		},
	}

	assertion := &LEDAssertion{
		Name:     "ERROR",
		Expected: LEDState{On: &off},
	}

	result := assertion.Evaluate(obs)
	if !result.Passed {
		t.Errorf("Expected LED OFF assertion to pass, got: %s", result.Message)
	}
}

func TestLEDAssertion_BlinkRate(t *testing.T) {
	blinkHz := 2.0
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.LEDSignal{
				Name:       "STATUS",
				On:         true,
				BlinkHz:    2.0,
				Confidence: 0.92,
			},
		},
	}

	assertion := &LEDAssertion{
		Name:     "STATUS",
		Expected: LEDState{BlinkHz: &blinkHz},
	}

	result := assertion.Evaluate(obs)
	if !result.Passed {
		t.Errorf("Expected LED blink assertion to pass, got: %s", result.Message)
	}
}

func TestLEDAssertion_BlinkRate_WithTolerance(t *testing.T) {
	blinkHz := 2.0
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.LEDSignal{
				Name:       "STATUS",
				On:         true,
				BlinkHz:    2.05, // Within 10% tolerance
				Confidence: 0.92,
			},
		},
	}

	assertion := &LEDAssertion{
		Name:     "STATUS",
		Expected: LEDState{BlinkHz: &blinkHz},
	}

	result := assertion.Evaluate(obs)
	if !result.Passed {
		t.Errorf("Expected LED blink with tolerance to pass, got: %s", result.Message)
	}
}

func TestLEDAssertion_BlinkRate_OutsideTolerance(t *testing.T) {
	blinkHz := 2.0
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.LEDSignal{
				Name:       "STATUS",
				On:         true,
				BlinkHz:    3.0, // Outside 10% tolerance
				Confidence: 0.92,
			},
		},
	}

	assertion := &LEDAssertion{
		Name:     "STATUS",
		Expected: LEDState{BlinkHz: &blinkHz},
	}

	result := assertion.Evaluate(obs)
	if result.Passed {
		t.Error("Expected LED blink assertion to fail when outside tolerance")
	}
}

func TestLEDAssertion_BlinkRate_NotBlinking(t *testing.T) {
	blinkHz := 2.0
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.LEDSignal{
				Name:       "STATUS",
				On:         true,
				BlinkHz:    0, // Not blinking
				Confidence: 0.92,
			},
		},
	}

	assertion := &LEDAssertion{
		Name:     "STATUS",
		Expected: LEDState{BlinkHz: &blinkHz},
	}

	result := assertion.Evaluate(obs)
	if result.Passed {
		t.Error("Expected LED blink assertion to fail when LED is not blinking")
	}
}

func TestLEDAssertion_Color(t *testing.T) {
	color := core.RGB{R: 255, G: 0, B: 0}
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.LEDSignal{
				Name:       "RGB1",
				On:         true,
				Color:      core.RGB{R: 255, G: 0, B: 0},
				Confidence: 0.88,
			},
		},
	}

	assertion := &LEDAssertion{
		Name:     "RGB1",
		Expected: LEDState{Color: &color},
	}

	result := assertion.Evaluate(obs)
	if !result.Passed {
		t.Errorf("Expected LED color assertion to pass, got: %s", result.Message)
	}
}

func TestLEDAssertion_Color_WithTolerance(t *testing.T) {
	color := core.RGB{R: 255, G: 0, B: 0}
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.LEDSignal{
				Name:       "RGB1",
				On:         true,
				Color:      core.RGB{R: 252, G: 3, B: 2}, // Within tolerance
				Confidence: 0.88,
			},
		},
	}

	assertion := &LEDAssertion{
		Name:     "RGB1",
		Expected: LEDState{Color: &color},
	}

	result := assertion.Evaluate(obs)
	if !result.Passed {
		t.Errorf("Expected LED color with tolerance to pass, got: %s", result.Message)
	}
}

func TestLEDAssertion_Color_Mismatch(t *testing.T) {
	color := core.RGB{R: 255, G: 0, B: 0}
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.LEDSignal{
				Name:       "RGB1",
				On:         true,
				Color:      core.RGB{R: 0, G: 255, B: 0}, // Green instead of red
				Confidence: 0.88,
			},
		},
	}

	assertion := &LEDAssertion{
		Name:     "RGB1",
		Expected: LEDState{Color: &color},
	}

	result := assertion.Evaluate(obs)
	if result.Passed {
		t.Error("Expected LED color assertion to fail for mismatched color")
	}
}

func TestLEDAssertion_Color_NoColorInfo(t *testing.T) {
	color := core.RGB{R: 255, G: 0, B: 0}
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.LEDSignal{
				Name:       "LED1",
				On:         true,
				Color:      core.RGB{R: 0, G: 0, B: 0}, // No color info
				Confidence: 0.88,
			},
		},
	}

	assertion := &LEDAssertion{
		Name:     "LED1",
		Expected: LEDState{Color: &color},
	}

	result := assertion.Evaluate(obs)
	if result.Passed {
		t.Error("Expected LED color assertion to fail when LED has no color info")
	}
}

func TestLEDAssertion_NotFound(t *testing.T) {
	on := true
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals:  []core.Signal{},
	}

	assertion := &LEDAssertion{
		Name:     "POWER",
		Expected: LEDState{On: &on},
	}

	result := assertion.Evaluate(obs)
	if result.Passed {
		t.Error("Expected LED assertion to fail when LED not found")
	}
}

func TestLEDAssertion_SingleLEDFallback(t *testing.T) {
	on := true
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.LEDSignal{
				Name:       "LED1",
				On:         true,
				Confidence: 0.95,
			},
		},
	}

	assertion := &LEDAssertion{
		Name:     "DIFFERENT_NAME",
		Expected: LEDState{On: &on},
	}

	result := assertion.Evaluate(obs)
	if !result.Passed {
		t.Errorf("Expected single LED fallback to work, got: %s", result.Message)
	}
}

func TestLEDAssertion_CaseInsensitive(t *testing.T) {
	on := true
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.LEDSignal{
				Name:       "power",
				On:         true,
				Confidence: 0.95,
			},
		},
	}

	assertion := &LEDAssertion{
		Name:     "POWER",
		Expected: LEDState{On: &on},
	}

	result := assertion.Evaluate(obs)
	if !result.Passed {
		t.Errorf("Expected case-insensitive LED name match, got: %s", result.Message)
	}
}

func TestLEDAssertion_StateMismatch(t *testing.T) {
	on := true
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.LEDSignal{
				Name:       "POWER",
				On:         false,
				Confidence: 0.95,
			},
		},
	}

	assertion := &LEDAssertion{
		Name:     "POWER",
		Expected: LEDState{On: &on},
	}

	result := assertion.Evaluate(obs)
	if result.Passed {
		t.Error("Expected LED state assertion to fail for state mismatch")
	}
}

func TestLEDAssertion_String(t *testing.T) {
	on := true
	blinkHz := 2.5
	color := core.RGB{R: 255, G: 128, B: 0}

	tests := []struct {
		name     string
		assertion *LEDAssertion
		expected string
	}{
		{
			name: "ON state",
			assertion: &LEDAssertion{
				Name:     "POWER",
				Expected: LEDState{On: &on},
			},
			expected: "LED.POWER ON",
		},
		{
			name: "Blink rate",
			assertion: &LEDAssertion{
				Name:     "STATUS",
				Expected: LEDState{BlinkHz: &blinkHz},
			},
			expected: "LED.STATUS 2.50 Hz",
		},
		{
			name: "Color",
			assertion: &LEDAssertion{
				Name:     "RGB1",
				Expected: LEDState{Color: &color},
			},
			expected: "LED.RGB1 RGB(255,128,0)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.assertion.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// Timing Assertion Tests

func TestTimingAssertion_Pass(t *testing.T) {
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.BootTimingSignal{
				DurationMs: 2000,
				Confidence: 0.90,
			},
		},
	}

	assertion := &TimingAssertion{
		MaxDurationMs: 3000,
	}

	result := assertion.Evaluate(obs)
	if !result.Passed {
		t.Errorf("Expected timing assertion to pass, got: %s", result.Message)
	}
}

func TestTimingAssertion_Fail(t *testing.T) {
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.BootTimingSignal{
				DurationMs: 4000,
				Confidence: 0.90,
			},
		},
	}

	assertion := &TimingAssertion{
		MaxDurationMs: 3000,
	}

	result := assertion.Evaluate(obs)
	if result.Passed {
		t.Error("Expected timing assertion to fail when duration exceeds max")
	}
}

func TestTimingAssertion_ExactLimit(t *testing.T) {
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.BootTimingSignal{
				DurationMs: 3000,
				Confidence: 0.90,
			},
		},
	}

	assertion := &TimingAssertion{
		MaxDurationMs: 3000,
	}

	result := assertion.Evaluate(obs)
	if !result.Passed {
		t.Errorf("Expected timing assertion to pass at exact limit, got: %s", result.Message)
	}
}

func TestTimingAssertion_SignalNotFound(t *testing.T) {
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals:  []core.Signal{},
	}

	assertion := &TimingAssertion{
		MaxDurationMs: 3000,
	}

	result := assertion.Evaluate(obs)
	if result.Passed {
		t.Error("Expected timing assertion to fail when signal not found")
	}

	if !strings.Contains(result.Message, "boot timing signal not present") {
		t.Errorf("Expected helpful error message, got: %s", result.Message)
	}
}

func TestTimingAssertion_String(t *testing.T) {
	assertion := &TimingAssertion{
		MaxDurationMs: 5000,
	}

	expected := "BootTime < 5000 ms"
	if assertion.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, assertion.String())
	}
}

// DisplayAssertion Tests

func TestDisplayAssertion_Pass(t *testing.T) {
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.DisplaySignal{
				Name:       "LCD",
				Text:       "System Ready",
				Confidence: 0.92,
			},
		},
	}

	assertion := &DisplayAssertion{
		Name:     "LCD",
		Expected: "Ready",
	}

	result := assertion.Evaluate(obs)
	if !result.Passed {
		t.Errorf("Expected display assertion to pass, got: %s", result.Message)
	}
}

func TestDisplayAssertion_Fail(t *testing.T) {
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals: []core.Signal{
			&core.DisplaySignal{
				Name:       "LCD",
				Text:       "Error",
				Confidence: 0.92,
			},
		},
	}

	assertion := &DisplayAssertion{
		Name:     "LCD",
		Expected: "Ready",
	}

	result := assertion.Evaluate(obs)
	if result.Passed {
		t.Error("Expected display assertion to fail for text mismatch")
	}
}

func TestDisplayAssertion_NotFound(t *testing.T) {
	obs := &core.Observation{
		ID:       "test-1",
		DeviceID: "test-device",
		Signals:  []core.Signal{},
	}

	assertion := &DisplayAssertion{
		Name:     "LCD",
		Expected: "Ready",
	}

	result := assertion.Evaluate(obs)
	if result.Passed {
		t.Error("Expected display assertion to fail when display not found")
	}
}

func TestDisplayAssertion_String(t *testing.T) {
	assertion := &DisplayAssertion{
		Name:     "LCD",
		Expected: "Hello World",
	}

	expected := `Display.LCD "Hello World"`
	if assertion.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, assertion.String())
	}
}

// Parser Tests

func TestParse_LED_ON(t *testing.T) {
	assertion, err := Parse("LED.POWER ON")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	ledAssert, ok := assertion.(*LEDAssertion)
	if !ok {
		t.Fatalf("Expected *LEDAssertion, got %T", assertion)
	}

	if ledAssert.Name != "POWER" {
		t.Errorf("Expected name 'POWER', got '%s'", ledAssert.Name)
	}

	if ledAssert.Expected.On == nil || !*ledAssert.Expected.On {
		t.Error("Expected On state to be true")
	}
}

func TestParse_LED_OFF(t *testing.T) {
	assertion, err := Parse("LED.ERROR OFF")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	ledAssert, ok := assertion.(*LEDAssertion)
	if !ok {
		t.Fatalf("Expected *LEDAssertion, got %T", assertion)
	}

	if ledAssert.Expected.On == nil || *ledAssert.Expected.On {
		t.Error("Expected On state to be false")
	}
}

func TestParse_LED_Blink(t *testing.T) {
	assertion, err := Parse("LED.STATUS BLINK 2.5Hz")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	ledAssert, ok := assertion.(*LEDAssertion)
	if !ok {
		t.Fatalf("Expected *LEDAssertion, got %T", assertion)
	}

	if ledAssert.Expected.BlinkHz == nil || *ledAssert.Expected.BlinkHz != 2.5 {
		t.Errorf("Expected blink rate 2.5, got %v", ledAssert.Expected.BlinkHz)
	}
}

func TestParse_LED_BlinkCaseInsensitive(t *testing.T) {
	assertion, err := Parse("LED.STATUS blink 1.0 Hz")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	ledAssert, ok := assertion.(*LEDAssertion)
	if !ok {
		t.Fatalf("Expected *LEDAssertion, got %T", assertion)
	}

	if ledAssert.Expected.BlinkHz == nil || *ledAssert.Expected.BlinkHz != 1.0 {
		t.Errorf("Expected blink rate 1.0, got %v", ledAssert.Expected.BlinkHz)
	}
}

func TestParse_LED_Color(t *testing.T) {
	assertion, err := Parse("LED.RGB1 COLOR RGB(255, 128, 0)")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	ledAssert, ok := assertion.(*LEDAssertion)
	if !ok {
		t.Fatalf("Expected *LEDAssertion, got %T", assertion)
	}

	if ledAssert.Expected.Color == nil {
		t.Fatal("Expected color to be set")
	}

	color := *ledAssert.Expected.Color
	if color.R != 255 || color.G != 128 || color.B != 0 {
		t.Errorf("Expected RGB(255,128,0), got RGB(%d,%d,%d)", color.R, color.G, color.B)
	}
}

func TestParse_LED_ColorNoSpaces(t *testing.T) {
	assertion, err := Parse("LED.RGB1 COLOR RGB(255,0,128)")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	ledAssert := assertion.(*LEDAssertion)
	color := *ledAssert.Expected.Color

	if color.R != 255 || color.G != 0 || color.B != 128 {
		t.Errorf("Expected RGB(255,0,128), got RGB(%d,%d,%d)", color.R, color.G, color.B)
	}
}

func TestParse_Timing(t *testing.T) {
	assertion, err := Parse("BootTime < 3000ms")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	timingAssert, ok := assertion.(*TimingAssertion)
	if !ok {
		t.Fatalf("Expected *TimingAssertion, got %T", assertion)
	}

	if timingAssert.MaxDurationMs != 3000 {
		t.Errorf("Expected max duration 3000, got %d", timingAssert.MaxDurationMs)
	}
}

func TestParse_TimingWithSpaces(t *testing.T) {
	assertion, err := Parse("BootTime < 5000 ms")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	timingAssert := assertion.(*TimingAssertion)
	if timingAssert.MaxDurationMs != 5000 {
		t.Errorf("Expected max duration 5000, got %d", timingAssert.MaxDurationMs)
	}
}

func TestParse_InvalidAssertion(t *testing.T) {
	_, err := Parse("Invalid assertion")
	if err == nil {
		t.Error("Expected error for invalid assertion")
	}
}

func TestParse_InvalidLEDState(t *testing.T) {
	_, err := Parse("LED.POWER INVALID_STATE")
	if err == nil {
		t.Error("Expected error for invalid LED state")
	}
}

func TestParse_InvalidTimingSyntax(t *testing.T) {
	_, err := Parse("BootTime 3000")
	if err == nil {
		t.Error("Expected error for invalid timing syntax")
	}
}
