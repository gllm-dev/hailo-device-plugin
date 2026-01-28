<p align="center">
  <img src=".github/banner.png" alt="Hailo Device Plugin for Kubernetes">
</p>

# Hailo Device Plugin for Kubernetes

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gllm-dev/hailo-device-plugin)](go.mod)
[![Release](https://img.shields.io/github/v/release/gllm-dev/hailo-device-plugin)](https://github.com/gllm-dev/hailo-device-plugin/releases)

A Kubernetes [Device Plugin](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/device-plugins/) for [Hailo AI accelerators](https://hailo.ai/), enabling seamless scheduling of AI inference workloads on edge devices like the Raspberry Pi AI HAT+.

## Quick Start

### 1. Label Your Hailo Nodes

The plugin only deploys on nodes with this label, preventing unnecessary pods on nodes without Hailo devices:

```bash
kubectl label nodes <node-name> hailo.ai/device=present
```

### 2. Deploy the Plugin

```bash
kubectl apply -f https://raw.githubusercontent.com/gllm-dev/hailo-device-plugin/main/deploy/daemonset.yaml
```

### 3. Verify Installation

```bash
# Check pod is running
kubectl -n kube-system get pods -l app.kubernetes.io/name=hailo-device-plugin

# Check device is registered
kubectl get nodes -o custom-columns=NAME:.metadata.name,HAILO:.status.allocatable.hailo\\.ai/h10
```

Should show:
```
NAME          HAILO
<node-name>   1
```

### 4. Test Device Access

```bash
cat << 'EOF' | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: hailo-test
spec:
  containers:
  - name: test
    image: debian:bookworm-slim
    command: ["sh", "-c", "ls -la /dev/hailo* && sleep 3600"]
    resources:
      limits:
        hailo.ai/h10: 1
  restartPolicy: Never
EOF

# Check logs
kubectl logs hailo-test
```

Should show `/dev/hailo0`

### 5. Use in Your Workloads

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

The DaemonSet uses a node selector (`hailo.ai/device=present`) to deploy only on labeled nodes.

**1. Label nodes with Hailo devices:**

```bash
kubectl label nodes <node-name> hailo.ai/device=present
```

**2. Deploy the plugin:**

```bash
kubectl apply -f https://raw.githubusercontent.com/gllm-dev/hailo-device-plugin/main/deploy/daemonset.yaml
```

## Configuration

Configure the plugin using environment variables in the DaemonSet:

| Environment Variable   | Default        | Description                          |
|------------------------|----------------|--------------------------------------|
| `HAILO_RESOURCE_NAME`  | `hailo.ai/h10` | Kubernetes resource name to register |
| `HAILO_ARCHITECTURE`   | `HAILO10H`     | Device architecture identifier       |
| `HAILO_DEVICE_PATH`    | `/dev`         | Path to device directory             |
| `HAILO_DEVICE_PATTERN` | `hailo*`       | Glob pattern to match device files   |

### Example: Configuring for Hailo-8L

Edit the DaemonSet to add environment variables:

```bash
kubectl edit daemonset hailo-device-plugin -n kube-system
```

```yaml
spec:
  template:
    spec:
      containers:
        - name: hailo-device-plugin
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

## Building Inference Workloads

This plugin handles device scheduling and allocation. For inference workloads, you'll need to build containers with HailoRT on your Hailo nodes (the SDK is not publicly redistributable).

**Resources:**
- [HailoRT SDK](https://hailo.ai/developer-zone/) - Runtime library for Hailo devices
- [Hailo Model Zoo](https://github.com/hailo-ai/hailo_model_zoo) - Pre-trained models and examples
- [hailo_model_zoo_genai](https://github.com/hailo-ai/hailo_model_zoo_genai) - LLM support (Ollama-compatible API)

## Building from Source

```bash
git clone https://github.com/gllm-dev/hailo-device-plugin.git
cd hailo-device-plugin

go build -o hailo-device-plugin ./cmd/plugin

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
