# Cipher Hub - Implementation Style Guide

**Version**: 1.4  
**Last Updated**: Current Session

This style guide establishes consistent development practices for the Cipher Hub project. All code contributions must adhere to these standards to maintain security, quality, and maintainability.

---

## Go Code Standards

### Modern Go Idioms
- **Type Declarations**: Use `any` instead of `interface{}` for Go 1.18+
- **Error Handling**: All constructors return `(Type, error)` patterns, never panic
- **Error Wrapping**: Use `fmt.Errorf()` with `%w` verb for error chain preservation
- **Context Usage**: Pass `context.Context` as first parameter for all I/O operations
- **Nil Checks**: Always check for nil before dereferencing pointers

### Constructor Patterns

#### Structured Configuration Approach
```go
// Correct: Structured configuration with validation
type ServerConfig struct {
    Host            string        `json:"host"`
    Port            string        `json:"port"`
    ReadTimeout     time.Duration `json:"read_timeout"`
    WriteTimeout    time.Duration `json:"write_timeout"`
    IdleTimeout     time.Duration `json:"idle_timeout"`
    ShutdownTimeout time.Duration `json:"shutdown_timeout"`
}

func NewServer(config ServerConfig) (*Server, error) {
    // Apply defaults for zero values
    config.ApplyDefaults()
    
    // Validate configuration
    if err := config.Validate(); err != nil {
        return nil, err
    }
    
    // Create with validated configuration
    return &Server{config: config}, nil
}

// Incorrect: Individual parameters become unwieldy
func NewServer(host, port string, readTimeout, writeTimeout time.Duration) (*Server, error) {
    // This doesn't scale and becomes hard to extend
}
```

#### Configuration Validation Pattern
```go
// Establish consistent validation methods
func (c *ServerConfig) Validate() error {
    if err := c.validateHost(); err != nil {
        return fmt.Errorf("%s: %w", ServerConfigErrorPrefix, err)
    }
    if err := c.validatePort(); err != nil {
        return fmt.Errorf("%s: %w", ServerConfigErrorPrefix, err)
    }
    return nil
}

// Use named constants for error prefixes
const ServerConfigErrorPrefix = "ServerConfig"
```

#### Default Application Pattern
```go
// Provide secure defaults for zero values
func (c *ServerConfig) ApplyDefaults() {
    if c.ReadTimeout == 0 {
        c.ReadTimeout = DefaultReadTimeout
    }
    if c.WriteTimeout == 0 {
        c.WriteTimeout = DefaultWriteTimeout
    }
    // ... apply other defaults
}

// Define defaults as named constants
const (
    DefaultReadTimeout     = 15 * time.Second
    DefaultWriteTimeout    = 15 * time.Second
    DefaultIdleTimeout     = 60 * time.Second
    DefaultShutdownTimeout = 30 * time.Second
)
```

### Time Field Strategy
- **Required timestamps**: Use `time.Time` for semantic clarity
- **Optional timestamps**: Use `*time.Time` for JSON serialization benefits
- **Examples**: `CreatedAt time.Time`, `ExpiresAt *time.Time`

### Context Management Patterns

#### Typed Context Keys
```go
import (
    "context"
    "fmt"
    "log/slog"
    "net/url"
    "strings"
    "time"
)

// Correct: Use typed context keys for security
type contextKey string

const (
    loggerCtxKey    contextKey = "logger"
    requestIDCtxKey contextKey = "request_id"
    userCtxKey      contextKey = "user"
)

// Helper functions for type safety
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
    return context.WithValue(ctx, loggerCtxKey, logger)
}

func GetLogger(ctx context.Context) *slog.Logger {
    if logger, ok := ctx.Value(loggerCtxKey).(*slog.Logger); ok {
        return logger
    }
    return slog.Default() // Safe fallback
}

// Incorrect: String keys are error-prone
ctx = context.WithValue(ctx, "logger", logger) // Type unsafe
```

#### Context Timeout Management
```go
// Use WithTimeout for operations with time bounds
shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)

// Always provide cancellation cleanup
defer shutdownCancel()
```

---

## Security-First Development

### Input Validation Standards

