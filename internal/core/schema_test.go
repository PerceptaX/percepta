package core

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSchemaValidator_CurrentVersion(t *testing.T) {
	validator := NewSchemaValidator()

	// Create observation with current schema version
	obs := Observation{
		SchemaVersion: CurrentSchemaVersion,
		ID:            "test-1",
		DeviceID:      "test-device",
		Timestamp:     time.Now(),
		Signals: []Signal{
			LEDSignal{Name: "LED1", On: true, Confidence: 0.95},
		},
	}

	// Serialize
	data, err := json.Marshal(obs)
	if err != nil {
		t.Fatalf("Failed to marshal observation: %v", err)
	}

	// Validate (should succeed)
	validated, err := validator.ValidateAndMigrate(data)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if validated.SchemaVersion != CurrentSchemaVersion {
		t.Errorf("Expected schema version %s, got %s", CurrentSchemaVersion, validated.SchemaVersion)
	}
	if validated.ID != "test-1" {
		t.Errorf("Expected ID 'test-1', got '%s'", validated.ID)
	}
}

func TestSchemaValidator_LegacyObservation(t *testing.T) {
	validator := NewSchemaValidator()

	// Create legacy observation (no schema_version field)
	legacyJSON := `{
		"id": "legacy-1",
		"device_id": "legacy-device",
		"timestamp": "2026-02-13T00:00:00Z",
		"signals": [
			{
				"name": "LED1",
				"on": true,
				"confidence": 0.95
			}
		]
	}`

	// Validate (should inject current version)
	validated, err := validator.ValidateAndMigrate([]byte(legacyJSON))
	if err != nil {
		t.Fatalf("Legacy observation validation failed: %v", err)
	}

	if validated.SchemaVersion != CurrentSchemaVersion {
		t.Errorf("Expected schema version %s injected, got %s", CurrentSchemaVersion, validated.SchemaVersion)
	}
	if validated.ID != "legacy-1" {
		t.Errorf("Expected ID 'legacy-1', got '%s'", validated.ID)
	}
}

func TestSchemaValidator_UnsupportedVersion(t *testing.T) {
	validator := NewSchemaValidator()

	// Create observation with future unsupported version
	futureJSON := `{
		"schema_version": "2.0.0",
		"id": "future-1",
		"device_id": "future-device",
		"timestamp": "2026-02-13T00:00:00Z",
		"signals": []
	}`

	// Validate (should fail - no migration available)
	_, err := validator.ValidateAndMigrate([]byte(futureJSON))
	if err == nil {
		t.Fatal("Expected validation error for unsupported schema version, got nil")
	}

	expectedError := "unsupported schema version"
	if !contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing '%s', got '%s'", expectedError, err.Error())
	}
}

func TestSchemaValidator_WithMigration(t *testing.T) {
	validator := NewSchemaValidator()

	// Register a test migration from 0.9.0 to 1.0.0
	validator.migrations["0.9.0->1.0.0"] = func(data map[string]interface{}) (map[string]interface{}, error) {
		// Example migration: rename "leds" field to "signals"
		if leds, ok := data["leds"]; ok {
			data["signals"] = leds
			delete(data, "leds")
		}
		data["schema_version"] = "1.0.0"
		return data, nil
	}

	// Create observation with old schema
	oldJSON := `{
		"schema_version": "0.9.0",
		"id": "old-1",
		"device_id": "old-device",
		"timestamp": "2026-02-13T00:00:00Z",
		"leds": [
			{
				"name": "LED1",
				"on": true,
				"confidence": 0.95
			}
		]
	}`

	// Validate (should migrate successfully)
	validated, err := validator.ValidateAndMigrate([]byte(oldJSON))
	if err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	if validated.SchemaVersion != "1.0.0" {
		t.Errorf("Expected schema version 1.0.0 after migration, got %s", validated.SchemaVersion)
	}
	if validated.ID != "old-1" {
		t.Errorf("Expected ID 'old-1', got '%s'", validated.ID)
	}
	if len(validated.Signals) != 1 {
		t.Errorf("Expected 1 signal after migration, got %d", len(validated.Signals))
	}
}

func TestSchemaValidator_InvalidJSON(t *testing.T) {
	validator := NewSchemaValidator()

	// Invalid JSON
	invalidJSON := `{invalid json`

	// Validate (should fail)
	_, err := validator.ValidateAndMigrate([]byte(invalidJSON))
	if err == nil {
		t.Fatal("Expected validation error for invalid JSON, got nil")
	}

	expectedError := "invalid JSON"
	if !contains(err.Error(), expectedError) {
		t.Errorf("Expected error containing '%s', got '%s'", expectedError, err.Error())
	}
}

func TestEnsureSchemaVersion(t *testing.T) {
	// Test observation without schema version
	obs := &Observation{
		ID:       "test-1",
		DeviceID: "test-device",
	}

	EnsureSchemaVersion(obs)

	if obs.SchemaVersion != CurrentSchemaVersion {
		t.Errorf("Expected schema version %s, got %s", CurrentSchemaVersion, obs.SchemaVersion)
	}

	// Test observation with existing schema version (should not change)
	obs2 := &Observation{
		SchemaVersion: "custom-version",
		ID:            "test-2",
		DeviceID:      "test-device",
	}

	EnsureSchemaVersion(obs2)

	if obs2.SchemaVersion != "custom-version" {
		t.Errorf("Expected schema version unchanged at 'custom-version', got %s", obs2.SchemaVersion)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
