# Cipher Hub - Refined Development Roadmap

## Current Status: Phase 1 Foundation Complete ✅

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
- [ ] **Step 2.1.1.1**: Create basic HTTP server struct with configuration fields ⏳ **IMMEDIATE NEXT**
  - Define `Server` struct with host, port, and timeout fields
  - Add basic constructor `NewServer(host, port string) *Server`
  - Include context for graceful shutdown
- [ ] **Step 2.1.1.2**: Implement basic HTTP listener setup
  - Add `Start()` method that creates `http.Server` instance
  - Configure basic timeouts (read, write, idle)
  - Add port binding and listener creation
- [ ] **Step 2.1.1.3**: Add graceful shutdown mechanism
  - Implement `Shutdown(ctx context.Context) error` method
  - Handle `os.Signal` for SIGINT and SIGTERM
  - Ensure in-flight requests complete before shutdown
- [ ] **Step 2.1.1.4**: Create server configuration structure
  - Define `Config` struct for server settings
  - Add environment variable loading for host/port
  - Include validation for configuration values
- [ ] **Step 2.1.1.5**: Add basic server tests
  - Test server start/stop functionality
  - Test graceful shutdown behavior
  - Test configuration validation

#### **Task 2.1.2: Middleware Infrastructure** (`internal/server/middleware.go`)
- [ ] **Step 2.1.2.1**: Create middleware function signature pattern
  - Define `Middleware` type as `func(http.Handler) http.Handler`
  - Create `Chain` function to combine multiple middlewares
  - Add basic middleware application pattern
- [ ] **Step 2.1.2.2**: Implement request logging middleware
  - Add request ID generation using `crypto/rand`
  - Log HTTP method, path, and response status
  - Include request duration timing
- [ ] **Step 2.1.2.3**: Add CORS handling middleware
  - Support basic CORS headers (Origin, Methods, Headers)
  - Add configurable allowed origins
  - Handle preflight OPTIONS requests
- [ ] **Step 2.1.2.4**: Create error response formatting middleware
  - Standardize JSON error response format
  - Add error code mapping from internal errors
  - Ensure no sensitive data leaks in error responses
- [ ] **Step 2.1.2.5**: Implement security headers middleware
  - Add `X-Content-Type-Options: nosniff`
  - Include `X-Frame-Options: DENY`
  - Set `X-XSS-Protection: 1; mode=block`
- [ ] **Step 2.1.2.6**: Add middleware tests
  - Test request ID generation and propagation
  - Test CORS header setting
  - Test error response formatting

#### **Task 2.1.3: Health Check System** (`internal/handlers/health.go`)
- [ ] **Step 2.1.3.1**: Create basic health check handler structure
  - Define `HealthHandler` struct with dependencies
  - Add constructor `NewHealthHandler() *HealthHandler`
  - Set up basic handler method signatures
- [ ] **Step 2.1.3.2**: Implement liveness endpoint
  - Add `/health/live` endpoint that always returns 200 OK
  - Return simple JSON response `{"status": "alive"}`
  - Include timestamp in response
- [ ] **Step 2.1.3.3**: Implement readiness endpoint foundation
  - Add `/health/ready` endpoint structure
  - Return basic readiness status without dependencies
  - Use JSON format `{"status": "ready", "checks": []}`
- [ ] **Step 2.1.3.4**: Add health check response models
  - Create `HealthStatus` struct for JSON responses
  - Define `CheckResult` struct for individual checks
  - Add proper JSON tags and validation
- [ ] **Step 2.1.3.5**: Create health check tests
  - Test liveness endpoint returns 200
  - Test readiness endpoint basic functionality
  - Test JSON response format validation

#### **Task 2.1.4: Handler Framework** (`internal/handlers/handlers.go`)
- [ ] **Step 2.1.4.1**: Create base handler utilities
  - Add `writeJSON` utility function for consistent responses
  - Create `readJSON` utility for request parsing
  - Include proper content-type handling
- [ ] **Step 2.1.4.2**: Implement error response utilities
  - Add `writeError` function with status codes
  - Create `ValidationError` type for input errors
  - Ensure consistent error JSON format
- [ ] **Step 2.1.4.3**: Add request parsing helpers
  - Create URL parameter extraction utilities
  - Add query parameter parsing helpers
  - Include request body size limiting
