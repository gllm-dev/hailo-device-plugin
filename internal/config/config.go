package config

import (
	"os"
)

const (
	EnvResourceName  = "HAILO_RESOURCE_NAME"
	EnvArchitecture  = "HAILO_ARCHITECTURE"
	EnvDevicePath    = "HAILO_DEVICE_PATH"
	EnvDevicePattern = "HAILO_DEVICE_PATTERN"
)

const (
	DefaultResourceName  = "hailo.ai/h10"
	DefaultArchitecture  = "HAILO10H"
	DefaultDevicePath    = "/dev"
	DefaultDevicePattern = "hailo*"
)

// Config holds the plugin configuration.
type Config struct {
	ResourceName  string
	Architecture  string
	DevicePath    string
	DevicePattern string
}

// Load reads configuration from environment variables with defaults.
func Load() *Config {
	return &Config{
		ResourceName:  getEnv(EnvResourceName, DefaultResourceName),
		Architecture:  getEnv(EnvArchitecture, DefaultArchitecture),
		DevicePath:    getEnv(EnvDevicePath, DefaultDevicePath),
		DevicePattern: getEnv(EnvDevicePattern, DefaultDevicePattern),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
