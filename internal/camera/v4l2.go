//go:build linux

package camera

import (
	"fmt"

	"github.com/blackjack/webcam"
	"github.com/perceptumx/percepta/internal/core"
)

// V4L2Camera implements core.CameraDriver for Linux V4L2 devices
type V4L2Camera struct {
	devicePath string
	cam        *webcam.Webcam
}

// NewV4L2Camera creates a new Linux V4L2 camera driver
func NewV4L2Camera(devicePath string) core.CameraDriver {
	return &V4L2Camera{devicePath: devicePath}
}

func (c *V4L2Camera) Open() error {
	cam, err := webcam.Open(c.devicePath)
	if err != nil {
		return fmt.Errorf("failed to open camera %s: %w", c.devicePath, err)
	}
	c.cam = cam

	// Set format to MJPEG 1280x720 (balance quality vs size)
	// Get supported formats and find MJPEG
	formatDesc := c.cam.GetSupportedFormats()
	var mjpegFormat webcam.PixelFormat
	for format := range formatDesc {
		if formatDesc[format] == "Motion-JPEG" {
			mjpegFormat = format
			break
		}
	}
	if mjpegFormat == 0 {
		c.cam.Close()
		return fmt.Errorf("MJPEG format not supported")
	}

	_, _, _, err = c.cam.SetImageFormat(mjpegFormat, 1280, 720)
	if err != nil {
		c.cam.Close()
		return fmt.Errorf("failed to set image format: %w", err)
	}

	// Start streaming
	err = c.cam.StartStreaming()
	if err != nil {
		c.cam.Close()
		return fmt.Errorf("failed to start streaming: %w", err)
	}

	return nil
}

func (c *V4L2Camera) CaptureFrame() ([]byte, error) {
	if c.cam == nil {
		return nil, fmt.Errorf("camera not opened")
	}

	// Wait for frame (5 second timeout)
	err := c.cam.WaitForFrame(5)
	if err != nil {
		return nil, fmt.Errorf("timeout waiting for frame: %w", err)
	}

	frame, err := c.cam.ReadFrame()
	if err != nil {
		return nil, fmt.Errorf("failed to read frame: %w", err)
	}

	// Return JPEG bytes (MJPEG is already JPEG frames)
	return frame, nil
}

func (c *V4L2Camera) Close() error {
	if c.cam != nil {
		if err := c.cam.StopStreaming(); err != nil {
			// Log but continue to close - best effort cleanup
			_ = err
		}
		return c.cam.Close()
	}
	return nil
}