- [ ] **Step 2.1.4.4**: Implement response header utilities
  - Add common header setting functions
  - Include cache control helpers
  - Create security header application functions
- [ ] **Step 2.1.4.5**: Create handler framework tests
  - Test JSON response utilities
  - Test error response formatting
  - Test request parsing edge cases

**Technical Decisions Needed:**
- HTTP Router: Standard library `net/http` vs lightweight router
- Configuration: Environment variables, config files, or hybrid approach
- Logging: Standard library vs structured logging
- Middleware Pattern: Choose implementation approach

### Target 2.2: API Foundation

#### **Task 2.2.1: Service Registration Endpoints**
- [ ] **Step 2.2.1.1**: Create service registration handler structure
  - Define `ServiceHandler` struct with storage dependency
  - Add constructor with storage interface injection
  - Set up route method signatures
- [ ] **Step 2.2.1.2**: Implement GET /services endpoint
  - Add list services functionality
  - Include pagination parameters (limit, offset)
  - Return JSON array of service summaries
- [ ] **Step 2.2.1.3**: Implement GET /services/{id} endpoint
  - Add single service retrieval by ID
  - Handle service not found cases
  - Return full service details with participants
- [ ] **Step 2.2.1.4**: Implement POST /services endpoint
  - Add service creation from JSON request
  - Validate required fields (name, description)
  - Return created service with 201 status
- [ ] **Step 2.2.1.5**: Implement PUT /services/{id} endpoint
  - Add service update functionality
  - Validate update request format
  - Handle partial updates vs full replacement
- [ ] **Step 2.2.1.6**: Implement DELETE /services/{id} endpoint
  - Add service deletion with safety checks
  - Ensure no active participants exist
  - Return appropriate status codes
- [ ] **Step 2.2.1.7**: Add service endpoint tests
  - Test CRUD operations end-to-end
  - Test error conditions and edge cases
  - Test input validation behavior

#### **Task 2.2.2: Participant Management Endpoints**
- [ ] **Step 2.2.2.1**: Create participant handler structure
  - Define `ParticipantHandler` struct
  - Add service-scoped participant operations
  - Set up nested route patterns
- [ ] **Step 2.2.2.2**: Implement GET /services/{id}/participants
  - List participants within a service
  - Add participant type filtering
  - Include pagination support
- [ ] **Step 2.2.2.3**: Implement POST /services/{id}/participants
  - Add participant creation within service
  - Validate participant data and type
  - Return created participant details
- [ ] **Step 2.2.2.4**: Implement GET /services/{id}/participants/{pid}
  - Retrieve single participant by ID
  - Validate service and participant relationship
  - Return full participant details
- [ ] **Step 2.2.2.5**: Implement PUT /services/{id}/participants/{pid}
  - Update participant information
  - Validate metadata updates
  - Handle participant type changes
- [ ] **Step 2.2.2.6**: Implement DELETE /services/{id}/participants/{pid}
  - Remove participant from service
  - Check for dependent key relationships
  - Return appropriate confirmation
- [ ] **Step 2.2.2.7**: Add participant endpoint tests
  - Test nested resource operations
  - Test service-participant relationship validation
  - Test metadata handling

#### **Task 2.2.3: Request/Response Standards**
- [ ] **Step 2.2.3.1**: Define API response envelope structure
  - Create consistent response wrapper format
  - Include data, metadata, and error fields
  - Add response timestamp and request ID
- [ ] **Step 2.2.3.2**: Implement pagination response format
  - Add pagination metadata structure
  - Include total count, limit, offset fields
  - Create pagination link generation
- [ ] **Step 2.2.3.3**: Create input validation framework
  - Define validation error response format
  - Add field-level validation error details
  - Include validation rule descriptions
- [ ] **Step 2.2.3.4**: Implement API versioning strategy
  - Add version header handling
  - Create version-specific route prefixes
  - Plan backward compatibility approach
- [ ] **Step 2.2.3.5**: Add response standard tests
  - Test consistent response formatting
  - Test pagination metadata accuracy
  - Test validation error responses

#### **Task 2.2.4: Basic Configuration System**
- [ ] **Step 2.2.4.1**: Create configuration struct definitions
  - Define `AppConfig` with all necessary fields
  - Add server, database, and security sections
  - Include environment-specific defaults
- [ ] **Step 2.2.4.2**: Implement environment variable loading
  - Add `LoadFromEnv()` configuration method
  - Support standard environment variable patterns
  - Include fallback to default values
