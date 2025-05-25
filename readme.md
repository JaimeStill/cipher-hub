# Cipher Hub

**Cipher Hub** is a comprehensive, containerized key management service built in Go that serves as a centralized security layer for cryptographic operations across distributed systems. Designed as a sidecar component, it handles the complete lifecycle of encryption keys while providing standardized REST APIs for key operations, abstracting away cryptographic complexity from application services.

## Project Status: HTTP Server Infrastructure In Progress 🔄

**Current Development**: Step 2.1.1.1 Complete ✅ → Step 2.1.1.2 Next ⏳  
**Architecture Foundation**: Security-first HTTP server configuration with comprehensive validation  
**Go Version**: 1.24+ with enhanced routing patterns and standard library focus

### Development Progress Highlights

- **✅ Phase 1**: Foundation architecture with comprehensive data models and storage interface
- **✅ Step 2.1.1.1**: HTTP server configuration structure with security-first validation
- **⏳ Step 2.1.1.2**: HTTP listener implementation and server lifecycle management (current)
- **📋 Phase 2.1**: Complete HTTP server infrastructure with middleware and health checks
- **🔐 Phase 3**: Authentication and authorization framework with API key management
- **🔑 Phase 4**: Key generation and lifecycle management with rotation capabilities
- **🏗️ Phase 5**: Production readiness with persistent storage and monitoring
- **🚀 Phase 6**: Advanced security features and high availability

See [`roadmap.md`](./roadmap.md) for detailed development timeline with granular step-by-step progression.

## Architecture & Design Philosophy

### Security-First Development
- **Input Validation**: Comprehensive validation preventing injection attacks (path, script, resource exhaustion)
- **Secure Defaults**: All security features enabled by default with configurable overrides
- **Key Material Protection**: Cryptographic keys never exposed in logs, serialization, or memory dumps
- **Audit Everything**: Complete audit trails for all security-relevant operations

### Go Standard Library Focus
- **Minimal Dependencies**: Leverages Go's robust standard library to minimize security audit surface
- **Container-Native**: Built specifically for sidecar deployment with health checks and graceful shutdown
- **Modern Patterns**: Uses Go 1.22+ enhanced routing with structured configuration and context management

### Session-Based Development
- **Granular Progress**: 20-30 minute development steps with clear completion criteria
- **Incremental Quality**: Each step includes comprehensive testing and documentation
- **Architectural Consistency**: Established patterns guide all development decisions

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

# Run the service (HTTP server configuration ready, listener in development)
go run cmd/cipher-hub/main.go
```

### Development

#### Code Standards
All development follows the standards documented in [`style-guide.md`](./style-guide.md), including:
- **Go Best Practices**: Modern idioms with standard library focus and security-conscious patterns
- **Security Patterns**: Comprehensive input validation, key material protection, and secure coding
- **Testing Requirements**: >95% test coverage with security-focused edge cases and table-driven tests
- **Documentation**: Complete Go doc comments with security considerations for all public APIs

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
**Step 2.1.1.2**: Implementing HTTP listener setup with `Start()` method
- HTTP server creation using validated ServerConfig timeouts
- Port binding with proper error handling
- Integration with shutdown context for graceful lifecycle management

### Core Documents

The project documentation is organized into focused documents, each serving a specific purpose:

| Document | Primary Purpose | Content Focus | Version |
|----------|-----------------|---------------|---------|
| [`checkpoint.md`](./checkpoint.md) | Development Continuity | Current step state, immediate next tasks | Current |
| [`review.md`](./review.md) | Code Quality Status | Latest review findings, security assessment | Current |
| [`roadmap.md`](./roadmap.md) | Development Planning | Phases, milestones, granular steps | Current |
| [`spec.md`](./spec.md) | Technical Architecture | High-level design, capabilities, decisions | 1.5 |
| [`style-guide.md`](./style-guide.md) | Implementation Reference | All code examples & standards | **1.3** |

**💡 Development Tip**: The [`style-guide.md`](./style-guide.md) serves as the authoritative source for all implementation details, code examples, and established patterns from HTTP server development.

### Project Structure

```
cipher-hub/
├── cmd/cipher-hub/           # Application entry point
├── internal/
│   ├── models/              # Core data models (Phase 1 ✅)
│   ├── storage/             # Storage interface (Phase 1 ✅)  
│   ├── server/              # HTTP server infrastructure (Phase 2.1 🔄)
│   │   ├── server.go        # ServerConfig + Server struct (Step 2.1.1.1 ✅)
│   │   └── server_test.go   # Comprehensive security-focused tests (✅)
│   └── handlers/            # HTTP request handlers (Phase 2.1 📋)
├── checkpoint.md           # Development progress and next steps
├── go.mod                  # Go module definition
├── readme.md               # Project homepage and documentation
├── review.md               # Latest code review findings and quality status
├── roadmap.md              # Development roadmap with granular steps
├── spec.md                 # Technical specification (v1.5)
└── style-guide.md          # Implementation standards (primary reference v1.3)
```

## Implemented Features

### 🏗️ HTTP Server Infrastructure (Phase 2.1 - In Progress)
- **✅ ServerConfig Architecture**: Structured configuration with comprehensive validation
- **✅ Security-First Validation**: Input sanitization preventing injection attacks
- **✅ Timeout Management**: Configurable timeouts with security bounds (1s-5min)
- **✅ Context Integration**: Graceful shutdown with typed context keys
- **⏳ HTTP Listener**: Server lifecycle management (currently implementing)

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

## Development Workflow

### Quality Assurance
- **Pre-Commit Checklist**: Comprehensive validation using [`review.md`](./review.md)
- **Security Focus**: Every commit includes security verification steps
- **Testing Standards**: Table-driven tests with success/failure path coverage
- **Code Review**: Security-related changes require thorough peer review

### Progress Tracking
- **Session-Based Development**: 20-30 minute incremental development steps
- **Milestone Validation**: Clear completion criteria for each development step
- **Checkpoint Documentation**: Current step state and immediate next tasks
- **Architecture Decision Records**: Technical decisions documented with rationale

## Related Projects

- **Cipher Flux**: Secure data transfer service (follow-on project once **Cipher Hub** is complete)

---

**Built with ❤️ using Go standard library and security-first principles**  
*Current Focus: HTTP server infrastructure with security-first validation and lifecycle management*