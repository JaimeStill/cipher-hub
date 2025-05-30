# Step 2.1.2.3: Implement CORS Handling Middleware

## Overview

**Phase**: 2 (HTTP Server Infrastructure)  
**Target**: 2.1 (Basic Server Setup)  
**Task**: 2.1.2 (Middleware Infrastructure)  
**Step**: 2.1.2.3 (Add CORS handling middleware)

**Time Estimate**: 1 hour  
**Scope**: Implement environment-configurable CORS handling middleware with preflight support

## Step Objectives

### Primary Deliverables
- **CORS Configuration Structure**: Environment-driven CORS configuration with secure defaults
- **CORS Middleware Implementation**: Handle CORS headers and preflight OPTIONS requests
- **Environment Integration**: Add CORS variables to centralized configuration management
- **Request Correlation**: Integrate CORS events with existing request logging and correlation IDs
- **Server Integration**: Use conditional middleware application with `UseIf()` pattern
- **Comprehensive Testing**: Unit and integration tests following established patterns

### Implementation Requirements
- **Files Created**: CORS configuration in `internal/config/env.go`, CORS middleware implementation
- **Files Modified**: `internal/server/server.go` (add CORS integration example), tests for CORS middleware
- **Architecture Focus**: Environment-configurable CORS with secure defaults and preflight handling
- **Security Focus**: Restrictive CORS policy by default, explicit origin configuration required
- **Go Best Practices**: Follow established middleware patterns and configuration loading
- **Foundation Usage**: Leverage request logging correlation IDs and middleware infrastructure

---

## Implementation Requirements

### Technical Specifications

#### CORS Configuration Requirements
- **Environment Variables**: Use centralized configuration with `CIPHER_HUB_CORS_*` pattern
- **Secure Defaults**: Empty origins list means no CORS headers (secure by default)
- **Conditional Application**: Use `UseIf()` pattern to apply CORS only when origins are configured
- **Request Correlation**: Log CORS events with existing request ID correlation system

#### CORS Headers Requirements
- **Basic CORS Headers**: `Access-Control-Allow-Origin`, `Access-Control-Allow-Methods`, `Access-Control-Allow-Headers`
- **Preflight Support**: Handle OPTIONS requests with appropriate response headers
- **Configurable Origins**: Support comma-separated list of allowed origins from environment
- **Standard Methods**: Support `GET, POST, PUT, DELETE, OPTIONS` by default

#### Logging Requirements
- **Request Correlation**: Use existing `GetRequestID(r.Context())` for CORS event correlation
- **Structured Logging**: Use `slog.Info()` with consistent field naming patterns
- **CORS Events**: Log preflight requests and origin validation with security context

---

## Completion Criteria

### **Step 2.1.2.3 is complete when:**

1. **Enhanced CORS Environment Variables**:
   - `EnvCORSEnabled`, `EnvCORSOrigins`, `EnvCORSMethods`, `EnvCORSHeaders`, `EnvCORSMaxAge` constants added to `internal/config/env.go`
   - Environment variables follow established naming convention: `CIPHER_HUB_CORS_*`
   - Integration with existing centralized configuration management using helper functions

2. **Advanced CORS Configuration Structure**:
   - `CORSConfig` struct with `Enabled`, `Origins`, `Methods`, `Headers`, `MaxAge`, `Credentials` fields
   - `LoadFromEnv()` method using established configuration helper functions
   - `ApplyDefaults()` method with secure default values
   - `Validate()` method with comprehensive URL validation and security warnings
   - `IsOriginAllowed()` method with case-insensitive origin matching

3. **Security-Enhanced CORS Middleware Implementation**:
   - `CORSMiddleware()` function with configuration validation and error handling
   - `CORSMiddlewareWithConfig()` function for custom configuration
   - Preflight OPTIONS request handling with appropriate status codes
   - CORS header application based on validated origin matching
   - Security warnings for wildcard origin usage
   - Enhanced logging for security monitoring and audit trails

4. **Case-Insensitive Origin Matching**:
   - `normalizeOrigin()` function for RFC-compliant URL normalization
   - Lowercase scheme and host matching while preserving path case
   - Comprehensive test coverage for various case scenarios
   - Security-conscious exact matching preventing subdomain attacks

5. **Enhanced Request Correlation Integration**:
   - CORS events logged with request correlation IDs using `GetRequestID()`
   - Structured logging with CORS-specific field names and security context
   - Preflight request logging with detailed security information
   - Security event logging for rejected requests with correlation
   - Integration with existing request logging middleware

6. **Comprehensive Security Features**:
   - Secure defaults (no CORS headers unless origins explicitly configured)
   - Configuration validation preventing malformed origin URLs
   - Security warnings for wildcard origins with production recommendations
   - Preflight request validation with forbidden responses for disallowed origins
   - Enhanced logging for security monitoring and threat detection

7. **Advanced Server Integration**:
   - Conditional middleware application using `UseIf()` pattern
   - Integration with existing middleware stack and request logging
   - Configuration validation during middleware setup
   - Method chaining support for fluent configuration
   - Thread safety and performance standards maintained

8. **Enhanced Testing Coverage**:
   - Unit tests for enhanced CORS configuration loading and validation
   - Case-insensitive origin matching tests with comprehensive scenarios
   - Configuration validation tests including malformed URL handling
   - Middleware tests for various request scenarios and security edge cases
   - Integration tests with server lifecycle and other middleware
   - Security tests for origin validation, preflight handling, and wildcard warnings
   - Request correlation tests ensuring proper ID propagation

