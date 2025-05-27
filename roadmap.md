# Cipher Hub - Development Roadmap

## Current Status: Phase 2.1 HTTP Server Infrastructure In Progress 🔄

**Cipher Hub** development follows a structured, session-based approach with granular 20-30 minute implementation steps, comprehensive testing, and quality gates at each phase transition.

---

## Phase 1: Foundation Architecture ✅ **COMPLETE**

### Completed Infrastructure
- [x] **Project Structure & Organization** - Proper Go project layout with best practices
- [x] **Core Data Models** - ServiceRegistration, Participant, CryptoKey with full validation
- [x] **Error Handling System** - Graceful error handling with proper Go idioms
- [x] **Comprehensive Test Coverage** - Unit tests for all components with table-driven patterns
- [x] **Storage Interface Design** - Abstract storage layer ready for implementation
- [x] **Security-First Design** - Key material protection, validation everywhere

### Key Design Decisions Established
- **Participant Type Abstraction**: Metadata-driven approach using `Metadata["type"]` pattern
- **Algorithm Enum Strategy**: Strict enum aligned with implemented capabilities (`AlgorithmAES256`)
- **Error Handling Pattern**: All constructors return `(Type, error)` with proper error wrapping
- **Time Field Strategy**: Required fields use `time.Time`, optional fields use `*time.Time`
- **File Organization**: One primary type per file with co-located tests

---

## Phase 2: HTTP Server Infrastructure 🔄 **CURRENT PHASE**

### Target 2.1: Basic Server Setup ⏳ **CURRENT TARGET**

#### **Task 2.1.1: HTTP Server Creation** (`internal/server/server.go`) **CURRENT TASK**
- [x] **Step 2.1.1.1**: Create basic HTTP server struct with configuration fields ✅ **COMPLETE**
  - ✅ Defined `ServerConfig` struct with host, port, and timeout fields
  - ✅ Added structured constructor `NewServer(config ServerConfig) (*Server, error)`
  - ✅ Included context for graceful shutdown with configurable timeout
  - ✅ Implemented comprehensive validation and security-first design
  - ✅ Applied defaults pattern with `ApplyDefaults()` method
- [x] **Step 2.1.1.2**: Implement basic HTTP listener setup ✅ **COMPLETE**
  - ✅ Added `Start()` method that creates `http.Server` instance
  - ✅ Configured HTTP server with validated timeouts from ServerConfig
  - ✅ Added port binding and listener creation with comprehensive error handling
  - ✅ Integrated with shutdown context for lifecycle management
  - ✅ Implemented thread safety with `sync.RWMutex` for concurrent access
  - ✅ Enhanced port validation to support port "0" for dynamic assignment
  - ✅ Added channel-based readiness signaling for reliable startup coordination
- [ ] **Step 2.1.1.3**: Add graceful shutdown mechanism ⏳ **IMMEDIATE NEXT**
  - Enhance `Shutdown()` method to coordinate with `http.Server.Shutdown()`
  - Implement signal handling for SIGINT and SIGTERM
  - Ensure in-flight requests complete before shutdown
  - Use configured shutdown timeout from ServerConfig
  - Add proper resource cleanup and state management
- [x] **Step 2.1.1.4**: Create server configuration structure ✅ **COMPLETE** *(Integrated into 2.1.1.1)*
  - ✅ Defined `ServerConfig` struct for server settings with comprehensive fields
  - ✅ Added validation for configuration values with security bounds
  - ✅ Included environment variable loading foundation
  - ✅ Implemented default value application and timeout bounds checking
- [x] **Step 2.1.1.5**: Add basic server tests ✅ **COMPLETE**
  - ✅ Comprehensive test coverage for server struct and configuration
  - ✅ Test server configuration validation with security edge cases
  - ✅ Test constructor behavior and error handling
  - ✅ Validation of accessor methods and shutdown context management
  - ✅ HTTP server lifecycle testing with Start() method
  - ✅ Thread safety testing with concurrent access scenarios
  - ✅ Port binding and listener creation error handling tests

**Architectural Decisions Made:**
- **Structured Configuration**: `ServerConfig` pattern establishes foundation for environment loading
- **Security-First Validation**: Comprehensive input validation with injection attack prevention
- **Timeout Management**: Configurable timeouts with secure defaults and reasonable bounds
- **Context Integration**: Shutdown context with timeout for graceful lifecycle management
- **Thread Safety**: Production-ready concurrent access patterns with `sync.RWMutex`
- **Dynamic Port Support**: Enhanced port validation supporting both explicit and OS-assigned ports

#### **Task 2.1.2: Middleware Infrastructure** (`internal/server/middleware.go`)
- [ ] **Step 2.1.2.1**: Create middleware function signature pattern
  - Define `Middleware` type as `func(http.Handler) http.Handler`
  - Create enhanced middleware stack with conditional support
  - Add middleware application pattern with proper chaining
