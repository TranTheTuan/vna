// Package argon2_util provides password hashing and verification using argon2id.
// Uses OWASP-recommended parameters: memory=65536, time=3, threads=4.
package argon2_util

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// argon2id tuning parameters (OWASP recommended minimums)
const (
	memory     = 64 * 1024 // 64 MB
	iterations = 3
	threads    = 4
	keyLen     = 32
	saltLen    = 16
)

var (
	ErrInvalidHash  = errors.New("argon2_util: invalid hash format")
	ErrHashMismatch = errors.New("argon2_util: password does not match hash")
)

// HashPassword hashes a plaintext password using argon2id and returns a PHC-format encoded string.
func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("argon2_util: generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, threads, keyLen)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		memory, iterations, threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
	return encoded, nil
}

// VerifyPassword checks a plaintext password against a PHC-format argon2id hash.
// Returns nil if they match, ErrHashMismatch if not, or an error on invalid format.
func VerifyPassword(password, encoded string) error {
	// Parse PHC format: $argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return ErrInvalidHash
	}

	var ver int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &ver); err != nil {
		return ErrInvalidHash
	}

	var m, t, p uint32
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &m, &t, &p); err != nil {
		return ErrInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return ErrInvalidHash
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return ErrInvalidHash
	}

	// Re-derive hash with same params
	actualHash := argon2.IDKey([]byte(password), salt, t, m, uint8(p), uint32(len(expectedHash)))

	// Constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare(actualHash, expectedHash) != 1 {
		return ErrHashMismatch
	}
	return nil
}
