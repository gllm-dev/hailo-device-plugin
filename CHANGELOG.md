# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
## [0.0.1] - 2026-01-21
### Added
- Node selector label for Hailo devices
- Update README

## [0.0.1] - 2026-01-21

### Added

- Initial release of Hailo Device Plugin for Kubernetes
- Automatic discovery of Hailo devices (`/dev/hailo*`)
- Device health monitoring with configurable interval
- gRPC device plugin server implementation
- Kubelet registration and device allocation
- Structured JSON logging with `slog`
- Scratch-based Docker image for minimal footprint
- Environment variable configuration:
  - `HAILO_RESOURCE_NAME` - Kubernetes resource name (default: `hailo.ai/h10`)
  - `HAILO_ARCHITECTURE` - Device architecture (default: `HAILO10H`)
  - `HAILO_DEVICE_PATH` - Device directory path (default: `/dev`)
  - `HAILO_DEVICE_PATTERN` - Device glob pattern (default: `hailo*`)
- Support for multiple Hailo device types (Hailo-8, Hailo-8L, Hailo-10H)
- CI/CD pipeline with GitHub Actions (lint, test, build, release)
- golangci-lint configuration
