package knowledge

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestGraph(t *testing.T) (*Graph, func()) {
	// Create temporary database
	tmpDir := t.TempDir()

	// Override database path for testing
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	// Create .local/share/percepta directory
	dbDir := filepath.Join(tmpDir, ".local", "share", "percepta")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("failed to create test db dir: %v", err)
	}

	graph, err := NewGraph()
	if err != nil {
		t.Fatalf("failed to create test graph: %v", err)
	}

	cleanup := func() {
		graph.Close()
		os.Setenv("HOME", oldHome)
		os.RemoveAll(tmpDir)
	}

	return graph, cleanup
}

func TestGraph_AddPattern(t *testing.T) {
	graph, cleanup := setupTestGraph(t)
	defer cleanup()

	pattern := PatternNode{
		Code:           "void LED_Blink() { /* ... */ }",
		Spec:           "Blink LED at 1Hz",
		BoardType:      "esp32",
		StyleCompliant: true,
	}

	patternID, err := graph.AddPattern(pattern)
	if err != nil {
		t.Fatalf("AddPattern failed: %v", err)
	}

	if patternID == "" {
		t.Fatal("AddPattern returned empty ID")
	}

	// Verify pattern was stored in memory
	if len(graph.patterns) != 1 {
		t.Errorf("expected 1 pattern in memory, got %d", len(graph.patterns))
	}

	// Verify pattern can be retrieved
	storedPattern, err := graph.GetPattern(patternID)
	if err != nil {
		t.Fatalf("GetPattern failed: %v", err)
	}

	if storedPattern.Code != pattern.Code {
		t.Errorf("expected code %s, got %s", pattern.Code, storedPattern.Code)
	}

	if storedPattern.Spec != pattern.Spec {
		t.Errorf("expected spec %s, got %s", pattern.Spec, storedPattern.Spec)
	}

	if storedPattern.BoardType != pattern.BoardType {
		t.Errorf("expected board type %s, got %s", pattern.BoardType, storedPattern.BoardType)
	}

	if storedPattern.StyleCompliant != pattern.StyleCompliant {
		t.Errorf("expected style compliant %v, got %v", pattern.StyleCompliant, storedPattern.StyleCompliant)
	}
}

func TestGraph_AddObservation(t *testing.T) {
	graph, cleanup := setupTestGraph(t)
	defer cleanup()

	obs := ObservationNode{
		ID:       "obs-123",
		DeviceID: "esp32-dev-1",
		Firmware: "v1.0.0",
	}

	err := graph.AddObservation(obs)
	if err != nil {
		t.Fatalf("AddObservation failed: %v", err)
	}

	// Verify observation was stored
	storedObs, err := graph.GetObservation(obs.ID)
	if err != nil {
		t.Fatalf("GetObservation failed: %v", err)
	}

	if storedObs.DeviceID != obs.DeviceID {
		t.Errorf("expected device ID %s, got %s", obs.DeviceID, storedObs.DeviceID)
	}

	if storedObs.Firmware != obs.Firmware {
		t.Errorf("expected firmware %s, got %s", obs.Firmware, storedObs.Firmware)
	}
}

func TestGraph_AddStyleResult(t *testing.T) {
	graph, cleanup := setupTestGraph(t)
	defer cleanup()

	result := StyleResultNode{
		ID:        "style-123",
		Compliant: true,
		AutoFixed: false,
		ViolCount: 0,
	}

	err := graph.AddStyleResult(result)
	if err != nil {
		t.Fatalf("AddStyleResult failed: %v", err)
	}

	// Verify style result was stored
	storedResult, err := graph.GetStyleResult(result.ID)
	if err != nil {
		t.Fatalf("GetStyleResult failed: %v", err)
	}

	if storedResult.Compliant != result.Compliant {
		t.Errorf("expected compliant %v, got %v", result.Compliant, storedResult.Compliant)
	}

	if storedResult.ViolCount != result.ViolCount {
		t.Errorf("expected viol count %d, got %d", result.ViolCount, storedResult.ViolCount)
	}
}

