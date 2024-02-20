package mysql_test

import (
	"context"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/stretchr/testify/require"
	"github.com/twinj/uuid"
)

func validateTwoSessions(t *testing.T, sessionExpected dto.Session, sessionToCompare dto.Session) {
	require.NotZero(t, sessionToCompare.AccountID)
	require.Equal(t, sessionExpected.SessionUUID, sessionToCompare.SessionUUID)
	require.Equal(t, sessionExpected.RefreshToken, sessionToCompare.RefreshToken)
	require.Equal(t, sessionExpected.UserAgent, sessionToCompare.UserAgent)
	require.Equal(t, sessionExpected.ClientIP, sessionToCompare.ClientIP)
	require.Equal(t, sessionExpected.IsBlocked, sessionToCompare.IsBlocked)
	require.WithinDuration(t, sessionExpected.RefreshTokenExpiredAt, sessionToCompare.RefreshTokenExpiredAt, 2*time.Second)
}

func TestCreateAndGetSession(t *testing.T) {
	ctx := context.Background()
	account := createRandomAccount(t)

	session := dto.Session{
		SessionUUID:           uuid.NewV4().String(),
		AccountID:             account.ID,
		RefreshToken:          uuid.NewV4().String(),
		UserAgent:             "user-agent",
		ClientIP:              "client-ip",
		IsBlocked:             false,
		RefreshTokenExpiredAt: time.Now().Add(24 * time.Hour),
	}

	err := testMysql.Auth().CreateSession(ctx, session)
	require.NoError(t, err)

	session2, err := testMysql.Auth().GetSessionByUUID(ctx, session.SessionUUID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)
	validateTwoSessions(t, session, session2)
}
