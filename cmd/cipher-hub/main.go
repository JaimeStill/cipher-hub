package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cipher-hub/internal/server"
)

func main() {
	fmt.Println("Cipher Hub - Key Management Service")
	log.Println("Starting Cipher Hub...")

	// Create server configuration
	config := server.ServerConfig{
		Host: "localhost",
		Port: "8080",
	}

	// Create server instance
	srv, err := server.NewServer(config)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Create channel for shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start graceful shutdown handler (before server start to prevent races)
	go func() {
		sig := <-sigChan
		log.Printf("Received signal %v, initiating graceful shutdown...", sig)

		// Use server's configured timeout plus buffer for coordination
		shutdownTimeout := srv.ShutdownTimeout() + 5*time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// Channel to signal shutdown completion
		done := make(chan error, 1)

		// Perform shutdown in goroutine
		go func() {
			done <- srv.Shutdown()
		}()

		// Wait for shutdown completion or timeout
		select {
		case err := <-done:
			if err != nil {
				log.Printf("Server shutdown failed: %v", err)
				os.Exit(1)
			}
			log.Println("Graceful shutdown completed")
			os.Exit(0)
		case <-shutdownCtx.Done():
			log.Printf("Shutdown timeout (%v) exceeded, forcing exit", shutdownTimeout)
			os.Exit(1)
		}
	}()

	// Start the server
	log.Printf("Starting HTTP server on %s", srv.Address())
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Keep main goroutine alive
	log.Println("Server started successfully. Press Ctrl+C to stop.")

	// Block forever, shutdown handled by signal handler
	select {}
}
