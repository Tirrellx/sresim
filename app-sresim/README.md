# SRE Simulation Application

This is the main application component of the SRE Simulation Environment. It provides a comprehensive set of tools for simulating various SRE scenarios and testing system resilience.

## Prerequisites

- Go 1.21 or later
- Docker
- Kubernetes cluster
- Prometheus Operator
- Grafana

## Quick Start

### Local Build and Run
```bash
# Build the application
go build -o sresim ./cmd/main.go

# Run locally
./sresim
```

### Docker Builds

The application provides two Dockerfile variants for different use cases:

#### Production Build
```bash
# Build production image
docker build -f Dockerfile.prod -t sresim:prod .

# Run production container
docker run -p 8080:8080 sresim:prod
```

Features:
- Security-hardened distroless base image
- Static binary compilation
- Security scanning with Trivy
- Read-only root filesystem
- Non-root user execution
- Minimal attack surface

#### Development Build
```bash
# Build development image
docker build -f Dockerfile.dev -t sresim:dev .

# Run development container with debugging
docker run -p 8080:8080 -p 2345:2345 sresim:dev
```

Features:
- Debug-enabled distroless image
- Delve debugger integration
- Debug symbols preserved
- Remote debugging support
- Detailed logging
- Full stack traces

#### VS Code Debugging Configuration
Add the following to `.vscode/launch.json`:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Attach to Docker",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "/app",
            "port": 2345,
            "host": "127.0.0.1"
        }
    ]
}
```

### Docker Compose Development Environment
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Rebuild and restart a specific service
docker-compose up -d --build sresim
```

The services will be available at:
- SRE Simulator: http://localhost:8080
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin)
- Jaeger UI: http://localhost:16686

## API Reference

### Endpoints

- `GET /scenarios` - List available simulation scenarios
- `POST /scenarios/{name}/run` - Run a specific scenario
- `POST /scenarios/{name}/stop` - Stop a running scenario
- `GET /metrics` - Prometheus metrics endpoint
- `GET /health` - Health check endpoint

### Simulation Scenarios

#### High Latency
Simulates network latency by introducing artificial delays in request processing.
```bash
curl -X POST http://localhost:8080/scenarios/latency/run \
  -H "Content-Type: application/json" \
  -d '{"delay_ms": 1000}'
```
Parameters:
- `delay_ms`: Delay in milliseconds (default: 1000)

#### High Error Rate
Simulates service errors by randomly failing requests.
```bash
curl -X POST http://localhost:8080/scenarios/error_rate/run \
  -H "Content-Type: application/json" \
  -d '{"error_percentage": 50}'
```
Parameters:
- `error_percentage`: Percentage of requests to fail (default: 50)

#### Resource Exhaustion
Simulates CPU and memory exhaustion.
```bash
curl -X POST http://localhost:8080/scenarios/resource_exhaustion/run \
  -H "Content-Type: application/json" \
  -d '{"cpu_percentage": 90, "memory_percentage": 85}'
```
Parameters:
- `cpu_percentage`: Target CPU usage percentage (default: 90)
- `memory_percentage`: Target memory usage percentage (default: 85)

#### Circuit Breaker
Simulates circuit breaker pattern with configurable thresholds.
```bash
curl -X POST http://localhost:8080/scenarios/circuit_breaker/run \
  -H "Content-Type: application/json" \
  -d '{"threshold": 5, "timeout": 30}'
```
Parameters:
- `threshold`: Number of failures before opening circuit (default: 5)
- `timeout`: Time in seconds before attempting to close circuit (default: 30)

#### Rate Limiting
Simulates rate limiting with configurable request rates.
```bash
curl -X POST http://localhost:8080/scenarios/rate_limit/run \
  -H "Content-Type: application/json" \
  -d '{"requests_per_second": 10}'
```
Parameters:
- `requests_per_second`: Maximum requests allowed per second (default: 10)

#### Network Partition
Simulates network partition scenarios.
```bash
curl -X POST http://localhost:8080/scenarios/network_partition/run \
  -H "Content-Type: application/json" \
  -d '{"partition_duration": 60}'
```
Parameters:
- `partition_duration`: Duration of partition in seconds (default: 60)

