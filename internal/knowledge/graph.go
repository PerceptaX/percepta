package knowledge

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// Graph represents an in-memory graph with SQLite persistence
type Graph struct {
	patterns     map[string]*PatternNode
	observations map[string]*ObservationNode
	styleResults map[string]*StyleResultNode
	edges        map[string]*Edge
	db           *sql.DB
}

// NewGraph creates a new graph instance with SQLite persistence
// Database path: ~/.local/share/percepta/knowledge.db
func NewGraph() (*Graph, error) {
	// Get database path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dbDir := filepath.Join(homeDir, ".local", "share", "percepta")
	dbPath := filepath.Join(dbDir, "knowledge.db")

	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	graph := &Graph{
		patterns:     make(map[string]*PatternNode),
		observations: make(map[string]*ObservationNode),
		styleResults: make(map[string]*StyleResultNode),
		edges:        make(map[string]*Edge),
		db:           db,
	}

	// Initialize schema
	if err := graph.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Load existing data from database into memory
	if err := graph.loadFromDB(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to load graph data: %w", err)
	}

	return graph, nil
}

// initSchema creates the database schema for persisting graph data
func (g *Graph) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS patterns (
		id TEXT PRIMARY KEY,
		code TEXT NOT NULL,
		spec TEXT NOT NULL,
		board_type TEXT NOT NULL,
		style_compliant INTEGER NOT NULL,
		created_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS observations (
		id TEXT PRIMARY KEY,
		device_id TEXT NOT NULL,
		firmware TEXT NOT NULL,
		timestamp DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS style_results (
		id TEXT PRIMARY KEY,
		compliant INTEGER NOT NULL,
		auto_fixed INTEGER NOT NULL,
		viol_count INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS edges (
		id TEXT PRIMARY KEY,
		from_node TEXT NOT NULL,
		to_node TEXT NOT NULL,
		type TEXT NOT NULL,
		created DATETIME NOT NULL,
		metadata TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_patterns_board ON patterns(board_type);
	CREATE INDEX IF NOT EXISTS idx_patterns_spec ON patterns(spec);
	CREATE INDEX IF NOT EXISTS idx_edges_from ON edges(from_node, type);
	CREATE INDEX IF NOT EXISTS idx_edges_to ON edges(to_node, type);
	`

	_, err := g.db.Exec(schema)
	return err
}

// loadFromDB loads all graph data from SQLite into memory
func (g *Graph) loadFromDB() error {
	// Load patterns
	rows, err := g.db.Query("SELECT id, code, spec, board_type, style_compliant, created_at FROM patterns")
	if err != nil {
		return fmt.Errorf("failed to query patterns: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p PatternNode
		var styleCompliantInt int
		var createdAt string

		err := rows.Scan(&p.ID, &p.Code, &p.Spec, &p.BoardType, &styleCompliantInt, &createdAt)
		if err != nil {
			return fmt.Errorf("failed to scan pattern: %w", err)
		}

		p.StyleCompliant = styleCompliantInt == 1
		p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		g.patterns[p.ID] = &p
	}

	// Load observations
	rows, err = g.db.Query("SELECT id, device_id, firmware, timestamp FROM observations")
	if err != nil {
		return fmt.Errorf("failed to query observations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var o ObservationNode
		var timestamp string

		err := rows.Scan(&o.ID, &o.DeviceID, &o.Firmware, &timestamp)
		if err != nil {
			return fmt.Errorf("failed to scan observation: %w", err)
		}

		o.Timestamp, _ = time.Parse(time.RFC3339, timestamp)
		g.observations[o.ID] = &o
	}

	// Load style results
	rows, err = g.db.Query("SELECT id, compliant, auto_fixed, viol_count FROM style_results")
	if err != nil {
		return fmt.Errorf("failed to query style results: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var s StyleResultNode
		var compliantInt, autoFixedInt int

		err := rows.Scan(&s.ID, &compliantInt, &autoFixedInt, &s.ViolCount)
		if err != nil {
			return fmt.Errorf("failed to scan style result: %w", err)
		}

		s.Compliant = compliantInt == 1
		s.AutoFixed = autoFixedInt == 1
		g.styleResults[s.ID] = &s
	}

	// Load edges
	rows, err = g.db.Query("SELECT id, from_node, to_node, type, created, metadata FROM edges")
	if err != nil {
		return fmt.Errorf("failed to query edges: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var e Edge
		var created string
		var metadata sql.NullString

		err := rows.Scan(&e.ID, &e.From, &e.To, &e.Type, &created, &metadata)
		if err != nil {
			return fmt.Errorf("failed to scan edge: %w", err)
		}

		e.Created, _ = time.Parse(time.RFC3339, created)
		if metadata.Valid {
			e.Metadata = metadata.String
		}
		g.edges[e.ID] = &e
	}

	return nil
}

// AddPattern stores a pattern node in memory and persists to database
// Returns the generated ID for the pattern
func (g *Graph) AddPattern(pattern PatternNode) (string, error) {
	// Generate ID if not set
	if pattern.ID == "" {
		pattern.ID = generateID(pattern.Code + pattern.Spec + pattern.BoardType)
	}

	if pattern.CreatedAt.IsZero() {
		pattern.CreatedAt = time.Now()
	}

	// Store in memory
	g.patterns[pattern.ID] = &pattern

	// Persist to database
	query := `
	INSERT OR REPLACE INTO patterns (id, code, spec, board_type, style_compliant, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	styleCompliantInt := 0
	if pattern.StyleCompliant {
		styleCompliantInt = 1
	}

	_, err := g.db.Exec(query, pattern.ID, pattern.Code, pattern.Spec, pattern.BoardType,
		styleCompliantInt, pattern.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return "", fmt.Errorf("failed to insert pattern: %w", err)
	}

	return pattern.ID, nil
}

// AddObservation stores an observation node in memory and persists to database
func (g *Graph) AddObservation(obs ObservationNode) error {
	if obs.Timestamp.IsZero() {
		obs.Timestamp = time.Now()
	}

	// Store in memory
	g.observations[obs.ID] = &obs

	// Persist to database
	query := `
	INSERT OR REPLACE INTO observations (id, device_id, firmware, timestamp)
	VALUES (?, ?, ?, ?)
	`

	_, err := g.db.Exec(query, obs.ID, obs.DeviceID, obs.Firmware, obs.Timestamp.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("failed to insert observation: %w", err)
	}

	return nil
}

// AddStyleResult stores a style result node in memory and persists to database
func (g *Graph) AddStyleResult(result StyleResultNode) error {
	// Store in memory
	g.styleResults[result.ID] = &result

	// Persist to database
	query := `
	INSERT OR REPLACE INTO style_results (id, compliant, auto_fixed, viol_count)
	VALUES (?, ?, ?, ?)
	`

	compliantInt := 0
	if result.Compliant {
		compliantInt = 1
	}

	autoFixedInt := 0
	if result.AutoFixed {
		autoFixedInt = 1
	}

	_, err := g.db.Exec(query, result.ID, compliantInt, autoFixedInt, result.ViolCount)
	if err != nil {
		return fmt.Errorf("failed to insert style result: %w", err)
	}

	return nil
}

// AddEdge creates a relationship between two nodes
func (g *Graph) AddEdge(from, to string, relType Relationship) error {
	edge := Edge{
		ID:      generateID(from + string(relType) + to),
		From:    from,
		To:      to,
		Type:    relType,
		Created: time.Now(),
	}

	// Store in memory
	g.edges[edge.ID] = &edge

	// Persist to database
	query := `
	INSERT OR REPLACE INTO edges (id, from_node, to_node, type, created, metadata)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := g.db.Exec(query, edge.ID, edge.From, edge.To, string(edge.Type),
		edge.Created.Format(time.RFC3339), edge.Metadata)
	if err != nil {
		return fmt.Errorf("failed to insert edge: %w", err)
	}

	return nil
}

// GetPattern retrieves a pattern by ID
func (g *Graph) GetPattern(id string) (*PatternNode, error) {
	pattern, ok := g.patterns[id]
	if !ok {
		return nil, fmt.Errorf("pattern not found: %s", id)
	}
	return pattern, nil
}

// GetObservation retrieves an observation by ID
func (g *Graph) GetObservation(id string) (*ObservationNode, error) {
	obs, ok := g.observations[id]
	if !ok {
		return nil, fmt.Errorf("observation not found: %s", id)
	}
	return obs, nil
}

// GetStyleResult retrieves a style result by ID
func (g *Graph) GetStyleResult(id string) (*StyleResultNode, error) {
	result, ok := g.styleResults[id]
	if !ok {
		return nil, fmt.Errorf("style result not found: %s", id)
	}
	return result, nil
}

// QueryPatternsByBoard finds all patterns that run on a specific board type
func (g *Graph) QueryPatternsByBoard(boardType string) ([]*PatternNode, error) {
	var patterns []*PatternNode

	for _, pattern := range g.patterns {
		if pattern.BoardType == boardType {
			patterns = append(patterns, pattern)
		}
	}

	return patterns, nil
}

// QueryPatternsBySpec finds patterns with matching or similar specs
func (g *Graph) QueryPatternsBySpec(spec string) ([]*PatternNode, error) {
	var patterns []*PatternNode

	for _, pattern := range g.patterns {
		// Exact match for now - semantic search to be added in Phase 06-02
		if pattern.Spec == spec {
			patterns = append(patterns, pattern)
		}
	}

	return patterns, nil
}

// GetEdgesFrom retrieves all edges starting from a node
func (g *Graph) GetEdgesFrom(nodeID string, relType Relationship) ([]*Edge, error) {
	var edges []*Edge

	for _, edge := range g.edges {
		if edge.From == nodeID && (relType == "" || edge.Type == relType) {
			edges = append(edges, edge)
		}
	}

	return edges, nil
}

// GetEdgesTo retrieves all edges ending at a node
func (g *Graph) GetEdgesTo(nodeID string, relType Relationship) ([]*Edge, error) {
	var edges []*Edge

	for _, edge := range g.edges {
		if edge.To == nodeID && (relType == "" || edge.Type == relType) {
			edges = append(edges, edge)
		}
	}

	return edges, nil
}

// Stats returns statistics about the graph
func (g *Graph) Stats() map[string]int {
	return map[string]int{
		"patterns":      len(g.patterns),
		"observations":  len(g.observations),
		"style_results": len(g.styleResults),
		"edges":         len(g.edges),
	}
}

// Close closes the database connection
func (g *Graph) Close() error {
	return g.db.Close()
}

// generateID creates a deterministic ID from content
func generateID(content string) string {
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash[:16]) // Use first 16 bytes (32 hex chars)
}

// ToJSON serializes a node to JSON for metadata storage
func toJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
