# Step 2.1.2.4: Create Error Response Formatting Middleware

## Overview

**Phase**: 2 (HTTP Server Infrastructure)  
**Target**: 2.1 (Basic Server Setup)  
**Task**: 2.1.2 (Middleware Infrastructure)  
**Step**: 2.1.2.4 (Create error response formatting middleware)

**Time Estimate**: 1 hour  
**Scope**: Implement standardized JSON error response formatting with request correlation and security-conscious error handling

## Step Objectives

### Primary Deliverables
- **Error Response Format**: Standardized JSON error response structure with consistent field naming
- **Error Classification System**: Map internal errors to appropriate HTTP status codes and public error messages
- **Security-Conscious Error Handling**: Prevent sensitive information leakage while providing useful error context
- **Request Correlation Integration**: Include request IDs in error responses for tracing and debugging
- **Error Recovery Middleware**: Panic recovery with structured error responses
- **Comprehensive Testing**: Unit and integration tests following established patterns

### Implementation Requirements
- **Files Created**: Error response middleware implementation and tests
- **Files Modified**: Add error response utilities and server integration examples
- **Architecture Focus**: Security-first error handling with request correlation and standardized responses
- **Security Focus**: Never expose internal implementation details or sensitive information
- **Go Best Practices**: Follow established middleware patterns and error handling conventions
- **Foundation Usage**: Leverage request correlation IDs and structured logging infrastructure

---

## Implementation Requirements

### Technical Specifications

#### Error Response Structure Requirements
- **Standard JSON Format**: Consistent error response schema across all endpoints
- **Request Correlation**: Include request ID from existing request logging middleware
- **Timestamp Information**: ISO 8601 formatted timestamps for error occurrence
- **Error Classification**: Public error codes and categories separate from internal errors
- **Security Compliance**: No sensitive data or implementation details in responses

#### Error Classification Requirements
- **HTTP Status Mapping**: Map internal errors to appropriate HTTP status codes
- **Public Error Codes**: Application-specific error codes safe for external consumption
- **Error Categories**: Group errors by type (validation, authentication, internal, etc.)
- **Message Sanitization**: Safe, user-friendly error messages without internal details

#### Security Requirements
- **Information Disclosure Prevention**: Never expose internal error details or stack traces
- **Error Logging**: Log full internal errors while returning sanitized responses
- **Panic Recovery**: Gracefully handle panics with structured error responses
- **Input Validation**: Validate and sanitize any user input included in error responses

---

## Completion Criteria

### **Step 2.1.2.4 is complete when:**

1. **Core Error Classification in Models Package**:
   - `ErrorCode` type and constants in `internal/models/errors.go`
   - `ErrorCategory` type for consistent error grouping
   - `ClassifyError()` function to map internal errors to public error codes
   - `SanitizeErrorMessage()` function for security-conscious error messages
   - Framework-agnostic error utilities reusable across packages

2. **HTTP Error Response Structure in Server Package**:
   - `ErrorResponse` struct with fields: `Error`, `Message`, `RequestID`, `Timestamp`, optional `Details`
   - `ErrorDetail` struct for structured error information with code and category
   - JSON tags following established naming conventions
   - HTTP status code mapping based on error classifications

3. **Error Response Formatting Middleware**:
   - `ErrorResponseMiddleware()` function with default configuration
   - `ErrorResponseMiddlewareWithConfig()` for custom error handling configuration
   - Panic recovery with structured error responses
   - Integration with request logging for error correlation
   - Uses core error classification from models package

4. **Security-Conscious Error Handling**:
   - Internal error logging with full details using request correlation
   - External error responses with sanitized, user-safe messages from models package
   - Prevention of sensitive information leakage (no stack traces, internal paths, etc.)
   - Configurable error detail levels for different environments

5. **HTTP Status Code Mapping in Server Package**:
   - Map error classifications to appropriate HTTP status codes
   - Support for custom status code mapping
   - Integration with core error classification from models

6. **Request Correlation Integration**:
   - Use `GetRequestID()` from existing request logging infrastructure
   - Include request ID in all error responses for tracing
   - Log errors with request correlation for debugging
   - Maintain request context throughout error handling

7. **Middleware Integration**:
   - Position middleware to catch errors from handlers and other middleware
   - Integration with existing middleware stack patterns
   - Support for method chaining with other middleware
   - Proper error recovery and response completion

8. **Comprehensive Testing Coverage**:
   - Unit tests for core error classification in models package
   - Unit tests for HTTP error response structures and middleware in server package
   - Error mapping tests for different error types and status codes
   - Middleware tests for panic recovery and error formatting
   - Integration tests with server lifecycle and other middleware
   - Security tests ensuring no information leakage
   - Request correlation tests with error scenarios

9. **Production-Ready Code Quality**:
   - Complete Go doc comments for all public error handling functions in both packages
   - Security considerations documented in code comments
   - Usage examples for common error handling scenarios
   - Environment considerations for error detail levels
   - Code passes formatting (`go fmt`) and static analysis (`go vet`)
   - Maintains high test coverage (>95%) with security-focused testing

### **Files Created/Modified**
- `internal/models/errors.go` - Add core error classification and sanitization
- `internal/server/error_response.go` - HTTP error response middleware and utilities
- `internal/server/error_response_test.go` - Comprehensive error response testing
- `internal/models/errors_test.go` - Core error classification testing
- `internal/server/server_test.go` - Add error response integration tests

