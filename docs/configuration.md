# Configuration Guide

Complete reference for Percepta configuration options.

---

## Configuration File

**Location:** `~/.config/percepta/config.yaml`

Percepta uses YAML configuration for device definitions, vision settings, and storage paths.

**Auto-creation:**

Config file is created automatically when you add your first device:

```bash
percepta device add my-board
```

This creates `~/.config/percepta/config.yaml` if it doesn't exist.

---

## Example Configuration

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

  test-board:
    type: arduino
    camera: 0
    firmware: baseline

vision:
  frames: 5
  interval: 200ms
  confidence_threshold: 0.7

storage:
  path: ~/.local/share/percepta/percepta.db

knowledge:
  path: ~/.local/share/percepta/knowledge.db
```

---

## Configuration Sections

### Devices

**Format:**

```yaml
devices:
  <device-id>:
    type: <board-type>
    camera: <camera-path>
    firmware: <firmware-tag>
```

**Fields:**

**`device-id`** (required)
- Unique identifier for the device
- Used in all commands: `percepta observe <device-id>`
- Can be any string (no spaces recommended)
- Examples: `my-esp32`, `lab-board-1`, `test-device`

**`type`** (optional)
- Human-readable board type
- Not used functionally, just for documentation
- Examples: `esp32`, `stm32`, `arduino`, `fpga`

**`camera`** (required)
- Camera device path or index
- **Linux:** `/dev/video0`, `/dev/video1`, etc.
- **macOS:** `0` (built-in), `1` (external USB)
- **Windows:** `0`, `1`, `2` (camera index)

**`firmware`** (optional)
- Current firmware version tag
- Used for diffing: `percepta diff --from <tag> --to <tag>`
- Can be any string: `v1.0`, `baseline`, `abc123`, `feature-x`
- Set via: `percepta device set-firmware <device> <tag>`

**Examples:**

```yaml
# Minimal device (camera only)
devices:
  quick-test:
    camera: /dev/video0

# Full device config
devices:
  production-board:
    type: esp32-devkit
    camera: /dev/video1
    firmware: v2.1.3-release

# Multiple devices
devices:
  board-a:
    type: stm32f4
    camera: /dev/video0
    firmware: v1.0

  board-b:
    type: esp32
    camera: /dev/video2
    firmware: v1.0

  ci-test-board:
    type: generic
    camera: 0
    firmware: ci-build-1234
```

---

### Vision

Controls computer vision behavior (frame capture, confidence thresholds).

**Format:**

```yaml
vision:
  frames: <number>
  interval: <duration>
  confidence_threshold: <float>
```

**Fields:**

**`frames`** (optional, default: 5)
- Number of frames to capture per observation
- Range: 1-20
- Higher values = more reliable detection, slower observations
- Recommended: 5 for most use cases

**`interval`** (optional, default: 200ms)
- Time between frame captures
- Format: `<number>ms` or `<number>s`
- Examples: `100ms`, `200ms`, `1s`
- Total observation time = `frames * interval` (e.g., 5 frames * 200ms = 1 second)

**`confidence_threshold`** (optional, default: 0.7)
- Minimum confidence score for signals
- Range: 0.0 to 1.0
- Lower = more signals detected (may include false positives)
- Higher = fewer signals (may miss dim/occluded LEDs)
- Recommended: 0.7 for most lighting conditions

**Examples:**

```yaml
# Fast observation (fewer frames)
vision:
  frames: 3
  interval: 100ms
  confidence_threshold: 0.7

# Thorough observation (more frames, longer interval)
vision:
  frames: 10
  interval: 500ms
  confidence_threshold: 0.8

# Low-light conditions (lower threshold)
vision:
  frames: 5
  interval: 200ms
  confidence_threshold: 0.5
```

**Tuning guidelines:**

**Fast LEDs (>5Hz blink rate):**
```yaml
vision:
  frames: 10      # Capture more frames
  interval: 100ms # Faster sampling
```

**Slow observations:**
```yaml
vision:
  frames: 3       # Fewer frames
  interval: 100ms # Quick capture
```

**Low confidence issues:**
```yaml
vision:
  frames: 5
  interval: 200ms
  confidence_threshold: 0.5  # Lower threshold
```

---

### Storage

Controls SQLite database location for observations.

**Format:**

```yaml
storage:
  path: <database-path>
```

**Fields:**

**`path`** (optional, default: `~/.local/share/percepta/percepta.db`)
- SQLite database file path
- Stores all observations, timestamps, firmware tags
- Can use `~` for home directory
- Directory must exist (Percepta creates file if needed)

**Examples:**

```yaml
# Default location
storage:
  path: ~/.local/share/percepta/percepta.db

# Custom location
storage:
  path: /var/lib/percepta/observations.db

# Project-specific database
storage:
  path: ~/projects/my-firmware/percepta.db
