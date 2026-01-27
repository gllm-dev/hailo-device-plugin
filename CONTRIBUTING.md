# Contributing to Hailo Device Plugin

Thank you for your interest in contributing!

## Code of Conduct

Please be respectful and constructive in all interactions.

## How to Contribute

### Reporting Issues

- Check existing issues before creating a new one
- Include relevant details: Kubernetes version, Hailo device model, logs, and steps to reproduce

### Submitting Changes

1. Create a branch from `main`:
   ```bash
   git checkout -b feat/your-feature-name
   ```
2. Make your changes following the coding standards below
3. Test your changes thoroughly
4. Commit using [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/):
   ```bash
   git commit -m "feat: add support for Hailo-8"
   ```
5. Push and open a Pull Request

### Commit Message Format

This project follows [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/):

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```bash
feat(detector): add support for Hailo-8L architecture
fix(plugin): handle device disconnect during health check
docs(readme): update configuration examples
refactor(config): simplify environment variable loading
test(plugin): add allocation edge case tests
chore(deps): update grpc to v1.79.0
```

### Changelog

This project follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/). When submitting changes, update `CHANGELOG.md` under the `[Unreleased]` section:

- **Added** for new features
- **Changed** for changes in existing functionality
- **Deprecated** for soon-to-be removed features
- **Removed** for now removed features
- **Fixed** for bug fixes
- **Security** for vulnerability fixes

## Development Setup

### Prerequisites

- Go 1.25
- Docker
- Access to a Kubernetes cluster (minikube, kind, or real cluster)

### Building

```bash
# Build binary
go build -o hailo-device-plugin ./cmd/plugin

# Run tests
go test ./...

# Run linter
golangci-lint run ./...

# Build container
docker build -t hailo-device-plugin:dev .
```

## Coding Standards

- Follow standard Go conventions and `gofmt`
- Run `golangci-lint` before committing
- Use structured logging with `slog`
- Define errors in `internal/domain/errors.go`
- Add tests for new functionality

## Pull Request Guidelines

- Keep PRs focused on a single change
- Update `CHANGELOG.md` under `[Unreleased]`
- Ensure CI passes
