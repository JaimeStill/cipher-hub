package models

import (
	"testing"
)

func TestNewServiceRegistration(t *testing.T) {
	name := "Test Service"
	description := "A test service registration"

	sr, err := NewServiceRegistration(name, description)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if sr.ID == "" {
		t.Error("Expected generated ID, got empty string")
	}

	if sr.Name != name {
		t.Errorf("Expected name %q, got %q", name, sr.Name)
	}

	if sr.Description != description {
		t.Errorf("Expected description %q, got %q", description, sr.Description)
	}

	if sr.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set, got zero time")
	}

	if sr.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set, got zero time")
	}

	if sr.Participants == nil {
		t.Error("Expected Participants to be initialized, got nil")
	}

	if sr.Metadata == nil {
		t.Error("Expected Metadata to be initialized, got nil")
	}
}

func TestServiceRegistrationValidation(t *testing.T) {
	tests := []struct {
		name    string
		sr      *ServiceRegistration
		wantErr error
	}{
		{
			name: "valid service registration",
			sr: &ServiceRegistration{
				ID:   "test-id",
				Name: "Test Service",
			},
			wantErr: nil,
		},
		{
			name: "missing ID",
			sr: &ServiceRegistration{
				Name: "Test Service",
			},
			wantErr: ErrInvalidID,
		},
		{
			name: "missing name",
			sr: &ServiceRegistration{
				ID: "test-id",
			},
			wantErr: ErrInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sr.IsValid()
			if err != tt.wantErr {
				t.Errorf("ServiceRegistration.IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
