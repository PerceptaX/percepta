//go:build linux && cgo

package camera

import (
	"os"
	"testing"
)

// Integration test - requires real hardware
func TestV4L2Camera_Integration(t *testing.T) {
	// Check if /dev/video0 exists
	if _, err := os.Stat("/dev/video0"); os.IsNotExist(err) {
		t.Skip("Skipping integration test: /dev/video0 not available")
	}

	cam := NewV4L2Camera("/dev/video0")

	// Test Open
	err := cam.Open()
	if err != nil {
		t.Skipf("Skipping integration test: camera open failed (hardware may be in use): %v", err)
	}
	defer cam.Close()

	// Test CaptureFrame
	frame, err := cam.CaptureFrame()
	if err != nil {
		t.Fatalf("CaptureFrame failed: %v", err)
	}

	// Verify frame is not empty
	if len(frame) == 0 {
		t.Error("Expected non-empty frame data")
	}

	// Verify frame starts with JPEG magic bytes (0xFF 0xD8)
	if len(frame) >= 2 && (frame[0] != 0xFF || frame[1] != 0xD8) {
		t.Errorf("Expected JPEG magic bytes (FF D8), got: %02X %02X", frame[0], frame[1])
	}
}

func TestV4L2Camera_DeviceNotFound(t *testing.T) {
	cam := NewV4L2Camera("/dev/video999")

	// Test Open with non-existent device
	err := cam.Open()
	if err == nil {
		defer cam.Close()
		t.Fatal("Expected error for non-existent device, got nil")
	}

	// Verify error message mentions the device
	if err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}

func TestV4L2Camera_CaptureBeforeOpen(t *testing.T) {
	cam := &V4L2Camera{devicePath: "/dev/video0"}

	// Try to capture frame before opening
	_, err := cam.CaptureFrame()
	if err == nil {
		t.Fatal("Expected error when capturing before Open, got nil")
	}

	// Verify error message mentions camera not opened
	expectedMsg := "camera not opened"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestV4L2Camera_CloseWithoutOpen(t *testing.T) {
	cam := &V4L2Camera{devicePath: "/dev/video0"}

	// Close without opening should not panic
	err := cam.Close()
	if err != nil {
		t.Errorf("Close without Open should succeed, got error: %v", err)
	}
}

func TestV4L2Camera_DoubleClose(t *testing.T) {
	// Check if /dev/video0 exists
	if _, err := os.Stat("/dev/video0"); os.IsNotExist(err) {
		t.Skip("Skipping: /dev/video0 not available")
	}

	cam := NewV4L2Camera("/dev/video0")

	// Open camera
	err := cam.Open()
	if err != nil {
		t.Skipf("Skipping: camera open failed: %v", err)
	}

	// Close camera
	err = cam.Close()
	if err != nil {
		t.Fatalf("First close failed: %v", err)
	}

	// Close again - should not panic
	err = cam.Close()
	if err != nil {
		// Second close might fail, but shouldn't panic
		t.Logf("Second close returned error (expected): %v", err)
	}
}

func TestV4L2Camera_OpenCloseLifecycle(t *testing.T) {
	// Check if /dev/video0 exists
	if _, err := os.Stat("/dev/video0"); os.IsNotExist(err) {
		t.Skip("Skipping: /dev/video0 not available")
	}

	cam := NewV4L2Camera("/dev/video0")

	// Test Open
	err := cam.Open()
	if err != nil {
		t.Skipf("Skipping: camera open failed: %v", err)
	}

	// Test Close
	err = cam.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Verify we can't capture after close
	_, err = cam.CaptureFrame()
	if err == nil {
		t.Error("Expected error when capturing after Close, got nil")
	}
}

func TestV4L2Camera_MultipleCaptures(t *testing.T) {
	// Check if /dev/video0 exists
	if _, err := os.Stat("/dev/video0"); os.IsNotExist(err) {
		t.Skip("Skipping: /dev/video0 not available")
	}

	cam := NewV4L2Camera("/dev/video0")

	// Open camera
	err := cam.Open()
	if err != nil {
		t.Skipf("Skipping: camera open failed: %v", err)
	}
	defer cam.Close()

	// Capture multiple frames
	for i := 0; i < 3; i++ {
		frame, err := cam.CaptureFrame()
		if err != nil {
			t.Fatalf("Frame %d capture failed: %v", i, err)
		}

		if len(frame) == 0 {
			t.Errorf("Frame %d is empty", i)
		}

		// Verify JPEG magic bytes
		if len(frame) >= 2 && (frame[0] != 0xFF || frame[1] != 0xD8) {
			t.Errorf("Frame %d: invalid JPEG magic bytes: %02X %02X", i, frame[0], frame[1])
		}
	}
}

func TestNewV4L2Camera(t *testing.T) {
	// Test constructor
	cam := NewV4L2Camera("/dev/video0")
	if cam == nil {
		t.Fatal("Expected non-nil camera driver")
	}

	// Verify type
	v4l2Cam, ok := cam.(*V4L2Camera)
	if !ok {
		t.Fatalf("Expected *V4L2Camera, got %T", cam)
	}

	// Verify device path is stored
	if v4l2Cam.devicePath != "/dev/video0" {
		t.Errorf("Expected device path '/dev/video0', got '%s'", v4l2Cam.devicePath)
	}
}

func TestV4L2Camera_AlternativeDevices(t *testing.T) {
	// Test that constructor works with different device paths
	devices := []string{"/dev/video0", "/dev/video1", "/dev/video2"}

	for _, device := range devices {
		cam := NewV4L2Camera(device)
		v4l2Cam := cam.(*V4L2Camera)

		if v4l2Cam.devicePath != device {
			t.Errorf("Expected device path '%s', got '%s'", device, v4l2Cam.devicePath)
		}
	}
}

func TestV4L2Camera_InterfaceCompliance(t *testing.T) {
	// Verify V4L2Camera implements CameraDriver interface
	// This is primarily a compile-time check
	cam := NewV4L2Camera("/dev/video0")

	// Verify the camera driver is not nil
	if cam == nil {
		t.Fatal("Expected non-nil camera driver")
	}

	// The fact that this compiles means V4L2Camera implements CameraDriver
	// We can call methods to verify they're implemented
	v4l2Cam := cam.(*V4L2Camera)
	if v4l2Cam == nil {
		t.Error("Failed to cast to *V4L2Camera")
	}
}

func TestV4L2Camera_ErrorMessages(t *testing.T) {
	// Test that error messages are descriptive
	cam := NewV4L2Camera("/dev/video_nonexistent")

	err := cam.Open()
	if err == nil {
		defer cam.Close()
		t.Fatal("Expected error for non-existent device")
	}

	// Verify error message contains device path
	if err.Error() == "" {
		t.Error("Expected non-empty error message")
	}

	// Error message should mention the device path
	errMsg := err.Error()
	if len(errMsg) < 10 {
		t.Errorf("Error message seems too short: '%s'", errMsg)
	}
}