---

## Testing Requirements

### Unit Testing Requirements

#### Error Response Structure Testing
- Test `ErrorResponse` JSON serialization with all fields
- Test error response creation with different error types
- Test timestamp formatting and request ID inclusion
- Test error detail structure and optional fields

#### Error Classification Testing
- Test internal error mapping to public error codes
- Test HTTP status code mapping for different error types
- Test error message sanitization and security compliance
- Test error category assignment and consistency

#### Error Middleware Testing
- Test middleware with various error types from handlers
- Test panic recovery with structured error responses
- Test request correlation integration in error scenarios
- Test middleware chain integration and error propagation

### Integration Testing Requirements

#### Server Integration Testing
- Test error middleware integration with server lifecycle
- Test error handling with other middleware in the stack
- Test end-to-end error responses through complete request cycle
- Test error logging correlation with request logging middleware

#### Security Testing Requirements
- Test that sensitive information never appears in error responses
- Test error message sanitization prevents information disclosure
- Test panic recovery doesn't expose internal implementation details
- Test error responses under various attack scenarios

### Edge Case Testing Requirements
- Test error handling with malformed requests
- Test concurrent error scenarios
- Test memory exhaustion and resource limit errors
- Test error handling when logging systems fail

---

## Security Considerations

### Error Information Disclosure Prevention
```go
// Correct: Sanitized error response
func sanitizeError(err error) string {
    switch {
    case errors.Is(err, models.ErrInvalidID):
        return "Invalid identifier provided"
    case errors.Is(err, models.ErrInvalidName):
        return "Invalid name provided"
    default:
        return "An internal error occurred" // Never expose internal details
    }
}

// Incorrect: Exposing internal details
func exposeInternalError(err error) string {
    return err.Error() // Could expose internal paths, database details, etc.
}
```

### Panic Recovery Security
```go
// Correct: Secure panic recovery
func recoverPanic() interface{} {
    if r := recover(); r != nil {
        // Log full panic details internally
        slog.Error("Panic recovered", "panic", r, "stack", debug.Stack())
        // Return generic error to client
        return "Internal server error"
    }
    return nil
}

// Incorrect: Exposing panic details
func unsafeRecover() interface{} {
    if r := recover(); r != nil {
        return fmt.Sprintf("Panic: %v", r) // Exposes internal state
    }
    return nil
}
```

### Error Logging vs Response Security
```go
// Correct: Detailed internal logging, sanitized external response
func handleError(err error, requestID string) {
    // Log full error details with correlation
    slog.Error("Handler error occurred",
        LogFieldRequestID, requestID,
        "error", err,
        "error_type", fmt.Sprintf("%T", err))
    
    // Return sanitized response
    response := ErrorResponse{
        Error:     mapToPublicError(err),
        Message:   sanitizeErrorMessage(err),
        RequestID: requestID,
        Timestamp: time.Now().UTC(),
    }
}

// Incorrect: Same level of detail in logs and responses
func unsafeErrorHandling(err error) {
    message := err.Error() // Same message for both log and response
    log.Printf("Error: %s", message)
    // Send same detailed message to client - security risk
}
```

---

## Implementation

### Step 1: Extend Core Error Functionality in Models Package

**File**: `internal/models/errors.go` (modify existing file)

Add core error classification and sanitization functionality to the existing errors file:

```go
package models

import "errors"

// Existing error variables
var (
	ErrInvalidID                = errors.New("invalid ID")
	ErrInvalidName              = errors.New("invalid name")
	ErrInvalidServiceID         = errors.New("invalid service ID")
	ErrInvalidParticipantStatus = errors.New("invalid participant status")
	ErrInvalidAlgorithm         = errors.New("invalid algorithm")
	ErrInvalidKeyStatus         = errors.New("invalid key status")
	ErrInvalidVersion           = errors.New("invalid version")
)

// Public error codes safe for external consumption
type ErrorCode string

const (
	ErrorCodeValidation     ErrorCode = "VALIDATION_ERROR"
	ErrorCodeAuthentication ErrorCode = "AUTHENTICATION_ERROR"
	ErrorCodeAuthorization  ErrorCode = "AUTHORIZATION_ERROR"
	ErrorCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrorCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrorCodeBadRequest     ErrorCode = "BAD_REQUEST"
)

// ErrorCategory represents error groupings for classification
type ErrorCategory string

const (
	ErrorCategoryValidation ErrorCategory = "validation"
	ErrorCategoryAuth       ErrorCategory = "authentication"
	ErrorCategoryNotFound   ErrorCategory = "not_found"
	ErrorCategoryInternal   ErrorCategory = "internal"
)

// ErrorClassification holds the public representation of an error
type ErrorClassification struct {
	Code     ErrorCode     `json:"code"`
	Category ErrorCategory `json:"category"`
}

// ClassifyError maps internal errors to public error classifications
// This function is framework-agnostic and can be used across packages
func ClassifyError(err error) ErrorClassification {
	switch {
	case err == ErrInvalidID:
		return ErrorClassification{Code: ErrorCodeValidation, Category: ErrorCategoryValidation}
	case err == ErrInvalidName:
		return ErrorClassification{Code: ErrorCodeValidation, Category: ErrorCategoryValidation}
	case err == ErrInvalidServiceID:
		return ErrorClassification{Code: ErrorCodeValidation, Category: ErrorCategoryValidation}
	case err == ErrInvalidParticipantStatus:
		return ErrorClassification{Code: ErrorCodeValidation, Category: ErrorCategoryValidation}
	case err == ErrInvalidAlgorithm:
		return ErrorClassification{Code: ErrorCodeValidation, Category: ErrorCategoryValidation}
	case err == ErrInvalidKeyStatus:
		return ErrorClassification{Code: ErrorCodeValidation, Category: ErrorCategoryValidation}
	case err == ErrInvalidVersion:
		return ErrorClassification{Code: ErrorCodeValidation, Category: ErrorCategoryValidation}
	default:
		return ErrorClassification{Code: ErrorCodeInternal, Category: ErrorCategoryInternal}
	}
}

// SanitizeErrorMessage provides user-safe error messages without internal details
// This prevents sensitive information leakage in external responses
func SanitizeErrorMessage(err error) string {
	switch {
	case err == ErrInvalidID:
		return "Invalid identifier provided"
	case err == ErrInvalidName:
		return "Invalid name provided"
	case err == ErrInvalidServiceID:
		return "Invalid service identifier provided"
	case err == ErrInvalidParticipantStatus:
		return "Invalid participant status provided"
	case err == ErrInvalidAlgorithm:
		return "Invalid algorithm specified"
	case err == ErrInvalidKeyStatus:
		return "Invalid key status provided"
	case err == ErrInvalidVersion:
		return "Invalid version number provided"
	default:
		return "An internal error occurred"
	}
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	classification := ClassifyError(err)
	return classification.Category == ErrorCategoryValidation
}

// IsInternalError checks if an error should be treated as an internal error
func IsInternalError(err error) bool {
	classification := ClassifyError(err)
	return classification.Category == ErrorCategoryInternal
}
```

