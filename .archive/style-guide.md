# Cipher Hub - Implementation Style Guide

**Version**: 1.7  
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

#### Context Coordination vs Operation Timeout Pattern
```go
// Correct: Use WithCancel for coordination, separate timeout for operations
shutdownCtx, shutdownCancel := context.WithCancel(context.Background())

// In operation method, apply specific timeout
func (s *Server) Shutdown() error {
    // Use context with specific timeout for HTTP server shutdown
    shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
    defer cancel()
    
    if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
        return fmt.Errorf("graceful shutdown failed: %w", err)
    }
    return nil
}

// Incorrect: Using timeout context for coordination leads to complexity
shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
// This makes the coordination context timeout, which may not be desired
```

---

## Thread Safety Patterns

### Production-Ready Concurrent Access
```go
import (
    "sync"
)

// Server struct with thread safety
type Server struct {
    // Configuration (immutable after creation)
    config ServerConfig

    // Middleware stack for request processing
    middleware *MiddlewareStack

    // HTTP server instance for lifecycle management
    httpServer *http.Server

    // Lifecycle management
    shutdownCtx    context.Context
    shutdownCancel context.CancelFunc

    // Server state (protected by mutex)
    started  bool
    disposed bool // Prevents restart after shutdown
    mu       sync.RWMutex // Protects mutable state
}

// Thread-safe state access patterns
func (s *Server) IsStarted() bool {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.started
}

// Thread-safe state modification patterns
func (s *Server) Start() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Check disposed state before allowing operation
    if s.disposed {
        return fmt.Errorf("%s: cannot start server after shutdown", ServerErrorPrefix)
    }

    // Check state before modification
    if s.started {
        return fmt.Errorf("%s: server already started", ServerErrorPrefix)
    }

    // Perform operations...
    
    // Update state atomically with operation completion
    s.started = true
    return nil
}
```

### Goroutine Management Patterns
```go
// Channel-based coordination for goroutine startup
func (s *Server) Start() error {
    // ... state checks and setup ...

    // Channel for server readiness signaling
    ready := make(chan struct{})

    // Start server in goroutine with proper coordination
    go func() {
        defer func() {
            // Always update state on goroutine exit
            s.mu.Lock()
            s.started = false
            s.mu.Unlock()
            listener.Close()
        }()

        // Signal readiness before serving
        close(ready)

        // Serve with proper error handling
        if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
            log.Printf("HTTP server error: %v", err)
        }
    }()

    // Update state before waiting for readiness
    s.started = true

    // Wait for server to be ready (avoids race conditions)
    <-ready

    return nil
}
```

### Mutex Selection Guidelines
```go
// Use sync.RWMutex for read-heavy, write-occasional scenarios
type Server struct {
    // ... other fields ...
    mu sync.RWMutex // Many reads (IsStarted), few writes (Start/Stop)
}

// Use sync.Mutex for write-heavy or simple locking scenarios
type Counter struct {
    value int
    mu    sync.Mutex // Simple increment/decrement operations
}
```

---

## HTTP Server Lifecycle Patterns

### Complete Server Lifecycle Implementation
```go
// HTTP Server lifecycle with proper resource management
func (s *Server) Start() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Check disposed state first
    if s.disposed {
        return fmt.Errorf("%s: cannot start server after shutdown", ServerErrorPrefix)
    }

    // Idempotent behavior - prevent double start
    if s.started {
        return fmt.Errorf("%s: server already started", ServerErrorPrefix)
    }

    // Apply middleware to root handler
    finalHandler := s.middleware.Apply(s.rootHandler)

    // Create HTTP server with validated configuration
    s.httpServer = &http.Server{
        Addr:         s.config.Address(),
        Handler:      finalHandler, // Use middleware-wrapped handler
        ReadTimeout:  s.config.ReadTimeout,
        WriteTimeout: s.config.WriteTimeout,
        IdleTimeout:  s.config.IdleTimeout,
    }

    // Store address for error handling before potential cleanup
    addr := s.httpServer.Addr

    // Create listener with error handling
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        s.httpServer = nil // Clean up on failure
        return fmt.Errorf("%s: failed to create listener on %s: %w", 
            ServerErrorPrefix, addr, err)
    }

    // Channel-based readiness coordination
    ready := make(chan struct{})

    // Start server in goroutine with proper lifecycle management
    go func() {
        defer func() {
            s.mu.Lock()
            s.started = false
            s.mu.Unlock()
            listener.Close()
        }()

        // Signal readiness before serving
        close(ready)

        // Serve with proper error handling
        if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
            log.Printf("HTTP server error: %v", err)
        }
    }()

    // Update state atomically
    s.started = true

    // Wait for server readiness
    <-ready

    return nil
}
```

