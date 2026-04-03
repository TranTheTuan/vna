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
	ListByUser(ctx context.Context, userID string, limit int, cursor string) ([]*domain.Message, string, error)
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
		INSERT INTO messages(user_id, question, answer)
		VALUES($1, $2, $3)
		RETURNING id, user_id, question, answer, created_at`

	out := &domain.Message{}
	err := r.db.QueryRowContext(ctx, q, msg.UserID, msg.Question, msg.Answer).
		Scan(&out.ID, &out.UserID, &out.Question, &out.Answer, &out.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.message: save: %w", err)
	}
	return out, nil
}

// ListByUser returns up to limit messages for a user, ordered newest-first.
// cursor is the ID of the last seen message; empty string returns the first page.
// Returns the messages, the next cursor (last item's ID, empty if no more pages), and any error.
func (r *messageRepository) ListByUser(ctx context.Context, userID string, limit int, cursor string) ([]*domain.Message, string, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if cursor == "" {
		// First page: fetch most recent messages
		const q = `
			SELECT id, user_id, question, answer, created_at
			FROM messages
			WHERE user_id = $1
			ORDER BY created_at DESC, id DESC
			LIMIT $2`
		rows, err = r.db.QueryContext(ctx, q, userID, limit)
	} else {
		// Resolve the cursor row's created_at for keyset pagination
		var cursorTime time.Time
		const qCursor = `SELECT created_at FROM messages WHERE id = $1 AND user_id = $2`
		if scanErr := r.db.QueryRowContext(ctx, qCursor, cursor, userID).Scan(&cursorTime); scanErr != nil {
			if errors.Is(scanErr, sql.ErrNoRows) {
				// Invalid cursor — return empty page rather than error
				return []*domain.Message{}, "", nil
			}
			return nil, "", fmt.Errorf("repository.message: resolve cursor: %w", scanErr)
		}

		// Keyset: fetch messages strictly older than the cursor position
		const q = `
			SELECT id, user_id, question, answer, created_at
			FROM messages
			WHERE user_id = $1
			  AND (created_at, id) < ($2, $3)
			ORDER BY created_at DESC, id DESC
			LIMIT $4`
		rows, err = r.db.QueryContext(ctx, q, userID, cursorTime, cursor, limit)
	}

	if err != nil {
		return nil, "", fmt.Errorf("repository.message: list by user: %w", err)
	}
	defer rows.Close()

	var msgs []*domain.Message
	for rows.Next() {
		m := &domain.Message{}
		if err := rows.Scan(&m.ID, &m.UserID, &m.Question, &m.Answer, &m.CreatedAt); err != nil {
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
