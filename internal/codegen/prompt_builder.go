package codegen

import (
	"fmt"
	"strings"

	"github.com/perceptumx/percepta/internal/knowledge"
)

// PromptBuilder creates context-rich prompts for code generation
// Queries knowledge graph for similar validated patterns and includes BARR-C requirements
type PromptBuilder struct {
	patternStore *knowledge.PatternStore
}

// NewPromptBuilder creates a new prompt builder
func NewPromptBuilder(patternStore *knowledge.PatternStore) *PromptBuilder {
	return &PromptBuilder{
		patternStore: patternStore,
	}
}

// BuildSystemPrompt creates a system prompt with BARR-C requirements and validated patterns
// spec: Natural language specification (e.g., "Blink LED at 1Hz")
// boardType: Board type for filtering patterns (e.g., "esp32", "stm32")
func (p *PromptBuilder) BuildSystemPrompt(spec string, boardType string) (string, error) {
	// 1. Build BARR-C requirements section
	barrRequirements := `You are an expert embedded firmware engineer writing BARR-C compliant code.

BARR-C Style Requirements:
- Function names: Module_Function() format (e.g., LED_Init, Timer_Start)
- Variables: snake_case (e.g., led_state, timer_count)
- Constants: UPPER_SNAKE with meaningful names (e.g., LED_PIN, BLINK_PERIOD_MS)
- Types: Use stdint.h (uint8_t, uint16_t, uint32_t, not int/char)
- No magic numbers: Define all constants with #define
- Const correctness: const uint8_t* for read-only pointers
- Doxygen comments: /** ... */ for all functions
- Non-blocking: Use timers/interrupts, never blocking delays
- Error handling: Return error codes, check all returns
- Static allocation: No malloc/free
- Explicit casts: (uint16_t)value for type conversions

Code Structure:
1. Includes at top
2. Constants defined
3. Type definitions
4. Static variables
5. Function prototypes
6. Function implementations
7. Main function if required
`

	// 2. Search for similar validated patterns
	var patternExamples string
	patterns, err := p.patternStore.SearchSimilarPatterns(spec, boardType, 3)
	if err != nil {
		// If vector store not initialized or search fails, continue without patterns
		// This is graceful degradation - can still generate code without examples
		patternExamples = ""
	} else if len(patterns) > 0 {
		patternExamples = "\nValidated Patterns (these work on real hardware):\n\n"
		for i, result := range patterns {
			patternExamples += fmt.Sprintf("Example %d (%.0f%% similar, %.0f%% confidence):\n",
				i+1, result.Similarity*100, result.Confidence*100)
			patternExamples += fmt.Sprintf("Spec: %s\n", result.Pattern.Spec)
			patternExamples += fmt.Sprintf("Board: %s\n", result.Pattern.BoardType)

			// Show code preview (limit to reasonable size)
			codePreview := result.Pattern.Code
			lines := strings.Split(codePreview, "\n")
			if len(lines) > 50 {
				// Truncate long code examples
				codePreview = strings.Join(lines[:50], "\n") + "\n// ... (truncated)"
			}

			patternExamples += fmt.Sprintf("Code:\n```c\n%s\n```\n\n", codePreview)

			// Add observation metadata if available
			if result.Observation != nil && len(result.Observation.Signals) > 0 {
				patternExamples += fmt.Sprintf("Hardware validation: %d signals observed\n\n",
					len(result.Observation.Signals))
			}
		}
		patternExamples += "Use these patterns as reference. Adapt them to the current specification.\n"
	}

	// 3. Add board-specific notes
	boardNotes := getBoardSpecificNotes(boardType)

	// 4. Combine all sections
	systemPrompt := barrRequirements + patternExamples + boardNotes

	return systemPrompt, nil
}

// getBoardSpecificNotes returns board-specific API and setup notes
func getBoardSpecificNotes(boardType string) string {
	switch boardType {
	case "esp32":
		return `
ESP32-specific:
- GPIO pins: Use gpio_set_level() from driver/gpio.h
- GPIO config: Use gpio_config_t with gpio_config()
- Timers: Use esp_timer_create() for non-blocking timing
- No Arduino-style digitalWrite()
- Include: #include "driver/gpio.h" and #include "esp_timer.h"
`
	case "stm32":
		return `
STM32-specific:
- GPIO: Use HAL_GPIO_WritePin() and HAL_GPIO_ReadPin()
- GPIO init: Use HAL_GPIO_Init() with GPIO_InitTypeDef
- Timers: Use HAL timer functions (HAL_TIM_Base_Start_IT)
- Include: #include "stm32f4xx_hal.h" (adjust for your STM32 family)
`
	case "arduino":
		return `
Arduino-specific:
- GPIO: Use digitalWrite() and digitalRead()
- Timers: Use millis() for non-blocking timing
- Pin modes: Use pinMode(pin, OUTPUT) or pinMode(pin, INPUT)
- Include: #include <Arduino.h>
`
	default:
		return `
General embedded C:
- Use standard C library functions
- Include stdint.h for fixed-width integer types
- Implement non-blocking timing with software timers
- Follow BARR-C conventions above
`
	}
}
