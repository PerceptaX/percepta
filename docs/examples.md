# Example Workflows

Comprehensive collection of real-world Percepta workflows covering common embedded development scenarios.

---

## Table of Contents

**Basic Workflows:**
1. [Validate LED Blink Rate](#1-validate-led-blink-rate)
2. [Check Boot Time](#2-check-boot-time)
3. [Verify Display Content](#3-verify-display-content)
4. [Detect LED Color Change](#4-detect-led-color-change)
5. [Multiple LED Validation](#5-multiple-led-validation)

**Firmware Tracking:**
6. [Track Firmware Versions](#6-track-firmware-versions)
7. [Detect Boot Regression](#7-detect-boot-regression)
8. [Compare Feature Branches](#8-compare-feature-branches)
9. [Validate Bug Fix](#9-validate-bug-fix)
10. [Monitor Incremental Changes](#10-monitor-incremental-changes)

**Code Generation:**
11. [Generate Validated Firmware](#11-generate-validated-firmware)
12. [Generate with Pattern Search](#12-generate-with-pattern-search)
13. [Iterate on Generated Code](#13-iterate-on-generated-code)
14. [Fix Style Violations](#14-fix-style-violations)
15. [Store Validated Pattern](#15-store-validated-pattern)

**CI/CD Integration:**
16. [Automated Hardware Testing](#16-automated-hardware-testing)
17. [PR Validation Pipeline](#17-pr-validation-pipeline)
18. [Nightly Regression Suite](#18-nightly-regression-suite)
19. [Release Validation](#19-release-validation)
20. [Multi-Board Testing](#20-multi-board-testing)

**Advanced:**
21. [Debug LED Pattern Issues](#21-debug-led-pattern-issues)
22. [Validate Power State Transitions](#22-validate-power-state-transitions)
23. [Test Error Indicator Behavior](#23-test-error-indicator-behavior)
24. [Compare Production vs Debug Builds](#24-compare-production-vs-debug-builds)
25. [Hardware-in-Loop Testing](#25-hardware-in-loop-testing)

---

## Basic Workflows

### 1. Validate LED Blink Rate

**Scenario:** Verify LED blinks at exactly 2Hz after firmware update.

**Steps:**

```bash
# 1. Observe current behavior
percepta observe my-esp32

# 2. Assert expected blink rate (±10% tolerance)
percepta assert my-esp32 "led status blinks at 2Hz"

# Output: ✓ Assertion passed
```

**Expected output:**
```
✓ Assertion passed

Expected: LED 'status' blinks at 2.0 Hz
Actual:   LED 'status' blinks at 1.98 Hz
Confidence: 0.92
```

**Troubleshooting:**

If assertion fails:
```bash
# Check actual blink rate
percepta observe my-esp32

# Adjust tolerance if hardware variance expected
percepta assert my-esp32 "led status blinks at 2Hz" --tolerance 0.2  # ±20%
```

---

### 2. Check Boot Time

**Scenario:** Ensure device boots within 3 seconds.

**Steps:**

```bash
# 1. Power cycle device, then observe
percepta observe my-board

# 2. Assert boot time constraint
percepta assert my-board "boot time < 3000ms"

# Output: ✓ Assertion passed
```

**CI/CD usage:**
```bash
#!/bin/bash
# Reset device
./scripts/reset_device.sh

# Wait for boot
sleep 1

# Validate boot time
percepta observe my-board
percepta assert my-board "boot time < 3000ms" || exit 1
```

---

### 3. Verify Display Content

**Scenario:** Check LCD shows "Ready" after initialization.

**Steps:**

```bash
# 1. Observe display
percepta observe my-stm32

# 2. Assert display content
percepta assert my-stm32 "display LCD shows 'Ready'"

# Output: ✓ Assertion passed
```

**Partial match:**
```bash
# Use 'contains' for substring matching
percepta assert my-stm32 "display LCD contains 'Rdy'"
```

**Note:** OCR can be noisy. Use `contains` for partial matches rather than exact text.

---

### 4. Detect LED Color Change

**Scenario:** Verify status LED changes from red to green after connection.

**Before connection:**
```bash
percepta observe my-esp32
# Output: LED 'status': red solid

percepta assert my-esp32 "led status is red"
# Output: ✓ Assertion passed
```

**After connection:**
```bash
percepta observe my-esp32
# Output: LED 'status': green solid

percepta assert my-esp32 "led status is green"
# Output: ✓ Assertion passed
```

**RGB matching:**
```bash
# Exact RGB values (±5 tolerance)
percepta assert my-esp32 "led status color RGB(0,255,0)"
```

---

### 5. Multiple LED Validation

**Scenario:** Validate multi-LED state machine.

**Expected behavior:**
- Power LED: solid green
- Status LED: blinking blue at 1Hz
- Error LED: OFF

**Steps:**

```bash
# Observe all LEDs
percepta observe my-board

# Assert all conditions (space-separated)
percepta assert my-board \
  "led power is ON" \
  "led power is green" \
  "led status blinks at 1Hz" \
  "led status is blue" \
  "led error is OFF"

# Output: ✓ All assertions passed
```

---

## Firmware Tracking

### 6. Track Firmware Versions

**Scenario:** Track behavior across development cycle.

**Workflow:**

```bash
# 1. Tag initial version
percepta device set-firmware my-esp32 v1.0-baseline

# 2. Observe baseline behavior
percepta observe my-esp32

# 3. Develop and flash new firmware
# ... (your build and flash process) ...

# 4. Tag new version
percepta device set-firmware my-esp32 v1.1-feature-x

# 5. Observe new behavior
percepta observe my-esp32

# 6. Compare versions
percepta diff my-esp32 --from v1.0-baseline --to v1.1-feature-x
```

**Example diff output:**
```
Changes detected:

+ LED2: green blinking ~1Hz (ADDED)
~ LED1: blue solid → blue blinking 2Hz (MODIFIED)

Summary: 1 added, 0 removed, 1 modified
```

---

### 7. Detect Boot Regression

**Scenario:** Catch boot time regression introduced in new firmware.

**Steps:**

```bash
# Before optimization
percepta device set-firmware my-board v1.0
percepta observe my-board
# Boot time: 1200ms

# After optimization attempt
percepta device set-firmware my-board v1.1
percepta observe my-board
# Boot time: 2800ms

# Compare
percepta diff my-board --from v1.0 --to v1.1

# Output shows boot time change:
# ~ Boot timing: 1200ms → 2800ms (MODIFIED)
```

**Exit code check:**
```bash
percepta diff my-board --from v1.0 --to v1.1
if [ $? -eq 1 ]; then
  echo "⚠️  Behavioral change detected"
  exit 1
fi
```

---

### 8. Compare Feature Branches

**Scenario:** Validate feature branch doesn't break existing behavior.

**Git workflow integration:**

```bash
#!/bin/bash
# Compare feature branch against main

# 1. Flash main branch firmware
git checkout main
make flash
percepta device set-firmware my-board main
percepta observe my-board

# 2. Flash feature branch firmware
git checkout feature/new-sensor
make flash
percepta device set-firmware my-board feature/new-sensor
percepta observe my-board

# 3. Compare
percepta diff my-board --from main --to feature/new-sensor

# 4. Check for regressions
if [ $? -eq 1 ]; then
  echo "Behavioral changes detected. Review diff above."
  exit 1
fi
```

---

### 9. Validate Bug Fix

**Scenario:** Confirm bug fix works without side effects.

**Bug:** LED should blink at 1Hz but blinks at 0.5Hz

**Before fix:**
```bash
percepta device set-firmware my-esp32 v1.0-buggy
percepta observe my-esp32
# Output: LED 'status' blinks at 0.48Hz

percepta assert my-esp32 "led status blinks at 1Hz"
# Output: ✗ Assertion failed (expected 1.0Hz, got 0.48Hz)
```

**After fix:**
```bash
percepta device set-firmware my-esp32 v1.1-fixed
percepta observe my-esp32
# Output: LED 'status' blinks at 0.99Hz

percepta assert my-esp32 "led status blinks at 1Hz"
# Output: ✓ Assertion passed

# Verify no side effects
percepta diff my-esp32 --from v1.0-buggy --to v1.1-fixed
# Output should show ONLY the blink rate fix
```

---

### 10. Monitor Incremental Changes

**Scenario:** Track behavior changes across daily builds.

**Daily build script:**

```bash
#!/bin/bash
# daily-build.sh

DATE=$(date +%Y-%m-%d)
VERSION="daily-${DATE}"

# Build and flash
make clean && make
make flash

# Tag and observe
percepta device set-firmware my-board $VERSION
percepta observe my-board

# Compare to yesterday
YESTERDAY=$(date -d "yesterday" +%Y-%m-%d)
YESTERDAY_VERSION="daily-${YESTERDAY}"

percepta diff my-board --from $YESTERDAY_VERSION --to $VERSION

# Store diff results
percepta diff my-board --from $YESTERDAY_VERSION --to $VERSION > daily-diffs/${DATE}.txt
```

---

## Code Generation

### 11. Generate Validated Firmware

**Scenario:** Generate LED blink code and validate on hardware.

**Steps:**

```bash
# 1. Generate code
percepta generate "Blink LED at 1Hz" --board esp32 --output led_blink.c

# Output:
# ✓ Semantic search enabled
# Generating and validating code...
# ✓ Generated code
#
# Style check: PASSED (0 violations)
# Auto-fixes applied: 2 (naming corrections)
# Pattern stored: YES
#
# Saved to: led_blink.c

# 2. Flash to hardware (your build process)
# ... compile and flash led_blink.c ...

# 3. Validate behavior
percepta observe my-esp32

# 4. Assert expected behavior
percepta assert my-esp32 "led status blinks at 1Hz"

# Output: ✓ Assertion passed
```

**Full validation workflow:**

```bash
# Store validated pattern
percepta knowledge store "Blink LED at 1Hz" led_blink.c \
  --device my-esp32 \
  --firmware v1.0
```

---

### 12. Generate with Pattern Search

**Scenario:** Generate code similar to existing validated patterns.

**Steps:**

```bash
# 1. Search for similar patterns
percepta knowledge search "button debounce" --board esp32 --limit 3

# Output:
# 1. Button debounce 50ms (similarity: 0.89)
# 2. Button with LED toggle (similarity: 0.76)
# 3. Interrupt-based button handler (similarity: 0.68)

# 2. Generate with context from similar patterns
percepta generate "Debounce button with 50ms delay" --board esp32 --output button.c

# Output shows patterns used:
# ✓ Using 3 similar patterns:
#   - Button debounce 50ms (0.89 similarity)
#   - Button with LED toggle (0.76 similarity)
#   - Interrupt-based button handler (0.68 similarity)
#
# ✓ Generated code (style compliant)
```

**Without OPENAI_API_KEY:**

```bash
# Graceful degradation - still generates with BARR-C requirements
unset OPENAI_API_KEY

percepta generate "Blink LED at 1Hz" --board esp32 --output led.c

# Output:
# Note: OPENAI_API_KEY not set, semantic search disabled
# Continuing with basic BARR-C requirements...
#
# ✓ Generated code (style compliant)
```

---

### 13. Iterate on Generated Code

**Scenario:** Refine specification until code meets requirements.

**Iteration 1 - Too simple:**

```bash
percepta generate "LED control" --board esp32 --output led.c

# Review code - too generic, add details
```

**Iteration 2 - More specific:**

```bash
percepta generate "Toggle LED on GPIO 2 every second" --board esp32 --output led.c

# Review code - good, but want button trigger
```

**Iteration 3 - Final spec:**

```bash
percepta generate "Toggle LED on GPIO 2 when button on GPIO 0 pressed" --board esp32 --output led.c

# Flash and validate
# ... compile and flash ...
percepta observe my-esp32
percepta assert my-esp32 "led toggles on button press"
```

**Tip:** Be specific in specifications. Include:
- GPIO pins
- Timing requirements
- Trigger conditions
- Expected behavior

---

### 14. Fix Style Violations

**Scenario:** Generated code has minor violations, fix automatically.

**Steps:**

```bash
# 1. Generate code
percepta generate "Read sensor value" --board stm32 --output sensor.c

# 2. Check style (if violations remain)
percepta style-check sensor.c

# Output:
# sensor.c:15:9: error [types] Use 'uint8_t' instead of 'unsigned char'
# sensor.c:20:5: error [naming] Function 'readSensor' should be 'Sensor_Read'

# 3. Auto-fix
percepta style-check sensor.c --fix

# Output:
# ✓ Fixed 2 violations:
#   - types: 1 fix
#   - naming: 1 fix
#
# sensor.c is now BARR-C compliant

# 4. Verify
percepta style-check sensor.c

# Output: ✓ No violations (BARR-C compliant)
```

---

### 15. Store Validated Pattern

**Scenario:** Store hardware-validated pattern for future reuse.

**Requirements:**
1. Code is BARR-C compliant
2. Hardware validation exists (observation)

**Steps:**

```bash
# 1. Flash firmware
# ... compile and flash ...

# 2. Tag firmware
percepta device set-firmware my-esp32 v1.0

# 3. Validate on hardware
percepta observe my-esp32
percepta assert my-esp32 "led status blinks at 2Hz"
# Output: ✓ Assertion passed

# 4. Store pattern
percepta knowledge store "Blink LED at 2Hz" led_blink.c \
  --device my-esp32 \
  --firmware v1.0

# Output:
# ✓ Style check: PASSED
# ✓ Observation found: my-esp32 @ v1.0
# ✓ Pattern stored in knowledge graph
```

**Future reuse:**

```bash
# Search for similar patterns in future projects
percepta knowledge search "LED blink" --board esp32

# Output includes your validated pattern:
# 1. Blink LED at 2Hz (similarity: 0.95)
#    Board: esp32
#    Device: my-esp32
#    Firmware: v1.0
#    Style: Compliant
#    Observation: Present
```

---

## CI/CD Integration

### 16. Automated Hardware Testing

**Scenario:** Run hardware validation in CI pipeline.

**GitHub Actions example:**

```yaml
name: Hardware Validation

on: [push, pull_request]

jobs:
  hardware-test:
    runs-on: self-hosted  # Runner with hardware access
    steps:
      - uses: actions/checkout@v3

      - name: Build firmware
        run: |
          cd firmware
          make clean
          make

      - name: Flash firmware
        run: |
          make flash

      - name: Hardware validation
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
        run: |
          percepta device set-firmware test-board ${{ github.sha }}
          percepta observe test-board
          percepta assert test-board "led power is ON"
          percepta assert test-board "boot time < 3000ms"

      - name: Compare to main
        if: github.ref != 'refs/heads/main'
        run: |
          percepta diff test-board --from main --to ${{ github.sha }}
          if [ $? -eq 1 ]; then
            echo "⚠️  Behavioral changes detected"
          fi
```

---

### 17. PR Validation Pipeline

**Scenario:** Validate PR doesn't introduce regressions.

**Script: `.github/workflows/pr-validation.sh`**

```bash
#!/bin/bash
set -e

PR_SHA=$1
BASE_SHA=$2

echo "Validating PR: $PR_SHA against base: $BASE_SHA"

# Build and flash base
git checkout $BASE_SHA
make flash
percepta device set-firmware test-board base-$BASE_SHA
percepta observe test-board

# Build and flash PR
git checkout $PR_SHA
make flash
percepta device set-firmware test-board pr-$PR_SHA
percepta observe test-board

# Compare
echo "Comparing behavior..."
percepta diff test-board --from base-$BASE_SHA --to pr-$PR_SHA

DIFF_EXIT=$?

if [ $DIFF_EXIT -eq 0 ]; then
  echo "✓ No behavioral changes"
  exit 0
elif [ $DIFF_EXIT -eq 1 ]; then
  echo "⚠️  Behavioral changes detected. Review required."
  # Don't fail - just warn
  exit 0
else
  echo "✗ Error running diff"
  exit 1
fi
```

---

### 18. Nightly Regression Suite

**Scenario:** Run comprehensive regression tests nightly.

**Script: `nightly-test.sh`**

```bash
#!/bin/bash
# Comprehensive nightly hardware regression suite

DATE=$(date +%Y-%m-%d)
VERSION="nightly-${DATE}"

# Build firmware
echo "Building firmware..."
make clean && make

# Flash firmware
echo "Flashing firmware..."
make flash

# Tag version
percepta device set-firmware test-board $VERSION

# Run test suite
echo "Running hardware validation suite..."

# Test 1: Boot time
percepta observe test-board
percepta assert test-board "boot time < 3000ms" || {
  echo "FAIL: Boot time regression"
  exit 1
}

# Test 2: LED patterns
percepta assert test-board "led power is ON" || {
  echo "FAIL: Power LED not on"
  exit 1
}

percepta assert test-board "led status blinks at 1Hz" || {
  echo "FAIL: Status LED incorrect"
  exit 1
}

# Test 3: Display
percepta assert test-board "display LCD contains 'Ready'" || {
  echo "FAIL: Display not showing Ready"
  exit 1
}

# Compare to yesterday
YESTERDAY=$(date -d "yesterday" +%Y-%m-%d)
YESTERDAY_VERSION="nightly-${YESTERDAY}"

echo "Comparing to yesterday..."
percepta diff test-board --from $YESTERDAY_VERSION --to $VERSION > nightly-diffs/${DATE}.txt

echo "✓ All tests passed"

# Email results
cat nightly-diffs/${DATE}.txt | mail -s "Nightly test results ${DATE}" team@company.com
```

---

### 19. Release Validation

**Scenario:** Final validation before release.

**Script: `release-validation.sh`**

```bash
#!/bin/bash
RELEASE_VERSION=$1

if [ -z "$RELEASE_VERSION" ]; then
  echo "Usage: $0 <version>"
  exit 1
fi

echo "Validating release: $RELEASE_VERSION"

# Build release firmware
make clean
make release VERSION=$RELEASE_VERSION

# Flash to multiple test boards
for BOARD in board1 board2 board3; do
  echo "Testing $BOARD..."

  # Flash firmware
  make flash DEVICE=$BOARD

  # Tag version
  percepta device set-firmware $BOARD $RELEASE_VERSION

  # Observe
  percepta observe $BOARD

  # Run assertions
  percepta assert $BOARD "led power is ON" || {
    echo "FAIL: $BOARD power LED"
    exit 1
  }

  percepta assert $BOARD "boot time < 3000ms" || {
    echo "FAIL: $BOARD boot time"
    exit 1
  }

  echo "✓ $BOARD passed"
done

echo "✓ Release $RELEASE_VERSION validated on all boards"
```

---

### 20. Multi-Board Testing

**Scenario:** Test firmware across different board types.

**Script: `multi-board-test.sh`**

```bash
#!/bin/bash
VERSION=$1

# Board configs
declare -A BOARDS
BOARDS[esp32-dev]="esp32"
BOARDS[stm32-nucleo]="stm32"
BOARDS[arduino-uno]="arduino"

for BOARD in "${!BOARDS[@]}"; do
  BOARD_TYPE=${BOARDS[$BOARD]}

  echo "Testing $BOARD ($BOARD_TYPE)..."

  # Build for board type
  make clean
  make BOARD=$BOARD_TYPE

  # Flash
  make flash DEVICE=$BOARD

  # Tag and observe
  percepta device set-firmware $BOARD $VERSION
  percepta observe $BOARD

  # Board-specific assertions
  case $BOARD_TYPE in
    esp32)
      percepta assert $BOARD "led status blinks at 1Hz"
      ;;
    stm32)
      percepta assert $BOARD "led power is ON"
      percepta assert $BOARD "display LCD contains 'Ready'"
      ;;
    arduino)
      percepta assert $BOARD "led status blinks at 0.5Hz"
      ;;
  esac

  echo "✓ $BOARD passed"
done

echo "✓ All boards validated"
```

---

## Advanced

### 21. Debug LED Pattern Issues

**Scenario:** LED not behaving as expected, debug systematically.

**Steps:**

```bash
# 1. Capture current behavior
percepta observe my-board

# Output shows unexpected pattern:
# LED 'status': blue blinking ~0.47Hz

# Expected: 1Hz, Actual: 0.47Hz

# 2. Check if it's hardware variance
# Observe multiple times
for i in {1..5}; do
  echo "Observation $i:"
  percepta observe my-board | grep "LED 'status'"
  sleep 2
done

# Output shows consistent 0.47Hz - not variance

# 3. Review firmware code
# Check timer configuration, delay values

# 4. Generate reference implementation
percepta generate "Blink LED at 1Hz" --board esp32 --output reference.c

# 5. Compare to existing code
diff my_code.c reference.c

# 6. Fix and validate
# ... fix firmware ...
make flash
percepta observe my-board
# Output: LED 'status': blue blinking ~0.99Hz ✓
```

---

### 22. Validate Power State Transitions

**Scenario:** Verify device enters low-power mode correctly.

**Steps:**

```bash
# 1. Normal operation
percepta observe my-device
# Output: LED 'power': green solid, LED 'activity': blue blinking 2Hz

percepta assert my-device "led power is ON" "led activity blinks at 2Hz"
# Output: ✓ Assertions passed

# 2. Trigger low-power mode (via GPIO or command)
# ... trigger sleep mode ...

# 3. Observe low-power state
percepta observe my-device
# Output: LED 'power': green solid, LED 'activity': OFF

percepta assert my-device "led power is ON" "led activity is OFF"
# Output: ✓ Assertions passed

# 4. Wake device
# ... wake signal ...

# 5. Verify return to normal operation
percepta observe my-device
percepta assert my-device "led activity blinks at 2Hz"
# Output: ✓ Assertion passed
```

---

### 23. Test Error Indicator Behavior

**Scenario:** Validate error LED behavior under fault conditions.

**Normal operation:**

```bash
percepta observe my-board
# Output: LED 'error': OFF

percepta assert my-board "led error is OFF"
# Output: ✓ Assertion passed
```

**Inject fault (e.g., disconnect sensor):**

```bash
# ... disconnect I2C sensor ...

percepta observe my-board
# Output: LED 'error': red blinking 5Hz

percepta assert my-board "led error is ON" "led error blinks at 5Hz"
# Output: ✓ Assertions passed
```

**Clear fault:**

```bash
# ... reconnect sensor ...

percepta observe my-board
# Output: LED 'error': OFF

percepta assert my-board "led error is OFF"
# Output: ✓ Assertion passed
```

---

### 24. Compare Production vs Debug Builds

**Scenario:** Verify production build behaves identically to debug build.

**Steps:**

```bash
# 1. Flash debug build
make debug
make flash
percepta device set-firmware my-board debug-v1.0
percepta observe my-board

# 2. Flash production build
make production
make flash
percepta device set-firmware my-board production-v1.0
percepta observe my-board

# 3. Compare
percepta diff my-board --from debug-v1.0 --to production-v1.0

# Expected: No differences
# If differences found, investigate optimization issues
```

**Monitoring timing differences:**

```bash
# Debug build boot time
percepta observe my-board
# Boot time: 1500ms

# Production build boot time
percepta observe my-board
# Boot time: 1200ms

# Diff shows:
# ~ Boot timing: 1500ms → 1200ms (MODIFIED)
# ✓ Expected - optimizations improved boot time
```

---

### 25. Hardware-in-Loop Testing

**Scenario:** Automated HIL test with external stimulus.

**Setup:**
- Device under test (DUT)
- Stimulus generator (GPIO, button press, sensor simulator)
- Percepta camera monitoring DUT

**Test script:**

```bash
#!/bin/bash
# HIL test: Button press triggers LED

echo "HIL Test: Button Press LED Trigger"

# 1. Baseline - no button press
percepta observe dut
percepta assert dut "led status is OFF"
echo "✓ Baseline: LED off"

# 2. Simulate button press (via GPIO expander or relay)
./stimulus.sh button press GPIO0

# 3. Observe LED response
sleep 0.5  # Debounce delay
percepta observe dut

# 4. Assert LED turned on
percepta assert dut "led status is ON" || {
  echo "FAIL: LED did not turn on after button press"
  exit 1
}
echo "✓ LED turned on after button press"

# 5. Release button
./stimulus.sh button release GPIO0

# 6. Observe LED turned off
sleep 0.5
percepta observe dut
percepta assert dut "led status is OFF" || {
  echo "FAIL: LED did not turn off after button release"
  exit 1
}
echo "✓ LED turned off after button release"

echo "✓ HIL test passed"
```

---

## Summary

These 25 examples demonstrate Percepta's versatility across:

**Development workflows:**
- Observation and validation
- Firmware tracking
- Code generation
- Style enforcement

**Testing scenarios:**
- Unit-level LED/display validation
- Integration testing (multi-signal)
- Regression detection
- Cross-version comparison

**Production deployment:**
- CI/CD pipelines
- Release validation
- Multi-board testing
- Hardware-in-loop automation

**Best practices:**

1. **Start simple:** Basic observe → assert workflow
2. **Track versions:** Tag firmware for diff capability
3. **Automate:** Integrate into CI/CD early
4. **Generate with validation:** Use AI generation + hardware validation loop
5. **Store patterns:** Build knowledge base of validated code

**Common patterns:**

```bash
# Development cycle
percepta observe <device>         # Check current state
percepta assert <device> <expr>   # Validate behavior
percepta diff <device> --from X --to Y  # Compare versions

# Code generation cycle
percepta generate <spec> --board <type>  # Generate code
# ... flash to hardware ...
percepta observe <device>         # Validate behavior
percepta knowledge store <spec> <file>   # Store validated pattern

# CI/CD pattern
make flash                        # Deploy firmware
percepta observe <device>         # Capture behavior
percepta assert <device> <expr>   # Validate requirements
percepta diff <device> --from baseline --to $SHA  # Detect regressions
```

For more details, see [commands.md](commands.md) for complete command reference.