### Graceful Shutdown Implementation Pattern
```go
// Production-ready graceful shutdown with resource cleanup
func (s *Server) Shutdown() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Cancel coordination context for any waiting operations
    defer s.shutdownCancel()

    // Mark as disposed to prevent restart
    s.disposed = true

    // Check if server is already shut down (idempotent behavior)
    if !s.started || s.httpServer == nil {
        return nil
    }

    // Create timeout context for HTTP server shutdown
    shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
    defer cancel()

    // Perform graceful HTTP server shutdown
    if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
        // If graceful shutdown fails, clean up state anyway
        s.started = false
        s.httpServer = nil
        return fmt.Errorf(
            "%s: graceful shutdown failed after %v timeout: %w",
            ServerErrorPrefix,
            s.config.ShutdownTimeout,
            err,
        )
    }

    // Update server state after successful shutdown
    s.started = false
    s.httpServer = nil

    return nil
}
```

### Resource Cleanup Patterns
```go
// Always clean up resources on failure
func (s *Server) Start() error {
    // ... setup code ...
    
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        s.httpServer = nil // Clean up server instance
        return fmt.Errorf("%s: failed to create listener on %s: %w", 
            ServerErrorPrefix, addr, err)
    }
    
    // ... continue with success path ...
}

// Use defer for automatic cleanup in goroutines
go func() {
    defer func() {
        s.mu.Lock()
        s.started = false
        s.mu.Unlock()
        listener.Close() // Always close listener
    }()
    
    // ... server logic ...
}()
```

---

## Middleware Architecture Patterns

### Industry-Standard Middleware Pattern
```go
import (
    "net/http"
)

// Middleware defines the standard middleware function signature.
// A middleware function takes an http.Handler and returns an http.Handler,
// allowing for request/response processing before and after the wrapped handler.
//
// Example usage:
//
//	func LoggingMiddleware(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			log.Printf("Request: %s %s", r.Method, r.URL.Path)
//			next.ServeHTTP(w, r)
//		})
//	}
type Middleware func(http.Handler) http.Handler

// Enhanced MiddlewareStack with conditional support
type MiddlewareStack struct {
    middlewares []Middleware
}

func NewMiddlewareStack() *MiddlewareStack {
    return &MiddlewareStack{
        middlewares: make([]Middleware, 0),
    }
}

// Use adds middleware that will always be applied
func (ms *MiddlewareStack) Use(middleware Middleware) *MiddlewareStack {
    ms.middlewares = append(ms.middlewares, middleware)
    return ms
}

// UseIf conditionally adds middleware based on the provided condition
func (ms *MiddlewareStack) UseIf(condition bool, middleware Middleware) *MiddlewareStack {
    if condition {
        ms.middlewares = append(ms.middlewares, middleware)
    }
    return ms
}

// Apply wraps the provided handler with all registered middleware functions.
// Middleware is applied in registration order to achieve reverse execution order
// (last registered becomes outermost).
func (ms *MiddlewareStack) Apply(handler http.Handler) http.Handler {
    if handler == nil {
        handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            http.NotFound(w, r)
        })
    }

    result := handler
    
    // Apply middleware in registration order to make last registered outermost
    for i := 0; i < len(ms.middlewares); i++ {
        result = ms.middlewares[i](result)
    }
    
    return result
}
```

