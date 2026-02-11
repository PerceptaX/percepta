# Percepta

## What This Is

Percepta is the perception kernel for physical hardware. It uses computer vision to observe, validate, and compare real-world hardware behavior (LED states, display content, boot timing) without modifying firmware or hardware. Embedded developers use it to close the feedback loop in AI-driven firmware workflows — finally letting Claude Code, Embedder, and other tools see what the hardware actually does.

## Core Value

**observe() must work reliably.** If Percepta can accurately tell you "the LED is blinking at 1.98 Hz" with 95%+ confidence, everything else follows. Perception accuracy is the foundation; all other features are secondary.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] **Vision capture with Claude Vision API** — Webcam frame → Claude Vision → structured LED/Display/Boot signals with 95%+ accuracy
- [ ] **Basic CLI: observe** — `percepta observe <device>` captures current hardware state and displays structured output
- [ ] **Observation storage** — SQLite stores observations (signals + timestamp + firmware hash) for later querying
- [ ] **Basic CLI: assert** — `percepta assert <device> <dsl>` validates expected behavior using DSL (e.g., `led('STATUS').blink_hz(1.0, tolerance=0.1)`)
- [ ] **Firmware tracking** — Associate observations with git commit hash for version comparison
- [ ] **Basic CLI: diff** — `percepta diff <device> --from <hash> --to <hash>` compares behavior across firmware versions
- [ ] **Device management** — `percepta device add/list/set-firmware` for device registration and firmware tracking
- [ ] **Signal types: LED** — Detect on/off state, color (RGB), brightness, blink frequency with confidence scores
- [ ] **Signal types: Display** — OCR text extraction from OLED/LCD displays with confidence scores
- [ ] **Signal types: Boot timing** — Measure boot sequence duration
- [ ] **Manual camera configuration** — YAML config specifies camera device (`camera_id: /dev/video0`) per device
- [ ] **DSL assertion engine** — Deterministic DSL evaluation for LED/display/timing assertions
- [ ] **Cross-platform binary** — Runs on Linux, macOS, Windows
- [ ] **Installation flow** — Homebrew formula + binary releases for easy install
- [ ] **Alpha-quality documentation** — README, quickstart guide, example device configs

### Out of Scope

- **MCP server mode** — deferred to post-MVP (week 11-12). Percepta works standalone first.
- **TUI/watch mode** — deferred to post-MVP (week 13-14). CLI commands only for MVP.
- **Session recording (video)** — deferred to post-MVP (week 9-10). No video capture in v1.
- **Local vision models** — Claude Vision API only for MVP. moondream/alternatives deferred.
- **Raw image storage** — observations store only structured signals, not raw camera frames.
- **Assertion history tracking** — SQLite stores observations only, not assertion pass/fail history.
- **Camera auto-detection/wizard** — manual YAML config only. No setup wizard.
- **Natural language assertions** — DSL only. No "make sure the LED is blinking" → DSL compilation.
- **Button state signals** — only LED, Display, Boot timing in v1.
- **Replace HIL/probe-rs/GDB** — Percepta observes external behavior, not internal state or production testing.

## Context

**Market gap:** Embedded AI tooling (Embedder, Claude Code, tinymcp) can generate firmware and control hardware, but cannot observe physical behavior. Developers manually validate by staring at LEDs and displays. Percepta closes this observability gap with vision.

**Target user:** Individual embedded developers (ESP32/STM32) iterating on firmware at their desk. Comfortable with CLI tools, wants fast feedback loops.

**Technical environment:**
- Go language for cross-platform binary with zero dependencies
- Claude Vision API (Sonnet 4.5) for multimodal LLM vision
- SQLite for local-first observation storage
- MIT license, open source from day 1

**Prior work:** PRD researched competitive landscape — no existing tools provide general-purpose vision-based observation for embedded development. probe-rs handles internal debugging, HIL systems handle production testing ($50K+), but nothing fills the "development iteration" observability gap.

**Success criteria:** If 10 alpha users successfully validate hardware behavior with `percepta observe/assert/diff` and keep using it daily, MVP is successful.

## Constraints

- **Tech stack: Go** — Cross-platform binary, zero runtime dependencies, fast compile times
- **Vision: Claude Vision API** — No local models in v1. Pluggable driver architecture allows swapping later.
- **Storage: SQLite** — Local-first, no cloud dependencies, portable database file
- **License: MIT** — Open source from day 1, community-first distribution
- **Platforms: Linux/macOS/Windows** — Must work on all three, primary development on Linux
- **Timeline: 8 weeks to alpha** — Optimize for earliest real user. Cut scope aggressively to reach first alpha fast.
- **Cost constraint: Claude Vision API** — Estimated ~$0.27/observation. Acceptable for MVP, defer local models if cost becomes issue.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| CLI-first, MCP deferred | Percepta must work standalone before AI agent integration. Validates core value proposition independently. | — Pending |
| Claude Vision only (no local models) | Ship faster, proven 95%+ accuracy. Pluggable architecture allows adding moondream post-MVP if needed. | — Pending |
| No raw image storage in v1 | Reduces SQLite size, faster queries. Can add as opt-in feature if debugging requires it. | — Pending |
| DSL assertions only (no NL→DSL) | Deterministic evaluation more important than syntax sugar. NL compilation can be added as UX layer later. | — Pending |
| Manual camera config (no wizard) | Cuts 1-2 weeks of setup UX work. Alpha users comfortable editing YAML. | — Pending |
| Go language | Zero-dependency binary, cross-platform, strong ecosystem for CLI tools (Cobra, Viper). | — Pending |

---
*Last updated: 2026-02-11 after initialization*
