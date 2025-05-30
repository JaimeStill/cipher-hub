# Separation of Concerns Analysis

## Objective

Analyze the current Cipher Hub codebase to identify and document separation of concerns violations that should be addressed before the next development session. Use the error infrastructure refactoring as the blueprint for proper package responsibility separation.

## Analysis Framework

### Blueprint: Error Infrastructure Refactoring

**Before (Violation)**:
- Error classification logic in `internal/server/error_response.go`
- Core error utilities mixed with HTTP transport concerns
- Framework-specific code handling domain logic

**After (Proper Separation)**:
- Core error classification in `internal/models/errors.go` (domain layer)
- HTTP response structures in `internal/server/error_response.go` (transport layer)
- Clear dependency direction: server → models
- Reusable error utilities across packages

### Analysis Criteria

For each package, evaluate:

1. **Package Responsibility Clarity**:
   - Does the package have a single, well-defined responsibility?
   - Are domain concepts mixed with transport/framework concerns?
   - Are there utilities that should be framework-agnostic?

2. **Dependency Direction**:
   - Do dependencies flow in the correct direction (outer → inner)?
   - Are there circular dependencies or inappropriate coupling?
   - Should core logic be extracted to be reusable?

3. **Reusability Assessment**:
   - Is business logic tied to specific frameworks unnecessarily?
   - Could other packages benefit from extracted utilities?
   - Are there duplicated patterns that could be centralized?

## Packages to Analyze

### `internal/models/`
- **Current State**: Domain models with validation
- **Evaluate**: Are all domain utilities present? Missing reusable validation patterns?

### `internal/server/`
- **Current State**: HTTP server, middleware, transport concerns
- **Evaluate**: Any domain logic that should move to models? HTTP-specific vs generic utilities?

### `internal/config/`
- **Current State**: Environment variable handling and configuration
- **Evaluate**: Are configuration patterns reusable? Should validation move to models?

### `internal/storage/`
- **Current State**: Storage interface definition
- **Evaluate**: Are storage-specific validations mixed with domain validations?

## Documentation Requirements

For each identified violation, document:

1. **Current Location**: Where the code currently exists
2. **Proper Location**: Where it should be moved
3. **Justification**: Why the separation improves architecture
4. **Dependencies**: What other code would be affected
5. **Priority**: Impact on current development roadmap

## Output Format

Create a structured resolution guide artifact covering:

### High Priority Issues
- Violations that block current development
- Circular dependencies or major architectural problems
- Code that prevents proper testing or reusability

### Medium Priority Issues  
- Misplaced utilities that should be extracted
- Domain logic in transport layers
- Reusable patterns that are framework-specific

### Low Priority Issues
- Minor organizational improvements
- Cosmetic separation improvements
- Future extensibility enhancements

### Recommended Actions
- Specific refactoring steps for each priority level
- Order of operations to minimize disruption
- Integration points that need careful handling

## Success Criteria

The analysis is complete when:
- All packages have been evaluated against separation criteria
- Violations are documented with clear remediation paths
- Priority levels are assigned based on development impact
- Roadmap integration points are identified

## Guidelines

- Use precise but concise language that respects context and token limitations
  - Optimize the balance between providing enough context for clarity while removing any 
- Avoid emojis
- Artifact only - no commentary