- [ ] **Step 2.2.4.3**: Add configuration validation
  - Create `Validate()` method for configuration
  - Check required fields and value ranges
  - Validate interdependent configuration values
- [ ] **Step 2.2.4.4**: Implement configuration file support
  - Add YAML/JSON configuration file parsing
  - Support configuration file path via environment
  - Merge file and environment configurations
- [ ] **Step 2.2.4.5**: Create configuration tests
  - Test environment variable loading
  - Test validation error conditions
  - Test configuration merging logic

### Target 2.3: Initial Integration

#### **Task 2.3.1: In-Memory Storage Implementation**
- [ ] **Step 2.3.1.1**: Create in-memory storage struct
  - Define `MemoryStorage` struct with sync.RWMutex
  - Add maps for services, participants, and keys
  - Implement storage interface methods signatures
- [ ] **Step 2.3.1.2**: Implement service storage operations
  - Add service CRUD operations with thread safety
  - Include service ID generation and validation
  - Handle concurrent access with proper locking
- [ ] **Step 2.3.1.3**: Implement participant storage operations
  - Add participant operations within service context
  - Maintain service-participant relationships
  - Include participant lookup and filtering
- [ ] **Step 2.3.1.4**: Add storage operation tests
  - Test concurrent access patterns
  - Test data consistency under load
  - Test error conditions and edge cases
- [ ] **Step 2.3.1.5**: Implement storage interface compliance
  - Verify all interface methods implemented
  - Add context cancellation support
  - Include proper error handling patterns

#### **Task 2.3.2: API Integration Testing**
- [ ] **Step 2.3.2.1**: Create integration test framework
  - Set up test server with in-memory storage
  - Add test utilities for HTTP requests
  - Include test data setup and teardown
- [ ] **Step 2.3.2.2**: Implement service API integration tests
  - Test complete service lifecycle via API
  - Include error scenario testing
  - Verify response format compliance
- [ ] **Step 2.3.2.3**: Add participant API integration tests
  - Test nested resource operations via API
  - Verify service-participant relationship handling
  - Test concurrent participant operations
- [ ] **Step 2.3.2.4**: Create end-to-end workflow tests
  - Test complete service setup workflows
  - Include multi-step operations
  - Verify data consistency across operations
- [ ] **Step 2.3.2.5**: Add performance baseline tests
  - Measure response times for basic operations
  - Test concurrent request handling
  - Establish performance baseline metrics

---

## Phase 3: Security Foundation

### Target 3.1: Authentication System
#### **Task 3.1.1: API Key Authentication**
- [ ] **Step 3.1.1.1**: Create API key data structure
  - Define `APIKey` struct with ID, key hash, and metadata
  - Add key generation using `crypto/rand`
  - Include key expiration and status fields
- [ ] **Step 3.1.1.2**: Implement API key hashing
  - Add secure key hashing using bcrypt or similar
  - Create key comparison functions
  - Include salt generation and validation
- [ ] **Step 3.1.1.3**: Create API key storage interface
  - Define storage methods for API key operations
  - Add key lookup by hash functionality
  - Include key status management
- [ ] **Step 3.1.1.4**: Implement API key validation middleware
  - Add `Authorization` header parsing
  - Create key lookup and validation logic
  - Handle authentication failure responses
- [ ] **Step 3.1.1.5**: Add API key management endpoints
  - Create key generation endpoint
  - Add key listing and revocation endpoints
  - Include key status update functionality
- [ ] **Step 3.1.1.6**: Create API key tests
  - Test key generation and validation
  - Test authentication middleware
  - Test key management operations

#### **Task 3.1.2: Request Authentication**
- [ ] **Step 3.1.2.1**: Implement bearer token parsing
  - Add `Authorization: Bearer <token>` support
  - Create token extraction utilities
  - Handle malformed authorization headers
- [ ] **Step 3.1.2.2**: Add request context authentication
  - Store authenticated principal in request context
  - Create context extraction utilities
  - Include authentication status helpers
- [ ] **Step 3.1.2.3**: Create authentication error handling
  - Define authentication-specific error types
  - Add proper HTTP status code mapping
  - Ensure no sensitive data in error responses
- [ ] **Step 3.1.2.4**: Implement authentication bypass for health checks
  - Allow unauthenticated access to health endpoints
  - Create authentication exemption patterns
  - Test public endpoint accessibility
