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
