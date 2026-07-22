package auth

import (
	"crypto/subtle"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const passwordHashCost = bcrypt.DefaultCost

// HashPassword creates a bcrypt hash suitable for database storage.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordHashCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword supports a one-time migration from the legacy plaintext
// storage. needsUpgrade is true only after a legacy password matched.
func VerifyPassword(storedPassword, candidate string) (valid, needsUpgrade bool, err error) {
	if isBcryptHash(storedPassword) {
		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(candidate))
		if err == nil {
			return true, false, nil
		}
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, false, nil
		}
		return false, false, err
	}

	if subtle.ConstantTimeCompare([]byte(storedPassword), []byte(candidate)) == 1 {
		return true, true, nil
	}
	return false, false, nil
}

func isBcryptHash(value string) bool {
	return strings.HasPrefix(value, "$2a$") ||
		strings.HasPrefix(value, "$2b$") ||
		strings.HasPrefix(value, "$2y$")
}