### Step 2: Create Core Error Tests

**File**: `internal/models/errors_test.go` (modify existing file or create if needed)

Add tests for the new error classification functionality:

```go
package models

import (
	"fmt"
	"strings"
	"testing"
)

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name             string
		err              error
		expectedCode     ErrorCode
		expectedCategory ErrorCategory
	}{
		{
			name:             "invalid ID error",
			err:              ErrInvalidID,
			expectedCode:     ErrorCodeValidation,
			expectedCategory: ErrorCategoryValidation,
		},
		{
			name:             "invalid name error",
			err:              ErrInvalidName,
			expectedCode:     ErrorCodeValidation,
			expectedCategory: ErrorCategoryValidation,
		},
		{
			name:             "invalid service ID error",
			err:              ErrInvalidServiceID,
			expectedCode:     ErrorCodeValidation,
			expectedCategory: ErrorCategoryValidation,
		},
		{
			name:             "unknown error",
			err:              fmt.Errorf("unknown database error"),
			expectedCode:     ErrorCodeInternal,
			expectedCategory: ErrorCategoryInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClassifyError(tt.err)

			if result.Code != tt.expectedCode {
				t.Errorf("Code = %v, want %v", result.Code, tt.expectedCode)
			}

			if result.Category != tt.expectedCategory {
				t.Errorf("Category = %v, want %v", result.Category, tt.expectedCategory)
			}
		})
	}
}

func TestSanitizeErrorMessage(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		expectedMessage string
	}{
		{
			name:            "invalid ID error",
			err:             ErrInvalidID,
			expectedMessage: "Invalid identifier provided",
		},
		{
			name:            "invalid name error",
			err:             ErrInvalidName,
			expectedMessage: "Invalid name provided",
		},
		{
			name:            "internal error with sensitive details",
			err:             fmt.Errorf("database password authentication failed for user admin"),
			expectedMessage: "An internal error occurred",
		},
		{
			name:            "error with file path",
			err:             fmt.Errorf("failed to read /etc/passwd"),
			expectedMessage: "An internal error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeErrorMessage(tt.err)

			if result != tt.expectedMessage {
				t.Errorf("SanitizeErrorMessage() = %v, want %v", result, tt.expectedMessage)
			}

			// Verify no sensitive information is leaked
			if strings.Contains(result, "password") {
				t.Error("Sanitized message should not contain sensitive information like 'password'")
			}
			if strings.Contains(result, "/etc/") {
				t.Error("Sanitized message should not contain file paths")
			}
			if strings.Contains(result, "admin") {
				t.Error("Sanitized message should not contain usernames")
			}
		})
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "validation error",
			err:      ErrInvalidID,
			expected: true,
		},
		{
			name:     "another validation error",
			err:      ErrInvalidName,
			expected: true,
		},
		{
			name:     "internal error",
			err:      fmt.Errorf("database connection failed"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidationError(tt.err)
			if result != tt.expected {
				t.Errorf("IsValidationError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsInternalError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "validation error",
			err:      ErrInvalidID,
			expected: false,
		},
		{
			name:     "internal error",
			err:      fmt.Errorf("database connection failed"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInternalError(tt.err)
			if result != tt.expected {
				t.Errorf("IsInternalError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestErrorClassificationSecurity(t *testing.T) {
	// Test that sensitive information is never exposed
	sensitiveErrors := []error{
		fmt.Errorf("database password authentication failed for user admin"),
		fmt.Errorf("failed to read /etc/passwd"),
		fmt.Errorf("SQL injection detected: DROP TABLE users"),
		fmt.Errorf("internal file path: /home/app/secrets/api_keys.txt"),
	}

	for _, err := range sensitiveErrors {
		t.Run(err.Error(), func(t *testing.T) {
			// Test sanitized message
			message := SanitizeErrorMessage(err)
			if message != "An internal error occurred" {
				t.Errorf("Sensitive error should be sanitized, got: %v", message)
			}

			// Test classification
			classification := ClassifyError(err)
			if classification.Code != ErrorCodeInternal {
				t.Errorf("Sensitive error should be classified as internal, got: %v", classification.Code)
			}

			// Verify no sensitive patterns in sanitized message
			sensitivePatterns := []string{"password", "admin", "/etc/", "DROP TABLE", "/home/app/"}
			for _, pattern := range sensitivePatterns {
				if strings.Contains(message, pattern) {
					t.Errorf("Sanitized message should not contain sensitive pattern %q, got: %v", pattern, message)
				}
			}
		})
	}
}
```

