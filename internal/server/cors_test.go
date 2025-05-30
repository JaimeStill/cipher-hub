package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCORSConfig_LoadFromEnv(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	originallySet := make(map[string]bool)
	envVars := []string{
		"CIPHER_HUB_CORS_ENABLED",
		"CIPHER_HUB_CORS_ORIGINS",
		"CIPHER_HUB_CORS_METHODS",
		"CIPHER_HUB_CORS_HEADERS",
		"CIPHER_HUB_CORS_MAX_AGE",
		"CIPHER_HUB_CORS_CREDENTIALS",
	}

	for _, key := range envVars {
		val, exists := os.LookupEnv(key)
		originalEnv[key] = val
		originallySet[key] = exists
	}

	// Clean up environment after test
	defer func() {
		for _, key := range envVars {
			if originallySet[key] {
				os.Setenv(key, originalEnv[key])
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	tests := []struct {
		name     string
		envVars  map[string]string
		expected CORSConfig
	}{
		{
			name:    "default values",
			envVars: map[string]string{},
			expected: CORSConfig{
				Enabled:     false,
				Origins:     nil,
				Methods:     DefaultCORSMethods,
				Headers:     DefaultCORSHeaders,
				MaxAge:      DefaultCORSMaxAge,
				Credentials: false,
			},
		},
		{
			name: "enabled with origins and custom config",
			envVars: map[string]string{
				"CIPHER_HUB_CORS_ENABLED":     "true",
				"CIPHER_HUB_CORS_ORIGINS":     "https://app.example.com,https://admin.example.com",
				"CIPHER_HUB_CORS_METHODS":     "GET, POST, PUT",
				"CIPHER_HUB_CORS_HEADERS":     "Content-Type, Authorization",
				"CIPHER_HUB_CORS_MAX_AGE":     "3600",
				"CIPHER_HUB_CORS_CREDENTIALS": "true",
			},
			expected: CORSConfig{
				Enabled:     true,
				Origins:     []string{"https://app.example.com", "https://admin.example.com"},
				Methods:     "GET, POST, PUT",
				Headers:     "Content-Type, Authorization",
				MaxAge:      "3600",
				Credentials: true,
			},
		},
		{
			name: "origins without explicit enabled",
			envVars: map[string]string{
				"CIPHER_HUB_CORS_ORIGINS": "https://app.example.com",
			},
			expected: CORSConfig{
				Enabled:     true, // Should be enabled automatically when origins are set
				Origins:     []string{"https://app.example.com"},
				Methods:     DefaultCORSMethods,
				Headers:     DefaultCORSHeaders,
				MaxAge:      DefaultCORSMaxAge,
				Credentials: false,
			},
		},
		{
			name: "flexible boolean values",
			envVars: map[string]string{
				"CIPHER_HUB_CORS_ENABLED":     "1",
				"CIPHER_HUB_CORS_ORIGINS":     "https://app.example.com",
				"CIPHER_HUB_CORS_CREDENTIALS": "yes",
			},
			expected: CORSConfig{
				Enabled:     true,
				Origins:     []string{"https://app.example.com"},
				Methods:     DefaultCORSMethods,
				Headers:     DefaultCORSHeaders,
				MaxAge:      DefaultCORSMaxAge,
				Credentials: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// CLEAR all environment variables first
			for _, key := range envVars {
				os.Unsetenv(key)
			}

			// Set environment variables for this test
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			config := CORSConfig{}
			config.ApplyDefaults()
			config.LoadFromEnv()

			if config.Enabled != tt.expected.Enabled {
				t.Errorf("Enabled = %v, want %v", config.Enabled, tt.expected.Enabled)
			}

			if len(config.Origins) != len(tt.expected.Origins) {
				t.Errorf("Origins length = %v, want %v", len(config.Origins), len(tt.expected.Origins))
			} else {
				for i, origin := range config.Origins {
					if origin != tt.expected.Origins[i] {
						t.Errorf("Origins[%d] = %v, want %v", i, origin, tt.expected.Origins[i])
					}
				}
			}

			if config.Methods != tt.expected.Methods {
				t.Errorf("Methods = %v, want %v", config.Methods, tt.expected.Methods)
			}

			if config.Headers != tt.expected.Headers {
				t.Errorf("Headers = %v, want %v", config.Headers, tt.expected.Headers)
			}

			if config.MaxAge != tt.expected.MaxAge {
				t.Errorf("MaxAge = %v, want %v", config.MaxAge, tt.expected.MaxAge)
			}

			if config.Credentials != tt.expected.Credentials {
				t.Errorf("Credentials = %v, want %v", config.Credentials, tt.expected.Credentials)
			}
		})
	}
}

func TestCORSConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  CORSConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid configuration",
			config: CORSConfig{
				Origins: []string{"https://app.example.com", "http://localhost:3000"},
			},
			wantErr: false,
		},
		{
			name: "wildcard origin (valid but warns)",
			config: CORSConfig{
				Origins: []string{"*"},
			},
			wantErr: false, // Valid but should warn
		},
		{
			name: "invalid origin URL",
			config: CORSConfig{
				Origins: []string{"not-a-valid-url"},
			},
			wantErr: true,
			errMsg:  "invalid origin URL",
		},
		{
			name: "mixed valid and invalid origins",
			config: CORSConfig{
				Origins: []string{"https://app.example.com", "invalid-url"},
			},
			wantErr: true,
			errMsg:  "invalid origin URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error, got nil")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestCORSConfig_IsOriginAllowed(t *testing.T) {
	tests := []struct {
		name     string
		config   CORSConfig
		origin   string
		expected bool
	}{
		{
			name: "no origins configured",
			config: CORSConfig{
				Origins: []string{},
			},
			origin:   "https://app.example.com",
			expected: false,
		},
		{
			name: "origin allowed (exact match)",
			config: CORSConfig{
				Origins: []string{"https://app.example.com", "http://localhost:3000"},
			},
			origin:   "https://app.example.com",
			expected: true,
		},
		{
			name: "origin not allowed",
			config: CORSConfig{
				Origins: []string{"https://app.example.com"},
			},
			origin:   "http://evil.com",
			expected: false,
		},
		{
			name: "wildcard origin",
			config: CORSConfig{
				Origins: []string{"*"},
			},
			origin:   "https://anywhere.com",
			expected: true,
		},
		{
			name: "case insensitive scheme matching",
			config: CORSConfig{
				Origins: []string{"https://app.example.com"},
			},
			origin:   "HTTPS://app.example.com",
			expected: true,
		},
		{
			name: "case insensitive host matching",
			config: CORSConfig{
				Origins: []string{"https://app.example.com"},
			},
			origin:   "https://APP.EXAMPLE.COM",
			expected: true,
		},
		{
			name: "case insensitive full URL matching",
			config: CORSConfig{
				Origins: []string{"http://localhost:3000"},
			},
			origin:   "HTTP://LOCALHOST:3000",
			expected: true,
		},
		{
			name: "different scheme should not match",
			config: CORSConfig{
				Origins: []string{"https://app.example.com"},
			},
			origin:   "http://app.example.com",
			expected: false,
		},
		{
			name: "different port should not match",
			config: CORSConfig{
				Origins: []string{"http://localhost:3000"},
			},
			origin:   "http://localhost:8080",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsOriginAllowed(tt.origin)
			if result != tt.expected {
				t.Errorf("IsOriginAllowed(%q) = %v, want %v", tt.origin, result, tt.expected)
			}
		})
	}
}

func TestNormalizeOrigin(t *testing.T) {
	tests := []struct {
		name     string
		origin   string
		expected string
	}{
		{
			name:     "lowercase scheme and host",
			origin:   "HTTPS://APP.EXAMPLE.COM",
			expected: "https://app.example.com",
		},
		{
			name:     "preserve port",
			origin:   "HTTP://LOCALHOST:3000",
			expected: "http://localhost:3000",
		},
		{
			name:     "preserve path case",
			origin:   "https://API.EXAMPLE.COM/API/v1",
			expected: "https://api.example.com/API/v1",
		},
		{
			name:     "already lowercase",
			origin:   "https://app.example.com",
			expected: "https://app.example.com",
		},
		{
			name:     "malformed URL fallback",
			origin:   "not-a-valid-url",
			expected: "not-a-valid-url", // Fallback to lowercase
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeOrigin(tt.origin)
			if result != tt.expected {
				t.Errorf("normalizeOrigin(%q) = %q, want %q", tt.origin, result, tt.expected)
			}
		})
	}
}

