# Getting Started with Percepta

This guide walks you through your first observation, assertion, and firmware diff. By the end, you'll understand Percepta's core workflow and be ready to integrate it into your development cycle.

**Prerequisites:**
- Percepta installed ([Installation Guide](installation.md))
- ANTHROPIC_API_KEY environment variable set
- Webcam connected and accessible
- Hardware device with visible LEDs or display

**Time to complete:** ~10 minutes

---

## Step 1: Configure Your First Device

Add a device to Percepta's config:

```bash
percepta device add fpga
```

You'll be prompted for:
- **Device type:** e.g., `fpga`, `esp32`, `stm32` (human-readable identifier)
- **Camera device path:** e.g., `/dev/video0` (Linux), `0` (macOS/Windows built-in)
- **Firmware version:** (optional) e.g., `v1`, `main`, `abc123` (leave blank for now)

Example:

```bash
$ percepta device add fpga
Device type (e.g., fpga, esp32, stm32): fpga
Camera device path (default: /dev/video0): /dev/video0
Firmware version (optional, press Enter to skip):

✓ Device 'fpga' added successfully
```

**Verify configuration:**

```bash
percepta device list
```

Output:

```
Configured devices:

fpga
  Type: fpga
  Camera: /dev/video0
```

---

## Step 2: Test Camera Setup

Before observing hardware, verify your camera is working:

**Linux:**
```bash
# List available cameras
ls /dev/video*

# Test camera access (should not error)
v4l2-ctl --device=/dev/video0 --all
```

**macOS/Windows:**
- Open native camera app to verify camera is detected
- Camera index is typically `0` for built-in, `1` for external USB

**Adjust camera position:**
1. Point webcam directly at your hardware device
2. Ensure LEDs/display are clearly visible in frame
3. Avoid glare or backlighting (impacts color detection)
4. Recommended distance: 15-30cm from hardware

---

## Step 3: Your First Observation

With hardware powered on and visible to camera, run:

```bash
percepta observe fpga
```

**Expected output:**

```
Signals (3):
LED1: blue blinking ~2.0Hz
LED2: purple blinking ~0.8Hz
LED3: red solid

Observation complete (748ms)
```

**What just happened:**
1. Percepta captured a frame from `/dev/video0`
2. Sent the image to Claude Vision API
3. Parsed structured signals (LED state, color, blink frequency)
4. Stored observation in SQLite (`~/.local/share/percepta/percepta.db`)

**Troubleshooting:**
- **No signals detected:** Ensure LEDs are bright enough and in frame. Try adjusting camera position.
- **Wrong colors:** Verify no color filters on camera. Claude Vision works best with clear, direct lighting.
- **Slow response (>3s):** First API call can be slow. Subsequent calls are faster (~500-1000ms).

---

## Step 4: Your First Assertion

Assertions validate expected hardware behavior using a DSL (Domain-Specific Language).

**Example:** Verify LED1 is blinking and blue:

```bash
percepta assert fpga "led('LED1').blinks() && led('LED1').color_rgb(0,0,255)"
```

**Expected output:**

```
✓ All assertions passed
```

**If assertion fails:**

```bash
$ percepta assert fpga "led('LED1').color_rgb(255,0,0)"

✗ Assertion failed: expected LED1 color (255,0,0), got (0,0,255)
```

**Common assertion patterns:**

```bash
# LED is ON
percepta assert fpga "led('LED1').is_on()"

# LED is blinking at ~2Hz (±10% tolerance)
percepta assert fpga "led('LED1').blinks() && led('LED1').blink_hz(2.0, tolerance=0.2)"

# Display shows specific text
percepta assert fpga "display('LCD').contains('Ready')"

# Multiple conditions (AND)
percepta assert fpga "led('STATUS').is_on() && led('ERROR').is_off()"
```

**Assertion DSL Reference:**

**LED assertions:**
- `led('NAME').is_on()` — LED is ON
- `led('NAME').is_off()` — LED is OFF
- `led('NAME').blinks()` — LED is blinking
- `led('NAME').color_rgb(r, g, b)` — LED color matches RGB (±5 tolerance)
- `led('NAME').blink_hz(hz, tolerance=0.1)` — Blink frequency matches (±tolerance)

**Display assertions:**
- `display('NAME').contains('text')` — Display text contains substring (case-sensitive)

