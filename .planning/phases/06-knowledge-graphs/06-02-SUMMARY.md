---
phase: 06-knowledge-graphs
plan: 02
subsystem: knowledge-graphs
tags: [vector-store, embeddings, semantic-search, openai, sqlite, cli, cobra]

# Dependency graph
requires:
  - phase: 06-01-knowledge-graphs
    provides: Graph database, PatternStore API, schema for patterns/observations/styles
  - phase: 05-style-infrastructure
    provides: StyleChecker for BARR-C validation
  - phase: 01-03
    provides: Observation storage (SQLite)
provides:
  - In-memory vector store with SQLite persistence for code embeddings
  - OpenAI embeddings provider (text-embedding-ada-002)
  - Cosine similarity search for pattern matching
  - SearchSimilarPatterns API with board filtering and confidence scoring
  - CLI commands: percepta knowledge store/search/list
  - Semantic search integration with PatternStore
affects: [07-code-generation, 08-public-launch]

# Tech tracking
tech-stack:
  added: [OpenAI embeddings API, in-memory vector store, cosine similarity]
  patterns: [Vector store with SQLite persistence, mock embedder for testing, semantic search with filtering]

key-files:
  created:
    - internal/knowledge/embeddings.go
    - internal/knowledge/vector_store.go
    - internal/knowledge/vector_store_test.go
    - internal/knowledge/semantic_search.go
    - internal/knowledge/semantic_search_test.go
    - cmd/percepta/knowledge.go
  modified:
    - internal/knowledge/pattern_store.go (added vectorStore field, embedding storage)
    - cmd/percepta/main.go (registered knowledge command)

key-decisions:
  - "In-memory vector store + SQLite persistence (pure Go, matches Phase 3 decision)"
  - "OpenAI embeddings API (text-embedding-ada-002) for semantic similarity"
  - "Cosine similarity for pattern matching (simple, effective for MVP)"
  - "Mock embedder for testing (NewVectorStoreWithEmbedder constructor)"
  - "Only return validated patterns (StyleCompliant=true filter in search)"
  - "Confidence scoring = similarity + signal boost (max 10% boost)"
  - "CLI graceful degradation when OPENAI_API_KEY not set"
  - "Board type filtering on both search and list commands"

patterns-established:
  - "Vector store: in-memory with SQLite persistence pattern"
  - "Embedder interface for pluggable embedding providers"
  - "Semantic search: vector similarity + graph filtering + confidence scoring"
  - "CLI commands follow Cobra patterns from device.go, style.go"
  - "Clear, actionable CLI output with code previews and stats"
  - "Automatic embedding storage when patterns added to graph"

issues-created: []

# Metrics
duration: 90min
completed: 2026-02-12
---

# Phase 06-02: Semantic Search and CLI Summary

**In-memory vector store with OpenAI embeddings, semantic pattern search API, and CLI commands for knowledge graph management**

## Performance

- **Duration:** 90 min
- **Started:** 2026-02-12T22:00:00Z
- **Completed:** 2026-02-12T23:30:00Z
- **Tasks:** 3
- **Files created:** 6
- **Files modified:** 2
- **Tests:** 28 total (all passing)

## Accomplishments
- Vector store with code embeddings (in-memory + SQLite persistence)
- OpenAI embeddings provider using text-embedding-ada-002
- Cosine similarity search for semantic pattern matching
- SearchSimilarPatterns API with board filtering and confidence scoring
- CLI commands: `percepta knowledge store/search/list`
- Full integration with PatternStore from Phase 06-01
- Comprehensive test coverage (11 new tests, all passing)

## Task Commits

Each task was committed atomically:

1. **Task 1: Set up vector store for code embeddings** - `5c6bd32` (feat)
2. **Task 2: Add semantic search API for pattern retrieval** - `514b565` (feat)
3. **Task 3: Add CLI commands for pattern management** - `762b7e5` (feat)

## Files Created/Modified

### Created Files
- `internal/knowledge/embeddings.go` - OpenAI embeddings provider, cosine similarity
- `internal/knowledge/vector_store.go` - In-memory vector store with SQLite persistence
- `internal/knowledge/vector_store_test.go` - Vector store tests (6 tests)
- `internal/knowledge/semantic_search.go` - SearchSimilarPatterns API, confidence scoring
- `internal/knowledge/semantic_search_test.go` - Semantic search tests (5 tests)
- `cmd/percepta/knowledge.go` - CLI commands for pattern management

### Modified Files
- `internal/knowledge/pattern_store.go` - Added vectorStore field, automatic embedding storage
- `cmd/percepta/main.go` - Registered knowledge command

## Architecture

### Vector Store Design

**Storage strategy:** In-memory embeddings with SQLite persistence (pure Go, no CGO)

