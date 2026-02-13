package vision

import (
	"testing"
	"time"

	"github.com/perceptumx/percepta/internal/core"
)

func TestLEDAggregator_SteadyOn(t *testing.T) {
	agg := &ledAggregator{
		name: "LED1",
		observations: []core.LEDSignal{
			{Name: "LED1", On: true, Confidence: 0.9},
			{Name: "LED1", On: true, Confidence: 0.92},
			{Name: "LED1", On: true, Confidence: 0.88},
			{Name: "LED1", On: true, Confidence: 0.91},
			{Name: "LED1", On: true, Confidence: 0.89},
		},
	}

	result := agg.aggregate()

	if result.Name != "LED1" {
		t.Errorf("expected LED1, got %s", result.Name)
	}
	if !result.On {
		t.Errorf("expected LED to be on")
	}
	if result.BlinkHz != 0 {
		t.Errorf("expected BlinkHz 0, got %f", result.BlinkHz)
	}

	// Confidence should be average
	expectedConf := (0.9 + 0.92 + 0.88 + 0.91 + 0.89) / 5.0
	if result.Confidence != expectedConf {
		t.Errorf("expected confidence %f, got %f", expectedConf, result.Confidence)
	}
}

func TestLEDAggregator_SteadyOff(t *testing.T) {
	agg := &ledAggregator{
		name: "LED2",
		observations: []core.LEDSignal{
			{Name: "LED2", On: false, Confidence: 0.85},
			{Name: "LED2", On: false, Confidence: 0.87},
			{Name: "LED2", On: false, Confidence: 0.86},
		},
	}

	result := agg.aggregate()

	if result.On {
		t.Errorf("expected LED to be off")
	}
	if result.BlinkHz != 0 {
		t.Errorf("expected BlinkHz 0, got %f", result.BlinkHz)
	}
}

func TestLEDAggregator_Blinking(t *testing.T) {
	agg := &ledAggregator{
		name: "LED3",
		observations: []core.LEDSignal{
			{Name: "LED3", On: true, Confidence: 0.9},
			{Name: "LED3", On: false, Confidence: 0.88},
			{Name: "LED3", On: true, Confidence: 0.91},
			{Name: "LED3", On: false, Confidence: 0.89},
			{Name: "LED3", On: true, Confidence: 0.90},
		},
	}

	result := agg.aggregate()

	if !result.On {
		t.Errorf("expected blinking LED to be logically on")
	}

	// With 4 transitions (on->off, off->on, on->off, off->on), frequency = 4/2 = 2 Hz
	if result.BlinkHz != 2.0 {
		t.Errorf("expected BlinkHz 2.0, got %f", result.BlinkHz)
	}
}

func TestAggregateDisplays_StaticText(t *testing.T) {
	base := time.Now()
	frames := []FrameResult{
		{
			Signals:    []core.Signal{core.DisplaySignal{Name: "LCD", Text: "Ready", Confidence: 0.90}},
			CapturedAt: base,
		},
		{
			Signals:    []core.Signal{core.DisplaySignal{Name: "LCD", Text: "Ready", Confidence: 0.92}},
			CapturedAt: base.Add(200 * time.Millisecond),
		},
		{
			Signals:    []core.Signal{core.DisplaySignal{Name: "LCD", Text: "Ready", Confidence: 0.88}},
			CapturedAt: base.Add(400 * time.Millisecond),
		},
	}

	displays := AggregateDisplays(frames)

	if len(displays) != 1 {
		t.Fatalf("expected 1 display, got %d", len(displays))
	}

	d := displays[0]
	if d.Name != "LCD" {
		t.Errorf("expected name LCD, got %s", d.Name)
	}
	if d.Text != "Ready" {
		t.Errorf("expected text 'Ready', got '%s'", d.Text)
	}
	if d.Changed {
		t.Errorf("expected Changed=false for static text")
	}
	if len(d.History) != 0 {
		t.Errorf("expected nil History for static text, got %d entries", len(d.History))
	}
}

