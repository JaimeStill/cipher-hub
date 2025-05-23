package storage

import (
	"cipher-hub/internal/models"
	"context"
)

// Storage defines the interface for persisting and retrieving key management data
type Storage interface {
	// Service Registration operations
	CreateServiceRegistration(ctx context.Context, sr *models.ServiceRegistration) error
	GetServiceRegistration(ctx context.Context, id string) (*models.ServiceRegistration, error)
	UpdateServiceRegistration(ctx context.Context, sr *models.ServiceRegistration) error
	DeleteServiceRegistration(ctx context.Context, id string) error
	ListServiceRegistrations(ctx context.Context) ([]*models.ServiceRegistration, error)

	// Participant operations
	CreateParticipant(ctx context.Context, p *models.Participant) error
	GetParticipant(ctx context.Context, id string) (*models.Participant, error)
	UpdateParticipant(ctx context.Context, p *models.Participant) error
	DeleteParticipant(ctx context.Context, id string) error
	ListParticipantsByService(ctx context.Context, serviceID string) ([]*models.Participant, error)

	// Key operations
	CreateKey(ctx context.Context, key *models.CryptoKey) error
	GetKey(ctx context.Context, id string) (*models.CryptoKey, error)
	UpdateKey(ctx context.Context, key *models.CryptoKey) error
	DeleteKey(ctx context.Context, id string) error
	ListKeysByService(ctx context.Context, serviceID string) ([]*models.CryptoKey, error)
}
