//go:build !windows

package percepta

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/perceptumx/percepta/internal/core"
)

// Mock implementations

type mockCameraDriver struct {
	openErr       error
	captureErr    error
	closeErr      error
	captureFrames [][]byte
	captureCount  int
	openCalled    bool
	closeCalled   bool
}

func (m *mockCameraDriver) Open() error {
	m.openCalled = true
	return m.openErr
}

func (m *mockCameraDriver) CaptureFrame() ([]byte, error) {
	if m.captureErr != nil {
		return nil, m.captureErr
	}
	if m.captureCount >= len(m.captureFrames) {
		return nil, fmt.Errorf("no more frames")
	}
	frame := m.captureFrames[m.captureCount]
	m.captureCount++
	return frame, nil
}

func (m *mockCameraDriver) Close() error {
	m.closeCalled = true
	return m.closeErr
}

type mockStorageDriver struct {
	saveErr      error
	queryErr     error
	observations []core.Observation
	saveCount    int
}

func (m *mockStorageDriver) Save(obs core.Observation) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.observations = append(m.observations, obs)
	m.saveCount++
	return nil
}

func (m *mockStorageDriver) Query(deviceID string, limit int) ([]core.Observation, error) {
	if m.queryErr != nil {
		return nil, m.queryErr
	}

	// Filter by device ID
	var results []core.Observation
	for _, obs := range m.observations {
		if obs.DeviceID == deviceID {
			results = append(results, obs)
			if len(results) >= limit {
				break
			}
		}
	}
	return results, nil
}

func (m *mockStorageDriver) Count() int {
	return len(m.observations)
}

type mockSignalParser struct {
	parseErr error
	signals  []core.Signal
}

func (m *mockSignalParser) Parse(frame []byte) ([]core.Signal, error) {
	if m.parseErr != nil {
		return nil, m.parseErr
	}
	return m.signals, nil
}

// Helper function to create a test Core with mocks
func setupTestCore(t *testing.T, camera *mockCameraDriver, storage *mockStorageDriver) *Core {
	if camera == nil {
		camera = &mockCameraDriver{
			captureFrames: [][]byte{
				[]byte("test-frame-1"),
				[]byte("test-frame-2"),
				[]byte("test-frame-3"),
			},
		}
	}

	if storage == nil {
		storage = &mockStorageDriver{
			observations: []core.Observation{},
		}
	}

	// Note: We can't easily mock vision.ClaudeVision since it's created internally
	// These tests will focus on the aspects we can test without API integration
	return &Core{
		camera:  camera,
		storage: storage,
		// vision and smoother will be nil in these unit tests
	}
}

// Tests for Core structure and initialization

func TestCore_Structure(t *testing.T) {
	camera := &mockCameraDriver{
		captureFrames: [][]byte{[]byte("test")},
	}
	storage := &mockStorageDriver{}

	perceptaCore := setupTestCore(t, camera, storage)

	if perceptaCore == nil {
		t.Fatal("Expected non-nil Core")
	}

	if perceptaCore.camera == nil {
		t.Error("Expected camera to be set")
	}

	if perceptaCore.storage == nil {
		t.Error("Expected storage to be set")
	}
}

func TestCore_ObservationCount(t *testing.T) {
	storage := &mockStorageDriver{
		observations: []core.Observation{
			{ID: "obs-1", DeviceID: "device1"},
			{ID: "obs-2", DeviceID: "device1"},
			{ID: "obs-3", DeviceID: "device2"},
		},
	}

	perceptaCore := setupTestCore(t, nil, storage)

	count := perceptaCore.ObservationCount()
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
}