- [ ] **Step 2.1.2.2**: Implement request logging middleware
  - Add request ID generation using `crypto/rand`
  - Log HTTP method, path, and response status with structured logging
  - Include request duration timing and correlation IDs
- [ ] **Step 2.1.2.3**: Add CORS handling middleware
  - Support environment-configurable CORS origins (not hard-coded)
  - Add configurable allowed origins based on deployment environment
  - Handle preflight OPTIONS requests with proper headers
- [ ] **Step 2.1.2.4**: Create error response formatting middleware
  - Standardize JSON error response format with request tracing
  - Add error code mapping from internal errors with security consciousness
  - Ensure no sensitive data leaks in error responses
- [ ] **Step 2.1.2.5**: Implement security headers middleware
  - Add comprehensive security headers (HSTS, CSP, X-Frame-Options)
  - Include `X-Content-Type-Options: nosniff` and XSS protection
  - Apply conditional HSTS only for HTTPS deployments
- [ ] **Step 2.1.2.6**: Add middleware tests
  - Test middleware chaining and conditional application
  - Test request ID generation and propagation with security validation
  - Test error response formatting and CORS header setting

#### **Task 2.1.3: Health Check System** (`internal/handlers/health.go`)
- [ ] **Step 2.1.3.1**: Create basic health check handler structure
  - Define `HealthHandler` struct with dependency health checkers
  - Add constructor `NewHealthHandler() *HealthHandler` with interface support
  - Set up health check interface pattern for extensibility
- [ ] **Step 2.1.3.2**: Implement liveness endpoint
  - Add `/health/live` endpoint that always returns 200 OK
  - Return simple JSON response `{"status": "alive", "timestamp": "..."}`
  - Include basic server operational status
- [ ] **Step 2.1.3.3**: Implement readiness endpoint with dependency checks
  - Add `/health/ready` endpoint with actual dependency validation
  - Return detailed readiness status with individual component checks
  - Use JSON format `{"status": "ready", "checks": {...}, "timestamp": "..."}`
- [ ] **Step 2.1.3.4**: Add health check response models
  - Create `HealthStatus` struct for consistent JSON responses
  - Define `CheckResult` interface for extensible health checks
  - Add proper JSON tags and validation with error handling
- [ ] **Step 2.1.3.5**: Create health check tests
  - Test liveness endpoint returns 200 with proper JSON format
  - Test readiness endpoint dependency checking behavior
  - Test health check interface compliance and error scenarios

#### **Task 2.1.4: Handler Framework** (`internal/handlers/handlers.go`)
- [ ] **Step 2.1.4.1**: Create base handler utilities
  - Add `writeJSON` utility function for consistent responses
  - Create `readJSON` utility for request parsing with size limits
  - Include proper content-type handling and validation
- [ ] **Step 2.1.4.2**: Implement error response utilities
  - Add `writeError` function with structured error responses and request tracing
  - Create `ValidationError` type for input errors with detailed field information
  - Ensure consistent error JSON format following established patterns
- [ ] **Step 2.1.4.3**: Add request parsing helpers
  - Create URL parameter extraction utilities with validation
  - Add query parameter parsing helpers with type conversion
  - Include request body size limiting and malformed JSON handling
- [ ] **Step 2.1.4.4**: Implement response header utilities
  - Add common header setting functions for security and caching
  - Include cache control helpers for different resource types
  - Create security header application functions
- [ ] **Step 2.1.4.5**: Create handler framework tests
  - Test JSON response utilities with various data types
  - Test error response formatting and request tracing integration
  - Test request parsing edge cases and security validation

**Key Architectural Decisions Required:**
- **HTTP Router**: Standard library `net/http` with Go 1.22+ enhanced routing (recommended)
- **Configuration**: Environment variables with structured validation (foundation established)
- **Logging**: Standard library `log/slog` with structured JSON output
- **Middleware Pattern**: Enhanced function chaining with conditional support

### Target 2.2: API Foundation

#### **Task 2.2.1: Service Registration Endpoints**
- [ ] **Step 2.2.1.1**: Create service registration handler structure
- [ ] **Step 2.2.1.2**: Implement GET /services endpoint
- [ ] **Step 2.2.1.3**: Implement GET /services/{id} endpoint
- [ ] **Step 2.2.1.4**: Implement POST /services endpoint
- [ ] **Step 2.2.1.5**: Implement PUT /services/{id} endpoint
- [ ] **Step 2.2.1.6**: Implement DELETE /services/{id} endpoint
- [ ] **Step 2.2.1.7**: Add service endpoint tests