9. **Production-Ready Documentation and Code Quality**:
   - Complete Go doc comments for all public CORS functions and types
   - Security warnings documented in code comments and validation
   - Usage examples for both simple and advanced CORS configuration
   - Environment variable examples with production security guidance
   - Code passes formatting (`go fmt`) and static analysis (`go vet`)
   - High test coverage maintained (>95%) with enhanced security testing

### **Files Created/Modified**
- `internal/config/env.go` - Added CORS environment variable constants
- `internal/server/cors.go` - Complete CORS middleware implementation
- `internal/server/cors_test.go` - Comprehensive CORS testing
- `internal/server/server_test.go` - Added CORS integration tests

---

## Testing Requirements

### Unit Testing Requirements

#### CORS Configuration Testing
- Test `LoadFromEnv()` with various environment variable combinations
- Test `ApplyDefaults()` with secure default behavior validation
- Test `IsOriginAllowed()` with exact matching and security edge cases
- Test environment variable parsing with malformed input handling

#### CORS Middleware Testing
- Test middleware with enabled/disabled configuration
- Test origin validation with allowed and disallowed origins
- Test preflight OPTIONS request handling with proper status codes
- Test CORS header application for various request types
- Test wildcard origin handling (if supported)

#### Edge Case Testing
- Test requests without Origin header
- Test requests with empty or malformed Origin headers
- Test case sensitivity in origin matching
- Test middleware behavior with empty origins configuration

### Integration Testing Requirements

#### Server Integration Testing
- Test CORS middleware integration with server lifecycle
- Test conditional middleware application using `UseIf()` pattern
- Test middleware chain execution with request logging + CORS
- Test CORS with various handler types and response scenarios

#### Request Correlation Testing
- Test CORS events are logged with request correlation IDs
- Test request ID propagation through CORS middleware
- Test structured logging output includes CORS-specific fields
- Test preflight request logging with correlation

### Security Testing Requirements
- Test CORS policy enforcement prevents unauthorized cross-origin requests
- Test preflight request handling rejects disallowed origins
- Test exact origin matching prevents subdomain attacks
- Test secure defaults prevent accidental permissive configuration

---

## Security Considerations

### CORS Security Patterns
```go
// Correct: Restrictive CORS policy by default
config := CORSConfig{
    Origins: []string{"https://app.example.com"}, // Specific origins only
}

// Incorrect: Overly permissive CORS policy
config := CORSConfig{
    Origins: []string{"*"}, // Allows all origins - security risk
}
```

### Origin Validation Security
```go
// Correct: Exact origin matching
func (c *CORSConfig) IsOriginAllowed(origin string) bool {
    for _, allowedOrigin := range c.Origins {
        if allowedOrigin == origin { // Exact match required
            return true
        }
    }
    return false
}

// Incorrect: Substring matching (vulnerable to subdomain attacks)
if strings.Contains(allowedOrigin, origin) { // NEVER do this
    return true
}
```

### Preflight Request Handling Security
```go
// Correct: Secure preflight handling with origin validation
if r.Method == "OPTIONS" {
    if originAllowed {
        w.WriteHeader(http.StatusOK)
    } else {
        w.WriteHeader(http.StatusForbidden) // Reject disallowed origins
    }
    return
}

// Incorrect: Always allowing preflight requests
if r.Method == "OPTIONS" {
    w.WriteHeader(http.StatusOK) // Security risk - no origin validation
    return
}
```

### Environment Configuration Security
```go
// Correct: Secure defaults with explicit configuration required
func (c *CORSConfig) ApplyDefaults() {
    // Secure default: disabled unless explicitly configured
    if !c.Enabled && len(c.Origins) == 0 {
        c.Enabled = false // No CORS headers by default
    }
}

// Incorrect: Permissive defaults
func (c *CORSConfig) ApplyDefaults() {
    c.Enabled = true                    // Dangerous default
    c.Origins = []string{"*"}           // Allows all origins
}
```

---

## Implementation

### Step 1: Add CORS Environment Variables to Configuration

**File**: `internal/config/env.go` (modify existing file)

Add CORS environment variables to the existing constants in the security configuration section:

```go
// Update the existing security configuration section
const (
    // ... existing server and logging constants ...

    // Security configuration (existing section - add new CORS variables here)
    EnvCORSEnabled = "CIPHER_HUB_CORS_ENABLED"    // Add this line
    EnvCORSOrigins = "CIPHER_HUB_CORS_ORIGINS"    // Already exists - don't duplicate
    EnvCORSMethods = "CIPHER_HUB_CORS_METHODS"    // Add this line
    EnvCORSHeaders = "CIPHER_HUB_CORS_HEADERS"    // Add this line  
    EnvCORSMaxAge  = "CIPHER_HUB_CORS_MAX_AGE"    // Add this line
    EnvTLSCertFile = "CIPHER_HUB_TLS_CERT_FILE"   // Existing
    EnvTLSKeyFile  = "CIPHER_HUB_TLS_KEY_FILE"    // Existing

    // ... rest of existing constants ...
)
```

### Step 2: Create CORS Configuration Structure

**File**: `internal/server/cors.go` (new file)

