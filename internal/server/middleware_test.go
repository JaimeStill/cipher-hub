package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware_TypeDefinition(t *testing.T) {
	// Test that Middleware type can be used as expected
	var middleware Middleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "middleware")
			next.ServeHTTP(w, r)
		})
	}

	// Create a simple handler to wrap
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Apply middleware
	wrappedHandler := middleware(handler)

	// Test the wrapped handler
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(w, req)

	// Verify middleware was applied
	if w.Header().Get("X-Test") != "middleware" {
		t.Error("Middleware was not applied correctly")
	}

	if w.Body.String() != "test response" {
		t.Errorf("Handler response incorrect: got %q", w.Body.String())
	}
}

func TestNewMiddlewareStack(t *testing.T) {
	stack := NewMiddlewareStack()

	if stack == nil {
		t.Fatal("NewMiddlewareStack() returned nil")
	}

	if stack.Count() != 0 {
		t.Errorf("New middleware stack should be empty, got count %d", stack.Count())
	}

	if stack.middlewares == nil {
		t.Error("Middleware slice should be initialized")
	}
}

func TestMiddlewareStack_Use(t *testing.T) {
	stack := NewMiddlewareStack()

	// Create test middleware
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-1", "applied")
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-2", "applied")
			next.ServeHTTP(w, r)
		})
	}

	// Test method chaining
	result := stack.Use(middleware1).Use(middleware2)

	// Verify chaining returns same instance
	if result != stack {
		t.Error("Use() should return same instance for chaining")
	}

	// Verify middleware count
	if stack.Count() != 2 {
		t.Errorf("Expected 2 middleware, got %d", stack.Count())
	}
}

func TestMiddlewareStack_UseIf(t *testing.T) {
	tests := []struct {
		name          string
		condition     bool
		expectedCount int
		expectHeader  bool
	}{
		{
			name:          "condition true",
			condition:     true,
			expectedCount: 1,
			expectHeader:  true,
		},
		{
			name:          "condition false",
			condition:     false,
			expectedCount: 0,
			expectHeader:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := NewMiddlewareStack()

			middleware := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("X-Conditional", "applied")
					next.ServeHTTP(w, r)
				})
			}

			// Test conditional addition
			result := stack.UseIf(tt.condition, middleware)

			// Verify chaining
			if result != stack {
				t.Error("UseIf() should return same instance for chaining")
			}

			// Verify count
			if stack.Count() != tt.expectedCount {
				t.Errorf("Expected %d middleware, got %d", tt.expectedCount, stack.Count())
			}

			// Test application
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			finalHandler := stack.Apply(handler)
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			finalHandler.ServeHTTP(w, req)

			// Verify header presence
			hasHeader := w.Header().Get("X-Conditional") == "applied"
			if hasHeader != tt.expectHeader {
				t.Errorf("Expected header present: %v, got: %v", tt.expectHeader, hasHeader)
			}
		})
	}
}

func TestMiddlewareStack_Apply(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.Handler
		expectBody string
		expectCode int
	}{
		{
			name: "with valid handler",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("test response"))
			}),
			expectBody: "test response",
			expectCode: http.StatusOK,
		},
		{
			name:       "with nil handler",
			handler:    nil,
			expectBody: "404 page not found\n",
			expectCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := NewMiddlewareStack()

			// Add test middleware to verify application
			stack.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("X-Applied", "true")
					next.ServeHTTP(w, r)
				})
			})

			finalHandler := stack.Apply(tt.handler)

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			finalHandler.ServeHTTP(w, req)

			// Verify middleware was applied
			if w.Header().Get("X-Applied") != "true" {
				t.Error("Middleware was not applied")
			}

			// Verify response
			if w.Code != tt.expectCode {
				t.Errorf("Expected status %d, got %d", tt.expectCode, w.Code)
			}

			if w.Body.String() != tt.expectBody {
				t.Errorf("Expected body %q, got %q", tt.expectBody, w.Body.String())
			}
		})
	}
}

