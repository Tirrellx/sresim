package metrics

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

var (
	// HTTP metrics
	requestDuration = prom.NewHistogramVec(
		prom.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prom.DefBuckets,
		},
		[]string{"handler", "method", "status"},
	)

	requestTotal = prom.NewCounterVec(
		prom.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"handler", "method"},
	)

	errorTotal = prom.NewCounterVec(
		prom.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"handler", "method"},
	)

	// Scenario metrics
	activeScenarios = prom.NewGaugeVec(
		prom.GaugeOpts{
			Name: "sresim_active_scenarios",
			Help: "Number of currently active simulation scenarios",
		},
		[]string{"scenario_type"},
	)

	scenarioDuration = prom.NewHistogramVec(
		prom.HistogramOpts{
			Name:    "sresim_scenario_duration_seconds",
			Help:    "Duration of simulation scenarios",
			Buckets: prom.DefBuckets,
		},
		[]string{"scenario_type"},
	)

	scenarioErrors = prom.NewCounterVec(
		prom.CounterOpts{
			Name: "sresim_scenario_errors_total",
			Help: "Total number of scenario errors",
		},
		[]string{"scenario_type", "error_type"},
	)

	// Resource metrics
	cpuUsage = prom.NewGaugeVec(
		prom.GaugeOpts{
			Name: "sresim_cpu_usage_percent",
			Help: "CPU usage percentage",
		},
		[]string{"scenario_type"},
	)

	memoryUsage = prom.NewGaugeVec(
		prom.GaugeOpts{
			Name: "sresim_memory_usage_bytes",
			Help: "Memory usage in bytes",
		},
		[]string{"scenario_type"},
	)

	diskIO = prom.NewCounterVec(
		prom.CounterOpts{
			Name: "sresim_disk_io_bytes_total",
			Help: "Total disk I/O operations in bytes",
		},
		[]string{"scenario_type", "operation_type"},
	)

	// Network metrics
	networkLatency = prom.NewHistogramVec(
		prom.HistogramOpts{
			Name:    "sresim_network_latency_seconds",
			Help:    "Network latency in seconds",
			Buckets: prom.DefBuckets,
		},
		[]string{"scenario_type"},
	)

	networkErrors = prom.NewCounterVec(
		prom.CounterOpts{
			Name: "sresim_network_errors_total",
			Help: "Total number of network errors",
		},
		[]string{"scenario_type", "error_type"},
	)

	// Circuit breaker metrics
	circuitBreakerState = prom.NewGaugeVec(
		prom.GaugeOpts{
			Name: "sresim_circuit_breaker_state",
			Help: "Current state of circuit breaker (0: closed, 1: open, 2: half-open)",
		},
		[]string{"scenario_type"},
	)

	circuitBreakerFailures = prom.NewCounterVec(
		prom.CounterOpts{
			Name: "sresim_circuit_breaker_failures_total",
			Help: "Total number of circuit breaker failures",
		},
		[]string{"scenario_type"},
	)

	// Rate limiting metrics
	rateLimitHits = prom.NewCounterVec(
		prom.CounterOpts{
			Name: "sresim_rate_limit_hits_total",
			Help: "Total number of rate limit hits",
		},
		[]string{"scenario_type"},
	)

	rateLimitCurrent = prom.NewGaugeVec(
		prom.GaugeOpts{
			Name: "sresim_rate_limit_current",
			Help: "Current rate limit value",
		},
		[]string{"scenario_type"},
	)
)

func init() {
	// Register all metrics
	prom.MustRegister(requestDuration)
	prom.MustRegister(requestTotal)
	prom.MustRegister(errorTotal)
	prom.MustRegister(activeScenarios)
	prom.MustRegister(scenarioDuration)
	prom.MustRegister(scenarioErrors)
	prom.MustRegister(cpuUsage)
	prom.MustRegister(memoryUsage)
	prom.MustRegister(diskIO)
	prom.MustRegister(networkLatency)
	prom.MustRegister(networkErrors)
	prom.MustRegister(circuitBreakerState)
	prom.MustRegister(circuitBreakerFailures)
	prom.MustRegister(rateLimitHits)
	prom.MustRegister(rateLimitCurrent)
}

