# Command Reference

Complete reference for all Percepta CLI commands.

## percepta observe

Observe hardware state via computer vision.

**Usage:**
```bash
percepta observe <device> [flags]
```

**Description:**

Captures frames from the device's camera, analyzes them with Claude Vision API, and stores the observation in SQLite. Detects LED states, display content, and boot timing.

**Examples:**
```bash
# Basic observation
percepta observe my-board

# Extended observation (longer duration)
percepta observe my-board --duration 10s

# Use specific camera
percepta observe my-board --camera /dev/video1
```

**Output:**

Shows detected signals with confidence scores:
- LED states (ON/OFF, blinking, color, frequency)
- Display content (LCD text via OCR)
- Boot timing (milliseconds)

**Exit codes:**
- `0` - Observation successful
- `1` - Error (device not found, camera error, API error)

---

## percepta assert

Validate hardware state against expected behavior.

**Usage:**
```bash
percepta assert <device> <assertion> [flags]
```

**Description:**

Captures an observation and evaluates the assertion DSL expression. Returns exit code 0 if passed, 1 if failed. Useful for CI/CD pipelines and automated testing.

**Examples:**
```bash
# LED is ON
percepta assert my-board "led power is ON"

# LED blinks at specific rate (±10% tolerance)
percepta assert my-board "led status blinks at 2Hz"

# LED color check (RGB ±5 tolerance)
percepta assert my-board "led status is blue"

# Display contains text
percepta assert my-board "display LCD shows 'Ready'"

# Multiple conditions (space-separated)
percepta assert my-board "led power is ON" "led error is OFF"
```

**Assertion DSL Syntax:**

**LED assertions:**
- `led <name> is ON` - LED is illuminated
- `led <name> is OFF` - LED is not illuminated
- `led <name> blinks at <hz>Hz` - LED blinks at frequency (±10% tolerance)
- `led <name> is <color>` - LED color matches (red, green, blue, etc.)

**Display assertions:**
- `display <name> shows '<text>'` - Display contains exact text
- `display <name> contains '<substring>'` - Display contains substring

**Timing assertions:**
- `boot time < <ms>ms` - Boot completes within time

**Exit codes:**
- `0` - Assertion passed
- `1` - Assertion failed
- `2` - Error (device not found, invalid syntax)

---

## percepta diff

Compare hardware behavior across firmware versions.

**Usage:**
```bash
percepta diff <device> --from <firmware1> --to <firmware2>
```

**Description:**

Retrieves observations for two firmware versions and compares them signal-by-signal. Shows added, removed, and modified signals. Useful for regression detection and feature validation.

**Examples:**
```bash
# Compare two firmware versions
percepta diff my-esp32 --from v1.0 --to v1.1

# Check for regressions from baseline
percepta diff my-board --from baseline --to feature-branch

# Validate behavior change
percepta diff test-device --from before --to after
```

**Output:**

Shows changes categorized as:
- `+` Added signals (new LEDs, displays)
- `-` Removed signals (LEDs turned off, displays cleared)
- `~` Modified signals (blink rate change, color change, text change)

**Example output:**
```
Comparing firmware versions:
FROM: v1.0 (2026-02-11 10:30:00)
TO:   v1.1 (2026-02-11 11:45:00)

Device: my-esp32

Changes detected:

+ LED2: purple blinking ~0.8Hz (ADDED)
- LED3: red solid (REMOVED)
~ LED1: blue blinking 2.0Hz → 2.5Hz (MODIFIED)

Summary: 1 added, 1 removed, 1 modified
```

**Exit codes:**
- `0` - No changes detected (identical behavior)
- `1` - Changes detected
- `2` - Error (device not found, firmware tag missing)

---

## percepta device

Manage device configurations.

**Usage:**
```bash
percepta device <subcommand>
```

**Subcommands:**
- `list` - List all configured devices
- `add <name>` - Add a new device
- `set-firmware <device> <version>` - Update firmware tag

### percepta device list

List all configured devices.

**Usage:**
```bash
percepta device list
```

**Example output:**
```
Configured devices:

my-esp32
  Type: esp32
  Camera: /dev/video0
  Firmware: v1.0

lab-stm32
  Camera: /dev/video2
  Firmware: v2.1.0
```

### percepta device add

Add a new device configuration.

**Usage:**
```bash
percepta device add <name>
```

**Interactive prompts:**
1. Device type (e.g., esp32, stm32, arduino)
2. Camera device path (e.g., /dev/video0)
3. Firmware version (optional)