func TestMiddlewareStack_Apply_Order(t *testing.T) {
	stack := NewMiddlewareStack()

	var executionOrder []string

	// Add middleware that tracks execution order
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware1-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware1-after")
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware2-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware2-after")
		})
	}

	middleware3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware3-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware3-after")
		})
	}

	// Add middleware in order
	stack.Use(middleware1).Use(middleware2).Use(middleware3)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executionOrder = append(executionOrder, "handler")
		w.WriteHeader(http.StatusOK)
	})

	finalHandler := stack.Apply(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	finalHandler.ServeHTTP(w, req)

	// Verify execution order: first registered middleware executes first (conventional)
	expectedOrder := []string{
		"middleware1-before", // First registered, outermost
		"middleware2-before",
		"middleware3-before", // Last registered, innermost
		"handler",
		"middleware3-after", // Last registered, innermost
		"middleware2-after",
		"middleware1-after", // First registered, outermost
	}

	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d execution steps, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("Execution order[%d]: expected %q, got %q", i, expected, executionOrder[i])
		}
	}
}

func TestMiddlewareStack_Count(t *testing.T) {
	stack := NewMiddlewareStack()

	// Initially empty
	if stack.Count() != 0 {
		t.Errorf("New stack should have count 0, got %d", stack.Count())
	}

	// Add middleware and verify count
	testMiddleware := func(next http.Handler) http.Handler { return next }

	stack.Use(testMiddleware)
	if stack.Count() != 1 {
		t.Errorf("After one Use(), count should be 1, got %d", stack.Count())
	}

	stack.Use(testMiddleware)
	if stack.Count() != 2 {
		t.Errorf("After two Use(), count should be 2, got %d", stack.Count())
	}

	// Test UseIf with false condition
	stack.UseIf(false, testMiddleware)
	if stack.Count() != 2 {
		t.Errorf("After UseIf(false), count should remain 2, got %d", stack.Count())
	}

	// Test UseIf with true condition
	stack.UseIf(true, testMiddleware)
	if stack.Count() != 3 {
		t.Errorf("After UseIf(true), count should be 3, got %d", stack.Count())
	}
}

func TestMiddlewareStack_Clear(t *testing.T) {
	stack := NewMiddlewareStack()

	// Add some middleware
	testMiddleware := func(next http.Handler) http.Handler { return next }
	stack.Use(testMiddleware).Use(testMiddleware)

	if stack.Count() != 2 {
		t.Errorf("Expected count 2 before clear, got %d", stack.Count())
	}

	// Clear and verify
	stack.Clear()

	if stack.Count() != 0 {
		t.Errorf("Expected count 0 after clear, got %d", stack.Count())
	}

	// Verify we can still add middleware after clear
	stack.Use(testMiddleware)
	if stack.Count() != 1 {
		t.Errorf("Expected count 1 after clear and add, got %d", stack.Count())
	}
}

func TestMiddlewareStack_ChainedUsage(t *testing.T) {
	stack := NewMiddlewareStack()

	// Test complex chaining scenario
	result := stack.
		Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Middleware-1", "applied")
				next.ServeHTTP(w, r)
			})
		}).
		UseIf(true, func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Middleware-2", "applied")
				next.ServeHTTP(w, r)
			})
		}).
		UseIf(false, func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Middleware-3", "should-not-apply")
				next.ServeHTTP(w, r)
			})
		}).
		Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Middleware-4", "applied")
				next.ServeHTTP(w, r)
			})
		})

	// Verify chaining returns same instance
	if result != stack {
		t.Error("Chained methods should return same instance")
	}

	// Verify correct count (3 middleware, 1 skipped due to false condition)
	if stack.Count() != 3 {
		t.Errorf("Expected 3 middleware after chaining, got %d", stack.Count())
	}

	// Test application
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	finalHandler := stack.Apply(handler)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	finalHandler.ServeHTTP(w, req)

	// Verify applied middleware
	if w.Header().Get("X-Middleware-1") != "applied" {
		t.Error("Middleware 1 should be applied")
	}
	if w.Header().Get("X-Middleware-2") != "applied" {
		t.Error("Middleware 2 should be applied")
	}
	if w.Header().Get("X-Middleware-3") != "" {
		t.Error("Middleware 3 should not be applied (UseIf false)")
	}
	if w.Header().Get("X-Middleware-4") != "applied" {
		t.Error("Middleware 4 should be applied")
	}
}