```
~/.local/share/percepta/
├── percepta.db        # Observations (from Phase 1-4)
├── knowledge.db       # Knowledge graph (from Phase 06-01)
└── embeddings.db      # Vector embeddings (NEW)
    └── embeddings table: pattern_id, vector (JSON), created_at
```

**Why in-memory + SQLite:**
- Pure Go architecture (matches Phase 3 modernc.org/sqlite decision)
- Fast in-memory vector search (cosine similarity)
- Persistent across restarts
- No external services or CGO dependencies
- Can upgrade to Qdrant/Milvus later if needed (Phase 8+)

### Embeddings Provider

**OpenAI API integration:**
```go
type OpenAIEmbeddings struct {
    apiKey string
    model  string  // "text-embedding-ada-002"
}

func (o *OpenAIEmbeddings) Embed(text string) ([]float32, error)
```

**Key features:**
- API key from `OPENAI_API_KEY` environment variable
- Standard OpenAI embeddings API endpoint
- Returns 1536-dimensional float32 vectors
- Error handling for API failures

**Cost:** ~$0.0001 per 1K tokens (acceptable for MVP)

### Cosine Similarity

**Algorithm:**
```go
func cosineSimilarity(a, b []float32) float32 {
    dotProduct := sum(a[i] * b[i])
    normA := sqrt(sum(a[i] * a[i]))
    normB := sqrt(sum(b[i] * b[i]))
    return dotProduct / (normA * normB)
}
```

**Returns:** Value between -1 (opposite) and 1 (identical)

**Performance:** O(n*d) where n=patterns, d=dimensions (1536)
- Acceptable for <10k patterns (in-memory)
- Can optimize with HNSW/IVF if needed later

### Semantic Search API

**SearchSimilarPatterns pipeline:**
```
1. Generate query embedding via OpenAI API
2. Compute cosine similarity with all stored embeddings
3. Sort by similarity (highest first)
4. Filter by board type (if specified)
5. Filter by StyleCompliant=true (safety check)
6. Calculate confidence scores
7. Return top K results
```

**Confidence scoring:**
```go
confidence = similarity + signalBoost
where signalBoost = min(0.1, signalCount / 100)
```

**Rationale:** High similarity + more validation signals = higher confidence

### CLI Commands

**percepta knowledge store:**
```bash
percepta knowledge store "Blink LED at 1Hz" led.c --device esp32-dev --firmware v1.0.0
```

**Output:**
```
✓ Pattern stored successfully
  ID:       237f462af0137e60...
  Spec:     Blink LED at 1Hz
  Device:   esp32-dev
  Firmware: v1.0.0
  File:     led.c

Knowledge graph stats:
  Patterns:      1
  Observations:  1
```

**percepta knowledge search:**
```bash
percepta knowledge search "blink LED" --board esp32 --limit 5
```

**Output:**
```
Found 3 similar pattern(s):

1. Blink LED at 1Hz
   Board:      esp32
   Similarity: 95%
   Confidence: 97%
   Style:      BARR-C compliant
   Code:
     #include <stdint.h>

     void LED_Blink(void) {
     ... (4 more lines)
   Signals:    2 observed
```

**percepta knowledge list:**
```bash
percepta knowledge list --board esp32
```

**Output:**
```
Validated patterns (5 total):

1. Blink LED at 1Hz
   Board:    esp32
   Style:    BARR-C compliant
   Created:  2026-02-12 23:32:23
   Code:
     #include <stdint.h>

     ... (4 more lines)
```

## Usage Examples

### Store a validated pattern

```bash
# 1. First, ensure observation exists for device+firmware
percepta observe esp32-dev

# 2. Store pattern (code must be BARR-C compliant)
percepta knowledge store "Blink LED at 1Hz" led.c \
  --device esp32-dev \
  --firmware v1.0.0
```

### Search for similar patterns

```bash
# Semantic search (uses OpenAI embeddings)
export OPENAI_API_KEY=sk-...
percepta knowledge search "toggle LED" --board esp32 --limit 3

# Results ranked by similarity + confidence
```

### List all patterns

```bash
# All patterns
percepta knowledge list

# Filter by board type
percepta knowledge list --board esp32
```

## Test Coverage

**28 tests total, all passing:**

### Vector Store Tests (6 tests)
- `TestCosineSimilarity` - Similarity calculation correctness
- `TestVectorStore_StoreAndRetrieve` - Basic storage and retrieval
- `TestVectorStore_FindSimilar` - Semantic search ranking
- `TestVectorStore_Persistence` - SQLite persistence across restarts
- `TestVectorStore_ReplaceEmbedding` - Update existing embeddings
- `TestVectorStore_EmptyQuery` - Handle empty store gracefully

