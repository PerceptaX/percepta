---
phase: 07-code-generation-engine
plan: 01
subsystem: code-generation
tags: [claude-api, anthropic, prompt-engineering, knowledge-graph, barr-c, firmware-generation, cli, cobra]

# Dependency graph
requires:
  - phase: 06-knowledge-graphs
    provides: PatternStore API, SearchSimilarPatterns, semantic search, validated patterns
  - phase: 05-style-infrastructure
    provides: StyleChecker for BARR-C validation, style requirements
  - phase: 04-polish-alpha
    provides: Cobra CLI framework
provides:
  - Claude API client for code generation (ClaudeClient)
  - Pattern-based prompt engineering (PromptBuilder)
  - percepta generate CLI command
  - Context-rich prompts with BARR-C requirements and validated patterns
  - Board-specific API guidance (ESP32, STM32, Arduino)
affects: [07-02-validation-pipeline, 08-public-launch]

# Tech tracking
tech-stack:
  added: [anthropic-sdk-go, Claude API integration]
  patterns: [Prompt builder with pattern context, graceful degradation without vector store, board-specific API guidance]

key-files:
  created:
    - internal/codegen/claude_client.go
    - internal/codegen/claude_client_test.go
    - internal/codegen/prompt_builder.go
    - internal/codegen/prompt_builder_test.go
    - cmd/percepta/generate.go
  modified:
    - cmd/percepta/main.go

key-decisions:
  - "Use Anthropic SDK directly (already in dependencies) instead of custom HTTP client"
  - "Model: claude-sonnet-4-5-20250929 (latest Claude Sonnet 4.5)"
  - "Temperature: 0.3 for deterministic code generation (vs 1.0 creative)"
  - "Max tokens: 4096 default for firmware code generation"
  - "System prompt includes BARR-C requirements + top 3 similar patterns + board-specific APIs"
  - "Code truncation: 50 lines max per pattern example to avoid over-long prompts"
  - "Graceful degradation: Generate without semantic search if OPENAI_API_KEY not set"
  - "Board-specific API guidance: ESP32 (driver/gpio.h), STM32 (HAL), Arduino (digitalWrite)"
  - "Extract code from markdown blocks (```c ... ```)"
  - "CLI output: Save to file or print to stdout"

patterns-established:
  - "ClaudeClient wrapper pattern for API abstraction"
  - "PromptBuilder queries knowledge graph for context injection"
  - "Graceful degradation without vector store (BARR-C only)"
  - "Board-specific API guidance in system prompt"
  - "CLI command structure: <verb> <spec> --board <type> [--output <file>]"
  - "Progress indicators and next steps guidance"

issues-created: []

# Metrics
duration: 45min
completed: 2026-02-13
---

# Phase 07-01: Code Generation Engine Summary

**Claude API code generator with pattern-based prompts, BARR-C requirements, and board-specific API guidance via percepta generate CLI**

## Performance

- **Duration:** 45 min
- **Started:** 2026-02-13T01:30:00Z
- **Completed:** 2026-02-13T02:15:00Z
- **Tasks:** 3
- **Files created:** 6
- **Tests:** 10 (all passing)

## Accomplishments
- Claude API client for firmware code generation (ClaudeClient)
- Pattern-based prompt builder queries knowledge graph for similar validated patterns
- System prompts include BARR-C requirements, pattern examples, and board-specific APIs
- CLI command: percepta generate <spec> --board <type> [--output <file>]
- Graceful degradation without semantic search (BARR-C requirements still work)
- All tests passing (10 passed)

## Task Commits

Each task was committed atomically:

1. **Task 1: Create Claude API client for code generation** - `330a666` (feat)
2. **Task 2: Add pattern-based prompt engineering** - `e0b5812` (feat)
3. **Task 3: Add percepta generate CLI command** - `ef4eabc` (feat)

## Files Created/Modified

### Created Files
- `internal/codegen/claude_client.go` - Anthropic API wrapper for code generation
- `internal/codegen/claude_client_test.go` - Tests for ClaudeClient (5 tests)
- `internal/codegen/prompt_builder.go` - Pattern-based prompt engineering
- `internal/codegen/prompt_builder_test.go` - Tests for PromptBuilder (5 tests)
- `cmd/percepta/generate.go` - CLI command for firmware generation

### Modified Files
- `cmd/percepta/main.go` - Registered generateCmd

## Architecture

### Claude API Client

**ClaudeClient design:**
```go
type ClaudeClient struct {
    apiKey string
    model  string  // claude-sonnet-4-5-20250929
    client anthropic.Client
}

func (c *ClaudeClient) GenerateCode(
    spec string,           // "Blink LED at 1Hz"
    boardType string,      // "esp32"
    systemPrompt string,   // BARR-C + patterns
    maxTokens int,         // 4096
) (string, error)
```