### Step 3: Create HTTP Error Response Structures and Middleware

**File**: `internal/server/error_response.go` (new file)

```go
package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"cipher-hub/internal/models"
)

// Error response constants
const (
	// Error handling
	ErrorResponsePrefix = "ErrorResponse"

	// Log field names for consistency
	LogFieldErrorType     = "error_type"
	LogFieldErrorCode     = "error_code"
	LogFieldErrorCategory = "error_category"
	LogFieldPanicStack    = "panic_stack"
)

// ErrorResponse provides standardized JSON error response structure
type ErrorResponse struct {
	Error     models.ErrorClassification `json:"error"`
	Message   string                     `json:"message"`
	RequestID string                     `json:"request_id"`
	Timestamp time.Time                  `json:"timestamp"`
	Details   interface{}                `json:"details,omitempty"`
}

// ErrorResponseConfig holds configuration for error response middleware
type ErrorResponseConfig struct {
	IncludeDetails bool `json:"include_details"` // Include additional error details
	LogFullErrors  bool `json:"log_full_errors"` // Log complete error information
}

// NewErrorResponse creates a standardized error response with request correlation
// Uses core error classification from models package
func NewErrorResponse(err error, requestID string) *ErrorResponse {
	errorClassification := models.ClassifyError(err)
	message := models.SanitizeErrorMessage(err)

	return &ErrorResponse{
		Error:     errorClassification,
		Message:   message,
		RequestID: requestID,
		Timestamp: time.Now().UTC(),
	}
}

// NewErrorResponseWithDetails creates an error response with additional details
func NewErrorResponseWithDetails(err error, requestID string, details interface{}) *ErrorResponse {
	response := NewErrorResponse(err, requestID)
	response.Details = details
	return response
}

// mapErrorToStatusCode maps error classifications to appropriate HTTP status codes
func mapErrorToStatusCode(err error) int {
	if models.IsValidationError(err) {
		return http.StatusBadRequest
	}
	
	// Add more specific mappings as needed
	switch models.ClassifyError(err).Code {
	case models.ErrorCodeAuthentication:
		return http.StatusUnauthorized
	case models.ErrorCodeAuthorization:
		return http.StatusForbidden
	case models.ErrorCodeNotFound:
		return http.StatusNotFound
	case models.ErrorCodeValidation, models.ErrorCodeBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// writeErrorResponse writes a standardized JSON error response
func writeErrorResponse(w http.ResponseWriter, err error, requestID string) {
	errorResponse := NewErrorResponse(err, requestID)
	statusCode := mapErrorToStatusCode(err)

	// Set content type and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Encode and write JSON response
	if encodeErr := json.NewEncoder(w).Encode(errorResponse); encodeErr != nil {
		// Fallback to plain text if JSON encoding fails
		slog.Error("Failed to encode error response",
			LogFieldRequestID, requestID,
			"encode_error", encodeErr,
			"original_error", err)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal server error")
	}
}
```

### Step 4: Implement Error Response Middleware

**Continue in**: `internal/server/error_response.go`

