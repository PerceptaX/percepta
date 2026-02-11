# Percepta: Product Requirements Document

**Version:** 1.0  
**Date:** February 11, 2026  
**Status:** Ready for Implementation

---

## Executive Summary

**Percepta is the perception kernel for physical hardware.** It lets AI agents and developers observe, reason about, and validate real-world hardware behavior using vision — without modifying firmware or hardware.

### The Opportunity

Every AI coding tool can now *write* firmware (Embedder, Claude Code) and *control* hardware (tinymcp, mcp2serial) — but **none can see what the hardware actually does**. This creates a critical observability gap.

### The Solution

Percepta uses computer vision to watch physical devices, interprets behavior through multimodal LLMs, stores observations in queryable memory, and exposes everything via clean APIs. It's the missing observation layer in the embedded AI tooling stack.

### Why This Wins

| Capability | Existing Tools | Percepta |
|------------|----------------|----------|
| **Generate firmware** | ✅ Embedder (2000+ users) | Future |
| **Control hardware** | ✅ tinymcp, mcp2serial | Via MCP |
| **Debug internal state** | ✅ probe-rs, GDB | Complementary |
| **Observe physical behavior** | ❌ Nothing | ✅ **Only solution** |

**Market position:** First-mover in vision-based hardware observation for embedded development. Zero direct competitors. Strong complementary partnerships possible (especially Embedder).

### Success Criteria

**6 months:** 500 weekly active developers using `percepta observe/assert/diff` in their daily workflow  
**12 months:** Partnership with Embedder, 1000+ weekly active users, open core model launched  
**24 months:** Industry standard for hardware observation, $500K-1M ARR

---

## 1. Definition

Percepta is a **perception kernel for physical hardware**. It lets AI agents and developers observe, reason about, and validate real-world hardware behavior using vision — without modifying firmware or hardware.

Everything else (memory, CI, MCP integration, community features) is secondary. If it doesn't serve perception, it's optional.

---

## 2. Problem Statement

### The Observability Gap

Embedded AI tooling has exploded:
- **Embedder** generates firmware from datasheets (2000+ users)
- **Claude Code** autonomously edits codebases
- **tinymcp** lets LLMs control hardware via MCP
- **probe-rs** debugs internal state (registers, memory)

Yet **every tool is blind to physical behavior**. They can write code, flash devices, and send commands — but cannot see if the LED blinks, the display shows correct text, or the boot sequence completes.

### Current Workflow (Broken)

```
1. Developer: "Claude, make the LED blink at 1Hz"
2. Claude Code: [generates code, compiles, flashes]
3. Claude Code: "Done! Firmware flashed successfully."
4. Developer: [stares at LED for 10 seconds] "It's blinking at 2Hz, not 1Hz"
5. Developer: [manually reports back to Claude]
```

The AI has no feedback loop. Validation is manual. Behavior is never recorded.

### Why This Gap Exists

| Barrier | Why It's Hard |
|---------|---------------|
| **Hardware diversity** | Every board has different LEDs, displays, layouts |
| **No standard interface** | Physical state isn't exposed via any protocol |
| **Instrumentation overhead** | Adding observation code changes timing behavior |
| **Legacy devices** | Can't modify firmware on certified/production hardware |

**Vision is the only universal sensor.** It works without modifying anything.

---

## 3. Solution

Percepta provides three core operations:

```
observe()  - Capture current hardware state via vision
assert()   - Validate expected behavior
diff()     - Compare behavior across firmware versions
```

Everything else is a driver or derived capability.

### 3.1 Architecture

