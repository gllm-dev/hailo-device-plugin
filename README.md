<p align="center">
  <img src=".github/banner.png" alt="Hailo Device Plugin for Kubernetes">
</p>

# Hailo Device Plugin for Kubernetes

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gllm-dev/hailo-device-plugin)](go.mod)
[![Release](https://img.shields.io/github/v/release/gllm-dev/hailo-device-plugin)](https://github.com/gllm-dev/hailo-device-plugin/releases)

A Kubernetes [Device Plugin](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/device-plugins/) for [Hailo AI accelerators](https://hailo.ai/), enabling seamless scheduling of AI inference workloads on edge devices like the Raspberry Pi AI HAT+.

## Quick Start

### 1. Deploy the Plugin

```bash
kubectl apply -f https://raw.githubusercontent.com/gllm-dev/hailo-device-plugin/main/deploy/daemonset.yaml
```

### 2. Verify Installation

```bash
kubectl get nodes -o json | jq '.items[].status.allocatable["hailo.ai/h10"]'
```

### 3. Use in Your Workloads

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: inference
spec:
  containers:
    - name: model
      image: your-inference-image
      resources:
        limits:
          hailo.ai/h10: 1
```

## Prerequisites

- Kubernetes v1.26+
- Hailo driver installed on worker nodes ([installation guide](https://hailo.ai/developer-zone/))
- Hailo device accessible at `/dev/hailo*`

## Installation

### Option 1: DaemonSet (Recommended)

Deploy as a DaemonSet to automatically run on all nodes with Hailo devices:

```bash
kubectl apply -f deploy/daemonset.yaml
```

Or directly from the repository:

```bash
kubectl apply -f https://raw.githubusercontent.com/gllm-dev/hailo-device-plugin/main/deploy/daemonset.yaml
```

### Option 2: Helm (Coming Soon)

```bash
helm repo add gllm https://gllm-dev.github.io/charts
helm install hailo-device-plugin gllm/hailo-device-plugin
```

## Configuration

Configure the plugin using environment variables:

| Environment Variable   | Default        | Description                          |
|------------------------|----------------|--------------------------------------|
| `HAILO_RESOURCE_NAME`  | `hailo.ai/h10` | Kubernetes resource name to register |
| `HAILO_ARCHITECTURE`   | `HAILO10H`     | Device architecture identifier       |
| `HAILO_DEVICE_PATH`    | `/dev`         | Path to device directory             |
| `HAILO_DEVICE_PATTERN` | `hailo*`       | Glob pattern to match device files   |

### Example: Hailo-8L Configuration

```yaml
env:
  - name: HAILO_RESOURCE_NAME
    value: "hailo.ai/h8l"
  - name: HAILO_ARCHITECTURE
    value: "HAILO8L"
```

## Supported Devices

| Device                  | Resource Name    | Architecture | Default |
|-------------------------|------------------|--------------|---------|
| AI HAT+ 2 (Hailo-10H)   | `hailo.ai/h10`   | `HAILO10H`   | Yes     |
| AI HAT+ (Hailo-8L)      | `hailo.ai/h8l`   | `HAILO8L`    | No      |
| Hailo-8L                | `hailo.ai/h8l`   | `HAILO8L`    | No      |
| Hailo-8                 | `hailo.ai/h8`    | `HAILO8`     | No      |

## Building from Source

```bash
# Clone the repository
git clone https://github.com/gllm-dev/hailo-device-plugin.git
cd hailo-device-plugin

# Build binary
go build -o hailo-device-plugin ./cmd/plugin

# Build container
docker build -t hailo-device-plugin:local .
```

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting a Pull Request.

This project follows [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) and [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Hailo](https://hailo.ai/) for their AI accelerator technology
- [Kubernetes Device Plugin framework](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/device-plugins/)