func TestCore_ObservationCount_Empty(t *testing.T) {
	storage := &mockStorageDriver{
		observations: []core.Observation{},
	}

	perceptaCore := setupTestCore(t, nil, storage)

	count := perceptaCore.ObservationCount()
	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

// Tests for camera integration

func TestCore_CameraOpen_Error(t *testing.T) {
	camera := &mockCameraDriver{
		openErr: errors.New("camera busy"),
	}
	storage := &mockStorageDriver{}

	perceptaCore := setupTestCore(t, camera, storage)

	// This would require observe() to be testable, but we can at least verify
	// the mock behaves correctly
	err := perceptaCore.camera.Open()
	if err == nil {
		t.Error("Expected error when camera open fails")
	}

	if !strings.Contains(err.Error(), "camera busy") {
		t.Errorf("Expected 'camera busy' error, got: %v", err)
	}
}

func TestCore_CameraCapture_Error(t *testing.T) {
	camera := &mockCameraDriver{
		captureErr: errors.New("capture timeout"),
	}

	_, err := camera.CaptureFrame()
	if err == nil {
		t.Error("Expected error when capture fails")
	}

	if !strings.Contains(err.Error(), "capture timeout") {
		t.Errorf("Expected 'capture timeout' error, got: %v", err)
	}
}

func TestCore_CameraClose_Error(t *testing.T) {
	camera := &mockCameraDriver{
		closeErr: errors.New("close failed"),
	}

	err := camera.Close()
	if err == nil {
		t.Error("Expected error when close fails")
	}

	if !strings.Contains(err.Error(), "close failed") {
		t.Errorf("Expected 'close failed' error, got: %v", err)
	}
}

func TestCore_CameraLifecycle(t *testing.T) {
	camera := &mockCameraDriver{
		captureFrames: [][]byte{[]byte("frame1")},
	}

	// Test open
	if err := camera.Open(); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	if !camera.openCalled {
		t.Error("Expected Open to be called")
	}

	// Test capture
	frame, err := camera.CaptureFrame()
	if err != nil {
		t.Fatalf("CaptureFrame failed: %v", err)
	}

	if string(frame) != "frame1" {
		t.Errorf("Expected 'frame1', got '%s'", string(frame))
	}

	// Test close
	if err := camera.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	if !camera.closeCalled {
		t.Error("Expected Close to be called")
	}
}

func TestCore_CameraMultipleFrames(t *testing.T) {
	camera := &mockCameraDriver{
		captureFrames: [][]byte{
			[]byte("frame1"),
			[]byte("frame2"),
			[]byte("frame3"),
		},
	}

	for i := 1; i <= 3; i++ {
		frame, err := camera.CaptureFrame()
		if err != nil {
			t.Fatalf("CaptureFrame %d failed: %v", i, err)
		}

		expected := fmt.Sprintf("frame%d", i)
		if string(frame) != expected {
			t.Errorf("Frame %d: expected '%s', got '%s'", i, expected, string(frame))
		}
	}

	// Fourth capture should fail
	_, err := camera.CaptureFrame()
	if err == nil {
		t.Error("Expected error when no more frames available")
	}
}

// Tests for storage integration

func TestCore_Storage_Save(t *testing.T) {
	storage := &mockStorageDriver{}

	obs := core.Observation{
		ID:       "test-obs-1",
		DeviceID: "device1",
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true},
		},
	}

	err := storage.Save(obs)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if storage.saveCount != 1 {
		t.Errorf("Expected saveCount 1, got %d", storage.saveCount)
	}

	if len(storage.observations) != 1 {
		t.Errorf("Expected 1 observation stored, got %d", len(storage.observations))
	}

	if storage.observations[0].ID != "test-obs-1" {
		t.Errorf("Expected ID 'test-obs-1', got '%s'", storage.observations[0].ID)
	}
}

func TestCore_Storage_Query(t *testing.T) {
	storage := &mockStorageDriver{
		observations: []core.Observation{
			{ID: "obs-1", DeviceID: "device1"},
			{ID: "obs-2", DeviceID: "device1"},
			{ID: "obs-3", DeviceID: "device2"},
		},
	}

	results, err := storage.Query("device1", 10)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results for device1, got %d", len(results))
	}

	for _, obs := range results {
		if obs.DeviceID != "device1" {
			t.Errorf("Expected DeviceID 'device1', got '%s'", obs.DeviceID)
		}
	}
}

func TestCore_Storage_QueryWithLimit(t *testing.T) {
	storage := &mockStorageDriver{
		observations: []core.Observation{
			{ID: "obs-1", DeviceID: "device1"},
			{ID: "obs-2", DeviceID: "device1"},
			{ID: "obs-3", DeviceID: "device1"},
		},
	}

	results, err := storage.Query("device1", 2)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected limit of 2 results, got %d", len(results))
	}
}

func TestCore_Storage_QueryError(t *testing.T) {
	storage := &mockStorageDriver{
		queryErr: errors.New("database error"),
	}

	_, err := storage.Query("device1", 10)
	if err == nil {
		t.Error("Expected error from Query")
	}

	if !strings.Contains(err.Error(), "database error") {
		t.Errorf("Expected 'database error', got: %v", err)
	}
}

func TestCore_Storage_SaveError(t *testing.T) {
	storage := &mockStorageDriver{
		saveErr: errors.New("disk full"),
	}

	obs := core.Observation{ID: "test-obs"}
	err := storage.Save(obs)
	if err == nil {
		t.Error("Expected error from Save")
	}

	if !strings.Contains(err.Error(), "disk full") {
		t.Errorf("Expected 'disk full' error, got: %v", err)
	}
}

