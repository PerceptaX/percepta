package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/perceptumx/percepta/internal/core"
	_ "modernc.org/sqlite"
)

// SQLiteStorage provides persistent storage for observations
type SQLiteStorage struct {
	db *sql.DB
}

// NewSQLiteStorage creates a new SQLite storage instance
// Database path: ~/.local/share/percepta/percepta.db
func NewSQLiteStorage() (*SQLiteStorage, error) {
	// Get database path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dbDir := filepath.Join(homeDir, ".local", "share", "percepta")
	dbPath := filepath.Join(dbDir, "percepta.db")

	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &SQLiteStorage{db: db}

	// Initialize schema
	if err := storage.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return storage, nil
}

// initSchema creates the observations table and indexes
func (s *SQLiteStorage) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS observations (
		id TEXT PRIMARY KEY,
		device_id TEXT NOT NULL,
		firmware TEXT NOT NULL DEFAULT '',
		timestamp DATETIME NOT NULL,
		signals_json TEXT NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_device_firmware
	ON observations(device_id, firmware, timestamp);
	`

	_, err := s.db.Exec(schema)
	return err
}

// Save stores an observation in the database
func (s *SQLiteStorage) Save(obs core.Observation) error {
	// Serialize signals to JSON
	signalsJSON, err := json.Marshal(obs.Signals)
	if err != nil {
		return fmt.Errorf("failed to marshal signals: %w", err)
	}

	query := `
	INSERT INTO observations (id, device_id, firmware, timestamp, signals_json)
	VALUES (?, ?, ?, ?, ?)
	`

	_, err = s.db.Exec(query, obs.ID, obs.DeviceID, obs.FirmwareHash, obs.Timestamp, string(signalsJSON))
	if err != nil {
		return fmt.Errorf("failed to insert observation: %w", err)
	}

	return nil
}

// Query retrieves observations for a device, optionally limited
func (s *SQLiteStorage) Query(deviceID string, limit int) ([]core.Observation, error) {
	query := `
	SELECT id, device_id, firmware, timestamp, signals_json
	FROM observations
	WHERE device_id = ? OR ? = ''
	ORDER BY timestamp ASC
	`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := s.db.Query(query, deviceID, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query observations: %w", err)
	}
	defer rows.Close()

	return s.scanObservations(rows)
}

// QueryByFirmware retrieves observations for a specific device and firmware version
func (s *SQLiteStorage) QueryByFirmware(deviceID, firmware string, limit int) ([]core.Observation, error) {
	query := `
	SELECT id, device_id, firmware, timestamp, signals_json
	FROM observations
	WHERE device_id = ? AND firmware = ?
	ORDER BY timestamp ASC
	`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := s.db.Query(query, deviceID, firmware)
	if err != nil {
		return nil, fmt.Errorf("failed to query observations by firmware: %w", err)
	}
	defer rows.Close()

	return s.scanObservations(rows)
}

// GetLatestForFirmware retrieves the most recent observation for a device+firmware combination
func (s *SQLiteStorage) GetLatestForFirmware(deviceID, firmware string) (*core.Observation, error) {
	query := `
	SELECT id, device_id, firmware, timestamp, signals_json
	FROM observations
	WHERE device_id = ? AND firmware = ?
	ORDER BY timestamp DESC
	LIMIT 1
	`

	row := s.db.QueryRow(query, deviceID, firmware)

	var id, devID, fw, timestamp, signalsJSON string

	err := row.Scan(&id, &devID, &fw, &timestamp, &signalsJSON)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no observation found for device %s with firmware %s", deviceID, firmware)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan observation: %w", err)
	}

	// Build complete JSON for schema validation
	fullJSON := fmt.Sprintf(`{
		"id": "%s",
		"device_id": "%s",
		"firmware_hash": "%s",
		"timestamp": "%s",
		"signals": %s
	}`, id, devID, fw, timestamp, signalsJSON)

	// Validate and migrate if needed
	validator := core.NewSchemaValidator()
	obs, err := validator.ValidateAndMigrate([]byte(fullJSON))
	if err != nil {
		return nil, fmt.Errorf("schema validation failed: %w", err)
	}

	return obs, nil
}

// Count returns the total number of observations in the database
func (s *SQLiteStorage) Count() int {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM observations").Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

// scanObservations is a helper to scan multiple rows into observations
func (s *SQLiteStorage) scanObservations(rows *sql.Rows) ([]core.Observation, error) {
	var observations []core.Observation
	validator := core.NewSchemaValidator()

	for rows.Next() {
		var rawData string
		var timestamp string

		// Scan raw JSON data (we'll reconstruct the full JSON for validation)
		var id, deviceID, firmware string
		err := rows.Scan(&id, &deviceID, &firmware, &timestamp, &rawData)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Build complete JSON for schema validation
		fullJSON := fmt.Sprintf(`{
			"id": "%s",
			"device_id": "%s",
			"firmware_hash": "%s",
			"timestamp": "%s",
			"signals": %s
		}`, id, deviceID, firmware, timestamp, rawData)

		// Validate and migrate if needed
		obs, err := validator.ValidateAndMigrate([]byte(fullJSON))
		if err != nil {
			// Log invalid observation but continue (graceful degradation)
			fmt.Printf("Warning: invalid observation schema (id=%s): %v\n", id, err)
			continue
		}

		observations = append(observations, *obs)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return observations, nil
}
