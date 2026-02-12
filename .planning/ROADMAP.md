# Roadmap: Percepta

## Overview

Percepta evolves from perception kernel (v1.0) to complete firmware platform (v2.0+). Phase 1 built vision-based hardware observation; Phase 2 adds AI code generation with hardware validation; Phase 3 scales to cloud platform with community features.

## Domain Expertise

None

## Milestones

- ✅ **v1.0 Perception MVP** - Phases 1-4 (shipped 2026-02-12)
- 🚧 **v2.0 Code Generation** - Phases 5-8 (in progress)

## Phases

<details>
<summary>✅ v1.0 Perception MVP (Phases 1-4) - SHIPPED 2026-02-12</summary>

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

### Phase 1: Core + Vision
**Goal**: percepta observe <device> works end-to-end with 95%+ accuracy on LED/display/boot signals

**Depends on**: Nothing (first phase)

**Research**: Likely (Claude Vision API integration, Go camera capture)

**Research topics**: Anthropic Go SDK usage, Go webcam libraries (gocv vs native options), SQLite schema design for time-series observations

**Plans**: 3 plans

Plans:
- [x] 01-01: Core types and in-memory storage (SQLite deferred)
- [x] 01-02: Claude Vision driver and camera capture
- [x] 01-03: CLI observe command and output formatting

### Phase 2: Assertions
**Goal**: percepta assert <device> <dsl> validates expected behavior deterministically

**Depends on**: Phase 1 (needs observe() working)

**Research**: Unlikely (internal DSL parser, deterministic evaluation logic)

**Plans**: 2 plans

Plans:
- [x] 02-01: DSL parser and assertion types (LED, display, timing)
- [x] 02-02: CLI assert command and result formatting

### Phase 2.5: Multi-LED Signal Identity (INSERTED)
**Goal**: Fix parser to extract ALL LEDs (not just first match) with deterministic identity (LED1, LED2, LED3)

**Depends on**: Phase 2 (needs assertions working to validate fix)

**Research**: None (refactor existing parser)

**Plans**: 1 plan

Plans:
- [x] 2.5-01: Refactor parser to extract all LEDs with index-based naming

**Why this was critical:**
Parser only extracted first LED match, causing unstable signal identity. Fixed to establish object permanence (LED1 = LED1 consistently), enabling reliable diff.

### Phase 3: Diff + Firmware Tracking
**Goal**: percepta diff --from X --to Y compares behavior across firmware versions

**Depends on**: Phase 2 (needs observations + assertions), Phase 2.5 (needs stable signal identity)

**Research**: None (straightforward SQLite + git integration)

**Plans**: 2 plans

Plans:
- [x] 03-01: SQLite storage with manual firmware tagging (modernc.org/sqlite, config.Device.Firmware, no git integration)
- [x] 03-02: Exact signal comparison and diff command (diff.Compare(), CLI with --from/--to flags, normalized BlinkHz)

### Phase 4: Polish + Alpha
**Goal**: Ship to 10 alpha users with installation in <10 minutes

**Depends on**: Phase 3 (needs all three verbs working)

**Research**: Unlikely (installation tooling, documentation)

**Plans**: 2 plans

Plans:
- [x] 04-01: Device management CLI (device list/add/set-firmware)
- [x] 04-02: Documentation, build script, alpha release

</details>

### 🚧 v2.0 Code Generation (In Progress)

**Milestone Goal:** Generate professional, hardware-validated firmware that passes code review. Only AI tool that delivers code that actually works on real hardware with BARR-C style compliance.

**Key innovation:** Automatic validation loop (generate → flash → observe → validate) beats Embedder's "95% compiles, hope it works" with "100% works after validation."

#### Phase 5: Style Infrastructure

**Goal**: BARR-C style checker + enforcement engine ensures generated code follows professional embedded coding standards

**Depends on**: v1.0 complete (needs perception kernel for validation)

**Research**: Likely (BARR-C standard specification, C AST parsing in Go, auto-fix strategies)

**Research topics**: BARR-C Embedded C Coding Standard rules, Go C parser libraries (tree-sitter vs custom), style violation auto-fix patterns, integration with code generation workflow

**Plans**: 2 plans

Plans:
- [x] 05-01: BARR-C rule engine with tree-sitter C parser
- [x] 05-02: Auto-fix engine and percepta style-check CLI