#### Comprehensive Hostname Validation
```go
// Security-conscious hostname validation
func isValidHostname(hostname string) bool {
    // Basic length checks
    if len(hostname) == 0 || len(hostname) > 253 {
        return false
    }

    // Allow localhost
    if hostname == "localhost" {
        return true
    }

    // Use URL parsing for strict validation
    testURL := "http://" + hostname
    parsedURL, err := url.Parse(testURL)
    if err != nil {
        return false
    }

    // Verify hostname matches parsed result
    if parsedURL.Hostname() != hostname {
        return false
    }

    // Security checks for malicious input
    if strings.Contains(hostname, "..") ||
        strings.Contains(hostname, "<") ||
        strings.Contains(hostname, ">") ||
        strings.Contains(hostname, "'") ||
        strings.Contains(hostname, "\"") {
        return false
    }

    // Additional RFC compliance validation...
    return true
}
```

#### Port Validation with Security Bounds
```go
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
```

#### Timeout Bounds Validation
```go
// Define security bounds for timeouts
const (
    MinTimeout         = 1 * time.Second
    MaxTimeout         = 5 * time.Minute
    MaxShutdownTimeout = 2 * time.Minute
)

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

        // Check bounds
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
```

### Key Material Protection
- **Serialization**: Always use `json:"-"` tags on key material fields
- **Logging**: Never log key material, even in debug mode
- **Memory**: Clear sensitive data from memory when no longer needed
- **Error Messages**: Never include key material in error strings

```go
type CryptoKey struct {
    ID      string `json:"id"`
    KeyData []byte `json:"-"` // NEVER serialize key material
    // ... other fields
}
```

### Error Handling Security

#### Consistent Error Prefixes
```go
// Use typed error prefixes for consistent categorization
const (
    ServerConfigErrorPrefix = "ServerConfig"
    ServerErrorPrefix      = "Server"
    HandlerErrorPrefix     = "Handler"
    StorageErrorPrefix     = "Storage"
)

// Apply consistently in error handling
func (c *ServerConfig) Validate() error {
    if err := c.validateHost(); err != nil {
        return fmt.Errorf("%s: %w", ServerConfigErrorPrefix, err)
    }
    return nil
}
```

#### Structured Error Responses
```go
// Standard error response structure for APIs
type ErrorResponse struct {
    Error     string            `json:"error"`           // Human-readable message
    Code      string            `json:"code"`            // Machine-readable code
    Details   map[string]any    `json:"details,omitempty"` // Additional context
    RequestID string            `json:"request_id"`      // For tracing
    Timestamp time.Time         `json:"timestamp"`       // When error occurred
}

// Error codes constants
const (
    ErrCodeValidation     = "VALIDATION_ERROR"
    ErrCodeNotFound       = "NOT_FOUND"
    ErrCodeUnauthorized   = "UNAUTHORIZED"
    ErrCodeInternalError  = "INTERNAL_ERROR"
)
```

### Environment-Based Configuration

#### General Configuration Pattern
Configuration should be environment-driven with secure defaults and comprehensive validation. All configuration parameters should support environment variable loading with fallback to default values.

