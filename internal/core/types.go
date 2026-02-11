package core

import "time"

// Signal interface - LED, Display, or Boot timing
type Signal interface {
	Type() string
	State() interface{}
}

// LEDSignal represents LED state observation
type LEDSignal struct {
	Name       string  `json:"name"`
	On         bool    `json:"on"`
	Color      RGB     `json:"color,omitempty"`
	Brightness uint8   `json:"brightness,omitempty"`
	BlinkHz    float64 `json:"blink_hz,omitempty"`
	Confidence float64 `json:"confidence"`
}

func (l LEDSignal) Type() string        { return "led" }
func (l LEDSignal) State() interface{} { return l }

// RGB color
type RGB struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

// DisplaySignal represents display content observation
type DisplaySignal struct {
	Name       string  `json:"name"`
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
}

func (d DisplaySignal) Type() string        { return "display" }
func (d DisplaySignal) State() interface{} { return d }

// BootTimingSignal represents boot sequence timing
type BootTimingSignal struct {
	DurationMs int64   `json:"duration_ms"`
	Confidence float64 `json:"confidence"`
}

func (b BootTimingSignal) Type() string        { return "boot_timing" }
func (b BootTimingSignal) State() interface{} { return b }

// Observation is a snapshot of hardware state at a point in time
type Observation struct {
	ID           string    `json:"id"`
	DeviceID     string    `json:"device_id"`
	FirmwareHash string    `json:"firmware_hash,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	Signals      []Signal  `json:"signals"`
}
