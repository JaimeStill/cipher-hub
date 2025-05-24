# Cipher Hub - Development Style Guide

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

#### Structured Error Responses (Planned)
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

#### CORS Configuration Pattern
```go
// Environment-configurable CORS origins (not hard-coded)
func (c *Config) GetCORSOrigins() []string {
    switch c.Environment {
    case "development":
        return []string{"http://localhost:3000", "http://localhost:8080"}
    case "staging":
        return []string{"https://staging.cipher-hub.com"}
    case "production":
        return []string{"https://cipher-hub.com"}
    default:
        return []string{} // Deny all for unknown environments
    }
}
```

---

## Testing Requirements

### Test Coverage Standards
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

### Test Organization Best Practices
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

### Integration Testing
- **HTTP API Tests**: End-to-end API testing for all endpoints
- **Storage Tests**: Test storage implementations with real backends
- **Security Tests**: Verify authentication and authorization enforcement

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

## Development Workflow Standards

### Session-Based Development Pattern
- **Step Granularity**: Each development step should be completable in 20-30 minutes
- **Incremental Progress**: Each step builds upon previous completed work
- **Validation Requirements**: Each step includes verification and testing requirements
- **Documentation Updates**: Concurrent documentation updates with implementation

### Pre-Commit Standards
```bash
# Required commands before every commit
go fmt ./...
go build ./...
go test ./...
go vet ./...
go mod tidy

# Verify no security issues
go list -json -m all | nancy sleuth
```

### Commit Message Standards
```
type(scope): description

feat(server): implement HTTP server configuration with security validation
fix(validation): prevent hostname injection attacks in ServerConfig
docs(api): update API documentation for health check endpoints
test(security): add comprehensive input validation tests
```

---

## Git Workflow

### Branch Strategy
- **Feature Branches**: Create feature branches for each roadmap step
- **Naming**: Use descriptive branch names (e.g., `step/2.1.1.2-http-listener`)
- **Main Protection**: Main branch requires pull request reviews
- **No Direct Commits**: Never commit directly to main branch

### Code Review Requirements
- **Security Changes**: All security-related changes require thorough review
- **Test Coverage**: New code must include appropriate tests
- **Documentation**: Public API changes require documentation updates
- **Performance**: Performance-sensitive changes need benchmarks

---

## Security Checklist

### Before Every Commit
- [ ] No key material in logs or error messages
- [ ] Input validation on all user-provided data with injection prevention
- [ ] Proper error handling without information leakage
- [ ] Authentication and authorization checks in place
- [ ] Audit logging for security-relevant operations
- [ ] Timeout bounds validation to prevent resource exhaustion
- [ ] Environment variable configuration without hard-coded secrets

### Code Review Focus Areas
- [ ] Key material protection and secure handling
- [ ] Input validation and sanitization with security bounds
- [ ] Error handling and information disclosure prevention
- [ ] Authentication and authorization implementation
- [ ] Logging and audit trail completeness without sensitive data
- [ ] Configuration security and environment-based settings

---

*Style Guide Version: 1.3*  
*Last Updated: Current Session*  
*Status: Updated with Phase 2.1 HTTP Server Infrastructure Patterns*