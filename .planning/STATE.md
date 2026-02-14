# Project State

## Project Reference

See: .planning/PROJECT.md (updated 2026-02-11)

**Core value:** observe() must work reliably. If Percepta can accurately tell you "the LED is blinking at 1.98 Hz" with 95%+ confidence, everything else follows.

**Current focus:** v2.0 Code Generation milestone COMPLETE — Phase 8 complete, ready for public launch

## Current Position

Phase: 8 of 8 (Public Launch) — COMPLETE
Plan: 2/2 complete
Status: Phase 8 complete, v2.0 milestone COMPLETE
Last activity: 2026-02-13 — Plan 08-02 complete (Marketing + Launch Campaign)

Progress: ██████████ 100% (All phases complete, v2.0 ready for public launch)

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

**v2.0 Code Generation (COMPLETE):**
- Total plans completed: 14
- Status: All phases complete (5, 6, 6.1, 7, 8)
- Milestone: SHIPPED 2026-02-13
- Total execution time: 4.6 hours

**By Phase (v2.0):**

| Phase | Plans | Total | Avg/Plan |
|-------|-------|-------|----------|
| 5     | 2/2   | 90 min | 45 min   |
| 6     | 2/2   | 150 min | 75 min   |
| 6.1   | 2/2   | 90 min | 45 min   |
| 7     | 2/2   | 80 min | 40 min   |
| 8     | 2/2   | 29 min | 14.5 min |

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
| 6     | In-memory vector store + SQLite persistence (pure Go, matches Phase 3 decision) | Ship faster without external services, can upgrade to Qdrant later |
| 6     | OpenAI embeddings API (text-embedding-ada-002) for semantic similarity | Industry standard, proven accuracy, pluggable architecture for local models later |
| 6     | Cosine similarity for pattern matching | Simple, effective, well-understood for MVP |
| 6     | Mock embedder for testing (NewVectorStoreWithEmbedder) | Enables deterministic testing without API keys |
| 6     | Confidence scoring = similarity + signal boost | Combines vector similarity with validation metadata |
| 6     | CLI graceful degradation when OPENAI_API_KEY not set | Pattern storage works without vector store, semantic search fails gracefully |
| 6.1   | Use Anthropic tool use for structured output | Deterministic LCD OCR extraction, eliminates regex brittleness |
| 6.1   | Keep RegexParser as fallback | Graceful degradation when tool use fails, maintains robustness |
| 6.1   | 5 frames over 1 second for multi-frame capture | Balances completeness (detects all LEDs) with latency (1s acceptable) |
| 6.1   | Calibrate confidence dynamically | Adjust scores based on detection rate, color presence, text quality |
| 6.1   | Blink frequency from transition count | Simple algorithm works for typical embedded LED rates (0.5-5 Hz) |
| 6.1   | 5-second time window with 2/3 agreement for temporal smoothing | Balances noise filtering with state change detection |
| 6.1   | Schema version locked at 1.0.0 with migration framework | Future-proofs for schema changes, ensures compatibility |
| 6.1   | Graceful degradation on storage/validation failures | Smoothing returns unfiltered, validation logs warnings but continues |
| 7     | Use Anthropic SDK directly (anthropic-sdk-go) | Already in dependencies, simplifies implementation vs custom HTTP client |
| 7     | Model: claude-sonnet-4-5-20250929 (latest Sonnet 4.5) | Best balance of performance and quality for code generation |
| 7     | Temperature 0.3 for code generation | Lower than default (1.0) for deterministic, consistent code |
| 7     | Max tokens 4096 for firmware code | Suitable for typical firmware (50-200 lines) |
| 7     | Top 3 similar patterns in prompt | Balances context richness with prompt length and cost |
| 7     | Code truncation at 50 lines per example | Prevents over-long prompts while providing useful context |
| 7     | Graceful degradation without semantic search | Generate code with BARR-C requirements only if OPENAI_API_KEY not set |
| 7     | Board-specific API guidance hardcoded | Common boards (ESP32, STM32, Arduino) to prevent API mistakes |
| 7     | Style validation only for MVP | Hardware validation deferred to Phase 8 - faster shipping while maintaining BARR-C compliance |
| 7     | Auto-fix violations automatically | Deterministic fixes (naming, types) applied without user confirmation - speeds up workflow |
| 7     | Store only compliant patterns | Patterns must pass style check to be stored - ensures knowledge graph quality |
| 7     | Graceful degradation on storage failure | Generation succeeds even if pattern storage fails - user gets code regardless |
| 7     | Device ID from config first device | Simple MVP approach - load first device from config.Devices or fallback to 'unknown-device' |
| 7     | Detailed validation reporting | Transparent reporting of validation steps - shows style status, auto-fix, pattern storage |
| 7     | Iterate count = 1 for MVP | Future enhancement: multi-iteration refinement with feedback loop |
| 7     | Re-check style after auto-fix | Shows remaining violations after fixes - user knows what needs manual attention |
| 8     | UserError type with structured fields | Message/Suggestion/DocsURL pattern for consistent, actionable error messages |
| 8     | Progress spinners use stderr | Avoids polluting stdout, enables piping command output |
| 8     | Help text includes Examples section | Users learn by copy-paste; show exact commands to run |
| 8     | Documentation organized by user journey | Installation → Getting Started → Commands → Examples → Configuration → Troubleshooting → API |
| 8     | 25+ example workflows | Covers basic usage, firmware tracking, code generation, CI/CD, and advanced scenarios |
| 8     | Marketing materials in docs/launch/ directory | Centralized location for all launch materials (user-specified) |
| 8     | "Better than Embedder" positioning | Hardware validation + BARR-C compliance vs compilation-only tools |
| 8     | Complement Embedder (not attack) | "Prototype with Embedder, productionize with Percepta" - builds credibility |
| 8     | Honest benchmark reporting | Report real tradeoffs: 45s validated vs 10s unvalidated, but guaranteed working |
| 8     | Launch timing: Tuesday-Thursday 9-11am PT | Optimal HN visibility, sustained engagement first 4 hours |
| 8     | Metrics aligned with PRD Part VIII | 1500 WAU, $10K MRR, 200 paying customers at Month 12 |

### Deferred Issues

- **ISS-001**: ✅ RESOLVED in Phase 6.1 Plan 01. Multi-frame capture (5 frames, 200ms interval) detects all LEDs including blinking ones. Object permanence maintained.

### Blockers/Concerns Carried Forward

None - starting fresh with v2.0 milestone.

### Roadmap Evolution

- Phase 2.5 inserted after Phase 2: Fix multi-LED signal identity extraction (BLOCKING - required before Phase 3)
- Milestone v2.0 Code Generation created: AI firmware generation with hardware validation, 4 phases (Phase 5-8)
- Phase 6.1 inserted after Phase 6: Perception Enhancements (URGENT - required before Phase 7 validation loop). Addresses ISS-001, adds LCD OCR robustness, temporal smoothing, and schema stability.

## Session Continuity

Last session: 2026-02-15T09:30:00Z
Stopped at: v2.0.0 RELEASED - GitHub release published
Resume file: None

**Completed:**
- ✅ Synced dev with main
- ✅ Tagged v2.0.0 and pushed to origin
- ✅ Created GitHub release: https://github.com/Perceptax/percepta/releases/tag/v2.0.0

**Next:** Manual launch tasks (HN/Reddit posts, metrics monitoring)
**Status:** v2.0.0 publicly released, awaiting community launch announcements