**Timing assertions:**
- `boot_time_ms() < 5000` — Boot time under 5 seconds

---

## Step 5: Firmware Tracking Workflow

Percepta tracks hardware behavior across firmware versions using manual tags.

**Workflow:**

**1. Tag initial firmware version:**

```bash
percepta device set-firmware fpga v1
```

**2. Run first observation:**

```bash
percepta observe fpga
```

This observation is now associated with firmware `v1`.

**3. Update firmware on your hardware:**

Flash new firmware to your device using your normal workflow (e.g., `idf.py flash`, `st-flash`, `openocd`, etc.).

**4. Tag new firmware version:**

```bash
percepta device set-firmware fpga v2
```

**5. Run second observation:**

```bash
percepta observe fpga
```

This observation is now associated with firmware `v2`.

**6. Compare versions:**

```bash
percepta diff fpga --from v1 --to v2
```

**Expected output:**

```
Comparing firmware versions:
FROM: v1 (2026-02-11 10:30:00)
TO:   v2 (2026-02-11 11:45:00)

Device: fpga

Changes detected:

+ LED2: purple blinking ~0.8Hz (ADDED)
- LED3: red solid (REMOVED)
~ LED1: blue blinking 2.0Hz → 2.5Hz (MODIFIED)

Summary: 1 added, 1 removed, 1 modified
```

**What the diff shows:**
- `+` — Signal exists in `v2` but not in `v1` (new LED turned on)
- `-` — Signal exists in `v1` but not in `v2` (LED no longer visible)
- `~` — Signal exists in both, but properties changed (color, blink rate, state)

**Exit codes:**
- `0` — No changes detected (identical behavior)
- `1` — Changes detected
- `2` — Error (device not found, firmware tag missing, etc.)

---

## Step 6: Integration into Development Workflow

**Typical workflow:**

```bash
# 1. Start with baseline firmware
percepta device set-firmware esp32 baseline
percepta observe esp32

# 2. Make firmware changes (edit code, rebuild, flash)
# ... (your normal development cycle) ...

# 3. Tag new version and observe
percepta device set-firmware esp32 test-feature-x
percepta observe esp32

# 4. Validate expected behavior
percepta assert esp32 "led('STATUS').blinks() && led('ERROR').is_off()"

# 5. Compare to baseline (detect regressions)
percepta diff esp32 --from baseline --to test-feature-x
```

**CI/CD integration example:**

```bash
#!/bin/bash
# Flash firmware
idf.py flash

# Wait for boot
sleep 2

# Validate hardware state
percepta device set-firmware esp32 ${CI_COMMIT_SHA}
percepta observe esp32
percepta assert esp32 "led('STATUS').is_on()"

if [ $? -eq 0 ]; then
  echo "Hardware validation passed"
else
  echo "Hardware validation failed"
  exit 1
fi
```

---

## Troubleshooting

**"Device 'xyz' not found"**
- Run `percepta device list` to see configured devices
- Add device with `percepta device add xyz`

**"No signals detected"**
- Verify hardware is powered and LEDs/display are visible
- Check camera position and lighting
- Try running observation multiple times (some LEDs may be blinking off during capture)

**"Assertion timeout: signal 'LEDX' not found"**
- Signal name doesn't match observation output
- Check exact LED name with `percepta observe <device>` first
- LED names are case-insensitive but must match identifier (LED1, LED2, etc.)

**"Diff shows unexpected changes"**
- Single-frame capture may miss blinking LEDs (known limitation: ISS-001)
- Blink frequency can vary slightly (~0.1Hz) across observations
- Verify firmware tags are correct: `percepta device list`

**"Observation is slow (>3s)"**
- First API call is slower (cold start)
- Check internet connection
- Verify Claude API key has credits/quota

---

## Next Steps

**Explore examples:**
- [ESP32 with status LED](examples/esp32.yaml)
- [STM32 with LCD display](examples/stm32.yaml)
- [FPGA dev board with multiple LEDs](examples/fpga.yaml)

**Advanced topics:**
- DSL assertion syntax (full reference coming soon)
- Multi-device configurations
- Observation history querying

**Join alpha testing:**
- Report issues: [GitHub Issues](https://github.com/Perceptax/percepta/issues)
- Share feedback: [GitHub Discussions](https://github.com/Perceptax/percepta/discussions)

---

**You're ready to use Percepta!** Close the feedback loop and watch your firmware run.
