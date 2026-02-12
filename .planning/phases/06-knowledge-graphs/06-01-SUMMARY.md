---
phase: 06-knowledge-graphs
plan: 01
subsystem: knowledge-graphs
tags: [graph-database, pattern-storage, barr-c, sqlite, pure-go, knowledge-graph]

# Dependency graph
requires:
  - phase: 05-style-infrastructure
    provides: StyleChecker for BARR-C validation, tree-sitter parser
  - phase: 01-03
    provides: Observation storage (SQLite)
provides:
  - Graph database with in-memory + SQLite persistence
  - PatternNode, ObservationNode, StyleResultNode, Edge schema
  - Relationship types: IMPLEMENTED_BY, RUNS_ON, PRODUCES, VALIDATED_BY
  - PatternStore API for validated pattern storage
  - Query APIs: QueryPatternsByBoard, QueryPatternsBySpec
affects: [06-02-semantic-search, 07-code-generation]

# Tech tracking
tech-stack:
  added: [knowledge graph (in-memory), modernc.org/sqlite (knowledge.db)]
  patterns: [Graph pattern with persistence, PatternStore integration layer, relationship edges]

key-files:
  created:
    - internal/knowledge/schema.go
    - internal/knowledge/graph.go
    - internal/knowledge/graph_test.go
    - internal/knowledge/pattern_store.go
    - internal/knowledge/pattern_store_test.go

key-decisions:
  - "In-memory graph with SQLite persistence (pure Go, matches Phase 3 modernc.org/sqlite decision)"
  - "Store only validated patterns (StyleCompliant=true AND has observation)"
  - "Full relationship graph: spec->pattern->board->observation->style_result"
  - "Database path: ~/.local/share/percepta/knowledge.db (alongside percepta.db)"
  - "AddPattern returns generated ID for caller convenience"
  - "PatternStore integrates StyleChecker, Graph, and SQLite storage"
  - "Reject patterns without observation (hardware validation required)"

patterns-established:
  - "Graph database: in-memory nodes with SQLite persistence"
  - "Node types: PatternNode, ObservationNode, StyleResultNode"
  - "Edge types: IMPLEMENTED_BY, RUNS_ON, PRODUCES, VALIDATED_BY, SIMILAR_TO (ready for 06-02)"
  - "PatternStore: validation -> storage -> relationship creation pipeline"
  - "Query patterns: by board type (exact match), by spec (exact match for now)"
  - "getBoardType heuristic: extract board family from device ID"

issues-created: []

# Metrics
duration: 60min
completed: 2026-02-13
---

# Phase 06-01: Knowledge Graph Storage Summary

**In-memory graph database with SQLite persistence + PatternStore API for validated firmware patterns**

## Performance

- **Duration:** 60 min
- **Started:** 2026-02-13T01:00:00Z
- **Completed:** 2026-02-13T02:00:00Z
- **Tasks:** 2
- **Files created:** 5
- **Tests:** 16 (all passing)

## Accomplishments
- Graph database with in-memory representation and SQLite persistence
- Schema: PatternNode (code, spec, board), ObservationNode, StyleResultNode, Edge
- Relationship types: IMPLEMENTED_BY, RUNS_ON, PRODUCES, VALIDATED_BY
- PatternStore API integrates StyleChecker + Graph + SQLite storage
- Only stores validated patterns (BARR-C compliant + has observation)
- Query APIs: QueryPatternsByBoard, QueryPatternsBySpec
- Full persistence to ~/.local/share/percepta/knowledge.db

## Task Commits

Each task was committed atomically:

1. **Task 1: Set up graph database with schema** - `5dcb806` (feat)
2. **Task 2: Create PatternStore API for validated code** - `3616b4c` (feat)

## Files Created/Modified
- `internal/knowledge/schema.go` - Node and edge type definitions
- `internal/knowledge/graph.go` - Graph database with SQLite persistence
- `internal/knowledge/graph_test.go` - Graph tests (9 tests)
- `internal/knowledge/pattern_store.go` - PatternStore API with validation
- `internal/knowledge/pattern_store_test.go` - PatternStore tests (7 tests)

## Architecture

### Graph Database Design

**Storage strategy:** In-memory graph with SQLite persistence (pure Go, no CGO)

```
~/.local/share/percepta/
├── percepta.db        # Observations (from Phase 1-4)
└── knowledge.db       # Knowledge graph (NEW)
    ├── patterns       # Validated code patterns
    ├── observations   # Observation references
    ├── style_results  # Style check results
    └── edges          # Relationships between nodes
```

