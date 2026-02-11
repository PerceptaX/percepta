package core

// CameraDriver captures frames from physical camera
// Implementation must be platform-specific (V4L2 on Linux, AVFoundation on macOS, etc.)
type CameraDriver interface {
	Open() error
	CaptureFrame() ([]byte, error) // Returns JPEG bytes
	Close() error
}

// VisionDriver converts camera frames to structured observations
type VisionDriver interface {
	Observe(deviceID string, frame []byte) (*Observation, error)
}

// StorageDriver persists observations
type StorageDriver interface {
	Save(obs Observation) error
	Query(deviceID string, limit int) ([]Observation, error)
	Count() int
}
