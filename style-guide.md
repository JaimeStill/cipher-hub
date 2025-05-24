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
```go
// Correct: Returns error for validation failures
func NewServiceRegistration(name, description string) (*ServiceRegistration, error) {
    if name == "" {
        return nil, fmt.Errorf("invalid name: %w", ErrInvalidName)
    }updated document draft
    // ... implementation
}

// Incorrect: Panics on invalid input
func NewServiceRegistration(name, description string) *ServiceRegistration {
    if name == "" {
        panic("invalid name")
    }
    // ... implementation
}
```

### Time Field Strategy
- **Required timestamps**: Use `time.Time` for semantic clarity
- **Optional timestamps**: Use `*time.Time` for JSON serialization benefits
- **Examples**: `CreatedAt time.Time`, `ExpiresAt *time.Time`

---

## Security-First Development

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

### Input Validation
- **Constructor Validation**: Validate all inputs at object creation
- **Runtime Validation**: Provide `IsValid()` methods for ongoing validation
- **Error Types**: Use defined error variables (e.g., `ErrInvalidID`) for consistent error handling
- **Sanitization**: Sanitize all user inputs before processing

### Security Headers and Practices
- **HTTP Headers**: Implement comprehensive security headers
- **Authentication**: Validate authentication on every protected endpoint
- **Authorization**: Check permissions for every resource access
- **Audit Logging**: Log all security-relevant operations

---

## Testing Requirements

### Test Coverage Standards
- **Unit Tests**: Every public function and method must have unit tests
- **Table-Driven Tests**: Use table-driven tests for multiple input scenarios
- **Error Cases**: Test both success and failure paths
- **Edge Cases**: Include boundary conditions and edge cases

### Test Organization
```go
func TestNewServiceRegistration(t *testing.T) {
    tests := []struct {
        name        string
        inputName   string
        inputDesc   string
        wantErr     error
    }{
        {
            name:      "valid input",
            inputName: "Test Service",
            inputDesc: "Test Description",
            wantErr:   nil,
        },
        {
            name:      "empty name",
            inputName: "",
            inputDesc: "Test Description",
            wantErr:   ErrInvalidName,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ... test implementation
        })
    }
}
```

### Integration Testing
- **HTTP API Tests**: End-to-end API testing for all endpoints
- **Storage Tests**: Test storage implementations with real backends
- **Security Tests**: Verify authentication and authorization enforcement

---

## Documentation Standards

### Go Doc Comments
- **Public APIs**: All public functions, types, and methods require doc comments
- **Package Documentation**: Every package needs a package-level doc comment
- **Security Notes**: Document security considerations for sensitive operations
- **Usage Examples**: Include examples for complex APIs

```go
// NewServiceRegistration creates a new service registration with the given name and description.
// It generates a cryptographically secure ID and initializes timestamps.
// Returns an error if ID generation fails or if validation fails.
//
// Security: The service registration acts as a security boundary for related participants.
func NewServiceRegistration(name, description string) (*ServiceRegistration, error) {
    // ... implementation
}
```

### README Standards
- **Getting Started**: Clear setup and running instructions
- **API Documentation**: Basic API usage examples
- **Security Considerations**: Highlight security-relevant configuration
- **Development Setup**: Instructions for development environment

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

## Git Workflow

### Branch Strategy
- **Feature Branches**: Create feature branches for each roadmap item
- **Naming**: Use descriptive branch names (e.g., `feature/http-server-setup`)
- **Main Protection**: Main branch requires pull request reviews
- **No Direct Commits**: Never commit directly to main branch

### Commit Standards
- **Message Format**: Use conventional commit format
  ```
  type(scope): description
  
  Optional body explaining the change
  ```
- **Types**: Use `feat`, `fix`, `docs`, `test`, `refactor`, `security`
- **Reference Issues**: Reference roadmap items or issues in commit messages
- **Atomic Commits**: Each commit should represent a single logical change

### Code Review Requirements
- **Security Changes**: All security-related changes require thorough review
- **Test Coverage**: New code must include appropriate tests
- **Documentation**: Public API changes require documentation updates
- **Performance**: Performance-sensitive changes need benchmarks

---

## Project Management

### Milestone Tracking
- **Phase Completion**: Each phase requires completion criteria validation
- **Progress Updates**: Regular updates to roadmap and checkpoint documents
- **Technical Decisions**: Document major technical decisions in code comments
- **Architecture Records**: Maintain ADRs for significant architectural changes

### Quality Gates
- **Pre-Commit**: Run tests and linting before committing
- **CI/CD Pipeline**: Automated testing and security scanning
- **Code Coverage**: Maintain minimum test coverage thresholds
- **Security Scanning**: Regular security vulnerability scanning

### Documentation Maintenance
- **Living Documents**: Keep roadmap and specifications current
- **Change Documentation**: Update relevant docs when making changes
- **Checkpoint Updates**: Update checkpoint document at major milestones
- **API Documentation**: Keep API docs synchronized with implementation

---

## Code Quality Enforcement

### Static Analysis
- **go fmt**: All code must be formatted with `go fmt`
- **go vet**: Code must pass `go vet` analysis
- **Linting**: Use `golangci-lint` with security-focused rules
- **Import Organization**: Use `goimports` for import management

### Pre-Commit Checklist
- [ ] Code formatted with `go fmt`
- [ ] All tests pass (`go test ./...`)
- [ ] No security vulnerabilities in dependencies
- [ ] Documentation updated for public API changes
- [ ] Commit message follows conventional format

### Performance Standards
- [ ] No unnecessary allocations in hot paths
- [ ] Proper resource cleanup (defer statements)
- [ ] Context cancellation handling
- [ ] Efficient error handling without performance impact

---

## Security Checklist

### Before Every Commit
- [ ] No key material in logs or error messages
- [ ] Input validation on all user-provided data
- [ ] Proper error handling without information leakage
- [ ] Authentication and authorization checks in place
- [ ] Audit logging for security-relevant operations

### Code Review Focus Areas
- [ ] Key material protection and secure handling
- [ ] Input validation and sanitization
- [ ] Error handling and information disclosure
- [ ] Authentication and authorization implementation
- [ ] Logging and audit trail completeness

---

*Style Guide Version: 1.0*  
*Last Updated: Current Session*  
*Status: Aligned with Foundation Phase Standards*