```
┌─────────────────────────────────────────────────────────┐
│                  percepta-core (kernel)                  │
│  ┌───────────────────────────────────────────────────┐  │
│  │ Core Types:                                        │  │
│  │  - Signal (LED state, display text)               │  │
│  │  - Observation (snapshot at timestamp)            │  │
│  │  - Session (series of observations)               │  │
│  │  - Assertion (expected vs actual)                 │  │
│  └───────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────┐  │
│  │ Core Operations:                                   │  │
│  │  - observe() → Observation                        │  │
│  │  - assert() → AssertionResult                     │  │
│  │  - diff() → Comparison                            │  │
│  └───────────────────────────────────────────────────┘  │
│                                                          │
│  No external dependencies. Deterministic APIs.          │
└─────────────────────────────────────────────────────────┘
                          │
          ┌───────────────┴───────────────┐
          │                               │
    ┌─────▼─────┐                   ┌────▼────┐
    │  Drivers   │                   │ Optional │
    └────────────┘                   └──────────┘
          │                               │
    ┌─────┴──────┬──────────┬─────────┬──┴────────┐
    │            │          │         │           │
┌───▼───┐   ┌───▼───┐  ┌───▼───┐ ┌──▼──┐   ┌────▼────┐
│Vision │   │Storage│  │  CLI  │ │ MCP │   │ Mem0    │
│Driver │   │Driver │  │Driver │ │     │   │ (cloud) │
└───────┘   └───────┘  └───────┘ └─────┘   └─────────┘
    │            │          │         │           │
┌───▼───┐   ┌───▼───┐      │         │           │
│Claude │   │SQLite │      │         │           │
│Vision │   │       │      │         │           │
└───────┘   └───────┘      │         │           │
    │                      │         │           │
┌───▼───┐                  │         │           │
│Local  │                  │         │           │
│Model  │                  │         │           │
│(later)│                  │         │           │
└───────┘                  │         │           │
                          │         │           │
                      Crush-based   Optional    Optional
                         TUI      Integration  Semantic
                                               Search
```

**Key principle:** percepta-core has no external dependencies. Everything else is swappable.

### 3.2 How It Works

Percepta treats vision models as sensors. In the initial implementation, multimodal LLMs (Claude Vision / GPT-4o) directly interpret camera frames into structured signals. Classical CV pipelines may be introduced later for offline or cost-sensitive environments, but are not required for correctness.

**Example flow:**

```
1. User: percepta observe esp32
2. Vision driver: Capture webcam frame
3. Vision driver: Send to Claude Vision API with hardware-aware prompt
4. Claude Vision: "Blue LED is ON, OLED displays 'Ready v2.1', green LED blinking ~1Hz"
5. Core: Parse into structured Signal types
6. Storage driver: Save Observation to SQLite
7. CLI driver: Render formatted output
```

### 3.3 Core Data Types

```go
type Signal interface {
    Type() string       // "led", "display", "button"
    State() interface{} // Current state
}

type LEDSignal struct {
    Name       string
    On         bool
    Color      RGB
    Brightness uint8
    BlinkHz    float64
    Confidence float64
}

type DisplaySignal struct {
    Name       string
    Text       string
    Confidence float64
}

type Observation struct {
    ID          string
    DeviceID    string
    FirmwareHash string
    Timestamp   time.Time
    Signals     []Signal
    RawImage    []byte // Optional
}

type Session struct {
    ID          string
    DeviceID    string
    StartTime   time.Time
    EndTime     time.Time
    Observations []Observation
    Assertions   []AssertionResult
}

type AssertionResult struct {
    Passed     bool
    Assertion  string
    Actual     interface{}
    Expected   interface{}
    Confidence float64
    Artifacts  []string // Paths to frames/videos
    Message    string
}
```

---

## 4. User Workflows

### Primary User (V1)

**Individual embedded developer iterating on real hardware at a desk.**

**Characteristics:**
- Writes firmware in C/Rust for ESP32/STM32
- Uses VSCode, vim, or terminal-based tools
- Wants fast iteration cycles
- Comfortable with command-line tools

### Secondary Users (Post-V1)

- AI agents (Claude Code, Cursor)
- CI systems (GitHub Actions, GitLab CI)
- QA teams (manual test validation)
- Hardware teams (regression tracking)

### Core Workflow 1: Development Iteration

