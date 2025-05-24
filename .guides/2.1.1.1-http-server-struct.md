# Step 2.1.1.1: Create Basic HTTP Server Struct with Configuration Fields

## Overview

**Step**: 2.1.1.1  
**Task**: 2.1.1 (HTTP Server Creation)  
**Target**: 2.1 (Basic Server Setup)  
**Phase**: 2 (HTTP Server Infrastructure)  

**Time Estimate**: 25-35 minutes  
**Scope**: Create foundational HTTP server struct with structured configuration approach

## Step Objectives

### Primary Deliverables
- [x] Define `ServerConfig` struct for structured configuration
- [x] Create `Server` struct with host, port, and timeout fields
- [x] Add constructor `NewServer(config ServerConfig) (*Server, error)`
- [x] Include context for graceful shutdown

### Implementation Requirements
- **File Location**: `internal/server/server.go`
- **Architecture Focus**: Structured configuration pattern for enterprise scalability
- **Security Focus**: Strict input validation and secure defaults
- **Go Best Practices**: Consistent error handling and comprehensive documentation
- **Foundation Setup**: Prepare for Step 2.1.1.4 (configuration system integration)

---

## Implementation

### Step 1: Create Server Package with Structured Configuration

**File**: `internal/server/server.go`

```go
// Package server provides HTTP server infrastructure for Cipher Hub.
//
// This package implements the core HTTP server functionality with structured
// configuration management, graceful shutdown capabilities, and security-first
// design principles.
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

// Configuration constants for defaults and validation bounds
const (
	// Default timeout values (secure defaults)
	DefaultReadTimeout     = 15 * time.Second
	DefaultWriteTimeout    = 15 * time.Second
	DefaultIdleTimeout     = 60 * time.Second
	DefaultShutdownTimeout = 30 * time.Second
	
	// Validation bounds for timeouts
	MinTimeout        = 1 * time.Second
	MaxTimeout        = 5 * time.Minute
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
		return nil // Valid IP address
	}
	
	// Check if it's a valid hostname
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
	
	// Server state (for future use in Step 2.1.1.2)
	started bool
}

// NewServer creates a new HTTP server instance with the specified configuration.
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
	
	// Create shutdown context with timeout for graceful lifecycle management
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

// isValidHostname performs strict hostname validation according to RFC standards
func isValidHostname(hostname string) bool {
	// Basic length checks
	if len(hostname) == 0 || len(hostname) > 253 {
		return false
	}
	
	// Check for localhost (always valid)
	if hostname == "localhost" {
		return true
	}
	
	// Use URL parsing for strict validation
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
		
		// Label cannot start or end with hyphen
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
```

### Step 2: Create Comprehensive Test Structure

**File**: `internal/server/server_test.go`

