// Package domain contains domain entities and errors.
package domain

import "errors"

// Domain errors for the Hailo device plugin.
var (
	ErrNoDevicesFound          = errors.New("no Hailo devices found")
	ErrFailedToDiscoverDevices = errors.New("failed to discover devices")
	ErrFailedToRemoveSocket    = errors.New("failed to remove socket")
	ErrFailedToListenOnSocket  = errors.New("failed to listen on socket")
	ErrFailedToCreateClient    = errors.New("failed to create kubelet client")
	ErrFailedToRegister        = errors.New("failed to register with kubelet")
	ErrFailedToGlobDevices     = errors.New("failed to glob hailo devices")
)
