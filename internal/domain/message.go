package domain

import "time"

// Message represents a single Q&A exchange stored for a user within a thread.
type Message struct {
	ID        string    // UUID
	UserID    string    // UUID of owning user
	ThreadID  string    // UUID of owning thread
	Question  string    // user's original message
	Answer    string    // response from OpenResponses API
	CreatedAt time.Time
}