// Init initializes all metrics
func Init() error {
	// Register all metrics
	prom.MustRegister(requestDuration)
	prom.MustRegister(requestTotal)
	prom.MustRegister(errorTotal)
	prom.MustRegister(activeScenarios)
	prom.MustRegister(scenarioDuration)
	prom.MustRegister(scenarioErrors)
	prom.MustRegister(cpuUsage)
	prom.MustRegister(memoryUsage)
	prom.MustRegister(diskIO)
	prom.MustRegister(networkLatency)
	prom.MustRegister(networkErrors)
	prom.MustRegister(circuitBreakerState)
	prom.MustRegister(circuitBreakerFailures)
	prom.MustRegister(rateLimitHits)
	prom.MustRegister(rateLimitCurrent)

	// Initialize OpenTelemetry metrics
	return InitMetrics()
}

// InitMetrics initializes Prometheus and OpenTelemetry metrics
func InitMetrics() error {
	// Initialize Prometheus exporter
	exporter, err := otelprom.New()
	if err != nil {
		return err
	}

	// Create metric provider
	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)

	return nil
}

// MetricsMiddleware wraps HTTP handlers with metrics collection
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer that captures the status code
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		requestDuration.WithLabelValues(r.URL.Path, r.Method, wrapped.status).Observe(duration)
		requestTotal.WithLabelValues(r.URL.Path, r.Method).Inc()

		statusCode, _ := strconv.Atoi(wrapped.status)
		if statusCode >= 400 {
			errorTotal.WithLabelValues(r.URL.Path, r.Method).Inc()
		}
	})
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows us to track the status code
type responseWriter struct {
	http.ResponseWriter
	status string
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: "200"}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = strconv.Itoa(code)
	rw.ResponseWriter.WriteHeader(code)
}

// MetricsHandler returns a handler for the /metrics endpoint
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// ScenarioMetrics provides methods to update scenario-specific metrics
type ScenarioMetrics struct {
	scenarioType string
}

// NewScenarioMetrics creates a new ScenarioMetrics instance
func NewScenarioMetrics(scenarioType string) *ScenarioMetrics {
	return &ScenarioMetrics{scenarioType: scenarioType}
}

// SetActive updates the active scenarios gauge
func (sm *ScenarioMetrics) SetActive(active bool) {
	if active {
		activeScenarios.WithLabelValues(sm.scenarioType).Set(1)
	} else {
		activeScenarios.WithLabelValues(sm.scenarioType).Set(0)
	}
}

// RecordDuration records the duration of a scenario
func (sm *ScenarioMetrics) RecordDuration(duration time.Duration) {
	scenarioDuration.WithLabelValues(sm.scenarioType).Observe(duration.Seconds())
}

// RecordError records a scenario error
func (sm *ScenarioMetrics) RecordError(errorType string) {
	scenarioErrors.WithLabelValues(sm.scenarioType, errorType).Inc()
}

// UpdateCPUUsage updates the CPU usage metric
func (sm *ScenarioMetrics) UpdateCPUUsage(percentage float64) {
	cpuUsage.WithLabelValues(sm.scenarioType).Set(percentage)
}

// UpdateMemoryUsage updates the memory usage metric
func (sm *ScenarioMetrics) UpdateMemoryUsage(bytes int64) {
	memoryUsage.WithLabelValues(sm.scenarioType).Set(float64(bytes))
}

// RecordDiskIO records disk I/O operations
func (sm *ScenarioMetrics) RecordDiskIO(bytes int64, operationType string) {
	diskIO.WithLabelValues(sm.scenarioType, operationType).Add(float64(bytes))
}

