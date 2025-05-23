package models

import (
	"testing"
)

func TestNewParticipant(t *testing.T) {
	serviceID := "service-123"
	name := "Test Participant"

	p, err := NewParticipant(serviceID, name)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if p.ID == "" {
		t.Error("Expected generated ID, got empty string")
	}

	if p.ServiceID != serviceID {
		t.Errorf("Expected ServiceID %q, got %q", serviceID, p.ServiceID)
	}

	if p.Name != name {
		t.Errorf("Expected name %q, got %q", name, p.Name)
	}

	if p.Status != ParticipantStatusActive {
		t.Errorf("Expected status %q, got %q", ParticipantStatusActive, p.Status)
	}

	if p.Metadata == nil {
		t.Error("Expected Metadata to be initialized, got nil")
	}
}

func TestParticipantValidation(t *testing.T) {
	tests := []struct {
		name    string
		p       *Participant
		wantErr error
	}{
		{
			name: "valid participant",
			p: &Participant{
				ID:        "test-id",
				ServiceID: "service-id",
				Name:      "Test Participant",
				Status:    ParticipantStatusActive,
			},
			wantErr: nil,
		},
		{
			name: "missing ID",
			p: &Participant{
				ServiceID: "service-id",
				Name:      "Test Participant",
				Status:    ParticipantStatusActive,
			},
			wantErr: ErrInvalidID,
		},
		{
			name: "missing service ID",
			p: &Participant{
				ID:     "test-id",
				Name:   "Test Participant",
				Status: ParticipantStatusActive,
			},
			wantErr: ErrInvalidServiceID,
		},
		{
			name: "missing name",
			p: &Participant{
				ID:        "test-id",
				ServiceID: "service-id",
				Status:    ParticipantStatusActive,
			},
			wantErr: ErrInvalidName,
		},
		{
			name: "invalid status",
			p: &Participant{
				ID:        "test-id",
				ServiceID: "service-id",
				Name:      "Test Participant",
				Status:    "invalid",
			},
			wantErr: ErrInvalidParticipantStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.p.IsValid()
			if err != tt.wantErr {
				t.Errorf("Participant.IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParticipantMetadata(t *testing.T) {
	serviceID := "service-123"
	name := "Test Participant"

	p, err := NewParticipant(serviceID, name)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Test metadata functionality
	p.Metadata["type"] = "user"
	p.Metadata["department"] = "engineering"
	p.Metadata["auth_method"] = "certificate"

	if p.Metadata["type"] != "user" {
		t.Errorf("Expected metadata type 'user', got %v", p.Metadata["type"])
	}

	if p.Metadata["department"] != "engineering" {
		t.Errorf("Expected metadata department 'engineering', got %v", p.Metadata["department"])
	}

	if p.Metadata["auth_method"] != "certificate" {
		t.Errorf("Expected metadata auth_method 'certificate', got %v", p.Metadata["auth_method"])
	}
}

func TestParticipantStatusValidation(t *testing.T) {
	tests := []struct {
		status ParticipantStatus
		want   bool
	}{
		{ParticipantStatusActive, true},
		{ParticipantStatusInactive, true},
		{ParticipantStatusSuspended, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.want {
				t.Errorf("ParticipantStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
