package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cipher-hub/internal/models"
)

func TestNewErrorResponse(t *testing.T) {
	tests := []struct {
		name             string
		err              error
		requestID        string
		expectedCode     models.ErrorCode
		expectedCategory models.ErrorCategory
		expectedMessage  string
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