#### Memory Leak
Simulates memory leak by continuously allocating memory.
```bash
curl -X POST http://localhost:8080/scenarios/memory_leak/run \
  -H "Content-Type: application/json" \
  -d '{"leak_rate_mb_per_second": 10, "duration_seconds": 300}'
```
Parameters:
- `leak_rate_mb_per_second`: Memory leak rate in MB/s (default: 10)
- `duration_seconds`: Duration of leak in seconds (default: 300)

#### CPU Spike
Simulates sudden CPU usage spikes.
```bash
curl -X POST http://localhost:8080/scenarios/cpu_spike/run \
  -H "Content-Type: application/json" \
  -d '{"spike_percentage": 95, "duration_seconds": 30, "interval_seconds": 60}'
```
Parameters:
- `spike_percentage`: CPU usage during spike (default: 95)
- `duration_seconds`: Duration of each spike (default: 30)
- `interval_seconds`: Time between spikes (default: 60)

#### Disk I/O Saturation
Simulates high disk I/O operations.
```bash
curl -X POST http://localhost:8080/scenarios/disk_io/run \
  -H "Content-Type: application/json" \
  -d '{"io_operations_per_second": 1000, "file_size_mb": 100}'
```
Parameters:
- `io_operations_per_second`: Number of I/O operations per second (default: 1000)
- `file_size_mb`: Size of test file in MB (default: 100)

#### Connection Pool Exhaustion
Simulates database connection pool exhaustion.
```bash
curl -X POST http://localhost:8080/scenarios/connection_pool_exhaustion/run \
  -H "Content-Type: application/json" \
  -d '{"max_connections": 10, "hold_time_seconds": 30}'
```
Parameters:
- `max_connections`: Maximum number of connections (default: 10)
- `hold_time_seconds`: Time to hold connections (default: 30)

#### Cascading Failure
Simulates cascading failures across services.
```bash
curl -X POST http://localhost:8080/scenarios/cascading_failure/run \
  -H "Content-Type: application/json" \
  -d '{"failure_chain_length": 3, "delay_between_failures_seconds": 5}'
```
Parameters:
- `failure_chain_length`: Number of services in failure chain (default: 3)
- `delay_between_failures_seconds`: Delay between failures (default: 5)

#### Thundering Herd
Simulates thundering herd problem with concurrent requests.
```bash
curl -X POST http://localhost:8080/scenarios/thundering_herd/run \
  -H "Content-Type: application/json" \
  -d '{"concurrent_requests": 100, "cache_miss_percentage": 80}'
```
Parameters:
- `concurrent_requests`: Number of concurrent requests (default: 100)
- `cache_miss_percentage`: Percentage of cache misses (default: 80)

## Monitoring

### Grafana Dashboard

The dashboard provides real-time visualization of various metrics:

1. **Request Rate Panel**
   - Shows HTTP request rate over time
   - Filtered by handler and method
   - 5-minute rate window

2. **Average Response Time Panel**
   - Displays average HTTP response time
   - Filtered by handler
   - Calculated from histogram metrics

3. **Error Rate Panel**
   - Shows HTTP error rate
   - Filtered by handler and method
   - 5-minute rate window

4. **CPU Usage Panel**
   - Real-time CPU usage by scenario
   - Percentage-based visualization
   - Threshold indicators

5. **Memory Usage Panel**
   - Memory consumption by scenario
   - Byte-based visualization
   - Memory leak detection

6. **Network Latency Panel**
   - Network latency measurements
   - Filtered by scenario type
   - 5-minute average

7. **Circuit Breaker Panel**
   - Circuit breaker state changes
   - Failure rate visualization
   - State transitions

8. **Rate Limit Panel**
   - Rate limit hits over time
   - Current rate limit values
   - Burst detection

### Prometheus Metrics

Key metrics available:

1. **HTTP Metrics**
   - `http_request_duration_seconds`: Request duration histogram
   - `http_requests_total`: Total request counter
   - `http_errors_total`: Error counter

2. **Scenario Metrics**
   - `sresim_active_scenarios`: Active scenario gauge
   - `sresim_scenario_duration_seconds`: Scenario duration histogram
   - `sresim_scenario_errors_total`: Scenario error counter

3. **Resource Metrics**
   - `sresim_cpu_usage_percent`: CPU usage gauge
   - `sresim_memory_usage_bytes`: Memory usage gauge
   - `sresim_disk_io_bytes_total`: Disk I/O counter

4. **Network Metrics**
   - `sresim_network_latency_seconds`: Network latency histogram
   - `sresim_network_errors_total`: Network error counter