```go
import (
    "fmt"
    "net/url"
    "os"
    "strings"
    "time"
)

// Comprehensive configuration structure
type ServerConfig struct {
    // Network configuration
    Host        string        `json:"host"`
    Port        string        `json:"port"`
    
    // Timeout configurations
    ReadTimeout     time.Duration `json:"read_timeout"`
    WriteTimeout    time.Duration `json:"write_timeout"`
    IdleTimeout     time.Duration `json:"idle_timeout"`
    ShutdownTimeout time.Duration `json:"shutdown_timeout"`
    
    // Security configuration
    CORSOrigins     []string `json:"cors_origins"`
    TLSCertFile     string   `json:"tls_cert_file"`
    TLSKeyFile      string   `json:"tls_key_file"`
    
    // Application configuration
    Environment     string `json:"environment"`
    LogLevel        string `json:"log_level"`
    DatabaseURL     string `json:"database_url"`
}

// LoadFromEnv populates configuration from environment variables
func (c *ServerConfig) LoadFromEnv() error {
    // Network configuration
    if host := os.Getenv("CIPHER_HUB_HOST"); host != "" {
        c.Host = host
    }
    if port := os.Getenv("CIPHER_HUB_PORT"); port != "" {
        c.Port = port
    }
    
    // Timeout configuration with parsing
    if readTimeout := os.Getenv("CIPHER_HUB_READ_TIMEOUT"); readTimeout != "" {
        if duration, err := time.ParseDuration(readTimeout); err == nil {
            c.ReadTimeout = duration
        }
    }
    if writeTimeout := os.Getenv("CIPHER_HUB_WRITE_TIMEOUT"); writeTimeout != "" {
        if duration, err := time.ParseDuration(writeTimeout); err == nil {
            c.WriteTimeout = duration
        }
    }
    
    // Security configuration
    if corsOriginsEnv := os.Getenv("CIPHER_HUB_CORS_ORIGINS"); corsOriginsEnv != "" {
        c.CORSOrigins = parseCommaSeparatedList(corsOriginsEnv)
    }
    if certFile := os.Getenv("CIPHER_HUB_TLS_CERT_FILE"); certFile != "" {
        c.TLSCertFile = certFile
    }
    if keyFile := os.Getenv("CIPHER_HUB_TLS_KEY_FILE"); keyFile != "" {
        c.TLSKeyFile = keyFile
    }
    
    // Application configuration
    if env := os.Getenv("CIPHER_HUB_ENVIRONMENT"); env != "" {
        c.Environment = env
    }
    if logLevel := os.Getenv("CIPHER_HUB_LOG_LEVEL"); logLevel != "" {
        c.LogLevel = logLevel
    }
    if dbURL := os.Getenv("CIPHER_HUB_DATABASE_URL"); dbURL != "" {
        c.DatabaseURL = dbURL
    }
    
    return c.validateEnvironmentConfig()
}

// Helper function for parsing comma-separated lists
func parseCommaSeparatedList(value string) []string {
    items := strings.Split(value, ",")
    result := make([]string, 0, len(items))
    for _, item := range items {
        if trimmed := strings.TrimSpace(item); trimmed != "" {
            result = append(result, trimmed)
        }
    }
    return result
}

// Validate environment-specific configuration
func (c *ServerConfig) validateEnvironmentConfig() error {
    // Validate CORS origins
    for _, origin := range c.CORSOrigins {
        if _, err := url.Parse(origin); err != nil {
            return fmt.Errorf("invalid CORS origin format: %s", origin)
        }
    }
    
    // Validate TLS configuration
    if (c.TLSCertFile != "") != (c.TLSKeyFile != "") {
        return fmt.Errorf("both TLS cert and key files must be provided")
    }
    
    // Validate environment
    validEnvironments := []string{"development", "staging", "production"}
    if c.Environment != "" && !contains(validEnvironments, c.Environment) {
        return fmt.Errorf("invalid environment: %s", c.Environment)
    }
    
    // Validate log level
    validLogLevels := []string{"debug", "info", "warn", "error"}
    if c.LogLevel != "" && !contains(validLogLevels, c.LogLevel) {
        return fmt.Errorf("invalid log level: %s", c.LogLevel)
    }
    
    return nil
}

// Helper function for slice contains check
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

#### Environment Variable Naming Convention
Use consistent naming patterns for environment variables:
- **Prefix**: All variables start with `CIPHER_HUB_`
- **Format**: `CIPHER_HUB_<COMPONENT>_<SETTING>`
- **Case**: Use SCREAMING_SNAKE_CASE for environment variables
- **Examples**: 
  - `CIPHER_HUB_HOST=localhost`
  - `CIPHER_HUB_READ_TIMEOUT=30s`
  - `CIPHER_HUB_CORS_ORIGINS=http://localhost:3000,https://app.example.com`

#### Configuration Usage Pattern
```go
// Initialize configuration with environment loading
func NewServerFromEnv() (*Server, error) {
    config := ServerConfig{}
    
    // Apply secure defaults first
    config.ApplyDefaults()
    
    // Load from environment variables
    if err := config.LoadFromEnv(); err != nil {
        return nil, fmt.Errorf("failed to load environment config: %w", err)
    }
    
    // Create server with loaded configuration
    return NewServer(config)
}
```

---

## Testing Standards

### Test Coverage Requirements
- **Unit Tests**: Every public function and method must have unit tests
- **Table-Driven Tests**: Use table-driven tests for multiple input scenarios
- **Error Cases**: Test both success and failure paths
- **Edge Cases**: Include boundary conditions and edge cases
- **Security Tests**: Test malicious input and injection attempts

### Security-Focused Testing Patterns
```go
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
                Host: "localhost",
                Port: "8080",
                ReadTimeout: 15 * time.Second,
            },
            wantErr: false,
        },
        {
            name: "malicious hostname with path injection",
            config: ServerConfig{
                Host: "../../../etc/passwd",
                Port: "8080",
            },
            wantErr: true,
            errMsg:  "ServerConfig: invalid host format",
        },
        {
            name: "script injection attempt",
            config: ServerConfig{
                Host: "<script>alert('xss')</script>",
                Port: "8080",
            },
            wantErr: true,
            errMsg:  "ServerConfig: invalid host format",
        },
        {
            name: "timeout below minimum",
            config: ServerConfig{
                Host: "localhost",
                Port: "8080",
                ReadTimeout: 500 * time.Millisecond, // Below MinTimeout
            },
            wantErr: true,
            errMsg:  "ServerConfig: read_timeout must be at least",
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
```

