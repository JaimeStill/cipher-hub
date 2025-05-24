# Cipher Hub

**Cipher Hub** is a comprehensive, containerized key management service built in Go that serves as a centralized security layer for cryptographic operations across distributed systems. Designed as a sidecar component, it handles the complete lifecycle of encryption keys while providing standardized REST APIs for key operations, abstracting away cryptographic complexity from application services.

## Project Status: Phase 1 Foundation Complete ✅

**Current Phase**: Phase 2.1 - HTTP Server Infrastructure  
**Architecture**: Standard library foundation with security-first design  
**Go Version**: 1.24+ with enhanced routing patterns

### Roadmap Highlights

- **✅ Phase 1**: Foundation architecture with core data models
- **🔄 Phase 2.1**: HTTP server infrastructure (current)
- **📋 Phase 2.2**: Core API endpoints and service management
- **🔐 Phase 3**: Authentication and authorization framework
- **🔑 Phase 4**: Key generation and lifecycle management
- **🏗️ Phase 5**: Production readiness with persistent storage
- **🚀 Phase 6**: Advanced security features and multi-algorithm support

See [`roadmap.md`](./roadmap.md) for detailed development timeline and milestones.

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

# Run the service (basic placeholder currently)
go run cmd/cipher-hub/main.go
```

### Development

#### Code Standards
All development follows the standards documented in [`style-guide.md`](./style-guide.md), including:
- **Go Best Practices**: Modern idioms with standard library focus
- **Security Patterns**: Key material protection and secure coding
- **Testing Requirements**: Comprehensive test coverage with table-driven tests
- **Documentation**: Complete Go doc comments for all public APIs

#### Pre-Commit Workflow
```bash
# Format and validate
go fmt ./...
go build ./...
go test ./...
go vet ./...
go mod tidy

# Use the pre-commit checklist
# See review.md for complete verification steps
```

### Core Documents

The project documentation is organized into focused documents, each serving a specific purpose:

| Document | Primary Purpose | Content Focus | Version |
|----------|-----------------|---------------|---------|
| [`spec.md`](./spec.md) | Technical Architecture | High-level design, capabilities | 1.4 |
| [`roadmap.md`](./roadmap.md) | Development Planning | Phases, milestones, tasks | Current |
| [`style-guide.md`](./style-guide.md) | Implementation Reference | All code examples & standards | **1.2** |
| [`review.md`](./review.md) | Quality Verification | Checklists, verification steps | 1.2 |

**💡 Development Tip**: The [`style-guide.md`](./style-guide.md) serves as the authoritative source for all implementation details, code examples, and technical patterns.

### Project Structure

```
cipher-hub/
├── cmd/cipher-hub/           # Application entry point
├── internal/
│   ├── models/              # Core data models (Phase 1 ✅)
│   ├── storage/             # Storage interface (Phase 1 ✅)  
│   ├── server/              # HTTP server infrastructure (Phase 2.1 🔄)
│   └── handlers/            # HTTP request handlers (Phase 2.1 🔄)
├── docs/                    # Additional documentation
├── spec.md                  # Technical specification
├── roadmap.md              # Development roadmap
├── style-guide.md          # Implementation standards (primary reference)
├── review.md               # Quality assurance checklist
└── README.md               # This file
```

## Key Features

### 🔑 Key Lifecycle Management (Planned)
- **Secure Generation**: Cryptographically secure key creation using `crypto/rand`
- **Safe Storage**: Encryption at rest with proper access controls
- **Controlled Distribution**: Authenticated key retrieval with audit trails
- **Automated Rotation**: Configurable rotation policies with version management

### 🏢 Service Management (Planned)
- **Service Registration**: Logical containers for related participants
- **Participant Management**: Flexible participant types (user/device/service)
- **Access Control**: Fine-grained permissions and authorization
- **Audit Trails**: Comprehensive logging of all operations

### 🛡️ Security First (Planned)
- **Authentication**: Multi-layered authentication with API keys and JWT
- **Authorization**: Role-based access control (RBAC)
- **Encryption**: All sensitive data encrypted at rest and in transit
- **Compliance**: Structured audit logs for regulatory requirements

*Note: Features marked as "Planned" are part of future development phases. See [`roadmap.md`](./roadmap.md) for implementation timeline.*

## Related Projects

- **Cipher Flux**: Secure data transfer service (follow-on project once **Cipher Hub** is complete)

---

**Built with ❤️ using Go standard library and security-first principles**