### Conditional Middleware Patterns
```go
// Environment-based middleware configuration
func setupMiddleware(server *Server, config AppConfig) {
    server.Middleware().
        Use(RequestIDMiddleware()).              // Always generate request IDs
        Use(RequestLoggingMiddleware()).         // Always log requests
        UseIf(config.EnableCORS, CORSMiddleware(config.CORSOrigins)). // Environment-specific CORS
        UseIf(config.IsProduction, HSTSMiddleware()).                 // HSTS only in production
        UseIf(config.EnableAuth, AuthMiddleware()).                   // Auth when enabled
        Use(SecurityHeadersMiddleware())         // Always apply security headers
}

// Correct: Environment-configurable CORS middleware
func CORSMiddleware(origins []string) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if len(origins) > 0 {
                w.Header().Set("Access-Control-Allow-Origin", strings.Join(origins, ","))
                w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
                w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            }
            
            // Handle preflight requests
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// Incorrect: Hard-coded CORS origins
func CORSMiddleware() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Hard-coded
            next.ServeHTTP(w, r)
        })
    }
}
```

### Request Logging Middleware Pattern
```go
import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "log/slog"
    "net/http"
    "time"
)

// Secure request ID generation
func generateRequestID() (string, error) {
    bytes := make([]byte, 8)
    if _, err := rand.Read(bytes); err != nil {
        return "", fmt.Errorf("failed to generate request ID: %w", err)
    }
    return hex.EncodeToString(bytes), nil
}

// Request logging middleware with structured logging
func RequestLoggingMiddleware() Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Generate secure request ID
            requestID, err := generateRequestID()
            if err != nil {
                slog.Error("Failed to generate request ID", "error", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }
            
            // Add request ID to context
            ctx := context.WithValue(r.Context(), requestIDCtxKey, requestID)
            r = r.WithContext(ctx)
            
            // Add request ID to response headers
            w.Header().Set("X-Request-ID", requestID)
            
            // Log request start
            slog.Info("Request started",
                "request_id", requestID,
                "method", r.Method,
                "path", r.URL.Path,
                "remote_addr", r.RemoteAddr,
                "user_agent", r.UserAgent())
            
            // Wrap ResponseWriter to capture status code
            wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
            
            // Call next handler
            next.ServeHTTP(wrapped, r)
            
            // Log request completion
            duration := time.Since(start)
            slog.Info("Request completed",
                "request_id", requestID,
                "method", r.Method,
                "path", r.URL.Path,
                "status_code", wrapped.statusCode,
                "duration_ms", duration.Milliseconds(),
                "bytes_written", wrapped.bytesWritten)
        })
    }
}

// ResponseWriter wrapper to capture status code and bytes written
type responseWriter struct {
    http.ResponseWriter
    statusCode   int
    bytesWritten int64
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    n, err := rw.ResponseWriter.Write(b)
    rw.bytesWritten += int64(n)
    return n, err
}
```

### Security Headers Middleware Pattern
```go
// Security headers middleware with conditional HSTS
func SecurityHeadersMiddleware(isHTTPS bool) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Always apply these security headers
            w.Header().Set("X-Content-Type-Options", "nosniff")
            w.Header().Set("X-Frame-Options", "DENY")
            w.Header().Set("X-XSS-Protection", "1; mode=block")
            w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
            
            // Content Security Policy
            w.Header().Set("Content-Security-Policy", 
                "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")
            
            // Apply HSTS only for HTTPS
            if isHTTPS {
                w.Header().Set("Strict-Transport-Security", 
                    "max-age=31536000; includeSubDomains; preload")
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// Usage with conditional application
server.Middleware().
    UseIf(config.IsHTTPS, SecurityHeadersMiddleware(true)).
    UseIf(!config.IsHTTPS, SecurityHeadersMiddleware(false))
```

