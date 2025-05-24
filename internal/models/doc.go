// Package models defines the core data types for Cipher Hub key management service.
//
// This package provides the fundamental data structures for managing cryptographic
// keys across distributed systems: ServiceRegistration, Participant, and CryptoKey.
// All types implement security-first design patterns with comprehensive validation,
// secure serialization, and audit-friendly metadata.
//
// Key Security Features:
//   - Cryptographic key material never appears in serialized output
//   - Comprehensive input validation at construction and runtime
//   - Secure ID generation using crypto/rand
//   - Metadata-driven extensibility without compromising type safety
//
// Main Types:
//   - ServiceRegistration: Logical container for related participants
//   - Participant: Any entity (user/device/service) that can access keys
//   - CryptoKey: Cryptographic key with lifecycle management and metadata
//
// All constructors follow the pattern (Type, error) and include full validation.
// Runtime validation is available through IsValid() methods on all types.
package models
