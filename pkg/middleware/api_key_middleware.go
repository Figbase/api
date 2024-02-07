package middleware

import (
	"crypto/sha256"
	"crypto/subtle"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

const (
	apiKey = "my-super-secret-key"
)

// AuthMiddleware creates a new instance of key-based authentication middleware.
func AuthMiddleware(apiKey string) func(*fiber.Ctx) error {
	// Validator function for key-based authentication.
	validator := func(c *fiber.Ctx, key string) (bool, error) {
		hashedAPIKey := sha256.Sum256([]byte(apiKey))
		hashedKey := sha256.Sum256([]byte(key))

		if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
			return true, nil
		}
		return false, keyauth.ErrMissingOrMalformedAPIKey
	}

	// Create config for key-based authentication middleware.
	config := keyauth.Config{
		Validator: validator,
	}

	// Return new key-based authentication middleware instance.
	return keyauth.New(config)
}
