package models

import (
	"fmt"
	"time"
)

// Algorithm defines the cryptographic algorithm types supported
type Algorithm string

const (
	AlgorithmAES256 Algorithm = "aes-256"
	// Future algorithms can be added here:
	// AlgorithmChaCha20Poly1305 Algorithm = "chacha20-poly1305"
	// AlgorithmRSA2048 Algorithm = "rsa-2048"
	// AlgorithmECDSP256 Algorithm = "ecdsa-p256"
)

// KeyStatus defines the lifecycle status of a key
type KeyStatus string

const (
	KeyStatusActive     KeyStatus = "active"
	KeyStatusRotating   KeyStatus = "rotating"
	KeyStatusDeprecated KeyStatus = "deprecated"
	KeyStatusRevoked    KeyStatus = "revoked"
)

// RotationInfo contains information about key rotation schedules and history
type RotationInfo struct {
	ScheduleEnabled bool          `json:"schedule_enabled"`
	RotationPeriod  time.Duration `json:"rotation_period,omitempty"`
	NextRotation    *time.Time    `json:"next_rotation,omitempty"`
	LastRotation    *time.Time    `json:"last_rotation,omitempty"`
	PreviousKeyID   string        `json:"previous_key_id,omitempty"`
}

// CryptoKey represents a cryptographic key with its metadata and lifecycle information
type CryptoKey struct {
	ID           string         `json:"id"`
	ServiceID    string         `json:"service_id"`
	Name         string         `json:"name"`
	Algorithm    Algorithm      `json:"algorithm"`
	KeyData      []byte         `json:"-"` // never serialize key material
	Version      int            `json:"version"`
	Status       KeyStatus      `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	ExpiresAt    *time.Time     `json:"expires_at,omitempty"`
	RotationInfo *RotationInfo  `json:"rotation_info,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

func NewCryptoKey(serviceID, name string, algorithm Algorithm) (*CryptoKey, error) {
	id, err := generateID()
	if err != nil {
		return nil, fmt.Errorf("failed to create crypto key: %w", err)
	}

	now := time.Now()
	return &CryptoKey{
		ID:        id,
		ServiceID: serviceID,
		Name:      name,
		Algorithm: algorithm,
		Version:   1,
		Status:    KeyStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
		Metadata:  make(map[string]any),
	}, nil
}

// IsValid validates a CryptoKey
func (ck *CryptoKey) IsValid() error {
	if ck.ID == "" {
		return ErrInvalidID
	}
	if ck.ServiceID == "" {
		return ErrInvalidServiceID
	}
	if ck.Name == "" {
		return ErrInvalidName
	}
	if !ck.Algorithm.IsValid() {
		return ErrInvalidAlgorithm
	}
	if !ck.Status.IsValid() {
		return ErrInvalidKeyStatus
	}
	if ck.Version < 1 {
		return ErrInvalidVersion
	}
	return nil
}

// IsValid validates an Algorithm
func (a Algorithm) IsValid() bool {
	switch a {
	case AlgorithmAES256:
		return true
	default:
		return false
	}
}

// IsValid validates a KeyStatus
func (ks KeyStatus) IsValid() bool {
	switch ks {
	case KeyStatusActive, KeyStatusRotating, KeyStatusDeprecated, KeyStatusRevoked:
		return true
	default:
		return false
	}
}