```go
// ErrorResponseMiddleware creates middleware that handles error responses and panic recovery.
// Uses default configuration with comprehensive error logging enabled.
//
// Returns:
//   - Middleware: Configured error response middleware function
func ErrorResponseMiddleware() Middleware {
	config := ErrorResponseConfig{
		IncludeDetails: false, // Secure default - no details
		LogFullErrors:  true,  // Log full errors for debugging
	}
	return ErrorResponseMiddlewareWithConfig(config)
}

// ErrorResponseMiddlewareWithConfig creates middleware with custom error response configuration.
// Provides standardized JSON error responses, panic recovery, and request correlation.
//
// Features:
//   - Standardized JSON error response format with request correlation
//   - Panic recovery with structured error responses and stack trace logging
//   - Security-conscious error message sanitization using models package
//   - HTTP status code mapping based on error classifications
//   - Integration with request logging for error correlation and debugging
//   - Configurable error detail inclusion for different environments
//   - Comprehensive error logging with request context
//
// Security: Prevents sensitive information leakage by using sanitized error messages
// from the models package. All errors are logged with full details for internal debugging
// while returning safe, user-friendly messages to clients.
//
// Parameters:
//   - config: ErrorResponseConfig with error handling behavior settings
//
// Returns:
//   - Middleware: Configured error response middleware function
//
// Example usage:
//
//	config := ErrorResponseConfig{
//		IncludeDetails: false, // Production setting
//		LogFullErrors:  true,  // Enable full error logging
//	}
//	server.Middleware().Use(ErrorResponseMiddlewareWithConfig(config))
func ErrorResponseMiddlewareWithConfig(config ErrorResponseConfig) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get request ID for correlation
			requestID := GetRequestID(r.Context())

			// Panic recovery with structured error response
			defer func() {
				if rec := recover(); rec != nil {
					// Log panic with full details and stack trace
					if config.LogFullErrors {
						slog.Error("Panic recovered in error response middleware",
							LogFieldRequestID, requestID,
							"panic", rec,
							LogFieldPanicStack, string(debug.Stack()),
							"url", r.URL.String(),
							"method", r.Method)
					}

					// Create generic error for panic (uses models.ClassifyError internally)
					panicErr := fmt.Errorf("internal server error")
					writeErrorResponse(w, panicErr, requestID)
				}
			}()

			// Create response writer wrapper to capture status
			wrapped := newResponseWriter(w)

			// Call next handler
			next.ServeHTTP(wrapped, r)

			// Note: This middleware handles panics but doesn't intercept
			// explicit error responses from handlers. Handlers should use
			// HandleError() directly for explicit error handling.
		})
	}
}

// HandleError is a utility function for handlers to write standardized error responses
// This function should be used by handlers when they need to return error responses
//
// Uses core error classification and sanitization from models package
//
// Parameters:
//   - w: HTTP response writer
//   - r: HTTP request for context
//   - err: Error to handle and format
//
// Example usage in a handler:
//
//	func MyHandler(w http.ResponseWriter, r *http.Request) {
//		if err := someOperation(); err != nil {
//			HandleError(w, r, err)
//			return
//		}
//		// Success path...
//	}
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	requestID := GetRequestID(r.Context())
	classification := models.ClassifyError(err)

	// Log full error details for debugging
	slog.Error("Handler error occurred",
		LogFieldRequestID, requestID,
		"error", err,
		LogFieldErrorType, fmt.Sprintf("%T", err),
		LogFieldErrorCode, classification.Code,
		LogFieldErrorCategory, classification.Category,
		"url", r.URL.String(),
		"method", r.Method)

	// Write standardized error response
	writeErrorResponse(w, err, requestID)
}

// HandleErrorWithDetails is a utility function for handlers to write error responses with additional details
// Should only be used when additional context is safe to expose to clients
//
// Parameters:
//   - w: HTTP response writer
//   - r: HTTP request for context
//   - err: Error to handle and format
//   - details: Additional details safe for client consumption
func HandleErrorWithDetails(w http.ResponseWriter, r *http.Request, err error, details interface{}) {
	requestID := GetRequestID(r.Context())
	classification := models.ClassifyError(err)

	// Log full error details for debugging
	slog.Error("Handler error occurred with details",
		LogFieldRequestID, requestID,
		"error", err,
		LogFieldErrorType, fmt.Sprintf("%T", err),
		LogFieldErrorCode, classification.Code,
		LogFieldErrorCategory, classification.Category,
		"details", details,
		"url", r.URL.String(),
		"method", r.Method)

	// Create error response with details
	errorResponse := NewErrorResponseWithDetails(err, requestID, details)
	statusCode := mapErrorToStatusCode(err)

	// Set content type and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Encode and write JSON response
	if encodeErr := json.NewEncoder(w).Encode(errorResponse); encodeErr != nil {
		// Fallback to standard error response
		slog.Error("Failed to encode error response with details",
			LogFieldRequestID, requestID,
			"encode_error", encodeErr,
			"original_error", err)

		writeErrorResponse(w, err, requestID)
	}
}
```

### Step 5: Create HTTP Error Response Tests

**File**: `internal/server/error_response_test.go` (new file)

