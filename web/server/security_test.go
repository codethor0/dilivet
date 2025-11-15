// DiliVet – ML-DSA diagnostics and vetting toolkit
// Author: Thor "Thor Thor" (codethor@gmail.com, https://www.linkedin.com/in/thor-thor0)

package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMaxBodySizeMiddleware(t *testing.T) {
	handler := maxBodySizeMiddleware(100)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	t.Run("small body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader("small"))
		req.ContentLength = 5
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("large body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", strings.NewReader(strings.Repeat("a", 200)))
		req.ContentLength = 200
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("Expected 413, got %d", w.Code)
		}
	})
}

func TestCORSMiddleware(t *testing.T) {
	t.Run("no origins configured", func(t *testing.T) {
		handler := corsMiddleware([]string{})(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "https://evil.com")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		// Should pass through (no CORS headers, but not blocked)
		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("allowed origin", func(t *testing.T) {
		handler := corsMiddleware([]string{"https://example.com"})(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
		if w.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
			t.Errorf("Expected CORS header, got %s", w.Header().Get("Access-Control-Allow-Origin"))
		}
	})

	t.Run("preflight request", func(t *testing.T) {
		handler := corsMiddleware([]string{"https://example.com"})(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		)

		req := httptest.NewRequest("OPTIONS", "/", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusNoContent {
			t.Errorf("Expected 204, got %d", w.Code)
		}
		if w.Header().Get("Access-Control-Allow-Methods") == "" {
			t.Error("Expected CORS methods header")
		}
	})
}

func TestAuthMiddleware(t *testing.T) {
	t.Run("auth disabled", func(t *testing.T) {
		handler := authMiddleware(false, "token123")(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		)

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("auth enabled, no token", func(t *testing.T) {
		handler := authMiddleware(true, "token123")(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		)

		req := httptest.NewRequest("GET", "/api/test", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", w.Code)
		}
	})

	t.Run("auth enabled, correct token", func(t *testing.T) {
		handler := authMiddleware(true, "token123")(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		)

		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("Authorization", "Bearer token123")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("auth enabled, wrong token", func(t *testing.T) {
		handler := authMiddleware(true, "token123")(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		)

		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("Authorization", "Bearer wrongtoken")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", w.Code)
		}
	})

	t.Run("auth enabled, invalid format", func(t *testing.T) {
		handler := authMiddleware(true, "token123")(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
		)

		req := httptest.NewRequest("GET", "/api/test", nil)
		req.Header.Set("Authorization", "token123")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", w.Code)
		}
	})
}

func TestLoadSecurityConfig(t *testing.T) {
	// Save original env
	origRequireAuth := os.Getenv("REQUIRE_AUTH")
	origAuthToken := os.Getenv("AUTH_TOKEN")
	origOrigins := os.Getenv("ALLOWED_ORIGINS")

	defer func() {
		os.Setenv("REQUIRE_AUTH", origRequireAuth)
		os.Setenv("AUTH_TOKEN", origAuthToken)
		os.Setenv("ALLOWED_ORIGINS", origOrigins)
	}()

	t.Run("defaults", func(t *testing.T) {
		os.Unsetenv("REQUIRE_AUTH")
		os.Unsetenv("AUTH_TOKEN")
		os.Unsetenv("ALLOWED_ORIGINS")
		cfg := loadSecurityConfig()
		if cfg.requireAuth {
			t.Error("Expected requireAuth=false by default")
		}
		if len(cfg.allowedOrigins) != 0 {
			t.Error("Expected no allowed origins by default")
		}
	})

	t.Run("auth enabled", func(t *testing.T) {
		os.Setenv("REQUIRE_AUTH", "true")
		os.Setenv("AUTH_TOKEN", "testtoken")
		cfg := loadSecurityConfig()
		if !cfg.requireAuth {
			t.Error("Expected requireAuth=true")
		}
		if cfg.authToken != "testtoken" {
			t.Errorf("Expected authToken=testtoken, got %s", cfg.authToken)
		}
	})

	t.Run("allowed origins", func(t *testing.T) {
		os.Setenv("ALLOWED_ORIGINS", "https://example.com, https://test.com")
		cfg := loadSecurityConfig()
		if len(cfg.allowedOrigins) != 2 {
			t.Errorf("Expected 2 origins, got %d", len(cfg.allowedOrigins))
		}
		if cfg.allowedOrigins[0] != "https://example.com" {
			t.Errorf("Expected first origin=https://example.com, got %s", cfg.allowedOrigins[0])
		}
	})
}

func TestSanitizeError(t *testing.T) {
	t.Run("normal error", func(t *testing.T) {
		err := io.EOF
		sanitized := sanitizeError(err)
		if sanitized != "EOF" {
			t.Errorf("Expected 'EOF', got %s", sanitized)
		}
	})

	t.Run("long error", func(t *testing.T) {
		longMsg := strings.Repeat("a", 300)
		err := &testError{msg: longMsg}
		sanitized := sanitizeError(err)
		if len(sanitized) > 203 { // 200 + "..."
			t.Errorf("Expected truncated error, got length %d", len(sanitized))
		}
		if !strings.HasSuffix(sanitized, "...") {
			t.Error("Expected truncated error to end with '...'")
		}
	})
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestHexLengthValidation(t *testing.T) {
	t.Run("normal hex", func(t *testing.T) {
		_, err := decodeHex("deadbeef", "test")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("too long hex", func(t *testing.T) {
		longHex := strings.Repeat("a", 100001)
		_, err := decodeHex(longHex, "test")
		if err == nil {
			t.Error("Expected error for too long hex")
		}
		if !strings.Contains(err.Error(), "too long") {
			t.Errorf("Expected 'too long' error, got %v", err)
		}
	})
}

func TestLargeRequestBody(t *testing.T) {
	// Create a handler with max body size limit
	handler := maxBodySizeMiddleware(100)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to read body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(body)
		}),
	)

	t.Run("request exceeds limit", func(t *testing.T) {
		largeBody := bytes.NewReader(make([]byte, 200))
		req := httptest.NewRequest("POST", "/", largeBody)
		req.ContentLength = 200
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("Expected 413, got %d", w.Code)
		}
	})
}

