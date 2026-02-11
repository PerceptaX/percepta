package percepta

import (
	"fmt"

	"github.com/perceptumx/percepta/internal/camera"
	"github.com/perceptumx/percepta/internal/core"
	"github.com/perceptumx/percepta/internal/storage"
	"github.com/perceptumx/percepta/internal/vision"
)

type Core struct {
	camera  core.CameraDriver
	vision  core.VisionDriver
	storage *storage.MemoryStorage
}

func NewCore(cameraPath string) (*Core, error) {
	// Initialize camera driver (Linux V4L2 for now)
	cameraDriver := camera.NewV4L2Camera(cameraPath)

	// Initialize vision driver
	visionDriver, err := vision.NewClaudeVision()
	if err != nil {
		return nil, fmt.Errorf("vision init failed: %w", err)
	}

	// Initialize in-memory storage
	storageDriver := storage.NewMemoryStorage()

	return &Core{
		camera:  cameraDriver,
		vision:  visionDriver,
		storage: storageDriver,
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

	// Save to storage
	if err := c.storage.Save(*obs); err != nil {
		return nil, fmt.Errorf("storage failed: %w", err)
	}

	return obs, nil
}

func (c *Core) ObservationCount() int {
	return c.storage.Count()
}
