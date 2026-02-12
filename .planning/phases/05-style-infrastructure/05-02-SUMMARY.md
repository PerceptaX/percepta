---
phase: 05-style-infrastructure
plan: 02
subsystem: code-quality
tags: [barr-c, auto-fix, cli, cobra, embedded-c, code-standards]

# Dependency graph
requires:
  - phase: 05-01
    provides: BARR-C rule engine, StyleChecker, NamingChecker, TypesChecker, tree-sitter parser
provides:
  - Auto-fix engine for deterministic BARR-C violations (naming, types)
  - StyleFixer orchestrator with NamingFixer and TypesFixer
  - percepta style-check CLI command with --fix flag
  - Automatic #include <stdint.h> injection
  - Standard linter output format (file:line:col severity [rule] message)
affects: [06-knowledge-graphs, 07-code-generation, firmware-validation, ci-integration]

# Tech tracking
tech-stack:
  added: [cobra CLI framework integration for style command]
  patterns: [Fixer interface pattern, line-accurate source replacement, CLI error handling with exit codes]

key-files:
  created:
    - internal/style/fixer.go
    - internal/style/fixer_test.go
    - cmd/percepta/style.go
  modified:
    - cmd/percepta/main.go

key-decisions:
  - "Auto-fix only deterministic violations (naming, types), not magic numbers or const correctness"
  - "Apply fixes in category order (types first, naming second) to avoid breaking cascading fixes"
  - "Automatic #include <stdint.h> injection when type fixes applied and header missing"
  - "Standard linter output format for CI integration: file:line:col severity [rule] message"
  - "Exit code 0 if clean, non-zero if violations remain (CI-friendly)"
  - "Directory traversal finds all .c and .h files recursively"

patterns-established:
  - "Fixer interface: Fix(violation, source) ([]byte, bool) for all fixers"
  - "StyleFixer orchestration: applies fixes in category order with fix tracking"
  - "Header injection: insert #include at file top or after existing includes"
  - "CLI output: clear summary with fix count and remaining violations"
  - "Error messages: actionable suggestions with arrows → for fixes"

issues-created: []

# Metrics
duration: 45min
completed: 2026-02-12
---

# Phase 05-02: Auto-Fix Engine and CLI Summary

**Auto-fix engine with NamingFixer and TypesFixer, plus percepta style-check CLI command with --fix flag**

## Performance

- **Duration:** 45 min
- **Started:** 2026-02-12T23:45:00Z
- **Completed:** 2026-02-13T00:30:00Z
- **Tasks:** 2
- **Files modified:** 5

## Accomplishments
- Auto-fix engine applies deterministic corrections (function naming, type replacements)
- percepta style-check CLI command validates C code against BARR-C standard
- --fix flag auto-corrects violations with clear reporting
- Automatic #include <stdint.h> injection when types fixed
- Directory traversal for batch checking all .c/.h files
- Standard linter output format for CI integration

## Task Commits

Each task was committed atomically:

1. **Task 1: Build auto-fix engine** - `081d57c` (feat)
2. **Task 2: Add percepta style-check CLI command** - `c50164f` (feat)

## Files Created/Modified
- `internal/style/fixer.go` - Auto-fix engine with NamingFixer, TypesFixer, StyleFixer orchestrator
- `internal/style/fixer_test.go` - Comprehensive tests for all fixers (13 tests, all passing)
- `cmd/percepta/style.go` - CLI command for style checking with --fix flag
- `cmd/percepta/main.go` - Registered style-check command

## Auto-Fix Strategy

### Deterministic Fixes (Automated)
1. **Function naming**: Converts to Module_Function format (e.g., initLED → Module_InitLED)
2. **Type replacements**: unsigned char → uint8_t, unsigned short → uint16_t, etc.
3. **Header injection**: Adds #include <stdint.h> when type fixes applied

