package domain

import "time"

// User represents an authenticated user of the platform.
type User struct {
	ID           string    // UUID
	Email        string
	PasswordHash string    // argon2id PHC encoded string
	CreatedAt    time.Time
}
