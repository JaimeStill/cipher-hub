# Cipher Hub - Refined Development Roadmap

## Current Status: Phase 2.1 HTTP Server Infrastructure In Progress 🔄

**Cipher Hub** is a Go-based key management service designed as a centralized security layer for cryptographic operations across distributed systems, acting as a sidecar component for complete key lifecycle management.

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
- [ ] **Step 2.1.1.2**: Implement basic HTTP listener setup ⏳ **IMMEDIATE NEXT**
  - Add `Start()` method that creates `http.Server` instance
  - Configure HTTP server with validated timeouts from ServerConfig
  - Add port binding and listener creation with error handling
  - Integrate with shutdown context for lifecycle management
- [ ] **Step 2.1.1.3**: Add graceful shutdown mechanism
  - Implement `Shutdown(ctx context.Context) error` method
  - Handle `os.Signal` for SIGINT and SIGTERM
  - Ensure in-flight requests complete before shutdown
  - Use configured shutdown timeout from ServerConfig
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

**Architectural Decisions Made:**
- **Structured Configuration**: `ServerConfig` pattern establishes foundation for environment loading
- **Security-First Validation**: Comprehensive input validation with injection attack prevention
- **Timeout Management**: Configurable timeouts with secure defaults and reasonable bounds
- **Context Integration**: Shutdown context with timeout for graceful lifecycle management

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
  - Define `ServiceHandler` struct with storage dependency injection
  - Add constructor with storage interface and validation
  - Set up route method signatures following established patterns
- [ ] **Step 2.2.1.2**: Implement GET /services endpoint
  - Add list services functionality with pagination support
  - Include pagination parameters (limit, offset) with validation
  - Return JSON array of service summaries with proper headers
- [ ] **Step 2.2.1.3**: Implement GET /services/{id} endpoint
  - Add single service retrieval by ID with validation
  - Handle service not found cases with proper error responses
  - Return full service details with participants using established JSON patterns
- [ ] **Step 2.2.1.4**: Implement POST /services endpoint
  - Add service creation from JSON request with comprehensive validation
  - Validate required fields (name, description) and sanitize input
  - Return created service with 201 status and proper location headers
- [ ] **Step 2.2.1.5**: Implement PUT /services/{id} endpoint
  - Add service update functionality with validation
  - Validate update request format and handle partial updates
  - Include proper timestamp updating and audit logging
- [ ] **Step 2.2.1.6**: Implement DELETE /services/{id} endpoint
  - Add service deletion with safety checks for dependent resources
  - Ensure no active participants exist before deletion
  - Return appropriate status codes with audit logging
- [ ] **Step 2.2.1.7**: Add service endpoint tests
  - Test CRUD operations end-to-end with integration testing
  - Test error conditions and edge cases including security scenarios
  - Test input validation behavior and error response formatting

#### **Task 2.2.2: Participant Management Endpoints**
- [ ] **Step 2.2.2.1**: Create participant handler structure
  - Define `ParticipantHandler` struct with storage and validation
  - Add service-scoped participant operations following REST patterns
  - Set up nested route patterns for service relationships
- [ ] **Step 2.2.2.2**: Implement GET /services/{id}/participants
  - List participants within a service with filtering support
  - Add participant type filtering using metadata patterns
  - Include pagination support following established patterns
- [ ] **Step 2.2.2.3**: Implement POST /services/{id}/participants
  - Add participant creation within service context
  - Validate participant data and flexible type metadata
  - Return created participant details with proper relationships
- [ ] **Step 2.2.2.4**: Implement GET /services/{id}/participants/{pid}
  - Retrieve single participant by ID with relationship validation
  - Validate service and participant relationship integrity
  - Return full participant details with metadata
- [ ] **Step 2.2.2.5**: Implement PUT /services/{id}/participants/{pid}
  - Update participant information with validation
  - Validate metadata updates using flexible metadata patterns
  - Handle participant type changes through metadata
- [ ] **Step 2.2.2.6**: Implement DELETE /services/{id}/participants/{pid}
  - Remove participant from service with dependency checking
  - Check for dependent key relationships and prevent orphaned keys
  - Return appropriate confirmation with audit logging
- [ ] **Step 2.2.2.7**: Add participant endpoint tests
  - Test nested resource operations with service relationships
  - Test service-participant relationship validation and integrity
  - Test metadata handling and flexible type management