func TestGraph_AddEdge(t *testing.T) {
	graph, cleanup := setupTestGraph(t)
	defer cleanup()

	// Add pattern and observation first
	pattern := PatternNode{
		Code:      "void test() {}",
		Spec:      "Test spec",
		BoardType: "esp32",
	}
	patternID, err := graph.AddPattern(pattern)
	if err != nil {
		t.Fatalf("AddPattern failed: %v", err)
	}

	obs := ObservationNode{
		ID:       "obs-123",
		DeviceID: "esp32-dev-1",
		Firmware: "v1.0.0",
	}
	err = graph.AddObservation(obs)
	if err != nil {
		t.Fatalf("AddObservation failed: %v", err)
	}

	// Create edge
	err = graph.AddEdge(patternID, obs.ID, PRODUCES)
	if err != nil {
		t.Fatalf("AddEdge failed: %v", err)
	}

	// Verify edge was created
	edges, err := graph.GetEdgesFrom(patternID, PRODUCES)
	if err != nil {
		t.Fatalf("GetEdgesFrom failed: %v", err)
	}

	if len(edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(edges))
	}

	if edges[0].From != patternID {
		t.Errorf("expected from %s, got %s", patternID, edges[0].From)
	}

	if edges[0].To != obs.ID {
		t.Errorf("expected to %s, got %s", obs.ID, edges[0].To)
	}

	if edges[0].Type != PRODUCES {
		t.Errorf("expected type %s, got %s", PRODUCES, edges[0].Type)
	}
}

func TestGraph_QueryPatternsByBoard(t *testing.T) {
	graph, cleanup := setupTestGraph(t)
	defer cleanup()

	// Add patterns for different boards
	esp32Pattern := PatternNode{
		Code:      "void esp_test() {}",
		Spec:      "ESP32 test",
		BoardType: "esp32",
	}
	_, err := graph.AddPattern(esp32Pattern)
	if err != nil {
		t.Fatalf("AddPattern failed: %v", err)
	}

	stm32Pattern := PatternNode{
		Code:      "void stm_test() {}",
		Spec:      "STM32 test",
		BoardType: "stm32",
	}
	_, err = graph.AddPattern(stm32Pattern)
	if err != nil {
		t.Fatalf("AddPattern failed: %v", err)
	}

	// Query ESP32 patterns
	esp32Patterns, err := graph.QueryPatternsByBoard("esp32")
	if err != nil {
		t.Fatalf("QueryPatternsByBoard failed: %v", err)
	}

	if len(esp32Patterns) != 1 {
		t.Fatalf("expected 1 ESP32 pattern, got %d", len(esp32Patterns))
	}

	if esp32Patterns[0].BoardType != "esp32" {
		t.Errorf("expected board type esp32, got %s", esp32Patterns[0].BoardType)
	}

	// Query STM32 patterns
	stm32Patterns, err := graph.QueryPatternsByBoard("stm32")
	if err != nil {
		t.Fatalf("QueryPatternsByBoard failed: %v", err)
	}

	if len(stm32Patterns) != 1 {
		t.Fatalf("expected 1 STM32 pattern, got %d", len(stm32Patterns))
	}
}

func TestGraph_QueryPatternsBySpec(t *testing.T) {
	graph, cleanup := setupTestGraph(t)
	defer cleanup()

	// Add patterns with different specs
	pattern1 := PatternNode{
		Code:      "void blink() {}",
		Spec:      "Blink LED at 1Hz",
		BoardType: "esp32",
	}
	_, err := graph.AddPattern(pattern1)
	if err != nil {
		t.Fatalf("AddPattern failed: %v", err)
	}

	pattern2 := PatternNode{
		Code:      "void blink2() {}",
		Spec:      "Blink LED at 2Hz",
		BoardType: "esp32",
	}
	_, err = graph.AddPattern(pattern2)
	if err != nil {
		t.Fatalf("AddPattern failed: %v", err)
	}

	// Query by spec
	patterns, err := graph.QueryPatternsBySpec("Blink LED at 1Hz")
	if err != nil {
		t.Fatalf("QueryPatternsBySpec failed: %v", err)
	}

	if len(patterns) != 1 {
		t.Fatalf("expected 1 pattern, got %d", len(patterns))
	}

	if patterns[0].Spec != "Blink LED at 1Hz" {
		t.Errorf("expected spec 'Blink LED at 1Hz', got %s", patterns[0].Spec)
	}
}

