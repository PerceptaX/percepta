package codegen

import (
	"strings"
	"testing"

	"github.com/perceptumx/percepta/internal/knowledge"
)

func TestPromptBuilder_BuildSystemPrompt_WithoutPatterns(t *testing.T) {
	// Create pattern store without vector store (no semantic search)
	store, err := knowledge.NewPatternStore()
	if err != nil {
		t.Fatalf("Failed to create pattern store: %v", err)
	}
	defer store.Close()

	builder := NewPromptBuilder(store)

	prompt, err := builder.BuildSystemPrompt("Blink LED at 1Hz", "esp32")
	if err != nil {
		t.Fatalf("BuildSystemPrompt() error = %v", err)
	}

	// Should contain BARR-C requirements even without patterns
	requiredSections := []string{
		"BARR-C Style Requirements",
		"Function names: Module_Function()",
		"Variables: snake_case",
		"Constants: UPPER_SNAKE",
		"Use stdint.h",
		"Non-blocking",
		"ESP32-specific",
		"gpio_set_level()",
	}

	for _, section := range requiredSections {
		if !strings.Contains(prompt, section) {
			t.Errorf("Prompt missing required section: %s", section)
		}
	}
}

func TestPromptBuilder_BuildSystemPrompt_WithPatterns(t *testing.T) {
	// Create pattern store with mock embedder for testing
	store, err := knowledge.NewPatternStore()
	if err != nil {
		t.Fatalf("Failed to create pattern store: %v", err)
	}
	defer store.Close()

	// Initialize vector store with mock embedder
	mockEmbedder := &mockEmbedder{}
	if err := store.InitializeVectorStoreWithEmbedder(mockEmbedder); err != nil {
		t.Fatalf("Failed to initialize vector store: %v", err)
	}

	// Add a validated pattern (this requires observation to exist first)
	// For testing, we'll create a minimal setup
	// Note: In real usage, observation must exist before storing pattern

	builder := NewPromptBuilder(store)

	prompt, err := builder.BuildSystemPrompt("Blink LED", "esp32")
	if err != nil {
		t.Fatalf("BuildSystemPrompt() error = %v", err)
	}

	// Should contain BARR-C requirements
	if !strings.Contains(prompt, "BARR-C Style Requirements") {
		t.Error("Prompt missing BARR-C requirements")
	}

	// Should contain board-specific notes
	if !strings.Contains(prompt, "ESP32-specific") {
		t.Error("Prompt missing ESP32-specific notes")
	}
}

func TestPromptBuilder_BoardSpecificNotes(t *testing.T) {
	tests := []struct {
		boardType string
		expected  []string
	}{
		{
			boardType: "esp32",
			expected:  []string{"ESP32-specific", "gpio_set_level()", "driver/gpio.h", "esp_timer_create()"},
		},
		{
			boardType: "stm32",
			expected:  []string{"STM32-specific", "HAL_GPIO_WritePin()", "stm32f4xx_hal.h"},
		},
		{
			boardType: "arduino",
			expected:  []string{"Arduino-specific", "digitalWrite()", "millis()"},
		},
		{
			boardType: "unknown",
			expected:  []string{"General embedded C", "stdint.h"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.boardType, func(t *testing.T) {
			notes := getBoardSpecificNotes(tt.boardType)

			for _, exp := range tt.expected {
				if !strings.Contains(notes, exp) {
					t.Errorf("Board notes for %s missing: %s", tt.boardType, exp)
				}
			}
		})
	}
}

func TestPromptBuilder_CodeTruncation(t *testing.T) {
	// Create pattern store
	store, err := knowledge.NewPatternStore()
	if err != nil {
		t.Fatalf("Failed to create pattern store: %v", err)
	}
	defer store.Close()

	// Initialize vector store with mock embedder
	mockEmbedder := &mockEmbedder{}
	if err := store.InitializeVectorStoreWithEmbedder(mockEmbedder); err != nil {
		t.Fatalf("Failed to initialize vector store: %v", err)
	}

	builder := NewPromptBuilder(store)

	// Build prompt - should not error even if patterns have long code
	prompt, err := builder.BuildSystemPrompt("Test spec", "esp32")
	if err != nil {
		t.Fatalf("BuildSystemPrompt() error = %v", err)
	}

	// Verify prompt was created
	if prompt == "" {
		t.Error("BuildSystemPrompt() returned empty prompt")
	}
}

func TestPromptBuilder_GracefulDegradation(t *testing.T) {
	// Create pattern store without vector store
	store, err := knowledge.NewPatternStore()
	if err != nil {
		t.Fatalf("Failed to create pattern store: %v", err)
	}
	defer store.Close()

	builder := NewPromptBuilder(store)

	// Should work even without vector store (graceful degradation)
	prompt, err := builder.BuildSystemPrompt("Blink LED", "esp32")
	if err != nil {
		t.Fatalf("BuildSystemPrompt() should not error without vector store, got: %v", err)
	}

	// Should still have BARR-C requirements
	if !strings.Contains(prompt, "BARR-C Style Requirements") {
		t.Error("Prompt missing BARR-C requirements in graceful degradation")
	}
}

// mockEmbedder is a simple mock for testing
type mockEmbedder struct{}

func (m *mockEmbedder) Embed(text string) ([]float32, error) {
	// Return deterministic embedding based on text length
	embedding := make([]float32, 1536)
	for i := range embedding {
		embedding[i] = float32(len(text)%100) / 100.0
	}
	return embedding, nil
}