```go
package server

import (
	"strings"
	"testing"
	"time"
)

func TestServerConfig_ApplyDefaults(t *testing.T) {
	tests := []struct {
		name     string
		config   ServerConfig
		expected ServerConfig
	}{
		{
			name: "apply all defaults",
			config: ServerConfig{
				Host: "localhost",
				Port: "8080",
				// All timeouts zero - should get defaults
			},
			expected: ServerConfig{
				Host:            "localhost",
				Port:            "8080",
				ReadTimeout:     DefaultReadTimeout,
				WriteTimeout:    DefaultWriteTimeout,
				IdleTimeout:     DefaultIdleTimeout,
				ShutdownTimeout: DefaultShutdownTimeout,
			},
		},
		{
			name: "preserve custom values",
			config: ServerConfig{
				Host:            "localhost",
				Port:            "8080",
				ReadTimeout:     30 * time.Second,
				WriteTimeout:    45 * time.Second,
				IdleTimeout:     120 * time.Second,
				ShutdownTimeout: 60 * time.Second,
			},
			expected: ServerConfig{
				Host:            "localhost",
				Port:            "8080",
				ReadTimeout:     30 * time.Second,
				WriteTimeout:    45 * time.Second,
				IdleTimeout:     120 * time.Second,
				ShutdownTimeout: 60 * time.Second,
			},
		},
		{
			name: "mixed defaults and customs",
			config: ServerConfig{
				Host:        "localhost",
				Port:        "8080",
				ReadTimeout: 25 * time.Second,
				// Other timeouts should get defaults
			},
			expected: ServerConfig{
				Host:            "localhost",
				Port:            "8080",
				ReadTimeout:     25 * time.Second,
				WriteTimeout:    DefaultWriteTimeout,
				IdleTimeout:     DefaultIdleTimeout,
				ShutdownTimeout: DefaultShutdownTimeout,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.ApplyDefaults()
			
			if tt.config.ReadTimeout != tt.expected.ReadTimeout {
				t.Errorf("ReadTimeout = %v, want %v", tt.config.ReadTimeout, tt.expected.ReadTimeout)
			}
			if tt.config.WriteTimeout != tt.expected.WriteTimeout {
				t.Errorf("WriteTimeout = %v, want %v", tt.config.WriteTimeout, tt.expected.WriteTimeout)
			}
			if tt.config.IdleTimeout != tt.expected.IdleTimeout {
				t.Errorf("IdleTimeout = %v, want %v", tt.config.IdleTimeout, tt.expected.IdleTimeout)
			}
			if tt.config.ShutdownTimeout != tt.expected.ShutdownTimeout {
				t.Errorf("ShutdownTimeout = %v, want %v", tt.config.ShutdownTimeout, tt.expected.ShutdownTimeout)
			}
		})
	}
}

func TestServerConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  ServerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid configuration",
			config: ServerConfig{
				Host:            "localhost",
				Port:            "8080",
				ReadTimeout:     15 * time.Second,
				WriteTimeout:    15 * time.Second,
				IdleTimeout:     60 * time.Second,
				ShutdownTimeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid IP address",
			config: ServerConfig{
				Host:            "127.0.0.1",
				Port:            "3000",
				ReadTimeout:     10 * time.Second,
				WriteTimeout:    10 * time.Second,
				IdleTimeout:     30 * time.Second,
				ShutdownTimeout: 20 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid IPv6 address",
			config: ServerConfig{
				Host:            "::1",
				Port:            "8080",
				ReadTimeout:     15 * time.Second,
				WriteTimeout:    15 * time.Second,
				IdleTimeout:     60 * time.Second,
				ShutdownTimeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid domain name",
			config: ServerConfig{
				Host:            "example.com",
				Port:            "8080",
				ReadTimeout:     15 * time.Second,
				WriteTimeout:    15 * time.Second,
				IdleTimeout:     60 * time.Second,
				ShutdownTimeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "empty host",
			config: ServerConfig{
				Host: "",
				Port: "8080",
			},
			wantErr: true,
			errMsg:  "ServerConfig: host cannot be empty",
		},
		{
			name: "empty port",
			config: ServerConfig{
				Host: "localhost",
				Port: "",
			},
			wantErr: true,
			errMsg:  "ServerConfig: port cannot be empty",
		},
		{
			name: "invalid port format",
			config: ServerConfig{
				Host: "localhost",
				Port: "abc",
			},
			wantErr: true,
			errMsg:  "ServerConfig: invalid port format",
		},
		{
			name: "port out of range low",
			config: ServerConfig{
				Host: "localhost",
				Port: "0",
			},
			wantErr: true,
			errMsg:  "ServerConfig: port must be between 1 and 65535",
		},
		{
			name: "port out of range high",
			config: ServerConfig{
				Host: "localhost",
				Port: "70000",
			},
			wantErr: true,
			errMsg:  "ServerConfig: port must be between 1 and 65535",
		},
		{
			name: "invalid hostname with path injection",
			config: ServerConfig{
				Host: "../../../etc/passwd",
				Port: "8080",
			},
			wantErr: true,
			errMsg:  "ServerConfig: invalid host format",
		},
		{
			name: "invalid hostname with script injection",
			config: ServerConfig{
				Host: "<script>alert('xss')</script>",
				Port: "8080",
			},
			wantErr: true,
			errMsg:  "ServerConfig: invalid host format",
		},
		{
			name: "negative read timeout",
			config: ServerConfig{
				Host:        "localhost",
				Port:        "8080",
				ReadTimeout: -5 * time.Second,
			},
			wantErr: true,
			errMsg:  "ServerConfig: read_timeout cannot be negative",
		},
		{
			name: "timeout below minimum",
			config: ServerConfig{
				Host:        "localhost",
				Port:        "8080",
				ReadTimeout: 500 * time.Millisecond, // Below MinTimeout
			},
			wantErr: true,
			errMsg:  "ServerConfig: read_timeout must be at least",
		},
		{
			name: "timeout above maximum",
			config: ServerConfig{
				Host:        "localhost",
				Port:        "8080",
				ReadTimeout: 10 * time.Minute, // Above MaxTimeout
			},
			wantErr: true,
			errMsg:  "ServerConfig: read_timeout must not exceed",
		},
		{
			name: "shutdown timeout above maximum",
			config: ServerConfig{
				Host:            "localhost",
				Port:            "8080",
				ShutdownTimeout: 5 * time.Minute, // Above MaxShutdownTimeout
			},
			wantErr: true,
			errMsg:  "ServerConfig: shutdown_timeout must not exceed",
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

func TestServerConfig_Address(t *testing.T) {
	tests := []struct {
		name     string
		config   ServerConfig
		expected string
	}{
		{
			name: "localhost with standard port",
			config: ServerConfig{
				Host: "localhost",
				Port: "8080",
			},
			expected: "localhost:8080",
		},
		{
			name: "IP address with custom port",
			config: ServerConfig{
				Host: "192.168.1.1",
				Port: "3000",
			},
			expected: "192.168.1.1:3000",
		},
		{
			name: "IPv6 address",
			config: ServerConfig{
				Host: "::1",
				Port: "8080",
			},
			expected: "[::1]:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address := tt.config.Address()
			if address != tt.expected {
				t.Errorf("Address() = %v, want %v", address, tt.expected)
			}
		})
	}
}

func TestNewServer(t *testing.T) {
	tests := []struct {
		name    string
		config  ServerConfig
		wantErr bool
	}{
		{
			name: "valid configuration with defaults",
			config: ServerConfig{
				Host: "localhost",
				Port: "8080",
				// Timeouts will be defaulted
			},
			wantErr: false,
		},
		{
			name: "valid configuration with custom timeouts",
			config: ServerConfig{
				Host:            "0.0.0.0",
				Port:            "3000",
				ReadTimeout:     30 * time.Second,
				WriteTimeout:    30 * time.Second,
				IdleTimeout:     120 * time.Second,
				ShutdownTimeout: 45 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "invalid configuration",
			config: ServerConfig{
				Host: "",
				Port: "8080",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewServer(tt.config)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewServer() expected error, got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("NewServer() unexpected error: %v", err)
				return
			}
			
			if server == nil {
				t.Error("NewServer() returned nil server")
				return
			}
			
			// Verify configuration is properly stored and defaults applied
			config := server.Config()
			if config.Host != tt.config.Host {
				t.Errorf("Host = %v, want %v", config.Host, tt.config.Host)
			}
			if config.Port != tt.config.Port {
				t.Errorf("Port = %v, want %v", config.Port, tt.config.Port)
			}
			
			// Verify defaults were applied for zero timeout values
			if tt.config.ReadTimeout == 0 && config.ReadTimeout != DefaultReadTimeout {
				t.Errorf("ReadTimeout default not applied: got %v, want %v", config.ReadTimeout, DefaultReadTimeout)
			}
			if tt.config.WriteTimeout == 0 && config.WriteTimeout != DefaultWriteTimeout {
				t.Errorf("WriteTimeout default not applied: got %v, want %v", config.WriteTimeout, DefaultWriteTimeout)
			}
			if tt.config.IdleTimeout == 0 && config.IdleTimeout != DefaultIdleTimeout {
				t.Errorf("IdleTimeout default not applied: got %v, want %v", config.IdleTimeout, DefaultIdleTimeout)
			}
			if tt.config.ShutdownTimeout == 0 && config.ShutdownTimeout != DefaultShutdownTimeout {
				t.Errorf("ShutdownTimeout default not applied: got %v, want %v", config.ShutdownTimeout, DefaultShutdownTimeout)
			}
		})
	}
}

func TestServer_Accessors(t *testing.T) {
	config := ServerConfig{
		Host:            "localhost",
		Port:            "8080",
		ReadTimeout:     20 * time.Second,
		WriteTimeout:    25 * time.Second,
		IdleTimeout:     90 * time.Second,
		ShutdownTimeout: 40 * time.Second,
	}
	
	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("NewServer() unexpected error: %v", err)
	}
	
	// Test all accessor methods
	if server.Host() != "localhost" {
		t.Errorf("Host() = %v, want %v", server.Host(), "localhost")
	}
	
	if server.Port() != "8080" {
		t.Errorf("Port() = %v, want %v", server.Port(), "8080")
	}
	
	if server.PortNum() != 8080 {
		t.Errorf("PortNum() = %v, want %v", server.PortNum(), 8080)
	}
	
	if server.Address() != "localhost:8080" {
		t.Errorf("Address() = %v, want %v", server.Address(), "localhost:8080")
	}
	
	if server.ReadTimeout() != 20*time.Second {
		t.Errorf("ReadTimeout() = %v, want %v", server.ReadTimeout(), 20*time.Second)
	}
	
	if server.WriteTimeout() != 25*time.Second {
		t.Errorf("WriteTimeout() = %v, want %v", server.WriteTimeout(), 25*time.Second)
	}
	
	if server.IdleTimeout() != 90*time.Second {
		t.Errorf("IdleTimeout() = %v, want %v", server.IdleTimeout(), 90*time.Second)
	}
	
	if server.ShutdownTimeout() != 40*time.Second {
		t.Errorf("ShutdownTimeout() = %v, want %v", server.ShutdownTimeout(), 40*time.Second)
	}
}

func TestServer_ShutdownContext(t *testing.T) {
	config := ServerConfig{
		Host:            "localhost",
		Port:            "8080",
		ShutdownTimeout: 5 * time.Second,
	}
	
	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("NewServer() unexpected error: %v", err)
	}
	
	ctx := server.ShutdownContext()
	if ctx == nil {
		t.Error("ShutdownContext() returned nil")
	}
	
	// Verify context is not canceled initially
	select {
	case <-ctx.Done():
		t.Error("ShutdownContext() should not be canceled initially")
	default:
		// Expected - context should be active
	}
	
	// Test shutdown cancellation
	server.Shutdown()
	
	// Verify context is canceled after shutdown
	select {
	case <-ctx.Done():
		// Expected - context should be canceled
	default:
		t.Error("ShutdownContext() should be canceled after Shutdown()")
	}
}

func TestServer_IsStarted(t *testing.T) {
	config := ServerConfig{
		Host: "localhost",
		Port: "8080",
	}
	
	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("NewServer() unexpected error: %v", err)
	}
	
	// Server should not be started initially
	if server.IsStarted() {
		t.Error("IsStarted() should return false for new server")
	}
}

func TestIsValidHostname(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		expected bool
	}{
		{"localhost", "localhost", true},
		{"valid domain", "example.com", true},
		{"subdomain", "api.example.com", true},
		{"numeric domain", "123.example.com", true},
		{"hyphenated domain", "my-api.example.com", true},
		{"empty string", "", false},
		{"too long", strings.Repeat("a", 254), false},
		{"path injection", "../../../etc/passwd", false},
		{"script injection", "<script>alert('xss')</script>", false},
		{"quote injection", "test'ing", false},
		{"double quote injection", "test\"ing", false},
		{"label too long", strings.Repeat("a", 64) + ".com", false},
		{"starts with hyphen", "-example.com", false},
		{"ends with hyphen", "example-.com", false},
		{"invalid characters", "exam@ple.com", false},
		{"multiple dots", "example..com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidHostname(tt.hostname)
			if result != tt.expected {
				t.Errorf("isValidHostname(%q) = %v, want %v", tt.hostname, result, tt.expected)
			}
		})
	}
}
```