### Semantic Search Tests (5 tests)
- `TestSearchSimilarPatterns_Success` - End-to-end search flow
- `TestSearchSimilarPatterns_BoardFilter` - Board type filtering
- `TestSearchSimilarPatterns_TopK` - Result limit enforcement
- `TestSearchSimilarPatterns_NoVectorStore` - Graceful failure without vector store
- `TestCalculateConfidence` - Confidence scoring correctness

### Integration with Phase 06-01 (17 tests from previous phase)
- All graph and PatternStore tests still passing
- Vector store seamlessly integrated

## Decisions Made

1. **In-memory + SQLite over external vector DB:** Pure Go, zero dependencies, matches Phase 3 decision. Can upgrade to Qdrant/Milvus if performance needed (>10k patterns).

2. **OpenAI embeddings over local models:** Ship faster with proven accuracy. text-embedding-ada-002 is industry standard. Pluggable architecture allows adding moondream/BERT later if cost becomes issue.

3. **Cosine similarity over more complex metrics:** Simple, effective, well-understood. Sufficient for MVP. Can add dot product, L2 distance if needed.

4. **Mock embedder for testing:** NewVectorStoreWithEmbedder constructor allows testing without API keys. Deterministic embeddings based on text content enable reliable tests.

5. **Confidence = similarity + signal boost:** Combines vector similarity with validation metadata. More signals = more validation = higher confidence. Capped at 10% boost to avoid over-weighting.

6. **CLI graceful degradation:** If OPENAI_API_KEY not set, pattern storage still works (graph only), semantic search fails gracefully with clear error message. User can still use exact-match queries.

7. **Board type filtering in both API and CLI:** Common use case ("show me ESP32 patterns"). Implemented at API level for reusability.

8. **Code preview in output:** Show first 2-3 lines of code for quick scanning. Full code available via graph query if needed.

## Deviations from Plan

None - plan executed exactly as written.

All tasks completed successfully:
- ✅ Task 1: Vector store with embeddings (6 tests passing)
- ✅ Task 2: Semantic search API (5 tests passing)
- ✅ Task 3: CLI commands (manual testing successful)

No bugs encountered, no scope creep, no architectural changes needed.

## Issues Encountered

None - all components integrated smoothly.

**Smooth integration points:**
- PatternStore already had clean API (from 06-01)
- StyleChecker integration already working (from 05-02)
- Observation storage already functional (from 03-01)
- Cobra CLI patterns already established (from 04-01, 05-02)

## Integration Notes for Phase 07

**Ready for code generation:**
- Semantic search can find similar validated patterns
- Confidence scores help prioritize high-quality examples
- Board type filtering ensures relevant patterns
- CLI allows manual inspection of knowledge graph

**How Phase 07 will use this:**
```go
// Code generation workflow
func generateFirmware(spec string, board string) (string, error) {
    // 1. Find similar validated patterns
    patterns := store.SearchSimilarPatterns(spec, board, 5)

    // 2. Extract code patterns for LLM context
    examples := extractCodePatterns(patterns)

    // 3. Generate with LLM (using examples as context)
    code := llm.Generate(spec, examples)

    // 4. Validate style + behavior
    // 5. Store as new pattern
}
```

**Current capabilities available:**
- `SearchSimilarPatterns(query, board, topK)` - Semantic search
- `QueryPatternsByBoard(board)` - Exact board match
- `GetPatternWithObservation(id)` - Full pattern + observation
- `Stats()` - Knowledge graph statistics

## Performance Notes

**Vector search speed:**
- Current: O(n) linear scan (acceptable for <10k patterns)
- Tested with 10 patterns: <20ms
- Estimated 1000 patterns: <200ms
- Estimated 10k patterns: ~2s

**Optimization strategies if needed (Phase 8+):**
- HNSW index for approximate nearest neighbor
- Quantization (float32 → int8) for 4x memory reduction
- Migrate to Qdrant/Milvus for distributed search
- Batch embedding generation (reduce API calls)

**Storage size:**
- Each embedding: 1536 floats × 4 bytes = 6KB
- 1000 patterns: ~6MB
- 10k patterns: ~60MB
- Acceptable for in-memory MVP

## Next Phase Readiness

**Phase 07 (Code Generation) is ready to start:**
- ✅ Semantic search operational
- ✅ Validated pattern storage working
- ✅ CLI commands for manual testing
- ✅ Confidence scoring functional
- ✅ Board type filtering ready

**No blockers for Phase 07.**

**Phase 6 complete - both plans done:**
- Plan 06-01: Knowledge Graph Storage ✅
- Plan 06-02: Semantic Search + CLI ✅

**Total Phase 6 stats:**
- Duration: 150 min (06-01: 60min, 06-02: 90min)
- Files created: 11
- Tests: 28 (all passing)
- Ready for AI code generation with hardware validation

---
*Phase: 06-knowledge-graphs*
*Completed: 2026-02-12*
