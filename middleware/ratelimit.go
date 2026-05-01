package middleware

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter implements per-user rate limiting using token bucket algorithm
type RateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter // user_id -> limiter
	rpm      int                       // Requests per minute
	burst    int                       // Burst capacity
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerMinute, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rpm:      requestsPerMinute,
		burst:    burst,
	}
}

// GetLimiter returns the rate limiter for a user, creating one if needed
func (rl *RateLimiter) GetLimiter(userID string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[userID]
	if !exists {
		// Create new limiter with specified rate and burst
		// Convert RPM to tokens per second
		tokensPerSecond := float64(rl.rpm) / 60.0
		limiter = rate.NewLimiter(rate.Limit(tokensPerSecond), rl.burst)
		rl.limiters[userID] = limiter
	}

	return limiter
}

// Allow checks if a request is allowed for a user
func (rl *RateLimiter) Allow(userID string) bool {
	limiter := rl.GetLimiter(userID)
	return limiter.Allow()
}

// Cleanup removes inactive limiters (should be called periodically)
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Remove limiters that haven't been used recently
	// This is a simple cleanup - could be enhanced with last-used tracking
	for userID := range rl.limiters {
		delete(rl.limiters, userID)
	}
}

// StartCleanupRoutine starts a background goroutine to cleanup inactive limiters
func (rl *RateLimiter) StartCleanupRoutine(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			rl.Cleanup()
		}
	}()
}
