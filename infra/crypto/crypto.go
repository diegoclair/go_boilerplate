package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Client struct{}

// NewCrypto returns a new crypto client
func NewCrypto() *Client {
	return &Client{}
}

// HashPassword returns the bcrypt hash of the password
func (c *Client) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func (c *Client) CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
