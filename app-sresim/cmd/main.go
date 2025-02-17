package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tirrellx/sresim/pkg/handlers"
	"github.com/tirrellx/sresim/pkg/middleware"
	"github.com/google/uuid"
)

func main() {
	// Create a new HTTP multiplexer.
	mux := http.NewServeMux()

	// Define the /simulate endpoint that your handler will process.
	mux.HandleFunc("/simulate", handlers.SimulateHandler)

	// Wrap the multiplexer with our chaos middleware.
	chaosHandler := middleware.ChaosMiddleware(mux)

	// Start the HTTP server.
	log.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", chaosHandler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
