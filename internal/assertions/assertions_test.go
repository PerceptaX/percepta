package assertions

import (
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
