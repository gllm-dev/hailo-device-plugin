package plugin

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"

	"github.com/gllm-dev/hailo-device-plugin/internal/config"
	"github.com/gllm-dev/hailo-device-plugin/internal/domain"
)

const (
	SocketName          = "hailo.sock"
	KubeletSocket       = "/var/lib/kubelet/device-plugins/kubelet.sock"
	DevicePluginPath    = "/var/lib/kubelet/device-plugins/"
	HealthCheckInterval = 10 * time.Second
	ServerStartupDelay  = time.Second
	DevicePermissions   = "rw"
	HailoDeviceEnvVar   = "HAILO_DEVICE"
)

// HailoPlugin implements the Kubernetes device plugin interface for Hailo devices.
type HailoPlugin struct {
	pluginapi.UnimplementedDevicePluginServer

	config   *config.Config
	detector domain.DeviceDetector
	devices  []*pluginapi.Device
	socket   string
	server   *grpc.Server
	stop     chan struct{}
	health   chan *pluginapi.Device
	logger   *slog.Logger
}

// New creates a new HailoPlugin with the given configuration, detector, and logger.
func New(cfg *config.Config, detector domain.DeviceDetector, logger *slog.Logger) *HailoPlugin {
	return &HailoPlugin{
		config:   cfg,
		detector: detector,
		socket:   filepath.Join(DevicePluginPath, SocketName),
		stop:     make(chan struct{}),
		health:   make(chan *pluginapi.Device),
		logger:   logger,
	}
}

// Start discovers devices, starts the gRPC server, and registers with kubelet.
func (p *HailoPlugin) Start() error {
	if err := p.discoverDevices(); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrFailedToDiscoverDevices, err)
	}

	if len(p.devices) == 0 {
		return domain.ErrNoDevicesFound
	}

	p.logger.Info("discovered devices", slog.Int("count", len(p.devices)))

	if err := os.Remove(p.socket); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("%w: %v", domain.ErrFailedToRemoveSocket, err)
	}

	listener, err := net.Listen("unix", p.socket)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrFailedToListenOnSocket, err)
	}

	p.server = grpc.NewServer()
	pluginapi.RegisterDevicePluginServer(p.server, p)

	go func() {
		p.logger.Info("starting gRPC server", slog.String("socket", p.socket))
		if err := p.server.Serve(listener); err != nil {
			p.logger.Error("gRPC server error", slog.Any("error", err))
		}
	}()

	time.Sleep(ServerStartupDelay)

	if err := p.register(); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrFailedToRegister, err)
	}

	p.logger.Info("registered with kubelet", slog.String("resource", p.config.ResourceName))

	go p.healthCheck()

	return nil
}

func (p *HailoPlugin) discoverDevices() error {
	hailoDevices, err := p.detector.Detect()
	if err != nil {
		return err
	}

	p.devices = make([]*pluginapi.Device, len(hailoDevices))
	for i, d := range hailoDevices {
		health := pluginapi.Unhealthy
		if d.Healthy {
			health = pluginapi.Healthy
		}

		p.devices[i] = &pluginapi.Device{
			ID:     d.ID,
			Health: health,
		}
		p.logger.Info("discovered device",
			slog.String("id", d.ID),
			slog.String("architecture", d.Architecture),
			slog.String("health", health),
		)
	}

	return nil
}

func (p *HailoPlugin) register() error {
	conn, err := grpc.NewClient(
		"unix://"+KubeletSocket,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrFailedToCreateClient, err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			p.logger.Error("failed to close connection", slog.Any("error", err))
		}
	}()

	client := pluginapi.NewRegistrationClient(conn)

	req := &pluginapi.RegisterRequest{
		Version:      pluginapi.Version,
		Endpoint:     SocketName,
		ResourceName: p.config.ResourceName,
		Options: &pluginapi.DevicePluginOptions{
			PreStartRequired: false,
		},
	}

	_, err = client.Register(context.Background(), req)
	return err
}

func (p *HailoPlugin) healthCheck() {
	ticker := time.NewTicker(HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.stop:
			return
		case <-ticker.C:
			for _, device := range p.devices {
				healthy := p.detector.IsHealthy(domain.HailoDevice{
					Path: filepath.Join(p.config.DevicePath, device.ID),
				})

				newHealth := pluginapi.Unhealthy
				if healthy {
					newHealth = pluginapi.Healthy
				}

				if device.Health != newHealth {
					device.Health = newHealth
					p.health <- device
				}
			}
		}
	}
}

// Stop gracefully stops the plugin and gRPC server.
func (p *HailoPlugin) Stop() {
	close(p.stop)
	if p.server != nil {
		p.server.Stop()
	}
}

// GetDevicePluginOptions returns the device plugin options.
func (p *HailoPlugin) GetDevicePluginOptions(_ context.Context, _ *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{
		PreStartRequired: false,
	}, nil
}

// ListAndWatch streams the list of devices and updates when health changes.
func (p *HailoPlugin) ListAndWatch(_ *pluginapi.Empty, stream pluginapi.DevicePlugin_ListAndWatchServer) error {
	p.logger.Info("ListAndWatch called")

	if err := stream.Send(&pluginapi.ListAndWatchResponse{Devices: p.devices}); err != nil {
		return err
	}

	for {
		select {
		case <-p.stop:
			return nil
		case d := <-p.health:
			p.logger.Info("device health changed",
				slog.String("device_id", d.ID),
				slog.String("health", d.Health),
			)
			if err := stream.Send(&pluginapi.ListAndWatchResponse{Devices: p.devices}); err != nil {
				return err
			}
		}
	}
}

// Allocate allocates devices for container requests.
func (p *HailoPlugin) Allocate(_ context.Context, req *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	p.logger.Info("allocate called", slog.Int("containers", len(req.ContainerRequests)))

	responses := make([]*pluginapi.ContainerAllocateResponse, len(req.ContainerRequests))

	for i, containerReq := range req.ContainerRequests {
		p.logger.Info("container requesting devices",
			slog.Int("container_index", i),
			slog.Any("device_ids", containerReq.DevicesIds),
		)

		devices := make([]*pluginapi.DeviceSpec, len(containerReq.DevicesIds))
		for j, deviceID := range containerReq.DevicesIds {
			devicePath := filepath.Join(p.config.DevicePath, deviceID)
			devices[j] = &pluginapi.DeviceSpec{
				ContainerPath: devicePath,
				HostPath:      devicePath,
				Permissions:   DevicePermissions,
			}
		}

		responses[i] = &pluginapi.ContainerAllocateResponse{
			Devices: devices,
			Envs: map[string]string{
				HailoDeviceEnvVar: containerReq.DevicesIds[0],
			},
		}
	}

	return &pluginapi.AllocateResponse{ContainerResponses: responses}, nil
}

// PreStartContainer is called before container start; no-op for this plugin.
func (p *HailoPlugin) PreStartContainer(_ context.Context, _ *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}

// GetPreferredAllocation returns preferred device allocation; no-op for this plugin.
func (p *HailoPlugin) GetPreferredAllocation(_ context.Context, _ *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return &pluginapi.PreferredAllocationResponse{}, nil
}
