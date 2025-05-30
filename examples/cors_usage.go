package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"cipher-hub/internal/server"
)

func main() {
	// Configure structured logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	config := server.ServerConfig{
		Host: "localhost",
		Port: "8080",
	}

	srv, err := server.NewServer(config)
	if err != nil {
		panic(err)
	}

	// Configure CORS with specific origins
	corsConfig := server.CORSConfig{
		Enabled: true,
		Origins: []string{"http://localhost:3000", "https://app.example.com"},
	}

	// Configure middleware with conditional CORS
	srv.Middleware().
		Use(server.RequestLoggingMiddleware()).                                         // Always log requests
		UseIf(len(corsConfig.Origins) > 0, server.CORSMiddlewareWithConfig(corsConfig)) // CORS when origins configured

	// Set handler
	srv.SetHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := server.GetRequestID(r.Context())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hello! Request ID: %s", requestID)))
	}))

	fmt.Println("CORS middleware integration example compiled successfully")
	fmt.Printf("Middleware count: %d\n", srv.Middleware().Count())
	fmt.Printf("CORS origins configured: %v\n", corsConfig.Origins)
}