func TestCORSMiddleware(t *testing.T) {
	middleware := CORSMiddleware()

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Apply middleware
	wrappedHandler := middleware(handler)

	// Create test request with Origin header
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	w := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(w, req)

	// Verify response (no CORS headers since no origins configured by default)
	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("CORS headers should not be set without configured origins")
	}

	// Verify response body
	if w.Code != http.StatusOK {
		t.Errorf("Response status = %d, want %d", w.Code, http.StatusOK)
	}

	if w.Body.String() != "test response" {
		t.Errorf("Response body = %q, want %q", w.Body.String(), "test response")
	}
}

func TestCORSMiddlewareWithConfig(t *testing.T) {
	tests := []struct {
		name           string
		config         CORSConfig
		requestOrigin  string
		method         string
		expectCORS     bool
		expectedOrigin string
		expectedStatus int
	}{
		{
			name: "CORS disabled",
			config: CORSConfig{
				Enabled: false,
				Origins: []string{},
			},
			requestOrigin:  "http://localhost:3000",
			method:         "GET",
			expectCORS:     false,
			expectedOrigin: "",
			expectedStatus: http.StatusOK,
		},
		{
			name: "allowed origin GET request",
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
				Methods: DefaultCORSMethods,
				Headers: DefaultCORSHeaders,
			},
			requestOrigin:  "http://localhost:3000",
			method:         "GET",
			expectCORS:     true,
			expectedOrigin: "http://localhost:3000",
			expectedStatus: http.StatusOK,
		},
		{
			name: "disallowed origin GET request",
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
				Methods: DefaultCORSMethods,
				Headers: DefaultCORSHeaders,
			},
			requestOrigin:  "http://evil.com",
			method:         "GET",
			expectCORS:     false,
			expectedOrigin: "",
			expectedStatus: http.StatusOK,
		},
		{
			name: "allowed origin OPTIONS preflight",
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
				Methods: DefaultCORSMethods,
				Headers: DefaultCORSHeaders,
			},
			requestOrigin:  "http://localhost:3000",
			method:         "OPTIONS",
			expectCORS:     true,
			expectedOrigin: "http://localhost:3000",
			expectedStatus: http.StatusOK,
		},
		{
			name: "disallowed origin OPTIONS preflight",
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
				Methods: DefaultCORSMethods,
				Headers: DefaultCORSHeaders,
			},
			requestOrigin:  "http://evil.com",
			method:         "OPTIONS",
			expectCORS:     false,
			expectedOrigin: "",
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "wildcard origin",
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"*"},
				Methods: DefaultCORSMethods,
				Headers: DefaultCORSHeaders,
			},
			requestOrigin:  "http://anywhere.com",
			method:         "GET",
			expectCORS:     true,
			expectedOrigin: "http://anywhere.com",
			expectedStatus: http.StatusOK,
		},
		{
			name: "credentials enabled",
			config: CORSConfig{
				Enabled:     true,
				Origins:     []string{"http://localhost:3000"},
				Methods:     DefaultCORSMethods,
				Headers:     DefaultCORSHeaders,
				Credentials: true,
			},
			requestOrigin:  "http://localhost:3000",
			method:         "GET",
			expectCORS:     true,
			expectedOrigin: "http://localhost:3000",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := CORSMiddlewareWithConfig(tt.config)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("handler response"))
			})

			wrappedHandler := middleware(handler)

			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.requestOrigin != "" {
				req.Header.Set("Origin", tt.requestOrigin)
			}

			w := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(w, req)

			// Verify CORS headers
			corsOrigin := w.Header().Get("Access-Control-Allow-Origin")
			if tt.expectCORS {
				if corsOrigin != tt.expectedOrigin {
					t.Errorf("Access-Control-Allow-Origin = %q, want %q", corsOrigin, tt.expectedOrigin)
				}

				if methods := w.Header().Get("Access-Control-Allow-Methods"); methods == "" {
					t.Error("Access-Control-Allow-Methods should be set for allowed origins")
				}

				if headers := w.Header().Get("Access-Control-Allow-Headers"); headers == "" {
					t.Error("Access-Control-Allow-Headers should be set for allowed origins")
				}

				// Check credentials header if enabled
				if tt.config.Credentials {
					if creds := w.Header().Get("Access-Control-Allow-Credentials"); creds != "true" {
						t.Error("Access-Control-Allow-Credentials should be 'true' when credentials enabled")
					}
				}
			} else {
				if corsOrigin != "" {
					t.Errorf("Access-Control-Allow-Origin should not be set, got %q", corsOrigin)
				}
			}

			// Verify response status
			if w.Code != tt.expectedStatus {
				t.Errorf("Response status = %d, want %d", w.Code, tt.expectedStatus)
			}

			// Verify handler was called for non-OPTIONS or allowed OPTIONS requests
			if tt.method != "OPTIONS" || (tt.method == "OPTIONS" && tt.expectCORS) {
				if tt.method != "OPTIONS" && w.Body.String() != "handler response" {
					t.Errorf("Handler response = %q, want %q", w.Body.String(), "handler response")
				}
			}
		})
	}
}

