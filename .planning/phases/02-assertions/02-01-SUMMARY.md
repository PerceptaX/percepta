---
phase: 02-assertions
plan: 01
subsystem: assertions
tags: [dsl, parser, validation, testing]
requires:
  - phase: 01-core-vision
    provides: [core-types, driver-interfaces, memory-storage, vision-driver]
provides:
  - Assertion types (LED, Display, Timing)
  - DSL parser for human-readable assertion syntax
  - Evaluation engine with confidence scores
  - AssertionResult with pass/fail and messages
affects: [02-02]
tech-stack:
  added: []
  patterns: [assertion-evaluation, case-insensitive-matching, contains-matching, graceful-failure]
key-files:
  created: [internal/assertions/types.go, internal/assertions/parser.go]
  modified: []
key-decisions:
  - "Case-insensitive LED matching with single-LED fallback for real-world usage"
  - "Display assertions use contains() instead of exact match (OCR is noisy)"
  - "Timing assertions fail gracefully with clear message if signal missing"
  - "10% tolerance on blink rate, ±5 tolerance on RGB values for robustness"
issues-created: []
duration: 2 min
completed: 2026-02-11
---

# Phase 2 Plan 1: DSL and Assertions Summary

**Assertion types with pragmatic matching: case-insensitive LED lookup + fallback, contains() for displays, graceful timing failures**

## Performance

- **Duration:** 2 min
- **Started:** 2026-02-11T13:22:33Z
- **Completed:** 2026-02-11T13:24:55Z
- **Tasks:** 2/2
- **Files created:** 2

## Accomplishments

- Assertion types for LED, Display, and Timing validation
- DSL parser for human-readable assertion syntax
- Evaluation logic with practical matching strategies
- AssertionResult with pass/fail, confidence, expected/actual, and messages
- Case-insensitive LED matching with single-LED fallback (handles "LED 'UNKNOWN'" scenario)
- Display assertions use contains() for noisy OCR
- Timing assertions fail gracefully when signal missing

## Task Commits

Each task was committed atomically:

1. **Task 1: Define assertion types and DSL structures** - `9b8fd54` (feat)
2. **Task 2: Implement DSL parser for assertion syntax** - `2c32ba9` (feat)

**Plan metadata:** `e58c6b7` (docs: complete plan)

## Files Created/Modified

**Created:**
- `internal/assertions/types.go` - Assertion types (LEDAssertion, DisplayAssertion, TimingAssertion) with Evaluate() methods
- `internal/assertions/parser.go` - DSL parser for LED/Display/Timing syntax

## Decisions Made

1. **Case-insensitive LED matching with fallback**: LED names matched case-insensitively, with fallback to single LED if name doesn't match. Addresses Phase 1's "UNKNOWN" LED issue.

2. **Display assertions use contains()**: Changed from exact match to `strings.Contains()` because OCR is noisy and exact matches are too brittle for real hardware.

3. **Timing assertions fail gracefully**: If BootTimingSignal missing, return clear message "did you capture from power-on?" instead of panicking.

4. **Tolerance for analog signals**: 10% tolerance on blink rate, ±5 tolerance on RGB values to handle real-world sensor noise.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Adjusted for non-pointer signal fields**
- **Found during:** Task 1 (types.go compilation)
- **Issue:** Plan assumed `Color` and `BlinkHz` were pointers, but core.LEDSignal uses struct values with zero-value semantics
- **Fix:** Changed nil checks to zero-value checks (RGB{0,0,0} and BlinkHz==0)
- **Files modified:** internal/assertions/types.go
- **Verification:** Build succeeds
- **Committed in:** 9b8fd54 (Task 1 commit)

---

**Total deviations:** 1 auto-fixed (blocking), 0 deferred
**Impact on plan:** Necessary fix to match actual signal types. No scope creep.

## Issues Encountered

None - execution proceeded smoothly with one compilation fix.

## Next Phase Readiness

Ready for 02-02: CLI assert command. Assertion engine operational with practical matching strategies, ready for command-line integration.

---
*Phase: 02-assertions*
*Completed: 2026-02-11*
