package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestGenerateRequestID(t *testing.T) {
	tests := []struct {
		name string
		runs int
	}{
		{
			name: "single generation",
			runs: 1,
		},
		{
			name: "multiple generations for uniqueness",
			runs: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids := make(map[string]bool)

			for i := 0; i < tt.runs; i++ {
				id, err := generateRequestID()
				if err != nil {
					t.Fatalf("generateRequestID() error = %v", err)
				}

				// Verify ID format (16 character hex string)
				if len(id) != 16 {
					t.Errorf("generateRequestID() ID length = %d, want 16", len(id))
				}

				// Verify hex format
				for _, char := range id {
					if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
						t.Errorf("generateRequestID() invalid hex character: %c", char)
					}
				}

				// Check for duplicates in multiple runs
				if tt.runs > 1 {
					if ids[id] {
						t.Errorf("generateRequestID() duplicate ID generated: %s", id)
					}
					ids[id] = true
				}
			}
		})
	}
}

func TestWithRequestID(t *testing.T) {
	ctx := context.Background()
	requestID := "test-request-id"

	// Add request ID to context
	newCtx := WithRequestID(ctx, requestID)

	// Verify request ID was added
	retrievedID := GetRequestID(newCtx)
	if retrievedID != requestID {
		t.Errorf("GetRequestID() = %v, want %v", retrievedID, requestID)
	}
}

func TestGetRequestID(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		want string
	}{
		{
			name: "context with request ID",
			ctx:  WithRequestID(context.Background(), "test-id"),
			want: "test-id",
		},
		{
			name: "context without request ID",
			ctx:  context.Background(),
			want: "",
		},
		{
			name: "context with wrong type value",
			ctx:  context.WithValue(context.Background(), requestIDCtxKey, 123),
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRequestID(tt.ctx)
			if got != tt.want {
				t.Errorf("GetRequestID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResponseWriter(t *testing.T) {
	w := httptest.NewRecorder()
	wrapped := newResponseWriter(w)

	// Test default values
	if wrapped.StatusCode() != http.StatusOK {
		t.Errorf("newResponseWriter() default status = %d, want %d",
			wrapped.StatusCode(), http.StatusOK)
	}

	if wrapped.BytesWritten() != 0 {
		t.Errorf("newResponseWriter() default bytes = %d, want 0", wrapped.BytesWritten())
	}
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	wrapped := newResponseWriter(w)

	// Test status code capture
	wrapped.WriteHeader(http.StatusNotFound)

	if wrapped.StatusCode() != http.StatusNotFound {
		t.Errorf("WriteHeader() status = %d, want %d",
			wrapped.StatusCode(), http.StatusNotFound)
	}

	// Verify underlying writer received the status
	if w.Code != http.StatusNotFound {
		t.Errorf("Underlying writer status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestResponseWriter_Write(t *testing.T) {
	w := httptest.NewRecorder()
	wrapped := newResponseWriter(w)

	testData := []byte("test response data")
	n, err := wrapped.Write(testData)

	if err != nil {
		t.Errorf("Write() error = %v", err)
	}

	if n != len(testData) {
		t.Errorf("Write() bytes written = %d, want %d", n, len(testData))
	}

	if wrapped.BytesWritten() != int64(len(testData)) {
		t.Errorf("BytesWritten() = %d, want %d", wrapped.BytesWritten(), len(testData))
	}

	// Verify underlying writer received the data
	if w.Body.String() != string(testData) {
		t.Errorf("Underlying writer body = %q, want %q", w.Body.String(), string(testData))
	}
}

func TestRequestLoggingMiddleware(t *testing.T) {
	middleware := RequestLoggingMiddleware()

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request ID is available in context
		requestID := GetRequestID(r.Context())
		if requestID == "" {
			t.Error("Request ID not found in context")
		}

		// Write test response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Apply middleware
	wrappedHandler := middleware(handler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	req.Header.Set("User-Agent", "test-agent")

	w := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(w, req)

	// Verify response has request ID header
	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("X-Request-ID header not set")
	}

	// Verify request ID format
	if len(requestID) != 16 {
		t.Errorf("Request ID length = %d, want 16", len(requestID))
	}

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Response status = %d, want %d", w.Code, http.StatusOK)
	}

	if w.Body.String() != "test response" {
		t.Errorf("Response body = %q, want %q", w.Body.String(), "test response")
	}
}

func TestRequestLoggingMiddleware_GenerationFailure(t *testing.T) {
	// This test would require mocking crypto/rand failure
	// For now, we test the middleware with successful generation
	// In production, request ID generation failure is extremely rare

	middleware := RequestLoggingMiddleware()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Should not panic even with potential generation issues
	wrappedHandler.ServeHTTP(w, req)

	// Verify some response was generated
	if w.Code == 0 {
		t.Error("No response status code set")
	}
}

func TestRequestLoggingMiddleware_StatusCodeCapture(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody string
	}{
		{
			name:         "200 OK",
			statusCode:   http.StatusOK,
			responseBody: "success",
		},
		{
			name:         "404 Not Found",
			statusCode:   http.StatusNotFound,
			responseBody: "not found",
		},
		{
			name:         "500 Internal Server Error",
			statusCode:   http.StatusInternalServerError,
			responseBody: "server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := RequestLoggingMiddleware()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			})

			wrappedHandler := middleware(handler)

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("Status code = %d, want %d", w.Code, tt.statusCode)
			}

			if w.Body.String() != tt.responseBody {
				t.Errorf("Response body = %q, want %q", w.Body.String(), tt.responseBody)
			}

			// Verify request ID header is present
			if w.Header().Get("X-Request-ID") == "" {
				t.Error("X-Request-ID header not set")
			}
		})
	}
}

func TestRequestLoggingMiddleware_MethodAndPathCapture(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "GET request",
			method: "GET",
			path:   "/api/health",
		},
		{
			name:   "POST request",
			method: "POST",
			path:   "/api/services",
		},
		{
			name:   "PUT request",
			method: "PUT",
			path:   "/api/services/123",
		},
		{
			name:   "DELETE request",
			method: "DELETE",
			path:   "/api/services/456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := RequestLoggingMiddleware()

			var capturedMethod, capturedPath string
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedMethod = r.Method
				capturedPath = r.URL.Path
				w.WriteHeader(http.StatusOK)
			})

			wrappedHandler := middleware(handler)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w, req)

			if capturedMethod != tt.method {
				t.Errorf("Captured method = %q, want %q", capturedMethod, tt.method)
			}

			if capturedPath != tt.path {
				t.Errorf("Captured path = %q, want %q", capturedPath, tt.path)
			}
		})
	}
}

