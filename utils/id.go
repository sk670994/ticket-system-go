package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// NewID generates a random 16-character hexadecimal ID.
func NewID() string {
	b := make([]byte, 8)

	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}
