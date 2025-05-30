# Cipher Hub

**Cipher Hub** is a comprehensive, containerized key management service built in Go that serves as a centralized security layer for cryptographic operations across distributed systems. Designed as a sidecar component, it handles the complete lifecycle of encryption keys while providing standardized REST APIs for key operations, abstracting away cryptographic complexity from application services.

## Project Status

- **Phase**: 2 (HTTP Server Infrastructure)
- **Target**: 2.1 (Basic Server Setup)
- **Task**: 2.1.2 (Middleware Infrastructure)
- **Step**: Step 2.1.2.4 (Create error response formatting middleware)
- **Blockers**: Separation of concerns architecture review required before next development session
    - Execute [`roadblock-separation-of-concerns.md`](./.prompts/roadblock-separation-of-concerns.md) to generate the resolution guide.

See [`roadmap.md`](./roadmap.md) for detailed development roadmap and current implementation details.

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

# Run the service
go run cmd/cipher-hub/main.go
```

### Pre-Commit Workflow

```bash
go fmt ./...
go build ./...
go test ./...
go vet ./...
go mod tidy
```

## Development Approach

Cipher Hub employs a sophisticated two-phase development workflow designed for consistent quality and sustainable progress:

### Development Phase Overview

```
Pre-Session → Session
     ↓           ↓
 step-guide → implementation + docs
```

**Prompt-Engineered Iterative Development**: Each development phase is executed through isolated chat sessions with an LLM using pre-engineered prompts located in the [`.prompts/`](.prompts/) directory, with developer feedback integration as indicated. The [`.workflows/`](.workflows/) directory contains the orchestration patterns that guide each development session through its specific objectives and deliverables.

**Session-Based Development Benefits:**
- **Focused Objectives**: Each phase has singular purpose and clear deliverables
- **Quality Gates**: Built-in validation at phase transitions prevents technical debt
- **Sustainable Pace**: Prevents cognitive overload through clear phase separation
- **AI-Assisted Quality**: Pre-engineered prompts ensure consistent code quality and standards

### Code Quality Standards

- **Go Best Practices**: Modern idioms with standard library focus and security-conscious patterns
- **Security Patterns**: Comprehensive input validation, key material protection, and secure coding
- **Testing Requirements**: >95% test coverage with security-focused edge cases and table-driven tests
- **Documentation**: Complete Go doc comments with security considerations for all public APIs
