package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"gogym-api/internal/domain/user"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordService implements password hashing and verification using Argon2id
type PasswordService struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// NewPasswordService creates a new password service with secure defaults
func NewPasswordService() *PasswordService {
	return &PasswordService{
		memory:      64 * 1024, // 64MB
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
}

// HashPassword hashes a password using Argon2id
func (s *PasswordService) HashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, s.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", user.NewDomainError(user.ErrInternal, "salt_generation_failed", "failed to generate salt")
	}

	// Generate the hash
	hash := argon2.IDKey([]byte(password), salt, s.iterations, s.memory, s.parallelism, s.keyLength)

	// Encode the hash in the format: $argon2id$v=19$m=memory,t=iterations,p=parallelism$salt$hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, s.memory, s.iterations, s.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// VerifyPassword verifies a password against its Argon2id hash
func (s *PasswordService) VerifyPassword(password, encodedHash string) error {
	// Parse the encoded hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return user.NewDomainError(user.ErrInvalidInput, "invalid_hash_format", "invalid hash format")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return user.NewDomainError(user.ErrInvalidInput, "invalid_version", "invalid hash version")
	}
	if version != argon2.Version {
		return user.NewDomainError(user.ErrInvalidInput, "incompatible_version", "incompatible hash version")
	}

	var memory, iterations uint32
	var parallelism uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
		return user.NewDomainError(user.ErrInvalidInput, "invalid_parameters", "invalid hash parameters")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return user.NewDomainError(user.ErrInvalidInput, "invalid_salt", "invalid salt encoding")
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return user.NewDomainError(user.ErrInvalidInput, "invalid_hash", "invalid hash encoding")
	}

	// Compute the hash of the provided password using the same parameters
	comparisonHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(decodedHash)))

	// Compare the hashes using constant-time comparison
	if subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1 {
		return nil
	}

	return user.NewDomainError(user.ErrUnauthorized, "invalid_password", "password verification failed")
}