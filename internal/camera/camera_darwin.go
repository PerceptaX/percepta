//go:build darwin && cgo

package camera

import "github.com/perceptumx/percepta/internal/core"

// NewCamera creates a platform-specific camera driver
func NewCamera(devicePath string) core.CameraDriver {
	return NewAVFoundationCamera(devicePath)
}
