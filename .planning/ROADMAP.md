# Roadmap: Percepta

## Overview

Percepta's 8-week journey from zero to alpha builds the perception kernel in four phases: establish the vision foundation (Core + Vision), add validation capabilities (Assertions), enable version comparison (Diff + Firmware Tracking), and ship to alpha users (Polish + Alpha). Each phase delivers a complete, verifiable capability that builds toward the core value: reliable observe() accuracy.

## Domain Expertise

None

## Phases

**Phase Numbering:**
- Integer phases (1, 2, 3): Planned milestone work
- Decimal phases (2.1, 2.2): Urgent insertions (marked with INSERTED)

Decimal phases appear between their surrounding integers in numeric order.

- [x] **Phase 1: Core + Vision** - Foundation types, in-memory storage, Claude Vision driver, observe command
- [ ] **Phase 2: Assertions** - DSL parser, LED/display/timing assertions, assert command
- [ ] **Phase 3: Diff + Firmware Tracking** - Firmware hash tracking, observation comparison, diff command
- [ ] **Phase 4: Polish + Alpha** - Device management, documentation, installation, alpha release

## Phase Details

### Phase 1: Core + Vision
**Goal**: percepta observe <device> works end-to-end with 95%+ accuracy on LED/display/boot signals

**Depends on**: Nothing (first phase)

**Research**: Likely (Claude Vision API integration, Go camera capture)

**Research topics**: Anthropic Go SDK usage, Go webcam libraries (gocv vs native options), SQLite schema design for time-series observations

**Plans**: 2-3 plans

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

### Phase 3: Diff + Firmware Tracking
**Goal**: percepta diff --from X --to Y compares behavior across firmware versions

**Depends on**: Phase 2 (needs observations + assertions)

**Research**: Unlikely (query patterns on existing SQLite schema)

**Plans**: 2 plans

Plans:
- [ ] 03-01: Firmware hash tracking and git integration
- [ ] 03-02: Observation comparison and diff command

### Phase 4: Polish + Alpha
**Goal**: Ship to 10 alpha users with installation in <10 minutes

**Depends on**: Phase 3 (needs all three verbs working)

**Research**: Unlikely (installation tooling, documentation)

**Plans**: 2 plans

Plans:
- [ ] 04-01: Device management commands and manual camera config
- [ ] 04-02: Documentation, Homebrew formula, alpha release

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Core + Vision | 3/3 | Complete | 2026-02-11 |
| 2. Assertions | 2/2 | Complete | 2026-02-11 |
| 3. Diff + Firmware Tracking | 0/2 | Not started | - |
| 4. Polish + Alpha | 0/2 | Not started | - |
