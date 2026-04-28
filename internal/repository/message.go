package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/TranTheTuan/vna/internal/domain"
)

// MessageRepository defines data access methods for chat messages.
type MessageRepository interface {
	Save(ctx context.Context, msg *domain.Message) (*domain.Message, error)
	ListByThread(ctx context.Context, threadID string, limit int, cursor string) ([]*domain.Message, string, error)
}

type messageRepository struct {
	db *sql.DB
}

// NewMessageRepository creates a MessageRepository backed by the given *sql.DB.
func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{db: db}
}

// Save inserts a new message and returns the persisted record with generated ID and timestamps.
func (r *messageRepository) Save(ctx context.Context, msg *domain.Message) (*domain.Message, error) {
	const q = `
		INSERT INTO messages(user_id, thread_id, question, answer)
		VALUES($1, $2, $3, $4)
		RETURNING id, user_id, thread_id, question, answer, created_at`

	out := &domain.Message{}
	err := r.db.QueryRowContext(ctx, q, msg.UserID, msg.ThreadID, msg.Question, msg.Answer).
		Scan(&out.ID, &out.UserID, &out.ThreadID, &out.Question, &out.Answer, &out.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.message: save: %w", err)
	}
	return out, nil
}

// ListByThread returns up to limit messages for a thread, ordered newest-first.
// cursor is the ID of the last seen message; empty string returns the first page.
// Thread ownership must be validated by the caller before invoking this method.
func (r *messageRepository) ListByThread(ctx context.Context, threadID string, limit int, cursor string) ([]*domain.Message, string, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if cursor == "" {
		// First page: fetch most recent messages for the thread
		const q = `
			SELECT id, user_id, thread_id, question, answer, created_at
			FROM messages
			WHERE thread_id = $1
			ORDER BY created_at DESC, id DESC
			LIMIT $2`
		rows, err = r.db.QueryContext(ctx, q, threadID, limit)
	} else {
		// Resolve the cursor row's created_at for keyset pagination
		var cursorTime time.Time
		const qCursor = `SELECT created_at FROM messages WHERE id = $1 AND thread_id = $2`
		if scanErr := r.db.QueryRowContext(ctx, qCursor, cursor, threadID).Scan(&cursorTime); scanErr != nil {
			if errors.Is(scanErr, sql.ErrNoRows) {
				// Invalid cursor — return empty page rather than error
				return []*domain.Message{}, "", nil
			}
			return nil, "", fmt.Errorf("repository.message: resolve cursor: %w", scanErr)
		}

		// Keyset: fetch messages strictly older than the cursor position
		const q = `
			SELECT id, user_id, thread_id, question, answer, created_at
			FROM messages
			WHERE thread_id = $1
			  AND (created_at, id) < ($2, $3)
			ORDER BY created_at DESC, id DESC
			LIMIT $4`
		rows, err = r.db.QueryContext(ctx, q, threadID, cursorTime, cursor, limit)
	}

	if err != nil {
		return nil, "", fmt.Errorf("repository.message: list by thread: %w", err)
	}
	defer rows.Close()

	var msgs []*domain.Message
	for rows.Next() {
		m := &domain.Message{}
		if err := rows.Scan(&m.ID, &m.UserID, &m.ThreadID, &m.Question, &m.Answer, &m.CreatedAt); err != nil {
			return nil, "", fmt.Errorf("repository.message: scan row: %w", err)
		}
		msgs = append(msgs, m)
	}
	if err := rows.Err(); err != nil {
		return nil, "", fmt.Errorf("repository.message: rows error: %w", err)
	}

	// Determine next cursor: if we got a full page, there may be more
	nextCursor := ""
	if len(msgs) == limit {
		nextCursor = msgs[len(msgs)-1].ID
	}

	return msgs, nextCursor, nil
}
