// Package main provides the entry point for the Cipher Hub key management service.
//
// Cipher Hub is a containerized, security-first key management service designed
// to act as a centralized cryptographic layer for distributed systems. It provides
// secure key generation, storage, distribution, and lifecycle management through
// a RESTful HTTP API.
//
// The service is designed for sidecar deployment patterns within container
// orchestration platforms, handling all cryptographic operations for application
// services without requiring changes to application code.
//
// Key Capabilities:
//   - Secure cryptographic key generation and storage
//   - Service registration and participant management
//   - RESTful API for key operations with comprehensive authentication
//   - Container-native design with health checks and graceful shutdown
//   - Comprehensive audit logging for all key operations
//
// Signal Handling:
// The service handles SIGINT and SIGTERM signals for graceful shutdown:
//   - SIGINT (Ctrl+C): Initiates graceful shutdown with configured timeout
//   - SIGTERM: Container orchestration graceful shutdown signal
//   - Shutdown allows in-flight requests to complete within timeout bounds
//
// Usage:
//
//	cipher-hub [flags]
//
// The service reads configuration from environment variables and provides
// health check endpoints for container orchestration integration.
//
// Security Considerations:
// This service handles sensitive cryptographic material and should be deployed
// with appropriate security controls including TLS, authentication, and
// network isolation.
package main
