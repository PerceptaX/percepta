# Percepta

**Computer vision for hardware observability.** Watch your firmware run.

Percepta uses Claude Vision to observe, validate, and compare physical hardware behavior (LEDs, displays, boot timing) without modifying firmware or hardware. Close the feedback loop in AI-driven firmware workflows — finally let Claude Code, Embedder, and other tools see what the hardware actually does.

## What is Percepta?

Percepta is the perception kernel for embedded development. Point a webcam at your ESP32/STM32/FPGA, and Percepta tells you exactly what it's doing:

```bash
$ percepta observe fpga

Signals (3):
LED1: blue blinking ~2.0Hz
LED2: purple blinking ~0.8Hz
LED3: red solid

Observation complete (748ms)
```

Then validate it with assertions:

```bash
$ percepta assert fpga "led('LED1').blinks() && led('LED1').color_rgb(0,0,255)"
✓ All assertions passed
```

And track regressions across firmware versions:

```bash
$ percepta diff fpga --from v1 --to v2
~ LED1: blue blinking 2.0Hz → 2.5Hz (MODIFIED)
```

## Quick Start

**1. Install:**
```bash
# Download binary from GitHub releases, or build from source:
go install github.com/Perceptax/percepta/cmd/percepta@latest
```

**2. Set API key:**
```bash
export ANTHROPIC_API_KEY="your-api-key-here"
```

**3. Add a device:**
```bash
$ percepta device add fpga
Device type (e.g., fpga, esp32, stm32): fpga
Camera device path (default: /dev/video0): /dev/video0
Firmware version (optional, press Enter to skip): v1

✓ Device 'fpga' added successfully
```

**4. Observe:**
```bash
$ percepta observe fpga
```

**That's it.** Point your webcam at the hardware and run `observe`.

## Commands

- **`percepta observe <device>`** — Capture current hardware state (LEDs, displays, boot timing)
- **`percepta assert <device> <dsl>`** — Validate expected behavior using assertion DSL
- **`percepta diff <device> --from <fw1> --to <fw2>`** — Compare behavior across firmware versions
- **`percepta device list/add/set-firmware`** — Manage device configurations

## Documentation

- **[Installation](docs/installation.md)** — Binary installation, building from source, environment setup
- **[Getting Started](docs/getting-started.md)** — Step-by-step walkthrough with examples
- **[Example Configs](docs/examples/)** — ESP32, STM32, FPGA device configurations

## Requirements

- **Claude API key** (ANTHROPIC_API_KEY environment variable)
- **Webcam** (USB camera or built-in, `/dev/video0` on Linux)
- **Go 1.20+** (for building from source)

Supported platforms: Linux, macOS, Windows

## Status

**Alpha** — Expect rough edges. Contributions welcome!

Current capabilities:
- ✅ LED detection (state, color, blink frequency)
- ✅ Display OCR (text extraction from OLED/LCD)
- ✅ Boot timing measurement
- ✅ DSL assertions (LED/display/timing validation)
- ✅ Firmware diff (version comparison)
- ✅ SQLite observation storage

See [ROADMAP.md](.planning/ROADMAP.md) for future plans.

## Why Percepta?

**The problem:** Embedded AI tools (Claude Code, Embedder) can generate firmware and control hardware, but cannot observe physical behavior. Developers manually validate by staring at LEDs and displays. Percepta closes this observability gap.

**The solution:** Computer vision + LLM perception = automated hardware observation. Point a camera at your hardware, get structured data about what it's actually doing.

**Target users:** Individual embedded developers (ESP32/STM32/FPGA) iterating on firmware at their desk. Comfortable with CLI tools, want fast feedback loops.

## License

MIT License - see repository for details.

## Contributing

Found a bug? Have a feature request? Open an issue or PR!

For alpha testing inquiries, reach out via GitHub issues.
