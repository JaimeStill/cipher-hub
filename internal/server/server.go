package server

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Configuration constnats for defaults and validation bounds
const (
	// Default timeout values (secure defaults)
	DefaultReadTimeout     = 15 * time.Second
	DefaultWriteTimeout    = 15 * time.Second
	DefaultIdleTimeout     = 60 * time.Second
	DefaultShutdownTimeout = 30 * time.Second

	// Validation bounds for timeouts
	MinTimeout         = 1 * time.Second
	MaxTimeout         = 5 * time.Minute
	MaxShutdownTimeout = 2 * time.Minute

	// Error message prefix for consistent error handling
	ServerConfigErrorPrefix = "ServerConfig"
)

// ServerConfig holds all configuration parameters for the HTTP server.
// This struct provides a structured approach to server configuration that can
// be extended in future steps with environment variable loading and validation.
type ServerConfig struct {
	// Network configuration
	Host string `json:"host"`
	Port string `json:"port"`

	// Timeout configurations with zero values indicating defaults should be used
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
}

// ApplyDefaults applies secure default values to any zero-value timeout fields.
// This ensures the server has reasonable security defaults while allowing
// configuration override when needed.
func (c *ServerConfig) ApplyDefaults() {
	if c.ReadTimeout == 0 {
		c.ReadTimeout = DefaultReadTimeout
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = DefaultWriteTimeout
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = DefaultIdleTimeout
	}
	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = DefaultShutdownTimeout
	}
}

// Validate performs comprehensive validation of the server configuration.
// Returns detailed error messages for invalid configuration values following
// consistent error handling patterns.
func (c *ServerConfig) Validate() error {
	// Validate required fields
	if err := c.validateHost(); err != nil {
		return fmt.Errorf("%s: %w", ServerConfigErrorPrefix, err)
	}

	if err := c.validatePort(); err != nil {
		return fmt.Errorf("%s: %w", ServerConfigErrorPrefix, err)
	}

	// Validate timeout values
	if err := c.validateTimeouts(); err != nil {
		return fmt.Errorf("%s: %w", ServerConfigErrorPrefix, err)
	}

	return nil
}

// validateHost validates the host field using strict hostname validation
func (c *ServerConfig) validateHost() error {
	if c.Host == "" {
		return fmt.Errorf("host cannot be empty")
	}

	// Check if it's a valid IP address
	if ip := net.ParseIP(c.Host); ip != nil {
		return nil
	}

	if !isValidHostname(c.Host) {
		return fmt.Errorf("invalid host format: %s", c.Host)
	}

	return nil
}

// validatePort validates the port field
func (c *ServerConfig) validatePort() error {
	if c.Port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	portNum, err := strconv.Atoi(c.Port)
	if err != nil {
		return fmt.Errorf("invalid port format: %w", err)
	}

	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", portNum)
	}

	return nil
}

// validateTimeouts validates all timeout fields with sensible bounds
func (c *ServerConfig) validateTimeouts() error {
	timeouts := map[string]time.Duration{
		"read_timeout":     c.ReadTimeout,
		"write_timeout":    c.WriteTimeout,
		"idle_timeout":     c.IdleTimeout,
		"shutdown_timeout": c.ShutdownTimeout,
	}

	for name, timeout := range timeouts {
		if timeout < 0 {
			return fmt.Errorf("%s cannot be negative: %v", name, timeout)
		}

		// Check maximum bounds (shutdown timeout has higher limit)
		maxAllowed := MaxTimeout
		if name == "shutdown_timeout" {
			maxAllowed = MaxShutdownTimeout
		}

		if timeout > 0 && timeout < MinTimeout {
			return fmt.Errorf("%s must be at least %v, got %v", name, MinTimeout, timeout)
		}

		if timeout > maxAllowed {
			return fmt.Errorf("%s must not exceed %v, got %v", name, maxAllowed, timeout)
		}
	}

	return nil
}

// Address returns the full server address in host:port format
func (c *ServerConfig) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