#### **Task 2.2.3: Request/Response Standards**
- [ ] **Step 2.2.3.1**: Define API response envelope structure
  - Create consistent response wrapper format with metadata
  - Include data, metadata, error, and tracing fields
  - Add response timestamp and request ID for correlation
- [ ] **Step 2.2.3.2**: Implement pagination response format
  - Add pagination metadata structure with total counts
  - Include total count, limit, offset, and navigation fields
  - Create pagination link generation for hypermedia support
- [ ] **Step 2.2.3.3**: Create input validation framework
  - Define validation error response format with detailed field errors
  - Add field-level validation error details with user-friendly messages
  - Include validation rule descriptions and suggestions
- [ ] **Step 2.2.3.4**: Implement API versioning strategy
  - Add version header handling for API evolution
  - Create version-specific route prefixes and content negotiation
  - Plan backward compatibility approach for future versions
- [ ] **Step 2.2.3.5**: Add response standard tests
  - Test consistent response formatting across all endpoints
  - Test pagination metadata accuracy and link generation
  - Test validation error responses and message quality

#### **Task 2.2.4: Enhanced Configuration System**
- [ ] **Step 2.2.4.1**: Extend ServerConfig with application settings
  - Add database connection, logging, and security configuration
  - Include environment-specific defaults and validation
  - Extend validation framework to cover new configuration areas
- [ ] **Step 2.2.4.2**: Implement environment variable loading
  - Add `LoadFromEnv()` configuration method with comprehensive parsing
  - Support standard environment variable patterns (`CIPHER_HUB_*`)
  - Include fallback to default values and environment validation
- [ ] **Step 2.2.4.3**: Add CORS origins environment configuration
  - Implement environment-configurable CORS origins (per earlier decision)
  - Support comma-separated origin lists in environment variables
  - Add validation for origin format and security
- [ ] **Step 2.2.4.4**: Implement configuration file support
  - Add YAML/JSON configuration file parsing capabilities
  - Support configuration file path via environment variable
  - Merge file and environment configurations with precedence rules
- [ ] **Step 2.2.4.5**: Create enhanced configuration tests
  - Test environment variable loading and parsing
  - Test validation error conditions and edge cases
  - Test configuration merging logic and precedence handling

### Target 2.3: Initial Integration

#### **Task 2.3.1: In-Memory Storage Implementation**
- [ ] **Step 2.3.1.1**: Create in-memory storage struct
  - Define `MemoryStorage` struct with sync.RWMutex for thread safety
  - Add maps for services, participants, and keys with proper indexing
  - Implement storage interface method signatures with full compliance
- [ ] **Step 2.3.1.2**: Implement service storage operations
  - Add service CRUD operations with thread safety guarantees
  - Include service ID generation and comprehensive validation
  - Handle concurrent access with proper locking strategies
- [ ] **Step 2.3.1.3**: Implement participant storage operations
  - Add participant operations within service context and relationships
  - Maintain service-participant relationships with referential integrity
  - Include participant lookup and filtering with efficient indexing
- [ ] **Step 2.3.1.4**: Add storage operation tests
  - Test concurrent access patterns and thread safety
  - Test data consistency under load with race condition detection
  - Test error conditions and edge cases including boundary conditions
- [ ] **Step 2.3.1.5**: Implement storage interface compliance
  - Verify all interface methods implemented with proper signatures
  - Add context cancellation support throughout storage operations
  - Include proper error handling patterns following established conventions

#### **Task 2.3.2: API Integration Testing**
- [ ] **Step 2.3.2.1**: Create integration test framework
  - Set up test server with in-memory storage and proper lifecycle
  - Add test utilities for HTTP requests and response validation
  - Include test data setup and teardown with proper isolation
- [ ] **Step 2.3.2.2**: Implement service API integration tests
  - Test complete service lifecycle via API endpoints
  - Include error scenario testing and edge case handling
  - Verify response format compliance with established standards
- [ ] **Step 2.3.2.3**: Add participant API integration tests
  - Test nested resource operations via API with relationship validation
  - Verify service-participant relationship handling and integrity
  - Test concurrent participant operations and consistency
- [ ] **Step 2.3.2.4**: Create end-to-end workflow tests
  - Test complete service setup workflows across multiple endpoints
  - Include multi-step operations with transaction-like behavior
  - Verify data consistency across operations and proper state management
