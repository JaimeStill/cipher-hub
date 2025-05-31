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
	Details   any                        `json:"details,omitempty"`
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
func NewErrorResponseWithDetails(err error, requestID string, details any) *ErrorResponse {
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
// # Uses core error classification and sanitization from models package
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
func HandleErrorWithDetails(w http.ResponseWriter, r *http.Request, err error, details any) {
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
