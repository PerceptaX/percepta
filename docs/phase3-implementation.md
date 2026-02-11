# Phase 3 Implementation: Diff + Firmware Tracking

## Overview

Phase 3 adds firmware version tracking and diff capabilities to Percepta, enabling users to compare hardware behavior across different firmware versions.

## Architecture

### Storage Layer

**SQLite Storage** (`internal/storage/sqlite.go`)
- Pure Go implementation using `modernc.org/sqlite` (no CGO dependencies)
- Database location: `~/.local/share/percepta/percepta.db`
- Schema:
  ```sql
  CREATE TABLE observations (
      id TEXT PRIMARY KEY,
      device_id TEXT NOT NULL,
      firmware TEXT NOT NULL DEFAULT '',
      timestamp DATETIME NOT NULL,
      signals_json TEXT NOT NULL
  );

  CREATE INDEX idx_device_firmware
  ON observations(device_id, firmware, timestamp);
  ```

**Key Methods:**
- `Save(obs)` - Persists observations with firmware tag
- `Query(deviceID, limit)` - Retrieves observations by device
- `QueryByFirmware(deviceID, firmware, limit)` - Retrieves observations for specific firmware
- `GetLatestForFirmware(deviceID, firmware)` - Gets most recent observation for firmware version

### Manual Firmware Tagging

**Configuration** (`internal/config/config.go`)
- Added `Firmware` field to `DeviceConfig`
- User-specified string tags (e.g., "v1", "abc123", "main")
- No git integration - manual tagging only

**Example config** (`~/.config/percepta/config.yaml`):
```yaml
devices:
  fpga:
    type: fpga
    camera_id: /dev/video0
    firmware: v1
```

**Workflow:**
1. User sets firmware tag in config
2. Observations are captured with that tag
3. User changes tag when firmware changes
4. New observations use new tag

### Diff Logic

**Comparison Engine** (`internal/diff/compare.go`)
- Exact comparison of signals (no tolerances)
- Detects: ADDED, REMOVED, MODIFIED signals
- Comparison rules:
  - State (ON/OFF): Exact match
  - Color (RGB): Exact match
  - Text: Exact match
  - BlinkHz: **Normalized to 1 decimal** (handles Claude Vision fluctuations)
  - Confidence: Ignored completely

**BlinkHz Normalization:**
```go
// 2.04 → 2.0
// 2.05 → 2.1
// Handles Claude Vision's slight variations (2.0 vs 2.1)
func normalizeBlinkHz(hz float64) float64 {
    if hz == 0 {
        return 0
    }
    return math.Round(hz*10) / 10
}
```

**Diff Types** (`internal/diff/types.go`)
- `ChangeAdded`: Signal exists in 'to' but not 'from'
- `ChangeRemoved`: Signal exists in 'from' but not 'to'
- `ChangeModified`: Signal exists in both but with different state

### CLI Command

**Usage:**
```bash
percepta diff <device> --from <firmware1> --to <firmware2>
```

**Example:**
```bash
percepta diff fpga --from v1 --to v2
```

**Output Format:**
```
Comparing firmware versions:
FROM: v1 (2026-02-11 10:30:00)
TO:   v2 (2026-02-11 11:45:00)

Device: fpga

Changes detected:

+ LED2: purple blinking 0.8Hz (ADDED)
- LED3: red solid (REMOVED)
~ LED1: blue blinking 2.0Hz → 2.5Hz (MODIFIED)

Summary: 1 added, 1 removed, 1 modified
```

**Exit Codes:**
- `0`: No changes (firmware behavior identical)
- `1`: Changes detected
- `2`: Error occurred

## Key Design Decisions

### 1. Manual Firmware Tags (NOT git auto-integration)

**Rationale:**
- Git coupling breaks FPGA workflows (binary formats, non-repo users)
- CI/CD environments may have detached HEAD or shallow clones
- Users want explicit control over firmware versioning
- Simpler architecture without `os/exec` dependencies

**Benefits:**
- Works with any workflow (git, binaries, manual updates)
- No assumptions about repository structure
- Deterministic behavior

### 2. Pure Go SQLite (modernc.org/sqlite)

**Rationale:**
- Maintains zero CGO dependencies
- Cross-platform without C compiler
- Consistent with project's portability goals

**Verification:**
```bash
CGO_ENABLED=0 go build -o percepta ./cmd/percepta
# Build succeeds without CGO
```

### 3. Exact Diff (NO tolerances)

**Rationale:**
- Assertions already handle fuzz/tolerances
- Diff must be deterministic and explicit
- Users expect "what changed" not "what might have changed"

