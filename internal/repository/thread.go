package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/TranTheTuan/vna/internal/domain"
)

// ThreadRepository defines data access methods for chat threads.
type ThreadRepository interface {
	Create(ctx context.Context, userID string) (*domain.Thread, error)
	ListByUser(ctx context.Context, userID string) ([]*domain.Thread, error)
	// GetByIDAndUser fetches a thread only if it belongs to userID — prevents IDOR.
	GetByIDAndUser(ctx context.Context, threadID, userID string) (*domain.Thread, error)
	// Rename updates the name of a thread, returning the updated record.
	Rename(ctx context.Context, threadID, name string) (*domain.Thread, error)
	// Delete removes a thread by ID — used to clean up orphan threads on upstream failure.
	Delete(ctx context.Context, threadID string) error
}

type threadRepository struct{ db *sql.DB }

// NewThreadRepository creates a ThreadRepository backed by the given *sql.DB.
func NewThreadRepository(db *sql.DB) ThreadRepository {
	return &threadRepository{db: db}
}

// Create inserts a new thread for the user with the default name "New Chat".
func (r *threadRepository) Create(ctx context.Context, userID string) (*domain.Thread, error) {
	const q = `
		INSERT INTO threads(user_id)
		VALUES($1)
		RETURNING id, user_id, name, created_at`
	t := &domain.Thread{}
	err := r.db.QueryRowContext(ctx, q, userID).
		Scan(&t.ID, &t.UserID, &t.Name, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.thread: create: %w", err)
	}
	return t, nil
}

// ListByUser returns all threads for the user ordered newest-first.
func (r *threadRepository) ListByUser(ctx context.Context, userID string) ([]*domain.Thread, error) {
	const q = `
		SELECT id, user_id, name, created_at
		FROM threads
		WHERE user_id = $1
		ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("repository.thread: list: %w", err)
	}
	defer rows.Close()

	var threads []*domain.Thread
	for rows.Next() {
		t := &domain.Thread{}
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("repository.thread: scan: %w", err)
		}
		threads = append(threads, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository.thread: rows error: %w", err)
	}
	return threads, nil
}

// GetByIDAndUser fetches a thread by ID, returning sql.ErrNoRows if not found or not owned by userID.
func (r *threadRepository) GetByIDAndUser(ctx context.Context, threadID, userID string) (*domain.Thread, error) {
	const q = `SELECT id, user_id, name, created_at FROM threads WHERE id=$1 AND user_id=$2`
	t := &domain.Thread{}
	err := r.db.QueryRowContext(ctx, q, threadID, userID).
		Scan(&t.ID, &t.UserID, &t.Name, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.thread: get: %w", err)
	}
	return t, nil
}

// Delete removes a thread by ID. Used to clean up orphan threads when upstream AI call fails.
func (r *threadRepository) Delete(ctx context.Context, threadID string) error {
	const q = `DELETE FROM threads WHERE id=$1`
	if _, err := r.db.ExecContext(ctx, q, threadID); err != nil {
		return fmt.Errorf("repository.thread: delete: %w", err)
	}
	return nil
}

// Rename updates the name of a thread and returns the updated record.
func (r *threadRepository) Rename(ctx context.Context, threadID, name string) (*domain.Thread, error) {
	const q = `
		UPDATE threads SET name=$2
		WHERE id=$1
		RETURNING id, user_id, name, created_at`
	t := &domain.Thread{}
	err := r.db.QueryRowContext(ctx, q, threadID, name).
		Scan(&t.ID, &t.UserID, &t.Name, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.thread: rename: %w", err)
	}
	return t, nil
}
