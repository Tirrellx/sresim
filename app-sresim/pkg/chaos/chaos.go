package chaos

import (
	"math/rand"
	"time"
)

// init seeds the random number generator.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// ShouldFail randomly returns true based on a set probability (e.g., 20%).
func ShouldFail() bool {
	return rand.Float32() < 0.2
}

// ShouldDelay randomly decides to delay the response (e.g., 30% chance).
func ShouldDelay() bool {
	return rand.Float32() < 0.3
}

// RandomDelay returns a random delay duration between 100ms and 1s.
func RandomDelay() time.Duration {
	return time.Duration(100+rand.Intn(900)) * time.Millisecond
}