- [ ] **Step 3.1.2.5**: Add authentication integration tests
  - Test authenticated request workflows
  - Test authentication failure scenarios
  - Test context propagation through handlers

### Target 3.2: Authorization Framework
#### **Task 3.2.1: Role-Based Access Control (RBAC)**
- [ ] **Step 3.2.1.1**: Define permission and role structures
  - Create `Permission` enum for resource operations
  - Define `Role` struct with permission collections
  - Add role assignment and inheritance patterns
- [ ] **Step 3.2.1.2**: Implement permission checking functions
  - Add `HasPermission(user, resource, operation)` function
  - Create resource-specific permission validation
  - Include role-based permission resolution
- [ ] **Step 3.2.1.3**: Create authorization middleware
  - Add permission checking to request pipeline
  - Include resource identification from request
  - Handle authorization failure responses
- [ ] **Step 3.2.1.4**: Add role management functionality
  - Create role assignment operations
  - Include role update and revocation
  - Add role inheritance processing
- [ ] **Step 3.2.1.5**: Implement resource-level authorization
  - Add service-specific access controls
  - Include participant-level permissions
  - Create key operation authorization
- [ ] **Step 3.2.1.6**: Create authorization tests
  - Test permission checking logic
  - Test authorization middleware
  - Test role management operations

### Target 3.3: Secure Key Storage
#### **Task 3.3.1: Encryption at Rest**
- [ ] **Step 3.3.1.1**: Create master key management
  - Implement master key generation and storage
  - Add key derivation for data encryption keys
  - Include master key rotation capabilities
- [ ] **Step 3.3.1.2**: Implement key material encryption
  - Add encryption before storage persistence
  - Create decryption on key retrieval
  - Include encrypted field tagging
- [ ] **Step 3.3.1.3**: Add secure memory handling
  - Implement secure memory allocation for keys
  - Add memory clearing after key usage
  - Include memory protection from dumps
- [ ] **Step 3.3.1.4**: Create key serialization protection
  - Ensure encrypted keys never serialize as plaintext
  - Add secure key marshaling/unmarshaling
  - Include key material protection in logs
- [ ] **Step 3.3.1.5**: Implement audit trail encryption
  - Add audit log encryption capabilities
  - Create tamper-evident audit records
  - Include audit integrity verification

---

## Phase 4: Key Lifecycle Management

### Target 4.1: Key Generation Engine
#### **Task 4.1.1: AES-256 Key Generation**
- [ ] **Step 4.1.1.1**: Implement secure random key generation
  - Use `crypto/rand` for cryptographically secure randomness
  - Add AES-256 key length validation (32 bytes)
  - Include key format standardization
- [ ] **Step 4.1.1.2**: Create key generation request handling
  - Add key generation API endpoint
  - Include algorithm specification in requests
  - Validate generation parameters
- [ ] **Step 4.1.1.3**: Implement key quality validation
  - Add entropy checking for generated keys
  - Include weak key detection
  - Create key strength verification
- [ ] **Step 4.1.1.4**: Add key generation audit logging
  - Log all key generation events
  - Include generation metadata and context
  - Ensure no key material in audit logs
- [ ] **Step 4.1.1.5**: Create key generation tests
  - Test key randomness and uniqueness
  - Test key format compliance
  - Test generation error handling

#### **Task 4.1.2: Key Metadata Management**
- [ ] **Step 4.1.2.1**: Define key metadata structure
  - Create comprehensive key metadata fields
  - Include version, status, and usage tracking
  - Add creation and modification timestamps
- [ ] **Step 4.1.2.2**: Implement key versioning system
  - Add version tracking for key updates
  - Create version history maintenance
  - Include version-specific operations
- [ ] **Step 4.1.2.3**: Add key status lifecycle management
  - Define key status states (active, deprecated, revoked)
  - Implement status transition validation
  - Include status change audit logging
- [ ] **Step 4.1.2.4**: Create usage statistics tracking
  - Add key usage counters and metrics
  - Include last access timestamp tracking
  - Create usage pattern analysis
- [ ] **Step 4.1.2.5**: Implement metadata persistence
  - Add metadata storage operations
  - Include metadata query capabilities
  - Create metadata update validation

