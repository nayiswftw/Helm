package api

import (
	"crypto/subtle"
	"log/slog"
	"net/http"
	"strings"
)

// apiKeyAuth returns middleware that validates API key authentication.
// The key can be provided via "Authorization: Bearer <token>" header
// or via "X-API-Key: <token>" header (for ESP32 simplicity).
//
// If apiKey is empty, the middleware is a no-op passthrough (development mode).
func apiKeyAuth(apiKey string, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// No API key configured — skip authentication (dev mode).
		if apiKey == "" {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)

			if token == "" {
				logger.Warn("unauthenticated request",
					"method", r.Method,
					"path", r.URL.Path,
					"remote", r.RemoteAddr,
				)
				respondError(w, http.StatusUnauthorized, "unauthorized", "API key is required")
				return
			}

			if subtle.ConstantTimeCompare([]byte(token), []byte(apiKey)) != 1 {
				logger.Warn("invalid API key",
					"method", r.Method,
					"path", r.URL.Path,
					"remote", r.RemoteAddr,
				)
				respondError(w, http.StatusUnauthorized, "unauthorized", "Invalid API key")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// extractToken retrieves the API key from the request.
// Checks Authorization header first (Bearer scheme), then X-API-Key header.
func extractToken(r *http.Request) string {
	// Try Authorization: Bearer <token>
	if auth := r.Header.Get("Authorization"); auth != "" {
		if strings.HasPrefix(auth, "Bearer ") {
			return strings.TrimPrefix(auth, "Bearer ")
		}
	}

	// Try X-API-Key: <token>
	if key := r.Header.Get("X-API-Key"); key != "" {
		return key
	}

	return ""
}