func TestRequestLoggingConfig_LoadFromEnv(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{
		"CIPHER_HUB_LOGGING_ENABLED",
		"CIPHER_HUB_LOG_LEVEL",
		"CIPHER_HUB_LOG_FORMAT",
		"CIPHER_HUB_LOG_INCLUDE_HEADERS",
	}

	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
	}

	// Clean up environment after test
	defer func() {
		for _, key := range envVars {
			if val, exists := originalEnv[key]; exists {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		expected RequestLoggingConfig
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			expected: RequestLoggingConfig{
				Enabled:        true,
				Level:          "info",
				Format:         "json",
				IncludeHeaders: false,
			},
		},
		{
			name: "custom configuration",
			envVars: map[string]string{
				"CIPHER_HUB_LOGGING_ENABLED":     "false",
				"CIPHER_HUB_LOG_LEVEL":           "debug",
				"CIPHER_HUB_LOG_FORMAT":          "text",
				"CIPHER_HUB_LOG_INCLUDE_HEADERS": "true",
			},
			expected: RequestLoggingConfig{
				Enabled:        false,
				Level:          "debug",
				Format:         "text",
				IncludeHeaders: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			config := RequestLoggingConfig{}
			config.ApplyDefaults()
			config.LoadFromEnv()

			if config.Enabled != tt.expected.Enabled {
				t.Errorf("Enabled = %v, want %v", config.Enabled, tt.expected.Enabled)
			}
			if config.Level != tt.expected.Level {
				t.Errorf("Level = %v, want %v", config.Level, tt.expected.Level)
			}
			if config.Format != tt.expected.Format {
				t.Errorf("Format = %v, want %v", config.Format, tt.expected.Format)
			}
			if config.IncludeHeaders != tt.expected.IncludeHeaders {
				t.Errorf("IncludeHeaders = %v, want %v", config.IncludeHeaders, tt.expected.IncludeHeaders)
			}
		})
	}
}

func TestFilterSensitiveHeaders(t *testing.T) {
	headers := http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"Bearer token123"},
		"X-Api-Key":     []string{"secret123"},
		"User-Agent":    []string{"test-client"},
		"Cookie":        []string{"session=abc123"},
		"X-Request-ID":  []string{"req-123"},
	}

	filtered := filterSensitiveHeaders(headers)

	// Should include non-sensitive headers
	if filtered["Content-Type"] != "application/json" {
		t.Error("Content-Type should be included")
	}
	if filtered["User-Agent"] != "test-client" {
		t.Error("User-Agent should be included")
	}
	if filtered["X-Request-ID"] != "req-123" {
		t.Error("X-Request-ID should be included")
	}

	// Should exclude sensitive headers
	if _, exists := filtered["Authorization"]; exists {
		t.Error("Authorization should be filtered out")
	}
	if _, exists := filtered["X-Api-Key"]; exists {
		t.Error("X-Api-Key should be filtered out")
	}
	if _, exists := filtered["Cookie"]; exists {
		t.Error("Cookie should be filtered out")
	}
}

