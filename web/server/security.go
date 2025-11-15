// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// Default max request body size: 10MB
	defaultMaxBodySize = 10 * 1024 * 1024
	// Default request timeout: 30 seconds
	defaultRequestTimeout = 30 * time.Second
)

// Security configuration
type securityConfig struct {
	requireAuth    bool
	authToken      string
	allowedOrigins []string
	maxBodySize    int64
	requestTimeout time.Duration
}

func loadSecurityConfig() securityConfig {
	cfg := securityConfig{
		requireAuth:    os.Getenv("REQUIRE_AUTH") == "true",
		authToken:      os.Getenv("AUTH_TOKEN"),
		maxBodySize:    defaultMaxBodySize,
		requestTimeout: defaultRequestTimeout,
	}

	// Parse allowed origins
	if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
		cfg.allowedOrigins = strings.Split(origins, ",")
		for i := range cfg.allowedOrigins {
			cfg.allowedOrigins[i] = strings.TrimSpace(cfg.allowedOrigins[i])
		}
	}

	// Parse max body size
	if sizeStr := os.Getenv("MAX_BODY_SIZE"); sizeStr != "" {
		var size int64
		if _, err := fmt.Sscanf(sizeStr, "%d", &size); err == nil && size > 0 {
			cfg.maxBodySize = size
		}
	}

	// Parse request timeout
	if timeoutStr := os.Getenv("REQUEST_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil && timeout > 0 {
			cfg.requestTimeout = timeout
		}
	}

	return cfg
}

// Middleware to enforce request size limits
func maxBodySizeMiddleware(maxSize int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxSize {
				respondError(w, http.StatusRequestEntityTooLarge, "Request body too large")
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, maxSize)
			next.ServeHTTP(w, r)
		})
	}
}

// Middleware to enforce request timeouts
func timeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Middleware for CORS
func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if wildcard is enabled (allow all origins)
			allowAll := false
			for _, allowedOrigin := range allowedOrigins {
				if allowedOrigin == "*" {
					allowAll = true
					break
				}
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				if allowAll {
					// Allow all origins
					if origin != "" {
						w.Header().Set("Access-Control-Allow-Origin", origin)
					} else {
						w.Header().Set("Access-Control-Allow-Origin", "*")
					}
				} else if len(allowedOrigins) > 0 {
					// Check if origin is allowed
					allowed := false
					for _, allowedOrigin := range allowedOrigins {
						if origin == allowedOrigin {
							allowed = true
							break
						}
					}
					if !allowed {
						http.Error(w, "CORS policy violation", http.StatusForbidden)
						return
					}
					w.Header().Set("Access-Control-Allow-Origin", origin)
				} else {
					// No CORS allowed by default
					http.Error(w, "CORS not allowed", http.StatusForbidden)
					return
				}
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Max-Age", "3600")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// Handle actual requests
			if allowAll {
				// Allow all origins
				if origin != "" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				} else {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				}
			} else if len(allowedOrigins) > 0 && origin != "" {
				allowed := false
				for _, allowedOrigin := range allowedOrigins {
					if origin == allowedOrigin {
						allowed = true
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
				if !allowed {
					http.Error(w, "CORS policy violation", http.StatusForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Middleware for optional token authentication
func authMiddleware(requireAuth bool, authToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !requireAuth {
				next.ServeHTTP(w, r)
				return
			}

			if authToken == "" {
				log.Printf("[security] AUTH_TOKEN not set but REQUIRE_AUTH=true")
				respondError(w, http.StatusInternalServerError, "Authentication misconfigured")
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			// Expect "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondError(w, http.StatusUnauthorized, "Invalid authorization format")
				return
			}

			if parts[1] != authToken {
				logSecurityEvent("auth_failed", r.URL.Path, "invalid_token")
				respondError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			logSecurityEvent("auth_success", r.URL.Path, "")
			next.ServeHTTP(w, r)
		})
	}
}

// Sanitized logging: log security events without sensitive data
func logSecurityEvent(event, path, detail string) {
	log.Printf("[security] event=%s path=%s detail=%s", event, path, detail)
}

// Sanitize error message: remove sensitive data
func sanitizeError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	// Truncate long messages (might contain keys/sigs)
	if len(msg) > 200 {
		msg = msg[:200] + "..."
	}
	return msg
}

