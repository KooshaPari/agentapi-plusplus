// Package middleware provides standardized middleware for Phenotype Go services.
// This package consolidates CORS, rate limiting, JWT auth, and logging.
package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth/v5"
)

// Config holds middleware configuration
type Config struct {
	// CORS settings
	CORSAllowedOrigins   []string
	CORSAllowedMethods   []string
	CORSAllowedHeaders   []string
	CORSExposedHeaders   []string
	CORSAllowCredentials bool
	CORSMaxAge           int

	// Rate limiting settings
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// JWT settings
	JWTSecret string

	// Logging settings
	LogRequests bool
}

// DefaultConfig returns sensible defaults
func DefaultConfig() *Config {
	return &Config{
		// CORS defaults (allow all for development)
		CORSAllowedOrigins:   []string{"*"},
		CORSAllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		CORSAllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		CORSExposedHeaders:   []string{},
		CORSAllowCredentials: false,
		CORSMaxAge:           300,

		// Rate limiting defaults
		RateLimitRequests: 100,
		RateLimitWindow:   time.Minute,

		// JWT (empty by default, must be configured)
		JWTSecret: "",

		// Logging
		LogRequests: true,
	}
}

// NewCORS returns CORS middleware handler
func NewCORS(cfg *Config) func(http.Handler) http.Handler {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:     cfg.CORSAllowedOrigins,
		AllowedMethods:     cfg.CORSAllowedMethods,
		AllowedHeaders:     cfg.CORSAllowedHeaders,
		ExposedHeaders:     cfg.CORSExposedHeaders,
		AllowCredentials:   cfg.CORSAllowCredentials,
		MaxAge:             cfg.CORSMaxAge,
		OptionsPassthrough: false,
		Debug:              false,
	})
}

// NewRateLimiter returns rate limiting middleware
func NewRateLimiter(cfg *Config) func(http.Handler) http.Handler {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// 0 means disabled
	if cfg.RateLimitRequests == 0 {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return httprate.LimitByIP(cfg.RateLimitRequests, cfg.RateLimitWindow)
}

// NewJWT returns JWT authentication middleware
func NewJWT(cfg *Config) func(http.Handler) http.Handler {
	if cfg == nil || cfg.JWTSecret == "" {
		// Return passthrough if no secret configured
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	// Create JWT authenticator
	ja := jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	// Return verifier + decoder middleware
	return func(next http.Handler) http.Handler {
		return jwtauth.Verifier(ja)(jwtauth.Authenticator(next))
	}
}

// Chain chains multiple middleware handlers
func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Apply in reverse order so first middleware is outermost
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// New creates a complete middleware stack
func New(cfg *Config) []func(http.Handler) http.Handler {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	middlewares := []func(http.Handler) http.Handler{}

	// CORS (should be first)
	middlewares = append(middlewares, NewCORS(cfg))

	// Rate limiting
	middlewares = append(middlewares, NewRateLimiter(cfg))

	// JWT auth (if configured)
	if cfg.JWTSecret != "" {
		middlewares = append(middlewares, NewJWT(cfg))
	}

	return middlewares
}

// Helper types for common middleware

// RequestIDMiddleware adds a request ID to each request
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = generateRequestID()
		}
		
		ctx := r.Context()
		ctx = WithRequestID(ctx, reqID)
		
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Context key type
type contextKey string

const requestIDKey contextKey = "request_id"

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// GetRequestID gets request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

// Simple request ID generator (replace with uuid in production)
func generateRequestID() string {
	return time.Now().Format("20060102150405.000000")
}