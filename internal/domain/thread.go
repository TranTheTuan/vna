package domain

import "time"

// Thread represents a named chat conversation belonging to a user.
type Thread struct {
	ID        string
	UserID    string
	Name      string
	CreatedAt time.Time
}