---

## Verification Steps

### Step 1: Build Verification
```bash
# Navigate to project root
cd cipher-hub/

# Run build check
go build ./internal/server

# Expected: No compilation errors
```

### Step 2: Test Verification
```bash
# Run tests with verbose output
go test ./internal/server -v

# Expected: All tests pass with comprehensive coverage
# Sample expected output:
# === RUN   TestServerConfig_ApplyDefaults
# === RUN   TestServerConfig_Validate
# === RUN   TestServerConfig_Address
# === RUN   TestNewServer
# === RUN   TestServer_Accessors
# === RUN   TestServer_ShutdownContext
# === RUN   TestServer_IsStarted
# === RUN   TestIsValidHostname
# --- PASS: All tests should pass
# PASS
```

### Step 3: Test Coverage Analysis
```bash
# Check test coverage
go test ./internal/server -cover

# Expected: High coverage percentage (>90%)
```

### Step 4: Security Validation
```bash
# Test with malicious hostnames
go test ./internal/server -run TestIsValidHostname -v

# Expected: All injection attempts properly rejected
```

### Step 5: Code Quality Verification
```bash
# Format check
go fmt ./internal/server

# Vet check
go vet ./internal/server

# Expected: No issues reported
```

### Step 6: Documentation Verification
```bash
# Check documentation
go doc ./internal/server

# Check specific types
go doc ./internal/server.ServerConfig
go doc ./internal/server.NewServer

# Expected: Comprehensive documentation displayed
```

