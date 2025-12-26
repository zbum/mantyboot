package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/zbum/mantyboot/http/mux"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) isAllowed(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Clean old requests
	if times, exists := rl.requests[key]; exists {
		var validTimes []time.Time
		for _, t := range times {
			if t.After(windowStart) {
				validTimes = append(validTimes, t)
			}
		}
		rl.requests[key] = validTimes
	}

	// Check if limit exceeded
	if len(rl.requests[key]) >= rl.limit {
		return false
	}

	// Add current request
	rl.requests[key] = append(rl.requests[key], now)
	return true
}

func RateLimit(limiter *RateLimiter, keyFunc func(*http.Request) string) mux.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)

			if !limiter.isAllowed(key) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next(w, r)
		}
	}
}

// Default key functions
func IPKeyFunc(r *http.Request) string {
	return r.RemoteAddr
}

func UserAgentKeyFunc(r *http.Request) string {
	return r.UserAgent()
}

func HeaderKeyFunc(headerName string) func(*http.Request) string {
	return func(r *http.Request) string {
		return r.Header.Get(headerName)
	}
}

// Convenience function for IP-based rate limiting
func RateLimitByIP(limit int, window time.Duration) mux.Middleware {
	limiter := NewRateLimiter(limit, window)
	return RateLimit(limiter, IPKeyFunc)
}
