package models

// User represents a registered user.
// PasswordHash is never returned in API responses.
type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}