func TestCORSMiddleware_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		setupReq  func() *http.Request
		config    CORSConfig
		expectErr bool
	}{
		{
			name: "no origin header",
			setupReq: func() *http.Request {
				return httptest.NewRequest("GET", "/test", nil)
			},
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
			},
			expectErr: false,
		},
		{
			name: "empty origin header",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Origin", "")
				return req
			},
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
			},
			expectErr: false,
		},
		{
			name: "malformed origin header",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Origin", "not-a-valid-url")
				return req
			},
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
			},
			expectErr: false,
		},
		{
			name: "case sensitive origin matching",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/test", nil)
				req.Header.Set("Origin", "HTTP://LOCALHOST:3000")
				return req
			},
			config: CORSConfig{
				Enabled: true,
				Origins: []string{"http://localhost:3000"},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := CORSMiddlewareWithConfig(tt.config)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("handler response"))
			})

			wrappedHandler := middleware(handler)
			req := tt.setupReq()
			w := httptest.NewRecorder()

			// Should not panic or fail regardless of request content
			wrappedHandler.ServeHTTP(w, req)

			if tt.expectErr {
				if w.Code < 400 {
					t.Errorf("Expected error response, got status %d", w.Code)
				}
			} else {
				if w.Code != http.StatusOK {
					t.Errorf("Expected success response, got status %d", w.Code)
				}
			}
		})
	}
}

func TestCORSMiddleware_RequestCorrelation(t *testing.T) {
	config := CORSConfig{
		Enabled: true,
		Origins: []string{"http://localhost:3000"},
		Methods: DefaultCORSMethods,
		Headers: DefaultCORSHeaders,
	}

	middleware := CORSMiddlewareWithConfig(config)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request ID is available for correlation
		requestID := GetRequestID(r.Context())
		if requestID == "" {
			t.Error("Request ID should be available for CORS correlation")
		}

		// Verify request ID format
		if len(requestID) != RequestIDHexLength {
			t.Errorf("Request ID length = %d, want %d", len(requestID), RequestIDHexLength)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("correlation test"))
	})

	// We need to wrap with request logging middleware to generate request ID
	wrappedHandler := RequestLoggingMiddleware()(middleware(handler))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)

	// Verify request ID was added by request logging middleware
	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("X-Request-ID header should be present from request logging middleware")
	}

	// Verify CORS headers were applied
	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Error("CORS headers should be applied for allowed origin")
	}

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Response status = %d, want %d", w.Code, http.StatusOK)
	}

	if w.Body.String() != "correlation test" {
		t.Errorf("Response body = %q, want %q", w.Body.String(), "correlation test")
	}
}