```go
package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"cipher-hub/internal/config"
)

// CORS configuration constants
const (
	// CORS error handling
	CORSErrorPrefix = "CORS"

	// CORS log field names for consistency
	LogFieldCORSOrigin     = "cors_origin"
	LogFieldCORSMethod     = "cors_method"
	LogFieldCORSPreflight  = "cors_preflight"
	LogFieldCORSAllowed    = "cors_allowed"

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

// LoadFromEnv populates CORS configuration from environment variables using established helper functions
func (c *CORSConfig) LoadFromEnv() {
	// Use established configuration helpers
	c.Enabled = config.GetEnvBool(config.EnvCORSEnabled, c.Enabled)
	c.Origins = config.GetEnvStringSlice(config.EnvCORSOrigins, ",", c.Origins)
	c.Methods = config.GetEnvString(config.EnvCORSMethods, c.Methods)
	c.Headers = config.GetEnvString(config.EnvCORSHeaders, c.Headers)
	c.MaxAge = config.GetEnvString(config.EnvCORSMaxAge, c.MaxAge)
}

// ApplyDefaults sets secure default values for CORS configuration
func (c *CORSConfig) ApplyDefaults() {
	// Secure default: disabled unless explicitly configured
	if !c.Enabled && len(c.Origins) == 0 {
		c.Enabled = false
	}

	// Enable CORS if origins are configured
	if len(c.Origins) > 0 {
		c.Enabled = true
	}

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
}

// Validate performs comprehensive validation of CORS configuration
func (c *CORSConfig) Validate() error {
	for _, origin := range c.Origins {
		if origin == "*" {
			// Allow wildcard but warn about security implications
			slog.Warn("CORS wildcard origin configured - major security risk in production",
				"origin", "*",
				"recommendation", "use specific origins in production")
			continue
		}
		
		// Validate origin URL format
		if _, err := url.Parse(origin); err != nil {
			return fmt.Errorf("%s: invalid origin URL %q: %w", 
				CORSErrorPrefix, origin, err)
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
			slog.Warn("CORS wildcard origin matched - security risk in production",
				"wildcard_origin", "*",
				"actual_origin", origin,
				"recommendation", "use specific origins in production")
			return true
		}
		
		if normalizeOrigin(allowedOrigin) == normalizedOrigin {
			return true
		}
	}
	return false
}
```

### Step 3: Implement CORS Middleware

**Continue in**: `internal/server/cors.go`

```go
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
				slog.Info("CORS request received",
					LogFieldRequestID, requestID,
					LogFieldCORSOrigin, origin,
					LogFieldCORSMethod, r.Method,
					LogFieldCORSAllowed, originAllowed)
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
				slog.Info("CORS preflight request",
					LogFieldRequestID, requestID,
					LogFieldCORSOrigin, origin,
					LogFieldCORSPreflight, true,
					LogFieldCORSAllowed, originAllowed)

				if originAllowed {
					w.WriteHeader(http.StatusOK)
				} else {
					// Log security event for disallowed preflight
					slog.Warn("CORS preflight request rejected",
						LogFieldRequestID, requestID,
						LogFieldCORSOrigin, origin,
						"reason", "origin not allowed")
					w.WriteHeader(http.StatusForbidden)
				}
				return
			}

			// Continue to next handler for non-preflight requests
			next.ServeHTTP(w, r)
		})
	}
}
```

### Step 4: Create Comprehensive CORS Tests

**File**: `internal/server/cors_test.go` (new file)

