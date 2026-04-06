package dto

import "time"

// SendMessageRequest is the body for POST /api/v1/messages.
type SendMessageRequest struct {
	Message string `json:"message"`
}

// MessageResponse is returned for a single message (send or list item).
type MessageResponse struct {
	ID        string    `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
}

// ListResponse is returned by GET /api/v1/messages.
type ListResponse struct {
	Data       []*MessageResponse `json:"data"`
	NextCursor string             `json:"next_cursor"` // empty string when no more pages
}

// StreamDeltaEvent is emitted as SSE "delta" events during POST /api/v1/messages/stream.
type StreamDeltaEvent struct {
	Delta string `json:"delta"`
}