### Constructor Testing Pattern
```go
// Comprehensive constructor testing
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

            // Verify defaults were applied
            config := server.Config()
            if tt.config.ReadTimeout == 0 && config.ReadTimeout != DefaultReadTimeout {
                t.Errorf("ReadTimeout default not applied: got %v, want %v", 
                    config.ReadTimeout, DefaultReadTimeout)
            }
        })
    }
}
```

### Integration Testing Pattern
```go
import (
    "fmt"
    "net/http"
    "testing"
)

// HTTP API integration testing
func TestHealthEndpoint(t *testing.T) {
    // Setup test server
    config := ServerConfig{
        Host: "localhost",
        Port: "0", // Use random port
    }
    server, err := NewServer(config)
    if err != nil {
        t.Fatalf("Failed to create server: %v", err)
    }

    // Start server in goroutine
    go func() {
        if err := server.Start(); err != nil {
            t.Errorf("Server failed to start: %v", err)
        }
    }()
    defer server.Shutdown()

    // Test health endpoint
    resp, err := http.Get(fmt.Sprintf("http://%s/health", server.Address()))
    if err != nil {
        t.Fatalf("Health check failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status 200, got %d", resp.StatusCode)
    }
}
```

---

## HTTP Server Architecture Patterns

### Middleware Design Pattern
```go
// Enhanced middleware stack with conditional support
type MiddlewareStack struct {
    middlewares []Middleware
}

func NewMiddlewareStack() *MiddlewareStack {
    return &MiddlewareStack{}
}

func (ms *MiddlewareStack) Use(middleware Middleware) *MiddlewareStack {
    ms.middlewares = append(ms.middlewares, middleware)
    return ms
}

func (ms *MiddlewareStack) UseIf(condition bool, middleware Middleware) *MiddlewareStack {
    if condition {
        ms.middlewares = append(ms.middlewares, middleware)
    }
    return ms
}

func (ms *MiddlewareStack) Apply(handler http.Handler) http.Handler {
    result := handler
    for i := len(ms.middlewares) - 1; i >= 0; i-- {
        result = ms.middlewares[i](result)
    }
    return result
}
```

### Health Check Interface Pattern
```go
import (
    "context"
    "time"
    
    "cipher-hub/internal/storage"
)

// Extensible health check pattern
type HealthChecker interface {
    Name() string
    Check(ctx context.Context) CheckResult
}

type CheckResult struct {
    Status  string        `json:"status"`
    Message string        `json:"message,omitempty"`
    Latency time.Duration `json:"latency_ms"`
    Details map[string]any `json:"details,omitempty"`
}

// Implementation example
type StorageHealthChecker struct {
    storage storage.Storage
}

func (s *StorageHealthChecker) Name() string {
    return "storage"
}

func (s *StorageHealthChecker) Check(ctx context.Context) CheckResult {
    start := time.Now()
    
    // Test storage operation
    _, err := s.storage.ListServiceRegistrations(ctx)
    latency := time.Since(start)
    
    if err != nil {
        return CheckResult{
            Status:  "unhealthy",
            Message: "Storage connectivity failed",
            Latency: latency,
            Details: map[string]any{"error": err.Error()},
        }
    }
    
    return CheckResult{
        Status:  "healthy",
        Message: "Storage operational",
        Latency: latency,
    }
}
```

### Request ID Generation Pattern
```go
import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "net/http"
)

// Secure request ID generation
func generateRequestID() (string, error) {
    bytes := make([]byte, 8)
    if _, err := rand.Read(bytes); err != nil {
        return "", fmt.Errorf("failed to generate request ID: %w", err)
    }
    return hex.EncodeToString(bytes), nil
}

// Use in middleware
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID, err := generateRequestID()
        if err != nil {
            // Handle error appropriately
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        
        ctx := WithRequestID(r.Context(), requestID)
        w.Header().Set("X-Request-ID", requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### JSON Handler Utilities
```go
import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// Standard JSON response utilities
func WriteJSON(w http.ResponseWriter, status int, data any) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    
    if err := json.NewEncoder(w).Encode(data); err != nil {
        return fmt.Errorf("failed to encode JSON response: %w", err)
    }
    return nil
}