```go
package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCORSConfig_LoadFromEnv(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{
		"CIPHER_HUB_CORS_ENABLED",
		"CIPHER_HUB_CORS_ORIGINS",
		"CIPHER_HUB_CORS_METHODS",
		"CIPHER_HUB_CORS_HEADERS",
		"CIPHER_HUB_CORS_MAX_AGE",
	}

	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
	}

	// Clean up environment after test
	defer func() {
		for _, key := range envVars {
			if val, exists := originalEnv[key]; exists {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		expected CORSConfig
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			expected: CORSConfig{
				Enabled: false,
				Origins: nil,
				Methods: DefaultCORSMethods,
				Headers: DefaultCORSHeaders,
				MaxAge:  DefaultCORSMaxAge,
			},
		},
		{
			name: "enabled with origins and custom config",
			envVars: map[string]string{
				"CIPHER_HUB_CORS_ENABLED": "true",
				"CIPHER_HUB_CORS_ORIGINS": "https://app.example.com,https://admin.example.com",
				"CIPHER_HUB_CORS_METHODS": "GET, POST, PUT",
				"CIPHER_HUB_CORS_HEADERS": "Content-Type, Authorization",
				"CIPHER_HUB_CORS_MAX_AGE": "3600",
			},
			expected: CORSConfig{
				Enabled: true,
				Origins: []string{"https://app.example.com", "https://admin.example.com"},
				Methods: "GET, POST, PUT",
				Headers: "Content-Type, Authorization",
				MaxAge:  "3600",
			},
		},
		{
			name: "origins without explicit enabled",
			envVars: map[string]string{
				"CIPHER_HUB_CORS_ORIGINS": "https://app.example.com",
			},
			expected: CORSConfig{
				Enabled: true, // Should be enabled automatically when origins are set
				Origins: []string{"https://app.example.com"},
				Methods: DefaultCORSMethods,
				Headers: DefaultCORSHeaders,
				MaxAge:  DefaultCORSMaxAge,
			},
		},
		{
			name: "disabled explicitly with origins",
			envVars: map[string]string{
				"CIPHER_HUB_CORS_ENABLED": "false",
				"CIPHER_HUB_CORS_ORIGINS": "https://app.example.com",
			},
			expected: CORSConfig{
				Enabled: true, // Origins override disabled flag
				Origins: []string{"https://app.example.com"},
				Methods: DefaultCORSMethods,
				Headers: DefaultCORSHeaders,
				MaxAge:  DefaultCORSMaxAge,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			config := CORSConfig{}
			config.ApplyDefaults()
			config.LoadFromEnv()

			if config.Enabled != tt.expected.Enabled {
				t.Errorf("Enabled = %v, want %v", config.Enabled, tt.expected.Enabled)
			}

			if len(config.Origins) != len(tt.expected.Origins) {
				t.Errorf("Origins length = %v, want %v", len(config.Origins), len(tt.expected.Origins))
			} else {
				for i, origin := range config.Origins {
					if origin != tt.expected.Origins[i] {
						t.Errorf("Origins[%d] = %v, want %v", i, origin, tt.expected.Origins[i])
					}
				}
			}

			if config.Methods != tt.expected.Methods {
				t.Errorf("Methods = %v, want %v", config.Methods, tt.expected.Methods)
			}

			if config.Headers != tt.expected.Headers {
				t.Errorf("Headers = %v, want %v", config.Headers, tt.expected.Headers)
			}

			if config.MaxAge != tt.expected.MaxAge {
				t.Errorf("MaxAge = %v, want %v", config.MaxAge, tt.expected.MaxAge)
			}
		})
	}
}

func TestCORSConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  CORSConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid configuration",
			config: CORSConfig{
				Origins: []string{"https://app.example.com", "http://localhost:3000"},
			},
			wantErr: false,
		},
		{
			name: "wildcard origin (valid but warns)",
			config: CORSConfig{
				Origins: []string{"*"},
			},
			wantErr: false, // Valid but should warn
		},
		{
			name: "invalid origin URL",
			config: CORSConfig{
				Origins: []string{"not-a-valid-url"},
			},
			wantErr: true,
			errMsg:  "invalid origin URL",
		},
		{
			name: "mixed valid and invalid origins",
			config: CORSConfig{
				Origins: []string{"https://app.example.com", "invalid-url"},
			},
			wantErr: true,
			errMsg:  "invalid origin URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error, got nil")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestCORSConfig_IsOriginAllowed(t *testing.T) {
	tests := []struct {
		name     string
		config   CORSConfig
		origin   string
		expected bool
	}{
		{
			name: "no origins configured",
			config: CORSConfig{
				Origins: []string{},
			},
			origin:   "https://app.example.com",
			expected: false,
		},
		{
			name: "origin allowed (exact match)",
			config: CORSConfig{
				Origins: []string{"https://app.example.com", "http://localhost:3000"},
			},
			origin:   "https://app.example.com",
			expected: true,
		},
		{
			name: "origin not allowed",
			config: CORSConfig{
				Origins: []string{"https://app.example.com"},
			},
			origin:   "http://evil.com",
			expected: false,
		},
		{
			name: "wildcard origin",
			config: CORSConfig{
				Origins: []string{"*"},
			},
			origin:   "https://anywhere.com",
			expected: true,
		},
		{
			name: "case insensitive scheme matching",
			config: CORSConfig{
				Origins: []string{"https://app.example.com"},
			},
			origin:   "HTTPS://app.example.com",
			expected: true,
		},
		{
			name: "case insensitive host matching",
			config: CORSConfig{
				Origins: []string{"https://app.example.com"},
			},
			origin:   "https://APP.EXAMPLE.COM",
			expected: true,
		},
		{
			name: "case insensitive full URL matching",
			config: CORSConfig{
				Origins: []string{"http://localhost:3000"},
			},
			origin:   "HTTP://LOCALHOST:3000",
			expected: true,
		},
		{
			name: "different scheme should not match",
			config: CORSConfig{
				Origins: []string{"https://app.example.com"},
			},
			origin:   "http://app.example.com",
			expected: false,
		},
		{
			name: "different port should not match",
			config: CORSConfig{
				Origins: []string{"http://localhost:3000"},
			},
			origin:   "http://localhost:8080",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsOriginAllowed(tt.origin)
			if result != tt.expected {
				t.Errorf("IsOriginAllowed(%q) = %v, want %v", tt.origin, result, tt.expected)
			}
		})
	}
}

func TestNormalizeOrigin(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		expected string
	}{
		{
			name:     "lowercase scheme and host",
			origin:   "HTTPS://APP.EXAMPLE.COM",
			expected: "https://app.example.com",
		},
		{
			name:     "preserve port",
			origin:   "HTTP://LOCALHOST:3000",
			expected: "http://localhost:3000",
		},
		{
			name:     "preserve path case",
			origin:   "https://API.EXAMPLE.COM/API/v1",
			expected: "https://api.example.com/API/v1",
		},
		{
			name:     "already lowercase",
			origin:   "https://app.example.com",
			expected: "https://app.example.com",
		},
		{
			name:     "malformed URL fallback",
			origin:   "not-a-valid-url",
			expected: "not-a-valid-url", // Fallback to lowercase
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeOrigin(tt.origin)
			if result != tt.expected {
				t.Errorf("normalizeOrigin(%q) = %q, want %q", tt.origin, result, tt.expected)
			}
		})
	}
}

func TestCORSMiddleware(t *testing.T) {
	middleware := CORSMiddleware()

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request ID is available in context
		requestID := GetRequestID(r.Context())
		if requestID == "" {
			t.Error("Request ID not found in context")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Apply middleware
	wrappedHandler := middleware(handler)

	// Create test request with Origin header
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	w := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(w, req)

	// Verify response (no CORS headers since no origins configured by default)
	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("CORS headers should not be set without configured origins")
	}

	// Verify response body
	if w.Code != http.StatusOK {
		t.Errorf("Response status = %d, want %d", w.Code, http.StatusOK)
	}

	if w.Body.String() != "test response" {
		t.Errorf("Response body = %q, want %q", w.Body.String(), "test response")
	}
}

func TestCORSMiddlewareWithConfig(t *testing.T) {
	tests := []struct {
		name           string
		config         CORSConfig
		requestOrigin  string
		method         string
		expectCORS     bool
		expectedOrigin string
		expectedStatus int
	}{
		{
			name: "CORS disabled",
			config: CORSConfig{
				Enabled: false,
				Origins: []string{},
			},
			requestOrigin:  "http://localhost:3000",
			method:         "GET",
			expectCORS:     false,
			expectedOrigin: "",
			expectedStatus: http.StatusOK,
		},
		{
			name: "allowed origin GET request",
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
				Methods: DefaultCORSMethods,
				Headers: DefaultCORSHeaders,
			},
			requestOrigin:  "http://localhost:3000",
			method:         "GET",
			expectCORS:     true,
			expectedOrigin: "http://localhost:3000",
			expectedStatus: http.StatusOK,
		},
		{
			name: "disallowed origin GET request",
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
				Methods: DefaultCORSMethods,
				Headers: DefaultCORSHeaders,
			},
			requestOrigin:  "http://evil.com",
			method:         "GET",
			expectCORS:     false,
			expectedOrigin: "",
			expectedStatus: http.StatusOK,
		},
		{
			name: "allowed origin OPTIONS preflight",
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
				Methods: DefaultCORSMethods,
				Headers: DefaultCORSHeaders,
			},
			requestOrigin:  "http://localhost:3000",
			method:         "OPTIONS",
			expectCORS:     true,
			expectedOrigin: "http://localhost:3000",
			expectedStatus: http.StatusOK,
		},
		{
			name: "disallowed origin OPTIONS preflight",
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
				Methods: DefaultCORSMethods,
				Headers: DefaultCORSHeaders,
			},
			requestOrigin:  "http://evil.com",
			method:         "OPTIONS",
			expectCORS:     false,
			expectedOrigin: "",
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "wildcard origin",
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"*"},
				Methods: DefaultCORSMethods,
				Headers: DefaultCORSHeaders,
			},
			requestOrigin:  "http://anywhere.com",
			method:         "GET",
			expectCORS:     true,
			expectedOrigin: "http://anywhere.com",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := CORSMiddlewareWithConfig(tt.config)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestID := GetRequestID(r.Context())
				if requestID == "" {
					t.Error("Request ID not available in handler context")
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("handler response"))
			})

			wrappedHandler := middleware(handler)

			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.requestOrigin != "" {
				req.Header.Set("Origin", tt.requestOrigin)
			}

			w := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w, req)

			// Verify CORS headers
			corsOrigin := w.Header().Get("Access-Control-Allow-Origin")
			if tt.expectCORS {
				if corsOrigin != tt.expectedOrigin {
					t.Errorf("Access-Control-Allow-Origin = %q, want %q", corsOrigin, tt.expectedOrigin)
				}

				if methods := w.Header().Get("Access-Control-Allow-Methods"); methods == "" {
					t.Error("Access-Control-Allow-Methods should be set for allowed origins")
				}

				if headers := w.Header().Get("Access-Control-Allow-Headers"); headers == "" {
					t.Error("Access-Control-Allow-Headers should be set for allowed origins")
				}
			} else {
				if corsOrigin != "" {
					t.Errorf("Access-Control-Allow-Origin should not be set, got %q", corsOrigin)
				}
			}

			// Verify response status
			if w.Code != tt.expectedStatus {
				t.Errorf("Response status = %d, want %d", w.Code, tt.expectedStatus)
			}

			// Verify handler was called for non-OPTIONS or allowed OPTIONS requests
			if tt.method != "OPTIONS" || (tt.method == "OPTIONS" && tt.expectCORS) {
				if tt.method != "OPTIONS" && w.Body.String() != "handler response" {
					t.Errorf("Handler response = %q, want %q", w.Body.String(), "handler response")
				}
			}
		})
	}
}

func TestCORSMiddleware_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		setupReq  func() *http.Request
		config    CORSConfig
		expectErr bool
	}{
		{
			name: "no origin header",
			setupReq: func() *http.Request {
				return httptest.NewRequest("GET", "/test", nil)
			},
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
			},
			expectErr: false,
		},
		{
			name: "empty origin header",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Origin", "")
				return req
			},
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
			},
			expectErr: false,
		},
		{
			name: "malformed origin header",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Origin", "not-a-valid-url")
				return req
			},
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
			},
			expectErr: false,
		},
		{
			name: "case sensitive origin matching",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Origin", "HTTP://LOCALHOST:3000")
				return req
			},
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := CORSMiddlewareWithConfig(tt.config)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestID := GetRequestID(r.Context())
				if requestID == "" {
					t.Error("Request ID should be available")
				}
				w.WriteHeader(http.StatusOK)
			})

			wrappedHandler := middleware(handler)
			req := tt.setupReq()
			w := httptest.NewRecorder()

			// Should not panic or fail regardless of request content
			wrappedHandler.ServeHTTP(w, req)

			if tt.expectErr {
				if w.Code < 400 {
					t.Errorf("Expected error response, got status %d", w.Code)
				}
			} else {
				if w.Code != http.StatusOK {
					t.Errorf("Expected success response, got status %d", w.Code)
				}
			}
		})
	}
}

func TestCORSMiddleware_RequestCorrelation(t *testing.T) {
	config := CORSConfig{
		Enabled: true,
		Origins: []string{"http://localhost:3000"},
		Methods: DefaultCORSMethods,
		Headers: DefaultCORSHeaders,
	}

	middleware := CORSMiddlewareWithConfig(config)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request ID is available for correlation
		requestID := GetRequestID(r.Context())
		if requestID == "" {
			t.Error("Request ID should be available for CORS correlation")
		}

		// Verify request ID format
		if len(requestID) != RequestIDHexLength {
			t.Errorf("Request ID length = %d, want %d", len(requestID), RequestIDHexLength)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("correlation test"))
	})

	// We need to wrap with request logging middleware to generate request ID
	wrappedHandler := RequestLoggingMiddleware()(middleware(handler))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)

	// Verify request ID was added by request logging middleware
	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("X-Request-ID header should be present from request logging middleware")
	}

	// Verify CORS headers were applied
	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Error("CORS headers should be applied for allowed origin")
	}

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Response status = %d, want %d", w.Code, http.StatusOK)
	}

	if w.Body.String() != "correlation test" {
		t.Errorf("Response body = %q, want %q", w.Body.String(), "correlation test")
	}
}
```