**Exception:**
- BlinkHz normalized to 1 decimal (handles Claude Vision's natural variation)
- This is **normalization**, not tolerance - values are rounded before comparison

### 4. Storage Construction in CMD Layer

**Rationale:**
- `pkg/percepta` stays framework-agnostic (accepts `StorageDriver` interface)
- `cmd/percepta` handles concrete implementations (SQLite, firmware injection)
- Firmware tag applied in cmd layer before storage

**Flow:**
```
cmd/percepta/observe.go:
  1. Load config → get firmware tag
  2. Initialize SQLiteStorage
  3. Initialize Core with storage interface
  4. Capture observation
  5. Inject firmware tag
  6. Save to storage

pkg/percepta/percepta.go:
  - Accepts StorageDriver interface
  - No knowledge of SQLite or firmware tags
  - Framework-agnostic
```

## Testing

### Unit Tests

**Diff Logic** (`internal/diff/compare_test.go`):
- ✓ No changes detection
- ✓ Added/Removed/Modified signals
- ✓ BlinkHz normalization (2.04 → 2.0)
- ✓ Exact color comparison (255,0,0 ≠ 254,0,0)
- ✓ Confidence ignored
- ✓ Display text changes
- ✓ CountByType accuracy

**SQLite Storage** (`internal/storage/sqlite_test.go`):
- ✓ Save and query
- ✓ Query by firmware
- ✓ Get latest for firmware
- ✓ Empty firmware tag handling
- ✓ Signal serialization/deserialization
- ✓ Database path creation

**Run Tests:**
```bash
go test ./internal/diff/
go test ./internal/storage/
go test ./...
```

### Integration Test

**Automated Checks** (`test_phase3.sh`):
```bash
./test_phase3.sh
```

**Checks:**
- Database initialization
- Schema correctness
- No CGO dependencies
- Pure Go build (CGO_ENABLED=0)

### Manual Testing

```bash
# 1. Clean slate
rm -f ~/.local/share/percepta/percepta.db

# 2. Set firmware v1 in config
cat > ~/.config/percepta/config.yaml <<EOF
devices:
  fpga:
    firmware: v1
EOF

# 3. Capture observation with v1
./percepta observe fpga

# 4. Change to v2
sed -i 's/firmware: v1/firmware: v2/' ~/.config/percepta/config.yaml

# 5. Capture observation with v2
./percepta observe fpga

# 6. Compare
./percepta diff fpga --from v1 --to v2

# 7. Verify database
sqlite3 ~/.local/share/percepta/percepta.db \
  "SELECT device_id, firmware, timestamp FROM observations;"
```

## Files Created

```
internal/storage/sqlite.go          - SQLite storage implementation
internal/storage/sqlite_test.go     - SQLite unit tests
internal/diff/types.go              - Diff type definitions
internal/diff/compare.go            - Signal comparison logic
internal/diff/compare_test.go       - Comparison unit tests
cmd/percepta/diff.go                - Diff CLI command
test_phase3.sh                      - Verification script
docs/phase3-implementation.md       - This file
```

## Files Modified

```
internal/config/config.go           - Added Device.Firmware field
internal/core/interfaces.go         - Added StorageDriver interface
pkg/percepta/percepta.go            - Accept storage via interface
cmd/percepta/observe.go             - Use SQLite, inject firmware
cmd/percepta/assert.go              - Use SQLite, inject firmware
cmd/percepta/main.go                - Register diff command
go.mod                              - Added modernc.org/sqlite
```

## Dependencies Added

```
modernc.org/sqlite v1.45.0          - Pure Go SQLite driver
└── Dependencies (all pure Go):
    ├── github.com/dustin/go-humanize
    ├── github.com/google/uuid
    ├── github.com/mattn/go-isatty
    ├── github.com/ncruces/go-strftime
    ├── github.com/remyoudompheng/bigfft
    ├── golang.org/x/exp
    ├── modernc.org/libc
    ├── modernc.org/mathutil
    └── modernc.org/memory
```

**Zero CGO dependencies maintained ✓**

## API Stability

**Public API (unchanged):**
- `percepta observe <device>` - Still works, now persists to SQLite
- `percepta assert <device> <assertion>` - Still works with SQLite

**New API:**
- `percepta diff <device> --from <fw1> --to <fw2>` - Compare firmware versions

## Limitations

1. **Single observation comparison**: Diff compares latest observation for each firmware version, not historical trends
2. **No firmware auto-detection**: User must manually update config when firmware changes
3. **No migration from MemoryStorage**: Existing in-memory observations are not migrated (fresh start)

## Future Enhancements (Out of Scope for Phase 3)

- Firmware auto-detection from git/build metadata
- Historical trend analysis (multiple observations per firmware)
- Migration tools for existing data
- Web UI for diff visualization
- Tolerance configuration per signal type

## Success Criteria ✓

- [x] SQLite storage persists observations across runs
- [x] Firmware tags are user-specified strings (not git hashes)
- [x] Diff shows exact differences (no tolerance fuzz)
- [x] No CGO dependencies (pure Go build works)
- [x] Exit codes: 0 if identical, 1 if changes, 2 if error
- [x] Human-readable output with +/-/~ indicators
- [x] Database at `~/.local/share/percepta/percepta.db`
- [x] All tests pass

## Verification Results

```bash
$ ./test_phase3.sh
=== Phase 3 Verification Test ===

1. Cleaning database...
   ✓ Database cleaned

2. Testing database initialization...
   ✓ Database created and schema initialized

3. Verifying schema...
   ✓ observations table created
   ✓ firmware column exists
   ✓ index created

4. Checking for CGO dependencies...
   ✓ Using modernc.org/sqlite (pure Go)

5. Testing build without CGO...
   ✓ Build successful without CGO

=== All Phase 3 checks passed! ===
```

```bash
$ go test ./internal/diff/ ./internal/storage/
ok  	github.com/perceptumx/percepta/internal/diff	0.002s
ok  	github.com/perceptumx/percepta/internal/storage	0.010s
```

## Summary

Phase 3 successfully implements:
1. **SQLite storage** - Persistent observations with pure Go driver
2. **Manual firmware tagging** - User-specified firmware versions in config
3. **Exact signal diff** - Deterministic comparison with normalized BlinkHz
4. **CLI diff command** - Human-readable firmware comparison

**No git, no tolerances, no CGO, no cleverness. Simple, explicit, deterministic.**
