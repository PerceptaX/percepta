# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-02-11)

**Core value:** observe() must work reliably. If Percepta can accurately tell you "the LED is blinking at 1.98 Hz" with 95%+ confidence, everything else follows.

**Current focus:** Phase 6 — Knowledge Graphs (v2.0 Code Generation milestone)

## Current Position

Phase: 6 of 8 (Knowledge Graphs)
Plan: 06-01 complete (Knowledge Graph Storage)
Status: In progress - 1/2 plans complete
Last activity: 2026-02-13 — Plan 06-01 completed (Knowledge Graph Storage)

Progress: ██████░░░░ 62.5% (5/8 phases complete, Phase 6 in progress)

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
- Total plans completed: 3
- Status: Phase 5 complete, Phase 6 in progress (1/2 plans done)

**By Phase (v2.0):**

| Phase | Plans | Total | Avg/Plan |
|-------|-------|-------|----------|
| 5     | 2     | 90 min | 45 min   |
| 6     | 1/2   | 60 min | 60 min   |

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
| 5     | Use tree-sitter-c for Go instead of custom parser | Industry standard, robust, well-maintained C grammar |
| 5     | Checker interface pattern for extensible rule system | Allows adding new checkers easily, follows Go interface idioms |
| 5     | Global const uses UPPER_SNAKE, local const uses snake_case | BARR-C scope-aware naming - matches professional embedded coding standards |
| 5     | Descriptive error messages with auto-fix suggestions | Actionable feedback better than generic violations |
| 5     | Auto-fix only deterministic violations (naming, types) | Magic numbers and const correctness require manual review |
| 5     | Apply fixes in category order (types first, naming second) | Avoids breaking cascading replacements |
| 5     | Automatic #include <stdint.h> injection when types fixed | Ensures header available without manual intervention |
| 5     | Standard linter output format (file:line:col severity [rule] message) | Enables CI integration, familiar to developers |
| 5     | Directory traversal finds all .c and .h files recursively | Batch processing for entire codebases |
| 6     | In-memory graph with SQLite persistence (pure Go, matches Phase 3 decision) | Avoids external services, maintains zero-dependency architecture |
| 6     | Store only validated patterns (StyleCompliant=true AND has observation) | Quality moat - only code that works on real hardware |
| 6     | Full relationship graph: spec->pattern->board->observation->style_result | Enables context injection for code generation |
| 6     | Database path: ~/.local/share/percepta/knowledge.db (alongside percepta.db) | Separates knowledge from perception data |
| 6     | PatternStore integrates StyleChecker, Graph, and SQLite storage | Single API for validated pattern storage |
| 6     | Reject patterns without observation (hardware validation required) | Ensures patterns are hardware-verified, not theoretical |

### Deferred Issues

- **ISS-001**: Single-frame capture misses blinking LEDs (discovered Phase 2.5). Object permanence IS working (LED1 = LED1), but single frame only captures LEDs that are ON at that instant. Multi-frame/video capture needed for complete LED detection. Deferred to post-v2.0.

### Blockers/Concerns Carried Forward

None - starting fresh with v2.0 milestone.

### Roadmap Evolution

- Phase 2.5 inserted after Phase 2: Fix multi-LED signal identity extraction (BLOCKING - required before Phase 3)
- Milestone v2.0 Code Generation created: AI firmware generation with hardware validation, 4 phases (Phase 5-8)

## Session Continuity

Last session: 2026-02-13T00:30:00Z
Stopped at: Phase 5 complete (Style Infrastructure)
Resume file: None

**Next:** Phase 6 planning (Knowledge Graphs) - break down into plans
