package dto

import "time"

// SendMessageRequest is the body for POST /api/v1/messages and POST /api/v1/messages/stream.
type SendMessageRequest struct {
	Message  string `json:"message"`
	ThreadID string `json:"thread_id"` // empty = create new thread implicitly
}

// MessageResponse is returned for a single message (send or list item).
type MessageResponse struct {
	ID        string    `json:"id"`
	ThreadID  string    `json:"thread_id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	CreatedAt time.Time `json:"created_at"`
}

// ListResponse is returned by GET /api/v1/messages.
type ListResponse struct {
	Data       []*MessageResponse `json:"data"`
	NextCursor string             `json:"next_cursor"` // empty string when no more pages
}

// StreamMetaEvent is the first SSE event emitted on POST /api/v1/messages/stream.
// It carries the thread_id so the client can associate the stream with a thread
// before any content delta arrives.
type StreamMetaEvent struct {
	ThreadID string `json:"thread_id"`
}

// StreamDeltaEvent is emitted as SSE "delta" events during POST /api/v1/messages/stream.
type StreamDeltaEvent struct {
	Delta string `json:"delta"`
}