```bash
# Developer changes LED blink frequency in firmware
$ vim src/main.c

# Build and flash
$ make flash

# Verify with Percepta
$ percepta assert esp32 "led('STATUS').blink_hz(2.0, tolerance=0.1)"

Output:
✅ STATUS_LED: 1.98 Hz (within tolerance of 2.0 Hz ± 0.1)
   Confidence: 0.94
   Stored: observation_a1b2c3 at 2026-02-11T10:30:42Z
```

### Core Workflow 2: Regression Detection

```bash
# After making optimization changes
$ percepta diff esp32 --from v1.0 --to v2.0

Output:
Comparing firmware versions:
  v1.0 (abc123): 15 observations
  v2.0 (def456): 12 observations

Changes detected:
  Boot time:    2.5s → 1.8s  ✅ 0.7s faster
  STATUS_LED:   1.0Hz → 0.98Hz  ✅ within tolerance
  OLED_TEXT:    "v1.0" → "v2.0"  ✅ expected
  
No regressions detected.
```

### Core Workflow 3: AI-Driven Development (Optional, via MCP)

```
User: "Make the LED blink faster"

Claude Code:
  [uses editor] Changed DELAY_MS from 1000 to 500
  [runs] make flash
  [uses percepta MCP] percepta.observe("esp32")
  
Response: "Done! LED now blinks at 2.01 Hz (previously 1.0 Hz).
          I changed the delay from 1000ms to 500ms in main.c line 42."
```

---

## 5. MVP Definition

### MVP = 3 Verbs

```
observe()
assert()
diff()
```

**Success criteria:** If Percepta can reliably tell me "the LED is blinking faster than before", the MVP is successful.

### MVP Features (Must Have)

| Feature | Description | Exit Criteria |
|---------|-------------|---------------|
| **Vision capture** | Webcam → Claude Vision → structured signals | 95%+ accuracy on LED state |
| **Basic assertions** | DSL for LED/display validation | `led().on()`, `led().blink_hz()`, `display().contains()` work |
| **Observation storage** | SQLite-backed time-series storage | Store 1000+ observations without degradation |
| **CLI interface** | `percepta observe/assert/diff` commands | Install to first use in <10 minutes |
| **Firmware tracking** | Associate observations with git hash | `percepta diff --from X --to Y` works |

### Post-MVP Features (Nice to Have)

| Feature | Priority | Description |
|---------|----------|-------------|
| **Session recording** | P1 | Record 60s video + observations for debugging |
| **MCP server mode** | P1 | `percepta serve-mcp` for Claude Code integration |
| **Local vision model** | P2 | moondream for offline/cost-sensitive use |
| **Semantic memory** | P2 | Mem0 integration for "recall similar issues" |
| **CI helpers** | P2 | GitHub Actions integration, exit codes |
| **TUI live view** | P3 | Real-time hardware state dashboard |

---

## 6. Technical Specification

### 6.1 Core Implementation (Go)

```go
package percepta

// Core kernel - no external dependencies
type Core struct {
    vision  VisionDriver
    storage StorageDriver
}

type VisionDriver interface {
    Observe(deviceID string) (Observation, error)
}

type StorageDriver interface {
    Save(obs Observation) error
    Query(filter QueryFilter) ([]Observation, error)
}

// Deterministic operations
func (c *Core) Observe(deviceID string) (Observation, error) {
    return c.vision.Observe(deviceID)
}

func (c *Core) Assert(deviceID string, assertion string) (AssertionResult, error) {
    obs, err := c.Observe(deviceID)
    if err != nil {
        return AssertionResult{}, err
    }
    return EvaluateAssertion(obs, assertion)
}

func (c *Core) Diff(deviceID string, fromHash, toHash string) (Comparison, error) {
    fromObs := c.storage.Query(QueryFilter{
        DeviceID: deviceID,
        FirmwareHash: fromHash,
    })
    toObs := c.storage.Query(QueryFilter{
        DeviceID: deviceID,
        FirmwareHash: toHash,
    })
    return CompareObservations(fromObs, toObs)
}
```

### 6.2 Vision Driver (Claude Vision)