### Server Integration Pattern
```go
// Server struct with middleware integration
type Server struct {
    // Configuration
    config ServerConfig

    // Middleware stack for request processing
    middleware *MiddlewareStack

    // HTTP server instance for lifecycle management
    httpServer *http.Server

    // Root handler for middleware application
    rootHandler http.Handler

    // ... other fields
}

// NewServer with middleware initialization
func NewServer(config ServerConfig) (*Server, error) {
    // ... configuration validation ...

    server := &Server{
        config: config,

        // Initialize middleware stack
        middleware: NewMiddlewareStack(),

        // Root handler (will be set by user or default to NotFound)
        rootHandler: nil,

        // ... other field initialization
    }

    return server, nil
}

// Middleware accessor for configuration
func (s *Server) Middleware() *MiddlewareStack {
    return s.middleware
}

// SetHandler for root handler management
func (s *Server) SetHandler(handler http.Handler) {
    s.rootHandler = handler
}

// Start method with middleware application
func (s *Server) Start() error {
    // ... state checks ...

    // Apply middleware to root handler
    finalHandler := s.middleware.Apply(s.rootHandler)

    // Create HTTP server with middleware-wrapped handler
    s.httpServer = &http.Server{
        Addr:         s.config.Address(),
        Handler:      finalHandler, // Middleware applied here
        ReadTimeout:  s.config.ReadTimeout,
        WriteTimeout: s.config.WriteTimeout,
        IdleTimeout:  s.config.IdleTimeout,
    }

    // ... continue with server start logic
}
```

---

## Signal Handling Patterns

### Container-Native Signal Handling
```go
import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"
)

// Production-ready signal handling in main.go
func main() {
    // ... server setup ...

    // Create channel for shutdown signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    // Start graceful shutdown handler BEFORE server start to prevent races
    go func() {
        sig := <-sigChan
        log.Printf("Received signal %v, initiating graceful shutdown...", sig)
        
        // Use server's configured timeout plus buffer for coordination
        shutdownTimeout := srv.ShutdownTimeout() + 5*time.Second
        shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
        defer cancel()

        // Channel to signal shutdown completion
        done := make(chan error, 1)
        
        // Perform shutdown in goroutine
        go func() {
            done <- srv.Shutdown()
        }()

        // Wait for shutdown completion or timeout
        select {
        case err := <-done:
            if err != nil {
                log.Printf("Server shutdown failed: %v", err)
                os.Exit(1)
            }
            log.Println("Graceful shutdown completed")
            os.Exit(0)
        case <-shutdownCtx.Done():
            log.Printf("Shutdown timeout (%v) exceeded, forcing exit", shutdownTimeout)
            os.Exit(1)
        }
    }()

    // Start the server
    if err := srv.Start(); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }

    // Block forever, shutdown handled by signal handler
    select {}
}
```

### Signal Handler Design Principles
```go
// Correct: Signal handler setup before server start
func main() {
    // Create server
    srv, err := server.NewServer(config)
    
    // Set up signal handling BEFORE starting server
    setupSignalHandling(srv)
    
    // Start server after signal handler is ready
    srv.Start()
}

// Incorrect: Signal handler after server start creates race condition
func main() {
    srv, _ := server.NewServer(config)
    srv.Start() // Race condition: signals could arrive before handler setup
    setupSignalHandling(srv)
}
```