- [ ] **Step 2.3.2.5**: Add performance baseline tests
  - Measure response times for basic operations and establish baselines
  - Test concurrent request handling and resource utilization
  - Establish performance baseline metrics for future optimization

---

## Phase 3: Security Foundation

### Target 3.1: Authentication System
#### **Task 3.1.1: API Key Authentication**
- [ ] **Step 3.1.1.1**: Create API key data structure
  - Define `APIKey` struct with ID, key hash, and comprehensive metadata
  - Add key generation using `crypto/rand` with proper entropy
  - Include key expiration and status fields with lifecycle management
- [ ] **Step 3.1.1.2**: Implement API key hashing
  - Add secure key hashing using bcrypt or Argon2 with proper parameters
  - Create key comparison functions with timing attack protection
  - Include salt generation and validation with secure practices
- [ ] **Step 3.1.1.3**: Create API key storage interface
  - Define storage methods for API key operations with security focus
  - Add key lookup by hash functionality with efficient indexing
  - Include key status management and audit logging integration
- [ ] **Step 3.1.1.4**: Implement API key validation middleware
  - Add `Authorization` header parsing with proper format validation
  - Create key lookup and validation logic with rate limiting
  - Handle authentication failure responses without information leakage
- [ ] **Step 3.1.1.5**: Add API key management endpoints
  - Create key generation endpoint with proper permissions
  - Add key listing and revocation endpoints with audit logging
  - Include key status update functionality with authorization checks
- [ ] **Step 3.1.1.6**: Create API key tests
  - Test key generation and validation with security edge cases
  - Test authentication middleware and rate limiting behavior
  - Test key management operations and audit logging

#### **Task 3.1.2: Request Authentication**
- [ ] **Step 3.1.2.1**: Implement bearer token parsing
  - Add `Authorization: Bearer <token>` support with format validation
  - Create token extraction utilities with proper error handling
  - Handle malformed authorization headers securely
- [ ] **Step 3.1.2.2**: Add request context authentication
  - Store authenticated principal in request context using typed keys
  - Create context extraction utilities following established patterns
  - Include authentication status helpers and user information access
- [ ] **Step 3.1.2.3**: Create authentication error handling
  - Define authentication-specific error types with proper categorization
  - Add proper HTTP status code mapping without information disclosure
  - Ensure no sensitive data in error responses with audit logging
- [ ] **Step 3.1.2.4**: Implement authentication bypass for health checks
  - Allow unauthenticated access to health endpoints for monitoring
  - Create authentication exemption patterns with security validation
  - Test public endpoint accessibility without compromising security
- [ ] **Step 3.1.2.5**: Add authentication integration tests
  - Test authenticated request workflows end-to-end
  - Test authentication failure scenarios and proper error handling
  - Test context propagation through handlers and middleware

### Target 3.2: Authorization Framework
#### **Task 3.2.1: Role-Based Access Control (RBAC)**
- [ ] **Step 3.2.1.1**: Define permission and role structures
  - Create `Permission` enum for resource operations with extensibility
  - Define `Role` struct with permission collections and inheritance
  - Add role assignment and inheritance patterns with validation
- [ ] **Step 3.2.1.2**: Implement permission checking functions
  - Add `HasPermission(user, resource, operation)` function with caching
  - Create resource-specific permission validation with context awareness
  - Include role-based permission resolution with inheritance support
- [ ] **Step 3.2.1.3**: Create authorization middleware
  - Add permission checking to request pipeline with performance optimization
  - Include resource identification from request with proper validation
  - Handle authorization failure responses with audit logging
- [ ] **Step 3.2.1.4**: Add role management functionality
  - Create role assignment operations with proper authorization
  - Include role update and revocation with audit trails
  - Add role inheritance processing with cycle detection
- [ ] **Step 3.2.1.5**: Implement resource-level authorization
  - Add service-specific access controls with fine-grained permissions
  - Include participant-level permissions with relationship validation
  - Create key operation authorization with security boundary enforcement
- [ ] **Step 3.2.1.6**: Create authorization tests
  - Test permission checking logic with complex scenarios
  - Test authorization middleware and performance under load
  - Test role management operations and inheritance behavior

