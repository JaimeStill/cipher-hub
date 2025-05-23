package models

import (
	"fmt"
	"time"
)

// ParticipantStatus defines the status of a participant
type ParticipantStatus string

const (
	ParticipantStatusActive    ParticipantStatus = "active"
	ParticipantStatusInactive  ParticipantStatus = "inactive"
	ParticipantStatusSuspended ParticipantStatus = "suspended"
)

// Participant represents a user, device, or service that can access keys
// within a service registration context
type Participant struct {
	ID             string            `json:"id"`
	ServiceID      string            `json:"service_id"`
	Name           string            `json:"name"`
	Status         ParticipantStatus `json:"status"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	LastAccessedAt *time.Time        `json:"last_accessed_at,omitempty"`
	Metadata       map[string]any    `json:"metadata,omitempty"`
}

// NewParticipant creates a new participant with a generated ID,
// Returns an error if ID generation fails
func NewParticipant(serviceID, name string) (*Participant, error) {
	id, err := generateID()
	if err != nil {
		return nil, fmt.Errorf("failed to create participant: %w", err)
	}

	now := time.Now()

	return &Participant{
		ID:        id,
		ServiceID: serviceID,
		Name:      name,
		Status:    ParticipantStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
		Metadata:  make(map[string]any),
	}, nil
}

// IsValid validates a Participant
func (p *Participant) IsValid() error {
	if p.ID == "" {
		return ErrInvalidID
	}
	if p.ServiceID == "" {
		return ErrInvalidServiceID
	}
	if p.Name == "" {
		return ErrInvalidName
	}
	if !p.Status.IsValid() {
		return ErrInvalidParticipantStatus
	}
	return nil
}

// IsValid validates a ParticipantStatus
func (ps ParticipantStatus) IsValid() bool {
	switch ps {
	case ParticipantStatusActive, ParticipantStatusInactive, ParticipantStatusSuspended:
		return true
	default:
		return false
	}
}