**Key features:**
- API key from ANTHROPIC_API_KEY environment variable
- Model: claude-sonnet-4-5-20250929 (latest Sonnet 4.5)
- Temperature: 0.3 (lower for deterministic code)
- Max tokens: 4096 default (suitable for firmware code)
- Code extraction from markdown blocks (```c ... ```)
- Error handling for missing API key, empty responses

### Pattern-Based Prompt Engineering

**PromptBuilder design:**
```go
type PromptBuilder struct {
    patternStore *knowledge.PatternStore
}

func (p *PromptBuilder) BuildSystemPrompt(
    spec string,
    boardType string,
) (string, error)
```

**System prompt structure:**
1. **BARR-C requirements:** Naming conventions, types, structure
2. **Validated patterns:** Top 3 similar patterns from knowledge graph
3. **Board-specific APIs:** GPIO, timers, includes for target board

**Pattern context example:**
```
Validated Patterns (these work on real hardware):

Example 1 (95% similar, 97% confidence):
Spec: Blink LED at 1Hz
Board: esp32
Code:
```c
#include <stdint.h>
#include "driver/gpio.h"

void LED_Init(void) { ... }
```
Hardware validation: 2 signals observed
```

**Graceful degradation:**
- If vector store not initialized: Skip pattern search, continue with BARR-C only
- If semantic search fails: Log warning, continue with basic requirements
- User can generate code without OPENAI_API_KEY (degraded quality)

### Board-Specific API Guidance

**ESP32:**
- GPIO: `gpio_set_level()` from `driver/gpio.h`
- Timers: `esp_timer_create()` from `esp_timer.h`
- No Arduino-style `digitalWrite()`

**STM32:**
- GPIO: `HAL_GPIO_WritePin()`, `HAL_GPIO_ReadPin()`
- Timers: `HAL_TIM_Base_Start_IT()`
- Include: `stm32f4xx_hal.h`

**Arduino:**
- GPIO: `digitalWrite()`, `digitalRead()`
- Timers: `millis()` for non-blocking
- Include: `Arduino.h`

**Generic:**
- Standard C library
- `stdint.h` for fixed-width types
- Software timers

## CLI Usage

### Basic usage
```bash
percepta generate "Blink LED at 1Hz" --board esp32 --output led_blink.c
```

Output:
```
Generating firmware...
Spec: Blink LED at 1Hz
Board: esp32

✓ Semantic search enabled
Querying Claude API...
✓ Code generated (42 lines)

Saved to: led_blink.c

--- Next Steps ---
1. Validate style: percepta style-check led_blink.c
2. Review code for correctness
3. Flash to hardware and test
4. Observe behavior: percepta observe <device>
5. Store validated pattern: percepta knowledge store ...
```

### Without semantic search
```bash
# OPENAI_API_KEY not set
percepta generate "Toggle LED" --board stm32
```

Output:
```
Note: OPENAI_API_KEY not set, semantic search disabled
Continuing with basic BARR-C requirements...

Querying Claude API...
✓ Code generated (38 lines)
```

### Print to stdout
```bash
percepta generate "Button interrupt handler" --board arduino
```

Prints generated code to terminal.

### Help text
```bash
percepta generate --help
```

Shows usage, examples, requirements (ANTHROPIC_API_KEY).

## Test Coverage

**10 tests total, all passing:**

### ClaudeClient Tests (5 tests)
- `TestNewClaudeClient` - API key handling (explicit, environment, override)
- `TestExtractCode` - Code extraction from markdown blocks
- `TestClaudeClient_GenerateCode_NoAPIKey` - Error without API key
- `TestClaudeClient_GenerateCode_DefaultMaxTokens` - Default max tokens
- `TestClaudeClient_GenerateCode_Integration` - Integration test (skipped without API key)

### PromptBuilder Tests (5 tests)
- `TestPromptBuilder_BuildSystemPrompt_WithoutPatterns` - BARR-C requirements only
- `TestPromptBuilder_BuildSystemPrompt_WithPatterns` - With pattern context
- `TestPromptBuilder_BoardSpecificNotes` - Board-specific API guidance
- `TestPromptBuilder_CodeTruncation` - Long code truncation
- `TestPromptBuilder_GracefulDegradation` - Works without vector store

## Decisions Made

1. **Use Anthropic SDK directly:** Already in dependencies (anthropic-sdk-go v1.22.1), no need for custom HTTP client. Simplifies implementation and maintenance.

2. **Model selection:** claude-sonnet-4-5-20250929 (latest Claude Sonnet 4.5). Best balance of performance and quality for code generation.

3. **Temperature 0.3:** Lower than default (1.0) for more deterministic, consistent code generation. Still allows some creativity for different approaches.

4. **Max tokens 4096:** Suitable for firmware code (typically 50-200 lines). Can be adjusted by caller if needed.