### Step 5: Add Server Integration Tests

**File**: `internal/server/server_test.go` (add these tests to existing file)

```go
func TestServer_CORSMiddleware(t *testing.T) {
	config := ServerConfig{
		Host: "localhost",
		Port: "0",
	}

	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("NewServer() unexpected error: %v", err)
	}

	// Configure CORS middleware with specific origins
	corsConfig := CORSConfig{
		Enabled: true,
		Origins: []string{"http://localhost:3000", "https://app.example.com"},
		Methods: DefaultCORSMethods,
		Headers: DefaultCORSHeaders,
	}

	// Add middleware with conditional application
	server.Middleware().
		Use(RequestLoggingMiddleware()).
		UseIf(len(corsConfig.Origins) > 0, CORSMiddlewareWithConfig(corsConfig))

	// Set test handler
	server.SetHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request ID is available in handler
		requestID := GetRequestID(r.Context())
		if requestID == "" {
			t.Error("Request ID not available in handler context")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("CORS test response"))
	}))

	// Start server
	err = server.Start()
	if err != nil {
		t.Fatalf("Start() unexpected error: %v", err)
	}
	defer func() {
		if err := server.Shutdown(); err != nil {
			t.Logf("Cleanup shutdown error: %v", err)
		}
	}()

	// Test CORS request integration
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	w := httptest.NewRecorder()

	// Execute request through server
	server.httpServer.Handler.ServeHTTP(w, req)

	// Verify request ID header was added by request logging middleware
	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("X-Request-ID header not set by request logging middleware")
	}

	// Verify CORS headers were added by CORS middleware
	corsOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if corsOrigin != "http://localhost:3000" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q", corsOrigin, "http://localhost:3000")
	}

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Response status = %d, want %d", w.Code, http.StatusOK)
	}

	if w.Body.String() != "CORS test response" {
		t.Errorf("Response body = %q, want %q", w.Body.String(), "CORS test response")
	}
}

func TestServer_CORSMiddleware_Conditional(t *testing.T) {
	config := ServerConfig{
		Host: "localhost",
		Port: "0",
	}

	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("NewServer() unexpected error: %v", err)
	}

	// Configure CORS middleware without origins (should not be applied)
	corsConfig := CORSConfig{
		Enabled: false,
		Origins: []string{}, // Empty origins
	}

	// Add middleware with conditional application
	server.Middleware().
		Use(RequestLoggingMiddleware()).
		UseIf(len(corsConfig.Origins) > 0, CORSMiddlewareWithConfig(corsConfig)) // Should not apply

	// Set test handler
	server.SetHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := GetRequestID(r.Context())
		if requestID == "" {
			t.Error("Request ID not available in handler context")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("no CORS test"))
	}))

	// Start server
	err = server.Start()
	if err != nil {
		t.Fatalf("Start() unexpected error: %v", err)
	}
	defer func() {
		if err := server.Shutdown(); err != nil {
			t.Logf("Cleanup shutdown error: %v", err)
		}
	}()

	// Test request without CORS
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	w := httptest.NewRecorder()

	server.httpServer.Handler.ServeHTTP(w, req)

	// Verify request ID header was added by request logging middleware
	if w.Header().Get("X-Request-ID") == "" {
		t.Error("X-Request-ID header should be present from request logging")
	}

	// Verify NO CORS headers were added (middleware not applied)
	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("CORS headers should not be present when middleware not applied")
	}

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Response status = %d, want %d", w.Code, http.StatusOK)
	}

	if w.Body.String() != "no CORS test" {
		t.Errorf("Response body = %q, want %q", w.Body.String(), "no CORS test")
	}
}

func TestServer_CORSMiddleware_Preflight(t *testing.T) {
	config := ServerConfig{
		Host: "localhost",
		Port: "0",
	}

	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("NewServer() unexpected error: %v", err)
	}

	// Configure CORS middleware with origins
	corsConfig := CORSConfig{
		Enabled: true,
		Origins: []string{"http://localhost:3000"},
		Methods: DefaultCORSMethods,
		Headers: DefaultCORSHeaders,
	}

	// Add middleware chain
	server.Middleware().
		Use(RequestLoggingMiddleware()).
		UseIf(len(corsConfig.Origins) > 0, CORSMiddlewareWithConfig(corsConfig))

	// Set test handler (should not be called for OPTIONS preflight)
	server.SetHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for preflight OPTIONS request")
		w.WriteHeader(http.StatusOK)
	}))

	// Start server
	err = server.Start()
	if err != nil {
		t.Fatalf("Start() unexpected error: %v", err)
	}
	defer func() {
		if err := server.Shutdown(); err != nil {
			t.Logf("Cleanup shutdown error: %v", err)
		}
	}()

	// Test preflight OPTIONS request
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	w := httptest.NewRecorder()

	server.httpServer.Handler.ServeHTTP(w, req)

	// Verify request ID header was added
	if w.Header().Get("X-Request-ID") == "" {
		t.Error("X-Request-ID header should be present")
	}

	// Verify CORS preflight headers
	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Error("CORS origin header should be set for preflight")
	}

	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("CORS methods header should be set for preflight")
	}

	if w.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Error("CORS headers header should be set for preflight")
	}

	// Verify preflight response status
	if w.Code != http.StatusOK {
		t.Errorf("Preflight response status = %d, want %d", w.Code, http.StatusOK)
	}
}
```

