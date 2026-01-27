package plugin

import (
	"context"
	"log/slog"
	"os"
	"testing"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"

	"github.com/gllm-dev/hailo-device-plugin/internal/config"
	"github.com/gllm-dev/hailo-device-plugin/internal/domain"
)

type mockDetector struct {
	devices []domain.HailoDevice
	err     error
	healthy bool
}

func (m *mockDetector) Detect() ([]domain.HailoDevice, error) {
	return m.devices, m.err
}

func (m *mockDetector) IsHealthy(_ domain.HailoDevice) bool {
	return m.healthy
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func newTestConfig() *config.Config {
	return &config.Config{
		ResourceName:  "hailo.ai/h10",
		Architecture:  "HAILO10H",
		DevicePath:    "/dev",
		DevicePattern: "hailo*",
	}
}

func TestDiscoverDevices(t *testing.T) {
	tests := []struct {
		name           string
		devices        []domain.HailoDevice
		err            error
		wantErr        bool
		wantCount      int
		wantFirstID    string
		wantFirstState string
	}{
		{
			name: "discovers healthy and unhealthy devices",
			devices: []domain.HailoDevice{
				{ID: "hailo0", Path: "/dev/hailo0", Architecture: "HAILO10H", Healthy: true},
				{ID: "hailo1", Path: "/dev/hailo1", Architecture: "HAILO10H", Healthy: false},
			},
			wantCount:      2,
			wantFirstID:    "hailo0",
			wantFirstState: pluginapi.Healthy,
		},
		{
			name:      "propagates detector error",
			err:       domain.ErrFailedToGlobDevices,
			wantErr:   true,
			wantCount: 0,
		},
		{
			name:      "handles no devices",
			devices:   []domain.HailoDevice{},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := &mockDetector{devices: tt.devices, err: tt.err}
			p := New(newTestConfig(), detector, newTestLogger())

			err := p.discoverDevices()

			if (err != nil) != tt.wantErr {
				t.Errorf("discoverDevices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(p.devices) != tt.wantCount {
				t.Errorf("got %d devices, want %d", len(p.devices), tt.wantCount)
			}
			if tt.wantFirstID != "" && p.devices[0].ID != tt.wantFirstID {
				t.Errorf("first device ID = %s, want %s", p.devices[0].ID, tt.wantFirstID)
			}
			if tt.wantFirstState != "" && p.devices[0].Health != tt.wantFirstState {
				t.Errorf("first device health = %s, want %s", p.devices[0].Health, tt.wantFirstState)
			}
		})
	}
}

func TestAllocate(t *testing.T) {
	p := New(newTestConfig(), &mockDetector{}, newTestLogger())

	req := &pluginapi.AllocateRequest{
		ContainerRequests: []*pluginapi.ContainerAllocateRequest{
			{DevicesIds: []string{"hailo0"}},
			{DevicesIds: []string{"hailo1", "hailo2"}},
		},
	}

	resp, err := p.Allocate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.ContainerResponses) != 2 {
		t.Fatalf("expected 2 container responses, got %d", len(resp.ContainerResponses))
	}

	cr0 := resp.ContainerResponses[0]
	if len(cr0.Devices) != 1 {
		t.Fatalf("expected 1 device for container 0, got %d", len(cr0.Devices))
	}
	if cr0.Devices[0].HostPath != "/dev/hailo0" {
		t.Errorf("expected host path /dev/hailo0, got %s", cr0.Devices[0].HostPath)
	}
	if cr0.Devices[0].ContainerPath != "/dev/hailo0" {
		t.Errorf("expected container path /dev/hailo0, got %s", cr0.Devices[0].ContainerPath)
	}
	if cr0.Devices[0].Permissions != DevicePermissions {
		t.Errorf("expected permissions %s, got %s", DevicePermissions, cr0.Devices[0].Permissions)
	}
	if cr0.Envs[HailoDeviceEnvVar] != "hailo0" {
		t.Errorf("expected HAILO_DEVICE=hailo0, got %s", cr0.Envs[HailoDeviceEnvVar])
	}

	cr1 := resp.ContainerResponses[1]
	if len(cr1.Devices) != 2 {
		t.Fatalf("expected 2 devices for container 1, got %d", len(cr1.Devices))
	}
	if cr1.Envs[HailoDeviceEnvVar] != "hailo1" {
		t.Errorf("expected HAILO_DEVICE=hailo1 (first device), got %s", cr1.Envs[HailoDeviceEnvVar])
	}
}
