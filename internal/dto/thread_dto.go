package dto

import "time"

// RenameThreadRequest is the body for PATCH /api/v1/threads/:id.
type RenameThreadRequest struct {
	Name string `json:"name"`
}

// ThreadResponse represents a single thread in API responses.
type ThreadResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ListThreadsResponse is returned by GET /api/v1/threads.
type ListThreadsResponse struct {
	Data []*ThreadResponse `json:"data"`
}
