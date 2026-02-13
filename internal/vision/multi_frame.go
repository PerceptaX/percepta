package vision

import (
	"fmt"
	"time"

	"github.com/perceptumx/percepta/internal/core"
)

// MultiFrameCapture captures multiple frames and aggregates LED detections
type MultiFrameCapture struct {
	camera     core.CameraDriver
	parser     SignalParser
	frameCount int           // Number of frames to capture
	interval   time.Duration // Time between frames
}

func NewMultiFrameCapture(camera core.CameraDriver, parser SignalParser) *MultiFrameCapture {
	return &MultiFrameCapture{
		camera:     camera,
		parser:     parser,
		frameCount: 5,                      // Capture 5 frames
		interval:   200 * time.Millisecond, // 200ms apart (1 second total)
	}
}

func NewMultiFrameCaptureWithOptions(camera core.CameraDriver, parser SignalParser, frameCount int, interval time.Duration) *MultiFrameCapture {
	return &MultiFrameCapture{
		camera:     camera,
		parser:     parser,
		frameCount: frameCount,
		interval:   interval,
	}
}

type FrameResult struct {
	Signals    []core.Signal
	CapturedAt time.Time
}

func (m *MultiFrameCapture) Capture() ([]FrameResult, error) {
	var results []FrameResult

	for i := 0; i < m.frameCount; i++ {
		// Capture frame
		frame, err := m.camera.CaptureFrame()
		if err != nil {
			return nil, fmt.Errorf("frame %d capture failed: %w", i, err)
		}

		// Parse signals
		signals, err := m.parser.Parse(frame)
		if err != nil {
			// Log but continue with other frames
			continue
		}

		results = append(results, FrameResult{
			Signals:    signals,
			CapturedAt: time.Now(),
		})

		// Wait before next frame (except last)
		if i < m.frameCount-1 {
			time.Sleep(m.interval)
		}
	}

	return results, nil
}

// AggregateLEDs combines LED detections across frames
func AggregateLEDs(frames []FrameResult) []core.LEDSignal {
	// Map LED name → aggregated state
	ledMap := make(map[string]*ledAggregator)
	calibrator := NewConfidenceCalibrator()

	for _, frame := range frames {
		for _, signal := range frame.Signals {
			if led, ok := signal.(core.LEDSignal); ok {
				if agg, exists := ledMap[led.Name]; exists {
					agg.addObservation(led)
				} else {
					ledMap[led.Name] = &ledAggregator{
						name:         led.Name,
						observations: []core.LEDSignal{led},
					}
				}
			}
		}
	}

	var leds []core.LEDSignal
	for _, agg := range ledMap {
		led := agg.aggregate()

		// Calibrate confidence based on detection rate
		detectionRate := float64(len(agg.observations)) / float64(len(frames))
		led = calibrator.CalibrateLED(led, detectionRate)

		leds = append(leds, led)
	}

	return leds
}

// AggregateDisplays combines display detections across frames, tracking text changes
func AggregateDisplays(frames []FrameResult) []core.DisplaySignal {
	if len(frames) == 0 {
		return nil
	}

	type displayObs struct {
		text       string
		confidence float64
		offsetMs   int64
	}

	displayMap := make(map[string][]displayObs)
	var baseTime time.Time
	if len(frames) > 0 {
		baseTime = frames[0].CapturedAt
	}

	for _, frame := range frames {
		offsetMs := frame.CapturedAt.Sub(baseTime).Milliseconds()
		for _, signal := range frame.Signals {
			if d, ok := signal.(core.DisplaySignal); ok {
				displayMap[d.Name] = append(displayMap[d.Name], displayObs{
					text:       d.Text,
					confidence: d.Confidence,
					offsetMs:   offsetMs,
				})
			}
		}
	}

	var displays []core.DisplaySignal
	for name, observations := range displayMap {
		if len(observations) == 0 {
			continue
		}

		// Deduplicate: only store transitions (consecutive different text)
		var transitions []core.DisplayTextEntry
		transitions = append(transitions, core.DisplayTextEntry{
			OffsetMs:   observations[0].offsetMs,
			Text:       observations[0].text,
			Confidence: observations[0].confidence,
		})
		for i := 1; i < len(observations); i++ {
			if observations[i].text != observations[i-1].text {
				transitions = append(transitions, core.DisplayTextEntry{
					OffsetMs:   observations[i].offsetMs,
					Text:       observations[i].text,
					Confidence: observations[i].confidence,
				})
			}
		}

		// Average confidence across all observations
		totalConf := 0.0
		for _, obs := range observations {
			totalConf += obs.confidence
		}
		avgConf := totalConf / float64(len(observations))

		// Latest text for backward compat
		latest := observations[len(observations)-1]

		changed := len(transitions) > 1
		display := core.DisplaySignal{
			Name:       name,
			Text:       latest.text,
			Confidence: avgConf,
			Changed:    changed,
		}
		if changed {
			display.History = transitions
		}

		displays = append(displays, display)
	}

	return displays
}

type ledAggregator struct {
	name         string
	observations []core.LEDSignal
}

func (a *ledAggregator) addObservation(led core.LEDSignal) {
	a.observations = append(a.observations, led)
}

func (a *ledAggregator) aggregate() core.LEDSignal {
	if len(a.observations) == 0 {
		return core.LEDSignal{}
	}

	// Calculate blink frequency from on/off transitions
	onCount := 0
	for _, obs := range a.observations {
		if obs.On {
			onCount++
		}
	}
	offCount := len(a.observations) - onCount

	// Determine state and blink frequency
	led := a.observations[0] // Start with first observation
	led.Name = a.name

	if onCount > 0 && offCount > 0 {
		// Blinking detected (transitions between on/off)
		// Estimate frequency: transitions per second
		// With 5 frames over 1 second, transitions ≈ blink_hz
		transitionCount := 0
		for i := 1; i < len(a.observations); i++ {
			if a.observations[i].On != a.observations[i-1].On {
				transitionCount++
			}
		}
		led.BlinkHz = float64(transitionCount) / 2.0 // Each cycle has 2 transitions
		led.On = true                                // Blinking LED is "on" logically
	} else if onCount == len(a.observations) {
		// Steady on
		led.On = true
		led.BlinkHz = 0
	} else {
		// Steady off
		led.On = false
		led.BlinkHz = 0
	}

	// Aggregate confidence (average)
	totalConf := 0.0
	for _, obs := range a.observations {
		totalConf += obs.Confidence
	}
	led.Confidence = totalConf / float64(len(a.observations))

	return led
}