#### Phase 6: Knowledge Graphs

**Goal**: Behavioral and style knowledge graphs store validated patterns (what actually works on hardware) to guide future code generation

**Depends on**: Phase 5 (style checker needed to validate patterns)

**Research**: Likely (Graph database selection, vector search for semantic similarity, pattern extraction from validated code)

**Research topics**: Neo4j vs alternatives for behavioral graph, Qdrant vs alternatives for vector search, embedding models for code similarity, knowledge graph schema design for hardware quirks

**Plans**: 2 plans

Plans:
- [x] 06-01: Knowledge Graph Storage (graph DB setup, PatternStore API)
- [x] 06-02: Semantic Search + CLI (vector store, semantic search, CLI commands)

#### Phase 6.1: Perception Enhancements (INSERTED)

**Goal**: Enhance vision system reliability for hardware validation loop (LCD OCR, multi-object tracking, temporal smoothing, schema stability)

**Depends on**: Phase 6 (needs complete knowledge graph before improving perception)

**Research**: Unlikely (internal improvements to existing vision system)

**Plans**: 2 plans

Plans:
- [x] 6.1-01: Vision System Enhancements (LCD OCR, multi-frame capture, confidence)
- [x] 6.1-02: Data Stability (temporal smoothing, schema lock)

**Why this was inserted:**
Phase 7's hardware validation loop requires robust perception. Current v1.0 perception has limitations:
- ISS-001: Single-frame capture misses blinking LEDs
- LCD OCR needs robustness improvements
- No temporal smoothing for noisy observations
- Observation schema not locked (breaking changes possible)

**Required features:**
- ✅ LCD OCR solid
- ✅ Multi-object tracking (address ISS-001)
- ✅ Confidence calibration
- ✅ Temporal smoothing
- ✅ JSON schema lock

#### Phase 7: Code Generation Engine

**Goal**: LLM-based firmware generator with automatic hardware validation loop (generate → flash → observe → validate → iterate)

**Depends on**: Phase 6.1 (needs robust perception for hardware validation loop)

**Research**: Likely (LLM fine-tuning on embedded code, validation loop architecture, feedback engineering for failed validations)

**Research topics**: Claude API fine-tuning vs LoRA adapters, prompt engineering for BARR-C compliance, validation feedback loop design (how to give LLM actionable error context), iteration limits and fallback strategies

**Plans**: 2 plans

Plans:
- [x] 07-01: Code Generator + Pattern Retrieval (Claude API, prompt engineering, CLI)
- [x] 07-02: Validation Pipeline (style validation, pattern storage)

#### Phase 8: Public Launch

**Goal**: Polish UX, launch marketing campaign, ship to 200 paying users with "Better than Embedder" positioning

**Depends on**: Phase 7 (needs end-to-end code generation working)

**Research**: Unlikely (integration of existing components, marketing execution)

**Plans**: 2 plans

Plans:
- [x] 08-01: UX Polish + Documentation (error messages, help text, installation guide, getting started, examples)
- [ ] 08-02: Marketing + Launch Campaign (positioning, demo materials, blog post, HN launch, metrics tracking)

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 2.5 → 3 → 4 → 5 → 6 → 6.1 → 7 → 8

| Phase | Milestone | Plans Complete | Status | Completed |
|-------|-----------|----------------|--------|-----------|
| 1. Core + Vision | v1.0 | 3/3 | Complete | 2026-02-11 |
| 2. Assertions | v1.0 | 2/2 | Complete | 2026-02-11 |
| 2.5. Multi-LED Identity (INSERTED) | v1.0 | 1/1 | Complete | 2026-02-11 |
| 3. Diff + Firmware Tracking | v1.0 | 2/2 | Complete | 2026-02-11 |
| 4. Polish + Alpha | v1.0 | 2/2 | Complete | 2026-02-12 |
| 5. Style Infrastructure | v2.0 | 2/2 | Complete | 2026-02-12 |
| 6. Knowledge Graphs | v2.0 | 2/2 | Complete | 2026-02-12 |
| 6.1. Perception Enhancements (INSERTED) | v2.0 | 2/2 | Complete | 2026-02-13 |
| 7. Code Generation Engine | v2.0 | 2/2 | Complete | 2026-02-13 |
| 8. Public Launch | v2.0 | 1/2 | In Progress | - |
