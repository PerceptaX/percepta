package vision

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/perceptumx/percepta/internal/core"
)

// RegexParser uses regex to extract signals from unstructured text
// This is a fallback implementation for when structured output fails
type RegexParser struct{}

func NewRegexParser() *RegexParser {
	return &RegexParser{}
}

func (p *RegexParser) Parse(frame []byte) ([]core.Signal, error) {
	// Call Claude Vision API with text prompt (no tool use)
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))

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
	message, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
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

	// Parse text with regex
	return p.parseText(responseText), nil
}

func (p *RegexParser) parseText(text string) []core.Signal {
	var signals []core.Signal

	// Parse LED signals
	ledRegex := regexp.MustCompile(`(?i)([a-z0-9_-]+)?\s*(LED|led)(?:[^.]*)(on|off|blinking)(?:[^.]*?)(?:(\d+(?:\.\d+)?)\s*Hz)?`)
	ledMatches := ledRegex.FindAllStringSubmatch(text, -1)

	for i, match := range ledMatches {
		// Assign deterministic name by index (object permanence)
		name := fmt.Sprintf("LED%d", i+1)
		stateStr := strings.ToLower(match[3])

		led := core.LEDSignal{
			Name:       name,
			On:         stateStr == "on" || stateStr == "blinking",
			Confidence: 0.85,
		}

		if match[4] != "" {
			if freq, err := strconv.ParseFloat(match[4], 64); err == nil {
				led.BlinkHz = freq
			}
		}

		// Extract color from matched LED segment only (not entire response)
		// This prevents false positives on multi-LED boards
		segment := strings.ToLower(match[0])
		colors := []string{"red", "green", "blue", "yellow", "white", "orange"}
		for _, color := range colors {
			if strings.Contains(segment, color) {
				switch color {
				case "red":
					led.Color = core.RGB{R: 255, G: 0, B: 0}
				case "green":
					led.Color = core.RGB{R: 0, G: 255, B: 0}
				case "blue":
					led.Color = core.RGB{R: 0, G: 0, B: 255}
				case "yellow":
					led.Color = core.RGB{R: 255, G: 255, B: 0}
				case "white":
					led.Color = core.RGB{R: 255, G: 255, B: 255}
				case "orange":
					led.Color = core.RGB{R: 255, G: 165, B: 0}
				}
				break
			}
		}

		signals = append(signals, led)
	}

	// Parse Display signals
	displayRegex := regexp.MustCompile(`(?i)(OLED|LCD|Display)(?:[^\"]*)\"([^\"]+)\"`)
	displayMatches := displayRegex.FindAllStringSubmatch(text, -1)

	for _, match := range displayMatches {
		signals = append(signals, core.DisplaySignal{
			Name:       match[1],
			Text:       match[2],
			Confidence: 0.90,
		})
	}

	return signals
}
