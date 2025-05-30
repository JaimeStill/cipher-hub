package server

import (
	"cipher-hub/internal/config"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// CORS configuration constants
const (
	// CORS error handling
	CORSErrorPrefix = "CORS"

	// CORS log field names for consistency
	LogFieldCORSOrigin    = "cors_origin"
	LogFieldCORSMethod    = "cors_method"
	LogFieldCORSPreflight = "cors_preflight"
	LogFieldCORSAllowed   = "cors_allowed"

	// Default CORS configuration
	DefaultCORSMethods = "GET, POST, PUT, DELETE, OPTIONS"
	DefaultCORSHeaders = "Content-Type, Authorization, X-Request-ID"
	DefaultCORSMaxAge  = "86400" // 24 hours
)

// CORSConfig holds configuration for CORS middleware
type CORSConfig struct {
	Enabled     bool     `json:"enabled"`
	Origins     []string `json:"origins"`
	Methods     string   `json:"methods"`
	Headers     string   `json:"headers"`
	MaxAge      string   `json:"max_age"`
	Credentials bool     `json:"credentials"`
}

// LoadFromEnv populates CORS configuration from environment variables using centralized helper functions
func (c *CORSConfig) LoadFromEnv() {
	// Load origins first so we can use them for auto-enable logic
	c.Origins = config.GetEnvStringSlice(config.EnvCORSOrigins, ",", c.Origins)

	// Load other settings
	c.Methods = config.GetEnvString(config.EnvCORSMethods, c.Methods)
	c.Headers = config.GetEnvString(config.EnvCORSHeaders, c.Headers)
	c.MaxAge = config.GetEnvString(config.EnvCORSMaxAge, c.MaxAge)
	c.Credentials = config.GetEnvBool(config.EnvCORSCredentials, c.Credentials)

	// Handle Enabled with auto-enable logic:
	// If CORS_ENABLED was explicitly set in environment, use that value
	// Otherwise, auto-enable if origins are configured
	if envEnabled := os.Getenv(config.EnvCORSEnabled); envEnabled != "" {
		c.Enabled = config.GetEnvBool(config.EnvCORSEnabled, c.Enabled)
	} else {
		// Auto-enable CORS if origins are configured, disable if not
		c.Enabled = len(c.Origins) > 0
	}
}

// ApplyDefaults sets secure default values for CORS configuration
func (c *CORSConfig) ApplyDefaults() {
	// DON'T auto-enable here - origins haven't been loaded yet
	// The auto-enable logic is now in LoadFromEnv()

	// Set default methods, headers, and max age
	if c.Methods == "" {
		c.Methods = DefaultCORSMethods
	}
	if c.Headers == "" {
		c.Headers = DefaultCORSHeaders
	}
	if c.MaxAge == "" {
		c.MaxAge = DefaultCORSMaxAge
	}

	// Note: Credentials defaults to false (secure default - no credentials unless explicitly enabled)
	// Note: Enabled defaults to false and is auto-enabled in LoadFromEnv() based on Origins
}

// IsOriginAllowed checks if the given origin is in the allowed origins list.
// Performs case-insensitive matching for URL schemes and hostnames following RFC standards.
//
// Security Warning: Using "*" as an origin allows ALL origins and is a major security risk.
// Only use "*" in development environments, never in production.
func (c *CORSConfig) IsOriginAllowed(origin string) bool {
	if len(c.Origins) == 0 {
		return false // No origins configured = no CORS
	}

	normalizedOrigin := normalizeOrigin(origin)

	for _, allowedOrigin := range c.Origins {
		if allowedOrigin == "*" {
			// Log security warning for wildcard usage
			slog.Warn(
				"CORS wildcard origin matched - security risk in production",
				"wildcard_origin", "*",
				"actual_origin", origin,
				"recommendation", "use specific origins in production",
			)
			return true
		}

		if normalizeOrigin(allowedOrigin) == normalizedOrigin {
			return true
		}
	}
	return false
}

// Validate performs comprehensive validation of CORS configuration
func (c *CORSConfig) Validate() error {
	for _, origin := range c.Origins {
		if origin == "*" {
			// Allow wildcard but warn about security implications
			slog.Warn(
				"CORS wildcard origin configured - major security risk in production",
				"origin", "*",
				"recommendation", "use specific origins in production",
			)
			continue
		}

		// Validate origin URL format
		parsedURL, err := url.Parse(origin)
		if err != nil {
			return fmt.Errorf(
				"%s: invalid origin URL %q: %w",
				CORSErrorPrefix,
				origin,
				err,
			)
		}

		// Check that it has a scheme and host (required for valid origins)
		if parsedURL.Scheme == "" || parsedURL.Host == "" {
			return fmt.Errorf(
				"%s: invalid origin URL %q: missing scheme or host",
				CORSErrorPrefix,
				origin,
			)
		}
	}
	return nil
}

// normalizeOrigin converts origin to lowercase scheme and host (RFC compliant)
func normalizeOrigin(origin string) string {
	if u, err := url.Parse(origin); err == nil {
		u.Scheme = strings.ToLower(u.Scheme)
		u.Host = strings.ToLower(u.Host)
		return u.String()
	}
	// Fallback for malformed URLs
	return strings.ToLower(origin)
}

// CORSMiddleware creates middleware that handles CORS (Cross-Origin Resource Sharing) requests.
// Uses default configuration with environment variable loading.
//
// Returns:
//   - Middleware: Configured CORS middleware function
func CORSMiddleware() Middleware {
	config := CORSConfig{}
	config.ApplyDefaults()
	config.LoadFromEnv()

	// Validate configuration and log any issues
	if err := config.Validate(); err != nil {
		slog.Error("CORS configuration validation failed", "error", err)
		// Return pass-through middleware on validation failure
		return func(next http.Handler) http.Handler { return next }
	}

	return CORSMiddlewareWithConfig(config)
}

// CORSMiddlewareWithConfig creates middleware with custom CORS configuration.
// Handles CORS headers for cross-origin requests and processes preflight OPTIONS requests.
//
// Features:
//   - Environment-configurable allowed origins with secure defaults
//   - Case-insensitive origin matching following RFC standards
//   - Preflight OPTIONS request handling with proper response headers
//   - Request correlation logging for CORS events and security monitoring
//   - Conditional CORS header application based on origin validation
//   - Support for configurable HTTP methods, headers, and max-age
//   - Comprehensive security warnings for wildcard origins
//
// Security: Uses restrictive defaults (no CORS headers unless origins configured).
// Empty origins list results in no CORS headers being applied, following secure-by-default principle.
// Wildcard "*" origins trigger security warnings and should only be used in development.
//
// Parameters:
//   - config: CORSConfig with CORS behavior settings
//
// Returns:
//   - Middleware: Configured CORS middleware function
//
// Example usage:
//
//	config := CORSConfig{
//		Enabled: true,
//		Origins: []string{"https://app.example.com", "https://admin.example.com"},
//		Methods: "GET, POST, PUT, DELETE, OPTIONS",
//		Headers: "Content-Type, Authorization, X-Request-ID",
//	}
//	server.Middleware().UseIf(len(config.Origins) > 0, CORSMiddlewareWithConfig(config))
func CORSMiddlewareWithConfig(config CORSConfig) Middleware {
	return func(next http.Handler) http.Handler {
		// Skip CORS entirely if disabled or no origins configured
		if !config.Enabled || len(config.Origins) == 0 {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			requestID := GetRequestID(r.Context())

			// Check if origin is allowed (includes case-insensitive matching)
			originAllowed := config.IsOriginAllowed(origin)

			// Log CORS request for monitoring and security analysis
			if origin != "" {
				slog.Info(
					"CORS request received",
					LogFieldRequestID, requestID,
					LogFieldCORSOrigin, origin,
					LogFieldCORSMethod, r.Method,
					LogFieldCORSAllowed, originAllowed,
				)
			}

			// Set CORS headers if origin is allowed
			if originAllowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", config.Methods)
				w.Header().Set("Access-Control-Allow-Headers", config.Headers)
				w.Header().Set("Access-Control-Max-Age", config.MaxAge)

				// Set credentials header if configured
				if config.Credentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}

			// Handle preflight OPTIONS requests
			if r.Method == "OPTIONS" {
				// Log preflight request with security context
				slog.Info(
					"CORS preflight request",
					LogFieldRequestID, requestID,
					LogFieldCORSOrigin, origin,
					LogFieldCORSPreflight, true,
					LogFieldCORSAllowed, originAllowed,
				)

				if originAllowed {
					w.WriteHeader(http.StatusOK)
				} else {
					// Log security event for disallowed preflight
					slog.Warn(
						"CORS preflight request rejected",
						LogFieldRequestID, requestID,
						LogFieldCORSOrigin, origin,
						"reason", "origin not allowed",
					)
					w.WriteHeader(http.StatusForbidden)
				}
				return
			}

			// Continue to next handler for non-preflight requests
			next.ServeHTTP(w, r)
		})
	}
}
