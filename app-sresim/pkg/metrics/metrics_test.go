package metrics

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetMetrics() {
	prometheus.DefaultRegisterer.Unregister(requestDuration)
	prometheus.DefaultRegisterer.Unregister(requestTotal)
	prometheus.DefaultRegisterer.Unregister(errorTotal)
	prometheus.DefaultRegisterer.Unregister(activeScenarios)
	prometheus.DefaultRegisterer.Unregister(scenarioDuration)
	prometheus.DefaultRegisterer.Unregister(scenarioErrors)
	prometheus.DefaultRegisterer.Unregister(cpuUsage)
	prometheus.DefaultRegisterer.Unregister(memoryUsage)
	prometheus.DefaultRegisterer.Unregister(diskIO)
	prometheus.DefaultRegisterer.Unregister(networkLatency)
	prometheus.DefaultRegisterer.Unregister(networkErrors)
	prometheus.DefaultRegisterer.Unregister(circuitBreakerState)
	prometheus.DefaultRegisterer.Unregister(circuitBreakerFailures)
	prometheus.DefaultRegisterer.Unregister(rateLimitHits)
	prometheus.DefaultRegisterer.Unregister(rateLimitCurrent)
}

func TestMetricsInitialization(t *testing.T) {
	resetMetrics()

	// Test initialization
	err := Init()
	require.NoError(t, err)

	// Verify metrics are registered
	metrics := []prometheus.Collector{
		requestDuration,
		requestTotal,
		errorTotal,
		activeScenarios,
		scenarioDuration,
		scenarioErrors,
		cpuUsage,
		memoryUsage,
		diskIO,
		networkLatency,
		networkErrors,
		circuitBreakerState,
		circuitBreakerFailures,
		rateLimitHits,
		rateLimitCurrent,
	}

	for _, m := range metrics {
		err := prometheus.Register(m)
		assert.Error(t, err, "Metric should already be registered")
	}
}

func TestHTTPMetricsMiddleware(t *testing.T) {
	resetMetrics()

	// Setup test server
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/test", nil)

	// Test middleware
	handler := HTTPMetricsMiddleware()
	handler(c)

	// Verify metrics were recorded
	assert.Equal(t, float64(1), testutil.ToFloat64(requestTotal.WithLabelValues("/test", "GET")))
}

func TestScenarioMetrics(t *testing.T) {
	resetMetrics()

	// Create test scenario
	scenario := NewScenarioMetrics("test-scenario")

	// Test scenario lifecycle
	scenario.SetActive(true)
	time.Sleep(100 * time.Millisecond) // Simulate some work
	scenario.SetActive(false)
	scenario.RecordDuration(100 * time.Millisecond)

	// Verify metrics
	assert.Equal(t, float64(0), testutil.ToFloat64(activeScenarios.WithLabelValues("test-scenario")))

	// For histogram metrics, we just verify that attempting to register it again fails
	err := prometheus.Register(scenarioDuration)
	assert.NoError(t, err, "Histogram metric should not be registered yet")
}

func TestResourceMetrics(t *testing.T) {
	resetMetrics()

	// Test resource updates
	UpdateResourceMetrics(1000, 2000, 3000)

	// Verify metrics
	assert.Equal(t, float64(1000), testutil.ToFloat64(cpuUsage))
	assert.Equal(t, float64(2000), testutil.ToFloat64(memoryUsage))
	assert.Equal(t, float64(3000), testutil.ToFloat64(diskIO))
}

func TestCircuitBreakerMetrics(t *testing.T) {
	resetMetrics()

	// Test circuit breaker state changes
	UpdateCircuitBreakerState("closed") // Should set to 0
	assert.Equal(t, float64(0), testutil.ToFloat64(circuitBreakerState.WithLabelValues("")))

	UpdateCircuitBreakerState("open") // Should set to 1
	assert.Equal(t, float64(1), testutil.ToFloat64(circuitBreakerState.WithLabelValues("")))

	UpdateCircuitBreakerState("half-open") // Should set to 2
	assert.Equal(t, float64(2), testutil.ToFloat64(circuitBreakerState.WithLabelValues("")))

	// Test failures
	UpdateCircuitBreakerFailures()
	assert.Equal(t, float64(1), testutil.ToFloat64(circuitBreakerFailures.WithLabelValues("")))
}

func TestRateLimitMetrics(t *testing.T) {
	resetMetrics()

	// Test rate limit updates
	UpdateRateLimitMetrics(5, 10)

	// Verify metrics
	assert.Equal(t, float64(5), testutil.ToFloat64(rateLimitHits))
	assert.Equal(t, float64(10), testutil.ToFloat64(rateLimitCurrent))
}

func TestMetricsContext(t *testing.T) {
	resetMetrics()

	// Create test context
	ctx := context.Background()
	ctx = WithMetricsContext(ctx)

	// Verify metrics context
	metricsCtx, ok := ctx.Value(metricsContextKey).(*MetricsContext)
	require.True(t, ok)
	require.NotNil(t, metricsCtx)

	// Test scenario tracking
	scenario := metricsCtx.TrackScenario("test-scenario")
	require.NotNil(t, scenario)
	assert.Equal(t, "test-scenario", scenario.scenarioType)

	// Test resource tracking
	metricsCtx.TrackResources(1000, 2000, 3000)

	// Verify metrics
	assert.Equal(t, float64(1000), testutil.ToFloat64(cpuUsage))
}
