package percepta

import (
	"fmt"

	"github.com/perceptumx/percepta/internal/camera"
	"github.com/perceptumx/percepta/internal/core"
	"github.com/perceptumx/percepta/internal/vision"
)

type Core struct {
	camera  core.CameraDriver
	vision  core.VisionDriver
	storage core.StorageDriver
}

func NewCore(cameraPath string, storage core.StorageDriver) (*Core, error) {
	// Initialize camera driver (Linux V4L2 for now)
	cameraDriver := camera.NewV4L2Camera(cameraPath)

	// Initialize vision driver
	visionDriver, err := vision.NewClaudeVision()
	if err != nil {
		return nil, fmt.Errorf("vision init failed: %w", err)
	}

	return &Core{
		camera:  cameraDriver,
		vision:  visionDriver,
		storage: storage,
	}, nil
}

func (c *Core) Observe(deviceID string) (*core.Observation, error) {
	// Open camera
	if err := c.camera.Open(); err != nil {
		return nil, fmt.Errorf("camera open failed: %w", err)
	}
	defer c.camera.Close()

	// Capture frame
	frame, err := c.camera.CaptureFrame()
	if err != nil {
		return nil, fmt.Errorf("camera capture failed: %w", err)
	}

	// Analyze frame with vision
	obs, err := c.vision.Observe(deviceID, frame)
	if err != nil {
		return nil, fmt.Errorf("vision analysis failed: %w", err)
	}

	// Note: firmware tag injection and storage save happens in cmd layer
	// This keeps the core framework-agnostic

	return obs, nil
}

func (c *Core) ObservationCount() int {
	return c.storage.Count()
}
