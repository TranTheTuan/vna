// Package repository provides PostgreSQL-backed data access for the application.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/TranTheTuan/vna/internal/domain"
)

// RefreshTokenRow holds the data returned when looking up a refresh token.
type RefreshTokenRow struct {
	UserID    string
	ExpiresAt time.Time
	RevokedAt *time.Time // nil means valid (not revoked)
}

// UserRepository defines data access methods for users and refresh tokens.
type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	SaveRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, tokenHash string) (*RefreshTokenRow, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a UserRepository backed by the given *sql.DB.
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user and returns the created record.
func (r *userRepository) Create(ctx context.Context, email, passwordHash string) (*domain.User, error) {
	const q = `
		INSERT INTO users(email, password_hash)
		VALUES($1, $2)
		RETURNING id, email, password_hash, created_at`

	u := &domain.User{}
	err := r.db.QueryRowContext(ctx, q, email, passwordHash).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.user: create: %w", err)
	}
	return u, nil
}

// FindByEmail retrieves a user by email. Returns sql.ErrNoRows if not found.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, created_at
		FROM users WHERE email = $1`

	u := &domain.User{}
	err := r.db.QueryRowContext(ctx, q, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("repository.user: find by email: %w", err)
	}
	return u, nil
}

// FindByID retrieves a user by UUID. Returns sql.ErrNoRows if not found.
func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	const q = `
		SELECT id, email, password_hash, created_at
		FROM users WHERE id = $1`

	u := &domain.User{}
	err := r.db.QueryRowContext(ctx, q, id).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("repository.user: find by id: %w", err)
	}
	return u, nil
}

// SaveRefreshToken stores a hashed refresh token associated with the given user.
func (r *userRepository) SaveRefreshToken(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error {
	const q = `
		INSERT INTO refresh_tokens(user_id, token_hash, expires_at)
		VALUES($1, $2, $3)`

	if _, err := r.db.ExecContext(ctx, q, userID, tokenHash, expiresAt); err != nil {
		return fmt.Errorf("repository.user: save refresh token: %w", err)
	}
	return nil
}

// FindRefreshToken looks up a refresh token by its hash. Returns sql.ErrNoRows if not found.
func (r *userRepository) FindRefreshToken(ctx context.Context, tokenHash string) (*RefreshTokenRow, error) {
	const q = `
		SELECT user_id, expires_at, revoked_at
		FROM refresh_tokens WHERE token_hash = $1`

	row := &RefreshTokenRow{}
	err := r.db.QueryRowContext(ctx, q, tokenHash).
		Scan(&row.UserID, &row.ExpiresAt, &row.RevokedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, fmt.Errorf("repository.user: find refresh token: %w", err)
	}
	return row, nil
}

// RevokeRefreshToken marks a refresh token as revoked by setting revoked_at to now.
func (r *userRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	const q = `
		UPDATE refresh_tokens SET revoked_at = NOW()
		WHERE token_hash = $1`

	if _, err := r.db.ExecContext(ctx, q, tokenHash); err != nil {
		return fmt.Errorf("repository.user: revoke refresh token: %w", err)
	}
	return nil
}