**Example:**
```bash
$ percepta device add my-board
Device type (e.g., fpga, esp32, stm32): esp32
Camera device path (default: /dev/video0): /dev/video0
Firmware version (optional, press Enter to skip): v1.0

✓ Device 'my-board' added successfully
```

### percepta device set-firmware

Update firmware tag for a device.

**Usage:**
```bash
percepta device set-firmware <device> <firmware>
```

**Examples:**
```bash
percepta device set-firmware my-esp32 v1.0
percepta device set-firmware my-board baseline
percepta device set-firmware test-board abc123
```

**Use case:**

Run this before observations to associate them with a specific firmware version. Enables firmware diffing with `percepta diff`.

---

## percepta generate

Generate BARR-C compliant firmware from specification.

**Usage:**
```bash
percepta generate <spec> --board <type> [--output <file>]
```

**Description:**

Uses Claude AI with knowledge from validated patterns to create professional, BARR-C compliant embedded C code. Generated code follows established working patterns and includes proper error handling, non-blocking architecture, and static allocation.

**Examples:**
```bash
# Generate LED blink code
percepta generate "Blink LED at 1Hz" --board esp32 --output led_blink.c

# Generate sensor reading code
percepta generate "Read temperature sensor every 2 seconds" --board stm32

# Generate button handler
percepta generate "Toggle LED on button press" --board arduino --output button_led.c
```

**Flags:**
- `--board` (required) - Target board type (esp32, stm32, arduino, atmega, generic)
- `--output` (optional) - Output file path (prints to stdout if not set)

**Requirements:**
- `ANTHROPIC_API_KEY` environment variable (Claude API)

**Optional:**
- `OPENAI_API_KEY` for semantic pattern search (graceful degradation without it)

**Workflow:**

1. Searches knowledge graph for similar validated patterns
2. Generates code using Claude with BARR-C requirements
3. Validates against style checker
4. Auto-fixes deterministic violations
5. Stores compliant patterns back to knowledge graph

**Exit codes:**
- `0` - Generation successful
- `1` - Generation failed (API error, invalid spec)

---

## percepta style-check

Check C code for BARR-C compliance.

**Usage:**
```bash
percepta style-check <file-or-directory> [--fix]
```

**Description:**

Validates code against BARR-C Embedded C Coding Standard. Checks naming conventions, type usage, const correctness, and magic numbers.

**Examples:**
```bash
# Check single file
percepta style-check led_blink.c

# Check entire directory
percepta style-check ./src

# Auto-fix violations
percepta style-check led_blink.c --fix
```

**Flags:**
- `--fix` - Auto-fix deterministic violations (naming, types)

**BARR-C Rules:**

**Naming:**
- Functions: `Module_Function` (e.g., `LED_Init`, `UART_Write`)
- Variables: `snake_case` (e.g., `button_state`, `timer_count`)
- Global constants: `UPPER_SNAKE` (e.g., `MAX_RETRIES`, `BUFFER_SIZE`)
- Local constants: `snake_case` (e.g., `timeout_ms`, `baud_rate`)

**Types:**
- Use `stdint.h` types: `uint8_t`, `int16_t`, `uint32_t`
- Avoid primitives: `unsigned char`, `short`, `long`

**Const correctness:**
- Read-only parameters: `const uint8_t* data`
- Read-only pointers: `uint8_t* const ptr`

**Magic numbers:**
- No bare numeric literals (except 0, 1, -1)
- Use named constants or enums

**Output format:**

Standard linter format for CI integration:
```
led_blink.c:10:5: error [naming] Function 'blinkLED' should be 'LED_Blink' (Module_Function)
led_blink.c:15:9: error [types] Use 'uint8_t' instead of 'unsigned char'
led_blink.c:20:12: warning [magic-number] Magic number '500' should be named constant
```

**Exit codes:**
- `0` - No violations (BARR-C compliant)
- `1` - Violations found

---

## percepta knowledge

Manage validated pattern knowledge graph.

**Usage:**
```bash
percepta knowledge <subcommand>
```

**Subcommands:**
- `store` - Store a validated pattern
- `search` - Search for similar patterns
- `list` - List all patterns

### percepta knowledge store

Store validated pattern in knowledge graph.

**Usage:**
```bash
percepta knowledge store <spec> <file.c> --device <device-id> --firmware <tag>
```

**Description:**

Validates code against BARR-C, links to device observation, and stores in graph. Only stores patterns that are:
1. BARR-C compliant (no style violations)
2. Hardware-validated (observation exists for device+firmware)
3. Successfully linked in graph

