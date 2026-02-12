package core

import (
	"encoding/json"
	"fmt"
	"time"
)

// SchemaValidator validates observation schema versions
type SchemaValidator struct {
	currentVersion string
	migrations     map[string]MigrationFunc
}

// MigrationFunc upgrades observation from old version to new version
type MigrationFunc func(data map[string]interface{}) (map[string]interface{}, error)

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator() *SchemaValidator {
	return &SchemaValidator{
		currentVersion: CurrentSchemaVersion,
		migrations:     map[string]MigrationFunc{
			// Future migrations go here
			// "0.9.0->1.0.0": migrateV0ToV1,
		},
	}
}

// ValidateAndMigrate checks observation schema version and migrates if needed
func (v *SchemaValidator) ValidateAndMigrate(data []byte) (*Observation, error) {
	// Parse JSON to check schema version
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Check schema version
	schemaVer, ok := raw["schema_version"].(string)
	if !ok || schemaVer == "" {
		// Missing schema_version → assume legacy format (pre-v1.0.0)
		// Inject current version and continue
		raw["schema_version"] = v.currentVersion
	} else if schemaVer != v.currentVersion {
		// Old schema → check if migration exists
		migrationKey := schemaVer + "->" + v.currentVersion
		if migration, exists := v.migrations[migrationKey]; exists {
			migrated, err := migration(raw)
			if err != nil {
				return nil, fmt.Errorf("migration failed: %w", err)
			}
			raw = migrated
		} else {
			return nil, fmt.Errorf("unsupported schema version: %s (current: %s, no migration available)", schemaVer, v.currentVersion)
		}
	}

	// Build Observation manually to handle Signal interface deserialization
	obs := &Observation{
		SchemaVersion: raw["schema_version"].(string),
	}

	// Parse basic fields
	if id, ok := raw["id"].(string); ok {
		obs.ID = id
	}
	if deviceID, ok := raw["device_id"].(string); ok {
		obs.DeviceID = deviceID
	}
	if firmwareHash, ok := raw["firmware_hash"].(string); ok {
		obs.FirmwareHash = firmwareHash
	}
	if timestamp, ok := raw["timestamp"].(string); ok {
		t, err := parseTimestamp(timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}
		obs.Timestamp = t
	}

	// Parse signals (handle Signal interface)
	if signals, ok := raw["signals"].([]interface{}); ok {
		parsedSignals, err := parseSignals(signals)
		if err != nil {
			return nil, fmt.Errorf("failed to parse signals: %w", err)
		}
		obs.Signals = parsedSignals
	}

	return obs, nil
}

// parseTimestamp parses a timestamp string in various formats
func parseTimestamp(ts string) (time.Time, error) {
	// Try RFC3339 first (standard format)
	t, err := time.Parse(time.RFC3339, ts)
	if err == nil {
		return t, nil
	}

	// Try RFC3339Nano as fallback
	t, err = time.Parse(time.RFC3339Nano, ts)
	if err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unsupported timestamp format: %s", ts)
}

// parseSignals deserializes signals from generic JSON data
func parseSignals(rawSignals []interface{}) ([]Signal, error) {
	signals := make([]Signal, 0, len(rawSignals))

	for _, raw := range rawSignals {
		signalMap, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		// Determine signal type
		signalType, ok := signalMap["type"].(string)
		if !ok {
			// Infer type from fields
			if _, hasOn := signalMap["on"]; hasOn {
				signalType = "led"
			} else if _, hasText := signalMap["text"]; hasText {
				signalType = "display"
			} else if _, hasDuration := signalMap["duration_ms"]; hasDuration {
				signalType = "boot_timing"
			}
		}

		// Re-marshal and unmarshal to correct type
		rawBytes, err := json.Marshal(signalMap)
		if err != nil {
			return nil, fmt.Errorf("failed to re-marshal signal: %w", err)
		}

		switch signalType {
		case "led":
			var led LEDSignal
			if err := json.Unmarshal(rawBytes, &led); err != nil {
				return nil, fmt.Errorf("failed to unmarshal LED signal: %w", err)
			}
			signals = append(signals, led)
		case "display":
			var display DisplaySignal
			if err := json.Unmarshal(rawBytes, &display); err != nil {
				return nil, fmt.Errorf("failed to unmarshal display signal: %w", err)
			}
			signals = append(signals, display)
		case "boot_timing":
			var boot BootTimingSignal
			if err := json.Unmarshal(rawBytes, &boot); err != nil {
				return nil, fmt.Errorf("failed to unmarshal boot timing signal: %w", err)
			}
			signals = append(signals, boot)
		}
	}

	return signals, nil
}

// EnsureSchemaVersion sets schema version on new observations
func EnsureSchemaVersion(obs *Observation) {
	if obs.SchemaVersion == "" {
		obs.SchemaVersion = CurrentSchemaVersion
	}
}
