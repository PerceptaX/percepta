//go:build (!linux && !darwin) || !cgo

package camera

import (
	"fmt"

	"github.com/perceptumx/percepta/internal/core"
)

type stubCamera struct{}

// NewCamera returns a stub camera driver for unsupported platforms
func NewCamera(devicePath string) core.CameraDriver {
	return &stubCamera{}
}

func (s *stubCamera) Open() error {
	return fmt.Errorf("camera not supported on this platform")
}

func (s *stubCamera) CaptureFrame() ([]byte, error) {
	return nil, fmt.Errorf("camera not supported on this platform")
}

func (s *stubCamera) Close() error {
	return nil
}
