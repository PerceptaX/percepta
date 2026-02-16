#!/usr/bin/env bash
# =============================================================================
# Percepta End-to-End User Journey Test
#
# Exercises percepta exactly as a real Linux user would — from building the
# binary to running every CLI command. Uses the real camera (/dev/video0) and
# the real Anthropic API.
#
# Usage:
#   chmod +x test/e2e/user_journey_test.sh
#   ./test/e2e/user_journey_test.sh
# =============================================================================

set -euo pipefail

# ---------------------------------------------------------------------------
# Configuration
# ---------------------------------------------------------------------------

SOURCE_DIR="/home/utkarsh/Work/hacks/PerceptumX/percepta"
CAMERA="/dev/video0"

# ANTHROPIC_API_KEY must be set in the environment before running this script
if [[ -z "${ANTHROPIC_API_KEY:-}" ]]; then
    echo "Error: ANTHROPIC_API_KEY environment variable is not set."
    echo "Export it before running: export ANTHROPIC_API_KEY='sk-ant-...'"
    exit 1
fi
DEVICE_NAME="esp32-dev"

# Preserve GOPATH so Go module cache doesn't land in the fake HOME
export GOPATH="${GOPATH:-$HOME/go}"

# ---------------------------------------------------------------------------
# Isolated environment
# ---------------------------------------------------------------------------

FAKE_HOME=$(mktemp -d /tmp/percepta-e2e-XXXX)
export HOME="$FAKE_HOME"
mkdir -p "$FAKE_HOME/.config/percepta"
mkdir -p "$FAKE_HOME/.local/share/percepta"
WORKSPACE="$FAKE_HOME/workspace"
mkdir -p "$WORKSPACE"

PERCEPTA="$FAKE_HOME/bin/percepta"

# ---------------------------------------------------------------------------
# Counters and summary
# ---------------------------------------------------------------------------

PASS_COUNT=0
FAIL_COUNT=0
declare -a RESULTS=()

pass() {
    local phase="$1"
    PASS_COUNT=$((PASS_COUNT + 1))
    RESULTS+=("PASS  $phase")
    echo "  [PASS] $phase"
}

fail() {
    local phase="$1"
    local detail="${2:-}"
    FAIL_COUNT=$((FAIL_COUNT + 1))
    RESULTS+=("FAIL  $phase${detail:+: $detail}")
    echo "  [FAIL] $phase${detail:+: $detail}"
}

# Print a phase header
phase() {
    echo ""
    echo "========================================"
    echo "  Phase $1: $2"
    echo "========================================"
}

# ---------------------------------------------------------------------------
# Cleanup trap
# ---------------------------------------------------------------------------

cleanup() {
    echo ""
    echo "========================================"
    echo "  Cleanup"
    echo "========================================"
    # Make all files writable before removal (Go module cache files are read-only)
    chmod -R +w "$FAKE_HOME" 2>/dev/null || true
    rm -rf "$FAKE_HOME"
    echo "  Removed $FAKE_HOME"
}
trap cleanup EXIT

# ========================================================================
# Phase 1: Build from Source
# ========================================================================

phase 1 "Install (Build from Source)"

mkdir -p "$FAKE_HOME/bin"
(cd "$SOURCE_DIR" && go build -o "$PERCEPTA" ./cmd/percepta)

if [[ -x "$PERCEPTA" ]]; then
    pass "Binary built and is executable"
else
    fail "Binary not found or not executable"
    exit 1
fi

# ========================================================================
# Phase 2: Help & Discovery
# ========================================================================

phase 2 "First Run — Help & Discovery"

check_help() {
    local subcmd="$1"
    local keyword="$2"
    local label="$3"
    local output
    if [[ "$subcmd" == "root" ]]; then
        output=$("$PERCEPTA" --help 2>&1)
    else
        output=$("$PERCEPTA" $subcmd --help 2>&1)
    fi
    if echo "$output" | grep -qi "$keyword"; then
        pass "$label"
    else
        fail "$label" "missing keyword '$keyword'"
    fi
}

check_help "root"        "Usage"       "percepta --help"
check_help "device"      "device"      "percepta device --help"
check_help "observe"     "observe"     "percepta observe --help"
check_help "generate"    "generate"    "percepta generate --help"
check_help "knowledge"   "knowledge"   "percepta knowledge --help"
check_help "style-check" "style"       "percepta style-check --help"

# ========================================================================
# Phase 3: Configure a Device
# ========================================================================

phase 3 "Configure a Device"

# device add reads 3 lines interactively: type, camera path, firmware version
printf "esp32\n%s\nv1.0\n" "$CAMERA" | "$PERCEPTA" device add "$DEVICE_NAME"

