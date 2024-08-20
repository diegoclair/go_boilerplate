package crypto

import (
	"fmt"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"golang.org/x/crypto/bcrypt"
)

type crypto struct{}

func NewCrypto() contract.Crypto {
	return &crypto{}
}

// HashPassword returns the bcrypt hash of the password
func (c *crypto) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func (c *crypto) CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
