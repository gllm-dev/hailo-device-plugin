package domain

// HailoDevice represents a Hailo accelerator
type HailoDevice struct {
	ID           string
	Path         string
	Architecture string
	Healthy      bool
}

// DeviceDetector interface for detecting Hailo devices
type DeviceDetector interface {
	Detect() ([]HailoDevice, error)
	IsHealthy(device HailoDevice) bool
}