```go
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"cipher-hub/internal/models"
)

func TestNewErrorResponse(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		requestID       string
		expectedCode    models.ErrorCode
		expectedCategory models.ErrorCategory
		expectedMessage string
	}{
		{
			name:             "invalid ID error",
			err:              models.ErrInvalidID,
			requestID:        "test-request-123",
			expectedCode:     models.ErrorCodeValidation,
			expectedCategory: models.ErrorCategoryValidation,
			expectedMessage:  "Invalid identifier provided",
		},
		{
			name:             "invalid name error",
			err:              models.ErrInvalidName,
			requestID:        "test-request-456",
			expectedCode:     models.ErrorCodeValidation,
			expectedCategory: models.ErrorCategoryValidation,
			expectedMessage:  "Invalid name provided",
		},
		{
			name:             "internal error",
			err:              fmt.Errorf("database connection failed"),
			requestID:        "test-request-789",
			expectedCode:     models.ErrorCodeInternal,
			expectedCategory: models.ErrorCategoryInternal,
			expectedMessage:  "An internal error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := NewErrorResponse(tt.err, tt.requestID)

			if response.Error.Code != tt.expectedCode {
				t.Errorf("Error.Code = %v, want %v", response.Error.Code, tt.expectedCode)
			}

			if response.Error.Category != tt.expectedCategory {
				t.Errorf("Error.Category = %v, want %v", response.Error.Category, tt.expectedCategory)
			}

			if response.Message != tt.expectedMessage {
				t.Errorf("Message = %v, want %v", response.Message, tt.expectedMessage)
			}

			if response.RequestID != tt.requestID {
				t.Errorf("RequestID = %v, want %v", response.RequestID, tt.requestID)
			}

			// Check timestamp is recent (within last 5 seconds)
			if time.Since(response.Timestamp) > 5*time.Second {
				t.Errorf("Timestamp should be recent, got %v", response.Timestamp)
			}

			// Verify JSON serialization
			jsonData, err := json.Marshal(response)
			if err != nil {
				t.Errorf("Failed to marshal error response: %v", err)
			}

			var unmarshaled ErrorResponse
			if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
				t.Errorf("Failed to unmarshal error response: %v", err)
			}

			if unmarshaled.RequestID != tt.requestID {
				t.Errorf("Unmarshaled RequestID = %v, want %v", unmarshaled.RequestID, tt.requestID)
			}
		})
	}
}

func TestNewErrorResponseWithDetails(t *testing.T) {
	err := models.ErrInvalidID
	requestID := "test-request-123"
	details := map[string]string{"field": "user_id", "value": "invalid-format"}

	response := NewErrorResponseWithDetails(err, requestID, details)

	if response.Details == nil {
		t.Error("Details should not be nil")
	}

	// Verify details can be marshaled
	jsonData, err := json.Marshal(response)
	if err != nil {
		t.Errorf("Failed to marshal error response with details: %v", err)
	}

	var unmarshaled ErrorResponse
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Errorf("Failed to unmarshal error response with details: %v", err)
	}

	if unmarshaled.Details == nil {
		t.Error("Unmarshaled details should not be nil")
	}
}

func TestMapErrorToStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "validation error - invalid ID",
			err:            models.ErrInvalidID,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "validation error - invalid name",
			err:            models.ErrInvalidName,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "validation error - invalid algorithm",
			err:            models.ErrInvalidAlgorithm,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "internal error",
			err:            fmt.Errorf("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "unknown error",
			err:            fmt.Errorf("unexpected error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapErrorToStatusCode(tt.err)

			if result != tt.expectedStatus {
				t.Errorf("mapErrorToStatusCode() = %v, want %v", result, tt.expectedStatus)
			}
		})
	}
}

func TestWriteErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		requestID      string
		expectedStatus int
	}{
		{
			name:           "validation error",
			err:            models.ErrInvalidID,
			requestID:      "test-request-123",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "internal error",
			err:            fmt.Errorf("database error"),
			requestID:      "test-request-456",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			writeErrorResponse(w, tt.err, tt.requestID)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Status code = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Content-Type = %v, want application/json", contentType)
			}

			// Check JSON response structure
			var response ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Errorf("Failed to decode JSON response: %v", err)
			}

			if response.RequestID != tt.requestID {
				t.Errorf("Response RequestID = %v, want %v", response.RequestID, tt.requestID)
			}

			if response.Error.Code == "" {
				t.Error("Response should have error code")
			}

			if response.Message == "" {
				t.Error("Response should have error message")
			}
		})
	}
}

func TestErrorResponseMiddleware(t *testing.T) {
	middleware := ErrorResponseMiddleware()

	// Test normal operation (no error)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusOK)
	}

	if w.Body.String() != "success" {
		t.Errorf("Response body = %v, want success", w.Body.String())
	}
}

func TestErrorResponseMiddleware_PanicRecovery(t *testing.T) {
	middleware := ErrorResponseMiddleware()

	// Test panic recovery
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	wrappedHandler := middleware(handler)

	// Need request logging middleware to provide request ID
	fullMiddleware := RequestLoggingMiddleware()(wrappedHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Should not panic
	fullMiddleware.ServeHTTP(w, req)

	// Should return error response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusInternalServerError)
	}

	// Should be JSON response
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %v, want application/json", contentType)
	}

	// Should have structured error response
	var response ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode JSON response: %v", err)
	}

	if response.Error.Code != models.ErrorCodeInternal {
		t.Errorf("Error code = %v, want %v", response.Error.Code, models.ErrorCodeInternal)
	}

	if response.RequestID == "" {
		t.Error("Response should have request ID")
	}
}

func TestHandleError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedCode   models.ErrorCode
	}{
		{
			name:           "validation error",
			err:            models.ErrInvalidID,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   models.ErrorCodeValidation,
		},
		{
			name:           "internal error",
			err:            fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   models.ErrorCodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)

			// Add request ID to context
			ctx := WithRequestID(req.Context(), "test-request-123")
			req = req.WithContext(ctx)

			HandleError(w, req, tt.err)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Status code = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Check JSON response
			var response ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Errorf("Failed to decode JSON response: %v", err)
			}

			if response.Error.Code != tt.expectedCode {
				t.Errorf("Error code = %v, want %v", response.Error.Code, tt.expectedCode)
			}

			if response.RequestID != "test-request-123" {
				t.Errorf("Request ID = %v, want test-request-123", response.RequestID)
			}
		})
	}
}

func TestHandleErrorWithDetails(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)

	// Add request ID to context
	ctx := WithRequestID(req.Context(), "test-request-123")
	req = req.WithContext(ctx)

	err := models.ErrInvalidID
	details := map[string]string{"field": "user_id", "value": "invalid"}

	HandleErrorWithDetails(w, req, err, details)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Status code = %v, want %v", w.Code, http.StatusBadRequest)
	}

	// Check JSON response
	var response ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode JSON response: %v", err)
	}

	if response.Details == nil {
		t.Error("Response should include details")
	}

	if response.RequestID != "test-request-123" {
		t.Errorf("Request ID = %v, want test-request-123", response.RequestID)
	}
}

func TestErrorResponseIntegrationWithModels(t *testing.T) {
	// Test that server error responses correctly use models package functionality
	testErrors := []error{
		models.ErrInvalidID,
		models.ErrInvalidName,
		models.ErrInvalidServiceID,
		fmt.Errorf("some internal error"),
	}

	for _, err := range testErrors {
		t.Run(err.Error(), func(t *testing.T) {
			response := NewErrorResponse(err, "test-123")

			// Verify classification came from models package
			expectedClassification := models.ClassifyError(err)
			if response.Error.Code != expectedClassification.Code {
				t.Errorf("Error code = %v, want %v", response.Error.Code, expectedClassification.Code)
			}
			if response.Error.Category != expectedClassification.Category {
				t.Errorf("Error category = %v, want %v", response.Error.Category, expectedClassification.Category)
			}

			// Verify message came from models package
			expectedMessage := models.SanitizeErrorMessage(err)
			if response.Message != expectedMessage {
				t.Errorf("Message = %v, want %v", response.Message, expectedMessage)
			}
		})
	}
}
```