func TestGraph_Persistence(t *testing.T) {
	graph, cleanup := setupTestGraph(t)

	// Add pattern
	pattern := PatternNode{
		Code:           "void test() {}",
		Spec:           "Test pattern",
		BoardType:      "esp32",
		StyleCompliant: true,
	}
	patternID, err := graph.AddPattern(pattern)
	if err != nil {
		t.Fatalf("AddPattern failed: %v", err)
	}

	// Close and reopen graph
	graph.Close()

	// Don't cleanup yet - we need to reopen the database
	graph2, err := NewGraph()
	if err != nil {
		t.Fatalf("failed to reopen graph: %v", err)
	}
	defer cleanup()
	defer graph2.Close()

	// Verify pattern was persisted
	storedPattern, err := graph2.GetPattern(patternID)
	if err != nil {
		t.Fatalf("GetPattern failed after reload: %v", err)
	}

	if storedPattern.Code != pattern.Code {
		t.Errorf("expected code %s, got %s", pattern.Code, storedPattern.Code)
	}

	if storedPattern.Spec != pattern.Spec {
		t.Errorf("expected spec %s, got %s", pattern.Spec, storedPattern.Spec)
	}
}

func TestGraph_Stats(t *testing.T) {
	graph, cleanup := setupTestGraph(t)
	defer cleanup()

	// Add various nodes
	pattern := PatternNode{
		Code:      "void test() {}",
		Spec:      "Test",
		BoardType: "esp32",
	}
	patternID, _ := graph.AddPattern(pattern)

	obs := ObservationNode{
		ID:       "obs-1",
		DeviceID: "dev-1",
		Firmware: "v1",
	}
	graph.AddObservation(obs)

	result := StyleResultNode{
		ID:        "style-1",
		Compliant: true,
	}
	graph.AddStyleResult(result)

	graph.AddEdge(patternID, obs.ID, PRODUCES)

	// Get stats
	stats := graph.Stats()

	if stats["patterns"] != 1 {
		t.Errorf("expected 1 pattern, got %d", stats["patterns"])
	}

	if stats["observations"] != 1 {
		t.Errorf("expected 1 observation, got %d", stats["observations"])
	}

	if stats["style_results"] != 1 {
		t.Errorf("expected 1 style result, got %d", stats["style_results"])
	}

	if stats["edges"] != 1 {
		t.Errorf("expected 1 edge, got %d", stats["edges"])
	}
}

func TestGraph_MultipleEdgesFromNode(t *testing.T) {
	graph, cleanup := setupTestGraph(t)
	defer cleanup()

	// Add pattern
	pattern := PatternNode{
		Code:      "void test() {}",
		Spec:      "Test",
		BoardType: "esp32",
	}
	patternID, _ := graph.AddPattern(pattern)

	// Add multiple observations
	obs1 := ObservationNode{ID: "obs-1", DeviceID: "dev-1", Firmware: "v1"}
	obs2 := ObservationNode{ID: "obs-2", DeviceID: "dev-2", Firmware: "v1"}

	graph.AddObservation(obs1)
	graph.AddObservation(obs2)

	// Create edges
	graph.AddEdge(patternID, obs1.ID, PRODUCES)
	graph.AddEdge(patternID, obs2.ID, PRODUCES)

	// Query edges
	edges, err := graph.GetEdgesFrom(patternID, PRODUCES)
	if err != nil {
		t.Fatalf("GetEdgesFrom failed: %v", err)
	}

	if len(edges) != 2 {
		t.Fatalf("expected 2 edges, got %d", len(edges))
	}
}
