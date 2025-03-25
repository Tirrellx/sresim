package simulator

import (
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type ScenarioManager struct {
	activeScenarios map[string]bool
	stopChannels    map[string]chan struct{}
	mu              sync.RWMutex
}

var manager = &ScenarioManager{
	activeScenarios: make(map[string]bool),
	stopChannels:    make(map[string]chan struct{}),
}

// StartLatencySimulation simulates high network latency
func (sm *ScenarioManager) StartLatencySimulation(delayMs int) {
	sm.mu.Lock()
	sm.activeScenarios["latency"] = true
	stopCh := make(chan struct{})
	sm.stopChannels["latency"] = stopCh
	sm.mu.Unlock()

	go func() {
		for {
			select {
			case <-stopCh:
				return
			default:
				time.Sleep(time.Duration(delayMs) * time.Millisecond)
			}
		}
	}()
}

// StartErrorRateSimulation simulates high error rate
func (sm *ScenarioManager) StartErrorRateSimulation(errorPercentage int) {
	sm.mu.Lock()
	sm.activeScenarios["error_rate"] = true
	stopCh := make(chan struct{})
	sm.stopChannels["error_rate"] = stopCh
	sm.mu.Unlock()

	go func() {
		for {
			select {
			case <-stopCh:
				return
			default:
				if rand.Intn(100) < errorPercentage {
					// Simulate error by consuming CPU
					time.Sleep(time.Millisecond * 100)
				}
			}
		}
	}()
}

// StartResourceExhaustionSimulation simulates CPU and memory exhaustion
func (sm *ScenarioManager) StartResourceExhaustionSimulation(cpuPercentage, memoryPercentage int) {
	sm.mu.Lock()
	sm.activeScenarios["resource_exhaustion"] = true
	stopCh := make(chan struct{})
	sm.stopChannels["resource_exhaustion"] = stopCh
	sm.mu.Unlock()

	go func() {
		// Allocate memory
		memorySize := int(float64(memoryPercentage) / 100.0 * 1024 * 1024 * 1024) // GB
		memory := make([]byte, memorySize)

		// CPU intensive loop
		for {
			select {
			case <-stopCh:
				return
			default:
				if rand.Float64()*100 < float64(cpuPercentage) {
					// Consume CPU
					runtime.Gosched()
				}
			}
		}
	}()
}

// StartCircuitBreakerSimulation simulates circuit breaker pattern
func (sm *ScenarioManager) StartCircuitBreakerSimulation(threshold int, timeoutSeconds int) {
	sm.mu.Lock()
	sm.activeScenarios["circuit_breaker"] = true
	stopCh := make(chan struct{})
	sm.stopChannels["circuit_breaker"] = stopCh
	sm.mu.Unlock()

	go func() {
		failures := 0
		lastFailure := time.Now()

		for {
			select {
			case <-stopCh:
				return
			default:
				if time.Since(lastFailure).Seconds() > float64(timeoutSeconds) {
					failures = 0
				}
				if failures >= threshold {
					time.Sleep(time.Second)
				}
			}
		}
	}()
}

// StartRateLimitSimulation simulates rate limiting
func (sm *ScenarioManager) StartRateLimitSimulation(requestsPerSecond int) {
	sm.mu.Lock()
	sm.activeScenarios["rate_limit"] = true
	stopCh := make(chan struct{})
	sm.stopChannels["rate_limit"] = stopCh
	sm.mu.Unlock()

	go func() {
		interval := time.Second / time.Duration(requestsPerSecond)
		for {
			select {
			case <-stopCh:
				return
			default:
				time.Sleep(interval)
			}
		}
	}()
}

// StartNetworkPartitionSimulation simulates network partition
func (sm *ScenarioManager) StartNetworkPartitionSimulation(durationSeconds int) {
	sm.mu.Lock()
	sm.activeScenarios["network_partition"] = true
	stopCh := make(chan struct{})
	sm.stopChannels["network_partition"] = stopCh
	sm.mu.Unlock()

	go func() {
		time.Sleep(time.Duration(durationSeconds) * time.Second)
		sm.StopScenario("network_partition")
	}()
}

// StartMemoryLeakSimulation simulates memory leak
func (sm *ScenarioManager) StartMemoryLeakSimulation(leakRateMB int, durationSeconds int) {
	sm.mu.Lock()
	sm.activeScenarios["memory_leak"] = true
	stopCh := make(chan struct{})
	sm.stopChannels["memory_leak"] = stopCh
	sm.mu.Unlock()

	go func() {
		leakInterval := time.Second / time.Duration(leakRateMB)
		leakSize := 1024 * 1024 // 1MB
		leakedMemory := make([][]byte, 0)

		for {
			select {
			case <-stopCh:
				return
			case <-time.After(time.Duration(durationSeconds) * time.Second):
				sm.StopScenario("memory_leak")
				return
			default:
				leakedMemory = append(leakedMemory, make([]byte, leakSize))
				time.Sleep(leakInterval)
			}
		}
	}()
}

// StartCPUSpikeSimulation simulates CPU spikes
func (sm *ScenarioManager) StartCPUSpikeSimulation(spikePercentage int, durationSeconds int, intervalSeconds int) {
	sm.mu.Lock()
	sm.activeScenarios["cpu_spike"] = true
	stopCh := make(chan struct{})
	sm.stopChannels["cpu_spike"] = stopCh
	sm.mu.Unlock()

	go func() {
		for {
			select {
			case <-stopCh:
				return
			default:
				// Create CPU spike
				for i := 0; i < runtime.NumCPU(); i++ {
					go func() {
						for {
							if rand.Float64()*100 < float64(spikePercentage) {
								runtime.Gosched()
							}
						}
					}()
				}
				time.Sleep(time.Duration(durationSeconds) * time.Second)
				time.Sleep(time.Duration(intervalSeconds-durationSeconds) * time.Second)
			}
		}
	}()
}

// StartDiskIOSimulation simulates disk I/O saturation
func (sm *ScenarioManager) StartDiskIOSimulation(opsPerSecond int, fileSizeMB int) {
	sm.mu.Lock()
	sm.activeScenarios["disk_io"] = true
	stopCh := make(chan struct{})
	sm.stopChannels["disk_io"] = stopCh
	sm.mu.Unlock()

	go func() {
		interval := time.Second / time.Duration(opsPerSecond)
		data := make([]byte, fileSizeMB*1024*1024)
		rand.Read(data)

		for {
			select {
			case <-stopCh:
				return
			default:
				// Perform random I/O operations
				offset := rand.Int63n(int64(len(data)))
				length := rand.Intn(4096) + 1
				_ = data[offset : offset+int64(length)]
				time.Sleep(interval)
			}
		}
	}()
}

// StartConnectionPoolExhaustionSimulation simulates connection pool exhaustion
func (sm *ScenarioManager) StartConnectionPoolExhaustionSimulation(maxConnections int, holdTimeSeconds int) {
	sm.mu.Lock()
	sm.activeScenarios["connection_pool_exhaustion"] = true
	stopCh := make(chan struct{})
	sm.stopChannels["connection_pool_exhaustion"] = stopCh
	sm.mu.Unlock()

	go func() {
		connections := make([]chan struct{}, maxConnections)
		for i := 0; i < maxConnections; i++ {
			connections[i] = make(chan struct{})
		}

		for {
			select {
			case <-stopCh:
				return
			default:
				for _, conn := range connections {
					select {
					case conn <- struct{}{}:
						time.Sleep(time.Duration(holdTimeSeconds) * time.Second)
					default:
						// Connection pool is full
						time.Sleep(time.Millisecond * 100)
					}
				}
			}
		}
	}()
}

// StartCascadingFailureSimulation simulates cascading failures
func (sm *ScenarioManager) StartCascadingFailureSimulation(chainLength int, delaySeconds int) {
	sm.mu.Lock()
	sm.activeScenarios["cascading_failure"] = true
	stopCh := make(chan struct{})
	sm.stopChannels["cascading_failure"] = stopCh
	sm.mu.Unlock()

	go func() {
		for i := 0; i < chainLength; i++ {
			select {
			case <-stopCh:
				return
			default:
				// Simulate service failure
				time.Sleep(time.Duration(delaySeconds) * time.Second)
				// Trigger cascading effect
				runtime.Gosched()
			}
		}
		sm.StopScenario("cascading_failure")
	}()
}

// StartThunderingHerdSimulation simulates thundering herd problem
func (sm *ScenarioManager) StartThunderingHerdSimulation(concurrentRequests int, cacheMissPercentage int) {
	sm.mu.Lock()
	sm.activeScenarios["thundering_herd"] = true
	stopCh := make(chan struct{})
	sm.stopChannels["thundering_herd"] = stopCh
	sm.mu.Unlock()

	go func() {
		var wg sync.WaitGroup
		for {
			select {
			case <-stopCh:
				return
			default:
				for i := 0; i < concurrentRequests; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						if rand.Intn(100) < cacheMissPercentage {
							// Simulate cache miss and heavy computation
							time.Sleep(time.Millisecond * 100)
						}
					}()
				}
				wg.Wait()
				time.Sleep(time.Millisecond * 10)
			}
		}
	}()
}

// StopScenario stops a running simulation scenario
func (sm *ScenarioManager) StopScenario(scenarioName string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if stopCh, exists := sm.stopChannels[scenarioName]; exists {
		close(stopCh)
		delete(sm.stopChannels, scenarioName)
		delete(sm.activeScenarios, scenarioName)
	}
}

// IsScenarioActive checks if a scenario is currently running
func (sm *ScenarioManager) IsScenarioActive(scenarioName string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.activeScenarios[scenarioName]
}

// GetManager returns the singleton scenario manager
func GetManager() *ScenarioManager {
	return manager
}