### Step 6: Create Usage Example

**File**: `examples/cors_usage.go` (optional example file)

```go
package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"cipher-hub/internal/server"
)

func main() {
	// Configure structured logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	config := server.ServerConfig{
		Host: "localhost",
		Port: "8080",
	}

	srv, err := server.NewServer(config)
	if err != nil {
		panic(err)
	}

	// Configure CORS with specific origins
	corsConfig := server.CORSConfig{
		Enabled: true,
		Origins: []string{"http://localhost:3000", "https://app.example.com"},
	}

	// Configure middleware with conditional CORS
	srv.Middleware().
		Use(server.RequestLoggingMiddleware()).                                         // Always log requests
		UseIf(len(corsConfig.Origins) > 0, server.CORSMiddlewareWithConfig(corsConfig)) // CORS when origins configured

	// Set handler
	srv.SetHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := server.GetRequestID(r.Context())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hello! Request ID: %s", requestID)))
	}))

	fmt.Println("CORS middleware integration example compiled successfully")
	fmt.Printf("Middleware count: %d\n", srv.Middleware().Count())
	fmt.Printf("CORS origins configured: %v\n", corsConfig.Origins)
}
```

---

## Verification Steps

### Step 1: Build Verification
```bash
# Navigate to project root
cd cipher-hub/

# Verify clean build with CORS middleware
go build ./...

# Expected: No compilation errors
```