### Step 4: Add Server Integration Tests

**File**: `internal/server/server_test.go` (add these tests to existing file)

```go
func TestServer_ErrorResponseMiddleware(t *testing.T) {
	config := ServerConfig{
		Host: "localhost",
		Port: "0",
	}

	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("NewServer() unexpected error: %v", err)
	}

	// Add middleware stack with error response handling
	server.Middleware().
		Use(RequestLoggingMiddleware()).
		Use(ErrorResponseMiddleware())

	// Set test handler that returns an error
	server.SetHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a validation error
		HandleError(w, r, models.ErrInvalidID)
	}))

	// Start server
	err = server.Start()
	if err != nil {
		t.Fatalf("Start() unexpected error: %v", err)
	}
	defer func() {
		if err := server.Shutdown(); err != nil {
			t.Logf("Cleanup shutdown error: %v", err)
		}
	}()

	// Test error response integration
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Execute request through server
	server.httpServer.Handler.ServeHTTP(w, req)

	// Verify request ID header was added by request logging middleware
	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("X-Request-ID header not set by request logging middleware")
	}

	// Verify error response status
	if w.Code != http.StatusBadRequest {
		t.Errorf("Response status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	// Verify JSON error response
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %v, want application/json", contentType)
	}

	// Parse and verify error response structure
	var response ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode JSON error response: %v", err)
	}

	if response.RequestID != requestID {
		t.Errorf("Error response RequestID = %v, want %v", response.RequestID, requestID)
	}

	if response.Error.Code != ErrorCodeValidation {
		t.Errorf("Error code = %v, want %v", response.Error.Code, ErrorCodeValidation)
	}

	if response.Message != "Invalid identifier provided" {
		t.Errorf("Error message = %v, want 'Invalid identifier provided'", response.Message)
	}
}

func TestServer_ErrorResponseMiddleware_PanicRecovery(t *testing.T) {
	config := ServerConfig{
		Host: "localhost",
		Port: "0",
	}

	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("NewServer() unexpected error: %v", err)
	}

	// Add middleware stack with error response handling
	server.Middleware().
		Use(RequestLoggingMiddleware()).
		Use(ErrorResponseMiddleware())

	// Set test handler that panics
	server.SetHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic for error recovery")
	}))

	// Start server
	err = server.Start()
	if err != nil {
		t.Fatalf("Start() unexpected error: %v", err)
	}
	defer func() {
		if err := server.Shutdown(); err != nil {
			t.Logf("Cleanup shutdown error: %v", err)
		}
	}()

	// Test panic recovery
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Should not panic
	server.httpServer.Handler.ServeHTTP(w, req)

	// Verify request ID header was added
	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("X-Request-ID header should be present")
	}

	// Verify panic was recovered and returned as error response
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Response status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	// Verify JSON error response
	var response ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode JSON error response: %v", err)
	}

	if response.RequestID != requestID {
		t.Errorf("Error response RequestID = %v, want %v", response.RequestID, requestID)
	}

	if response.Error.Code != ErrorCodeInternal {
		t.Errorf("Error code = %v, want %v", response.Error.Code, ErrorCodeInternal)
	}

	// Verify panic details are not exposed in response
	if strings.Contains(response.Message, "test panic") {
		t.Error("Panic details should not be exposed in error response")
	}
}

func TestServer_ErrorResponseMiddleware_WithCORS(t *testing.T) {
	config := ServerConfig{
		Host: "localhost",
		Port: "0",
	}

	server, err := NewServer(config)
	if err != nil {
		t.Fatalf("NewServer() unexpected error: %v", err)
	}

	// Configure CORS
	corsConfig := CORSConfig{
		Enabled: true,
		Origins: []string{"http://localhost:3000"},
		Methods: DefaultCORSMethods,
		Headers: DefaultCORSHeaders,
	}

	// Add middleware stack with all middleware types
	server.Middleware().
		Use(RequestLoggingMiddleware()).
		UseIf(len(corsConfig.Origins) > 0, CORSMiddlewareWithConfig(corsConfig)).
		Use(ErrorResponseMiddleware())

	// Set test handler that returns an error
	server.SetHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleError(w, r, models.ErrInvalidName)
	}))

	// Start server
	err = server.Start()
	if err != nil {
		t.Fatalf("Start() unexpected error: %v", err)
	}
	defer func() {
		if err := server.Shutdown(); err != nil {
			t.Logf("Cleanup shutdown error: %v", err)
		}
	}()

	// Test error response with CORS
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	server.httpServer.Handler.ServeHTTP(w, req)

	// Verify all middleware applied correctly
	if w.Header().Get("X-Request-ID") == "" {
		t.Error("Request logging middleware should set request ID")
	}

	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Error("CORS middleware should set CORS headers")
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("Error response middleware should set status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Verify error response structure
	var response ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode JSON error response: %v", err)
	}

	if response.Error.Code != ErrorCodeValidation {
		t.Errorf("Error code = %v, want %v", response.Error.Code, ErrorCodeValidation)
	}
}
```

