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

## Development Workflow Standards

### Prompt-Engineered Iterative Development Process

Cipher Hub employs a sophisticated three-phase development workflow designed for consistent quality, comprehensive documentation, and sustainable progress. Each phase represents a distinct AI prompt interaction with specific deliverables and quality gates.

#### Development Phase Overview

```
Pre-Session → Session → Post-Session
     ↓           ↓          ↓
 step-guide → implementation → review + docs
```

**Phase Separation Benefits:**
- **Focused Objectives**: Each phase has singular purpose and clear deliverables
- **Quality Gates**: Built-in validation at phase transitions prevents accumulation of technical debt
- **Documentation Synchronization**: Ensures all project documentation remains current with implementation
- **Sustainable Pace**: Prevents cognitive overload by separating planning, execution, and review

### Pre-Session Phase (.workflows/pre-session.md)

**Objective**: Generate comprehensive step guide for the next development increment

**Process**:
1. **Step Preparation** ([`guide-prepare.md`](.prompts/guide-prepare.md))
   - Analyze current progress against `roadmap.md`
   - Identify the next **Step** tagged as `IMMEDIATE NEXT`
   - Review `checkpoint.md` and `review.md` for lingering issues or decisions required
   - Determine technical and administrative prerequisites

2. **Guide Generation** ([`guide-generate.md`](.prompts/guide-generate.md))
   - Create detailed `step-guide.md` with implementation instructions
   - Include code examples, architectural patterns, and security requirements
   - Define clear completion criteria and verification steps
   - Establish quality benchmarks and testing requirements

3. **Guide Validation** ([`guide-validate.md`](.prompts/guide-validate.md))
   - Analyze step guide for errors, security issues, and implementation inconsistencies
   - Identify opportunities for design pattern improvements
   - Validate alignment with established architectural decisions
   - Ensure Go best practices and security patterns are followed

**Deliverable**: Production-ready `step-guide.md` with comprehensive implementation details

**Quality Gate**: Guide must pass validation review before session phase begins

### Session Phase (.workflows/session.md)

**Objective**: Execute step guide through interactive collaboration and generate updated checkpoint

**Interactive Development Process**:
Work through `step-guide.md` systematically, engaging in collaborative dialogue when:

- **Clarification Needed**: Encountering ambiguous requirements or implementation details
- **Optimization Opportunities**: Discovering alternative approaches that may be superior
- **Problem Resolution**: Identifying potential issues, broken patterns, or non-optimal solutions
- **Creative Enhancement**: Exploring inspired ideas that emerge during implementation
- **Architecture Decisions**: Working through technical choices that impact future development

**Implementation Standards**:
```go
// Follow established patterns from style guide
func NewServer(config ServerConfig) (*Server, error) {
    // Apply validation-first approach
    config.ApplyDefaults()
    if err := config.Validate(); err != nil {
        return nil, err
    }
    
    // Security-conscious implementation
    // Comprehensive error handling
    // Full test coverage requirement
}
```

**Quality Requirements**:
- **Security First**: All implementations include comprehensive input validation
- **Test Coverage**: >95% coverage with security-focused edge cases
- **Documentation**: Complete Go doc comments with security considerations
- **Code Quality**: Passes all linting, formatting, and static analysis

**Session Completion**:
Generate updated `checkpoint.md` ([`checkpoint-generate.md`](.prompts/checkpoint-generate.md)):
- Capture current development state and completed work
- Document any unimplemented ideas or deferred decisions
- Establish context for next session continuation
- Record architectural decisions and their rationale

**Deliverable**: Implemented step with updated `checkpoint.md`

**Quality Gate**: All step completion criteria met, tests passing, documentation updated

### Post-Session Phase (.workflows/post-session.md)

**Objective**: Conduct comprehensive quality review and synchronize all project documentation

#### Phase 1: Code Review Process
Execute comprehensive pre-commit review ([`review-code.md`](.prompts/review-code.md)):

**Review Checklist Categories**:
- **🔧 Code Quality Verification**: Build, formatting, testing, linting compliance
- **🔒 Security Verification**: Key material protection, input validation, authentication
- **📋 Implementation Standards**: Go idioms, resource management, error handling
- **📚 Documentation**: Code comments, package documentation, API documentation
- **🗂️ Project Structure**: File organization, dependency management, compatibility
- **🚀 Performance**: Efficiency considerations, memory usage, optimization opportunities

**Security-Focused Review**:
```bash
# Mandatory security verification steps
go vet ./...                    # Static analysis
go test ./... -race            # Race condition detection
golangci-lint run              # Comprehensive linting
nancy sleuth                   # Dependency vulnerability scanning
```

