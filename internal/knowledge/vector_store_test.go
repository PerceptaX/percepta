//go:build !windows

package knowledge

import (
	"os"
	"path/filepath"
	"testing"
)

// MockEmbedder provides deterministic embeddings for testing
type MockEmbedder struct{}

func (m *MockEmbedder) Embed(text string) ([]float32, error) {
	// Generate simple embeddings based on text length and content
	// This is deterministic and allows testing similarity ranking
	vec := make([]float32, 128)

	// Use text characteristics to create different embeddings
	textLen := float32(len(text))
	for i := range vec {
		if i < len(text) {
			vec[i] = float32(text[i]) / 255.0
		} else {
			vec[i] = textLen / 1000.0
		}
	}

	return vec, nil
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float32
		b        []float32
		expected float32
	}{
		{
			name:     "identical vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{1, 0, 0},
			expected: 1.0,
		},
		{
			name:     "orthogonal vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{0, 1, 0},
			expected: 0.0,
		},
		{
			name:     "opposite vectors",
			a:        []float32{1, 0, 0},
			b:        []float32{-1, 0, 0},
			expected: -1.0,
		},
		{
			name:     "similar vectors",
			a:        []float32{1, 1, 0},
			b:        []float32{1, 0.9, 0},
			expected: 0.99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			similarity := cosineSimilarity(tt.a, tt.b)

			// Allow small floating point errors
			diff := similarity - tt.expected
			if diff < 0 {
				diff = -diff
			}

			if diff > 0.01 {
				t.Errorf("cosineSimilarity() = %f, expected ~%f", similarity, tt.expected)
			}
		})
	}
}

