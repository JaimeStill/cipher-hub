package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// generateID creates a cryptographically secure random ID
func generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random ID: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