### Target 3.3: Secure Key Storage
#### **Task 3.3.1: Encryption at Rest**
- [ ] **Step 3.3.1.1**: Create master key management
  - Implement master key generation and secure storage
  - Add key derivation for data encryption keys with proper parameters
  - Include master key rotation capabilities with backwards compatibility
- [ ] **Step 3.3.1.2**: Implement key material encryption
  - Add encryption before storage persistence with authenticated encryption
  - Create decryption on key retrieval with proper error handling
  - Include encrypted field tagging and secure serialization prevention
- [ ] **Step 3.3.1.3**: Add secure memory handling
  - Implement secure memory allocation for sensitive key material
  - Add memory clearing after key usage with explicit zeroing
  - Include memory protection from dumps and swap files
- [ ] **Step 3.3.1.4**: Create key serialization protection
  - Ensure encrypted keys never serialize as plaintext under any circumstances
  - Add secure key marshaling/unmarshaling with validation
  - Include key material protection in logs and error messages
- [ ] **Step 3.3.1.5**: Implement audit trail encryption
  - Add audit log encryption capabilities with tamper detection
  - Create tamper-evident audit records with integrity verification
  - Include audit integrity verification and secure log storage

---

## Phase 4: Key Lifecycle Management

### Target 4.1: Key Generation Engine
#### **Task 4.1.1: AES-256 Key Generation**
- [ ] **Step 4.1.1.1**: Implement secure random key generation
  - Use `crypto/rand` for cryptographically secure randomness with proper error handling
  - Add AES-256 key length validation (32 bytes) with format standardization
  - Include key format standardization and entropy validation
- [ ] **Step 4.1.1.2**: Create key generation request handling
  - Add key generation API endpoint with proper authentication and authorization
  - Include algorithm specification in requests with validation
  - Validate generation parameters and security requirements
- [ ] **Step 4.1.1.3**: Implement key quality validation
  - Add entropy checking for generated keys with statistical analysis
  - Include weak key detection and rejection mechanisms
  - Create key strength verification with industry standards
- [ ] **Step 4.1.1.4**: Add key generation audit logging
  - Log all key generation events with comprehensive metadata
  - Include generation metadata and context with correlation IDs
  - Ensure no key material in audit logs with secure logging practices
- [ ] **Step 4.1.1.5**: Create key generation tests
  - Test key randomness and uniqueness with statistical validation
  - Test key format compliance and security requirements
  - Test generation error handling and edge cases

#### **Task 4.1.2: Key Metadata Management**
- [ ] **Step 4.1.2.1**: Define key metadata structure
  - Create comprehensive key metadata fields with extensibility
  - Include version, status, and usage tracking with proper indexing
  - Add creation and modification timestamps with audit integration
- [ ] **Step 4.1.2.2**: Implement key versioning system
  - Add version tracking for key updates with proper sequencing
  - Create version history maintenance with efficient storage
  - Include version-specific operations and queries
- [ ] **Step 4.1.2.3**: Add key status lifecycle management
  - Define key status states (active, deprecated, revoked) with transitions
  - Implement status transition validation with business rules
  - Include status change audit logging with authorization checks
- [ ] **Step 4.1.2.4**: Create usage statistics tracking
  - Add key usage counters and metrics with efficient aggregation
  - Include last access timestamp tracking with performance optimization
  - Create usage pattern analysis and anomaly detection
- [ ] **Step 4.1.2.5**: Implement metadata persistence
  - Add metadata storage operations with indexing and querying
  - Include metadata query capabilities with filtering and pagination
  - Create metadata update validation with concurrency control

### Target 4.2: Key Distribution APIs
#### **Task 4.2.1: Secure Key Retrieval**
- [ ] **Step 4.2.1.1**: Create key retrieval endpoints
  - Add `/keys/{id}` endpoint for key access with proper validation
  - Include service-scoped key retrieval with relationship verification
  - Add participant-specific key operations with authorization
- [ ] **Step 4.2.1.2**: Implement key access authorization
  - Add permission checking for key retrieval with fine-grained controls
  - Include service membership validation with relationship integrity
  - Create participant key access controls with security boundaries
- [ ] **Step 4.2.1.3**: Add key retrieval audit logging
  - Log all key access attempts with comprehensive context
  - Include successful and failed access events with detailed information
  - Add access pattern monitoring and anomaly detection