func TestVectorStore_StoreAndRetrieve(t *testing.T) {
	// Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, ".local", "share", "percepta", "embeddings.db")

	// Override home directory for testing
	os.Setenv("HOME", tmpDir)
	os.MkdirAll(filepath.Join(tmpDir, ".local", "share", "percepta"), 0755)
	defer os.Unsetenv("HOME")

	// Create vector store with mock embedder
	store, err := NewVectorStoreWithEmbedder(&MockEmbedder{})
	if err != nil {
		t.Fatalf("Failed to create vector store: %v", err)
	}
	defer store.Close()

	// Store some embeddings
	patterns := map[string]string{
		"pattern1": "void LED_Blink(void) { /* blink LED */ }",
		"pattern2": "void LED_Toggle(void) { /* toggle LED */ }",
		"pattern3": "void UART_Send(uint8_t data) { /* send uart */ }",
	}

	for id, code := range patterns {
		if err := store.StoreEmbedding(id, code); err != nil {
			t.Errorf("Failed to store embedding for %s: %v", id, err)
		}
	}

	// Verify count
	if count := store.Count(); count != len(patterns) {
		t.Errorf("Count() = %d, expected %d", count, len(patterns))
	}

	// Verify HasEmbedding
	if !store.HasEmbedding("pattern1") {
		t.Error("HasEmbedding() = false, expected true for pattern1")
	}

	if store.HasEmbedding("nonexistent") {
		t.Error("HasEmbedding() = true, expected false for nonexistent pattern")
	}

	// Verify database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}

func TestVectorStore_FindSimilar(t *testing.T) {
	// Create temporary database
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	os.MkdirAll(filepath.Join(tmpDir, ".local", "share", "percepta"), 0755)
	defer os.Unsetenv("HOME")

	// Create vector store
	store, err := NewVectorStoreWithEmbedder(&MockEmbedder{})
	if err != nil {
		t.Fatalf("Failed to create vector store: %v", err)
	}
	defer store.Close()

	// Store LED-related patterns
	ledPatterns := map[string]string{
		"led_blink":  "void LED_Blink(void) { /* blink LED at 1Hz */ }",
		"led_toggle": "void LED_Toggle(void) { /* toggle LED state */ }",
		"led_pwm":    "void LED_PWM(uint8_t duty) { /* PWM LED */ }",
	}

	for id, code := range ledPatterns {
		if err := store.StoreEmbedding(id, code); err != nil {
			t.Fatalf("Failed to store embedding: %v", err)
		}
	}

	// Store unrelated pattern
	if err := store.StoreEmbedding("uart_send", "void UART_Send(uint8_t data) { /* send */ }"); err != nil {
		t.Fatalf("Failed to store embedding: %v", err)
	}

	// Query for LED blink
	query := "void LED_Blink(void) { /* blink LED at 1Hz */ }"
	matches, err := store.FindSimilar(query, 3)
	if err != nil {
		t.Fatalf("FindSimilar() failed: %v", err)
	}

	if len(matches) != 3 {
		t.Errorf("FindSimilar() returned %d matches, expected 3", len(matches))
	}

	// First match should be exact match (led_blink)
	if matches[0].PatternID != "led_blink" {
		t.Errorf("First match = %s, expected led_blink", matches[0].PatternID)
	}

	// Similarity should be very high (close to 1.0)
	if matches[0].Similarity < 0.99 {
		t.Errorf("First match similarity = %f, expected > 0.99", matches[0].Similarity)
	}

	// Results should be sorted by similarity (descending)
	for i := 1; i < len(matches); i++ {
		if matches[i].Similarity > matches[i-1].Similarity {
			t.Errorf("Results not sorted: match[%d].Similarity (%f) > match[%d].Similarity (%f)",
				i, matches[i].Similarity, i-1, matches[i-1].Similarity)
		}
	}
}

func TestVectorStore_Persistence(t *testing.T) {
	// Create temporary database
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	os.MkdirAll(filepath.Join(tmpDir, ".local", "share", "percepta"), 0755)
	defer os.Unsetenv("HOME")

	// Create first vector store
	store1, err := NewVectorStoreWithEmbedder(&MockEmbedder{})
	if err != nil {
		t.Fatalf("Failed to create vector store: %v", err)
	}

	// Store embeddings
	if err := store1.StoreEmbedding("pattern1", "code1"); err != nil {
		t.Fatalf("Failed to store embedding: %v", err)
	}
	if err := store1.StoreEmbedding("pattern2", "code2"); err != nil {
		t.Fatalf("Failed to store embedding: %v", err)
	}

	// Close first store
	store1.Close()

	// Create second vector store (should load from database)
	store2, err := NewVectorStoreWithEmbedder(&MockEmbedder{})
	if err != nil {
		t.Fatalf("Failed to create second vector store: %v", err)
	}
	defer store2.Close()

	// Verify embeddings were loaded
	if count := store2.Count(); count != 2 {
		t.Errorf("After reload, Count() = %d, expected 2", count)
	}

	if !store2.HasEmbedding("pattern1") {
		t.Error("After reload, pattern1 not found")
	}

	if !store2.HasEmbedding("pattern2") {
		t.Error("After reload, pattern2 not found")
	}
}

func TestVectorStore_ReplaceEmbedding(t *testing.T) {
	// Create temporary database
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	os.MkdirAll(filepath.Join(tmpDir, ".local", "share", "percepta"), 0755)
	defer os.Unsetenv("HOME")

	// Create vector store
	store, err := NewVectorStoreWithEmbedder(&MockEmbedder{})
	if err != nil {
		t.Fatalf("Failed to create vector store: %v", err)
	}
	defer store.Close()

	// Store initial embedding
	if err := store.StoreEmbedding("pattern1", "original code"); err != nil {
		t.Fatalf("Failed to store embedding: %v", err)
	}

	// Replace with new embedding
	if err := store.StoreEmbedding("pattern1", "updated code"); err != nil {
		t.Fatalf("Failed to replace embedding: %v", err)
	}

	// Count should still be 1
	if count := store.Count(); count != 1 {
		t.Errorf("After replacement, Count() = %d, expected 1", count)
	}
}

func TestVectorStore_EmptyQuery(t *testing.T) {
	// Create temporary database
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	os.MkdirAll(filepath.Join(tmpDir, ".local", "share", "percepta"), 0755)
	defer os.Unsetenv("HOME")

	// Create vector store
	store, err := NewVectorStoreWithEmbedder(&MockEmbedder{})
	if err != nil {
		t.Fatalf("Failed to create vector store: %v", err)
	}
	defer store.Close()

	// Query on empty store
	matches, err := store.FindSimilar("query", 5)
	if err != nil {
		t.Fatalf("FindSimilar() on empty store failed: %v", err)
	}

	if len(matches) != 0 {
		t.Errorf("FindSimilar() on empty store returned %d matches, expected 0", len(matches))
	}
}