#### **Task 2.2.2: Participant Management Endpoints**
- [ ] **Step 2.2.2.1**: Create participant handler structure
- [ ] **Step 2.2.2.2**: Implement GET /services/{id}/participants
- [ ] **Step 2.2.2.3**: Implement POST /services/{id}/participants
- [ ] **Step 2.2.2.4**: Implement GET /services/{id}/participants/{pid}
- [ ] **Step 2.2.2.5**: Implement PUT /services/{id}/participants/{pid}
- [ ] **Step 2.2.2.6**: Implement DELETE /services/{id}/participants/{pid}
- [ ] **Step 2.2.2.7**: Add participant endpoint tests

#### **Task 2.2.3: Request/Response Standards**
- [ ] **Step 2.2.3.1**: Define API response envelope structure
- [ ] **Step 2.2.3.2**: Implement pagination response format
- [ ] **Step 2.2.3.3**: Create input validation framework
- [ ] **Step 2.2.3.4**: Implement API versioning strategy
- [ ] **Step 2.2.3.5**: Add response standard tests

#### **Task 2.2.4: Enhanced Configuration System**
- [ ] **Step 2.2.4.1**: Extend ServerConfig with application settings
- [ ] **Step 2.2.4.2**: Implement environment variable loading
- [ ] **Step 2.2.4.3**: Add CORS origins environment configuration
- [ ] **Step 2.2.4.4**: Implement configuration file support
- [ ] **Step 2.2.4.5**: Create enhanced configuration tests

### Target 2.3: Initial Integration

#### **Task 2.3.1: In-Memory Storage Implementation**
- [ ] **Step 2.3.1.1**: Create in-memory storage struct
- [ ] **Step 2.3.1.2**: Implement service storage operations
- [ ] **Step 2.3.1.3**: Implement participant storage operations
- [ ] **Step 2.3.1.4**: Add storage operation tests
- [ ] **Step 2.3.1.5**: Implement storage interface compliance

#### **Task 2.3.2: API Integration Testing**
- [ ] **Step 2.3.2.1**: Create integration test framework
- [ ] **Step 2.3.2.2**: Implement service API integration tests
- [ ] **Step 2.3.2.3**: Add participant API integration tests
- [ ] **Step 2.3.2.4**: Create end-to-end workflow tests
- [ ] **Step 2.3.2.5**: Add performance baseline tests

---

## Phase 3: Security Foundation

### Target 3.1: Authentication System
#### **Task 3.1.1: API Key Authentication**
- [ ] **Step 3.1.1.1**: Create API key data structure
- [ ] **Step 3.1.1.2**: Implement API key hashing
- [ ] **Step 3.1.1.3**: Create API key storage interface
- [ ] **Step 3.1.1.4**: Implement API key validation middleware
- [ ] **Step 3.1.1.5**: Add API key management endpoints
- [ ] **Step 3.1.1.6**: Create API key tests

#### **Task 3.1.2: Request Authentication**
- [ ] **Step 3.1.2.1**: Implement bearer token parsing
- [ ] **Step 3.1.2.2**: Add request context authentication
- [ ] **Step 3.1.2.3**: Create authentication error handling
- [ ] **Step 3.1.2.4**: Implement authentication bypass for health checks
- [ ] **Step 3.1.2.5**: Add authentication integration tests

### Target 3.2: Authorization Framework
#### **Task 3.2.1: Role-Based Access Control (RBAC)**
- [ ] **Step 3.2.1.1**: Define permission and role structures
- [ ] **Step 3.2.1.2**: Implement permission checking functions
- [ ] **Step 3.2.1.3**: Create authorization middleware
- [ ] **Step 3.2.1.4**: Add role management functionality
- [ ] **Step 3.2.1.5**: Implement resource-level authorization
- [ ] **Step 3.2.1.6**: Create authorization tests

### Target 3.3: Secure Key Storage
#### **Task 3.3.1: Encryption at Rest**
- [ ] **Step 3.3.1.1**: Create master key management
- [ ] **Step 3.3.1.2**: Implement key material encryption
- [ ] **Step 3.3.1.3**: Add secure memory handling
- [ ] **Step 3.3.1.4**: Create key serialization protection
- [ ] **Step 3.3.1.5**: Implement audit trail encryption

---

## Phase 4: Key Lifecycle Management

### Target 4.1: Key Generation Engine
#### **Task 4.1.1: AES-256 Key Generation**
- [ ] **Step 4.1.1.1**: Implement secure random key generation
- [ ] **Step 4.1.1.2**: Create key generation request handling
- [ ] **Step 4.1.1.3**: Implement key quality validation
- [ ] **Step 4.1.1.4**: Add key generation audit logging
- [ ] **Step 4.1.1.5**: Create key generation tests

#### **Task 4.1.2: Key Metadata Management**
- [ ] **Step 4.1.2.1**: Define key metadata structure
- [ ] **Step 4.1.2.2**: Implement key versioning system
- [ ] **Step 4.1.2.3**: Add key status lifecycle management
- [ ] **Step 4.1.2.4**: Create usage statistics tracking
- [ ] **Step 4.1.2.5**: Implement metadata persistence