```

**Database management:**

**Backup database:**
```bash
cp ~/.local/share/percepta/percepta.db ~/backups/percepta-$(date +%Y%m%d).db
```

**Clear all observations:**
```bash
rm ~/.local/share/percepta/percepta.db
# Database will be recreated on next observation
```

**Inspect database:**
```bash
sqlite3 ~/.local/share/percepta/percepta.db
sqlite> .tables
sqlite> SELECT device_id, firmware, timestamp FROM observations LIMIT 5;
```

---

### Knowledge

Controls knowledge graph database location (validated patterns).

**Format:**

```yaml
knowledge:
  path: <database-path>
```

**Fields:**

**`path`** (optional, default: `~/.local/share/percepta/knowledge.db`)
- SQLite database for knowledge graph
- Stores validated patterns, relationships, embeddings
- Separate from observations database

**Examples:**

```yaml
# Default location
knowledge:
  path: ~/.local/share/percepta/knowledge.db

# Custom location
knowledge:
  path: ~/my-patterns/knowledge.db

# Shared knowledge base (team)
knowledge:
  path: /shared/percepta/team-knowledge.db
```

**Knowledge graph usage:**

```bash
# Store validated pattern
percepta knowledge store "Blink LED" led.c --device my-esp32 --firmware v1.0

# Search patterns
percepta knowledge search "button debounce"

# List all patterns
percepta knowledge list
```

---

## Environment Variables

Environment variables override config file settings (useful for CI/CD).

### API Keys

**`ANTHROPIC_API_KEY`** (required for observation and generation)
- Claude API key from [console.anthropic.com](https://console.anthropic.com)
- Used for Claude Vision API (observation) and code generation
- Format: `sk-ant-...`

**`OPENAI_API_KEY`** (optional for semantic search)
- OpenAI API key for embeddings
- Used for semantic pattern search in knowledge graph
- Format: `sk-...`
- Percepta gracefully degrades without it (no semantic search)

**Examples:**

```bash
# Linux/macOS (add to ~/.bashrc or ~/.zshrc)
export ANTHROPIC_API_KEY="sk-ant-api03-..."
export OPENAI_API_KEY="sk-..."

# Windows (PowerShell)
$env:ANTHROPIC_API_KEY="sk-ant-api03-..."
$env:OPENAI_API_KEY="sk-..."

# Windows (Command Prompt)
set ANTHROPIC_API_KEY=sk-ant-api03-...
set OPENAI_API_KEY=sk-...
```

**CI/CD secrets:**

GitHub Actions:
```yaml
env:
  ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
  OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
```

GitLab CI:
```yaml
variables:
  ANTHROPIC_API_KEY: $ANTHROPIC_API_KEY
  OPENAI_API_KEY: $OPENAI_API_KEY
```

### Path Overrides

**`PERCEPTA_CONFIG_PATH`** (optional)
- Override config file location
- Default: `~/.config/percepta/config.yaml`

**Example:**
```bash
export PERCEPTA_CONFIG_PATH=/custom/path/config.yaml
percepta observe my-board
```

**`PERCEPTA_STORAGE_PATH`** (optional)
- Override observations database location
- Default: `~/.local/share/percepta/percepta.db`

**Example:**
```bash
export PERCEPTA_STORAGE_PATH=/tmp/test-observations.db
percepta observe my-board
```

---

## Camera Configuration

Camera selection and troubleshooting.

### Linux

**List available cameras:**
```bash
ls -l /dev/video*
```

**Output:**
```
/dev/video0
/dev/video1
/dev/video2
```

**Identify camera:**
```bash
v4l2-ctl --list-devices
```

**Output:**
```
Integrated Camera (usb-0000:00:14.0):
	/dev/video0
	/dev/video1

USB Webcam (usb-0000:00:14.0):
	/dev/video2
	/dev/video3
```

**Permissions:**

Add user to `video` group if permission denied:
```bash
sudo usermod -a -G video $USER
# Log out and log back in
```

**Config:**
```yaml
devices:
  my-board:
    camera: /dev/video0  # Built-in camera
```

### macOS

**Camera index:**
- Built-in camera: `0`
- External USB camera: `1`

**Config:**
```yaml
devices:
  my-board:
    camera: 0  # Built-in FaceTime camera
```

**Permissions:**

macOS will prompt for camera permissions on first use. Grant access to Terminal or your IDE.

### Windows

**Camera index:**
- Built-in camera: `0`
- External USB camera: `1`, `2`, etc.

**Config:**
```yaml
devices:
  my-board:
    camera: 0  # Built-in camera
