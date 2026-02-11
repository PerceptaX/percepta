# Percepta Release Notes

## v0.1.0-alpha (2026-02-12)

**First alpha release** — Computer vision for hardware observability.

### What's New

**Core Features:**
- ✅ **Vision-based observation** — Point a webcam at your hardware, get structured signals (LEDs, displays, boot timing)
- ✅ **Assertion DSL** — Validate expected behavior with deterministic assertions
- ✅ **Firmware diff** — Compare hardware behavior across firmware versions
- ✅ **Device management** — Interactive CLI for device configuration
- ✅ **SQLite storage** — Persistent observation history

**Commands:**
- `percepta observe <device>` — Capture current hardware state
- `percepta assert <device> <dsl>` — Validate expected behavior
- `percepta diff <device> --from <fw> --to <fw>` — Compare firmware versions
- `percepta device list/add/set-firmware` — Manage device configurations

**Signal Types:**
- LED detection (state, color, blink frequency)
- Display OCR (OLED/LCD text extraction)
- Boot timing measurement

**Supported Platforms:**
- Linux (x86_64, ARM64)
- macOS (Intel, Apple Silicon)
- Windows (x86_64)

### Installation

**Binary installation:**

1. Download the binary for your platform from [Releases](https://github.com/Perceptax/percepta/releases)
2. Extract and move to PATH:
   ```bash
   tar -xzf percepta-*.tar.gz
   sudo mv percepta-* /usr/local/bin/percepta
   ```
3. Set API key:
   ```bash
   export ANTHROPIC_API_KEY="your-api-key-here"
   ```
4. Verify:
   ```bash
   percepta --help
   ```

**Build from source:**
```bash
go install github.com/Perceptax/percepta/cmd/percepta@latest
```

**Full installation guide:** [docs/installation.md](docs/installation.md)

### Getting Started

**Quick start (3 steps):**

```bash
# 1. Add device
percepta device add fpga

# 2. Observe hardware
percepta observe fpga

# 3. Assert expected behavior
percepta assert fpga "led('LED1').blinks()"
```

**Full walkthrough:** [docs/getting-started.md](docs/getting-started.md)

### Known Issues

- **ISS-001: Single-frame capture limitation** — Blinking LEDs may be missed if OFF at capture instant. Run observation multiple times or wait for multi-frame capture support (post-alpha).
- **First observation is slow (~2-3s)** — Claude Vision API has cold-start latency. Subsequent observations are faster (~500-1000ms).
- **Color detection accuracy** — Best results with direct lighting, avoid glare. RGB tolerance is ±5 per channel.

### Requirements

- **Claude API key** (ANTHROPIC_API_KEY environment variable)
- **Webcam** (USB or built-in, `/dev/video0` on Linux)
- **Go 1.20+** (for building from source)

### Breaking Changes

None (first release).

### Contributors

Built with [Claude Code](https://claude.com/claude-code) and the Get Shit Done (GSD) workflow.

### Feedback

- **Issues:** [GitHub Issues](https://github.com/Perceptax/percepta/issues)
- **Discussions:** [GitHub Discussions](https://github.com/Perceptax/percepta/discussions)
- **Alpha testing:** Looking for 10 embedded developers to validate the core workflow. Reach out via GitHub!

### What's Next

**Post-alpha roadmap:**
- Multi-frame/video capture (ISS-001)
- MCP server mode for AI agent integration
- TUI/watch mode for continuous monitoring
- Local vision models (moondream, etc.)
- Assertion history tracking

See [ROADMAP.md](.planning/ROADMAP.md) for full development plan.

---

**License:** MIT
