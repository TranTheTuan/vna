// Package service implements the business logic for the application.
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/TranTheTuan/vna/configs"
	"github.com/TranTheTuan/vna/internal/domain"
	"github.com/TranTheTuan/vna/internal/repository"
	"github.com/TranTheTuan/vna/pkg/argon2_util"
	"github.com/TranTheTuan/vna/pkg/jwt_util"
)

// Sentinel errors returned by UserService methods.
var (
	ErrDuplicateEmail     = errors.New("service.user: email already registered")
	ErrInvalidEmail       = errors.New("service.user: invalid email format")
	ErrPasswordTooShort   = errors.New("service.user: password must be at least 8 characters")
	ErrInvalidCredentials = errors.New("service.user: invalid email or password")
	ErrTokenInvalid       = errors.New("service.user: refresh token is invalid")
	ErrTokenExpired       = errors.New("service.user: refresh token has expired")
	ErrTokenRevoked       = errors.New("service.user: refresh token has been revoked")
)

// emailRegexp is a simple RFC-5322-lite validator.
var emailRegexp = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// UserService defines the auth operations.
type UserService interface {
	Register(ctx context.Context, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
	RefreshToken(ctx context.Context, rawRefreshToken string) (accessToken string, err error)
	Logout(ctx context.Context, rawRefreshToken string) error
}

type userService struct {
	cfg    *configs.Config
	repo   repository.UserRepository
	logger *slog.Logger
}

// NewUserService creates a UserService with the given config, repository, and structured logger.
func NewUserService(cfg *configs.Config, repo repository.UserRepository, logger *slog.Logger) UserService {
	return &userService{cfg: cfg, repo: repo, logger: logger}
}

// Register validates inputs, hashes the password, and creates a new user.
func (s *userService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	if !emailRegexp.MatchString(email) {
		return nil, ErrInvalidEmail
	}
	if len(password) < 8 {
		return nil, ErrPasswordTooShort
	}

	hash, err := argon2_util.HashPassword(password)
	if err != nil {
		s.logger.Error("hash password failed", "error", err)
		return nil, fmt.Errorf("service.user: hash password: %w", err)
	}

	user, err := s.repo.Create(ctx, email, hash)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, ErrDuplicateEmail
		}
		s.logger.Error("create user failed", "error", err)
		return nil, fmt.Errorf("service.user: register: %w", err)
	}
	return user, nil
}

// Login verifies credentials and returns a new access + refresh token pair.
func (s *userService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Run dummy verify to prevent timing-based user enumeration
			_ = argon2_util.VerifyPassword(password, dummyHash)
			return "", "", ErrInvalidCredentials
		}
		s.logger.Error("find user by email failed", "error", err)
		return "", "", fmt.Errorf("service.user: login: %w", err)
	}

	if err := argon2_util.VerifyPassword(password, user.PasswordHash); err != nil {
		return "", "", ErrInvalidCredentials
	}

	accessToken, err := jwt_util.GenerateAccessToken(user.ID, user.Email, s.cfg.Auth.JWTSecret, s.cfg.Auth.JWTAccessTTL)
	if err != nil {
		s.logger.Error("generate access token failed", "error", err)
		return "", "", fmt.Errorf("service.user: generate access token: %w", err)
	}

	rawRefresh, hashRefresh, err := jwt_util.GenerateRefreshToken()
	if err != nil {
		s.logger.Error("generate refresh token failed", "error", err)
		return "", "", fmt.Errorf("service.user: generate refresh token: %w", err)
	}

	expiresAt := time.Now().Add(s.cfg.Auth.JWTRefreshTTL)
	if err := s.repo.SaveRefreshToken(ctx, user.ID, hashRefresh, expiresAt); err != nil {
		s.logger.Error("save refresh token failed", "error", err)
		return "", "", fmt.Errorf("service.user: save refresh token: %w", err)
	}

	return accessToken, rawRefresh, nil
}

// RefreshToken validates a raw refresh token and issues a new access token.
func (s *userService) RefreshToken(ctx context.Context, rawRefreshToken string) (string, error) {
	tokenHash := jwt_util.HashRefreshToken(rawRefreshToken)

	row, err := s.repo.FindRefreshToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrTokenInvalid
		}
		s.logger.Error("find refresh token failed", "error", err)
		return "", fmt.Errorf("service.user: find refresh token: %w", err)
	}

	if row.RevokedAt != nil {
		return "", ErrTokenRevoked
	}
	if time.Now().After(row.ExpiresAt) {
		return "", ErrTokenExpired
	}

	user, err := s.repo.FindByID(ctx, row.UserID)
	if err != nil {
		s.logger.Error("find user by id failed during refresh", "error", err)
		return "", fmt.Errorf("service.user: find user for refresh: %w", err)
	}

	accessToken, err := jwt_util.GenerateAccessToken(user.ID, user.Email, s.cfg.Auth.JWTSecret, s.cfg.Auth.JWTAccessTTL)
	if err != nil {
		s.logger.Error("generate access token failed during refresh", "error", err)
		return "", fmt.Errorf("service.user: generate access token: %w", err)
	}

	return accessToken, nil
}

// Logout revokes the given refresh token, invalidating that session.
func (s *userService) Logout(ctx context.Context, rawRefreshToken string) error {
	tokenHash := jwt_util.HashRefreshToken(rawRefreshToken)
	if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		s.logger.Error("revoke refresh token failed", "error", err)
		return fmt.Errorf("service.user: logout: %w", err)
	}
	return nil
}

// isDuplicateKeyError detects PostgreSQL unique-constraint violation (SQLSTATE 23505).
func isDuplicateKeyError(err error) bool {
	return err != nil && (containsStr(err.Error(), "23505") || containsStr(err.Error(), "unique constraint") || containsStr(err.Error(), "duplicate key"))
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && stringContains(s, substr))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// dummyHash is a pre-computed argon2id hash used in constant-time non-user path during login.
// Prevents timing attacks that would reveal whether an email exists.
const dummyHash = "$argon2id$v=19$m=65536,t=3,p=4$AAAAAAAAAAAAAAAAAAAAAA$AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
