package dto

import "time"

type Session struct {
	SessionID             int64
	SessionUUID           string
	AccountID             int64
	RefreshToken          string
	UserAgent             string
	ClientIP              string
	IsBlocked             bool
	RefreshTokenExpiredAt time.Time
}
