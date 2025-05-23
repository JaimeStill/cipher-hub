package models

import (
	"testing"
)

func TestGenerateID(t *testing.T) {
	id1, err := generateID()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	id2, err := generateID()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id1 == "" {
		t.Error("Expected non-empty ID")
	}

	if id2 == "" {
		t.Error("Expected non-empty ID")
	}

	// ID should be 32 characters (16 bytes hex encoded)
	if len(id1) != 32 {
		t.Errorf("Expected ID length 32, got %d", len(id1))
	}
}