func TestAggregateDisplays_ChangingText(t *testing.T) {
	base := time.Now()
	frames := []FrameResult{
		{
			Signals:    []core.Signal{core.DisplaySignal{Name: "LCD", Text: "Boot", Confidence: 0.90}},
			CapturedAt: base,
		},
		{
			Signals:    []core.Signal{core.DisplaySignal{Name: "LCD", Text: "Init...", Confidence: 0.88}},
			CapturedAt: base.Add(400 * time.Millisecond),
		},
		{
			Signals:    []core.Signal{core.DisplaySignal{Name: "LCD", Text: "Ready", Confidence: 0.90}},
			CapturedAt: base.Add(800 * time.Millisecond),
		},
	}

	displays := AggregateDisplays(frames)

	if len(displays) != 1 {
		t.Fatalf("expected 1 display, got %d", len(displays))
	}

	d := displays[0]
	if !d.Changed {
		t.Errorf("expected Changed=true for changing text")
	}
	if d.Text != "Ready" {
		t.Errorf("expected latest text 'Ready', got '%s'", d.Text)
	}
	if len(d.History) != 3 {
		t.Fatalf("expected 3 history entries, got %d", len(d.History))
	}
	if d.History[0].Text != "Boot" {
		t.Errorf("expected history[0] text 'Boot', got '%s'", d.History[0].Text)
	}
	if d.History[1].Text != "Init..." {
		t.Errorf("expected history[1] text 'Init...', got '%s'", d.History[1].Text)
	}
	if d.History[2].Text != "Ready" {
		t.Errorf("expected history[2] text 'Ready', got '%s'", d.History[2].Text)
	}
}

func TestAggregateDisplays_Deduplication(t *testing.T) {
	base := time.Now()
	frames := []FrameResult{
		{
			Signals:    []core.Signal{core.DisplaySignal{Name: "LCD", Text: "Boot", Confidence: 0.90}},
			CapturedAt: base,
		},
		{
			Signals:    []core.Signal{core.DisplaySignal{Name: "LCD", Text: "Boot", Confidence: 0.88}},
			CapturedAt: base.Add(200 * time.Millisecond),
		},
		{
			Signals:    []core.Signal{core.DisplaySignal{Name: "LCD", Text: "Ready", Confidence: 0.90}},
			CapturedAt: base.Add(400 * time.Millisecond),
		},
		{
			Signals:    []core.Signal{core.DisplaySignal{Name: "LCD", Text: "Ready", Confidence: 0.91}},
			CapturedAt: base.Add(600 * time.Millisecond),
		},
	}

	displays := AggregateDisplays(frames)

	if len(displays) != 1 {
		t.Fatalf("expected 1 display, got %d", len(displays))
	}

	d := displays[0]
	if !d.Changed {
		t.Errorf("expected Changed=true")
	}
	// Should deduplicate consecutive identical text: Boot, Ready (not Boot, Boot, Ready, Ready)
	if len(d.History) != 2 {
		t.Fatalf("expected 2 history entries (deduplicated), got %d", len(d.History))
	}
	if d.History[0].Text != "Boot" {
		t.Errorf("expected history[0] 'Boot', got '%s'", d.History[0].Text)
	}
	if d.History[1].Text != "Ready" {
		t.Errorf("expected history[1] 'Ready', got '%s'", d.History[1].Text)
	}
}

func TestAggregateDisplays_MultipleDisplays(t *testing.T) {
	base := time.Now()
	frames := []FrameResult{
		{
			Signals: []core.Signal{
				core.DisplaySignal{Name: "LCD", Text: "Boot", Confidence: 0.90},
				core.DisplaySignal{Name: "OLED", Text: "Status: OK", Confidence: 0.85},
			},
			CapturedAt: base,
		},
		{
			Signals: []core.Signal{
				core.DisplaySignal{Name: "LCD", Text: "Ready", Confidence: 0.88},
				core.DisplaySignal{Name: "OLED", Text: "Status: OK", Confidence: 0.87},
			},
			CapturedAt: base.Add(200 * time.Millisecond),
		},
	}

	displays := AggregateDisplays(frames)

	if len(displays) != 2 {
		t.Fatalf("expected 2 displays, got %d", len(displays))
	}

	var lcd, oled *core.DisplaySignal
	for i := range displays {
		if displays[i].Name == "LCD" {
			lcd = &displays[i]
		} else if displays[i].Name == "OLED" {
			oled = &displays[i]
		}
	}

	if lcd == nil || oled == nil {
		t.Fatal("expected both LCD and OLED displays")
	}

	if !lcd.Changed {
		t.Errorf("LCD should have Changed=true")
	}
	if oled.Changed {
		t.Errorf("OLED should have Changed=false (static)")
	}
}