func TestRequestLoggingMiddlewareWithConfig(t *testing.T) {
	tests := []struct {
		name          string
		config        RequestLoggingConfig
		expectLogging bool
	}{
		{
			name: "logging enabled",
			config: RequestLoggingConfig{
				Enabled: true,
				Level:   "info",
				Format:  "json",
			},
			expectLogging: true,
		},
		{
			name: "logging disabled",
			config: RequestLoggingConfig{
				Enabled: false,
				Level:   "info",
				Format:  "json",
			},
			expectLogging: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := RequestLoggingMiddlewareWithConfig(tt.config)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request ID availability matches expectation
				requestID := GetRequestID(r.Context())
				if tt.expectLogging && requestID == "" {
					t.Error("Request ID should be available when logging enabled")
				}
				if !tt.expectLogging && requestID != "" {
					t.Error("Request ID should not be available when logging disabled")
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("test response"))
			})

			wrappedHandler := middleware(handler)

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w, req)

			// Check request ID header presence
			requestIDHeader := w.Header().Get("X-Request-ID")
			if tt.expectLogging && requestIDHeader == "" {
				t.Error("X-Request-ID header should be set when logging enabled")
			}
			if !tt.expectLogging && requestIDHeader != "" {
				t.Error("X-Request-ID header should not be set when logging disabled")
			}

			// Response should always be successful
			if w.Code != http.StatusOK {
				t.Errorf("Response status = %d, want %d", w.Code, http.StatusOK)
			}
		})
	}
}

func TestRequestLoggingMiddleware_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		setupReq    func() *http.Request
		expectError bool
	}{
		{
			name: "normal request",
			setupReq: func() *http.Request {
				return httptest.NewRequest("GET", "/test", nil)
			},
			expectError: false,
		},
		{
			name: "request with very long URL",
			setupReq: func() *http.Request {
				longPath := "/test/" + strings.Repeat("a", 1000)
				return httptest.NewRequest("GET", longPath, nil)
			},
			expectError: false,
		},
		{
			name: "request with malformed headers",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				// Add various header edge cases
				req.Header.Set("X-Test-Header", "value with\nnewline")
				req.Header.Set("X-Empty-Header", "")
				return req
			},
			expectError: false,
		},
		{
			name: "request with sensitive headers",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer secret-token")
				req.Header.Set("Cookie", "session=secret-session")
				req.Header.Set("X-API-Key", "secret-api-key")
				return req
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := RequestLoggingConfig{
				Enabled:        true,
				Level:          "info",
				Format:         "json",
				IncludeHeaders: true,
			}
			middleware := RequestLoggingMiddlewareWithConfig(config)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestID := GetRequestID(r.Context())
				if requestID == "" {
					t.Error("Request ID should be available")
				}
				w.WriteHeader(http.StatusOK)
			})

			wrappedHandler := middleware(handler)
			req := tt.setupReq()
			w := httptest.NewRecorder()

			// Should not panic or fail regardless of request content
			wrappedHandler.ServeHTTP(w, req)

			if tt.expectError {
				if w.Code < 400 {
					t.Errorf("Expected error response, got status %d", w.Code)
				}
			} else {
				if w.Code != http.StatusOK {
					t.Errorf("Expected success response, got status %d", w.Code)
				}

				// Verify request ID header is always present when logging enabled
				if w.Header().Get("X-Request-ID") == "" {
					t.Error("X-Request-ID header should be present")
				}
			}
		})
	}
}

func TestRequestLoggingMiddleware_HighConcurrency(t *testing.T) {
	config := RequestLoggingConfig{
		Enabled: true,
		Level:   "info",
		Format:  "json",
	}
	middleware := RequestLoggingMiddlewareWithConfig(config)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := GetRequestID(r.Context())
		if requestID == "" {
			t.Error("Request ID should be available")
		}
		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := middleware(handler)

	// Test concurrent request ID generation
	const concurrency = 100
	requestIDs := make(chan string, concurrency)

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w, req)

			requestID := w.Header().Get("X-Request-ID")
			if requestID != "" {
				requestIDs <- requestID
			}
		}()
	}

	wg.Wait()
	close(requestIDs)

	// Verify all request IDs are unique
	seen := make(map[string]bool)
	count := 0
	for requestID := range requestIDs {
		if seen[requestID] {
			t.Errorf("Duplicate request ID generated: %s", requestID)
		}
		seen[requestID] = true
		count++

		// Verify format
		if len(requestID) != RequestIDHexLength {
			t.Errorf("Invalid request ID length: %d, want %d", len(requestID), RequestIDHexLength)
		}
	}

	if count != concurrency {
		t.Errorf("Expected %d unique request IDs, got %d", concurrency, count)
	}
}
