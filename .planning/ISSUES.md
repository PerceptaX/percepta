# Deferred Issues

Issues logged during execution for future consideration.

---

## ISS-001: Single-frame capture misses blinking LEDs

**Discovered:** Phase 2.5 (multi-LED identity verification)  
**Severity:** Medium (limits observability but doesn't break core functionality)  
**Status:** Deferred to post-Phase 3

**Problem:**
Percepta captures a single frame, so only detects LEDs that are ON at that exact moment. LEDs blinking at different frequencies may be OFF during capture and become invisible.

**Current behavior:**
- FPGA has 3 LEDs blinking at different rates
- Single frame captures LED1 (ON) but misses LED2, LED3 (OFF at that instant)
- Subsequent captures may see different LEDs depending on timing

**Impact:**
- Cannot reliably detect all LEDs on a multi-LED board
- Diff may compare different subsets of LEDs
- Assertions cannot validate LEDs that happen to be OFF during capture

**Potential solutions:**
1. **Multi-frame capture** - Capture 3-5 frames over 2-3 seconds, merge signals
2. **Video analysis** - Capture 2-second video clip, extract all LEDs across frames
3. **Smart sampling** - Detect blink frequencies and time captures to catch all LEDs

**Why deferred:**
- Requires architectural change to camera driver (single frame → multi-frame/video)
- Phase 3 diff can work with whatever LEDs ARE captured (consistent subset)
- Object permanence IS working (LED1 = LED1 consistently)
- Can be addressed after core diff functionality proven

**Workaround for now:**
Test with boards where LEDs are solid or synchronously blinking, or accept that only visible LEDs are tracked.

---

## Next Issue: ISS-002