func TestAggregateDisplays_EmptyFrames(t *testing.T) {
	displays := AggregateDisplays(nil)
	if displays != nil {
		t.Errorf("expected nil for empty frames, got %d", len(displays))
	}

	displays = AggregateDisplays([]FrameResult{})
	if displays != nil {
		t.Errorf("expected nil for empty frames, got %d", len(displays))
	}
}

func TestAggregateDisplays_AverageConfidence(t *testing.T) {
	base := time.Now()
	frames := []FrameResult{
		{
			Signals:    []core.Signal{core.DisplaySignal{Name: "LCD", Text: "Hello", Confidence: 0.80}},
			CapturedAt: base,
		},
		{
			Signals:    []core.Signal{core.DisplaySignal{Name: "LCD", Text: "Hello", Confidence: 0.90}},
			CapturedAt: base.Add(200 * time.Millisecond),
		},
		{
			Signals:    []core.Signal{core.DisplaySignal{Name: "LCD", Text: "Hello", Confidence: 1.00}},
			CapturedAt: base.Add(400 * time.Millisecond),
		},
	}

	displays := AggregateDisplays(frames)
	if len(displays) != 1 {
		t.Fatalf("expected 1 display, got %d", len(displays))
	}

	expectedConf := (0.80 + 0.90 + 1.00) / 3.0
	if displays[0].Confidence != expectedConf {
		t.Errorf("expected confidence %f, got %f", expectedConf, displays[0].Confidence)
	}
}

func TestAggregateLEDs(t *testing.T) {
	frames := []FrameResult{
		{
			Signals: []core.Signal{
				core.LEDSignal{Name: "LED1", On: true, Confidence: 0.9},
				core.LEDSignal{Name: "LED2", On: false, Confidence: 0.85},
			},
			CapturedAt: time.Now(),
		},
		{
			Signals: []core.Signal{
				core.LEDSignal{Name: "LED1", On: true, Confidence: 0.92},
				core.LEDSignal{Name: "LED2", On: true, Confidence: 0.87},
			},
			CapturedAt: time.Now(),
		},
		{
			Signals: []core.Signal{
				core.LEDSignal{Name: "LED1", On: true, Confidence: 0.88},
				core.LEDSignal{Name: "LED2", On: false, Confidence: 0.86},
			},
			CapturedAt: time.Now(),
		},
	}

	leds := AggregateLEDs(frames)

	if len(leds) != 2 {
		t.Fatalf("expected 2 LEDs, got %d", len(leds))
	}

	// Find LED1 and LED2
	var led1, led2 *core.LEDSignal
	for i := range leds {
		if leds[i].Name == "LED1" {
			led1 = &leds[i]
		} else if leds[i].Name == "LED2" {
			led2 = &leds[i]
		}
	}

	if led1 == nil || led2 == nil {
		t.Fatalf("expected to find LED1 and LED2")
	}

	// LED1 should be steady on
	if !led1.On {
		t.Errorf("LED1 should be on")
	}
	if led1.BlinkHz != 0 {
		t.Errorf("LED1 should have BlinkHz 0, got %f", led1.BlinkHz)
	}

	// LED2 should be blinking (on in 1 frame, off in 2)
	if !led2.On {
		t.Errorf("LED2 should be logically on (blinking)")
	}
	// 2 transitions: off->on, on->off
	if led2.BlinkHz != 1.0 {
		t.Errorf("LED2 should have BlinkHz 1.0, got %f", led2.BlinkHz)
	}
}