---

## Verification Steps

### Step 1: Build Verification
```bash
# Navigate to project root
cd cipher-hub/

# Verify clean build with error response middleware
go build ./...

# Expected: No compilation errors
```

### Step 2: Unit Test Verification
```bash
# Run core error classification tests in models package
go test ./internal/models -run "TestClassifyError\|TestSanitizeErrorMessage\|TestIsValidationError\|TestIsInternalError\|TestErrorClassificationSecurity" -v

# Run HTTP error response tests in server package
go test ./internal/server -run "TestErrorResponse\|TestNewErrorResponse\|TestMapErrorToStatusCode\|TestWriteErrorResponse\|TestHandleError" -v

# Expected: All error response and classification tests pass
```

### Step 3: Security Test Verification
```bash
# Run security-focused tests across both packages
go test ./internal/models -run "TestErrorClassificationSecurity" -v
go test ./internal/server -run "TestErrorResponseIntegrationWithModels" -v

# Expected: Security tests pass, no sensitive information exposed
```

### Step 4: Integration Test Verification
```bash
# Run server integration tests with error response middleware
go test ./internal/server -run "TestServer_ErrorResponse" -v

# Expected: All integration tests pass
```

### Step 5: Complete Test Suite Verification
```bash
# Run all server tests to ensure no regressions
go test ./internal/server -v

# Expected: All existing and new tests pass
```

### Step 6: Error Response Format Test
```bash
# Create test program to verify error response format and models integration
cat > test_error_format.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
    "cipher-hub/internal/models"
    "cipher-hub/internal/server"
)

func main() {
    // Test various error types
    errors := []error{
        models.ErrInvalidID,
        models.ErrInvalidName,
        models.ErrInvalidAlgorithm,
        fmt.Errorf("database connection failed"),
    }
    
    for _, err := range errors {
        // Test core classification from models package
        classification := models.ClassifyError(err)
        message := models.SanitizeErrorMessage(err)
        
        fmt.Printf("Error: %v\n", err)
        fmt.Printf("  Classification: %s/%s\n", classification.Code, classification.Category)
        fmt.Printf("  Sanitized Message: %s\n", message)
        
        // Test HTTP response from server package
        response := server.NewErrorResponse(err, "test-123")
        jsonData, _ := json.MarshalIndent(response, "", "  ")
        fmt.Printf("  HTTP Response:\n%s\n\n", jsonData)
    }
}
EOF

go run test_error_format.go
rm test_error_format.go

# Expected: Well-formatted JSON error responses using models package classification
```

### Step 7: Security Verification Test
```bash
# Create test to verify no sensitive information leakage across both packages
cat > test_error_security.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
    "strings"
    "cipher-hub/internal/models"
    "cipher-hub/internal/server"
)

func main() {
    // Test sensitive errors
    sensitiveErrors := []error{
        fmt.Errorf("database password authentication failed for user admin"),
        fmt.Errorf("failed to read /etc/passwd"),
        fmt.Errorf("SQL injection: DROP TABLE users"),
        fmt.Errorf("API key leaked: sk-1234567890abcdef"),
    }
    
    for _, err := range sensitiveErrors {
        // Test models package sanitization
        message := models.SanitizeErrorMessage(err)
        classification := models.ClassifyError(err)
        
        // Test server package HTTP response
        response := server.NewErrorResponse(err, "test-123")
        jsonData, _ := json.Marshal(response)
        jsonString := string(jsonData)
        
        // Check for sensitive patterns
        sensitivePatterns := []string{
            "password", "admin", "/etc/", "DROP TABLE", "sk-", "API key",
        }
        
        leaked := false
        for _, pattern := range sensitivePatterns {
            if strings.Contains(message, pattern) || strings.Contains(jsonString, pattern) {
                fmt.Printf("SECURITY ISSUE: Pattern %q found in response\n", pattern)
                leaked = true
            }
        }
        
        if !leaked {
            fmt.Printf("✓ Safe: %v -> %s (%s)\n", err, message, classification.Code)
        }
    }
}
EOF

go run test_error_security.go
rm test_error_security.go

# Expected: All errors should be sanitized, no sensitive patterns exposed
```

### Step 8: Code Quality Verification
```bash
# Format and lint checks
go fmt ./...
go vet ./...

# Expected: No issues reported
```