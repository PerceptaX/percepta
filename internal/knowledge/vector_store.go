package knowledge

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	_ "modernc.org/sqlite"
)

// PatternMatch represents a pattern match with similarity score
type PatternMatch struct {
	PatternID  string
	Similarity float32
}

// VectorStore provides semantic search over code patterns using embeddings
type VectorStore struct {
	embeddings map[string][]float32 // patternID -> embedding vector
	embedder   EmbeddingProvider    // Embedding generation
	db         *sql.DB              // SQLite for persistence
	mu         sync.RWMutex         // Concurrent access control
}

// NewVectorStore creates a new vector store with SQLite persistence
func NewVectorStore() (*VectorStore, error) {
	// Create embeddings provider
	embedder, err := NewOpenAIEmbeddings()
	if err != nil {
		return nil, fmt.Errorf("failed to create embedder: %w", err)
	}

	return NewVectorStoreWithEmbedder(embedder)
}

// NewVectorStoreWithEmbedder creates a vector store with a custom embedder
// Useful for testing with mock embedders
func NewVectorStoreWithEmbedder(embedder EmbeddingProvider) (*VectorStore, error) {
	// Get database path
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home dir: %w", err)
	}

	dbDir := filepath.Join(home, ".local", "share", "percepta")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	dbPath := filepath.Join(dbDir, "embeddings.db")

	// Open database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create table if not exists
	if err := createEmbeddingsTable(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	// Load existing embeddings into memory
	embeddings, err := loadEmbeddings(db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to load embeddings: %w", err)
	}

	return &VectorStore{
		embeddings: embeddings,
		embedder:   embedder,
		db:         db,
	}, nil
}

// createEmbeddingsTable creates the embeddings table if it doesn't exist
func createEmbeddingsTable(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS embeddings (
		pattern_id TEXT PRIMARY KEY,
		vector TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_embeddings_created ON embeddings(created_at);
	`

	_, err := db.Exec(schema)
	return err
}

// loadEmbeddings loads all embeddings from database into memory
func loadEmbeddings(db *sql.DB) (map[string][]float32, error) {
	embeddings := make(map[string][]float32)

	rows, err := db.Query("SELECT pattern_id, vector FROM embeddings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var patternID, vectorJSON string
		if err := rows.Scan(&patternID, &vectorJSON); err != nil {
			return nil, err
		}

		var vector []float32
		if err := json.Unmarshal([]byte(vectorJSON), &vector); err != nil {
			return nil, fmt.Errorf("failed to unmarshal vector for %s: %w", patternID, err)
		}

		embeddings[patternID] = vector
	}

	return embeddings, rows.Err()
}

// StoreEmbedding stores an embedding for a pattern
// If the pattern already has an embedding, it is replaced
func (v *VectorStore) StoreEmbedding(patternID string, code string) error {
	// Generate embedding
	embedding, err := v.embedder.Embed(code)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Serialize embedding
	vectorJSON, err := json.Marshal(embedding)
	if err != nil {
		return fmt.Errorf("failed to marshal embedding: %w", err)
	}

	// Store in database (replace if exists)
	v.mu.Lock()
	defer v.mu.Unlock()

	_, err = v.db.Exec(
		"INSERT OR REPLACE INTO embeddings (pattern_id, vector) VALUES (?, ?)",
		patternID, string(vectorJSON),
	)
	if err != nil {
		return fmt.Errorf("failed to store embedding: %w", err)
	}

	// Update in-memory cache
	v.embeddings[patternID] = embedding

	return nil
}

// FindSimilar finds the most similar patterns to the given query code
// Returns top K matches sorted by similarity (highest first)
func (v *VectorStore) FindSimilar(queryCode string, topK int) ([]PatternMatch, error) {
	// Generate query embedding
	queryVec, err := v.embedder.Embed(queryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	v.mu.RLock()
	defer v.mu.RUnlock()

	// Compute similarity with all stored embeddings
	var matches []PatternMatch
	for patternID, embedding := range v.embeddings {
		similarity := cosineSimilarity(queryVec, embedding)
		matches = append(matches, PatternMatch{
			PatternID:  patternID,
			Similarity: similarity,
		})
	}

	// Sort by similarity (highest first)
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Similarity > matches[j].Similarity
	})

	// Return top K
	if len(matches) > topK {
		matches = matches[:topK]
	}

	return matches, nil
}

// HasEmbedding checks if a pattern has an embedding
func (v *VectorStore) HasEmbedding(patternID string) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	_, exists := v.embeddings[patternID]
	return exists
}

// Count returns the number of stored embeddings
func (v *VectorStore) Count() int {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return len(v.embeddings)
}

// Close closes the database connection
func (v *VectorStore) Close() error {
	return v.db.Close()
}
