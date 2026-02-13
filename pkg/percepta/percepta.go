//go:build linux || darwin

package percepta

import (
	"fmt"
	"time"

	"github.com/perceptumx/percepta/internal/camera"
	"github.com/perceptumx/percepta/internal/core"
	"github.com/perceptumx/percepta/internal/filter"
	"github.com/perceptumx/percepta/internal/vision"
)

type Core struct {
	camera   core.CameraDriver
	vision   *vision.ClaudeVision
	storage  core.StorageDriver
	smoother *filter.TemporalSmoother
}

func NewCore(cameraPath string, storage core.StorageDriver) (*Core, error) {
	// Initialize camera driver (platform-specific)
	cameraDriver := camera.NewCamera(cameraPath)

	// Initialize vision driver
	visionDriver, err := vision.NewClaudeVision()
	if err != nil {
		return nil, fmt.Errorf("vision init failed: %w", err)
	}

	return &Core{
		camera:   cameraDriver,
		vision:   visionDriver,
		storage:  storage,
		smoother: filter.NewTemporalSmoother(storage),
	}, nil
}

func (c *Core) Observe(deviceID string) (*core.Observation, error) {
	return c.observe(deviceID, 0, 0)
}

func (c *Core) ObserveWithOptions(deviceID string, frameCount int, interval time.Duration) (*core.Observation, error) {
	return c.observe(deviceID, frameCount, interval)
}

func (c *Core) observe(deviceID string, frameCount int, interval time.Duration) (*core.Observation, error) {
	// Open camera
	if err := c.camera.Open(); err != nil {
		return nil, fmt.Errorf("camera open failed: %w", err)
	}
	defer c.camera.Close()

	// Multi-frame capture for complete LED detection (fixes ISS-001)
	var multiFrame *vision.MultiFrameCapture
	if frameCount > 0 && interval > 0 {
		multiFrame = vision.NewMultiFrameCaptureWithOptions(c.camera, c.vision.GetParser(), frameCount, interval)
	} else {
		multiFrame = vision.NewMultiFrameCapture(c.camera, c.vision.GetParser())
	}
	frames, err := multiFrame.Capture()
	if err != nil {
		return nil, fmt.Errorf("multi-frame capture failed: %w", err)
	}

	if len(frames) == 0 {
		return nil, fmt.Errorf("no frames captured")
	}

	// Aggregate LED detections across frames
	leds := vision.AggregateLEDs(frames)

	// Aggregate display detections across frames (tracks text changes)
	aggregatedDisplays := vision.AggregateDisplays(frames)

	// Combine signals
	var signals []core.Signal
	for _, led := range leds {
		signals = append(signals, led)
	}
	for _, display := range aggregatedDisplays {
		signals = append(signals, display)
	}

	obs := &core.Observation{
		SchemaVersion: core.CurrentSchemaVersion,
		ID:            core.GenerateID(),
		DeviceID:      deviceID,
		Timestamp:     time.Now(),
		Signals:       signals,
	}

	// Apply temporal smoothing before returning
	smoothedObs, err := c.smoother.Smooth(obs)
	if err != nil {
		// Log but don't fail observation (graceful degradation)
		return obs, nil
	}

	return smoothedObs, nil
}

func (c *Core) ObservationCount() int {
	return c.storage.Count()
}
