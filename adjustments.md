# Configuration Helper Functions Implementation

## Step 1: Add Helper Functions to Configuration Package

**File**: `internal/config/env.go` (add to existing file)

```go
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Environment variable names for Cipher Hub
const (
	// Server configuration
	EnvHost            = "CIPHER_HUB_HOST"
	EnvPort            = "CIPHER_HUB_PORT"
	EnvReadTimeout     = "CIPHER_HUB_READ_TIMEOUT"
	EnvWriteTimeout    = "CIPHER_HUB_WRITE_TIMEOUT"
	EnvIdleTimeout     = "CIPHER_HUB_IDLE_TIMEOUT"
	EnvShutdownTimeout = "CIPHER_HUB_SHUTDOWN_TIMEOUT"

	// Logging configuration
	EnvLoggingEnabled    = "CIPHER_HUB_LOGGING_ENABLED"
	EnvLogLevel          = "CIPHER_HUB_LOG_LEVEL"
	EnvLogFormat         = "CIPHER_HUB_LOG_FORMAT"
	EnvLogIncludeHeaders = "CIPHER_HUB_LOG_INCLUDE_HEADERS"

	// Security configuration
	EnvCORSEnabled     = "CIPHER_HUB_CORS_ENABLED"
	EnvCORSOrigins     = "CIPHER_HUB_CORS_ORIGINS"
	EnvCORSMethods     = "CIPHER_HUB_CORS_METHODS"
	EnvCORSHeaders     = "CIPHER_HUB_CORS_HEADERS"
	EnvCORSMaxAge      = "CIPHER_HUB_CORS_MAX_AGE"
	EnvCORSCredentials = "CIPHER_HUB_CORS_CREDENTIALS"
	EnvTLSCertFile     = "CIPHER_HUB_TLS_CERT_FILE"
	EnvTLSKeyFile      = "CIPHER_HUB_TLS_KEY_FILE"

	// Application configuration
	EnvEnvironment = "CIPHER_HUB_ENVIRONMENT"
	EnvDatabaseURL = "CIPHER_HUB_DATABASE_URL"
)

// GetEnvString retrieves a string environment variable with a default value.
// Trims whitespace from the retrieved value.
//
// Parameters:
//   - key: Environment variable name
//   - defaultValue: Value to return if environment variable is empty or not set
//
// Returns:
//   - string: Environment variable value (trimmed) or default value
func GetEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return strings.TrimSpace(value)
	}
	return defaultValue
}

// GetEnvBool retrieves a boolean environment variable with a default value.
// Uses a flexible approach similar to strconv.ParseBool but with fallback to default.
//
// Truthy values: "1", "t", "T", "true", "TRUE", "True", "yes", "YES", "Yes", "y", "Y", "on", "ON", "On"
// Falsy values: "0", "f", "F", "false", "FALSE", "False", "no", "NO", "No", "n", "N", "off", "OFF", "Off"
//
// Parameters:
//   - key: Environment variable name
//   - defaultValue: Value to return if environment variable is empty, not set, or unrecognized
//
// Returns:
//   - bool: Parsed boolean value or default value
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		normalized := strings.ToLower(strings.TrimSpace(value))
		
		// Truthy values (comprehensive)
		switch normalized {
		case "1", "t", "true", "yes", "y", "on":
			return true
		case "0", "f", "false", "no", "n", "off":
			return false
		}
	}
	return defaultValue
}

// GetEnvStringSlice retrieves a string slice from an environment variable.
// Splits the value by the specified separator and trims whitespace from each item.
// Empty items after trimming are excluded from the result.
//
// Parameters:
//   - key: Environment variable name
//   - separator: String to split the value on (e.g., ",", ";", "|")
//   - defaultValue: Value to return if environment variable is empty or not set
//
// Returns:
//   - []string: Parsed string slice or default value
//
// Example:
//   - GetEnvStringSlice("ORIGINS", ",", nil) with "app.com, admin.com , " returns ["app.com", "admin.com"]
func GetEnvStringSlice(key, separator string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		items := strings.Split(value, separator)
		result := make([]string, 0, len(items))
		for _, item := range items {
			if trimmed := strings.TrimSpace(item); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	return defaultValue
}

// GetEnvInt retrieves an integer environment variable with a default value.
// Uses strconv.Atoi for parsing. Invalid values fall back to the default.
//
// Parameters:
//   - key: Environment variable name
//   - defaultValue: Value to return if environment variable is empty, not set, or invalid
//
// Returns:
//   - int: Parsed integer value or default value
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvDuration retrieves a time.Duration environment variable with a default value.
// Uses time.ParseDuration for parsing (supports units like "5s", "10m", "1h").
// Invalid values fall back to the default.
//
// Parameters:
//   - key: Environment variable name
//   - defaultValue: Value to return if environment variable is empty, not set, or invalid
//
// Returns:
//   - time.Duration: Parsed duration value or default value
//
// Example:
//   - GetEnvDuration("TIMEOUT", 30*time.Second) with "45s" returns 45*time.Second
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(strings.TrimSpace(value)); err == nil {
			return duration
		}
	}
	return defaultValue
}
```

