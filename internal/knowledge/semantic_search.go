package knowledge

import (
	"fmt"

	"github.com/perceptumx/percepta/internal/core"
)

// PatternResult represents a pattern match with similarity and confidence
type PatternResult struct {
	Pattern    *PatternNode
	Observation *core.Observation
	Similarity float32 // 0.0 to 1.0, from cosine similarity
	Confidence float32 // 0.0 to 1.0, adjusted for validation metadata
}

// SearchSimilarPatterns finds patterns similar to the query using semantic search
// Only returns validated patterns (StyleCompliant=true)
// Filters by board type if specified (empty string = all boards)
func (p *PatternStore) SearchSimilarPatterns(
	query string,
	boardType string,
	topK int,
) ([]PatternResult, error) {
	// Ensure vector store is initialized
	if p.vectorStore == nil {
		return nil, fmt.Errorf("vector store not initialized")
	}

	// Find similar patterns by embedding (get extra for filtering)
	matches, err := p.vectorStore.FindSimilar(query, topK*3)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// Filter and enrich results
	var results []PatternResult
	for _, match := range matches {
		// Get full pattern from graph
		pattern, err := p.graph.GetPattern(match.PatternID)
		if err != nil {
			// Pattern not found in graph, skip
			continue
		}

		// Filter by board type if specified
		if boardType != "" && pattern.BoardType != boardType {
			continue
		}

		// Safety check: only return validated patterns
		if !pattern.StyleCompliant {
			continue
		}

		// Get observation
		obs, err := p.getObservationForPattern(match.PatternID)
		if err != nil {
			// No observation, skip (shouldn't happen for validated patterns)
			continue
		}

		// Calculate confidence score
		confidence := calculateConfidence(match.Similarity, pattern, obs)

		results = append(results, PatternResult{
			Pattern:     pattern,
			Observation: obs,
			Similarity:  match.Similarity,
			Confidence:  confidence,
		})

		// Stop if we have enough results
		if len(results) >= topK {
			break
		}
	}

	return results, nil
}

// getObservationForPattern retrieves the observation linked to a pattern
func (p *PatternStore) getObservationForPattern(patternID string) (*core.Observation, error) {
	// Find PRODUCES edge
	edges, err := p.graph.GetEdgesFrom(patternID, PRODUCES)
	if err != nil {
		return nil, fmt.Errorf("failed to get edges: %w", err)
	}

	if len(edges) == 0 {
		return nil, fmt.Errorf("no observation linked to pattern")
	}

	// Get observation node
	obsNode, err := p.graph.GetObservation(edges[0].To)
	if err != nil {
		return nil, fmt.Errorf("failed to get observation node: %w", err)
	}

	// Get full observation from perception storage
	obs, err := p.storage.GetLatestForFirmware(obsNode.DeviceID, obsNode.Firmware)
	if err != nil {
		return nil, fmt.Errorf("failed to get observation: %w", err)
	}

	return obs, nil
}

// calculateConfidence computes confidence score based on similarity and validation metadata
// Returns value between 0.0 and 1.0
func calculateConfidence(similarity float32, pattern *PatternNode, obs *core.Observation) float32 {
	// Base confidence is the similarity score
	confidence := similarity

	// Boost for style-compliant patterns (already filtered, but explicit)
	if pattern.StyleCompliant {
		confidence *= 1.0 // No penalty
	} else {
		confidence *= 0.5 // Should never happen due to filtering
	}

	// Boost for patterns with observations
	if obs != nil && len(obs.Signals) > 0 {
		// More signals = more validation = higher confidence
		signalBoost := float32(len(obs.Signals)) / 10.0
		if signalBoost > 0.1 {
			signalBoost = 0.1 // Cap boost at 10%
		}
		confidence += signalBoost
	}

	// Ensure confidence stays in [0, 1]
	if confidence > 1.0 {
		confidence = 1.0
	}
	if confidence < 0.0 {
		confidence = 0.0
	}

	return confidence
}

// InitializeVectorStore sets up the vector store for semantic search
// Should be called after PatternStore is created
func (p *PatternStore) InitializeVectorStore() error {
	vectorStore, err := NewVectorStore()
	if err != nil {
		return fmt.Errorf("failed to initialize vector store: %w", err)
	}

	p.vectorStore = vectorStore
	return nil
}

// InitializeVectorStoreWithEmbedder sets up vector store with custom embedder
// Useful for testing
func (p *PatternStore) InitializeVectorStoreWithEmbedder(embedder EmbeddingProvider) error {
	vectorStore, err := NewVectorStoreWithEmbedder(embedder)
	if err != nil {
		return fmt.Errorf("failed to initialize vector store: %w", err)
	}

	p.vectorStore = vectorStore
	return nil
}

// CloseVectorStore closes the vector store if initialized
func (p *PatternStore) CloseVectorStore() error {
	if p.vectorStore != nil {
		return p.vectorStore.Close()
	}
	return nil
}
