package vision

import (
	"context"
	"encoding/base64"
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
	client *anthropic.Client
	parser SignalParser
}

// SignalParser converts vision API response text to structured signals
// Isolated for easy replacement (regex now, structured output later)
type SignalParser interface {
	Parse(responseText string) []core.Signal
}

func NewClaudeVision() (*ClaudeVision, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	return &ClaudeVision{
		client: &client,
		parser: NewRegexParser(), // Swappable parser
	}, nil
}

func (v *ClaudeVision) Observe(deviceID string, frame []byte) (*core.Observation, error) {
	// Encode to base64
	base64Frame := base64.StdEncoding.EncodeToString(frame)

	// Create image block with base64 source
	imageBlock := anthropic.NewImageBlockBase64(
		string(anthropic.Base64ImageSourceMediaTypeImageJPEG),
		base64Frame,
	)

	// Create text block with prompt
	textBlock := anthropic.NewTextBlock(HardwarePrompt)

	// Call Claude Vision API
	message, err := v.client.Messages.New(context.Background(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(imageBlock, textBlock),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("vision API call failed: %w", err)
	}

	// Extract text response
	responseText := ""
	for _, block := range message.Content {
		if block.Type == "text" {
			responseText += block.Text
		}
	}

	// Parse signals using isolated parser
	signals := v.parser.Parse(responseText)

	return &core.Observation{
		ID:        core.GenerateID(),
		DeviceID:  deviceID,
		Timestamp: time.Now(),
		Signals:   signals,
	}, nil
}
