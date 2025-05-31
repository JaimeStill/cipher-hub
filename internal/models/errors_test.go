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
