---
phase: 01-core-vision
plan: 02
subsystem: vision
tags: [anthropic-sdk, claude-vision, v4l2, webcam, regex-parser]
requires:
  - phase: 01-01
    provides: [core-types, driver-interfaces]
provides: [vision-driver, camera-driver, signal-parser]
affects: [01-03]
tech-stack:
  added: [github.com/anthropics/anthropic-sdk-go@v1.22.1, github.com/blackjack/webcam@v0.6.1]
  patterns: [parser-isolation, platform-specific-implementations]
key-files:
  created: [internal/camera/v4l2.go, internal/vision/claude.go, internal/vision/parser.go]
  modified: []
key-decisions:
  - "Camera behind CameraDriver interface: V4L2 for Linux, future: AVFoundation (macOS), Media Foundation (Windows)"
  - "Vision accepts frame bytes: Platform-agnostic, no coupling to camera API"
  - "Parser isolated behind SignalParser interface: Enables clean swap to structured output later"
  - "Regex for MVP: LED/Display extraction sufficient, will replace with tool use when stable"
  - "MJPEG format detection: Query supported formats instead of hardcoding constant"
issues-created: []
duration: 4 min
completed: 2026-02-11
---

# Phase 1 Plan 2: Vision Engine Summary

**Claude Vision API integration with Sonnet 4.5, Linux V4L2 camera capture, and isolated regex signal parser**

## Performance

- **Duration:** 4 min
- **Started:** 2026-02-11T18:13:42Z
- **Completed:** 2026-02-11T18:17:51Z
- **Tasks:** 2/2
- **Files created:** 3

## Accomplishments

- Linux V4L2 camera driver implemented behind core.CameraDriver interface
- MJPEG 1280x720 capture with format detection and proper device cleanup
- Claude Vision API integration with Anthropic SDK v1.22.1 behind core.VisionDriver interface
- Signal parser isolated behind SignalParser interface (regex for MVP, swappable)
- Platform-agnostic architecture: camera returns JPEG bytes, vision accepts bytes
- Color extraction scans matched LED segment only (prevents false positives on multi-LED boards)

## Task Commits

Each task was committed atomically:

1. **Task 1: Implement Linux V4L2 camera driver** - `003b4b3` (feat)
2. **Task 2: Implement Claude Vision driver with isolated parser** - `ca03bf1` (feat)

**Plan metadata:** (to be added after commit)

## Files Created/Modified

**Created:**
- `internal/camera/v4l2.go` - Linux V4L2 implementation of core.CameraDriver
- `internal/vision/claude.go` - Claude Vision implementation of core.VisionDriver
- `internal/vision/parser.go` - Isolated regex signal parser (RegexParser)

**Modified:**
- `go.mod`, `go.sum` - Added anthropic-sdk-go and webcam dependencies

## Decisions Made

1. **Camera behind interface**: V4L2 implementation in internal/camera. Future: Add avfoundation.go (macOS), mf.go (Windows) without refactor.

2. **Vision accepts frame bytes**: Platform-agnostic. No coupling to specific camera API.

3. **Parser isolation**: SignalParser interface allows clean swap to structured output prompting later.

4. **Regex for MVP**: Sufficient for LED/Display extraction. Will replace with tool use when structured output is reliable.

5. **MJPEG format detection**: Query GetSupportedFormats() instead of hardcoding V4L2_PIX_FMT_MJPEG constant (API compatibility).

6. **Anthropic SDK v1.22.1 API**: Use NewImageBlockBase64() and NewTextBlock() helper functions for correct ContentBlockParamUnion construction.

## Deviations from Plan

### Auto-fixed Issues

**1. [Rule 3 - Blocking] Fixed MJPEG constant reference**
- **Found during:** Task 1 (V4L2 camera implementation)
- **Issue:** webcam.V4L2_PIX_FMT_MJPEG constant undefined - build failed
- **Fix:** Query GetSupportedFormats() to detect MJPEG dynamically instead of hardcoding constant
- **Files modified:** internal/camera/v4l2.go
- **Verification:** Build succeeds, format detection logic works
- **Committed in:** 003b4b3 (Task 1 commit)

**2. [Rule 3 - Blocking] Fixed Anthropic SDK API usage**
- **Found during:** Task 2 (Claude Vision implementation)
- **Issue:** SDK v1.22.1 API different from plan pseudocode - multiple compilation errors
- **Fix:** Use NewImageBlockBase64() and NewTextBlock() helper functions, correct struct field names (OfImage, OfText)
- **Files modified:** internal/vision/claude.go
- **Verification:** Build succeeds, API call structure correct
- **Committed in:** ca03bf1 (Task 2 commit)

---

**Total deviations:** 2 auto-fixed (2 blocking), 0 deferred
**Impact on plan:** Both auto-fixes necessary to unblock compilation. No scope creep - just API compatibility adjustments.

## Issues Encountered

- Network timeouts during `go get` - retried successfully
- Anthropic SDK v1.22.1 API differs from typical REST patterns - used helper functions for correct union type construction

## Next Phase Readiness

**Ready for 01-03: CLI observe command**

Can now capture hardware state and convert to structured observations:
- CameraDriver interface complete with V4L2 implementation
- VisionDriver interface complete with Claude Sonnet 4.5 integration
- SignalParser extracting LED/Display signals from Vision API responses
- Platform-agnostic JPEG bytes flowing through the system

No blockers. Vision engine operational, ready for CLI integration.

---
*Phase: 01-core-vision*
*Completed: 2026-02-11*
