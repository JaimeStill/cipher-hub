package models

import (
	"fmt"
	"time"
)

// ServiceRegistration represents a logical container for related participants
// that share cryptographci contexts and security boundaries
type ServiceRegistration struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Participants []Participant  `json:"participants,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

// NewServiceRegistration creates a new service registration with a generated ID
func NewServiceRegistration(name, description string) (*ServiceRegistration, error) {
	id, err := generateID()
	if err != nil {
		return nil, fmt.Errorf("failed to create service registration: %w", err)
	}

	now := time.Now()
	return &ServiceRegistration{
		ID:           id,
		Name:         name,
		Description:  description,
		CreatedAt:    now,
		UpdatedAt:    now,
		Participants: make([]Participant, 0),
		Metadata:     make(map[string]any),
	}, nil
}

// IsValid validates a ServiceRegistration
func (sr *ServiceRegistration) IsValid() error {
	if sr.ID == "" {
		return ErrInvalidID
	}
	if sr.Name == "" {
		return ErrInvalidName
	}
	return nil
}