```

**Permissions:**

Windows will prompt for camera permissions. Grant access in Settings → Privacy → Camera.

---

## Multi-Device Setup

Configuration for multiple boards.

**Example:**

```yaml
devices:
  # Development boards (local desk)
  esp32-dev:
    type: esp32-devkit
    camera: /dev/video0
    firmware: dev-branch

  stm32-proto:
    type: stm32f4-disco
    camera: /dev/video1
    firmware: v2.0-rc1

  # Test boards (lab bench)
  lab-board-1:
    type: esp32
    camera: /dev/video2
    firmware: v1.5.0

  lab-board-2:
    type: stm32
    camera: /dev/video3
    firmware: v1.5.0

  # CI test board (automated test rig)
  ci-test-board:
    type: generic
    camera: 0
    firmware: ci-build-auto

vision:
  frames: 5
  interval: 200ms
  confidence_threshold: 0.7

storage:
  path: ~/.local/share/percepta/percepta.db

knowledge:
  path: ~/.local/share/percepta/knowledge.db
```

**Usage:**

```bash
# Work on development board
percepta observe esp32-dev

# Test on lab boards
percepta observe lab-board-1
percepta observe lab-board-2

# CI testing
percepta observe ci-test-board
```

---

## Configuration Best Practices

**1. Version control your config (without secrets):**

`.gitignore`:
```
# Percepta databases
*.db

# Percepta config (if contains machine-specific paths)
config.yaml
```

**Template config (commit this):**

`config.template.yaml`:
```yaml
devices:
  my-board:
    type: esp32
    camera: /dev/video0  # Update for your machine
    firmware: v1.0

vision:
  frames: 5
  interval: 200ms
  confidence_threshold: 0.7
```

**2. Use firmware tags consistently:**

```yaml
devices:
  my-board:
    firmware: v1.2.3  # Semantic versioning
```

Update before observations:
```bash
percepta device set-firmware my-board v1.2.3
percepta observe my-board
```

**3. Separate databases per project:**

```yaml
storage:
  path: ~/projects/my-firmware/observations.db

knowledge:
  path: ~/projects/my-firmware/knowledge.db
```

**4. Tune vision settings per use case:**

**Fast CI testing:**
```yaml
vision:
  frames: 3
  interval: 100ms
  confidence_threshold: 0.7
```

**Thorough validation:**
```yaml
vision:
  frames: 10
  interval: 200ms
  confidence_threshold: 0.8
```

**5. Document device types:**

```yaml
devices:
  board-a:
    type: esp32-devkit-v1  # Specific model for team reference
    camera: /dev/video0
    firmware: v2.0

  board-b:
    type: stm32f407-discovery  # Clear identification
    camera: /dev/video1
    firmware: v2.0
```

---

## Troubleshooting Configuration

**Config not found:**

```bash
percepta observe my-board
# Error: config load failed: no config file

# Solution: Create device
percepta device add my-board
```

**Device not found:**

```bash
percepta observe unknown-board
# Error: Device 'unknown-board' not found in config

# Solution: Check configured devices
percepta device list

# Add device if needed
percepta device add unknown-board
```

**Camera not found:**

```bash
percepta observe my-board
# Error: Camera '/dev/video0' not found

# Solution: Check available cameras
ls /dev/video*

# Update device config
percepta device add my-board  # Re-add with correct camera
```

**Vision settings too aggressive:**

```bash
percepta observe my-board
# Output: No signals detected

# Solution: Lower confidence threshold
```

Edit `~/.config/percepta/config.yaml`:
```yaml
vision:
  confidence_threshold: 0.5  # Lower from 0.7
```

**Storage permission denied:**

```bash
percepta observe my-board
# Error: storage init failed: permission denied

# Solution: Check directory permissions
ls -ld ~/.local/share/percepta/
mkdir -p ~/.local/share/percepta/
chmod 755 ~/.local/share/percepta/
```

---

## Advanced Configuration

### Custom Vision Settings Per Device

Currently, vision settings are global. To use different settings per device, use separate config files:

**Project A:**
```bash
export PERCEPTA_CONFIG_PATH=~/projects/projectA/config.yaml
percepta observe board-a
```

**Project B:**
```bash
export PERCEPTA_CONFIG_PATH=~/projects/projectB/config.yaml
percepta observe board-b
```

### Shared Knowledge Base (Team)

Use shared network path for knowledge database:

```yaml
knowledge:
  path: /nfs/shared/percepta/team-knowledge.db
```

All team members contribute to and benefit from shared validated patterns.

### CI/CD Environment Variables

Override config for CI:

```bash
#!/bin/bash
# CI test script

export PERCEPTA_CONFIG_PATH=./ci-config.yaml
export PERCEPTA_STORAGE_PATH=/tmp/ci-observations.db
export ANTHROPIC_API_KEY=$CI_ANTHROPIC_KEY

percepta observe ci-test-board
percepta assert ci-test-board "led power is ON"
```

---

## See Also

- [Installation Guide](installation.md) - Setup and API keys
- [Getting Started](getting-started.md) - First device setup
- [Troubleshooting](troubleshooting.md) - Common issues
- [Commands Reference](commands.md) - All CLI commands
