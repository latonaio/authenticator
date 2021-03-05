package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Encrypt encrypts password.
func Encrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash: %v", err)
	}
	return string(hash), nil
}

// CompareHashAndPassword compares a hashed password with its plaintext password.
// Returns nil on success, or an error on failure.
func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