**Deliverable**: Updated `review.md` with findings, recommendations, and quality metrics

#### Phase 2: Documentation Synchronization
Review and optimize all core project documents ([`review-artifacts.md`](.prompts/review-artifacts.md)):

**Core Document Updates**:
- **`readme.md`** ([`mod-readme.md`](.prompts/mod/mod-readme.md)): Project homepage, current status, getting started
- **`roadmap.md`** ([`mod-roadmap.md`](.prompts/mod/mod-roadmap.md)): Development timeline, next steps, architectural decisions
- **`spec.md`** ([`mod-spec.md`](.prompts/mod/mod-spec.md)): Technical architecture, capabilities, design philosophy
- **`style-guide.md`** ([`mod-style-guide.md`](.prompts/mod/mod-style-guide.md)): Implementation patterns, code examples, standards

**Document Optimization Process**:
1. **Overlap Analysis**: Identify redundant information between documents
2. **Content Optimization**: Ensure each document serves its specific purpose
3. **Consistency Verification**: Validate information consistency across all documents
4. **Currency Check**: Update version numbers, status indicators, and progress markers

**Quality Gate**: All documents current, consistent, and optimized for their intended purpose

### Document Lifecycle Management

#### Actively Maintained Documents
**Session-Generated**:
- `checkpoint.md` - Updated every session completion
- `step-guide.md` - Generated fresh for each development step
- `review.md` - Updated after every post-session review

**Periodically Synchronized**:
- Core project documents updated during post-session phase
- Package documentation (`doc.go`) updated as code evolves
- Architecture decision records maintained throughout development

#### Document Responsibilities
```
checkpoint.md    → Session phase (implementation context)
step-guide.md    → Pre-session phase (next step instructions)  
review.md        → Post-session phase (quality assessment)
readme.md        → Post-session phase (project overview)
roadmap.md       → Post-session phase (development planning)
spec.md          → Post-session phase (technical architecture)
style-guide.md   → Post-session phase (implementation reference)
```

### Quality Assurance Integration

#### Built-in Quality Gates
- **Pre-Session**: Step guide validation prevents flawed implementation plans
- **Session**: Interactive collaboration catches issues during development
- **Post-Session**: Comprehensive review ensures no quality regression

#### Continuous Quality Metrics
- **Code Quality Score**: Comprehensive assessment across multiple dimensions
- **Security Posture**: Regular security validation and threat assessment
- **Test Coverage**: Maintained >95% with security-focused test scenarios
- **Documentation Currency**: All documents synchronized with implementation state

#### Technical Debt Management
- **Issue Identification**: Systematic identification during post-session reviews
- **Priority Classification**: High/Medium/Low priority with clear remediation paths
- **Progress Tracking**: Integration with roadmap for systematic debt reduction

### Workflow Benefits and Outcomes

#### Development Velocity
- **Predictable Progress**: Granular 20-30 minute development steps
- **Reduced Context Switching**: Clear phase separation prevents cognitive overload
- **Quality Prevention**: Issues caught early rather than accumulated as technical debt

#### Documentation Excellence
- **Living Documentation**: All documents remain current with implementation
- **Multiple Perspectives**: Each document serves specific audience needs
- **Comprehensive Coverage**: From high-level architecture to implementation details

#### Quality Assurance
- **Multi-Layered Review**: Planning, implementation, and retrospective quality gates
- **Security Integration**: Security considerations embedded throughout workflow
- **Sustainable Quality**: Quality built-in rather than bolted-on

#### Knowledge Management
- **Context Preservation**: `checkpoint.md` ensures continuity across sessions
- **Decision Documentation**: Architectural decisions captured with rationale
- **Pattern Evolution**: `style-guide.md` evolves with established patterns

### Pre-Commit Standards
```bash
# Required commands before every commit
go fmt ./...                   # Code formatting
go build ./...                 # Compilation verification
go test ./...                  # Test execution
go vet ./...                   # Static analysis
go mod tidy                    # Dependency cleanup
```

### Commit Message Standards
```
type(scope): description

feat(server): implement HTTP server configuration with security validation
fix(validation): prevent hostname injection attacks in ServerConfig  
docs(api): update API documentation for health check endpoints
test(security): add comprehensive input validation tests
refactor(config): extract validation logic into smaller functions
```

## Related Projects

- **Cipher Flux**: Secure data transfer service (follow-on project once **Cipher Hub** is complete)

---

**Built with ❤️ using Go standard library and security-first principles**  
*Current Focus: HTTP server infrastructure with security-first validation and lifecycle management*