**Examples:**
```bash
percepta knowledge store "Blink LED at 1Hz" led.c --device esp32-dev --firmware v1.0.0
percepta knowledge store "Button debounce" button.c --device my-board --firmware baseline
```

**Flags:**
- `--device` (required) - Device ID from observations
- `--firmware` (required) - Firmware tag from device config

### percepta knowledge search

Search for similar validated patterns.

**Usage:**
```bash
percepta knowledge search <query> [--board <type>] [--limit <n>]
```

**Description:**

Semantic search using vector embeddings. Finds patterns by code similarity, not exact text matches.

**Examples:**
```bash
# Search for debounce patterns
percepta knowledge search "button debounce" --board esp32

# Find LED patterns
percepta knowledge search "LED blink" --limit 5

# General search
percepta knowledge search "sensor reading"
```

**Flags:**
- `--board` (optional) - Filter by board type
- `--limit` (optional) - Max results (default: 10)

**Requirements:**
- `OPENAI_API_KEY` for vector embeddings

### percepta knowledge list

List all validated patterns.

**Usage:**
```bash
percepta knowledge list
```

**Example output:**
```
Validated patterns:

1. Blink LED at 1Hz
   Board: esp32
   Device: esp32-dev
   Firmware: v1.0.0
   Style: Compliant

2. Button debounce 50ms
   Board: stm32
   Device: test-board
   Firmware: v2.0
   Style: Compliant
```

---

## Supported Boards

Percepta supports code generation for these board types:

**ESP32:**
- ESP32-WROOM
- ESP32-DevKit
- ESP32-S2/S3/C3

**STM32:**
- STM32F4 series
- STM32F1 series
- STM32L4 series

**Arduino:**
- Arduino Uno (ATmega328P)
- Arduino Mega (ATmega2560)
- Arduino Nano

**Generic:**
- Custom boards
- FPGA soft cores
- RISC-V boards

**Board-specific notes:**

Each board type gets tailored API guidance:
- ESP32: FreeRTOS, ESP-IDF APIs
- STM32: HAL, CMSIS
- Arduino: Arduino libraries
- Generic: ANSI C only

---

## Environment Variables

**Required for observation:**
- `ANTHROPIC_API_KEY` - Claude Vision API key from [console.anthropic.com](https://console.anthropic.com)

**Required for code generation:**
- `ANTHROPIC_API_KEY` - Claude API key

**Optional for semantic search:**
- `OPENAI_API_KEY` - OpenAI API key for embeddings

**Configuration:**

```bash
# Linux/macOS (add to ~/.bashrc or ~/.zshrc)
export ANTHROPIC_API_KEY="sk-ant-..."
export OPENAI_API_KEY="sk-..."

# Windows (PowerShell)
$env:ANTHROPIC_API_KEY="sk-ant-..."
$env:OPENAI_API_KEY="sk-..."
```

---

## Configuration File

Config file: `~/.config/percepta/config.yaml`

**Example:**

```yaml
devices:
  my-esp32:
    type: esp32
    camera: /dev/video0
    firmware: v1.2.0

  lab-stm32:
    type: stm32
    camera: /dev/video2
    firmware: v2.0.1

vision:
  frames: 5
  interval: 200ms
  confidence_threshold: 0.7

storage:
  path: ~/.local/share/percepta/percepta.db

knowledge:
  path: ~/.local/share/percepta/knowledge.db
```

**Fields:**

**devices:**
- `type` - Board type (optional, human-readable)
- `camera` - Camera device path
- `firmware` - Current firmware version tag

**vision:**
- `frames` - Number of frames to capture (default: 5)
- `interval` - Time between frames (default: 200ms)
- `confidence_threshold` - Minimum confidence for signals (default: 0.7)

**storage:**
- `path` - SQLite database location (default: ~/.local/share/percepta/percepta.db)

**knowledge:**
- `path` - Knowledge graph database (default: ~/.local/share/percepta/knowledge.db)

---

## Exit Codes

Percepta commands use standard exit codes:

- `0` - Success
- `1` - Failure (assertion failed, violations found, changes detected)
- `2` - Error (device not found, API error, invalid syntax)

**CI/CD integration:**

```bash
# Fail build on assertion failure
percepta assert my-board "led status is ON" || exit 1

# Fail build on behavioral changes
percepta diff my-board --from baseline --to ${CI_COMMIT_SHA}
if [ $? -eq 1 ]; then
  echo "Behavioral regression detected"
  exit 1
fi
```

---

## Common Workflows

See [examples.md](examples.md) for 20+ detailed workflow examples covering:
- LED validation
- Boot time monitoring
- Regression detection
- Code generation pipelines
- CI/CD integration
