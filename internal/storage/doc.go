// Package storage defines the persistence interface for Cipher Hub key management.
//
// This package provides the Storage interface that abstracts all persistence
// operations for ServiceRegistration, Participant, and CryptoKey entities.
// The interface is designed to support multiple backend implementations while
// maintaining consistent security and operational patterns.
//
// Key Design Principles:
//   - Context-aware operations for proper request lifecycle management
//   - Consistent error handling across all storage operations
//   - Security-conscious design with audit logging integration points
//   - Support for both single-entity and bulk operations
//
// Interface Design:
//   - All operations accept context.Context as first parameter
//   - CRUD operations provided for each core entity type
//   - List operations support service-scoped filtering
//   - Error handling follows Go best practices with proper error wrapping
//
// Backend implementations should ensure cryptographic key material is
// encrypted at rest and properly secured according to the organization's
// security policies.
package storage
