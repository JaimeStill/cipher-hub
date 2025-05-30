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