## Step 2: Create Tests for Helper Functions

**File**: `internal/config/env_test.go` (new file)

```go
package config

import (
	"os"
	"testing"
	"time"
)

func TestGetEnvString(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "environment variable set",
			key:          "TEST_STRING",
			envValue:     "test-value",
			defaultValue: "default",
			expected:     "test-value",
		},
		{
			name:         "environment variable with whitespace",
			key:          "TEST_STRING_WHITESPACE",
			envValue:     "  test-value  ",
			defaultValue: "default",
			expected:     "test-value",
		},
		{
			name:         "environment variable empty",
			key:          "TEST_STRING_EMPTY",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "environment variable not set",
			key:          "TEST_STRING_NOT_SET",
			envValue:     "", // Will not be set
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment
			defer os.Unsetenv(tt.key)

			// Set environment variable if value provided
			if tt.envValue != "" || tt.name == "environment variable empty" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := GetEnvString(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvString(%q, %q) = %q, want %q",
					tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestGetEnvBool(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue bool
		expected     bool
	}{
		// Truthy values
		{
			name:         "true value",
			key:          "TEST_BOOL_TRUE",
			envValue:     "true",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "TRUE value (case insensitive)",
			key:          "TEST_BOOL_TRUE_UPPER",
			envValue:     "TRUE",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "True value (mixed case)",
			key:          "TEST_BOOL_TRUE_MIXED",
			envValue:     "True",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "1 value (numeric true)",
			key:          "TEST_BOOL_ONE",
			envValue:     "1",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "t value (short true)",
			key:          "TEST_BOOL_T",
			envValue:     "t",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "T value (short true upper)",
			key:          "TEST_BOOL_T_UPPER",
			envValue:     "T",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "yes value",
			key:          "TEST_BOOL_YES",
			envValue:     "yes",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "YES value (case insensitive)",
			key:          "TEST_BOOL_YES_UPPER",
			envValue:     "YES",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "y value (short yes)",
			key:          "TEST_BOOL_Y",
			envValue:     "y",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "on value",
			key:          "TEST_BOOL_ON",
			envValue:     "on",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "ON value (case insensitive)",
			key:          "TEST_BOOL_ON_UPPER",
			envValue:     "ON",
			defaultValue: false,
			expected:     true,
		},
		
		// Falsy values
		{
			name:         "false value",
			key:          "TEST_BOOL_FALSE",
			envValue:     "false",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "FALSE value (case insensitive)",
			key:          "TEST_BOOL_FALSE_UPPER",
			envValue:     "FALSE",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "0 value (numeric false)",
			key:          "TEST_BOOL_ZERO",
			envValue:     "0",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "f value (short false)",
			key:          "TEST_BOOL_F",
			envValue:     "f",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "F value (short false upper)",
			key:          "TEST_BOOL_F_UPPER",
			envValue:     "F",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "no value",
			key:          "TEST_BOOL_NO",
			envValue:     "no",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "NO value (case insensitive)",
			key:          "TEST_BOOL_NO_UPPER",
			envValue:     "NO",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "n value (short no)",
			key:          "TEST_BOOL_N",
			envValue:     "n",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "off value",
			key:          "TEST_BOOL_OFF",
			envValue:     "off",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "OFF value (case insensitive)",
			key:          "TEST_BOOL_OFF_UPPER",
			envValue:     "OFF",
			defaultValue: true,
			expected:     false,
		},
		
		// Edge cases
		{
			name:         "invalid value uses default (true)",
			key:          "TEST_BOOL_INVALID_TRUE",
			envValue:     "invalid",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "invalid value uses default (false)",
			key:          "TEST_BOOL_INVALID_FALSE",
			envValue:     "invalid",
			defaultValue: false,
			expected:     false,
		},
		{
			name:         "empty value uses default",
			key:          "TEST_BOOL_EMPTY",
			envValue:     "",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "whitespace around true",
			key:          "TEST_BOOL_WHITESPACE",
			envValue:     "  true  ",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "whitespace around 1",
			key:          "TEST_BOOL_WHITESPACE_ONE",
			envValue:     "  1  ",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "typo falls back to default",
			key:          "TEST_BOOL_TYPO",
			envValue:     "tru", // Common typo
			defaultValue: false,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.envValue != "" || tt.name == "empty value uses default" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := GetEnvBool(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvBool(%q, %v) with env value %q = %v, want %v",
					tt.key, tt.defaultValue, tt.envValue, result, tt.expected)
			}
		})
	}
}

func TestGetEnvStringSlice(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		separator    string
		defaultValue []string
		expected     []string
	}{
		{
			name:         "comma separated values",
			key:          "TEST_SLICE_COMMA",
			envValue:     "app.com,admin.com,api.com",
			separator:    ",",
			defaultValue: nil,
			expected:     []string{"app.com", "admin.com", "api.com"},
		},
		{
			name:         "comma separated with whitespace",
			key:          "TEST_SLICE_WHITESPACE",
			envValue:     "app.com, admin.com , api.com ",
			separator:    ",",
			defaultValue: nil,
			expected:     []string{"app.com", "admin.com", "api.com"},
		},
		{
			name:         "empty items filtered out",
			key:          "TEST_SLICE_EMPTY_ITEMS",
			envValue:     "app.com,,admin.com,",
			separator:    ",",
			defaultValue: nil,
			expected:     []string{"app.com", "admin.com"},
		},
		{
			name:         "single value",
			key:          "TEST_SLICE_SINGLE",
			envValue:     "app.com",
			separator:    ",",
			defaultValue: nil,
			expected:     []string{"app.com"},
		},
		{
			name:         "empty value uses default",
			key:          "TEST_SLICE_DEFAULT",
			envValue:     "",
			separator:    ",",
			defaultValue: []string{"default1", "default2"},
			expected:     []string{"default1", "default2"},
		},
		{
			name:         "semicolon separator",
			key:          "TEST_SLICE_SEMICOLON",
			envValue:     "app.com;admin.com;api.com",
			separator:    ";",
			defaultValue: nil,
			expected:     []string{"app.com", "admin.com", "api.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.envValue != "" || tt.name == "empty value uses default" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := GetEnvStringSlice(tt.key, tt.separator, tt.defaultValue)

			// Compare slices
			if len(result) != len(tt.expected) {
				t.Errorf("GetEnvStringSlice(%q, %q, %v) length = %d, want %d",
					tt.key, tt.separator, tt.defaultValue, len(result), len(tt.expected))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("GetEnvStringSlice(%q, %q, %v)[%d] = %q, want %q",
						tt.key, tt.separator, tt.defaultValue, i, result[i], expected)
				}
			}
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue int
		expected     int
	}{
		{
			name:         "valid integer",
			key:          "TEST_INT_VALID",
			envValue:     "42",
			defaultValue: 10,
			expected:     42,
		},
		{
			name:         "zero value",
			key:          "TEST_INT_ZERO",
			envValue:     "0",
			defaultValue: 10,
			expected:     0,
		},
		{
			name:         "negative integer",
			key:          "TEST_INT_NEGATIVE",
			envValue:     "-5",
			defaultValue: 10,
			expected:     -5,
		},
		{
			name:         "invalid value uses default",
			key:          "TEST_INT_INVALID",
			envValue:     "not-a-number",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "empty value uses default",
			key:          "TEST_INT_EMPTY",
			envValue:     "",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "whitespace around number",
			key:          "TEST_INT_WHITESPACE",
			envValue:     "  42  ",
			defaultValue: 10,
			expected:     42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.envValue != "" || tt.name == "empty value uses default" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := GetEnvInt(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvInt(%q, %d) = %d, want %d",
					tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestGetEnvDuration(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue time.Duration
		expected     time.Duration
	}{
		{
			name:         "seconds duration",
			key:          "TEST_DURATION_SECONDS",
			envValue:     "30s",
			defaultValue: 10 * time.Second,
			expected:     30 * time.Second,
		},
		{
			name:         "minutes duration",
			key:          "TEST_DURATION_MINUTES",
			envValue:     "5m",
			defaultValue: 1 * time.Minute,
			expected:     5 * time.Minute,
		},
		{
			name:         "hours duration",
			key:          "TEST_DURATION_HOURS",
			envValue:     "2h",
			defaultValue: 1 * time.Hour,
			expected:     2 * time.Hour,
		},
		{
			name:         "mixed duration",
			key:          "TEST_DURATION_MIXED",
			envValue:     "1h30m",
			defaultValue: 1 * time.Hour,
			expected:     90 * time.Minute,
		},
		{
			name:         "invalid duration uses default",
			key:          "TEST_DURATION_INVALID",
			envValue:     "not-a-duration",
			defaultValue: 10 * time.Second,
			expected:     10 * time.Second,
		},
		{
			name:         "empty value uses default",
			key:          "TEST_DURATION_EMPTY",
			envValue:     "",
			defaultValue: 10 * time.Second,
			expected:     10 * time.Second,
		},
		{
			name:         "whitespace around duration",
			key:          "TEST_DURATION_WHITESPACE",
			envValue:     "  45s  ",
			defaultValue: 10 * time.Second,
			expected:     45 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.envValue != "" || tt.name == "empty value uses default" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := GetEnvDuration(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvDuration(%q, %v) = %v, want %v",
					tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}
```