### Step 2: Unit Test Verification
```bash
# Run CORS-specific tests
go test ./internal/server -run "TestCORS" -v

# Expected: All CORS tests pass
# Sample output:
# === RUN   TestCORSConfig_LoadFromEnv
# === RUN   TestCORSConfig_IsOriginAllowed
# === RUN   TestCORSMiddleware
# === RUN   TestCORSMiddlewareWithConfig
# --- PASS: All CORS tests should pass
```

### Step 3: Server Integration Test Verification
```bash
# Run server integration tests with CORS
go test ./internal/server -run "TestServer_CORS" -v

# Expected: All server integration tests pass
# === RUN   TestServer_CORSMiddleware
# === RUN   TestServer_CORSMiddleware_Conditional
# === RUN   TestServer_CORSMiddleware_Preflight
# --- PASS: All integration tests should pass
```

### Step 4: Complete Test Suite Verification
```bash
# Run all server tests to ensure no regressions
go test ./internal/server -v

# Expected: All existing and new tests pass
# Verify CORS middleware doesn't break existing functionality
```

### Step 5: Environment Variable Configuration Test
```bash
# Test enhanced CORS configuration loading
cat > test_cors_config.go << 'EOF'
package main

import (
    "fmt"
    "os"
    "cipher-hub/internal/server"
)

func main() {
    // Set comprehensive test environment variables
    os.Setenv("CIPHER_HUB_CORS_ENABLED", "true")
    os.Setenv("CIPHER_HUB_CORS_ORIGINS", "https://app.example.com,https://admin.example.com")
    os.Setenv("CIPHER_HUB_CORS_METHODS", "GET, POST, PUT, DELETE")
    os.Setenv("CIPHER_HUB_CORS_HEADERS", "Content-Type, Authorization, X-Custom-Header")
    os.Setenv("CIPHER_HUB_CORS_MAX_AGE", "7200")
    
    config := server.CORSConfig{}
    config.ApplyDefaults()
    config.LoadFromEnv()
    
    // Validate configuration
    if err := config.Validate(); err != nil {
        fmt.Printf("Configuration validation failed: %v\n", err)
        return
    }
    
    fmt.Printf("CORS Enabled: %v\n", config.Enabled)
    fmt.Printf("CORS Origins: %v\n", config.Origins)
    fmt.Printf("CORS Methods: %s\n", config.Methods)
    fmt.Printf("CORS Headers: %s\n", config.Headers)
    fmt.Printf("CORS Max Age: %s\n", config.MaxAge)
    
    // Test case-insensitive origin validation
    fmt.Printf("app.example.com (https) allowed: %v\n", config.IsOriginAllowed("https://app.example.com"))
    fmt.Printf("APP.EXAMPLE.COM (HTTPS) allowed: %v\n", config.IsOriginAllowed("HTTPS://APP.EXAMPLE.COM"))
    fmt.Printf("evil.com allowed: %v\n", config.IsOriginAllowed("https://evil.com"))
    
    // Test wildcard warning
    fmt.Println("\nTesting wildcard configuration...")
    wildcardConfig := server.CORSConfig{Origins: []string{"*"}}
    fmt.Printf("Wildcard (*) allows anything: %v\n", wildcardConfig.IsOriginAllowed("https://anywhere.com"))
}
EOF

go run test_cors_config.go
rm test_cors_config.go

# Expected output:
# CORS Enabled: true
# CORS Origins: [https://app.example.com https://admin.example.com]
# CORS Methods: GET, POST, PUT, DELETE
# CORS Headers: Content-Type, Authorization, X-Custom-Header
# CORS Max Age: 7200
# app.example.com (https) allowed: true
# APP.EXAMPLE.COM (HTTPS) allowed: true
# evil.com allowed: false
# 
# Testing wildcard configuration...
# Wildcard (*) allows anything: true
```

