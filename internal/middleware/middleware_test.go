// Package middleware_test provides tests for the middleware package.
package middleware_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/agentapi/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestApplyDefaultStack tests applying the default middleware stack.
func TestApplyDefaultStack(t *testing.T) {
	router := chi.NewRouter()
	err := middleware.ApplyDefaultStack(router)
	assert.NoError(t, err)

	// Register a simple test route
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Test the route works
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

// TestApplyCustomCORS tests applying custom CORS middleware.
func TestApplyCustomCORS(t *testing.T) {
	router := chi.NewRouter()
	options := middleware.CORSOptions{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedHosts:   []string{"localhost"},
	}
	middleware.ApplyCustomCORS(router, options)

	// Register a simple test route
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Test the route with CORS headers
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestHealthCheckRoute tests the health check endpoint.
func TestHealthCheckRoute(t *testing.T) {
	router := chi.NewRouter()
	middleware.HealthCheckRoute(router)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "ok")
}

// TestReadinessCheckRoute tests the readiness check endpoint.
func TestReadinessCheckRoute(t *testing.T) {
	router := chi.NewRouter()
	middleware.ReadinessCheckRoute(router)

	req := httptest.NewRequest("GET", "/readiness", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "ready")
}

// TestNewRequestIDHandler tests creating a new RequestIDHandler.
func TestNewRequestIDHandler(t *testing.T) {
	timeout := 5 * time.Second
	handler := middleware.NewRequestIDHandler(timeout)
	require.NotNil(t, handler)
}

// TestRequestIDHandler_WrapHandler tests wrapping a handler with RequestIDHandler.
func TestRequestIDHandler_WrapHandler(t *testing.T) {
	timeout := 5 * time.Second
	handler := middleware.NewRequestIDHandler(timeout)

	// Create a simple http.Handler to wrap
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	wrapped := handler.WrapHandler(inner)
	require.NotNil(t, wrapped)

	// Test the wrapped handler works
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

// TestRequestIDHandler_WrapHandler_Timeout tests that the wrapper handles timeouts.
func TestRequestIDHandler_WrapHandler_Timeout(t *testing.T) {
	timeout := 1 * time.Millisecond
	handler := middleware.NewRequestIDHandler(timeout)

	// Create a handler that takes longer than the timeout
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	wrapped := handler.WrapHandler(inner)
	require.NotNil(t, wrapped)

	// Test the wrapped handler times out
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	wrapped.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

// TestCORSOptions tests CORSOptions struct.
func TestCORSOptions(t *testing.T) {
	options := middleware.CORSOptions{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedHosts:   []string{"localhost", "127.0.0.1"},
	}

	assert.Len(t, options.AllowedOrigins, 2)
	assert.Len(t, options.AllowedHosts, 2)
}

// TestHealthAndReadinessEndpoints tests both health and readiness endpoints on same router.
func TestHealthAndReadinessEndpoints(t *testing.T) {
	router := chi.NewRouter()
	middleware.HealthCheckRoute(router)
	middleware.ReadinessCheckRoute(router)

	// Test health endpoint
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test readiness endpoint
	req = httptest.NewRequest("GET", "/readiness", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestApplyDefaultStackAndCORS tests combining default stack and CORS.
func TestApplyDefaultStackAndCORS(t *testing.T) {
	router := chi.NewRouter()
	err := middleware.ApplyDefaultStack(router)
	assert.NoError(t, err)

	options := middleware.CORSOptions{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedHosts:   []string{"localhost"},
	}
	middleware.ApplyCustomCORS(router, options)

	router.Get("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "API Response")
	})

	// Test the combined middleware stack
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "API Response", w.Body.String())
}