func ReadJSON(r *http.Request, dst any) error {
    // Limit request body size to prevent abuse
    r.Body = http.MaxBytesReader(nil, r.Body, 1048576) // 1MB limit
    
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields() // Strict parsing
    
    if err := decoder.Decode(dst); err != nil {
        return fmt.Errorf("failed to decode JSON request: %w", err)
    }
    
    return nil
}

// Error response helper
func WriteError(w http.ResponseWriter, status int, message string, requestID string) {
    response := ErrorResponse{
        Error:     message,
        Code:      getErrorCode(status),
        RequestID: requestID,
        Timestamp: time.Now(),
    }
    
    WriteJSON(w, status, response)
}

func getErrorCode(status int) string {
    switch status {
    case http.StatusBadRequest:
        return ErrCodeValidation
    case http.StatusNotFound:
        return ErrCodeNotFound
    case http.StatusUnauthorized:
        return ErrCodeUnauthorized
    default:
        return ErrCodeInternalError
    }
}
```

---

## Documentation Standards

### Go Doc Comments
- **Public APIs**: All public functions, types, and methods require doc comments
- **Package Documentation**: Every package needs a package-level doc comment
- **Security Notes**: Document security considerations for sensitive operations
- **Usage Examples**: Include examples for complex APIs

```go
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
// Security: Applies secure timeout defaults and validates all configuration parameters
// to prevent injection attacks and resource exhaustion.
func NewServer(config ServerConfig) (*Server, error) {
    // ... implementation
}
```

### Package Documentation Pattern
```go
// Package server provides HTTP server infrastructure for Cipher Hub.
//
// This package implements the core HTTP server functionality with structured
// configuration management, graceful shutdown capabilities, and security-first
// design principles.
//
// Key Security Features:
//   - Comprehensive input validation with injection prevention
//   - Configurable timeouts with security bounds
//   - Context-based lifecycle management
//   - Environment-configurable CORS origins
//
// Usage:
//
//	config := ServerConfig{
//		Host: "localhost",
//		Port: "8080",
//	}
//	server, err := NewServer(config)
//	if err != nil {
//		log.Fatal(err)
//	}
package server
```

---

---

## File Organization

### Naming Conventions
- **Files**: Use `snake_case` for file names (e.g., `service_registration.go`)
- **Types**: One primary type per file with co-located tests
- **Tests**: Co-locate tests as `type_test.go` alongside `type.go`
- **Packages**: Use descriptive, lowercase package names

### Directory Structure
```
cipher-hub/
├── cmd/cipher-hub/           # Application entry points
├── internal/
│   ├── models/              # Core data models
│   ├── storage/             # Storage interface and implementations
│   ├── server/              # HTTP server infrastructure
│   └── handlers/            # HTTP request handlers
├── pkg/                     # Public libraries (if any)
└── docs/                    # Documentation
```

### File Content Organization
- **Imports**: Group standard library, third-party, and internal imports
- **Constants**: Define package-level constants after imports
- **Types**: Define types before functions that use them
- **Constructors**: Place constructors immediately after type definitions
- **Methods**: Group methods by receiver type

---

## Security Implementation Checklist

### Code Security Requirements
- [ ] No key material in logs or error messages
- [ ] Input validation on all user-provided data with injection prevention
- [ ] Proper error handling without information leakage
- [ ] Authentication and authorization checks in place
- [ ] Audit logging for security-relevant operations
- [ ] Timeout bounds validation to prevent resource exhaustion
- [ ] Environment variable configuration without hard-coded secrets

### Key Material Protection Standards
- [ ] Use `json:"-"` tags on all key material fields
- [ ] Never log key material under any circumstances
- [ ] Clear sensitive data from memory when no longer needed
- [ ] Validate all key material operations with proper error handling
- [ ] Implement secure key generation using `crypto/rand`
- [ ] Protect against timing attacks in key comparison operations

### Input Validation Requirements
- [ ] Validate all configuration parameters with security bounds
- [ ] Prevent path injection attacks in hostname validation
- [ ] Prevent script injection attacks in all user inputs
- [ ] Implement proper bounds checking for numeric inputs
- [ ] Use structured validation with consistent error messages
- [ ] Apply timeout limits to prevent resource exhaustion

---

*Implementation Style Guide Version: 1.4*  
*Last Updated: Current Session*  
*Focus: Security-first implementation patterns with comprehensive validation and testing standards*