5. **Top 3 patterns:** Balance between context richness and prompt length. More patterns = more tokens = higher cost and latency.

6. **Code truncation at 50 lines:** Prevents over-long prompts while still providing useful examples. Full pattern available in knowledge graph if needed.

7. **Graceful degradation:** Generate code even without semantic search (OPENAI_API_KEY). BARR-C requirements alone still produce valid code, just without pattern context.

8. **Board-specific API guidance:** Hardcoded for common boards (ESP32, STM32, Arduino). Extensible for future boards. Prevents common mistakes (e.g., using digitalWrite() on ESP32).

9. **CLI output options:** Save to file (--output) or print to stdout. File output is more common workflow, but stdout useful for quick inspection.

10. **Next steps guidance:** Show user validation workflow (style-check, test, observe, store). Helps close the loop back to knowledge graph.

## Deviations from Plan

**None.** Plan executed exactly as written.

All tasks completed successfully:
- ✅ Task 1: Claude API client (5 tests passing)
- ✅ Task 2: Pattern-based prompt engineering (5 tests passing)
- ✅ Task 3: CLI command (manual testing successful)

No bugs encountered, no scope creep, no architectural changes needed.

## Issues Encountered

**1. Anthropic SDK API usage**
- **Problem:** Initial implementation used incorrect API patterns (`anthropic.F()`, `AsUnion()`)
- **Resolution:** Checked SDK documentation via `go doc`, corrected to proper API usage
- **Verification:** All tests passing, build successful

**Smooth integration points:**
- PatternStore API already clean (from Phase 06-02)
- SearchSimilarPatterns works as documented
- Cobra CLI patterns established (from Phase 04-01)
- All dependencies already installed (anthropic-sdk-go)

## Integration Notes

**Ready for Phase 07-02 (Validation Pipeline):**
- Code generation working end-to-end
- Pattern-based prompts functional
- CLI command integrated
- Next phase will add:
  - Style validation with auto-fix
  - Hardware validation loop
  - Pattern storage after validation

**How Phase 07-02 will extend this:**
```go
// Complete generation + validation workflow
func generateAndValidate(spec string, board string) error {
    // 1. Generate code (this phase)
    code := generateCmd.Run(spec, board)

    // 2. Validate style (Phase 07-02)
    violations := styleCheck.Check(code)
    if len(violations) > 0 {
        code = styleFix.Apply(code, violations)
    }

    // 3. Flash to hardware (Phase 07-02)
    flash(code, device)

    // 4. Observe behavior (Phase 07-02)
    obs := observe(device)

    // 5. Store validated pattern (Phase 07-02)
    store.StoreValidatedPattern(spec, code, device, firmware)
}
```

## Example Prompt

**For spec: "Blink LED at 1Hz", board: "esp32"**

System prompt includes:
```
You are an expert embedded firmware engineer writing BARR-C compliant code.

BARR-C Style Requirements:
- Function names: Module_Function() format (e.g., LED_Init, Timer_Start)
- Variables: snake_case (e.g., led_state, timer_count)
- Constants: UPPER_SNAKE with meaningful names (e.g., LED_PIN, BLINK_PERIOD_MS)
- Types: Use stdint.h (uint8_t, uint16_t, uint32_t, not int/char)
- No magic numbers: Define all constants with #define
...

Validated Patterns (these work on real hardware):

Example 1 (95% similar, 97% confidence):
Spec: Blink LED at 1Hz
Board: esp32
Code:
```c
#include <stdint.h>
#include "driver/gpio.h"
#include "esp_timer.h"

#define LED_PIN 2
#define BLINK_PERIOD_MS 1000

static uint8_t led_state = 0;

void LED_Toggle(void) {
    led_state = !led_state;
    gpio_set_level(LED_PIN, led_state);
}
```
Hardware validation: 2 signals observed

ESP32-specific:
- GPIO pins: Use gpio_set_level() from driver/gpio.h
- Timers: Use esp_timer_create() for non-blocking timing
- No Arduino-style digitalWrite()
...
```

User message:
```
Generate firmware for esp32 board:

Specification: Blink LED at 1Hz

Requirements:
- BARR-C compliant
- Use validated patterns provided in system prompt
- Non-blocking architecture (timers, not delays)
- Proper error handling
- Static allocation only

Output only the C source code, no explanations.
```

## Next Phase Readiness

**Phase 07-02 (Validation Pipeline) is ready to start:**
- ✅ Code generation functional
- ✅ Pattern-based prompts working
- ✅ CLI command integrated
- ✅ Graceful degradation tested
- ✅ Board-specific guidance available

**No blockers for Phase 07-02.**

**Plan 07-01 complete. Ready for validation pipeline integration.**

---
*Phase: 07-code-generation-engine*
*Completed: 2026-02-13*