- [ ] **Step 4.2.1.4**: Implement secure key response format
  - Create secure key serialization for responses with encryption
  - Add key material protection in transit with proper headers
  - Include key metadata in responses with security filtering
- [ ] **Step 4.2.1.5**: Create key retrieval tests
  - Test authorized key access with various permission scenarios
  - Test access control enforcement and security boundaries
  - Test audit logging functionality and completeness

#### **Task 4.2.2: Key Versioning System**
- [ ] **Step 4.2.2.1**: Implement version-aware key retrieval
  - Add version specification in key requests with validation
  - Create latest version resolution with proper caching
  - Include version history access with authorization checks
- [ ] **Step 4.2.2.2**: Add backward compatibility handling
  - Support requests for previous key versions with deprecation warnings
  - Include deprecated version warnings with proper communication
  - Create version migration assistance and guidance
- [ ] **Step 4.2.2.3**: Create version-specific operations
  - Add version tagging and labeling with metadata integration
  - Include version comparison utilities with semantic understanding
  - Create version dependency tracking and relationship management
- [ ] **Step 4.2.2.4**: Implement version cleanup policies
  - Add old version retention policies with configurable parameters
  - Create automated version cleanup with safety checks
  - Include version archival procedures with secure disposal

---

## Phase 5: Production Readiness

### Target 5.1: Persistent Storage Backends
#### **Task 5.1.1: Database Integration**
- [ ] **Step 5.1.1.1**: Create database connection management
  - Implement connection pool configuration with optimal parameters
  - Add database health checking with dependency monitoring
  - Include connection retry logic with exponential backoff
- [ ] **Step 5.1.1.2**: Define database schema
  - Create tables for services, participants, keys with proper normalization
  - Add proper indexes and constraints with performance optimization
  - Include audit table structures with tamper-evident design
- [ ] **Step 5.1.1.3**: Implement database migrations
  - Create migration framework with version control and rollback
  - Add schema version tracking with consistency verification
  - Include rollback capabilities with safety checks
- [ ] **Step 5.1.1.4**: Add database storage implementation
  - Implement storage interface with database backend and optimization
  - Include transaction management with proper isolation levels
  - Add query optimization and performance monitoring
- [ ] **Step 5.1.1.5**: Create database tests
  - Test connection management and pool behavior under load
  - Test transaction handling and isolation levels
  - Test migration procedures and rollback scenarios

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

## Development Standards

### Session-Based Development
- **Step Granularity**: Each numbered step completable in 20-30 minutes
- **Incremental Progress**: Steps build upon previous completed work
- **Testable Units**: Each step includes comprehensive test requirements
- **Documentation**: Step completion includes code documentation and validation

### Code Quality Standards
- **Modern Go Idioms**: Use `any` instead of `interface{}`, proper error handling patterns
- **Security First**: Key material protection always, comprehensive validation everywhere
- **Test-Driven Development**: Comprehensive test coverage with table-driven tests and security focus
- **Documentation**: Go doc comments for all public APIs with security considerations
- **File Organization**: `snake_case` file names, alphabetical organization, proper package structure

### Git Workflow
- Feature branches for logical step groups (e.g., 2.1.1.x steps)
- Comprehensive commit messages referencing step numbers and architectural decisions
- Code review requirements for all security-related changes
- Integration tests required before merging to main branch

### Milestone Tracking
- Each numbered step requires completion validation and verification
- Regular progress checkpoints at sub-phase completion with documentation updates
- Technical decision documentation for major choices with rationale
- Architecture decision records (ADRs) for significant changes affecting future development

---

## Technical Architecture Decisions

### Established in Phase 2.1
- **HTTP Router**: Standard library `net/http` with Go 1.22+ enhanced routing patterns
- **Configuration**: Structured configuration with environment variable support and validation
- **Constructor Pattern**: `NewServer(config ServerConfig) (*Server, error)` with comprehensive validation
- **Context Management**: Typed context keys with shutdown timeout configuration
- **Error Handling**: Consistent error prefixes with structured error responses
- **Security Validation**: Comprehensive input validation with injection attack prevention

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

---

*Roadmap Status: Updated to reflect current development progress*  
*Current Development Status: Step 2.1.1.2 - Implement HTTP listener setup*  
*Architecture Foundation: Established with security-first, structured configuration patterns*  
*Next Focus: Complete HTTP server infrastructure before API endpoint implementation*