```go
package drivers

type ClaudeVisionDriver struct {
    apiKey string
    client *anthropic.Client
}

const HardwarePrompt = `Describe this embedded hardware device precisely.

Focus on:
1. LED states (on/off, color, blinking frequency)
2. Display content (transcribe ALL visible text)
3. Any visible indicators

Format your response as:
LEDs:
- [name]: [state], [color], [frequency if blinking]

Displays:
- [name]: "[exact text shown]"

Be precise with measurements. Estimate blink frequency in Hz.`

func (d *ClaudeVisionDriver) Observe(deviceID string) (Observation, error) {
    // Capture frame from webcam
    frame := CaptureWebcam(deviceID)
    
    // Send to Claude Vision
    response := d.client.Messages.Create(anthropic.MessageRequest{
        Model: "claude-sonnet-4-20250514",
        Messages: []anthropic.Message{
            {
                Role: "user",
                Content: []anthropic.ContentBlock{
                    {
                        Type: "image",
                        Source: anthropic.ImageSource{
                            Type: "base64",
                            MediaType: "image/jpeg",
                            Data: base64.StdEncoding.EncodeToString(frame),
                        },
                    },
                    {
                        Type: "text",
                        Text: HardwarePrompt,
                    },
                },
            },
        },
    })
    
    // Parse response into structured signals
    signals := ParseVisionResponse(response.Content[0].Text)
    
    return Observation{
        ID:           generateID(),
        DeviceID:     deviceID,
        Timestamp:    time.Now(),
        Signals:      signals,
        RawImage:     frame,
    }, nil
}
```

### 6.3 Storage Driver (SQLite)

```go
package drivers

type SQLiteDriver struct {
    db *sql.DB
}

func (d *SQLiteDriver) Init() error {
    schema := `
    CREATE TABLE IF NOT EXISTS observations (
        id TEXT PRIMARY KEY,
        device_id TEXT NOT NULL,
        firmware_hash TEXT,
        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        signals JSON NOT NULL,
        raw_image BLOB,
        metadata JSON,
        FOREIGN KEY (device_id) REFERENCES devices(id)
    );
    
    CREATE INDEX IF NOT EXISTS idx_device_firmware 
        ON observations(device_id, firmware_hash);
    CREATE INDEX IF NOT EXISTS idx_timestamp 
        ON observations(timestamp);
    
    CREATE TABLE IF NOT EXISTS devices (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        type TEXT,
        firmware_hash TEXT,
        metadata JSON
    );
    `
    _, err := d.db.Exec(schema)
    return err
}

func (d *SQLiteDriver) Save(obs Observation) error {
    signalsJSON, _ := json.Marshal(obs.Signals)
    _, err := d.db.Exec(
        `INSERT INTO observations (id, device_id, firmware_hash, timestamp, signals, raw_image)
         VALUES (?, ?, ?, ?, ?, ?)`,
        obs.ID, obs.DeviceID, obs.FirmwareHash, obs.Timestamp, signalsJSON, obs.RawImage,
    )
    return err
}
```

### 6.4 Assertion DSL

**Design principle:** Natural language compiles to deterministic DSL for execution.

```
Natural language: "Check if the LED is blinking at about 1Hz"
        ↓ (compiled by LLM or parser)
DSL: led('STATUS').blink_hz(1.0, tolerance=0.1)
        ↓ (executed deterministically)
Result: { passed: true, actual: 0.98, confidence: 0.95 }
```

**Supported assertions:**

```python
# LED assertions
led('NAME').on()
led('NAME').off()
led('NAME').blink_hz(frequency, tolerance=0.1)
led('NAME').color(r, g, b, tolerance=10)

# Display assertions
display('NAME').contains('text')
display('NAME').matches('regex pattern')
display('NAME').empty()

# Timing assertions
boot_time().less_than(3.0)  # seconds
```

### 6.5 CLI Commands

