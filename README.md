# SRE Simulation Environment

A comprehensive environment for simulating and testing Site Reliability Engineering (SRE) scenarios. This project provides tools and frameworks for testing system resilience, monitoring, and observability in a controlled environment.

## Project Structure

```
sresim/
├── app-sresim/          # Main application for SRE scenario simulation
│   ├── cmd/            # Application entry points
│   ├── pkg/            # Core packages
│   ├── k8s/            # Kubernetes configurations
│   ├── Dockerfile.prod # Production Dockerfile
│   ├── Dockerfile.dev  # Development Dockerfile
│   └── docker-compose.yml
├── .github/            # GitHub workflows and templates
└── .devbox/            # Development environment configuration
```

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/yourusername/sresim.git
cd sresim
```

2. Set up the development environment:
```bash
# Using devbox
devbox install

# Or manually install dependencies
go install
```

3. Build and run the application:
```bash
cd app-sresim
go build -o sresim ./cmd/main.go
./sresim
```

## Features

- **Scenario Simulation**: Test various SRE scenarios including:
  - High latency
  - Error rates
  - Resource exhaustion
  - Network partitions
  - And more...

- **Monitoring & Observability**:
  - Prometheus metrics
  - OpenTelemetry integration
  - Grafana dashboards
  - Health checks
  - Resource usage tracking

- **Kubernetes Integration**:
  - Ready-to-use Kubernetes manifests
  - ConfigMap support
  - Service monitoring
  - Resource management

- **Security Features**:
  - Distroless container images
  - Security scanning with Trivy
  - Read-only root filesystem
  - Non-root user execution
  - Static binary compilation
  - Minimal attack surface

- **Development Tools**:
  - Remote debugging support
  - Development-specific Dockerfile
  - VS Code integration
  - Docker Compose environment
  - Detailed logging

## Docker Builds

The application provides two Dockerfile variants:

### Production Build
```bash
docker build -f Dockerfile.prod -t sresim:prod .
docker run -p 8080:8080 sresim:prod
```

Features:
- Security-hardened distroless base image
- Static binary compilation
- Security scanning with Trivy
- Read-only root filesystem
- Non-root user execution
- Minimal attack surface

### Development Build
```bash
docker build -f Dockerfile.dev -t sresim:dev .
docker run -p 8080:8080 -p 2345:2345 sresim:dev
```

Features:
- Debug-enabled distroless image
- Delve debugger integration
- Debug symbols preserved
- Remote debugging support
- Detailed logging
- Full stack traces

## Documentation

- [Application Documentation](app-sresim/README.md) - Detailed guide for the main application
- [Development Guide](.github/CONTRIBUTING.md) - Guidelines for contributing
- [Architecture](docs/architecture.md) - System architecture and design decisions

## Contributing

We welcome contributions! Please see our [Contributing Guide](.github/CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
