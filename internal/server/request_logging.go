package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"cipher-hub/internal/config"
)

// Request logging constants
const (
	// Request ID generation
	RequestIDBytes     = 8                  // 8 bytes for crypto/rand generation
	RequestIDHexLength = RequestIDBytes * 2 // 16 characters hex-encoded

	// Error handling
	RequestLoggingErrorPrefix = "RequestLogging"

	// Log field names for consistency
	LogFieldRequestID     = "request_id"
	LogFieldMethod        = "method"
	LogFieldPath          = "path"
	LogFieldStatusCode    = "status_code"
	LogFieldDurationMS    = "duration_ms"
	LogFieldBytesWritten  = "bytes_written"
	LogFieldRemoteAddr    = "remote_addr"
	LogFieldUserAgent     = "user_agent"
	LogFieldContentLength = "content_length"
)

// sensitive headers that should never be logged
var SensitiveHeaders = map[string]bool{
	"authorization":       true,
	"cookie":              true,
	"x-api-key":           true,
	"x-auth-token":        true,
	"proxy-authorization": true,
}

type contextKey string

const (
	requestIDCtxKey contextKey = LogFieldRequestID
)

// RequestLoggingConfig holds configuration for request logging middleware
type RequestLoggingConfig struct {
	Enabled        bool   `json:"enabled"`
	Level          string `json:"level"`           // "debug", "info", "warn", "error"
	Format         string `json:"format"`          // "json" or "text"
	IncludeHeaders bool   `json:"include_headers"` // Include non-sensitive headers
}

// LoadFromEnv populates logging configuration from environment variables
func (c *RequestLoggingConfig) LoadFromEnv() {
	if enabled := os.Getenv(config.EnvLoggingEnabled); enabled != "" {
		c.Enabled = strings.ToLower(enabled) == "true"
	}
	if level := os.Getenv(config.EnvLogLevel); level != "" {
		c.Level = strings.ToLower(level)
	}
	if format := os.Getenv(config.EnvLogFormat); format != "" {
		c.Format = strings.ToLower(format)
	}
	if headers := os.Getenv(config.EnvLogIncludeHeaders); headers != "" {
		c.IncludeHeaders = strings.ToLower(headers) == "true"
	}
}

// ApplyDefautls sets secure default values for logging configuration
func (c *RequestLoggingConfig) ApplyDefaults() {
	if c.Level == "" {
		c.Level = "info"
	}
	if c.Format == "" {
		c.Format = "json"
	}
	// Enabled defaults to true, IncludeHeaders defaults to false
	if !c.Enabled {
		c.Enabled = true
	}
}

// generateRequestID creates a cryptographically secure request ID for correlation tracking.
// Uses crypto/rand to generate random bytes, then hex-encodes to a string.
//
// Returns:
//   - string: Hex-encoded request ID of RequestIDHexLength characters
//   - error: Generation error if crypto/rand fails
//
// Security: uses cryptographically secure random generation for correlation safety.
// The request ID is not used for authentication and provides sufficient entropy for
// request correlcation in high-throughput scenarios.
func generateRequestID() (string, error) {
	bytes := make([]byte, RequestIDBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf(
			"%s: failed to generate request ID: %w",
			RequestLoggingErrorPrefix,
			err,
		)
	}
	return hex.EncodeToString(bytes), nil
}

// WithRequestID adds a request ID to the context using a typed context key.
// This enables safe request ID propagation throughout the request lifecycle.
//
// Parameters:
//   - ctx: Parent context
//   - requestID: Request correlation ID
//
// Returns:
//   - context.Context: Context with embedded request ID
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDCtxKey, requestID)
}

// GetRequestID retrieves the request ID from context with type safety.
// Returns empty string if no request ID is found, allowing graceful handling.
//
// Parameters:
//   - ctx: Context containing request ID
//
// Returns:
//   - string: Request ID if present, empty string otherwise
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDCtxKey).(string); ok {
		return requestID
	}
	return ""
}

// responseWriter wraps http.ResponseWriter to capture status codes and response metrics.
// This enables comprehensive request logging including response status and byte counts.
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
}

// newResponseWriter creates a wrapped ResponseWriter with default status code.
// Default status is 200 OK, which matches http.ResponseWriter behavior when
// WriteHeader is not explicitly called.
//
// Parameters:
//   - w: Original http.ResponseWriter to wrap
//
// Returns:
//   - *responseWriter: Wrapped writer with status and byte tracking
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status code
		bytesWritten:   0,
	}
}

// WriteHeader captures the status code before delegating to the wrapped writer.
// This method is called automatically by the HTTP server or can be called explicitly.
//
// Parameters:
//   - code: HTTP status code to set
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the number of bytes written while delegating to the wrapped writer.
// This enables tracking of response size for performance monitoring.
//
// Parameters:
//   - b: Bytes to write to resopnse
//
// Returns:
//   - int: Number of bytes written
//   - error: Write error if any
func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += int64(n)
	return n, err
}

