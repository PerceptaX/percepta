# Percepta

[![CI](https://github.com/Perceptax/percepta/actions/workflows/ci.yml/badge.svg)](https://github.com/Perceptax/percepta/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/perceptax/percepta)](https://goreportcard.com/report/github.com/perceptax/percepta)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**AI firmware development with hardware validation.** Generate code. Flash hardware. Validate behavior. All automated.

Percepta uses Claude Vision to observe, validate, and compare physical hardware behavior (LEDs, displays, boot timing) without modifying firmware or hardware. Generate BARR-C compliant firmware with AI, automatically validated on real hardware. Close the feedback loop in embedded development — the only AI tool that knows if your code actually works.

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

## Core Features

**Vision-Based Hardware Testing:** *(Linux/macOS only)*
- **`percepta observe <device>`** — Capture hardware state (LEDs, displays, boot timing)
- **`percepta assert <device> <expr>`** — Validate expected behavior
- **`percepta diff <device> --from <fw1> --to <fw2>`** — Compare firmware versions

**AI Code Generation:** *(All platforms)*
- **`percepta generate <spec> --board <type>`** — Generate BARR-C compliant firmware
- **`percepta style-check <file> --fix`** — Enforce embedded coding standards *(Linux only)*
- **`percepta knowledge store/search`** — Manage validated pattern library

**Device Management:** *(All platforms)*
- **`percepta device add/list/set-firmware`** — Configure hardware devices

## Documentation

**Getting Started:**
- **[Installation Guide](docs/installation.md)** — Binary installation, building from source, API keys
- **[Getting Started](docs/getting-started.md)** — First observation in 10 minutes
- **[Commands Reference](docs/commands.md)** — Complete command documentation

**Advanced Usage:**
- **[Examples](docs/examples.md)** — 25+ workflow examples (LED validation, CI/CD, code generation)
- **[Configuration Guide](docs/configuration.md)** — Config file, camera setup, multi-device
- **[Troubleshooting](docs/troubleshooting.md)** — Common issues and solutions
- **[API Integration](docs/api-integration.md)** — Go library, CI/CD, MCP server (planned)

## Requirements

- **Claude API key** (ANTHROPIC_API_KEY environment variable)
- **Webcam** (USB camera or built-in, Linux/macOS only)
- **Go 1.24+** (for building from source)

## Platform Support

| Feature | Linux | macOS | Windows |
|---------|-------|-------|---------|
| **Hardware Observation** (`observe`, `assert`) | ✅ V4L2 | ✅ AVFoundation | ❌ |
| **Firmware Diffing** (`diff`) | ✅ | ✅ | ✅ |
| **AI Code Generation** (`generate`) | ✅ | ✅ | ✅ |
| **Style Checking** (`style-check`) | ✅ | ❌ | ❌ |
| **Knowledge Management** | ✅ | ✅ | ✅ |
| **Device Management** | ✅ | ✅ | ✅ |

**Note:** Camera-based commands (`observe`, `assert`) require a webcam and are only available on Linux (V4L2) and macOS (AVFoundation). Style checking requires tree-sitter (Linux only). All other features work cross-platform.

## Status

**v2.0 - Code Generation** — Production-ready with hardware validation loop

**Perception Features:**
- ✅ LED detection (state, color, blink frequency)
- ✅ Display OCR (text extraction from OLED/LCD)
- ✅ Boot timing measurement
- ✅ Multi-frame capture (detects blinking LEDs)
- ✅ Temporal smoothing (noise filtering)
- ✅ Confidence calibration

**Code Generation Features:**
- ✅ AI firmware generation (Claude Sonnet 4.5)
- ✅ BARR-C style enforcement (professional embedded standards)
- ✅ Auto-fix violations (naming, types)
- ✅ Knowledge graph (validated patterns only)
- ✅ Semantic search (find similar working code)
- ✅ Hardware validation loop (generate → flash → observe → validate)

**Storage & Testing:**
- ✅ SQLite observation storage
- ✅ Firmware version tracking
- ✅ Behavioral diffing
- ✅ Assertion DSL
- ✅ CLI with progress indicators

**Supported Boards:**
- ESP32, STM32, Arduino, ATmega, Generic

See [ROADMAP.md](.planning/ROADMAP.md) for Phase 8 (public launch) plans.

## Why Percepta?

**The problem:** AI tools generate embedded code, but you still manually validate by staring at LEDs. Code "compiles" doesn't mean it "works." No feedback loop between AI generation and hardware behavior.

**The solution:** Percepta closes the loop. Generate firmware with AI, flash to hardware, automatically observe behavior, validate against requirements. Store only patterns that actually work. Build a knowledge base of hardware-validated code.

**What makes it different:**
- **Hardware validation:** Code isn't "done" until hardware behaves correctly
- **BARR-C compliance:** Professional embedded coding standards, auto-fixed
- **Knowledge graph:** Only stores patterns proven on real hardware
- **Better than Embedder:** 100% works after validation vs "95% compiles, hope it works"

**Target users:** Embedded developers (ESP32/STM32/FPGA) using AI code generation tools. Want fast iteration with confidence that code actually works.

## License

MIT License - see repository for details.

## Contributing

Found a bug? Have a feature request? Open an issue or PR!

For alpha testing inquiries, reach out via GitHub issues.
