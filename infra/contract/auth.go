package contract

import (
	"context"
	"time"
)

type TokenPayloadInput struct {
	AccountUUID string
	SessionUUID string
}

// What the token is for. A refresh token is long-lived on purpose, so letting
// one open a private route would undo the reason for having two.
type TokenKind string

const (
	TokenKindAccess  TokenKind = "access"
	TokenKindRefresh TokenKind = "refresh"
)

type TokenPayload struct {
	Kind         TokenKind
	AccountUUID  string
	SessionUUID  string
	RefreshToken string
	IssuedAt     time.Time
	ExpiredAt    time.Time
}

type AuthToken interface {
	CreateAccessToken(ctx context.Context, input TokenPayloadInput) (tokenString string, payload TokenPayload, err error)
	CreateRefreshToken(ctx context.Context, input TokenPayloadInput) (tokenString string, payload TokenPayload, err error)
	VerifyToken(ctx context.Context, token string) (payload TokenPayload, err error)
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hashedPassword string) error
}
