package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/TranTheTuan/vna/internal/domain"
	"github.com/TranTheTuan/vna/internal/repository"
)

// ErrThreadNotFound is returned when a thread does not exist or does not belong to the requesting user.
var ErrThreadNotFound = errors.New("service.thread: thread not found")

// ErrInvalidThreadName is returned when a rename request has an empty name.
var ErrInvalidThreadName = errors.New("service.thread: name cannot be empty")

// ThreadService defines operations for managing chat threads.
type ThreadService interface {
	Create(ctx context.Context, userID string) (*domain.Thread, error)
	ListByUser(ctx context.Context, userID string) ([]*domain.Thread, error)
	// Rename updates the thread name, validating that it belongs to userID.
	Rename(ctx context.Context, userID, threadID, name string) (*domain.Thread, error)
}

type threadService struct {
	repo repository.ThreadRepository
}

// NewThreadService creates a ThreadService backed by the given ThreadRepository.
func NewThreadService(repo repository.ThreadRepository) ThreadService {
	return &threadService{repo: repo}
}

func (s *threadService) Create(ctx context.Context, userID string) (*domain.Thread, error) {
	return s.repo.Create(ctx, userID)
}

func (s *threadService) ListByUser(ctx context.Context, userID string) ([]*domain.Thread, error) {
	return s.repo.ListByUser(ctx, userID)
}

// Rename updates the thread name, validating ownership to prevent IDOR.
func (s *threadService) Rename(ctx context.Context, userID, threadID, name string) (*domain.Thread, error) {
	if name == "" {
		return nil, ErrInvalidThreadName
	}
	// Validate ownership — returns ErrThreadNotFound if thread doesn't belong to user.
	if _, err := s.repo.GetByIDAndUser(ctx, threadID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrThreadNotFound
		}
		return nil, fmt.Errorf("service.thread: validate ownership: %w", err)
	}
	t, err := s.repo.Rename(ctx, threadID, name)
	if err != nil {
		return nil, fmt.Errorf("service.thread: rename: %w", err)
	}
	return t, nil
}
