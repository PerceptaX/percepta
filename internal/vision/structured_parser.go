package vision

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/perceptumx/percepta/internal/core"
)

// StructuredParser uses Claude tool use for deterministic signal extraction
type StructuredParser struct {
	client *anthropic.Client
}

func NewStructuredParser(client *anthropic.Client) *StructuredParser {
	return &StructuredParser{client: client}
}

// Tool definitions for signal extraction
func ledDetectionTool() anthropic.ToolParam {
	return anthropic.ToolParam{
		Name:        "report_led_signals",
		Description: anthropic.String("Report all detected LED signals with state, color, and blink frequency"),
		InputSchema: anthropic.ToolInputSchemaParam{
			Type: "object",
			Properties: map[string]interface{}{
				"leds": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name":       map[string]string{"type": "string", "description": "LED identifier (LED1, LED2, etc)"},
							"on":         map[string]string{"type": "boolean", "description": "True if LED is currently on"},
							"color":      map[string]string{"type": "string", "description": "Color name if visible (red/green/blue/yellow/white/orange)"},
							"blink_hz":   map[string]string{"type": "number", "description": "Blink frequency in Hz if blinking, 0 if steady"},
							"confidence": map[string]string{"type": "number", "description": "Confidence 0-1 in detection"},
						},
						"required": []string{"name", "on", "confidence"},
					},
				},
			},
			Required: []string{"leds"},
		},
	}
}

func displayDetectionTool() anthropic.ToolParam {
	return anthropic.ToolParam{
		Name:        "report_display_content",
		Description: anthropic.String("Report all detected display content with exact text"),
		InputSchema: anthropic.ToolInputSchemaParam{
			Type: "object",
			Properties: map[string]interface{}{
				"displays": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"name":       map[string]string{"type": "string", "description": "Display type (OLED/LCD/Display)"},
							"text":       map[string]string{"type": "string", "description": "Exact text shown on display"},
							"confidence": map[string]string{"type": "number", "description": "OCR confidence 0-1"},
						},
						"required": []string{"name", "text", "confidence"},
					},
				},
			},
			Required: []string{"displays"},
		},
	}
}

func (p *StructuredParser) Parse(frame []byte) ([]core.Signal, error) {
	// Encode frame to base64
	base64Frame := base64.StdEncoding.EncodeToString(frame)

	// Create tools
	ledTool := ledDetectionTool()
	displayTool := displayDetectionTool()

	// Create message with tool use
	message, err := p.client.Messages.New(context.Background(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		Tools: []anthropic.ToolUnionParam{
			{OfTool: &ledTool},
			{OfTool: &displayTool},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewImageBlockBase64(
					string(anthropic.Base64ImageSourceMediaTypeImageJPEG),
					base64Frame,
				),
				anthropic.NewTextBlock(`Analyze this embedded hardware device. Use the tools to report:
1. All detected LEDs (use report_led_signals tool)
2. All display content (use report_display_content tool)

Be precise with measurements.`),
			),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	// Extract signals from tool use responses
	var signals []core.Signal

	for _, block := range message.Content {
		if block.Type == "tool_use" {
			switch block.Name {
			case "report_led_signals":
				leds := parseLEDToolResponse(block.Input)
				signals = append(signals, leds...)
			case "report_display_content":
				displays := parseDisplayToolResponse(block.Input)
				signals = append(signals, displays...)
			}
		}
	}

	return signals, nil
}

func parseLEDToolResponse(input interface{}) []core.Signal {
	var signals []core.Signal

	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return signals
	}

	leds, ok := inputMap["leds"].([]interface{})
	if !ok {
		return signals
	}

	for _, ledData := range leds {
		led, ok := ledData.(map[string]interface{})
		if !ok {
			continue
		}

		signal := core.LEDSignal{
			Name:       getString(led, "name"),
			On:         getBool(led, "on"),
			Confidence: getFloat(led, "confidence"),
		}

		if colorStr := getString(led, "color"); colorStr != "" {
			signal.Color = parseColor(colorStr)
		}

		if blinkHz := getFloat(led, "blink_hz"); blinkHz > 0 {
			signal.BlinkHz = blinkHz
		}

		signals = append(signals, signal)
	}

	return signals
}

func parseDisplayToolResponse(input interface{}) []core.Signal {
	var signals []core.Signal

	inputMap, ok := input.(map[string]interface{})
	if !ok {
		return signals
	}

	displays, ok := inputMap["displays"].([]interface{})
	if !ok {
		return signals
	}

	for _, displayData := range displays {
		display, ok := displayData.(map[string]interface{})
		if !ok {
			continue
		}

		signals = append(signals, core.DisplaySignal{
			Name:       getString(display, "name"),
			Text:       getString(display, "text"),
			Confidence: getFloat(display, "confidence"),
		})
	}

	return signals
}

func parseColor(colorStr string) core.RGB {
	switch colorStr {
	case "red":
		return core.RGB{R: 255, G: 0, B: 0}
	case "green":
		return core.RGB{R: 0, G: 255, B: 0}
	case "blue":
		return core.RGB{R: 0, G: 0, B: 255}
	case "yellow":
		return core.RGB{R: 255, G: 255, B: 0}
	case "white":
		return core.RGB{R: 255, G: 255, B: 255}
	case "orange":
		return core.RGB{R: 255, G: 165, B: 0}
	default:
		return core.RGB{}
	}
}

// Helper functions for safe type conversion
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func getFloat(m map[string]interface{}, key string) float64 {
	// Handle both float64 and json.Number
	switch v := m[key].(type) {
	case float64:
		return v
	case json.Number:
		f, err := v.Float64()
		if err != nil {
			return 0.0
		}
		return f
	case int:
		return float64(v)
	case int64:
		return float64(v)
	}
	return 0.0
}
