package main

import (
	"log"
	"net/http"

	"github.com/localstack/sresim/app-sresim/pkg/handlers"
	"github.com/localstack/sresim/app-sresim/pkg/metrics"
	"github.com/localstack/sresim/app-sresim/pkg/middleware"
	"github.com/localstack/sresim/app-sresim/pkg/simulator"
)

func main() {
	// Initialize metrics
	if err := metrics.InitMetrics(); err != nil {
		log.Fatalf("Failed to initialize metrics: %v", err)
	}

	// Create a new HTTP multiplexer
	mux := http.NewServeMux()

	// Define endpoints
	mux.HandleFunc("/simulate", handlers.SimulateHandler)
	mux.HandleFunc("/health", handlers.HealthCheckHandler)
	mux.Handle("/metrics", metrics.MetricsHandler())

	// Simulation endpoints
	mux.HandleFunc("/scenarios", simulator.ListScenarios)
	mux.HandleFunc("/scenarios/run", simulator.RunScenario)
	mux.HandleFunc("/scenarios/stop", simulator.StopScenario)

	// Wrap the multiplexer with our middlewares
	handler := middleware.ChaosMiddleware(metrics.MetricsMiddleware(mux))

	// Start the HTTP server
	log.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
