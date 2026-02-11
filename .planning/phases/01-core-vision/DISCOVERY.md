# Phase 1 Discovery: Core + Vision

## Research Findings

### 1. Anthropic Go SDK

**Package:** `github.com/anthropics/anthropic-sdk-go`
**Version:** v1.22.1 (latest stable, 2026)
**Model:** `claude-sonnet-4-5-20250929` (Sonnet 4.5)

**Vision API pattern:**
```go
client := anthropic.NewClient(option.WithAPIKey(apiKey))
message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
    MaxTokens: 1024,
    Model: anthropic.ModelClaudeSonnet4_5_20250929,
    Messages: []anthropic.MessageParam{
        anthropic.NewUserMessage(
            anthropic.NewImageBlockFromBase64("image/jpeg", base64Image),
            anthropic.NewTextBlock("Describe this embedded hardware device..."),
        ),
    },
})
```

### 2. Camera Capture

**Decision:** Start with `github.com/blackjack/webcam` (Linux-first)

**Rationale:**
- Pure Go wrapper over V4L2
- Zero C dependencies for Linux
- Simple API for frame capture
- Percepta targets Linux devs iterating on hardware
- Can add GoCV for macOS/Windows post-MVP if needed

**Alternative for cross-platform:** `gocv.io/x/gocv` (requires OpenCV 4.x + cgo, heavier build)

### 3. SQLite

**Choice:** `modernc.org/sqlite` (pure Go)

**Rationale:**
- No cgo = easy cross-compilation
- Single binary distribution
- ~2x slower than mattn/go-sqlite3, but acceptable for observation storage
- Critical for Homebrew/binary release simplicity

### 4. CLI Framework

**Recommend:** `github.com/spf13/cobra` + `github.com/spf13/viper`

**Standard Go CLI stack:**
- Cobra: Command structure, flags, help
- Viper: Config file loading (YAML)

## Architecture Decisions

### Signal Parsing from Vision API

Claude Vision returns unstructured text. Need to parse into structured Signal types:

**Approach:** Regex + keyword matching
- LED patterns: "blue LED is ON", "STATUS LED blinking at 1Hz"
- Display patterns: "OLED displays 'Ready v2.1'"
- Boot timing: timestamp tracking via multiple observations

**Don't hand-roll:** NLP/LLM for parsing (overkill, adds latency). Simple regex sufficient for v1.

### SQLite Schema

```sql
CREATE TABLE observations (
    id TEXT PRIMARY KEY,
    device_id TEXT NOT NULL,
    firmware_hash TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    signals JSON NOT NULL,  -- Array of Signal objects
    metadata JSON
);

CREATE INDEX idx_device_firmware ON observations(device_id, firmware_hash);
CREATE INDEX idx_timestamp ON observations(timestamp);
```

**Key decision:** Store signals as JSON (not separate LED/Display tables). Flexible, simpler queries for MVP.

## Implementation Notes

### Camera Access (Linux V4L2)

```go
import "github.com/blackjack/webcam"

cam, err := webcam.Open("/dev/video0")
defer cam.Close()

// Set format
format := webcam.PixelFormat(webcam.V4L2_PIX_FMT_MJPEG)
cam.SetImageFormat(format, 1920, 1080)

// Capture frame
frame, err := cam.ReadFrame()
// Convert to JPEG if needed
```

### Config File (~/.config/percepta/config.yaml)

```yaml
vision:
  provider: claude
  api_key: ${ANTHROPIC_API_KEY}

storage:
  backend: sqlite
  path: ~/.percepta/observations.db

devices:
  esp32-dev:
    type: esp32-devkit-v1
    camera_id: /dev/video0
```

## Common Pitfalls (from research)

1. **Anthropic SDK:** Must use `anthropic.NewImageBlockFromBase64()` - cannot pass raw bytes
2. **Webcam:** Must call `cam.Close()` or device stays locked
3. **SQLite:** `modernc.org/sqlite` uses `database/sql` standard interface - same query patterns as mattn
4. **Base64 encoding:** Use `encoding/base64.StdEncoding` not `RawStdEncoding` for Vision API

## Next Steps

1. Initialize Go module with dependencies
2. Define core types (Signal, Observation, Session)
3. Implement SQLite driver
4. Implement Claude Vision driver
5. Wire up CLI command

---
*Research completed: 2026-02-11*
