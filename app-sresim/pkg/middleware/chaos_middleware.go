package middleware

import (
	"net/http"
	"time"

	"github.com/tirrellx/sresim/pkg/chaos"
)

// ChaosMiddleware intercepts HTTP requests and applies chaos
// by randomly failing or delaying the request.
func ChaosMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If chaos decides to fail, send an error response.
		if chaos.ShouldFail() {
			http.Error(w, "Simulated failure", http.StatusInternalServerError)
			return
		}

		// If chaos decides to delay, pause the request processing.
		if chaos.ShouldDelay() {
			time.Sleep(chaos.RandomDelay())
		}

		// Proceed with the next handler.
		next.ServeHTTP(w, r)
	})
}
