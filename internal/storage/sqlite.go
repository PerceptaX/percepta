package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

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

	var obs core.Observation
	var signalsJSON string
	var timestamp string

	err := row.Scan(&obs.ID, &obs.DeviceID, &obs.FirmwareHash, &timestamp, &signalsJSON)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no observation found for device %s with firmware %s", deviceID, firmware)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan observation: %w", err)
	}

	// Parse timestamp
	obs.Timestamp, err = time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	// Deserialize signals
	if err := s.unmarshalSignals(signalsJSON, &obs.Signals); err != nil {
		return nil, err
	}

	return &obs, nil
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

	for rows.Next() {
		var obs core.Observation
		var signalsJSON string
		var timestamp string

		err := rows.Scan(&obs.ID, &obs.DeviceID, &obs.FirmwareHash, &timestamp, &signalsJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Parse timestamp
		obs.Timestamp, err = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}

		// Deserialize signals
		if err := s.unmarshalSignals(signalsJSON, &obs.Signals); err != nil {
			return nil, err
		}

		observations = append(observations, obs)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return observations, nil
}

// unmarshalSignals deserializes JSON signals into the correct concrete types
func (s *SQLiteStorage) unmarshalSignals(signalsJSON string, signals *[]core.Signal) error {
	// First unmarshal to generic structure to determine types
	var rawSignals []map[string]interface{}
	if err := json.Unmarshal([]byte(signalsJSON), &rawSignals); err != nil {
		return fmt.Errorf("failed to unmarshal signals: %w", err)
	}

	*signals = make([]core.Signal, 0, len(rawSignals))

	// Re-marshal each signal and unmarshal to correct type
	for _, raw := range rawSignals {
		signalType, ok := raw["type"].(string)
		if !ok {
			// Try to infer type from fields
			if _, hasOn := raw["on"]; hasOn {
				signalType = "led"
			} else if _, hasText := raw["text"]; hasText {
				signalType = "display"
			} else if _, hasDuration := raw["duration_ms"]; hasDuration {
				signalType = "boot_timing"
			}
		}

		rawBytes, err := json.Marshal(raw)
		if err != nil {
			return fmt.Errorf("failed to re-marshal signal: %w", err)
		}

		switch signalType {
		case "led":
			var led core.LEDSignal
			if err := json.Unmarshal(rawBytes, &led); err != nil {
				return fmt.Errorf("failed to unmarshal LED signal: %w", err)
			}
			*signals = append(*signals, led)
		case "display":
			var display core.DisplaySignal
			if err := json.Unmarshal(rawBytes, &display); err != nil {
				return fmt.Errorf("failed to unmarshal display signal: %w", err)
			}
			*signals = append(*signals, display)
		case "boot_timing":
			var boot core.BootTimingSignal
			if err := json.Unmarshal(rawBytes, &boot); err != nil {
				return fmt.Errorf("failed to unmarshal boot timing signal: %w", err)
			}
			*signals = append(*signals, boot)
		}
	}

	return nil
}
