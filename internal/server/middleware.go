package server

import "net/http"

// Middleware defines the standard middleware function signature.
// A middleware function takes an http.Handler and returns an http.Handler,
// allowing for request / response processing before and after the wrapped handler.
//
// Example usage:
//
//	func LoggingMiddleware(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			log.Printf("Request: %s %s", r.Method, r.URL.Path)
//			next.ServeHTTP(w, r)
//		})
//	}
type Middleware func(http.Handler) http.Handler

// MiddlewareStack manages a collection of middleware functions with support
// for conditional applicatioin and ordered execution.
//
// Middleware is applied in the order it was added to the stack. The last
// middleware added will be the outermost middleware (executed first for requests,
// last for responses).
//
// Thread Safety: MiddlewareStack is safe for concurrent reads after setup
// is complete, but modifications (Use, UseIf) should only be performed
// during initialization phase before serving requests.
type MiddlewareStack struct {
	middlewares []Middleware
}

// NewMiddlewareStack creates a new empty middleware stack ready for use.
//
// Returns:
//   - *MiddlewareStack: Empty middleware stack ready for middleware registration
//
// Example:
//
//	stack := NewMiddlewareStack()
//	stack.
//		Use(RequestIDMiddleware()).
//		UseIf(config.EnableCORS, CORSMiddleware()).
//		Use(LoggingMiddleware())
func NewMiddlewareStack() *MiddlewareStack {
	return &MiddlewareStack{
		middlewares: make([]Middleware, 0),
	}
}

// Use adds a middleware function to the stack that will always be applied.
// Middleware is applied in registration order.
//
// Parameters:
//   - middleware: Middleware function to add to the stack
//
// Returns:
//   - *MiddlewareStack: The same stack instance for method chaining
//
// Example:
//
//	stack.Use(RequestIDMiddleware()).Use(LoggingMiddleware())
func (ms *MiddlewareStack) Use(middleware Middleware) *MiddlewareStack {
	ms.middlewares = append(ms.middlewares, middleware)
	return ms
}

// UseIf conditionally adds a middleware function to the stack based on the
// provided condition. If the condition is false, the middleware is not added.
//
// This is useful for environment-specific middleware or feature flags.
//
// Parameters:
//   - condition: Boolean condition determining whether to add the middleware
//   - middleware: Middleware function to add if condition is true
//
// Returns:
//   - *MiddlewareStack: The same stack instance for method chaining
//
// Example:
//
//	stack.
//		UseIf(config.EnableCORS, CORSMiddleware()).
//		UseIf(config.Environment == "development", DebugMiddleware())
func (ms *MiddlewareStack) UseIf(condition bool, middleware Middleware) *MiddlewareStack {
	if condition {
		ms.middlewares = append(ms.middlewares, middleware)
	}
	return ms
}

// Apply wraps the provided handler with all registered middleware functions.
// Middleware is applied in reverse order (last registered becomes outermost).
//
// This follows the standard middleware pattern where middleware closer to
// the registration point executes later in the request chain but earlier
// in the response chain.
//
// Performance: Middleware is applied once during server start for optimal
// runtime performance. The middleware chain is pre-built and reused for
// all requests, avoiding per-request overhead.
//
// Parameters:
//   - handler: The base handler to wrap the middleware
//
// Returns:
//   - http.Handler: Handler wrapped with all registered middleware
//
// Example:
//
//	finalHandler := stack.Apply(myBusinessLogicHandler)
//	http.ListenAndServe(":8080", finalHandler)
//
// Execution Flow Example:
//
//	stack.Use(A).Use(B).Use(C)
//	Request: C -> B -> A -> handler
//	Response: handler -> A -> B -> C
//
// This ensures middleware registered later can wrap and control middleware
// registered earlier, following standard middleware composition patterns.
func (ms *MiddlewareStack) Apply(handler http.Handler) http.Handler {
	if handler == nil {
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
	}

	result := handler

	// Apply middleware in reverse order for correct execution chain
	for i := len(ms.middlewares) - 1; i >= 0; i-- {
		result = ms.middlewares[i](result)
	}

	return result
}

// Count returns the number of middleware functions currently in the stack.
// This is useful for testing and debugging purposes.
//
// Returns:
//   - int: Number of middleware functions in the stack
func (ms *MiddlewareStack) Count() int {
	return len(ms.middlewares)
}

// Clear removes all middleware functions from the stack.
// This is primarily useful for testing scenarios.
func (ms *MiddlewareStack) Clear() {
	ms.middlewares = ms.middlewares[:0]
}