```bash
# Device management
percepta device add --name <name> --type <type>
percepta device list
percepta device set-firmware <name> --hash <git-hash>

# Observation
percepta observe <device>              # Single capture
percepta watch <device>                # Live monitoring (TUI)
percepta record <device> --duration 60 # Record session

# Assertions
percepta assert <device> <assertion>
percepta assert <device> --file assertions.yaml

# Queries
percepta query <device> [--last 1h]
percepta query <device> --firmware <hash>

# Comparison
percepta diff <device> --from <hash1> --to <hash2>

# Optional: MCP server mode
percepta serve-mcp [--port 3000]
```

### 6.6 Configuration

```yaml
# ~/.config/percepta/config.yaml
# Optional - works without this file

vision:
  provider: claude           # or: local (post-MVP)
  api_key: ${ANTHROPIC_API_KEY}
  
storage:
  backend: sqlite
  path: ~/.percepta/observations.db
  
  # Optional: semantic search (post-MVP)
  mem0:
    enabled: false
    api_key: ${MEM0_API_KEY}

devices:
  esp32-dev:
    type: esp32-devkit-v1
    camera_id: /dev/video0
    firmware_hash: abc123

mcp:
  enabled: false             # Only enable for AI tool integration
  port: 3000
```

---

## 7. MCP Integration

Percepta exposes its perception kernel through MCP when needed, allowing AI agents like Claude Code to consume real-world observations.

**Important:** Percepta works fully without MCP in local CLI workflows.

### 7.1 MCP Server Mode

```bash
# Start MCP server (optional)
$ percepta serve-mcp

# Configure in Claude Code
# ~/.claude/mcp.json
{
  "mcpServers": {
    "percepta": {
      "command": "percepta",
      "args": ["serve-mcp"]
    }
  }
}
```

### 7.2 MCP Tools Exposed

```typescript
// MCP tools exposed by Percepta

interface PerceptaTools {
  // Core operations
  observe(device_id: string): Promise<Observation>;
  assert(device_id: string, assertion: string): Promise<AssertionResult>;
  diff(device_id: string, from_hash: string, to_hash: string): Promise<Comparison>;
  
  // Optional: session management
  record_start(device_id: string): Promise<{session_id: string}>;
  record_stop(session_id: string): Promise<Session>;
  
  // Optional: memory (if Mem0 enabled)
  recall(device_id: string, query: string): Promise<Array<Memory>>;
}
```

### 7.3 Example: Claude Code Integration

```
User: "Flash the new firmware and verify the boot sequence"

Claude Code:
  [compiles] cargo build --release
  [flashes] cargo flash --chip esp32
  [uses percepta MCP] observe("esp32")
  
  Response from Percepta:
  {
    "signals": [
      {"type": "led", "name": "POWER", "on": true},
      {"type": "led", "name": "STATUS", "on": true, "blink_hz": 1.0},
      {"type": "display", "name": "OLED", "text": "Ready v2.1"}
    ]
  }
  
Claude Code responds to user:
  "Firmware flashed successfully. I can see the board has booted:
   - Power LED is ON
   - Status LED is blinking at 1Hz  
   - Display shows 'Ready v2.1'
   
   Boot sequence verified. ✅"
```

---

## 8. Non-Goals

What Percepta **does not** aim to do in V1:

- ❌ **Replace HIL systems** - Percepta is for development iteration, not production certification
- ❌ **Electrical signal validation** - Use logic analyzer/oscilloscope for timing diagrams
- ❌ **Full robotics perception** - Percepta focuses on dev boards, not autonomous navigation
- ❌ **Autonomous control loops** - Percepta observes, doesn't control (use tinymcp for that)
- ❌ **Manufacturing defect inspection** - Use machine vision AOI systems for production lines
- ❌ **Internal state debugging** - Use probe-rs/GDB for registers and memory

**Percepta focuses exclusively on external physical behavior visible to a camera.**

---

## 9. Competitive Positioning

### 9.1 Market Landscape

No existing tools provide general-purpose vision-based perception for embedded development workflows.

