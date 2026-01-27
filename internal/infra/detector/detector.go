package detector

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gllm-dev/hailo-device-plugin/internal/config"
	"github.com/gllm-dev/hailo-device-plugin/internal/domain"
)

// HailoDetector detects Hailo AI accelerator devices on the system.
type HailoDetector struct {
	devPath      string
	pattern      string
	architecture string
}

// New creates a new HailoDetector with the given configuration.
func New(cfg *config.Config) *HailoDetector {
	return &HailoDetector{
		devPath:      cfg.DevicePath,
		pattern:      cfg.DevicePattern,
		architecture: cfg.Architecture,
	}
}

// Detect discovers all Hailo devices in the configured device path.
func (d *HailoDetector) Detect() ([]domain.HailoDevice, error) {
	var devices []domain.HailoDevice

	matches, err := filepath.Glob(filepath.Join(d.devPath, d.pattern))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrFailedToGlobDevices, err)
	}

	for _, path := range matches {
		device := domain.HailoDevice{
			ID:           filepath.Base(path),
			Path:         path,
			Architecture: d.architecture,
			Healthy:      d.IsHealthy(domain.HailoDevice{Path: path}),
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// IsHealthy checks if a Hailo device is healthy and accessible.
func (d *HailoDetector) IsHealthy(device domain.HailoDevice) bool {
	info, err := os.Stat(device.Path)
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeDevice != 0
}