// RecordNetworkLatency records network latency
func (sm *ScenarioMetrics) RecordNetworkLatency(duration time.Duration) {
	networkLatency.WithLabelValues(sm.scenarioType).Observe(duration.Seconds())
}

// RecordNetworkError records a network error
func (sm *ScenarioMetrics) RecordNetworkError(errorType string) {
	networkErrors.WithLabelValues(sm.scenarioType, errorType).Inc()
}

// UpdateCircuitBreakerState updates the circuit breaker state
func (sm *ScenarioMetrics) UpdateCircuitBreakerState(state int) {
	circuitBreakerState.WithLabelValues(sm.scenarioType).Set(float64(state))
}

// RecordCircuitBreakerFailure records a circuit breaker failure
func (sm *ScenarioMetrics) RecordCircuitBreakerFailure() {
	circuitBreakerFailures.WithLabelValues(sm.scenarioType).Inc()
}

// RecordRateLimitHit records a rate limit hit
func (sm *ScenarioMetrics) RecordRateLimitHit() {
	rateLimitHits.WithLabelValues(sm.scenarioType).Inc()
}

// UpdateRateLimit updates the current rate limit
func (sm *ScenarioMetrics) UpdateRateLimit(limit float64) {
	rateLimitCurrent.WithLabelValues(sm.scenarioType).Set(limit)
}

// HTTPMetricsMiddleware returns a Gin middleware for HTTP metrics
func HTTPMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Seconds()
		requestDuration.WithLabelValues(c.Request.URL.Path, c.Request.Method, strconv.Itoa(c.Writer.Status())).Observe(duration)
		requestTotal.WithLabelValues(c.Request.URL.Path, c.Request.Method).Inc()
		if c.Writer.Status() >= 400 {
			errorTotal.WithLabelValues(c.Request.URL.Path, c.Request.Method).Inc()
		}
	}
}

// UpdateResourceMetrics updates CPU, memory, and disk I/O metrics
func UpdateResourceMetrics(cpuBytes, memoryBytes, diskBytes int64) {
	cpuUsage.WithLabelValues("").Set(float64(cpuBytes))
	memoryUsage.WithLabelValues("").Set(float64(memoryBytes))
	diskIO.WithLabelValues("", "read").Add(float64(diskBytes))
}

// UpdateCircuitBreakerState updates the circuit breaker state metric
func UpdateCircuitBreakerState(state string) {
	value := 0
	switch state {
	case "open":
		value = 1
	case "half-open":
		value = 2
	}
	circuitBreakerState.WithLabelValues("").Set(float64(value))
}

// UpdateCircuitBreakerFailures increments the circuit breaker failures counter
func UpdateCircuitBreakerFailures() {
	circuitBreakerFailures.WithLabelValues("").Inc()
}

// UpdateRateLimitMetrics updates rate limit metrics
func UpdateRateLimitMetrics(hits, current int64) {
	rateLimitHits.WithLabelValues("").Add(float64(hits))
	rateLimitCurrent.WithLabelValues("").Set(float64(current))
}

// MetricsContextKey is the key used to store metrics context in context.Context
type MetricsContextKey struct{}

var metricsContextKey = MetricsContextKey{}

// MetricsContext holds metrics-related data
type MetricsContext struct {
	scenario *ScenarioMetrics
}

// WithMetricsContext adds metrics context to the given context
func WithMetricsContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, metricsContextKey, &MetricsContext{})
}

// TrackScenario creates a new scenario metrics tracker
func (mc *MetricsContext) TrackScenario(name string) *ScenarioMetrics {
	mc.scenario = NewScenarioMetrics(name)
	return mc.scenario
}

// TrackResources updates resource metrics
func (mc *MetricsContext) TrackResources(cpuBytes, memoryBytes, diskBytes int64) {
	UpdateResourceMetrics(cpuBytes, memoryBytes, diskBytes)
}