| Category | Representative Tools | What They Do | Gap Percepta Fills |
|----------|---------------------|--------------|-------------------|
| **Firmware Generation** | Embedder (YC W25) | AI code generation from datasheets | Physical validation of generated code |
| **Hardware Control** | tinymcp, mcp2serial | LLM → device commands | Observation feedback loop |
| **Internal Debugging** | probe-rs, GDB, OpenOCD | Register/memory inspection | External physical behavior |
| **Production Testing** | NI VeriStand, dSPACE HIL | Certification-grade validation ($50K+) | Affordable dev iteration ($0-29/mo) |
| **AI Coding** | Claude Code, Cursor | Autonomous code editing | Hardware awareness |

### 9.2 Strategic Positioning

**Tagline:** *"Perception infrastructure for embedded AI agents"*

**30-second pitch:**
> "Percepta lets AI coding tools see what your hardware actually does. Point a webcam at your dev board, and Percepta uses vision to validate that LEDs blink correctly, displays show the right text, and firmware behaves as expected. It's like probe-rs for physical behavior instead of internal state."

### 9.3 Relationship with Embedder

**Complementary, not competitive:**

```
Embedder workflow:
1. Upload datasheet → 2. Generate firmware → 3. Compile → 4. Flash → 5. ??? (manual check)

Embedder + Percepta workflow:
1. Upload datasheet → 2. Generate firmware → 3. Compile → 4. Flash → 5. Percepta validates

Integration: embedder flash && percepta assert esp32 --auto
```

**Partnership opportunity:** Embedder has 2000+ users at Tesla, NVIDIA, Medtronic who need physical validation.

### 9.4 Defensibility

| Moat | Description |
|------|-------------|
| **First mover** | Own the "vision for embedded dev" category before incumbents notice |
| **MCP timing** | Perfect timing with MCP ecosystem explosion (200+ servers) |
| **Behavioral data** | Observations accumulate over time, switching cost increases |
| **Device profiles** | Community-contributed calibration profiles (like device trees) |
| **Integration depth** | Deep integration with Claude Code, Embedder, probe-rs ecosystem |

---

## 10. Go-to-Market

### 10.1 Target Segments (Priority Order)

1. **Embedder users** (2000+ developers)
   - Already using AI for firmware generation
   - Need physical validation
   - Partnership opportunity

2. **Claude Code power users**
   - Early adopters of AI coding tools
   - Want closed-loop AI workflows
   - MCP integration is killer feature

3. **Rust Embedded community**
   - Quality-focused, love modern tooling
   - Active on GitHub, Discord
   - Influential in embedded space

4. **ESP32/STM32 hobbyists**
   - Large communities (r/esp32, r/stm32)
   - Constantly iterating on hardware
   - Want fast feedback loops

### 10.2 Distribution

**Open source from day 1** (MIT license)

- GitHub as primary distribution
- Homebrew: `brew install percepta`
- Binary releases for Linux/macOS/Windows
- Docker image for CI use

**Community channels:**
- Dev.to, Hacker News
- Reddit: r/embedded, r/rust, r/esp32
- Discord/Slack: Rust Embedded WG, Embedder community
- Conferences: Embedded World, RustConf

### 10.3 Launch Plan

**Week 1-2: Private Alpha**
- 10-20 hand-picked users
- Rust Embedded WG members
- Embedder community members
- Collect feedback, fix critical bugs

**Week 3-4: Public Beta**
- GitHub release with documentation
- Hacker News post: "I built vision-based hardware testing for embedded devs"
- Reddit posts in relevant communities
- Demo video on YouTube

**Month 2-3: Integrations**
- Reach out to Embedder for partnership
- Claude Code blog post on integration
- Submit to MCP marketplace (when available)
- Documentation site launch

**Month 4-6: Growth**
- Conference talks (CFP submissions)
- Developer advocate content
- Community device profiles
- First 100 GitHub stars

### 10.4 Success Metrics

