package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2id parameters following OWASP recommendations:
// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
const (
	argonMemory      = 19 * 1024 // 19 MiB
	argonIterations  = 2
	argonParallelism = 1
	argonSaltLength  = 16
	argonKeyLength   = 32
)

type Client struct{}

// NewCrypto returns a new crypto client
func NewCrypto() *Client {
	return &Client{}
}

// HashPassword returns the Argon2id hash of the password.
func (c *Client) HashPassword(password string) (string, error) {
	salt := make([]byte, argonSaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, argonIterations, argonMemory, argonParallelism, argonKeyLength)

	// Format: $argon2id$v=19$m=19456,t=2,p=1$<salt>$<hash>
	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonMemory,
		argonIterations,
		argonParallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encoded, nil
}

// CheckPassword verifies a password against an Argon2id hashed password.
func (c *Client) CheckPassword(password, hashedPassword string) error {
	parts := strings.Split(hashedPassword, "$")
	// Expected format: $argon2id$v=19$m=X,t=X,p=X$salt$hash
	// Split produces: ["", "argon2id", "v=19", "m=X,t=X,p=X", "salt", "hash"]
	if len(parts) != 6 || parts[1] != "argon2id" {
		return fmt.Errorf("invalid argon2id hash format")
	}

	var memory uint32
	var iterations uint32
	var parallelism uint8
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return fmt.Errorf("invalid argon2id parameters: %w", err)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return fmt.Errorf("invalid argon2id salt: %w", err)
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return fmt.Errorf("invalid argon2id hash: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, uint32(len(expectedHash)))

	if subtle.ConstantTimeCompare(hash, expectedHash) != 1 {
		return fmt.Errorf("password does not match")
	}

	return nil
}
