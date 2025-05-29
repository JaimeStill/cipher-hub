package config

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

	// SEcurity configuration
	EnvCORSOrigins = "CIPHER_HUB_CORS_ORIGINS"
	EnvTLSCertFile = "CIPHER_HUB_TLS_CERT_FILE"
	EnvTLSKeyFile  = "CIPHER_HUB_TLS_KEY_FILE"

	// Application configuration
	EnvEnvironment = "CIPHER_HUB_ENVIRONMENT"
	EnvDatabaseURL = "CIPHER_HUB_DATABASE_URL"
)
