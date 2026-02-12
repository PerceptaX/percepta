package knowledge

import (
	"fmt"
	"time"

	"github.com/perceptumx/percepta/internal/core"
	"github.com/perceptumx/percepta/internal/storage"
	"github.com/perceptumx/percepta/internal/style"
)

// PatternStore integrates perception, style checking, and knowledge graph
// to store only validated firmware patterns
type PatternStore struct {
	graph      *Graph
	styleCheck *style.StyleChecker
	storage    *storage.SQLiteStorage
}

// NewPatternStore creates a new pattern store with all dependencies
func NewPatternStore() (*PatternStore, error) {
	graph, err := NewGraph()
	if err != nil {
		return nil, fmt.Errorf("failed to create graph: %w", err)
	}

	styleCheck := style.NewStyleChecker()

	sqliteStorage, err := storage.NewSQLiteStorage()
	if err != nil {
		graph.Close()
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	return &PatternStore{
		graph:      graph,
		styleCheck: styleCheck,
		storage:    sqliteStorage,
	}, nil
}

// StoreValidatedPattern stores a pattern only if it passes all validation
// Requirements:
// 1. Code must be BARR-C compliant (style checking)
// 2. Must have corresponding observation in storage
// 3. Creates full relationship graph
func (p *PatternStore) StoreValidatedPattern(
	spec string,
	code string,
	deviceID string,
	firmware string,
) (string, error) {
	// 1. Validate style compliance
	violations, err := p.styleCheck.CheckSource([]byte(code), "pattern.c")
	if err != nil {
		return "", fmt.Errorf("style check failed: %w", err)
	}

	if len(violations) > 0 {
		return "", fmt.Errorf("code not BARR-C compliant: %d violations found", len(violations))
	}

	// 2. Get observation for this firmware
	obs, err := p.storage.GetLatestForFirmware(deviceID, firmware)
	if err != nil {
		return "", fmt.Errorf("no observation found for device %s with firmware %s: %w",
			deviceID, firmware, err)
	}

	// 3. Create pattern node
	pattern := PatternNode{
		Code:           code,
		Spec:           spec,
		BoardType:      getBoardType(deviceID),
		StyleCompliant: true,
		CreatedAt:      time.Now(),
	}

	patternID, err := p.graph.AddPattern(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to add pattern: %w", err)
	}

	// 4. Create observation node in knowledge graph
	obsNode := ObservationNode{
		ID:        obs.ID,
		DeviceID:  obs.DeviceID,
		Firmware:  obs.FirmwareHash,
		Timestamp: obs.Timestamp,
	}

	if err := p.graph.AddObservation(obsNode); err != nil {
		return "", fmt.Errorf("failed to add observation node: %w", err)
	}

	// 5. Create style result node
	styleResult := StyleResultNode{
		ID:        generateID("style-" + patternID),
		Compliant: true,
		AutoFixed: false,
		ViolCount: 0,
	}

	if err := p.graph.AddStyleResult(styleResult); err != nil {
		return "", fmt.Errorf("failed to add style result: %w", err)
	}

	// 6. Create relationships
	// Spec -> Pattern (IMPLEMENTED_BY)
	specID := generateID("spec-" + spec)
	if err := p.graph.AddEdge(specID, patternID, IMPLEMENTED_BY); err != nil {
		return "", fmt.Errorf("failed to add IMPLEMENTED_BY edge: %w", err)
	}

	// Pattern -> Board (RUNS_ON)
	boardID := generateID("board-" + pattern.BoardType)
	if err := p.graph.AddEdge(patternID, boardID, RUNS_ON); err != nil {
		return "", fmt.Errorf("failed to add RUNS_ON edge: %w", err)
	}

	// Pattern -> Observation (PRODUCES)
	if err := p.graph.AddEdge(patternID, obs.ID, PRODUCES); err != nil {
		return "", fmt.Errorf("failed to add PRODUCES edge: %w", err)
	}

	// Pattern -> StyleResult (VALIDATED_BY)
	if err := p.graph.AddEdge(patternID, styleResult.ID, VALIDATED_BY); err != nil {
		return "", fmt.Errorf("failed to add VALIDATED_BY edge: %w", err)
	}

	return patternID, nil
}

// QueryPatternsByBoard retrieves all validated patterns that run on a specific board type
func (p *PatternStore) QueryPatternsByBoard(boardType string) ([]*PatternNode, error) {
	return p.graph.QueryPatternsByBoard(boardType)
}

// QueryPatternsBySpec finds patterns implementing a specific specification
func (p *PatternStore) QueryPatternsBySpec(spec string) ([]*PatternNode, error) {
	return p.graph.QueryPatternsBySpec(spec)
}

// GetPatternWithObservation retrieves a pattern along with its observation
func (p *PatternStore) GetPatternWithObservation(patternID string) (*PatternNode, *core.Observation, error) {
	// Get pattern
	pattern, err := p.graph.GetPattern(patternID)
	if err != nil {
		return nil, nil, fmt.Errorf("pattern not found: %w", err)
	}

	// Find observation edge
	edges, err := p.graph.GetEdgesFrom(patternID, PRODUCES)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get edges: %w", err)
	}

	if len(edges) == 0 {
		return pattern, nil, nil // Pattern exists but no observation linked
	}

	// Get observation from perception storage
	obsID := edges[0].To
	obsNode, err := p.graph.GetObservation(obsID)
	if err != nil {
		return pattern, nil, fmt.Errorf("failed to get observation node: %w", err)
	}

	// Retrieve full observation from perception storage
	obs, err := p.storage.GetLatestForFirmware(obsNode.DeviceID, obsNode.Firmware)
	if err != nil {
		return pattern, nil, fmt.Errorf("failed to get full observation: %w", err)
	}

	return pattern, obs, nil
}

// Stats returns statistics about stored patterns
func (p *PatternStore) Stats() map[string]int {
	return p.graph.Stats()
}

// Close closes all underlying connections
func (p *PatternStore) Close() error {
	if err := p.graph.Close(); err != nil {
		return err
	}
	return p.storage.Close()
}

// getBoardType extracts board type from device ID
// For now, simple heuristic - can be made more sophisticated later
func getBoardType(deviceID string) string {
	// Simple pattern matching for common boards
	// esp32-* -> esp32
	// stm32-* -> stm32
	// arduino-* -> arduino
	// etc.

	if len(deviceID) >= 5 {
		if deviceID[:5] == "esp32" {
			return "esp32"
		} else if deviceID[:5] == "stm32" {
			return "stm32"
		}
	}

	if len(deviceID) >= 7 {
		if deviceID[:7] == "arduino" {
			return "arduino"
		}
	}

	// Default: use first part before dash or entire device ID
	for i, c := range deviceID {
		if c == '-' {
			return deviceID[:i]
		}
	}

	return deviceID
}