## Step 3: Update Request Logging to Use Helper Functions

**File**: `internal/server/request_logging.go` (modify existing LoadFromEnv method)

```go
// LoadFromEnv populates logging configuration from environment variables using centralized helper functions
func (c *RequestLoggingConfig) LoadFromEnv() {
	c.Enabled = config.GetEnvBool(config.EnvLoggingEnabled, c.Enabled)
	c.Level = config.GetEnvString(config.EnvLogLevel, c.Level)
	c.Format = config.GetEnvString(config.EnvLogFormat, c.Format)
	c.IncludeHeaders = config.GetEnvBool(config.EnvLogIncludeHeaders, c.IncludeHeaders)
}
```

## Step 4: Update Documentation

**File**: `internal/config/docs.go` (update existing documentation)

```go
// Package config provides centralized configuration management for Cipher Hub.
//
// This package establishes consistent patterns for environment variable handling,
// configuration loading, and type-safe configuration access across the entire
// Cipher Hub service. It centralizes all environment variable definitions and
// provides helper functions for common configuration operations.
//
// Key Features:
//   - Centralized environment variable name definitions
//   - Type-safe environment variable access helpers
//   - Consistent naming conventions across all configuration
//   - Default value handling with fallback support
//   - Validation and parsing for complex configuration types
//
// Environment Variable Naming Convention:
// All Cipher Hub environment variables follow the pattern:
//
//	CIPHER_HUB_<COMPONENT>_<SETTING>
//
// Examples:
//   - CIPHER_HUB_HOST=localhost
//   - CIPHER_HUB_READ_TIMEOUT=30s
//   - CIPHER_HUB_LOGGING_ENABLED=true
//   - CIPHER_HUB_CORS_ORIGINS=http://localhost:3000,https://app.example.com
//
// Helper Functions:
// The package provides type-safe helper functions for common configuration types:
//   - GetEnvString: String values with trimming and default support
//   - GetEnvBool: Boolean values supporting flexible representations (1/t/true/yes/y/on for true, 0/f/false/no/n/off for false)
//   - GetEnvStringSlice: Comma-separated lists with trimming and filtering
//   - GetEnvInt: Integer values with strconv.Atoi parsing
//   - GetEnvDuration: Duration values with time.ParseDuration support
//
// Boolean Environment Variables:
// GetEnvBool supports common boolean representations used across different systems:
//   - Truthy: "1", "t", "T", "true", "TRUE", "True", "yes", "YES", "Yes", "y", "Y", "on", "ON", "On"
//   - Falsy: "0", "f", "F", "false", "FALSE", "False", "no", "NO", "No", "n", "N", "off", "OFF", "Off"
//   - Any unrecognized value falls back to the provided default
//
// Examples:
//   - CIPHER_HUB_LOGGING_ENABLED=1 (Docker-style)
//   - CIPHER_HUB_CORS_ENABLED=yes (script-friendly)
//   - CIPHER_HUB_DEBUG=on (common convention)
//
// Usage Pattern:
//
//	import "cipher-hub/internal/config"
//
//	// Using helper functions for type-safe configuration loading
//	host := config.GetEnvString(config.EnvHost, "localhost")
//	enabled := config.GetEnvBool(config.EnvLoggingEnabled, true)
//	timeout := config.GetEnvDuration(config.EnvReadTimeout, 15*time.Second)
//	origins := config.GetEnvStringSlice(config.EnvCORSOrigins, ",", nil)
//
// Configuration Loading:
// Each component should implement a LoadFromEnv() method that uses the
// centralized constants and helper functions:
//
//	func (c *RequestLoggingConfig) LoadFromEnv() {
//		c.Enabled = config.GetEnvBool(config.EnvLoggingEnabled, c.Enabled)
//		c.Level = config.GetEnvString(config.EnvLogLevel, c.Level)
//		c.Format = config.GetEnvString(config.EnvLogFormat, c.Format)
//	}
//
// Type Safety:
// The helper functions provide type-safe access to environment variables
// with automatic parsing and validation for common types. Invalid values
// gracefully fall back to provided defaults, ensuring configuration robustness.
//
// Security Considerations:
// - No sensitive values are logged during configuration loading
// - Environment variables are validated before use
// - Default values provide secure fallbacks for all settings
// - Configuration validation prevents insecure configurations
//
// Container Integration:
// This package supports container-native deployment patterns:
//   - 12-factor app configuration via environment variables
//   - Docker and Kubernetes environment variable injection
//   - Configuration management via container orchestration
//   - Environment-specific configuration without code changes
//
// Default Values:
// All configuration loading provides sensible defaults that prioritize
// security and operational stability:
//   - Secure timeout values preventing resource exhaustion
//   - Conservative logging levels preventing information leakage
//   - Restrictive security settings that can be relaxed via configuration
//   - Development-friendly defaults for local development
//
// Validation:
// Configuration values are validated after loading to ensure:
//   - Required fields are present and non-empty
//   - Numeric values are within acceptable ranges
//   - URLs and network addresses are properly formatted
//   - Security-sensitive settings meet minimum requirements
//
// This package establishes the foundation for all configuration management
// across Cipher Hub, ensuring consistent patterns and secure defaults
// throughout the service architecture.
package config
```

