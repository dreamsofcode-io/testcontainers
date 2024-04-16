package main

import (
	"net"
	"net/http"
	"strconv"

	"github.com/dreamsofcode-io/spellbook/ratelimit"
)

func main() {
	r := http.NewServeMux()

	// Simple health endpoint
	r.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Find All
	r.HandleFunc("GET /spells", func(w http.ResponseWriter, r *http.Request) {
	})

	// Create
	r.HandleFunc("POST /spells", func(w http.ResponseWriter, r *http.Request) {
	})

	// FindByID
	r.HandleFunc("GET /spells/{id}", func(w http.ResponseWriter, r *http.Request) {
	})

	// Update
	r.HandleFunc("PUT /spells/{id}", func(w http.ResponseWriter, r *http.Request) {
	})

	// Delete
	r.HandleFunc("DELETE /spells/{id}", func(w http.ResponseWriter, r *http.Request) {
	})
}

func rateLimiterMiddleware(limiter *ratelimit.RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := net.ParseIP(r.RemoteAddr)

		info, err := limiter.AddAndCheckIfExceeds(r.Context(), ip)
		if err != nil {
			// if an error occurs, lets just continue
			next.ServeHTTP(w, r)
		}

		w.Header().Add("x-ratelimit-limit", strconv.Itoa(int(info.Limit())))
		w.Header().Add("x-ratelimit-reset", strconv.Itoa(int(info.Resets().Seconds())))
		w.Header().Add("x-ratelimit-remaining", strconv.Itoa(int(info.Remaining())))

		if info.IsExceeded() {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
