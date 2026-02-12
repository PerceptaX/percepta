//go:build !windows

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
	// Initialize camera driver (Linux V4L2 for now)
	cameraDriver := camera.NewV4L2Camera(cameraPath)

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
	// Open camera
	if err := c.camera.Open(); err != nil {
		return nil, fmt.Errorf("camera open failed: %w", err)
	}
	defer c.camera.Close()

	// Multi-frame capture for complete LED detection (fixes ISS-001)
	multiFrame := vision.NewMultiFrameCapture(c.camera, c.vision.GetParser())
	frames, err := multiFrame.Capture()
	if err != nil {
		return nil, fmt.Errorf("multi-frame capture failed: %w", err)
	}

	if len(frames) == 0 {
		return nil, fmt.Errorf("no frames captured")
	}

	// Aggregate LED detections across frames
	leds := vision.AggregateLEDs(frames)

	// Get display signals from most recent frame (displays don't need aggregation)
	displays := getDisplaySignals(frames[len(frames)-1].Signals)

	// Combine signals
	var signals []core.Signal
	for _, led := range leds {
		signals = append(signals, led)
	}
	signals = append(signals, displays...)

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

func getDisplaySignals(signals []core.Signal) []core.Signal {
	var displays []core.Signal
	for _, signal := range signals {
		if _, ok := signal.(core.DisplaySignal); ok {
			displays = append(displays, signal)
		}
	}
	return displays
}

func (c *Core) ObservationCount() int {
	return c.storage.Count()
}