### Timeout Coordination Pattern
```go
// Correct: Server timeout + buffer prevents coordination timeout
shutdownTimeout := srv.ShutdownTimeout() + 5*time.Second
shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)

// The buffer prevents the coordination context from timing out before
// the actual shutdown operation completes

// Incorrect: Using same timeout for coordination and operation
shutdownTimeout := srv.ShutdownTimeout()
shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
// This can cause coordination timeout if shutdown takes exactly the configured time
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

#### Enhanced Port Validation with Dynamic Assignment Support
```go
func (c *ServerConfig) validatePort() error {
    if c.Port == "" {
        return fmt.Errorf("port cannot be empty")
    }

    portNum, err := strconv.Atoi(c.Port)
    if err != nil {
        return fmt.Errorf("invalid port format: %w", err)
    }

    // Support port "0" for OS dynamic assignment
    if portNum < 0 || portNum > 65535 {
        return fmt.Errorf("port must be between 0 and 65535, got %d", portNum)
    }

    return nil
}

// Document port semantics clearly
// PortNum returns the configured port as an integer.
// Port 0 indicates dynamic port assignment by the OS.
// This method assumes the port has been validated through Validate().
func (c *ServerConfig) PortNum() int {
    portNum, _ := strconv.Atoi(c.Port) // Safe after validation
    return portNum
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

#### Secure Error Information Handling
```go
// Correct: Safe error messages with necessary context
return fmt.Errorf("%s: failed to create listener on %s: %w", 
    ServerErrorPrefix, s.httpServer.Addr, err)

// Incorrect: Could leak sensitive configuration
return fmt.Errorf("failed to bind with config %+v: %w", s.config, err)
```

#### Shutdown Error Handling Pattern
```go
// Correct: Comprehensive shutdown error handling with cleanup
func (s *Server) Shutdown() error {
    // ... shutdown logic ...
    
    if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
        // Clean up state even on error to prevent inconsistent state
        s.started = false
        s.httpServer = nil
        return fmt.Errorf(
            "%s: graceful shutdown failed after %v timeout: %w",
            ServerErrorPrefix,
            s.config.ShutdownTimeout,
            err,
        )
    }
    
    // Success path cleanup
    s.started = false
    s.httpServer = nil
    return nil
}

// Incorrect: Incomplete cleanup on error
if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
    return err // State not cleaned up
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
- **Thread Safety Tests**: Test concurrent access patterns where applicable
- **Shutdown Tests**: Test graceful shutdown scenarios and resource cleanup
- **Middleware Tests**: Test middleware composition, execution order, and conditional application

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
            name: "port zero for dynamic assignment",
            config: ServerConfig{
                Host: "localhost",
                Port: "0", // Should be valid for OS dynamic assignment
            },
            wantErr: false,
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

### Middleware Testing Patterns
```go
func TestMiddlewareStack_Apply_Order(t *testing.T) {
    stack := NewMiddlewareStack()

    var executionOrder []string

    // Add middleware that tracks execution order
    middleware1 := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            executionOrder = append(executionOrder, "middleware1-before")
            next.ServeHTTP(w, r)
            executionOrder = append(executionOrder, "middleware1-after")
        })
    }

    middleware2 := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            executionOrder = append(executionOrder, "middleware2-before")
            next.ServeHTTP(w, r)
            executionOrder = append(executionOrder, "middleware2-after")
        })
    }

    middleware3 := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            executionOrder = append(executionOrder, "middleware3-before")
            next.ServeHTTP(w, r)
            executionOrder = append(executionOrder, "middleware3-after")
        })
    }

    // Add middleware in order
    stack.Use(middleware1).Use(middleware2).Use(middleware3)

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        executionOrder = append(executionOrder, "handler")
        w.WriteHeader(http.StatusOK)
    })

    finalHandler := stack.Apply(handler)

    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()

    finalHandler.ServeHTTP(w, req)

    // Verify execution order: last registered middleware executes first
    expectedOrder := []string{
        "middleware3-before", // Last registered, outermost
        "middleware2-before",
        "middleware1-before", // First registered, innermost
        "handler",
        "middleware1-after",  // First registered, innermost
        "middleware2-after",
        "middleware3-after",  // Last registered, outermost
    }

    if len(executionOrder) != len(expectedOrder) {
        t.Fatalf("Expected %d execution steps, got %d", len(expectedOrder), len(executionOrder))
    }

    for i, expected := range expectedOrder {
        if executionOrder[i] != expected {
            t.Errorf("Execution order[%d]: expected %q, got %q", i, expected, executionOrder[i])
        }
    }
}

