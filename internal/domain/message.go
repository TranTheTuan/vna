package domain

import "time"

// Message represents a single Q&A exchange stored for a user.
type Message struct {
	ID        string    // UUID
	UserID    string    // UUID of owning user
	Question  string    // user's original message
	Answer    string    // response from OpenResponses API
	CreatedAt time.Time
}