### Step 6: CORS Integration Test
```bash
# Create enhanced integration test program
cat > test_cors_integration.go << 'EOF'
package main

import (
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "cipher-hub/internal/server"
)

func main() {
    // Configure structured logging
    slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
    
    config := server.ServerConfig{
        Host: "localhost",
        Port: "8080",
    }
    
    srv, err := server.NewServer(config)
    if err != nil {
        panic(err)
    }
    
    // Configure enhanced CORS with specific origins and custom settings
    corsConfig := server.CORSConfig{
        Enabled: true,
        Origins: []string{"https://app.example.com", "https://admin.example.com"},
        Methods: "GET, POST, PUT, DELETE, OPTIONS",
        Headers: "Content-Type, Authorization, X-Request-ID, X-Custom-Header",
        MaxAge:  "7200",
    }
    
    // Validate CORS configuration
    if err := corsConfig.Validate(); err != nil {
        fmt.Printf("CORS configuration validation failed: %v\n", err)
        return
    }
    
    // Configure middleware with conditional CORS
    srv.Middleware().
        Use(server.RequestLoggingMiddleware()).
        UseIf(len(corsConfig.Origins) > 0, server.CORSMiddlewareWithConfig(corsConfig))
    
    // Set handler
    srv.SetHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := server.GetRequestID(r.Context())
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(fmt.Sprintf("Hello! Request ID: %s", requestID)))
    }))
    
    fmt.Println("Enhanced CORS middleware integration test compiled successfully")
    fmt.Printf("Middleware count: %d\n", srv.Middleware().Count())
    fmt.Printf("CORS enabled with origins: %v\n", corsConfig.Origins)
    fmt.Printf("CORS methods: %s\n", corsConfig.Methods)
    fmt.Printf("CORS headers: %s\n", corsConfig.Headers)
    fmt.Printf("CORS max age: %s seconds\n", corsConfig.MaxAge)
    
    // Test case-insensitive origin matching
    fmt.Printf("\nCase-insensitive origin matching:")
    fmt.Printf("\n  https://app.example.com allowed: %v", corsConfig.IsOriginAllowed("https://app.example.com"))
    fmt.Printf("\n  HTTPS://APP.EXAMPLE.COM allowed: %v", corsConfig.IsOriginAllowed("HTTPS://APP.EXAMPLE.COM"))
}
EOF

go run test_cors_integration.go
rm test_cors_integration.go

# Expected: "Enhanced CORS middleware integration test compiled successfully"
# Expected: "Middleware count: 2"
# Expected: "CORS enabled with origins: [https://app.example.com https://admin.example.com]"
# Expected: "CORS methods: GET, POST, PUT, DELETE, OPTIONS"
# Expected: "CORS headers: Content-Type, Authorization, X-Request-ID, X-Custom-Header"
# Expected: "CORS max age: 7200 seconds"
# Expected: Case-insensitive origin matching results
```

### Step 7: Code Quality Verification
```bash
# Format and lint checks
go fmt ./...
go vet ./...

# Expected: No issues reported
```

### Step 8: Documentation Verification
```bash
# Verify go doc generates proper documentation
go doc -all ./internal/server | grep -A 10 "func CORSMiddleware"
go doc -all ./internal/server | grep -A 5 "type CORSConfig"

# Expected: Complete documentation for CORS functions and types
```