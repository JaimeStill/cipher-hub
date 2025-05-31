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
