package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/diegoclair/go-boilerplate/util/config"
)

type Key string

func (k Key) String() string {
	return string(k)
}

const (
	AccountUUIDKey  Key = "AccountUUID"
	ContextTokenKey Key = "user-token"
	TokenHeaderName Key = "Token"
)

var (
	TokenSigningMethod = jwt.SigningMethodHS256
)

func GenerateToken(authCfg config.AuthConfig, claims jwt.Claims) (tokenString string, err error) {

	key := []byte(authCfg.PrivateKey)

	token := jwt.NewWithClaims(TokenSigningMethod, claims)
	tokenString, err = token.SignedString(key)
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}
