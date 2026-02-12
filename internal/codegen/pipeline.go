package codegen

import (
	"fmt"
	"log"

	"github.com/perceptumx/percepta/internal/knowledge"
	"github.com/perceptumx/percepta/internal/style"
)

// GenerationPipeline combines code generation with validation and storage
type GenerationPipeline struct {
	claudeClient  *ClaudeClient
	promptBuilder *PromptBuilder
	styleChecker  *style.StyleChecker
	styleFixer    *style.StyleFixer
	patternStore  *knowledge.PatternStore
}

// GenerationResult contains the results of generation and validation
type GenerationResult struct {
	Code           string
	StyleCompliant bool
	Violations     []style.Violation
	AutoFixed      bool
	PatternStored  bool
	IterationsUsed int
}

// NewGenerationPipeline creates a new generation pipeline
func NewGenerationPipeline(
	claudeClient *ClaudeClient,
	promptBuilder *PromptBuilder,
	styleChecker *style.StyleChecker,
	styleFixer *style.StyleFixer,
	patternStore *knowledge.PatternStore,
) *GenerationPipeline {
	return &GenerationPipeline{
		claudeClient:  claudeClient,
		promptBuilder: promptBuilder,
		styleChecker:  styleChecker,
		styleFixer:    styleFixer,
		patternStore:  patternStore,
	}
}

// Generate generates code with validation pipeline
// spec: Natural language specification (e.g., "Blink LED at 1Hz")
// boardType: Board type (e.g., "esp32", "stm32")
// deviceID: Device identifier for pattern storage linkage
func (p *GenerationPipeline) Generate(
	spec string,
	boardType string,
	deviceID string,
) (*GenerationResult, error) {
	result := &GenerationResult{}

	// 1. Build context-rich prompt
	systemPrompt, err := p.promptBuilder.BuildSystemPrompt(spec, boardType)
	if err != nil {
		return nil, fmt.Errorf("prompt building failed: %w", err)
	}

	// 2. Generate code
	code, err := p.claudeClient.GenerateCode(spec, boardType, systemPrompt, 4096)
	if err != nil {
		return nil, fmt.Errorf("code generation failed: %w", err)
	}
	result.Code = code
	result.IterationsUsed = 1

	// 3. Style validation
	violations, err := p.styleChecker.CheckSource([]byte(code), "generated.c")
	if err != nil {
		return nil, fmt.Errorf("style check failed: %w", err)
	}
	result.Violations = violations

	// 4. Auto-fix if violations exist
	if len(violations) > 0 {
		fixed, _ := p.styleFixer.ApplyFixes(violations, []byte(code))
		code = string(fixed)
		result.Code = code
		result.AutoFixed = true

		// Re-check after fixes
		violations, err = p.styleChecker.CheckSource([]byte(code), "generated.c")
		if err != nil {
			// Log but continue - we have the fixed code anyway
			log.Printf("Warning: style re-check failed: %v", err)
		}
		result.Violations = violations
	}

	// 5. Check if fully compliant
	result.StyleCompliant = len(violations) == 0

	// 6. Store pattern if style compliant (MVP: no hardware validation yet)
	if result.StyleCompliant {
		// Create mock observation for MVP (Phase 8 will add real hardware observation)
		// For now, we use a generated firmware tag to link the pattern
		// In Phase 8, this will be replaced with actual hardware validation
		firmware := "generated-v1"
		_, err = p.patternStore.StoreValidatedPattern(spec, code, deviceID, firmware)
		if err != nil {
			// Log but don't fail - storage is nice-to-have
			log.Printf("Warning: failed to store pattern: %v", err)
		} else {
			result.PatternStored = true
		}
	}

	return result, nil
}
