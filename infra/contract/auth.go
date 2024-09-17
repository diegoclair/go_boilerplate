package contract

import (
	"context"
	"time"
)

type TokenPayloadInput struct {
	AccountUUID string
	SessionUUID string
}

type TokenPayload struct {
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