### Target 4.2: Key Distribution APIs
#### **Task 4.2.1: Secure Key Retrieval**
- [ ] **Step 4.2.1.1**: Create key retrieval endpoints
- [ ] **Step 4.2.1.2**: Implement key access authorization
- [ ] **Step 4.2.1.3**: Add key retrieval audit logging
- [ ] **Step 4.2.1.4**: Implement secure key response format
- [ ] **Step 4.2.1.5**: Create key retrieval tests

#### **Task 4.2.2: Key Versioning System**
- [ ] **Step 4.2.2.1**: Implement version-aware key retrieval
- [ ] **Step 4.2.2.2**: Add backward compatibility handling
- [ ] **Step 4.2.2.3**: Create version-specific operations
- [ ] **Step 4.2.2.4**: Implement version cleanup policies

---

## Phase 5: Production Readiness

### Target 5.1: Persistent Storage Backends
#### **Task 5.1.1: Database Integration**
- [ ] **Step 5.1.1.1**: Create database connection management
- [ ] **Step 5.1.1.2**: Define database schema
- [ ] **Step 5.1.1.3**: Implement database migrations
- [ ] **Step 5.1.1.4**: Add database storage implementation
- [ ] **Step 5.1.1.5**: Create database tests

---

## Success Metrics

### Phase 2 Success Criteria
- HTTP server handles basic CRUD operations for all core entities
- Health checks integrate with container orchestration platforms
- Consistent API patterns established across all endpoints
- Integration tests cover full request/response cycles
- Configuration system supports environment-based deployment

### Phase 3 Success Criteria
- Secure authentication prevents unauthorized access with comprehensive testing
- Authorization controls enforce proper permissions with audit trails
- Key material never exposed in logs or serialization under any circumstances
- Audit trails capture all security-relevant events with integrity protection

### Long-term Success Criteria
- Handle 10,000+ concurrent key operations with sub-100ms response times
- 99.9% uptime in production deployment with comprehensive monitoring
- Zero key material exposure incidents with continuous security validation
- Complete audit trails for compliance requirements with tamper detection
- Seamless key rotation and lifecycle management with zero downtime

---

## Technical Architecture Decisions

### Established in Phase 2.1
- **HTTP Router**: Standard library `net/http` with Go 1.22+ enhanced routing patterns
- **Configuration**: Structured configuration with environment variable support and validation
- **Constructor Pattern**: `NewServer(config ServerConfig) (*Server, error)` with comprehensive validation
- **Context Management**: Typed context keys with shutdown timeout configuration
- **Error Handling**: Consistent error prefixes with structured error responses
- **Security Validation**: Comprehensive input validation with injection attack prevention
- **Thread Safety**: Production-ready concurrent access patterns with `sync.RWMutex`
- **Port Handling**: Enhanced port validation supporting both explicit and dynamic assignment

### Middleware Architecture
- **Pattern**: Enhanced function chaining with conditional middleware support
- **CORS Configuration**: Environment-configurable origins (not hard-coded)
- **Security Headers**: Comprehensive security header implementation with HTTPS detection
- **Request Tracing**: Correlation IDs throughout request lifecycle with audit integration

### Testing Strategy
- **Comprehensive Coverage**: >90% test coverage with security-focused edge case testing
- **Integration Testing**: End-to-end API testing with real HTTP server instances
- **Security Testing**: Input validation, injection prevention, and authentication/authorization testing
- **Performance Baseline**: Establish performance metrics and monitoring for optimization
- **Concurrency Testing**: Thread safety validation with concurrent access scenarios

---

## Session-Based Development Guidelines

### Step Granularity Standards
- **Time Requirement**: Each numbered step completable in 20-30 minutes
- **Incremental Progress**: Steps build upon previous completed work
- **Testable Units**: Each step includes comprehensive test requirements
- **Documentation**: Step completion includes code documentation and validation

### Quality Gates
- **Pre-Implementation**: Step guides validated before execution
- **During Implementation**: Interactive collaboration for optimization and problem-solving
- **Post-Implementation**: Comprehensive review and documentation synchronization

### Milestone Tracking
- Each numbered step requires completion validation and verification
- Regular progress checkpoints at sub-phase completion with documentation updates
- Technical decision documentation for major choices with rationale
- Architecture decision records (ADRs) for significant changes affecting future development

---

*Roadmap Status: Updated to reflect Step 2.1.1.2 completion and current progress*  
*Current Development Status: Step 2.1.1.3 - Implement graceful shutdown mechanism*  
*Architecture Foundation: HTTP server lifecycle management with thread safety complete*  
*Next Focus: Complete graceful shutdown to finish HTTP server infrastructure*