// PortNum returns the configured port as an integer.
// This method assumes the port has been validated through Validate().
func (c *ServerConfig) PortNum() int {
	portNum, _ := strconv.Atoi(c.Port) // Safe after validation
	return portNum
}

// Server represents the HTTP server with configuration and lifecycle management
type Server struct {
	// Configuration
	config ServerConfig

	// Lifecycle management
	shutdownCtx    context.Context
	shutdownCancel context.CancelFunc

	// Server state
	started bool
}

// NewServer creates a new HTTP server instance with the specificed configuration.
// It validates the configuration, applies secure defaults, and prepares the server
// for lifecycle management.
//
// Parameters:
//   - config: ServerConfig containing host, port, and timeout configuration
//
// Returns:
//   - *Server: Configured server instance
//   - error: Validation error if configuration is invalid
//
// Security: Applies secure timeout defaults and validates all configuration parameters.
func NewServer(config ServerConfig) (*Server, error) {
	// Apply defaults for any zero-value timeout fields
	config.ApplyDefaults()

	// Validate the complete configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Create shutdown context with timeout for gracefuly lifecycle management
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)

	// Create server with validated configuration
	server := &Server{
		config: config,

		// Lifecycle management
		shutdownCtx:    shutdownCtx,
		shutdownCancel: shutdownCancel,
		started:        false,
	}

	return server, nil
}

// Config returns a copy of the server configuration
func (s *Server) Config() ServerConfig {
	return s.config
}

// Address returns the full server address in host:port format
func (s *Server) Address() string {
	return s.config.Address()
}

// Host returns the configured host address
func (s *Server) Host() string {
	return s.config.Host
}

// Port returns the configured port as a string
func (s *Server) Port() string {
	return s.config.Port
}

// PortNum returns the configured port as an integer
func (s *Server) PortNum() int {
	return s.config.PortNum()
}

// ReadTimeout returns the configured read timeout
func (s *Server) ReadTimeout() time.Duration {
	return s.config.ReadTimeout
}

// WriteTimeout returns the configured write timeout
func (s *Server) WriteTimeout() time.Duration {
	return s.config.WriteTimeout
}

// IdleTimeout returns the configured idle timeout
func (s *Server) IdleTimeout() time.Duration {
	return s.config.IdleTimeout
}

// ShutdownTimeout returns the configured shutdown timeout
func (s *Server) ShutdownTimeout() time.Duration {
	return s.config.ShutdownTimeout
}

// ShutdownContext returns the context used for graceful shutdown coordination
func (s *Server) ShutdownContext() context.Context {
	return s.shutdownCtx
}

// IsStarted returns whether the server has been started (for future use)
func (s *Server) IsStarted() bool {
	return s.started
}

// Shutdown initiates graceful shutdown by canceling the shutdown context
func (s *Server) Shutdown() {
	s.shutdownCancel()
}

// isvalidHostname performs strict hostname validation according to RFC standards
func isValidHostname(hostname string) bool {
	// Basic length checks
	if len(hostname) == 0 || len(hostname) > 253 {
		return false
	}

	// Check for localhost (always valid)
	if hostname == "localhost" {
		return true
	}

	// use URL parsing for strict validation
	testURL := "http://" + hostname
	parsedURL, err := url.Parse(testURL)
	if err != nil {
		return false
	}

	// Verify the hostname matches what we parsed
	if parsedURL.Hostname() != hostname {
		return false
	}

	// Additional checks for malicious input
	if strings.Contains(hostname, "..") ||
		strings.Contains(hostname, "<") ||
		strings.Contains(hostname, ">") ||
		strings.Contains(hostname, "'") ||
		strings.Contains(hostname, "\"") {
		return false
	}

	// Split into labels and validate each
	labels := strings.Split(hostname, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return false
		}

		// Label cannot start or end wiht hyphen
		if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return false
		}

		// Label must contain only valid characters
		for _, char := range label {
			if !((char >= 'a' && char <= 'z') ||
				(char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') ||
				char == '-') {
				return false
			}
		}
	}

	return true
}
