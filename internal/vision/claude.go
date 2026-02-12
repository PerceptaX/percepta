package vision

import (
	"fmt"
	"os"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/perceptumx/percepta/internal/core"
)

const HardwarePrompt = `Describe this embedded hardware device precisely.

Focus on:
1. LED states (on/off, color, blinking frequency in Hz)
2. Display content (transcribe ALL visible text exactly)
3. Boot sequence indicators if visible

Format your response as:
LEDs:
- [name/color]: [on/off], [color if visible], [frequency in Hz if blinking]

Displays:
- [type]: "[exact text shown]"

Boot: [duration in seconds if measurable]

Be precise with measurements. Estimate blink frequency in Hz.`

// ClaudeVision implements core.VisionDriver using Claude Sonnet 4.5
type ClaudeVision struct {
	client           *anthropic.Client
	structuredParser SignalParser
	regexParser      SignalParser
}

// SignalParser converts vision API response to structured signals
// Isolated for easy replacement (regex now, structured output later)
type SignalParser interface {
	Parse(frame []byte) ([]core.Signal, error)
}

func NewClaudeVision() (*ClaudeVision, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	return &ClaudeVision{
		client:           &client,
		structuredParser: NewStructuredParser(&client), // Primary: structured output
		regexParser:      NewRegexParser(),             // Fallback: regex parsing
	}, nil
}

func (v *ClaudeVision) Observe(deviceID string, frame []byte) (*core.Observation, error) {
	// Try structured parser first (tool use for deterministic extraction)
	signals, err := v.structuredParser.Parse(frame)
	if err == nil && len(signals) > 0 {
		return &core.Observation{
			ID:        core.GenerateID(),
			DeviceID:  deviceID,
			Timestamp: time.Now(),
			Signals:   signals,
		}, nil
	}

	// Fallback to regex parser for robustness
	// This handles cases where tool use fails or returns no signals
	signals, err = v.regexParser.Parse(frame)
	if err != nil {
		return nil, fmt.Errorf("both parsers failed: %w", err)
	}

	return &core.Observation{
		ID:        core.GenerateID(),
		DeviceID:  deviceID,
		Timestamp: time.Now(),
		Signals:   signals,
	}, nil
}

// GetParser returns a parser that tries structured first, then falls back to regex
func (v *ClaudeVision) GetParser() SignalParser {
	return &fallbackParser{
		primary:  v.structuredParser,
		fallback: v.regexParser,
	}
}

// fallbackParser tries primary parser first, then falls back to secondary
type fallbackParser struct {
	primary  SignalParser
	fallback SignalParser
}

func (p *fallbackParser) Parse(frame []byte) ([]core.Signal, error) {
	signals, err := p.primary.Parse(frame)
	if err == nil && len(signals) > 0 {
		return signals, nil
	}

	return p.fallback.Parse(frame)
}
