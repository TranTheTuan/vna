// Package jwt_util provides JWT access token generation/parsing and opaque refresh token helpers.
package jwt_util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims holds the custom JWT payload for access tokens.
type Claims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

var (
	ErrTokenInvalid = errors.New("jwt_util: token is invalid")
	ErrTokenExpired = errors.New("jwt_util: token is expired")
)

// GenerateAccessToken creates a signed HS256 JWT with user identity claims.
func GenerateAccessToken(userID, email, secret string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("jwt_util: sign token: %w", err)
	}
	return signed, nil
}

// GenerateRefreshToken creates a cryptographically random opaque refresh token.
// Returns the raw base64url token (sent to client) and its SHA-256 hex hash (stored in DB).
func GenerateRefreshToken() (raw string, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", fmt.Errorf("jwt_util: generate refresh token: %w", err)
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	hash = hashToken(raw)
	return raw, hash, nil
}

// HashRefreshToken returns the SHA-256 hex hash of a raw refresh token.
// Used to look up a received token in the DB without storing the raw value.
func HashRefreshToken(raw string) string {
	return hashToken(raw)
}

// ParseAccessToken validates a signed JWT and returns its claims.
// Returns ErrTokenExpired or ErrTokenInvalid on failure.
func ParseAccessToken(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("jwt_util: unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}

// hashToken returns the lowercase SHA-256 hex digest of a string.
func hashToken(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
