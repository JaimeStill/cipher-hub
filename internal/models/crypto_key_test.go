package models

import (
	"testing"
)

func TestNewCryptoKey(t *testing.T) {
	serviceID := "service-123"
	name := "Test Key"
	algorithm := AlgorithmAES256

	ck, err := NewCryptoKey(serviceID, name, algorithm)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ck.ID == "" {
		t.Error("Expected generated ID, got empty string")
	}

	if ck.ServiceID != serviceID {
		t.Errorf("Expected ServiceID %q, got %q", serviceID, ck.ServiceID)
	}

	if ck.Name != name {
		t.Errorf("Expected name %q, got %q", name, ck.Name)
	}

	if ck.Algorithm != algorithm {
		t.Errorf("Expected algorithm %q, got %q", algorithm, ck.Algorithm)
	}

	if ck.Version != 1 {
		t.Errorf("Expected version 1, got %d", ck.Version)
	}

	if ck.Status != KeyStatusActive {
		t.Errorf("Expected status %q, got %q", KeyStatusActive, ck.Status)
	}
}

func TestCryptoKeyValidation(t *testing.T) {
	tests := []struct {
		name    string
		ck      *CryptoKey
		wantErr error
	}{
		{
			name: "valid crypto key",
			ck: &CryptoKey{
				ID:        "test-id",
				ServiceID: "service-id",
				Name:      "Test Key",
				Algorithm: AlgorithmAES256,
				Version:   1,
				Status:    KeyStatusActive,
			},
			wantErr: nil,
		},
		{
			name: "invalid algorithm",
			ck: &CryptoKey{
				ID:        "test-id",
				ServiceID: "service-id",
				Name:      "Test Key",
				Algorithm: "invalid",
				Version:   1,
				Status:    KeyStatusActive,
			},
			wantErr: ErrInvalidAlgorithm,
		},
		{
			name: "invalid version",
			ck: &CryptoKey{
				ID:        "test-id",
				ServiceID: "service-id",
				Name:      "Test Key",
				Algorithm: AlgorithmAES256,
				Version:   0,
				Status:    KeyStatusActive,
			},
			wantErr: ErrInvalidVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ck.IsValid()
			if err != tt.wantErr {
				t.Errorf("CryptoKey.IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAlgorithmValidation(t *testing.T) {
	tests := []struct {
		algorithm Algorithm
		want      bool
	}{
		{AlgorithmAES256, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.algorithm), func(t *testing.T) {
			if got := tt.algorithm.IsValid(); got != tt.want {
				t.Errorf("Algorithm.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyStatusValidation(t *testing.T) {
	tests := []struct {
		status KeyStatus
		want   bool
	}{
		{KeyStatusActive, true},
		{KeyStatusRotating, true},
		{KeyStatusDeprecated, true},
		{KeyStatusRevoked, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.want {
				t.Errorf("KeyStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