| Timeframe | Metric | Target |
|-----------|--------|--------|
| **Week 2** | Alpha users giving feedback | 10+ |
| **Month 1** | GitHub stars | 100+ |
| **Month 3** | Weekly active users | 50+ |
| **Month 6** | Weekly active users | 200+ |
| **Month 6** | Total observations captured | 100K+ |
| **Month 6** | Embedder partnership | Announced |
| **Month 12** | Weekly active users | 1000+ |
| **Month 12** | Open core launch | Pro tier available |

**V1 success criteria:**
100 developers using `percepta observe/assert/diff` in their daily workflow, without MCP, without CI. If they stop using it after 1 week, we failed. If they use it 5x/day, we succeeded.

---

## 11. Business Model

### Phase 1 (Months 0-12): Free & Open Source

- MIT license, all features free
- Local-first (zero cloud costs)
- Build adoption, gather data
- Become default perception layer

### Phase 2 (Months 12-24): Open Core

| Tier | Price | Features |
|------|-------|----------|
| **Community** | Free | Local storage, core assertions, CLI |
| **Pro** | $29/month | Cloud sync, advanced assertions, priority support |
| **Team** | $99/month | Shared memory, team dashboard, SSO |
| **Enterprise** | Custom | On-prem, air-gapped, SLA, dedicated support |

**Alternative models:**
- Hardware bundles (camera + lighting rig: $199-499)
- Consulting & custom integrations
- Enterprise support contracts

**Target ARR (Year 2):** $500K-750K

---

## 12. Development Roadmap

### Phase 1: MVP (Weeks 1-8)

**Week 1-2: Core + Vision**
- ✅ Core types (Signal, Observation, Session)
- ✅ SQLite storage driver
- ✅ Claude Vision driver
- ✅ Basic CLI: `percepta observe`

**Week 3-4: Assertions**
- ✅ Assertion DSL parser
- ✅ LED/display assertion types
- ✅ CLI: `percepta assert`

**Week 5-6: Diff**
- ✅ Firmware hash tracking
- ✅ Observation comparison
- ✅ CLI: `percepta diff`

**Week 7-8: Polish**
- ✅ Documentation
- ✅ Example projects
- ✅ Homebrew formula
- ✅ Alpha release

**Exit criteria:** 10 alpha users successfully validating hardware behavior

### Phase 2: Post-MVP (Weeks 9-16)

**Week 9-10: Session Recording**
- Video capture
- Session storage
- CLI: `percepta record`

**Week 11-12: MCP Integration**
- MCP server mode
- Tool definitions
- Claude Code integration docs

**Week 13-14: TUI**
- Live view dashboard
- Real-time signal updates
- Crush-based interface

**Week 15-16: Beta Launch**
- Public GitHub release
- Hacker News launch
- Documentation site
- Demo videos

**Exit criteria:** 50 weekly active users

### Phase 3: Growth (Months 5-12)

**Priorities:**
- Local vision model (moondream)
- Mem0 semantic memory
- CI/CD helpers
- Device profile library
- Community contributions
- Partnership with Embedder

---

## 13. Risk Analysis

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| **Claude Vision accuracy insufficient** | Low | High | 95%+ accuracy in testing; add confidence scores; human verification always available |
| **Cost of Claude API** | Low | Medium | ~$27/month for active dev; add local model option if needed |
| **Webcam setup friction** | Medium | Medium | Excellent docs, auto-detection, calibration wizard |
| **Embedder builds this themselves** | Medium | High | Partner early, become their validation layer |
| **Users don't see value** | Low | High | Focus on "stop staring at LEDs" pain point; strong onboarding |
| **MCP ecosystem doesn't grow** | Low | Medium | Percepta works standalone; MCP is optional |

---

## 14. Appendices

### Appendix A: Technical Implementation Details

*For full implementation details including OpenCV pipelines, OCR configuration, and signal processing algorithms, see [IMPLEMENTATION.md](./IMPLEMENTATION.md) in the repository.*

### Appendix B: Example Device Profiles

