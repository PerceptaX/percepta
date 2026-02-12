package codegen

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// ClaudeClient wraps the Anthropic API for code generation
type ClaudeClient struct {
	apiKey string
	model  string
	client anthropic.Client
}

// NewClaudeClient creates a new Claude API client
// API key is read from ANTHROPIC_API_KEY environment variable if not provided
func NewClaudeClient(apiKey string) *ClaudeClient {
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}

	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &ClaudeClient{
		apiKey: apiKey,
		model:  "claude-sonnet-4-5-20250929", // Latest Claude Sonnet 4.5 model
		client: client,
	}
}

// GenerateCode generates firmware code using Claude API
// spec: Natural language specification (e.g., "Blink LED at 1Hz")
// boardType: Board type (e.g., "esp32", "stm32")
// systemPrompt: Context including BARR-C requirements and validated patterns
// maxTokens: Maximum tokens for response (typically 4096 for code generation)
func (c *ClaudeClient) GenerateCode(
	spec string,
	boardType string,
	systemPrompt string,
	maxTokens int,
) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	if maxTokens <= 0 {
		maxTokens = 4096 // Default for code generation
	}

	// Build user message
	userMessage := fmt.Sprintf(`Generate firmware for %s board:

Specification: %s

Requirements:
- BARR-C compliant
- Use validated patterns provided in system prompt
- Non-blocking architecture (timers, not delays)
- Proper error handling
- Static allocation only

Output only the C source code, no explanations.`, boardType, spec)

	// Call Claude API
	ctx := context.Background()
	response, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(c.model),
		MaxTokens: int64(maxTokens),
		System: []anthropic.TextBlockParam{
			{
				Text: systemPrompt,
				Type: "text",
			},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userMessage)),
		},
		Temperature: anthropic.Float(0.3), // Lower temperature for more deterministic code
	})

	if err != nil {
		return "", fmt.Errorf("API call failed: %w", err)
	}

	// Extract code from response
	if len(response.Content) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	// Get text content from first content block
	contentBlock := response.Content[0]
	// Check if it's a text block
	if contentBlock.Type != "text" {
		return "", fmt.Errorf("unexpected response format: expected text block, got %s", contentBlock.Type)
	}

	text := contentBlock.Text

	// Extract code from markdown code blocks if present
	code := extractCode(text)

	return code, nil
}

// extractCode extracts C source code from text, removing markdown code blocks
func extractCode(text string) string {
	// Remove markdown code blocks (```c ... ``` or ``` ... ```)
	codeBlockRegex := regexp.MustCompile("(?s)```(?:c)?\\n?(.+?)```")
	matches := codeBlockRegex.FindStringSubmatch(text)

	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	// If no code blocks, return trimmed text
	return strings.TrimSpace(text)
}
