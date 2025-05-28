# Cipher Hub

**Cipher Hub** is a comprehensive, containerized key management service built in Go that serves as a centralized security layer for cryptographic operations across distributed systems. Designed as a sidecar component, it handles the complete lifecycle of encryption keys while providing standardized REST APIs for key operations, abstracting away cryptographic complexity from application services.

## Project Status: Task 2.1.1 HTTP Server Creation Complete ✅

**Current Development**: Task 2.1.1 Complete ✅ → Task 2.1.2 Current ⏳ → Step 2.1.2.1 Next ⏳  
**Architecture Foundation**: Production-ready HTTP server with graceful shutdown and signal handling  
**Go Version**: 1.24+ with enhanced routing patterns and standard library focus

### Development Progress Highlights

- **✅ Phase 1**: Foundation architecture with comprehensive data models and storage interface
- **🔄 Phase 2**: HTTP server infrastructure implementation
  - **⏳ Target 2.1**: Basic Server Setup (Task 2.1.1 Complete → Task 2.1.2 Current)
    - **✅ Task 2.1.1**: HTTP Server Creation (COMPLETE ✅)
      - **✅ Step 2.1.1.1**: HTTP server configuration structure with security-first validation
      - **✅ Step 2.1.1.2**: HTTP server lifecycle management with Start() method and thread safety
      - **✅ Step 2.1.1.3**: Graceful shutdown mechanism with signal handling (COMPLETE ✅)
    - **⏳ Task 2.1.2**: Middleware Infrastructure (CURRENT TASK ⏳)
      - **⏳ Step 2.1.2.1**: Middleware function signature pattern (IMMEDIATE NEXT ⏳)
    - **📋 Task 2.1.3**: Health Check System
    - **📋 Task 2.1.4**: Handler Framework
  - **📋 Target 2.2**: API foundation with service registration and participant endpoints
  - **📋 Target 2.3**: Initial Integration
- **📋 Phase 3**: Authentication and authorization framework with API key management
- **🔑 Phase 4**: Key generation and lifecycle management with rotation capabilities
- **🏗️ Phase 5**: Production readiness with persistent storage and monitoring
- **🚀 Phase 6**: Advanced security features and high availability

See [`roadmap.md`](./roadmap.md) for detailed development timeline with granular step-by-step progression.

## Getting Started

### Prerequisites
- Go 1.24 or later

### Quick Start
```bash
# Clone and setup
git clone https://github.com/JaimeStill/cipher-hub
cd cipher-hub

# Install dependencies and run tests
go mod tidy
go test ./...

# Run the service (complete HTTP server with graceful shutdown)
go run cmd/cipher-hub/main.go

# Test graceful shutdown
# Press Ctrl+C or send SIGTERM for graceful termination
```

### Development

#### Pre-Commit Workflow
```bash
# Format and validate
go fmt ./...
go build ./...
go test ./...
go vet ./...
go mod tidy

# Security and quality verification
# See review.md for complete pre-commit checklist
```

#### Current Development Focus
**Phase 2 → Target 2.1 → Task 2.1.2 → Step 2.1.2.1**: Implementing middleware function signature pattern
- Define `Middleware` type as `func(http.Handler) http.Handler`
- Create enhanced middleware stack with conditional support
- Implement middleware application pattern with proper chaining
- Leverage completed HTTP server lifecycle for middleware execution

## Core Documents

The project documentation is organized into focused documents, each serving a specific purpose:

| Document | Primary Purpose | Content Focus |
|----------|-----------------|---------------|
| [`checkpoint.md`](./checkpoint.md) | Development Continuity | Current step state, immediate next tasks |
| [`review.md`](./review.md) | Code Quality Status | Latest review findings, security assessment |
| [`roadmap.md`](./roadmap.md) | Development Planning | Phases, milestones, granular steps |
| [`spec.md`](./spec.md) | Technical Architecture | High-level design, capabilities, decisions |
| [`style-guide.md`](./style-guide.md) | Implementation Reference | All code examples & standards |

**💡 Development Tip**: The [`style-guide.md`](./style-guide.md) serves as the authoritative source for all implementation details, code examples, and established patterns.

## Project Structure

```
cipher-hub/
├── cmd/cipher-hub/           # Application entry point with signal handling
├── internal/
│   ├── models/              # Core data models (Phase 1 ✅)
│   ├── storage/             # Storage interface (Phase 1 ✅)  
│   ├── server/              # HTTP server infrastructure (Phase 2.1 ✅)
│   │   ├── server.go        # Complete HTTP server with graceful shutdown (✅)
│   │   └── server_test.go   # Comprehensive security and lifecycle testing (✅)
│   └── handlers/            # HTTP request handlers (Phase 2.1 📋)
├── checkpoint.md           # Development progress and next steps
├── go.mod                  # Go module definition
├── readme.md               # Project homepage and documentation
├── review.md               # Latest code review findings and quality status
├── roadmap.md              # Development roadmap with granular steps
├── spec.md                 # Technical specification
└── style-guide.md          # Implementation standards (primary reference)
```

## Implemented Features

