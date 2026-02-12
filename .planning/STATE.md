# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-02-11)

**Core value:** observe() must work reliably. If Percepta can accurately tell you "the LED is blinking at 1.98 Hz" with 95%+ confidence, everything else follows.

**Current focus:** Phase 5 — Style Infrastructure (v2.0 Code Generation milestone)

## Current Position

Phase: 5 of 8 (Style Infrastructure)
Plan: Not started
Status: Ready to plan
Last activity: 2026-02-12 — Milestone v2.0 Code Generation created

Progress: ░░░░░░░░░░ 0% (0/? plans in v2.0)

## Performance Metrics

**v1.0 Perception MVP (COMPLETED):**
- Total plans completed: 10
- Average duration: ~10 min
- Total execution time: 1.7 hours
- Shipped: 2026-02-12

**By Phase:**

| Phase | Plans | Total | Avg/Plan |
|-------|-------|-------|----------|
| 1     | 3     | 8 min | 2.7 min  |
| 2     | 2     | 3 min | 1.5 min  |
| 2.5   | 1     | 1 min | 1.0 min  |
| 3     | 2     | 75 min | 37.5 min |
| 4     | 2     | 7 min | 3.5 min  |

**v2.0 Code Generation (IN PROGRESS):**
- Total plans completed: 0
- Status: Starting Phase 5

## Accumulated Context

### Decisions

Decisions are logged in PROJECT.md Key Decisions table.
Historical decisions from v1.0:

| Phase | Decision | Rationale |
|-------|----------|-----------|
| 01    | Platform-agnostic interfaces (CameraDriver returns JPEG bytes) | Enables Linux/macOS/Windows without refactor |
| 01    | No StorageDriver interface | Premature abstraction - only MemoryStorage exists |
| 01    | In-memory storage for MVP | Focus on observe() accuracy, defer SQLite |
| 01    | Parser isolated behind SignalParser interface | Enables swap to structured output later |
| 01    | Regex for MVP signal parsing | Sufficient for LED/Display, replace when tool use stable |
| 01    | Human-readable output over JSON | Better alpha UX, JSON export later |
| 01    | Config file optional with defaults | Works without ~/.config/percepta/config.yaml |
| 02    | Case-insensitive LED matching with fallback | Handles real-world hardware (addresses "UNKNOWN" LED from Phase 1) |
| 02    | Display assertions use contains() not exact match | OCR is noisy, exact match too brittle |
| 02    | Timing assertions fail gracefully if signal missing | Better UX than panic, clear message to user |
| 02    | 10% tolerance on blink rate, ±5 on RGB | Handles real-world sensor noise |
| 2.5   | Index-based LED naming (LED1, LED2, LED3) | Establishes object permanence - stable identity enables diff |
| 2.5   | No spatial tracking in MVP | Appearance order sufficient, spatial clustering can be added later |
| 3     | Manual firmware tags (NOT git auto-integration) | Git coupling breaks FPGA workflows, binaries, CI, non-repo users |
| 3     | Use modernc.org/sqlite (NOT mattn/go-sqlite3) | Pure Go, zero CGO dependencies, maintains cross-platform architecture |
| 3     | Exact diff (NO tolerances except BlinkHz normalization) | Assertions handle fuzz, diff must be deterministic |
| 3     | Storage construction in cmd layer | pkg/percepta stays framework-agnostic with StorageDriver interface |
| 4     | Added yaml struct tags to DeviceConfig | Viper requires yaml tags for marshaling (separate from mapstructure tags) |

### Deferred Issues

- **ISS-001**: Single-frame capture misses blinking LEDs (discovered Phase 2.5). Object permanence IS working (LED1 = LED1), but single frame only captures LEDs that are ON at that instant. Multi-frame/video capture needed for complete LED detection. Deferred to post-v2.0.

### Blockers/Concerns Carried Forward

None - starting fresh with v2.0 milestone.

### Roadmap Evolution

- Phase 2.5 inserted after Phase 2: Fix multi-LED signal identity extraction (BLOCKING - required before Phase 3)
- Milestone v2.0 Code Generation created: AI firmware generation with hardware validation, 4 phases (Phase 5-8)

## Session Continuity

Last session: 2026-02-12T22:41:00Z
Stopped at: Milestone v2.0 Code Generation initialization
Resume file: None

**Ready to plan Phase 5: Style Infrastructure**