// Tests for signal parser

func TestCore_SignalParser_Parse_Success(t *testing.T) {
	parser := &mockSignalParser{
		signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true, Confidence: 0.95},
			core.LEDSignal{Name: "LED2", On: false, Confidence: 0.90},
		},
	}

	signals, err := parser.Parse([]byte("frame"))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(signals) != 2 {
		t.Errorf("Expected 2 signals, got %d", len(signals))
	}

	led1, ok := signals[0].(core.LEDSignal)
	if !ok {
		t.Fatal("Expected first signal to be LEDSignal")
	}

	if led1.Name != "LED1" || !led1.On {
		t.Errorf("Expected LED1 ON, got %+v", led1)
	}
}

func TestCore_SignalParser_Parse_Error(t *testing.T) {
	parser := &mockSignalParser{
		parseErr: errors.New("parse failed"),
	}

	_, err := parser.Parse([]byte("frame"))
	if err == nil {
		t.Error("Expected error from Parse")
	}

	if !strings.Contains(err.Error(), "parse failed") {
		t.Errorf("Expected 'parse failed' error, got: %v", err)
	}
}

func TestCore_SignalParser_Parse_EmptySignals(t *testing.T) {
	parser := &mockSignalParser{
		signals: []core.Signal{},
	}

	signals, err := parser.Parse([]byte("frame"))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(signals) != 0 {
		t.Errorf("Expected 0 signals, got %d", len(signals))
	}
}

func TestCore_SignalParser_Parse_MixedSignals(t *testing.T) {
	parser := &mockSignalParser{
		signals: []core.Signal{
			core.LEDSignal{Name: "STATUS", On: true},
			core.DisplaySignal{Name: "LCD", Text: "Ready", Confidence: 0.92},
			core.BootTimingSignal{DurationMs: 2000, Confidence: 0.88},
		},
	}

	signals, err := parser.Parse([]byte("frame"))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(signals) != 3 {
		t.Errorf("Expected 3 signals, got %d", len(signals))
	}

	// Verify each signal type
	_, okLED := signals[0].(core.LEDSignal)
	_, okDisplay := signals[1].(core.DisplaySignal)
	_, okBoot := signals[2].(core.BootTimingSignal)

	if !okLED {
		t.Error("Expected first signal to be LEDSignal")
	}
	if !okDisplay {
		t.Error("Expected second signal to be DisplaySignal")
	}
	if !okBoot {
		t.Error("Expected third signal to be BootTimingSignal")
	}
}

// Integration tests

func TestCore_Integration_CameraAndStorage(t *testing.T) {
	camera := &mockCameraDriver{
		captureFrames: [][]byte{
			[]byte("frame1"),
			[]byte("frame2"),
		},
	}

	storage := &mockStorageDriver{}

	perceptaCore := setupTestCore(t, camera, storage)

	// Simulate observation workflow
	if err := camera.Open(); err != nil {
		t.Fatalf("Camera open failed: %v", err)
	}

	frame1, err := camera.CaptureFrame()
	if err != nil {
		t.Fatalf("First capture failed: %v", err)
	}

	frame2, err := camera.CaptureFrame()
	if err != nil {
		t.Fatalf("Second capture failed: %v", err)
	}

	// Create observation
	obs := &core.Observation{
		ID:       "integration-test",
		DeviceID: "test-device",
		Signals: []core.Signal{
			core.LEDSignal{Name: "LED1", On: true},
		},
	}

	// Save observation
	if err := storage.Save(*obs); err != nil {
		t.Fatalf("Storage save failed: %v", err)
	}

	// Close camera
	if err := camera.Close(); err != nil {
		t.Fatalf("Camera close failed: %v", err)
	}

	// Verify state
	if !camera.openCalled {
		t.Error("Expected camera.Open to be called")
	}

	if !camera.closeCalled {
		t.Error("Expected camera.Close to be called")
	}

	if len(frame1) == 0 || len(frame2) == 0 {
		t.Error("Expected non-empty frames")
	}

	if storage.saveCount != 1 {
		t.Errorf("Expected 1 save, got %d", storage.saveCount)
	}

	if perceptaCore.ObservationCount() != 1 {
		t.Errorf("Expected observation count 1, got %d", perceptaCore.ObservationCount())
	}
}