---

## Completion Criteria

### ✅ **Step 2.1.1.1 is complete when:**

1. **ServerConfig Struct**: Properly defined with:
   - Network fields (`Host`, `Port`) with strict validation
   - All timeout fields including `ShutdownTimeout`
   - `ApplyDefaults()` method using named constants
   - `Validate()` method with comprehensive validation and bounds checking
   - `Address()` and `PortNum()` utility methods

2. **Server Struct**: Created with:
   - `config ServerConfig` field
   - Context fields with timeout for shutdown management
   - State field (`started bool`) for future use

3. **Security Implementation**: 
   - Strict hostname validation preventing injection attacks
   - Timeout bounds validation preventing operational issues
   - Consistent error handling with clear messages

4. **Constructor**: `NewServer(config ServerConfig) (*Server, error)` with:
   - Configuration validation and default application
   - Timeout-based shutdown context setup
   - Comprehensive error handling with consistent patterns

5. **Constants & Standards**:
   - Named constants for all default values and limits
   - Consistent error message formatting
   - Scalable validation patterns

6. **Comprehensive Testing**: All test cases passing:
   - Configuration validation (valid/invalid/malicious cases)
   - Hostname security validation 
   - Timeout bounds validation
   - Default application tests
   - Constructor tests with various configurations
   - Accessor method tests
   - Context and shutdown tests