5. **Circuit Breaker Metrics**
   - `sresim_circuit_breaker_state`: Circuit breaker state gauge
   - `sresim_circuit_breaker_failures_total`: Failure counter

6. **Rate Limiting Metrics**
   - `sresim_rate_limit_hits_total`: Rate limit hit counter
   - `sresim_rate_limit_current`: Current rate limit gauge

### Health Checks

The application provides health check endpoints:

1. **Basic Health Check**
```bash
curl http://localhost:8080/health
```
Response:
```json
{
  "status": "healthy",
  "timestamp": "2024-03-25T19:57:00Z",
  "active_scenarios": ["latency", "error_rate"]
}
```

2. **Detailed Health Check**
```bash
curl http://localhost:8080/health/detailed
```
Response:
```json
{
  "status": "healthy",
  "timestamp": "2024-03-25T19:57:00Z",
  "active_scenarios": ["latency", "error_rate"],
  "resource_usage": {
    "cpu_percent": 45.2,
    "memory_bytes": 256000000,
    "disk_io_bytes": 1024000
  },
  "network_status": {
    "latency_ms": 50,
    "error_rate": 0.1
  }
}
```

## Configuration

### ConfigMap Settings

```yaml
server:
  port: 8080
  read_timeout: 10s
  write_timeout: 10s
  idle_timeout: 60s

metrics:
  enabled: true
  path: /metrics
  port: 8080

scenarios:
  default_duration: 5m
  max_concurrent: 3
  cleanup_interval: 1m

circuit_breaker:
  failure_threshold: 5
  reset_timeout: 30s
  half_open_timeout: 10s

rate_limiter:
  requests_per_second: 10
  burst_size: 20

resource_limits:
  max_cpu_percent: 80
  max_memory_bytes: 512Mi
  max_disk_io_bytes: 1Gi

network:
  max_latency_ms: 1000
  error_rate_percent: 5
  partition_probability: 0.1
```

### Environment Variables

- `PROMETHEUS_MULTIPROC_DIR`: Directory for Prometheus multiprocess mode
- `OTEL_EXPORTER_OTLP_ENDPOINT`: OpenTelemetry collector endpoint
- `CONFIG_FILE`: Path to configuration file
- `LOG_LEVEL`: Logging level (debug, info, warn, error)
- `ENABLE_METRICS`: Enable/disable metrics collection
- `ENABLE_TRACING`: Enable/disable tracing

## Troubleshooting

### Common Issues

1. **Scenario Not Starting**
   - Check if scenario is already running
   - Verify scenario parameters
   - Check resource limits
   - Review logs for errors

2. **High Resource Usage**
   - Monitor CPU and memory usage
   - Check for memory leaks
   - Review scenario configurations
   - Adjust resource limits

3. **Metrics Not Showing**
   - Verify Prometheus configuration
   - Check ServiceMonitor labels
   - Ensure metrics endpoint is accessible
   - Review Grafana data source

4. **Network Issues**
   - Check network partition scenarios
   - Verify latency settings
   - Review error rate configurations
   - Check service connectivity

### Logging

Logs are available in the following formats:
- JSON structured logging
- Human-readable format
- Log levels: debug, info, warn, error

Example log output:
```json
{
  "level": "info",
  "time": "2024-03-25T19:57:00Z",
  "msg": "Starting scenario",
  "scenario": "latency",
  "parameters": {
    "delay_ms": 1000
  }
}
```

### Debug Mode

Enable debug mode for detailed logging:
```bash
export LOG_LEVEL=debug
```

## Development

### Project Structure

```
app-sresim/
├── cmd/
│   └── main.go
├── pkg/
│   ├── metrics/
│   │   └── metrics.go
│   ├── simulator/
│   │   ├── scenarios.go
│   │   └── implementations.go
│   └── health/
│       └── health.go
├── k8s/
│   ├── configmap.yaml
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── service-monitor.yaml
│   └── grafana-dashboard.yaml
└── README.md
```

### Adding New Scenarios

1. Define the scenario in `pkg/simulator/scenarios.go`
2. Implement the scenario in `pkg/simulator/implementations.go`
3. Add metrics in `pkg/metrics/metrics.go`
4. Update the Grafana dashboard if needed

### Testing

Run tests:
```bash
go test ./...
```

Run specific test:
```bash
go test ./pkg/simulator -run TestLatencyScenario
``` 