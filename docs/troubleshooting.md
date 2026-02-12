# Troubleshooting Guide

Solutions to common issues with Percepta.

---

## Table of Contents

**Installation & Setup:**
- [Command Not Found](#command-not-found)
- [API Key Errors](#api-key-errors)
- [Config File Issues](#config-file-issues)

**Camera Issues:**
- [Camera Not Found](#camera-not-found)
- [Permission Denied](#permission-denied-camera)
- [Wrong Camera Selected](#wrong-camera-selected)

**Observation Issues:**
- [No Signals Detected](#no-signals-detected)
- [Low Confidence Scores](#low-confidence-scores)
- [Slow Observations](#slow-observations)
- [Missing LEDs](#missing-leds)
- [Incorrect Blink Rate](#incorrect-blink-rate)
- [Display OCR Failures](#display-ocr-failures)

**Assertion Issues:**
- [Assertion Timeout](#assertion-timeout)
- [Assertion Fails Unexpectedly](#assertion-fails-unexpectedly)
- [Tolerance Issues](#tolerance-issues)

**Diff Issues:**
- [Firmware Tag Not Found](#firmware-tag-not-found)
- [Unexpected Changes](#unexpected-changes)
- [Missing Observations](#missing-observations)

**Code Generation Issues:**
- [Generation Failed](#generation-failed)
- [Style Violations Remain](#style-violations-remain)
- [Pattern Not Stored](#pattern-not-stored)

**Storage Issues:**
- [Storage Init Failed](#storage-init-failed)
- [Database Locked](#database-locked)
- [Disk Space Issues](#disk-space-issues)

---

## Installation & Setup

### Command Not Found

**Error:**
```bash
$ percepta --help
bash: percepta: command not found
```

**Cause:** Binary not in PATH.

**Solution:**

**Check if binary exists:**
```bash
ls -l /usr/local/bin/percepta
```

If not found, reinstall:

```bash
# Linux
sudo mv percepta /usr/local/bin/
sudo chmod +x /usr/local/bin/percepta

# macOS
sudo mv percepta /usr/local/bin/
sudo chmod +x /usr/local/bin/percepta

# Windows
# Move percepta.exe to a directory in PATH
# e.g., C:\Program Files\Percepta\
```

**Verify installation:**
```bash
percepta --help
```

---

### API Key Errors

**Error:**
```
Error: ANTHROPIC_API_KEY not set

Suggestion: Set ANTHROPIC_API_KEY environment variable or add to ~/.config/percepta/config.yaml
```

**Cause:** Claude API key not configured.

**Solution:**

**Set environment variable:**

```bash
# Linux/macOS (temporary)
export ANTHROPIC_API_KEY="sk-ant-api03-..."

# Add to shell profile for persistence
echo 'export ANTHROPIC_API_KEY="sk-ant-api03-..."' >> ~/.bashrc
source ~/.bashrc

# Windows (PowerShell)
$env:ANTHROPIC_API_KEY="sk-ant-api03-..."

# Windows (Command Prompt)
set ANTHROPIC_API_KEY=sk-ant-api03-...
```

**Verify:**
```bash
echo $ANTHROPIC_API_KEY  # Should show your key
```

**Get API key:**
1. Sign up at [console.anthropic.com](https://console.anthropic.com)
2. Navigate to API Keys
3. Create a new key
4. Copy and set as environment variable

**Check API key validity:**

```bash
# Test with simple observation
percepta observe my-board

# If authentication error:
# Error: API call failed: authentication error
# → API key is invalid or expired
```

**Invalid API key errors:**

**Error:**
```
Error: API call failed: authentication error
```

**Solutions:**
1. Verify key starts with `sk-ant-`
2. Check for extra spaces or quotes
3. Regenerate key in Anthropic console
4. Check account has active credits/quota

---

### Config File Issues

**Error:**
```
Error: No config file found at ~/.config/percepta/config.yaml
```

**Cause:** Config file doesn't exist yet.

**Solution:**

Create your first device (auto-creates config):
```bash
percepta device add my-board
```

**Manual config creation:**

```bash
mkdir -p ~/.config/percepta
cat > ~/.config/percepta/config.yaml <<EOF
devices:
  my-board:
    camera: /dev/video0
    firmware: v1.0

vision:
  frames: 5
  interval: 200ms
  confidence_threshold: 0.7
EOF
```

**Verify:**
```bash
percepta device list
```

---

## Camera Issues

### Camera Not Found

**Error:**
```
Error: Camera '/dev/video0' not found

Suggestion: Check available cameras: 'ls /dev/video*' (Linux) or try camera index 0-2 (macOS/Windows)
```

**Cause:** Camera device doesn't exist or incorrect path.

**Solution:**

**Linux:**

List available cameras:
```bash
ls -l /dev/video*
```

Output:
```
/dev/video0
/dev/video1
/dev/video2
```

Identify which is your camera:
```bash
v4l2-ctl --list-devices
```

Update device config:
```bash
percepta device add my-board
# Enter correct camera path when prompted
```

**macOS/Windows:**

Try camera indices 0-2:
```bash
# Try index 0 (usually built-in)
percepta device add my-board
# Enter: 0

# If that doesn't work, try 1
# ... edit ~/.config/percepta/config.yaml ...
# camera: 1
```

**USB camera not detected:**

1. Check USB connection
2. Check camera LED (if present) - should be lit
3. Test camera with native app (Photo Booth on macOS, Camera on Windows)
4. Replug USB cable
5. Try different USB port

---

### Permission Denied (Camera)

**Error:**
```
Error: failed to open camera: /dev/video0: permission denied
```

**Cause:** User doesn't have permission to access camera device.

**Solution:**

**Linux:**

Add user to `video` group:
```bash
sudo usermod -a -G video $USER
```

**Important:** Log out and log back in for changes to take effect.

Verify group membership:
```bash
groups | grep video
```

Alternative (temporary):
```bash
sudo chmod 666 /dev/video0
```

**macOS:**

Grant camera permissions:
1. System Preferences → Security & Privacy → Camera
2. Enable checkbox for Terminal (or your IDE)
3. Restart application

**Windows:**

Grant camera permissions:
1. Settings → Privacy → Camera
2. Enable "Allow apps to access your camera"
3. Enable for Command Prompt / PowerShell

---

### Wrong Camera Selected

**Problem:** Observing wrong camera (not pointed at hardware).

**Solution:**

**Identify cameras:**

**Linux:**
```bash
v4l2-ctl --list-devices

# Test each camera
for cam in /dev/video*; do
  echo "Testing $cam"
  ffmpeg -f v4l2 -i $cam -frames 1 test-$cam.jpg 2>&1 | grep -i "video"
done

# View captured frames to identify correct camera
```

**Update device config:**

```bash
# Re-add device with correct camera
percepta device add my-board
# Enter correct camera path

# Or edit config manually
nano ~/.config/percepta/config.yaml
```

```yaml
devices:
  my-board:
    camera: /dev/video2  # Update to correct camera
```

---

## Observation Issues

### No Signals Detected

**Output:**
```
✓ Observation captured: obs-12345
Device: my-board
Timestamp: 2026-02-13T10:30:00Z

No signals detected
```

**Causes:**
1. Hardware not visible to camera
2. LEDs too dim
3. Confidence threshold too high
4. Poor lighting conditions

**Solutions:**

**1. Check camera positioning:**

```bash
# Verify camera is pointed at hardware
# Distance: 15-30cm recommended
# Ensure LEDs/display are in frame
```

**2. Check hardware is powered:**

Verify LEDs are actually on/blinking.

**3. Improve lighting:**

- Avoid backlighting (light behind hardware)
- Use diffuse lighting (not direct harsh light)
- Avoid glare on displays

**4. Lower confidence threshold:**

Edit `~/.config/percepta/config.yaml`:
```yaml
vision:
  confidence_threshold: 0.5  # Lower from default 0.7
```

**5. Capture more frames:**

```yaml
vision:
  frames: 10  # Increase from default 5
  interval: 200ms
```

**6. Test with known-good hardware:**

Use a dev board with bright, obvious LEDs to verify camera setup is working.

---

### Low Confidence Scores

**Output:**
```
Signals (2):
  1. LED 'LED1': ON [confidence: 0.52]
  2. LED 'LED2': OFF [confidence: 0.48]
```

**Causes:**
- LEDs partially occluded
- Dim LEDs
- Poor lighting
- Camera out of focus

**Solutions:**

**1. Improve lighting conditions:**

- Use brighter ambient light
- Avoid shadows on hardware
- Position light source to side (not behind camera or hardware)

**2. Adjust camera:**

- Ensure hardware is in focus
- Move closer (15-20cm)
- Clean camera lens

**3. Check LED brightness:**

- Verify LED resistors aren't limiting brightness too much
- Test with brighter LEDs

**4. Lower threshold temporarily:**

```yaml
vision:
  confidence_threshold: 0.5
```

**Note:** Low confidence scores may indicate real ambiguity. If score < 0.6, verify LED state manually.

---

### Slow Observations

**Problem:** Observations take >5 seconds.

**Typical timing:** 1-2 seconds for 5 frames

**Causes:**
1. First API call (cold start)
2. Slow internet connection
3. Too many frames
4. API rate limiting

**Solutions:**

**1. First call slowness (expected):**

First observation is slower (~3-5 seconds) due to API cold start. Subsequent observations are faster (~1-2 seconds).

```bash
# First observation
percepta observe my-board  # ~3s

# Second observation
percepta observe my-board  # ~1s
```

**2. Reduce frames:**

```yaml
vision:
  frames: 3       # Reduce from 5
  interval: 100ms # Faster sampling
```

**3. Check internet connection:**

```bash
# Test API latency
time curl -s https://api.anthropic.com/v1/health > /dev/null
```

If >2 seconds, check network.

**4. API quota issues:**

Check Anthropic console for rate limits or quota exhaustion.

---

### Missing LEDs

**Problem:** Observation shows 2 LEDs, but hardware has 3.

**Causes:**
1. LED off during single-frame capture
2. LED too dim
3. LED occluded
4. Blinking LED missed (see ISS-001)

**Solutions:**

**1. Capture more frames:**

```yaml
vision:
  frames: 10      # Increase frames
  interval: 200ms # Longer observation window
```

Multi-frame capture detects blinking LEDs.

**2. Check LED is actually on:**

Visually verify all LEDs are lit.

**3. Ensure LEDs are visible:**

- Not blocked by other components
- Camera angle captures all LEDs
- LEDs not behind reflective surface

**4. Lower confidence threshold:**

```yaml
vision:
  confidence_threshold: 0.6
```

---

### Incorrect Blink Rate

**Output:**
```
LED 'status': blue blinking ~0.47Hz
```

**Expected:** 1Hz

**Causes:**
1. Actual firmware issue (not observing correctly)
2. Observation window too short
3. LED duty cycle issues

**Solutions:**

**1. Verify firmware code:**

Check timer configuration, delay values:
```c
// Ensure timer configured for correct frequency
// 1Hz = toggle every 500ms
```

**2. Extend observation window:**

```yaml
vision:
  frames: 10
  interval: 500ms  # Longer window for accurate frequency measurement
```

**3. Multiple observations:**

```bash
# Take several observations
for i in {1..5}; do
  echo "Observation $i:"
  percepta observe my-board | grep "blinks"
  sleep 1
done
```

If consistent, it's a firmware issue. If varies, increase frames/interval.

**4. Check duty cycle:**

Very short ON or OFF periods may appear as different frequency.

---

### Display OCR Failures

**Output:**
```
Display 'LCD': "" [confidence: 0.35]
```

**Expected:** Display shows "Ready"

**Causes:**
1. Display text too small
2. Display contrast too low
3. Glare on display
4. Display not in focus

**Solutions:**

**1. Camera positioning:**

- Move closer to display (10-15cm)
- Ensure display is in focus
- Angle camera perpendicular to display (not at angle)

**2. Lighting:**

- Avoid glare on LCD surface
- Use diffuse lighting
- Turn off overhead lights if causing glare

**3. Display contrast:**

Increase LCD contrast in firmware if possible.

**4. Use partial matching:**

If OCR is noisy, use `contains` instead of exact match:

```bash
# Instead of exact match
percepta assert my-board "display LCD shows 'Ready'"

# Use partial match
percepta assert my-board "display LCD contains 'Rdy'"
```

**5. Capture more frames:**

```yaml
vision:
  frames: 10
  interval: 200ms
```

Multiple frames improve OCR stability.

---

## Assertion Issues

### Assertion Timeout

**Error:**
```
Error: Assertion timeout: signal 'LED3' not found in observation

Suggestion: Run 'percepta observe <device>' first to see available signals
```

**Cause:** Signal name doesn't match observation output.

**Solution:**

**1. Check actual signal names:**

```bash
percepta observe my-board
```

Output:
```
Signals (2):
  1. LED 'LED1': ON
  2. LED 'LED2': blinking
```

**2. Use correct name:**

```bash
# LED names are LED1, LED2, not LED3
percepta assert my-board "led LED1 is ON"
```

**3. Case sensitivity:**

LED names are case-insensitive but must match identifier:
- Correct: `led LED1`, `led led1`
- Incorrect: `led Status`, `led Power` (unless named exactly)

---

### Assertion Fails Unexpectedly

**Output:**
```
✗ Assertion failed

Expected: LED 'LED1' is ON
Actual:   LED 'LED1' is OFF
Confidence: 0.85
```

**Causes:**
1. Hardware actually off (firmware issue)
2. LED off during capture (blinks too fast)
3. Low confidence detection

**Solutions:**

**1. Visual verification:**

Look at hardware - is LED actually on?

**2. Run observation first:**

```bash
percepta observe my-board
```

Check LED state before asserting.

**3. Capture more frames:**

For blinking LEDs:
```yaml
vision:
  frames: 10
  interval: 100ms
```

**4. Check confidence:**

If confidence < 0.7, observation may be unreliable.

---

### Tolerance Issues

**Problem:** Assertion fails due to minor variance.

**Error:**
```
✗ Assertion failed

Expected: LED 'LED1' blinks at 2.0 Hz
Actual:   LED 'LED1' blinks at 1.87 Hz
```

**Cause:** Hardware timing variance within acceptable range.

**Solution:**

Assertions have built-in tolerances:
- Blink rate: ±10% default
- RGB color: ±5 per channel

**If variance is expected:**

Accept minor differences - 1.87Hz is within 10% of 2.0Hz (1.8-2.2Hz range).

**If variance is too large:**

Indicates firmware timing issue. Check:
- Timer configuration
- Clock source accuracy
- System load affecting timing

---

## Diff Issues

### Firmware Tag Not Found

**Error:**
```
Error: failed to get observation for firmware 'v1.0': no observations found
```

**Cause:** No observations captured with that firmware tag.

**Solution:**

**1. Check available observations:**

```bash
sqlite3 ~/.local/share/percepta/percepta.db
sqlite> SELECT DISTINCT firmware FROM observations WHERE device_id='my-board';
```

**2. Verify firmware tag:**

```bash
percepta device list
```

Check current firmware tag matches.

**3. Capture observation with correct tag:**

```bash
percepta device set-firmware my-board v1.0
percepta observe my-board
```

Now diff will work:
```bash
percepta diff my-board --from v1.0 --to v1.1
```

---

### Unexpected Changes

**Output:**
```
Changes detected:

~ LED1: blue blinking 1.98Hz → 2.02Hz (MODIFIED)
```

**Cause:** Minor timing variance, not actual firmware change.

**Solution:**

**Understand diff behavior:**

Diff shows exact differences with minimal tolerance:
- BlinkHz: normalized to 0.01Hz precision
- RGB: exact values
- ON/OFF state: exact

**Variance is expected:**

Blink rates vary slightly (±0.1Hz) due to:
- Timer accuracy
- System load
- Measurement precision

**If change is significant:**

~ LED1: 1Hz → 5Hz (MODIFIED)

This indicates real firmware change - investigate.

**If change is minor:**

~ LED1: 1.98Hz → 2.02Hz (MODIFIED)

This is hardware variance, not a bug.

---

### Missing Observations

**Error:**
```
Error: failed to get observation for firmware 'v2.0': no observations found
```

**Cause:** Forgot to observe after flashing new firmware.

**Solution:**

**Workflow reminder:**

```bash
# 1. Tag firmware
percepta device set-firmware my-board v2.0

# 2. MUST observe (don't skip)
percepta observe my-board

# 3. Now diff works
percepta diff my-board --from v1.0 --to v2.0
```

**Check existing observations:**

```bash
sqlite3 ~/.local/share/percepta/percepta.db
sqlite> SELECT device_id, firmware, timestamp FROM observations ORDER BY timestamp DESC LIMIT 10;
```

---

## Code Generation Issues

### Generation Failed

**Error:**
```
Error: Code generation failed: API request failed: timeout
```

**Causes:**
1. Network timeout
2. API rate limit
3. Invalid specification

**Solutions:**

**1. Check API key:**

```bash
echo $ANTHROPIC_API_KEY
```

**2. Check network:**

```bash
curl -s https://api.anthropic.com/v1/health
```

**3. Retry:**

Code generation can timeout occasionally. Simply retry:
```bash
percepta generate "Blink LED at 1Hz" --board esp32 --output led.c
```

**4. Simplify specification:**

If spec is too complex, break into smaller pieces:
```bash
# Instead of:
percepta generate "Read sensor, process data, display on LCD, trigger alarm if threshold exceeded"

# Break into steps:
percepta generate "Read I2C sensor and return value" --board esp32 --output sensor.c
percepta generate "Display value on LCD" --board esp32 --output display.c
```

---

### Style Violations Remain

**Output:**
```
Style check: FAILED (3 violations)
Remaining violations:
  - magic-number: 2 violations
  - const-correctness: 1 violation

Suggestion: Fix manually or run 'percepta style-check --fix <file>'
```

**Cause:** Some violations require manual review (magic numbers, const correctness).

**Solution:**

**1. Auto-fix deterministic violations:**

```bash
percepta style-check led.c --fix
```

This fixes:
- Naming conventions
- Type usage (int → uint8_t)

**2. Fix magic numbers manually:**

Before:
```c
uint32_t delay = 500;  // Magic number
```

After:
```c
const uint32_t BLINK_DELAY_MS = 500;
uint32_t delay = BLINK_DELAY_MS;
```

**3. Fix const correctness manually:**

Before:
```c
void process(uint8_t* data) { /* ... */ }
```

After:
```c
void process(const uint8_t* data) { /* ... */ }
```

**4. Re-check:**

```bash
percepta style-check led.c
```

---

### Pattern Not Stored

**Output:**
```
✗ Pattern storage FAILED: no observation found for device+firmware
```

**Cause:** Missing hardware validation (observation).

**Solution:**

**Requirements for storing pattern:**
1. Code is BARR-C compliant
2. Observation exists for device+firmware combination

**Workflow:**

```bash
# 1. Flash code
# ... compile and flash ...

# 2. Tag firmware
percepta device set-firmware my-esp32 v1.0

# 3. Observe hardware
percepta observe my-esp32

# 4. NOW store pattern
percepta knowledge store "Blink LED" led.c --device my-esp32 --firmware v1.0

# Output: ✓ Pattern stored
```

---

## Storage Issues

### Storage Init Failed

**Error:**
```
Error: Failed to initialize storage: permission denied

Suggestion: Check that ~/.local/share/percepta/ directory is writable
```

**Cause:** Insufficient permissions for storage directory.

**Solution:**

```bash
# Create directory with correct permissions
mkdir -p ~/.local/share/percepta
chmod 755 ~/.local/share/percepta

# Verify
ls -ld ~/.local/share/percepta
```

**Custom storage path:**

Edit `~/.config/percepta/config.yaml`:
```yaml
storage:
  path: /path/with/permissions/percepta.db
```

---

### Database Locked

**Error:**
```
Error: database locked
```

**Cause:** Another Percepta process is accessing database.

**Solution:**

**1. Check for running processes:**

```bash
ps aux | grep percepta
```

Kill any stuck processes:
```bash
kill <pid>
```

**2. Close database connections:**

```bash
# Find processes with database open
lsof ~/.local/share/percepta/percepta.db

# Kill if necessary
```

**3. Restart:**

```bash
percepta observe my-board
```

---

### Disk Space Issues

**Error:**
```
Error: failed to save observation: disk full
```

**Cause:** Insufficient disk space.

**Solution:**

**1. Check disk space:**

```bash
df -h ~/.local/share/percepta/
```

**2. Clean old observations:**

```bash
# Backup first
cp ~/.local/share/percepta/percepta.db ~/backups/

# Remove database (observations will be lost)
rm ~/.local/share/percepta/percepta.db

# Database recreated on next observation
percepta observe my-board
```

**3. Archive old observations:**

```bash
# Export observations to JSON
sqlite3 ~/.local/share/percepta/percepta.db << EOF
.mode json
.output observations-backup.json
SELECT * FROM observations WHERE timestamp < datetime('now', '-30 days');
.quit
EOF

# Delete old observations
sqlite3 ~/.local/share/percepta/percepta.db << EOF
DELETE FROM observations WHERE timestamp < datetime('now', '-30 days');
VACUUM;
.quit
EOF
```

---

## Getting Help

If you encounter an issue not covered here:

**1. Check logs:**

Percepta prints detailed error messages. Read them carefully for suggestions.

**2. Enable debug output:**

```bash
# Set verbose mode (if available)
PERCEPTA_DEBUG=1 percepta observe my-board
```

**3. GitHub Issues:**

Report bugs with:
- Error message (full output)
- Command run
- Platform (Linux/macOS/Windows)
- Percepta version: `percepta --version`

[Create Issue](https://github.com/Perceptax/percepta/issues/new)

**4. Community Discussion:**

Ask questions: [GitHub Discussions](https://github.com/Perceptax/percepta/discussions)

---

## See Also

- [Installation Guide](installation.md) - Setup
- [Configuration Guide](configuration.md) - Config options
- [Commands Reference](commands.md) - All commands
- [Examples](examples.md) - Usage examples