// StatusCode returns the captured HTTP status code.
// Returns the status code set by WriteHeader or 200 OK by default.
//
// Returns:
//   - int: HTTP status code
func (rw *responseWriter) StatusCode() int {
	return rw.statusCode
}

// BytesWritten returns the total number of bytes written to the response.
// This provides response size metrics for performance monitoring.
//
// Returns:
//   - int64: Total bytes written
func (rw *responseWriter) BytesWritten() int64 {
	return rw.bytesWritten
}

// RequestLoggingMiddleware creates middleware that logs all HTTP requests with structured logging.
// Uses default configuration with info-level JSON logging enabled.
//
// Returns:
//   - Middleware: Configured request logging middleware function
func RequestLoggingMiddleware() Middleware {
	config := RequestLoggingConfig{}
	config.ApplyDefaults()
	config.LoadFromEnv()
	return RequestLoggingMiddlewareWithConfig(config)
}

// RequestLoggingMiddlewareWithConfig creates middleware with custom logging configuration.
// Generates secure request IDs for correlation, tracks request duration and response metrics,
// and uses structured JSON logging for production environments.
//
// Features:
//   - Cryptographically secure request ID generation for correlation tracking
//   - Structured logging with log/slog using JSON format
//   - Request duration timing for performance monitoring
//   - Status code and response size tracking
//   - Request ID propagation through context
//   - Security-conscious logging (no sensitive data)
//   - Configurable logging levels and format
//
// The middleware logs two events per request:
//  1. Request start: Method, path, remote address, user agent, request ID
//  2. Request completion: Duration,s tatus code, bytes written, request ID
//
// Parameters:
//   - config: RequestLoggingConfig with logging behavior settings
//
// Returns:
//   - Middleware: Configured request logging middleware function
//
// Security: Never logs sensitive data such as authentication tokens, request bodies,
// or any user-provided data that could contain secrets. Request IDs are correlation
// tokens only and are safe for logging.
//
// Performance: minimal overhead with efficient request ID generation and structured
// logging. Response writer wrapping has negligible performance impact.
//
// Example usage:
//
//	config := RequestLoggingConfig{
//		Enabled: true,
//		Level: "info",
//		Format: "json",
//	}
//	server.Middleware().Use(RequestLoggingMiddlewareWithConfig(config))
func RequestLoggingMiddlewareWithConfig(config RequestLoggingConfig) Middleware {
	return func(next http.Handler) http.Handler {
		if !config.Enabled {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Generate secure request ID
			requestID, err := generateRequestID()
			if err != nil {
				slog.Error("Failed to generate request ID", "error", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Add request ID to context for propagation
			ctx := WithRequestID(r.Context(), requestID)
			r = r.WithContext(ctx)

			// Add request ID to response headers for client correlation
			w.Header().Set("X-Request-ID", requestID)

			// Check log level before expensive operations
			if slog.Default().Enabled(context.Background(), slog.LevelInfo) {
				// Build log fields using consistent field names
				logFields := []any{
					LogFieldRequestID, requestID,
					LogFieldMethod, r.Method,
					LogFieldPath, r.URL.Path,
					LogFieldRemoteAddr, r.RemoteAddr,
					LogFieldUserAgent, r.UserAgent(),
					LogFieldContentLength, r.ContentLength,
				}

				if config.IncludeHeaders {
					headers := filterSensitiveHeaders(r.Header)
					if len(headers) > 0 {
						logFields = append(logFields, "headers", headers)
					}
				}

				// Log request start with structured data
				slog.Info("Request started", logFields...)
			}

			// Wrap ResponseWriter to capture metrics
			wrapped := newResponseWriter(w)

			// call next handler in middleware chain
			next.ServeHTTP(wrapped, r)

			// Calculate request duration
			duration := time.Since(start)

			// Check log level before completion logging
			if slog.Default().Enabled(context.Background(), slog.LevelInfo) {
				// Log request completion with performance metrics
				slog.Info(
					"Request completed",
					LogFieldRequestID, requestID,
					LogFieldMethod, r.Method,
					LogFieldPath, r.URL.Path,
					LogFieldStatusCode, wrapped.StatusCode(),
					LogFieldDurationMS, duration.Milliseconds(),
					LogFieldBytesWritten, wrapped.BytesWritten(),
					LogFieldRemoteAddr, r.RemoteAddr,
				)
			}
		})
	}
}

// filterSensitiveHeaders removes sensitive headers from logging
func filterSensitiveHeaders(headers http.Header) map[string]string {
	filtered := make(map[string]string)
	for key, values := range headers {
		lowerKey := strings.ToLower(key)
		if !SensitiveHeaders[lowerKey] {
			filtered[key] = strings.Join(values, ", ")
		}
	}
	return filtered
}