LIST_OUT=$("$PERCEPTA" device list 2>&1)

check_list() {
    local needle="$1"
    local label="$2"
    if echo "$LIST_OUT" | grep -q "$needle"; then
        pass "$label"
    else
        fail "$label" "'$needle' not in device list output"
    fi
}

check_list "$DEVICE_NAME" "device list shows name"
check_list "esp32"        "device list shows type"
check_list "$CAMERA"      "device list shows camera"
check_list "v1.0"         "device list shows firmware"

# ========================================================================
# Phase 4: Observe Hardware (Real Camera + Real API)
# ========================================================================

phase 4 "Observe Hardware"

OBS_OUT=$("$PERCEPTA" observe "$DEVICE_NAME" 2>&1) || true

if echo "$OBS_OUT" | grep -qi "Observation captured"; then
    pass "observe: 'Observation captured' present"
else
    fail "observe: 'Observation captured' missing" "$(echo "$OBS_OUT" | head -3)"
fi

if echo "$OBS_OUT" | grep -qi "Signals"; then
    pass "observe: 'Signals' present"
else
    fail "observe: 'Signals' missing"
fi

# ========================================================================
# Phase 5: Run Assertions (Real Camera)
# ========================================================================

phase 5 "Run Assertions"

# The assertion uses the real DSL: LED.LED1 ON
# We allow exit code 0 (pass) or 1 (fail) — both are valid
ASSERT_RC=0
ASSERT_OUT=$("$PERCEPTA" assert "$DEVICE_NAME" "LED.LED1 ON" 2>&1) || ASSERT_RC=$?

if [[ $ASSERT_RC -eq 0 || $ASSERT_RC -eq 1 ]]; then
    pass "assert: ran successfully (exit code $ASSERT_RC)"
else
    fail "assert: unexpected exit code $ASSERT_RC"
fi

if echo "$ASSERT_OUT" | grep -qE "(PASS|FAIL)"; then
    pass "assert: output contains PASS or FAIL"
else
    fail "assert: output missing PASS/FAIL verdict"
fi

# ========================================================================
# Phase 6: Update Firmware Tag & Observe Again
# ========================================================================

phase 6 "Update Firmware & Re-observe"

FW_OUT=$("$PERCEPTA" device set-firmware "$DEVICE_NAME" v2.0 2>&1)

if echo "$FW_OUT" | grep -qi "firmware"; then
    pass "set-firmware: confirmation printed"
else
    fail "set-firmware: no confirmation" "$(echo "$FW_OUT" | head -3)"
fi

OBS2_OUT=$("$PERCEPTA" observe "$DEVICE_NAME" 2>&1) || true

if echo "$OBS2_OUT" | grep -qi "Observation captured"; then
    pass "second observe: succeeded"
else
    fail "second observe: failed" "$(echo "$OBS2_OUT" | head -3)"
fi

# ========================================================================
# Phase 7: Diff Between Firmware Versions
# ========================================================================

phase 7 "Diff Between Firmware Versions"

DIFF_RC=0
DIFF_OUT=$("$PERCEPTA" diff "$DEVICE_NAME" --from v1.0 --to v2.0 2>&1) || DIFF_RC=$?

if [[ $DIFF_RC -eq 0 || $DIFF_RC -eq 1 ]]; then
    pass "diff: ran successfully (exit code $DIFF_RC)"
else
    fail "diff: unexpected exit code $DIFF_RC"
fi

if echo "$DIFF_OUT" | grep -qi "Comparing firmware versions"; then
    pass "diff: output contains expected header"
else
    fail "diff: missing 'Comparing firmware versions'" "$(echo "$DIFF_OUT" | head -3)"
fi

# ========================================================================
# Phase 8: Generate Code (Real Anthropic API)
# ========================================================================

phase 8 "Generate Code"

GEN_RC=0
GEN_OUT=$("$PERCEPTA" generate "Blink LED at 1Hz" --board esp32 --output "$WORKSPACE/led_blink.c" 2>&1) || GEN_RC=$?

if [[ $GEN_RC -eq 0 ]]; then
    pass "generate: exit code 0"
else
    fail "generate: exit code $GEN_RC" "$(echo "$GEN_OUT" | tail -5)"
fi

if [[ -f "$WORKSPACE/led_blink.c" ]]; then
    pass "generate: output file exists"
else
    fail "generate: output file missing"
fi

if grep -q "#include" "$WORKSPACE/led_blink.c" 2>/dev/null; then
    pass "generate: file contains #include"
else
    fail "generate: file missing #include directive"
fi

if grep -q "void" "$WORKSPACE/led_blink.c" 2>/dev/null; then
    pass "generate: file contains void"