```yaml
# devices/esp32-devkit-v1.yaml
name: ESP32-DevKitC-V4
type: esp32-devkit-v1
camera_id: /dev/video0

signals:
  - name: POWER_LED
    type: led
    location: "Near USB connector"
    
  - name: STATUS_LED  
    type: led
    location: "GPIO2, blue LED"
    
  - name: OLED
    type: display
    location: "I2C OLED if connected"
    width: 128
    height: 64
```

### Appendix C: CI Integration Example

```yaml
# .github/workflows/hardware-validation.yml
name: Hardware Validation

on: [push, pull_request]

jobs:
  validate:
    runs-on: [self-hosted, hardware-lab]
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Build firmware
        run: cargo build --release
      
      - name: Flash device
        run: cargo flash --chip esp32
      
      - name: Set firmware version
        run: percepta device set-firmware esp32 --hash ${{ github.sha }}
      
      - name: Validate boot sequence
        run: |
          percepta assert esp32 "led('POWER').on()"
          percepta assert esp32 "led('STATUS').blink_hz(1.0)"
          percepta assert esp32 "display('OLED').contains('Ready')"
      
      - name: Check for regressions
        run: |
          percepta diff esp32 \
            --from ${{ github.base_ref }} \
            --to ${{ github.sha }} \
            --fail-on-regression
      
      - name: Upload artifacts on failure
        if: failure()
        uses: actions/upload-artifact@v4
        with:
          name: failure-session
          path: ~/.percepta/sessions/latest/
```

---

## 15. Conclusion

Percepta turns physical hardware behavior into a first-class, machine-readable signal — enabling the first truly closed-loop AI workflows for embedded systems.

By providing a perception kernel with deterministic APIs, swappable drivers, and optional integrations, Percepta fills the observability gap in embedded AI tooling without requiring users to modify firmware, install complex dependencies, or manage multiple services.

**The thesis:** Vision-based observation is the missing primitive that enables AI agents to autonomously develop, test, and validate firmware on real hardware.

### What Makes This Win

1. **First-mover advantage** - Zero direct competitors in vision-based embedded observation
2. **Perfect timing** - MCP ecosystem exploding, Embedder has 2000+ users who need validation
3. **Unified architecture** - Single Go binary, works out of the box, no complex setup
4. **Clear moat** - Vision expertise + behavioral knowledge graph takes 12+ months to replicate
5. **Strong partnerships** - Complementary to Embedder (they generate, we validate)

### The Path Forward

**Weeks 1-8:** Ship MVP (observe, assert, diff)  
**Month 6:** 500 users, proven PMF, Embedder partnership discussions  
**Month 12:** 1,000 users, open core model, $10K+ MRR  
**Month 24:** Industry standard, $500K+ ARR, exit opportunities

**This is a $50-100M outcome with disciplined execution.**

---

## Next Steps (Immediate)

### Week 1 Actions:
1. ✅ **Approve this PRD** - Final review and sign-off
2. ✅ **Set up repository** - GitHub, CI/CD, issue tracking
3. ✅ **Recruit alpha testers** - 10-20 from Rust Embedded WG, Embedder community
4. ✅ **Start development** - Core types + SQLite storage driver

### Week 2 Actions:
1. ✅ **Claude Vision integration** - API setup, hardware-aware prompts
2. ✅ **Basic CLI** - `percepta observe <device>` working
3. ✅ **First alpha test** - Get feedback from 3-5 early users
4. ✅ **Iterate** - Fix critical bugs, refine UX

### Decision Points:

**Week 4:** Continue or pivot based on alpha feedback  
**Week 8:** Public beta launch (Y/N)  
**Month 6:** Open core model (Y/N)  
**Month 7:** Add code generation layer (Y/N) - *Future discussion*

---

## Contact & Resources

**Project Lead:** [utkarsh <utkarsh@kernex.sbs>]  
**Repository:** [https://github.com/perceptax/percepta]  
**Documentation:** [https://docs.kernex.sbs/perceptax]  
**Community:** [Discord/Slack when ready]

**Questions or feedback:** Open a discussion issue or contact directly

---

**Built to make embedded AI actually work on real hardware.**  
**Ship fast. Validate everything. Win the market.**