### Manual Review Required (Not Automated)
1. **Magic numbers**: Require semantic constant names (can't auto-generate)
2. **Const correctness**: Affects function signatures, needs careful review
3. **Variable naming**: Too risky without full context (scope, purpose)

## CLI Usage Examples

### Check violations only
```bash
percepta style-check firmware.c
```

Output:
```
firmware.c:
  5:1 error [Function Naming Convention] Function 'initLED' should use Module_Function format
    → Module_InitLED
  7:5 error [Stdint Type Usage] Type 'unsigned char' should be replaced with stdint.h type
    → Replace 'unsigned char brightness' with 'uint8_t brightness'

⚠️  2 violation(s) remain in 1 file(s).

Run with --fix to auto-correct fixable violations.
```

### Auto-fix violations
```bash
percepta style-check firmware.c --fix
```

Output:
```
Fixed in firmware.c:
  ✓ firmware.c:7:5 - Fixed: Stdint Type Usage
  ✓ firmware.c:5:1 - Fixed: Function Naming Convention

✅ Fixed 2 violation(s) automatically.
```

### Check entire directory
```bash
percepta style-check ./src --fix
```

Recursively processes all .c and .h files in the directory.

### Clean code validation
```bash
percepta style-check firmware.c
# Exit code 0
✅ No style violations found. Code is BARR-C compliant.
```

## Implementation Details

### NamingFixer
- Extracts function name from violation message
- Uses suggestion directly (e.g., "Module_InitLED")
- Applies word-boundary regex replacement to entire source
- Replaces ALL occurrences (function definitions and calls)

### TypesFixer
- Parses suggestion format: "Replace 'old type var' with 'new type var'"
- Applies line-accurate replacement (only fixes specific line)
- Handles all stdint.h mappings (uint8_t, uint16_t, uint32_t, int8_t, etc.)

### StyleFixer Orchestrator
- Applies fixes in category order: types first, naming second
- Tracks all fixes applied with file:line:col details
- Returns fixed source + list of applied fixes

### Header Injection
- Detects if #include <stdint.h> already present
- Inserts after last existing #include, or before first code line
- Only triggers if type fixes were applied

### CLI Implementation
- Standard Cobra command pattern matching other percepta commands
- File or directory detection with filepath.Walk
- Clear violation reporting with severity and suggestions
- Exit code 0 if clean, non-zero if violations remain (CI-friendly)

## Decisions Made

1. **Fix order matters**: Apply type fixes before naming fixes to avoid breaking cascading replacements
2. **Line-accurate type fixes**: Only replace on specific violation line to avoid over-aggressive changes
3. **Global function name replacement**: Safe to replace all occurrences since functions have unique names
4. **Header insertion strategy**: Place at top of file (after comments) or after existing includes
5. **Variable naming NOT auto-fixed**: Too risky without full context about scope and purpose
6. **CLI output format**: Matches standard linters for CI integration potential

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 1 - Bug] Fixed header insertion placing #include inside functions**
- **Found during:** Task 2 (CLI testing)
- **Issue:** Header was being inserted at first non-comment line, which could be inside a function
- **Fix:** Changed logic to insert after last existing #include, or before first code line
- **Files modified:** internal/style/fixer.go
- **Verification:** Manual testing with test C files, header now correctly placed at file top
- **Committed in:** c50164f (Task 2 commit)

---

**Total deviations:** 1 auto-fixed (bug fix)
**Impact on plan:** Bug fix essential for correct header injection. No scope creep.

## Issues Encountered

**1. Suggestion format mismatch between fixer and checker**
- **Problem**: Initial fixer expected "Replace with: uint8_t" but checker provides "Replace 'unsigned char x' with 'uint8_t x'"
- **Resolution**: Updated fixer regex patterns to match actual checker output format
- **Verification**: All tests passing (13/13)

**2. Header insertion logic bug**
- **Problem**: #include <stdint.h> was being inserted at first non-comment line, which could be inside a function body
- **Resolution**: Changed logic to find last existing #include or insert before first code line
- **Verification**: Manual testing with multiple C file formats

## Test Coverage

All tests passing (13/13 in internal/style package):

### Fixer Tests
- `TestTypesFixer_BasicFix` - Single type replacement
- `TestNamingFixer_FunctionRename` - Global function name replacement
- `TestStyleFixer_ApplyFixes` - Multiple fix orchestration
- `TestStyleFixer_EnsureStdintHeader` - Header injection logic (3 scenarios)
- `TestFixer_IntegrationWithChecker` - End-to-end checker + fixer workflow

### Manual CLI Tests
- Single file check: detects violations correctly ✓
- Single file fix: applies corrections and reports remaining ✓
- Directory check: traverses and processes all C files ✓
- Directory fix: batch fixes all violations ✓
- Clean code: reports BARR-C compliant status ✓
- Exit codes: 0 for clean, non-zero for violations ✓

## Next Phase Readiness

**Ready for Phase 6 (Knowledge Graphs):**
- Style checker can validate generated code before storing in knowledge base
- Auto-fix ensures stored patterns follow BARR-C standards
- CLI command available for manual validation workflows
- Standard output format ready for CI integration

**Blockers:** None

**Integration points for Phase 6:**
- Use StyleChecker.CheckSource() to validate patterns before storage
- Use StyleFixer.ApplyFixes() to ensure all patterns are BARR-C compliant
- Store violation-free patterns in knowledge graph with confidence

## Phase 5 Completion

**Phase 5 is now complete (both plans done):**
- ✅ 05-01: BARR-C rule engine with NamingChecker and TypesChecker
- ✅ 05-02: Auto-fix engine and CLI command

**Total time:** 90 min (45 min + 45 min)
**Total files created:** 9
**Total tests:** 21 (all passing)

**Ready for:** Phase 6 - Knowledge Graphs

---
*Phase: 05-style-infrastructure*
*Completed: 2026-02-12*