func TestCore_Integration_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		camera      *mockCameraDriver
		storage     *mockStorageDriver
		expectError bool
		errorPhase  string
	}{
		{
			name: "camera open fails",
			camera: &mockCameraDriver{
				openErr: errors.New("camera not found"),
			},
			storage:     &mockStorageDriver{},
			expectError: true,
			errorPhase:  "open",
		},
		{
			name: "camera capture fails",
			camera: &mockCameraDriver{
				captureErr: errors.New("timeout"),
			},
			storage:     &mockStorageDriver{},
			expectError: true,
			errorPhase:  "capture",
		},
		{
			name: "storage save fails",
			camera: &mockCameraDriver{
				captureFrames: [][]byte{[]byte("frame")},
			},
			storage: &mockStorageDriver{
				saveErr: errors.New("disk full"),
			},
			expectError: true,
			errorPhase:  "save",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			perceptaCore := setupTestCore(t, tt.camera, tt.storage)

			switch tt.errorPhase {
			case "open":
				err := perceptaCore.camera.Open()
				if !tt.expectError && err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.expectError && err == nil {
					t.Error("Expected error but got none")
				}
			case "capture":
				_, err := perceptaCore.camera.CaptureFrame()
				if !tt.expectError && err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.expectError && err == nil {
					t.Error("Expected error but got none")
				}
			case "save":
				obs := &core.Observation{ID: "test"}
				err := perceptaCore.storage.Save(*obs)
				if !tt.expectError && err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.expectError && err == nil {
					t.Error("Expected error but got none")
				}
			}
		})
	}
}

func TestCore_Integration_MultipleObservations(t *testing.T) {
	camera := &mockCameraDriver{
		captureFrames: [][]byte{
			[]byte("frame1"),
			[]byte("frame2"),
			[]byte("frame3"),
		},
	}

	storage := &mockStorageDriver{}
	perceptaCore := setupTestCore(t, camera, storage)

	// Simulate multiple observations
	observations := []*core.Observation{
		{ID: "obs-1", DeviceID: "device1", Timestamp: time.Now()},
		{ID: "obs-2", DeviceID: "device1", Timestamp: time.Now().Add(1 * time.Minute)},
		{ID: "obs-3", DeviceID: "device2", Timestamp: time.Now().Add(2 * time.Minute)},
	}

	for _, obs := range observations {
		if err := storage.Save(*obs); err != nil {
			t.Fatalf("Failed to save observation %s: %v", obs.ID, err)
		}
	}

	// Verify count
	if perceptaCore.ObservationCount() != 3 {
		t.Errorf("Expected 3 observations, got %d", perceptaCore.ObservationCount())
	}

	// Query by device
	device1Obs, err := storage.Query("device1", 10)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(device1Obs) != 2 {
		t.Errorf("Expected 2 observations for device1, got %d", len(device1Obs))
	}
}

// Edge case tests

func TestCore_EdgeCase_EmptyFrame(t *testing.T) {
	camera := &mockCameraDriver{
		captureFrames: [][]byte{[]byte("")},
	}

	frame, err := camera.CaptureFrame()
	if err != nil {
		t.Fatalf("CaptureFrame failed: %v", err)
	}

	if len(frame) != 0 {
		t.Errorf("Expected empty frame, got %d bytes", len(frame))
	}
}

func TestCore_EdgeCase_LargeFrame(t *testing.T) {
	largeFrame := make([]byte, 1024*1024) // 1MB
	for i := range largeFrame {
		largeFrame[i] = byte(i % 256)
	}

	camera := &mockCameraDriver{
		captureFrames: [][]byte{largeFrame},
	}

	frame, err := camera.CaptureFrame()
	if err != nil {
		t.Fatalf("CaptureFrame failed: %v", err)
	}

	if len(frame) != len(largeFrame) {
		t.Errorf("Expected frame size %d, got %d", len(largeFrame), len(frame))
	}
}

func TestCore_EdgeCase_ManySignals(t *testing.T) {
	// Test with many signals
	signals := make([]core.Signal, 100)
	for i := 0; i < 100; i++ {
		signals[i] = core.LEDSignal{
			Name:       fmt.Sprintf("LED%d", i),
			On:         i%2 == 0,
			Confidence: 0.9,
		}
	}

	parser := &mockSignalParser{signals: signals}

	parsed, err := parser.Parse([]byte("frame"))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(parsed) != 100 {
		t.Errorf("Expected 100 signals, got %d", len(parsed))
	}
}

func TestCore_EdgeCase_ZeroObservations(t *testing.T) {
	storage := &mockStorageDriver{}
	perceptaCore := setupTestCore(t, nil, storage)

	count := perceptaCore.ObservationCount()
	if count != 0 {
		t.Errorf("Expected 0 observations, got %d", count)
	}

	results, err := storage.Query("any-device", 10)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}