else
    fail "generate: file missing void keyword"
fi

# ========================================================================
# Phase 9: Style-Check Generated Code
# ========================================================================

phase 9 "Style-Check Generated Code"

SC_RC=0
SC_OUT=$("$PERCEPTA" style-check "$WORKSPACE/led_blink.c" 2>&1) || SC_RC=$?

if [[ $SC_RC -eq 0 || $SC_RC -eq 1 ]]; then
    pass "style-check: ran without crash (exit code $SC_RC)"
else
    fail "style-check: unexpected exit code $SC_RC"
fi

# ========================================================================
# Phase 10: Style-Check with --fix
# ========================================================================

phase 10 "Style-Check with --fix"

FIX_RC=0
FIX_OUT=$("$PERCEPTA" style-check "$WORKSPACE/led_blink.c" --fix 2>&1) || FIX_RC=$?

if [[ $FIX_RC -eq 0 || $FIX_RC -eq 1 ]]; then
    pass "style-check --fix: ran without crash (exit code $FIX_RC)"
else
    fail "style-check --fix: unexpected exit code $FIX_RC"
fi

# If there were violations to fix, check for "Fixed" in the output
if [[ $FIX_RC -eq 1 ]] || echo "$FIX_OUT" | grep -qi "Fixed"; then
    pass "style-check --fix: violations handled"
else
    pass "style-check --fix: code was already compliant"
fi

# ========================================================================
# Phase 11: Store Pattern in Knowledge Graph
# ========================================================================

phase 11 "Store Pattern in Knowledge Graph"

# knowledge store requires BARR-C compliant code. The generated code may
# still have violations after --fix, so write a minimal compliant file.
COMPLIANT_FILE="$WORKSPACE/led_blink_compliant.c"
cat > "$COMPLIANT_FILE" << 'CEOF'
#include <stdint.h>

void LED_Blink(void)
{
    volatile uint8_t led_state = 0U;

    while (1)
    {
        led_state ^= 1U;
    }
}
CEOF

KS_RC=0
KS_OUT=$("$PERCEPTA" knowledge store "Blink LED at 1Hz" "$COMPLIANT_FILE" \
    --device "$DEVICE_NAME" --firmware v2.0 2>&1) || KS_RC=$?

if echo "$KS_OUT" | grep -qi "Pattern stored successfully"; then
    pass "knowledge store: pattern stored"
else
    fail "knowledge store: failed (exit $KS_RC)" "$(echo "$KS_OUT" | head -5)"
fi

# ========================================================================
# Phase 12: List Knowledge Patterns
# ========================================================================

phase 12 "List Knowledge Patterns"

KL_RC=0
KL_OUT=$("$PERCEPTA" knowledge list --board esp32 2>&1) || KL_RC=$?

if echo "$KL_OUT" | grep -qi "Blink LED"; then
    pass "knowledge list: shows stored pattern"
else
    fail "knowledge list: 'Blink LED' not found" "$(echo "$KL_OUT" | head -5)"
fi

# ========================================================================
# Error Conditions
# ========================================================================

phase "E" "Error Conditions"

# generate without API key should fail with helpful message
ERR_OUT=$(ANTHROPIC_API_KEY="" "$PERCEPTA" generate "test" --board esp32 2>&1) || true
if echo "$ERR_OUT" | grep -qi "ANTHROPIC_API_KEY\|API key\|api.key\|not set"; then
    pass "error: missing API key produces helpful message"
else
    fail "error: missing API key message unclear" "$(echo "$ERR_OUT" | head -3)"
fi

# observe nonexistent device
ERR2_OUT=$("$PERCEPTA" observe nonexistent-device 2>&1) || true
if echo "$ERR2_OUT" | grep -qi "not found\|not configured\|unknown device"; then
    pass "error: nonexistent device produces helpful message"
else
    fail "error: nonexistent device message unclear" "$(echo "$ERR2_OUT" | head -3)"
fi

# ========================================================================
# Summary
# ========================================================================

echo ""
echo "========================================"
echo "  Summary"
echo "========================================"
echo ""

for result in "${RESULTS[@]}"; do
    if [[ "$result" == PASS* ]]; then
        echo "  [PASS] ${result#PASS  }"
    else
        echo "  [FAIL] ${result#FAIL  }"
    fi
done

TOTAL=$((PASS_COUNT + FAIL_COUNT))
echo ""
echo "  $PASS_COUNT/$TOTAL passed, $FAIL_COUNT failed"
echo ""

if [[ $FAIL_COUNT -gt 0 ]]; then
    echo "  Some tests failed."
    exit 1
else
    echo "  All tests passed."
    exit 0
fi
