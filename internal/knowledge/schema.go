package knowledge

import "time"

// Node types for the knowledge graph

// PatternNode represents a validated code pattern
type PatternNode struct {
	ID             string    `json:"id"`
	Code           string    `json:"code"`            // C source code
	Spec           string    `json:"spec"`            // Natural language spec
	BoardType      string    `json:"board_type"`      // "esp32", "stm32", etc.
	StyleCompliant bool      `json:"style_compliant"` // Passed BARR-C validation
	CreatedAt      time.Time `json:"created_at"`
}

// ObservationNode represents a hardware observation
type ObservationNode struct {
	ID        string    `json:"id"`
	DeviceID  string    `json:"device_id"`
	Firmware  string    `json:"firmware"`
	Timestamp time.Time `json:"timestamp"`
}

// StyleResultNode represents style checking results
type StyleResultNode struct {
	ID        string `json:"id"`
	Compliant bool   `json:"compliant"`  // Overall compliance status
	AutoFixed bool   `json:"auto_fixed"` // Whether violations were auto-fixed
	ViolCount int    `json:"viol_count"` // Number of violations found
}

// Relationship types define how nodes connect
type Relationship string

const (
	IMPLEMENTED_BY Relationship = "IMPLEMENTED_BY" // Spec -> Code
	RUNS_ON        Relationship = "RUNS_ON"        // Code -> Board
	PRODUCES       Relationship = "PRODUCES"       // Code -> Observation
	VALIDATED_BY   Relationship = "VALIDATED_BY"   // Code -> StyleResult
	SIMILAR_TO     Relationship = "SIMILAR_TO"     // Code -> Code (for semantic search)
)

// Edge represents a relationship between two nodes
type Edge struct {
	ID       string       `json:"id"`
	From     string       `json:"from"`     // Source node ID
	To       string       `json:"to"`       // Target node ID
	Type     Relationship `json:"type"`     // Relationship type
	Created  time.Time    `json:"created"`  // When relationship was created
	Metadata string       `json:"metadata"` // JSON metadata for relationship
}
