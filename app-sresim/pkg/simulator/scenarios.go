package simulator

import (
	"encoding/json"
	"net/http"
	"time"
)

type Scenario struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type ScenarioResponse struct {
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Timestamp  time.Time              `json:"timestamp"`
	Parameters map[string]interface{} `json:"parameters"`
}

var scenarios = map[string]Scenario{
	"latency": {
		Name:        "High Latency",
		Description: "Simulates high network latency",
		Parameters: map[string]interface{}{
			"delay_ms": 1000,
		},
	},
	"error_rate": {
		Name:        "High Error Rate",
		Description: "Simulates high error rate",
		Parameters: map[string]interface{}{
			"error_percentage": 50,
		},
	},
	"resource_exhaustion": {
		Name:        "Resource Exhaustion",
		Description: "Simulates CPU and memory exhaustion",
		Parameters: map[string]interface{}{
			"cpu_percentage":    90,
			"memory_percentage": 85,
		},
	},
	"circuit_breaker": {
		Name:        "Circuit Breaker",
		Description: "Simulates circuit breaker pattern",
		Parameters: map[string]interface{}{
			"threshold": 5,
			"timeout":   30,
		},
	},
	"rate_limit": {
		Name:        "Rate Limiting",
		Description: "Simulates rate limiting",
		Parameters: map[string]interface{}{
			"requests_per_second": 10,
		},
	},
	"network_partition": {
		Name:        "Network Partition",
		Description: "Simulates network partition",
		Parameters: map[string]interface{}{
			"partition_duration": 60,
		},
	},
	"memory_leak": {
		Name:        "Memory Leak",
		Description: "Simulates memory leak scenario",
		Parameters: map[string]interface{}{
			"leak_rate_mb_per_second": 10,
			"duration_seconds":        300,
		},
	},
	"cpu_spike": {
		Name:        "CPU Spike",
		Description: "Simulates sudden CPU usage spikes",
		Parameters: map[string]interface{}{
			"spike_percentage": 95,
			"duration_seconds": 30,
			"interval_seconds": 60,
		},
	},
	"disk_io": {
		Name:        "Disk I/O Saturation",
		Description: "Simulates high disk I/O operations",
		Parameters: map[string]interface{}{
			"io_operations_per_second": 1000,
			"file_size_mb":             100,
		},
	},
	"connection_pool_exhaustion": {
		Name:        "Connection Pool Exhaustion",
		Description: "Simulates database connection pool exhaustion",
		Parameters: map[string]interface{}{
			"max_connections":   10,
			"hold_time_seconds": 30,
		},
	},
	"cascading_failure": {
		Name:        "Cascading Failure",
		Description: "Simulates cascading failure across services",
		Parameters: map[string]interface{}{
			"failure_chain_length":           3,
			"delay_between_failures_seconds": 5,
		},
	},
	"thundering_herd": {
		Name:        "Thundering Herd",
		Description: "Simulates thundering herd problem",
		Parameters: map[string]interface{}{
			"concurrent_requests":   100,
			"cache_miss_percentage": 80,
		},
	},
}

// ListScenarios returns all available simulation scenarios
func ListScenarios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scenarios)
}

// RunScenario executes a specific simulation scenario
func RunScenario(w http.ResponseWriter, r *http.Request) {
	scenarioName := r.URL.Query().Get("scenario")
	scenario, exists := scenarios[scenarioName]
	if !exists {
		http.Error(w, "Scenario not found", http.StatusNotFound)
		return
	}

	// Get the scenario manager
	manager := GetManager()

	// Start the appropriate scenario based on the name
	switch scenarioName {
	case "latency":
		delayMs := int(scenario.Parameters["delay_ms"].(float64))
		manager.StartLatencySimulation(delayMs)
	case "error_rate":
		errorPercentage := int(scenario.Parameters["error_percentage"].(float64))
		manager.StartErrorRateSimulation(errorPercentage)
	case "resource_exhaustion":
		cpuPercentage := int(scenario.Parameters["cpu_percentage"].(float64))
		memoryPercentage := int(scenario.Parameters["memory_percentage"].(float64))
		manager.StartResourceExhaustionSimulation(cpuPercentage, memoryPercentage)
	case "circuit_breaker":
		threshold := int(scenario.Parameters["threshold"].(float64))
		timeout := int(scenario.Parameters["timeout"].(float64))
		manager.StartCircuitBreakerSimulation(threshold, timeout)
	case "rate_limit":
		requestsPerSecond := int(scenario.Parameters["requests_per_second"].(float64))
		manager.StartRateLimitSimulation(requestsPerSecond)
	case "network_partition":
		duration := int(scenario.Parameters["partition_duration"].(float64))
		manager.StartNetworkPartitionSimulation(duration)
	case "memory_leak":
		leakRate := int(scenario.Parameters["leak_rate_mb_per_second"].(float64))
		duration := int(scenario.Parameters["duration_seconds"].(float64))
		manager.StartMemoryLeakSimulation(leakRate, duration)
	case "cpu_spike":
		spikePercentage := int(scenario.Parameters["spike_percentage"].(float64))
		duration := int(scenario.Parameters["duration_seconds"].(float64))
		interval := int(scenario.Parameters["interval_seconds"].(float64))
		manager.StartCPUSpikeSimulation(spikePercentage, duration, interval)
	case "disk_io":
		opsPerSecond := int(scenario.Parameters["io_operations_per_second"].(float64))
		fileSize := int(scenario.Parameters["file_size_mb"].(float64))
		manager.StartDiskIOSimulation(opsPerSecond, fileSize)
	case "connection_pool_exhaustion":
		maxConnections := int(scenario.Parameters["max_connections"].(float64))
		holdTime := int(scenario.Parameters["hold_time_seconds"].(float64))
		manager.StartConnectionPoolExhaustionSimulation(maxConnections, holdTime)
	case "cascading_failure":
		chainLength := int(scenario.Parameters["failure_chain_length"].(float64))
		delay := int(scenario.Parameters["delay_between_failures_seconds"].(float64))
		manager.StartCascadingFailureSimulation(chainLength, delay)
	case "thundering_herd":
		concurrentRequests := int(scenario.Parameters["concurrent_requests"].(float64))
		cacheMissPercentage := int(scenario.Parameters["cache_miss_percentage"].(float64))
		manager.StartThunderingHerdSimulation(concurrentRequests, cacheMissPercentage)
	}

	response := ScenarioResponse{
		Status:     "running",
		Message:    "Scenario started",
		Timestamp:  time.Now(),
		Parameters: scenario.Parameters,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StopScenario stops a running simulation scenario
func StopScenario(w http.ResponseWriter, r *http.Request) {
	scenarioName := r.URL.Query().Get("scenario")
	_, exists := scenarios[scenarioName]
	if !exists {
		http.Error(w, "Scenario not found", http.StatusNotFound)
		return
	}

	manager := GetManager()
	manager.StopScenario(scenarioName)

	response := ScenarioResponse{
		Status:    "stopped",
		Message:   "Scenario stopped",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