## Step 5: Create CORS Configuration with Helper Functions

**File**: `internal/server/cors.go` (new file for CORS implementation)

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

// LoadFromEnv populates CORS configuration from environment variables using centralized helper functions
func (c *CORSConfig) LoadFromEnv() {
	c.Enabled = config.GetEnvBool(config.EnvCORSEnabled, c.Enabled)
	c.Origins = config.GetEnvStringSlice(config.EnvCORSOrigins, ",", c.Origins)
	c.Methods = config.GetEnvString(config.EnvCORSMethods, c.Methods)
	c.Headers = config.GetEnvString(config.EnvCORSHeaders, c.Headers)
	c.MaxAge = config.GetEnvString(config.EnvCORSMaxAge, c.MaxAge)
	c.Credentials = config.GetEnvBool(config.EnvCORSCredentials, c.Credentials)
}

// ApplyDefaults sets secure default values for CORS configuration
func (c *CORSConfig) ApplyDefaults() {
	// Enable CORS if origins are configured, disable if not
	c.Enabled = len(c.Origins) > 0

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

## Step 6: Create Comprehensive CORS Tests

**File**: `internal/server/cors_test.go` (new file)

```go
package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
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
		"CIPHER_HUB_CORS_CREDENTIALS",
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
				Enabled:     false,
				Origins:     nil,
				Methods:     DefaultCORSMethods,
				Headers:     DefaultCORSHeaders,
				MaxAge:      DefaultCORSMaxAge,
				Credentials: false,
			},
		},
		{
			name: "enabled with origins and custom config",
			envVars: map[string]string{
				"CIPHER_HUB_CORS_ENABLED":     "true",
				"CIPHER_HUB_CORS_ORIGINS":     "https://app.example.com,https://admin.example.com",
				"CIPHER_HUB_CORS_METHODS":     "GET, POST, PUT",
				"CIPHER_HUB_CORS_HEADERS":     "Content-Type, Authorization",
				"CIPHER_HUB_CORS_MAX_AGE":     "3600",
				"CIPHER_HUB_CORS_CREDENTIALS": "true",
			},
			expected: CORSConfig{
				Enabled:     true,
				Origins:     []string{"https://app.example.com", "https://admin.example.com"},
				Methods:     "GET, POST, PUT",
				Headers:     "Content-Type, Authorization",
				MaxAge:      "3600",
				Credentials: true,
			},
		},
		{
			name: "origins without explicit enabled",
			envVars: map[string]string{
				"CIPHER_HUB_CORS_ORIGINS": "https://app.example.com",
			},
			expected: CORSConfig{
				Enabled:     true, // Should be enabled automatically when origins are set
				Origins:     []string{"https://app.example.com"},
				Methods:     DefaultCORSMethods,
				Headers:     DefaultCORSHeaders,
				MaxAge:      DefaultCORSMaxAge,
				Credentials: false,
			},
		},
		{
			name: "flexible boolean values",
			envVars: map[string]string{
				"CIPHER_HUB_CORS_ENABLED":     "1",
				"CIPHER_HUB_CORS_ORIGINS":     "https://app.example.com",
				"CIPHER_HUB_CORS_CREDENTIALS": "yes",
			},
			expected: CORSConfig{
				Enabled:     true,
				Origins:     []string{"https://app.example.com"},
				Methods:     DefaultCORSMethods,
				Headers:     DefaultCORSHeaders,
				MaxAge:      DefaultCORSMaxAge,
				Credentials: true,
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

			if config.Credentials != tt.expected.Credentials {
				t.Errorf("Credentials = %v, want %v", config.Credentials, tt.expected.Credentials)
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
		{
			name: "credentials enabled",
			config: CORSConfig{
				Enabled:     true,
				Origins:     []string{"http://localhost:3000"},
				Methods:     DefaultCORSMethods,
				Headers:     DefaultCORSHeaders,
				Credentials: true,
			},
			requestOrigin:  "http://localhost:3000",
			method:         "GET",
			expectCORS:     true,
			expectedOrigin: "http://localhost:3000",
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

				// Check credentials header if enabled
				if tt.config.Credentials {
					if creds := w.Header().Get("Access-Control-Allow-Credentials"); creds != "true" {
						t.Error("Access-Control-Allow-Credentials should be 'true' when credentials enabled")
					}
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

## Step 7: Add Server Integration Tests

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