func TestMiddlewareStack_UseIf(t *testing.T) {
    tests := []struct {
        name           string
        condition      bool
        expectedCount  int
        expectHeader   bool
    }{
        {
            name:          "condition true",
            condition:     true,
            expectedCount: 1,
            expectHeader:  true,
        },
        {
            name:          "condition false",
            condition:     false,
            expectedCount: 0,
            expectHeader:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            stack := NewMiddlewareStack()

            middleware := func(next http.Handler) http.Handler {
                return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                    w.Header().Set("X-Conditional", "applied")
                    next.ServeHTTP(w, r)
                })
            }

            // Test conditional addition
            result := stack.UseIf(tt.condition, middleware)

            // Verify chaining
            if result != stack {
                t.Error("UseIf() should return same instance for chaining")
            }

            // Verify count
            if stack.Count() != tt.expectedCount {
                t.Errorf("Expected %d middleware, got %d", tt.expectedCount, stack.Count())
            }

            // Test application
            handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
            })

            finalHandler := stack.Apply(handler)
            req := httptest.NewRequest("GET", "/test", nil)
            w := httptest.NewRecorder()

            finalHandler.ServeHTTP(w, req)

            // Verify header presence
            hasHeader := w.Header().Get("X-Conditional") == "applied"
            if hasHeader != tt.expectHeader {
                t.Errorf("Expected header present: %v, got: %v", tt.expectHeader, hasHeader)
            }
        })
    }
}
```

### HTTP Server Lifecycle Testing Patterns
```go
func TestServer_Shutdown(t *testing.T) {
    tests := []struct {
        name        string
        config      ServerConfig
        startServer bool
        wantErr     bool
    }{
        {
            name: "successful shutdown of running server",
            config: ServerConfig{
                Host:            "localhost",
                Port:            "0",
                ShutdownTimeout: 5 * time.Second,
            },
            startServer: true,
            wantErr:     false,
        },
        {
            name: "shutdown of already stopped server",
            config: ServerConfig{
                Host:            "localhost",
                Port:            "0",
                ShutdownTimeout: 5 * time.Second,
            },
            startServer: false,
            wantErr:     false,
        },
        {
            name: "shutdown with short but reasonable timeout",
            config: ServerConfig{
                Host:            "localhost",
                Port:            "0",
                ShutdownTimeout: 1 * time.Second,
            },
            startServer: true,
            wantErr:     false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            server, err := NewServer(tt.config)
            if err != nil {
                t.Fatalf("NewServer() unexpected error: %v", err)
            }

            // Start server if requested
            if tt.startServer {
                err = server.Start()
                if err != nil {
                    t.Fatalf("Start() unexpected error: %v", err)
                }

                // Verify server is running
                if !server.IsStarted() {
                    t.Error("Server should be started before shutdown test")
                }
            }

            // Test shutdown
            err = server.Shutdown()

            if tt.wantErr {
                if err == nil {
                    t.Error("Shutdown() expected error, got nil")
                }
            } else {
                if err != nil {
                    t.Errorf("Shutdown() unexpected error: %v", err)
                }
            }

            // Verify server state after shutdown
            if server.IsStarted() {
                t.Error("Server should not be started after Shutdown()")
            }

            // Verify shutdown is idempotent
            err2 := server.Shutdown()
            if err2 != nil {
                t.Errorf("Second Shutdown() should be idempotent, got error: %v", err2)
            }
        })
    }
}
```

### Thread Safety Testing Patterns
```go
func TestServer_Shutdown_Concurrent(t *testing.T) {
    config := ServerConfig{
        Host:            "localhost",
        Port:            "0",
        ShutdownTimeout: 2 * time.Second,
    }

    server, err := NewServer(config)
    if err != nil {
        t.Fatalf("NewServer() unexpected error: %v", err)
    }

    // Start the server
    err = server.Start()
    if err != nil {
        t.Fatalf("Start() unexpected error: %v", err)
    }

    // Test concurrent shutdown calls
    var wg sync.WaitGroup
    errors := make(chan error, 3)

    // Launch multiple concurrent shutdown calls
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            errors <- server.Shutdown()
        }()
    }

    wg.Wait()
    close(errors)

    // Collect results
    var errorCount int
    for err := range errors {
        if err != nil {
            errorCount++
            t.Logf("Shutdown error: %v", err)
        }
    }

    // All shutdown calls should succeed (idempotent)
    if errorCount > 0 {
        t.Errorf("Expected all concurrent shutdowns to succeed, got %d errors", errorCount)
    }

    // Verify final state
    if server.IsStarted() {
        t.Error("Server should not be started after concurrent shutdown")
    }
}
```

### Complete Lifecycle Testing Pattern
```go
func TestServer_CompleteLifecycle(t *testing.T) {
    config := ServerConfig{
        Host:            "localhost",
        Port:            "0",
        ShutdownTimeout: 3 * time.Second,
    }

    server, err := NewServer(config)
    if err != nil {
        t.Fatalf("NewServer() unexpected error: %v", err)
    }

    // Test initial state
    if server.IsStarted() {
        t.Error("New server should not be started")
    }

    // Test start
    err = server.Start()
    if err != nil {
        t.Fatalf("Start() unexpected error: %v", err)
    }

    if !server.IsStarted() {
        t.Error("Server should be started after Start()")
    }

    // Brief pause to ensure server is fully operational
    time.Sleep(50 * time.Millisecond)

    // Test shutdown
    err = server.Shutdown()
    if err != nil {
        t.Errorf("Shutdown() unexpected error: %v", err)
    }

    if server.IsStarted() {
        t.Error("Server should not be started after Shutdown()")
    }

    // Test idempotent shutdown
    err = server.Shutdown()
    if err != nil {
        t.Errorf("Second Shutdown() should be idempotent: %v", err)
    }

    // Test that start after shutdown should fail
    err = server.Start()
    if err == nil {
        t.Error("Start() after Shutdown() should fail")
        // Clean up if it unexpectedly succeeded
        server.Shutdown()
    }

    // Verify the error message
    if err != nil && !strings.Contains(err.Error(), "cannot start server after shutdown") {
        t.Errorf("Expected 'cannot start server after shutdown' error, got: %v", err)
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

            // Verify middleware stack is initialized
            if server.Middleware() == nil {
                t.Error("Server middleware stack should be initialized")
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
// for lifecycle management with proper shutdown coordination.
//
// The server includes an initialized middleware stack ready for middleware
// registration and a root handler management system for request processing.
//
// Parameters:
//   - config: ServerConfig containing host, port, and timeout configuration
//
// Returns:
//   - *Server: Configured server instance ready for middleware and handler setup
//   - error: Validation error if configuration is invalid
//
// Security: Applies secure timeout defaults and validates all configuration parameters.
// The server is prepared but not started; call Start() to begin accepting connections.
//
// Middleware Pattern: The server includes a MiddlewareStack accessible via Middleware()
// method for adding request processing middleware before starting the server.
func NewServer(config ServerConfig) (*Server, error) {
    // ... implementation
}
```

### Package Documentation Pattern
```go
// Package server provides HTTP server infrastructure for Cipher Hub.
//
// This package implements the core HTTP server functionality with structured
// configuration management, graceful shutdown capabilities, middleware support,
// and security-first design principles.
//
// Key Features:
//   - Production-ready HTTP server lifecycle with graceful shutdown
//   - Industry-standard middleware infrastructure with conditional support
//   - Comprehensive input validation with injection prevention
//   - Context-based shutdown coordination and timeout management
//   - Thread-safe concurrent access patterns
//
// Middleware Infrastructure:
//   - Standard middleware signature: func(http.Handler) http.Handler
//   - Enhanced middleware stack with conditional application (UseIf)
//   - Method chaining for fluent API design
//   - Correct execution order (last registered becomes outermost)
//
// Container Integration:
//   - SIGINT and SIGTERM signal handling for orchestration platforms
//   - Health check endpoints for container health monitoring
//   - Environment variable configuration support
//   - Proper resource cleanup on shutdown
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
//	
//	// Configure middleware
//	server.Middleware().
//		Use(RequestLoggingMiddleware()).
//		UseIf(config.EnableCORS, CORSMiddleware(config.CORSOrigins))
//	
//	// Set handler and start server
//	server.SetHandler(myHandler)
//	go server.Start()
//	
//	// Handle graceful shutdown
//	sigChan := make(chan os.Signal, 1)
//	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
//	<-sigChan
//	server.Shutdown()
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
├── cmd/cipher-hub/           # Application entry point with signal handling
├── internal/
│   ├── models/              # Core data models (Phase 1 ✅)
│   ├── storage/             # Storage interface (Phase 1 ✅)  
│   ├── server/              # HTTP server infrastructure (Phase 2 → Target 2.1 ✅)
│   │   ├── server.go        # Complete HTTP server with graceful shutdown (✅)
│   │   ├── server_test.go   # Comprehensive security and lifecycle testing (✅)
│   │   ├── middleware.go    # Complete middleware infrastructure (✅)
│   │   └── middleware_test.go # Comprehensive middleware testing (✅)
│   └── handlers/            # HTTP request handlers (Phase 2 → Target 2.1 📋)
├── checkpoint.md           # Development progress and next steps
├── go.mod                  # Go module definition
├── readme.md               # Project homepage and documentation
├── review.md               # Latest code review findings and quality status
├── roadmap.md              # Development roadmap with granular steps
├── spec.md                 # Technical specification
└── style-guide.md          # Implementation standards (primary reference)
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
- [ ] Thread-safe operations preventing race conditions
- [ ] Graceful shutdown with proper resource cleanup
- [ ] Signal handling without race conditions
- [ ] Middleware security patterns with nil handler protection

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
- [ ] Support port "0" for OS dynamic assignment while maintaining security

### HTTP Server Security Requirements
- [ ] Thread-safe concurrent access patterns implemented
- [ ] Proper resource cleanup on all failure scenarios
- [ ] Idempotent operations preventing double-start issues
- [ ] Channel-based coordination avoiding race conditions
- [ ] Secure error handling without configuration leakage
- [ ] Proper HTTP server lifecycle management with cleanup
- [ ] Graceful shutdown with timeout enforcement
- [ ] Signal handling with proper coordination and cleanup

### Middleware Security Requirements
- [ ] Industry-standard middleware signature implemented
- [ ] Nil handler protection preventing runtime errors
- [ ] Correct middleware execution order ensuring standard composition
- [ ] Conditional middleware deployment maintains security posture
- [ ] Method chaining provides secure configuration patterns
- [ ] Server integration maintains thread safety during setup
- [ ] Comprehensive testing covers security scenarios and edge cases

### Shutdown Security Requirements
- [ ] In-flight requests complete within timeout bounds
- [ ] Resource cleanup handles all failure scenarios
- [ ] State consistency maintained during shutdown process
- [ ] Thread safety during concurrent shutdown attempts
- [ ] Signal handling prevents race conditions
- [ ] Context cancellation properly coordinated
- [ ] Error handling includes timeout context information

---

*Implementation Style Guide Version: 1.7*  
*Last Updated: Current Session*  
*Focus: Security-first implementation patterns with complete HTTP server lifecycle, middleware infrastructure, and graceful shutdown for Step 2.1.2.1 completion*