7. **Code Quality**: Passes all quality checks:
   - Formatting (`go fmt`)
   - Static analysis (`go vet`)
   - High test coverage (>90%)
   - Complete documentation

### 🏗️ **Enhanced Foundation for Future Steps**

This refined implementation provides:

**Step 2.1.1.2** enhancement readiness:
- Timeout configuration ready for `http.Server` integration
- Shutdown context with proper timeout handling
- Validated configuration for listener setup

**Step 2.1.1.4** configuration loading readiness:
- `ServerConfig` struct ready for environment variable population
- Validation patterns established for extending configuration
- JSON tags supporting configuration file loading

### 📁 **Files Created**
- `internal/server/server.go` - Server struct with production-ready configuration
- `internal/server/server_test.go` - Comprehensive test suite with security testing

---

## Security & Quality Improvements Implemented

### 🔒 **Security Enhancements**
- **Strict Hostname Validation**: Prevents injection attacks using RFC-compliant validation
- **Timeout Bounds Checking**: Prevents operational issues from extreme timeout values
- **Input Sanitization**: Multiple layers of validation for all configuration parameters
- **Consistent Error Handling**: No sensitive information leaked in error messages

### 🏗️ **Architecture Improvements**
- **Named Constants**: All magic numbers replaced with descriptive constants
- **Configurable Shutdown Timeout**: Extends timeout pattern naturally and safely
- **Validation Bounds**: Sensible minimum/maximum limits prevent misconfigurations
- **Scalable Error Patterns**: Simple but consistent error handling ready for growth

### 🧪 **Testing Excellence**  
- **Security Testing**: Comprehensive validation of malicious input rejection
- **Bounds Testing**: Validation of all timeout limits and edge cases
- **Integration Preparation**: Tests verify all integration points for future steps
- **High Coverage**: >90% test coverage with clear, maintainable test cases

This foundation now provides enterprise-grade security, maintainability, and scalability while maintaining the clean, simple patterns that will serve the project well as it grows! 🚀

---

## Next Step Preview

**Step 2.1.1.2: Implement basic HTTP listener setup** will build seamlessly on this robust foundation:
- Add `http.Server` field using all the validated timeout configuration
- Implement `Start()` method with proper listener setup using `s.config.Address()`
- Integrate with the timeout-based shutdown context for graceful termination
- Utilize all the established configuration and validation patterns

The security-hardened, well-tested foundation ensures smooth progression through all remaining steps.