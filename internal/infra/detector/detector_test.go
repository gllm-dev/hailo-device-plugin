package detector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gllm-dev/hailo-device-plugin/internal/config"
	"github.com/gllm-dev/hailo-device-plugin/internal/domain"
)

func TestDetect_NoDevices(t *testing.T) {
	cfg := &config.Config{
		DevicePath:    t.TempDir(),
		DevicePattern: "hailo*",
		Architecture:  "HAILO10H",
	}
	d := New(cfg)

	devices, err := d.Detect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(devices) != 0 {
		t.Errorf("expected 0 devices, got %d", len(devices))
	}
}

func TestDetect_WithDevices(t *testing.T) {
	tmpDir := t.TempDir()

	for _, name := range []string{"hailo0", "hailo1"} {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte{}, 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	cfg := &config.Config{
		DevicePath:    tmpDir,
		DevicePattern: "hailo*",
		Architecture:  "HAILO10H",
	}
	d := New(cfg)

	devices, err := d.Detect()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(devices) != 2 {
		t.Errorf("expected 2 devices, got %d", len(devices))
	}

	for _, dev := range devices {
		if dev.Architecture != "HAILO10H" {
			t.Errorf("expected architecture HAILO10H, got %s", dev.Architecture)
		}
	}
}

func TestIsHealthy_DeviceNotExists(t *testing.T) {
	cfg := &config.Config{
		DevicePath:    "/dev",
		DevicePattern: "hailo*",
		Architecture:  "HAILO10H",
	}
	d := New(cfg)
	device := domain.HailoDevice{Path: "/nonexistent/device"}

	if d.IsHealthy(device) {
		t.Error("expected nonexistent device to not be healthy")
	}
}