**Why in-memory + SQLite:**
- Pure Go architecture (matches Phase 3 modernc.org/sqlite decision)
- Fast in-memory queries
- Persistent across restarts
- No external services or CGO dependencies
- Can upgrade to dgraph later if needed (Phase 06-02+)

### Schema

**Node Types:**

```go
PatternNode {
    ID              string    // SHA256 hash of code+spec+board
    Code            string    // C source code (BARR-C compliant)
    Spec            string    // Natural language specification
    BoardType       string    // "esp32", "stm32", etc.
    StyleCompliant  bool      // Always true (validated before storage)
    CreatedAt       time.Time
}

ObservationNode {
    ID        string    // From perception storage
    DeviceID  string    // Device that produced observation
    Firmware  string    // Firmware hash
    Timestamp time.Time
}

StyleResultNode {
    ID        string
    Compliant bool      // Overall compliance
    AutoFixed bool      // Whether violations were auto-fixed
    ViolCount int       // Number of violations
}
```

**Relationship Types:**

```
IMPLEMENTED_BY:  Spec -> Pattern       (spec "Blink LED" -> code)
RUNS_ON:         Pattern -> Board      (code -> "esp32")
PRODUCES:        Pattern -> Observation (code -> hardware observation)
VALIDATED_BY:    Pattern -> StyleResult (code -> BARR-C check)
SIMILAR_TO:      Pattern -> Pattern    (for semantic search in 06-02)
```

**Graph structure:**

```
┌─────────┐
│  Spec   │
└────┬────┘
     │ IMPLEMENTED_BY
     ▼
┌─────────┐  RUNS_ON     ┌───────┐
│ Pattern ├─────────────>│ Board │
└────┬────┘              └───────┘
     │ PRODUCES
     ├──────────────────> Observation
     │ VALIDATED_BY
     └──────────────────> StyleResult
```

### PatternStore API

**Integration layer:** Combines StyleChecker + Graph + Observation Storage

```go
type PatternStore struct {
    graph      *Graph                  // Knowledge graph
    styleCheck *style.StyleChecker     // BARR-C validator
    storage    *storage.SQLiteStorage  // Observation storage
}
```

**Validation pipeline:**

```
StoreValidatedPattern(spec, code, deviceID, firmware)
    │
    ├─> 1. StyleChecker.CheckSource(code) ───> violations?
    │      └─> REJECT if not BARR-C compliant
    │
    ├─> 2. storage.GetLatestForFirmware(deviceID, firmware)
    │      └─> REJECT if no observation found
    │
    ├─> 3. graph.AddPattern(pattern)
    │
    ├─> 4. graph.AddObservation(obsNode)
    │
    ├─> 5. graph.AddStyleResult(styleResult)
    │
    └─> 6. Create relationships:
           - Spec -> Pattern (IMPLEMENTED_BY)
           - Pattern -> Board (RUNS_ON)
           - Pattern -> Observation (PRODUCES)
           - Pattern -> StyleResult (VALIDATED_BY)
```

**Only validated patterns are stored:**
- BARR-C compliant (no style violations)
- Has corresponding hardware observation
- Full relationship graph created

## Usage Examples

### Store a validated pattern

```go
store, _ := NewPatternStore()

// This code is BARR-C compliant
code := `#include <stdint.h>

void LED_Blink(void) {
    uint8_t status = 1;
}
`

// Must have observation for this device+firmware
patternID, err := store.StoreValidatedPattern(
    "Blink LED at 1Hz",              // spec
    code,                             // code
    "esp32-devkit-v1",               // device
    "v1.0.0",                        // firmware
)

// Success! Pattern stored with full relationship graph
```

### Query patterns by board type

```go
// Find all ESP32 patterns
patterns, _ := store.QueryPatternsByBoard("esp32")

for _, pattern := range patterns {
    fmt.Printf("Pattern: %s\n", pattern.Spec)
    fmt.Printf("Code:\n%s\n", pattern.Code)
}
```

### Query patterns by specification

```go
// Find patterns implementing specific spec
patterns, _ := store.QueryPatternsBySpec("Blink LED at 1Hz")

// Returns all patterns that exactly match spec
// (Semantic search coming in Phase 06-02)
```

### Get pattern with observation

```go
pattern, obs, _ := store.GetPatternWithObservation(patternID)

fmt.Printf("Pattern code: %s\n", pattern.Code)
fmt.Printf("Produced observation: %d signals\n", len(obs.Signals))
fmt.Printf("Observation ID: %s\n", obs.ID)
```