### Target 4.2: Key Distribution APIs
#### **Task 4.2.1: Secure Key Retrieval**
- [ ] **Step 4.2.1.1**: Create key retrieval endpoints
  - Add `/keys/{id}` endpoint for key access
  - Include service-scoped key retrieval
  - Add participant-specific key operations
- [ ] **Step 4.2.1.2**: Implement key access authorization
  - Add permission checking for key retrieval
  - Include service membership validation
  - Create participant key access controls
- [ ] **Step 4.2.1.3**: Add key retrieval audit logging
  - Log all key access attempts
  - Include successful and failed access events
  - Add access pattern monitoring
- [ ] **Step 4.2.1.4**: Implement secure key response format
  - Create secure key serialization for responses
  - Add key material protection in transit
  - Include key metadata in responses
- [ ] **Step 4.2.1.5**: Create key retrieval tests
  - Test authorized key access
  - Test access control enforcement
  - Test audit logging functionality

#### **Task 4.2.2: Key Versioning System**
- [ ] **Step 4.2.2.1**: Implement version-aware key retrieval
  - Add version specification in key requests
  - Create latest version resolution
  - Include version history access
- [ ] **Step 4.2.2.2**: Add backward compatibility handling
  - Support requests for previous key versions
  - Include deprecated version warnings
  - Create version migration assistance
- [ ] **Step 4.2.2.3**: Create version-specific operations
  - Add version tagging and labeling
  - Include version comparison utilities
  - Create version dependency tracking
- [ ] **Step 4.2.2.4**: Implement version cleanup policies
  - Add old version retention policies
  - Create automated version cleanup
  - Include version archival procedures

---

## Phase 5: Production Readiness

### Target 5.1: Persistent Storage Backends
#### **Task 5.1.1: Database Integration**
- [ ] **Step 5.1.1.1**: Create database connection management
  - Implement connection pool configuration
  - Add database health checking
  - Include connection retry logic
- [ ] **Step 5.1.1.2**: Define database schema
  - Create tables for services, participants, keys
  - Add proper indexes and constraints
  - Include audit table structures
- [ ] **Step 5.1.1.3**: Implement database migrations
  - Create migration framework
  - Add schema version tracking
  - Include rollback capabilities
- [ ] **Step 5.1.1.4**: Add database storage implementation
  - Implement storage interface with database backend
  - Include transaction management
  - Add query optimization
- [ ] **Step 5.1.1.5**: Create database tests
  - Test connection management
  - Test transaction handling
  - Test migration procedures

---

## Success Metrics

### Phase 2 Success Criteria
- HTTP server handles basic CRUD operations for all core entities
- Health checks integrate with container orchestration
- Consistent API patterns established across all endpoints
- Integration tests cover full request/response cycles

### Phase 3 Success Criteria
- Secure authentication prevents unauthorized access
- Authorization controls enforce proper permissions
- Key material never exposed in logs or serialization
- Audit trails capture all security-relevant events

### Long-term Success Criteria
- Handle 10,000+ concurrent key operations
- Sub-100ms response times for key retrieval
- 99.9% uptime in production deployment
- Zero key material exposure incidents
- Complete audit trails for compliance requirements

---

## Development Standards

### Session-Based Development
- **Task Granularity**: Each numbered task completable in 20-30 minutes
- **Incremental Progress**: Tasks build upon previous completed work
- **Testable Units**: Each task includes test requirements
- **Documentation**: Task completion includes code documentation

### Code Quality Standards
- **Modern Go Idioms**: Use `any` instead of `interface{}`, proper error handling
- **Security First**: Key material protection always, validation everywhere
- **Test-Driven Development**: Comprehensive test coverage with table-driven tests
- **Documentation**: Go doc comments for all public APIs
- **File Organization**: `snake_case` file names, alphabetical organization

### Git Workflow
- Feature branches for logical task groups (e.g., 2.1.1.x tasks)
- Comprehensive commit messages referencing task numbers
- Code review requirements for all security-related changes
- Integration tests required before merging to main

### Milestone Tracking
- Each numbered task requires completion validation
- Regular progress checkpoints at sub-phase completion
- Technical decision documentation for major choices
- Architecture decision records (ADRs) for significant changes

---

*Roadmap Status: Updated with granular task breakdown and hierarchical structure*  
*Next Milestone: Step 2.1.1.1 - Create basic HTTP server struct with configuration fields*  
*Current Focus: 20-30 minute development sessions with clear layered progression*