### 🏗️ HTTP Server Infrastructure (Phase 2 → Target 2.1 → Task 2.1.1 Complete ✅)
- **✅ ServerConfig Architecture**: Structured configuration with comprehensive validation
- **✅ Security-First Validation**: Input sanitization preventing injection attacks
- **✅ Timeout Management**: Configurable timeouts with security bounds (1s-5min)
- **✅ Context Integration**: Graceful shutdown with typed context keys
- **✅ HTTP Server Lifecycle**: Complete Start() method with thread safety and resource management
- **✅ Thread Safety**: Production-ready concurrent access patterns with sync.RWMutex
- **✅ Enhanced Port Validation**: Support for port "0" dynamic assignment and security bounds
- **✅ Graceful Shutdown**: Signal handling and HTTP server coordination (COMPLETE ✅)
- **✅ Signal Handling**: SIGINT and SIGTERM support for container orchestration
- **✅ Resource Management**: Complete cleanup on shutdown failure with proper error propagation
- **✅ Context Resolution**: Clear separation between coordination signaling and shutdown timeout

### 🗄️ Core Data Models (Phase 1 - Complete)
- **Service Registration**: Logical containers for related participants with metadata extensibility
- **Participant Management**: Flexible participant types using metadata-driven classification
- **Cryptographic Keys**: Secure key data structures with lifecycle management
- **Storage Interface**: Abstract persistence layer supporting multiple backends

### 🔒 Security Foundation (Phase 1 - Complete)
- **Input Validation**: Comprehensive validation with injection attack prevention
- **Secure Serialization**: Key material protection with `json:"-"` tags
- **Error Handling**: Structured error responses without information leakage
- **Audit-Ready Logging**: Comprehensive logging without sensitive data exposure
- **Thread Safety**: Production-ready concurrency patterns preventing race conditions

### 🚀 Container-Native Features (Phase 2 → Target 2.1 → Task 2.1.1 Complete ✅)
- **Signal Handling**: SIGINT and SIGTERM graceful shutdown for container orchestration
- **Graceful Termination**: In-flight request completion within configured timeout
- **Resource Cleanup**: Proper cleanup of listeners, connections, and server instances
- **Health Check Ready**: Foundation prepared for container health monitoring
- **Environment Configuration**: Ready for environment variable configuration loading

## Planned Features

### 🔑 Key Lifecycle Management (Phase 4)
- **Secure Generation**: Cryptographically secure key creation using `crypto/rand`
- **Safe Storage**: Encryption at rest with proper access controls
- **Controlled Distribution**: Authenticated key retrieval with audit trails
- **Automated Rotation**: Configurable rotation policies with version management

### 🏢 Service Management (Phase 2.2)
- **Service Registration**: RESTful APIs for service and participant management
- **Access Control**: Fine-grained permissions and authorization
- **Audit Trails**: Comprehensive logging of all operations
- **API Standards**: Consistent JSON API with structured error responses

### 🛡️ Enterprise Security (Phase 3)
- **Authentication**: Multi-layered authentication with API keys and JWT
- **Authorization**: Role-based access control (RBAC) with resource-level permissions
- **Encryption**: All sensitive data encrypted at rest and in transit
- **Compliance**: Structured audit logs for regulatory requirements

*Note: Features marked as "Planned" are part of future development phases. See [`roadmap.md`](./roadmap.md) for detailed implementation timeline.*

## Development Approach

Cipher Hub employs a sophisticated three-phase development workflow designed for consistent quality, comprehensive documentation, and sustainable progress:

### Development Phase Overview

```
Pre-Session → Session → Post-Session
     ↓           ↓          ↓
 step-guide → implementation → review + docs
```

**Prompt-Engineered Iterative Development**: Each development phase is executed through isolated chat sessions with an LLM using pre-engineered prompts located in [`.prompts/`](.prompts/) directory, with developer feedback integration as indicated. The [`.workflows/`](.workflows/) directory contains the orchestration patterns that guide each development session through its specific objectives and deliverables.

**Session-Based Development Benefits:**
- **Focused Objectives**: Each phase has singular purpose and clear deliverables
- **Quality Gates**: Built-in validation at phase transitions prevents technical debt
- **Documentation Synchronization**: Ensures all project documentation remains current
- **Sustainable Pace**: Prevents cognitive overload through clear phase separation
- **AI-Assisted Quality**: Pre-engineered prompts ensure consistent code quality and security standards

### Code Quality Standards
All development follows the standards documented in [`style-guide.md`](./style-guide.md), including:
- **Go Best Practices**: Modern idioms with standard library focus and security-conscious patterns
- **Security Patterns**: Comprehensive input validation, key material protection, and secure coding
- **Testing Requirements**: >95% test coverage with security-focused edge cases and table-driven tests
- **Documentation**: Complete Go doc comments with security considerations for all public APIs

## Production Deployment

### Container Deployment
Cipher Hub is designed for container-native deployment with production-ready features:

```bash
# Build the application
go build -o cipher-hub cmd/cipher-hub/main.go

# Run with environment configuration
CIPHER_HUB_HOST=0.0.0.0 \
CIPHER_HUB_PORT=8080 \
CIPHER_HUB_SHUTDOWN_TIMEOUT=30s \
./cipher-hub
```

### Signal Handling
The service properly handles container orchestration signals:
- **SIGINT (Ctrl+C)**: Initiates graceful shutdown with configured timeout
- **SIGTERM**: Container orchestration graceful shutdown signal
- **Graceful Termination**: In-flight requests complete before shutdown

### Health Check Integration
Ready for container orchestration health monitoring:
```yaml
# Kubernetes liveness probe (ready for implementation)
livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 30
```

## Related Projects

- **Cipher Flux**: Secure data transfer service (follow-on project once **Cipher Hub** is complete)

---

**Built with ❤️ using Go standard library and security-first principles**  
*Current Focus: Task 2.1.2 Middleware Infrastructure building on completed Task 2.1.1 HTTP Server Creation*