## Test Coverage

**16 tests total, all passing:**

### Graph Tests (9 tests)
- `TestGraph_AddPattern` - Store pattern node
- `TestGraph_AddObservation` - Store observation node
- `TestGraph_AddStyleResult` - Store style result node
- `TestGraph_AddEdge` - Create relationships
- `TestGraph_QueryPatternsByBoard` - Query by board type
- `TestGraph_QueryPatternsBySpec` - Query by spec
- `TestGraph_Persistence` - SQLite persistence across restarts
- `TestGraph_Stats` - Graph statistics
- `TestGraph_MultipleEdgesFromNode` - Multiple relationships

### PatternStore Tests (7 tests)
- `TestPatternStore_StoreValidatedPattern_Success` - Happy path
- `TestPatternStore_StoreValidatedPattern_StyleViolations` - Reject non-compliant
- `TestPatternStore_StoreValidatedPattern_NoObservation` - Reject without observation
- `TestPatternStore_QueryPatternsByBoard` - Board type queries
- `TestPatternStore_QueryPatternsBySpec` - Spec queries
- `TestPatternStore_GetPatternWithObservation` - Full retrieval
- `TestPatternStore_Stats` - Statistics
- `TestGetBoardType` - Board type extraction

## Decisions Made

1. **In-memory + SQLite over external graph DB:** Pure Go, zero dependencies, matches Phase 3 decision
2. **Store only validated patterns:** Ensures quality - no speculation, only hardware-verified code
3. **Observation required:** Pattern must have corresponding observation (hardware validation)
4. **Full relationship graph:** All connections created atomically in StoreValidatedPattern
5. **AddPattern returns ID:** Caller convenience for creating relationships
6. **Board type extraction heuristic:** Simple prefix matching (esp32-*, stm32-*) - can be enhanced later
7. **Exact match queries for now:** Semantic search deferred to Phase 06-02

## Deviations from Plan

**None.** Plan executed exactly as written.

Both tasks completed successfully:
- ✅ Task 1: Graph database with schema
- ✅ Task 2: PatternStore API for validated code

No bugs encountered, no scope creep, no architectural changes needed.

## Integration Notes for Phase 06-02

**Ready for semantic search:**
- SIMILAR_TO relationship type already defined
- Graph schema supports vector embeddings (can add metadata field)
- Query interface extensible (can add semantic search alongside exact match)

**Next phase will add:**
- Vector store (Qdrant or in-memory)
- Code embeddings (semantic similarity)
- `QuerySimilarPatterns(spec, board, topK)`
- CLI: `percepta knowledge search "Blink LED"`

**Current capabilities available:**
- `PatternStore.StoreValidatedPattern()` - Store after validation
- `PatternStore.QueryPatternsByBoard()` - Filter by board type
- `PatternStore.QueryPatternsBySpec()` - Exact spec match
- `PatternStore.GetPatternWithObservation()` - Full pattern + observation
- `PatternStore.Stats()` - Graph statistics

## Graph Statistics (after tests)

Example stats from test runs:
```
{
    "patterns": 5,
    "observations": 5,
    "style_results": 5,
    "edges": 20  // 4 edges per pattern (IMPLEMENTED_BY, RUNS_ON, PRODUCES, VALIDATED_BY)
}
```

## Future Enhancements (not in scope)

**Phase 06-02 will add:**
- Semantic search with vector embeddings
- SIMILAR_TO relationships between patterns
- `QuerySimilarPatterns()` API

**Phase 07 (code generation) will use:**
- Pattern retrieval by similarity
- Context injection from validated patterns
- Board-specific quirks from observation metadata

**Post-v2.0 enhancements:**
- Graph visualization
- Pattern clustering
- Automatic quirk detection from failure patterns
- Community pattern sharing

## Blockers/Concerns

**None.** All integration points working:
- ✅ StyleChecker integration (Phase 5)
- ✅ SQLite storage integration (Phase 3)
- ✅ Observation retrieval (Phase 1)
- ✅ All tests passing

## Next Phase Readiness

**Phase 06-02 (Semantic Search) is ready to start:**
- Graph database operational
- PatternStore API working
- Schema supports extensions (metadata field for embeddings)
- SIMILAR_TO relationship type defined

**No dependencies blocking next phase.**

---
*Phase: 06-knowledge-graphs*
*Completed: 2026-02-13*
