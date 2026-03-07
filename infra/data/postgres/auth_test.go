package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/internal/application/dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
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
		SessionUUID:           uuid.Must(uuid.NewV7()).String(),
		AccountID:             account.ID,
		RefreshToken:          uuid.Must(uuid.NewV7()).String(),
		UserAgent:             "user-agent",
		ClientIP:              "client-ip",
		IsBlocked:             false,
		RefreshTokenExpiredAt: time.Now().Add(24 * time.Hour),
	}

	sessionID, err := testDB.Auth().CreateSession(ctx, session)
	require.NoError(t, err)
	require.NotZero(t, sessionID)

	session2, err := testDB.Auth().GetSessionByUUID(ctx, session.SessionUUID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)
	validateTwoSessions(t, session, session2)
}

func TestSetSessionAsBlocked(t *testing.T) {
	ctx := context.Background()
	account := createRandomAccount(t)

	session1 := dto.Session{
		SessionUUID:           uuid.Must(uuid.NewV7()).String(),
		AccountID:             account.ID,
		RefreshToken:          uuid.Must(uuid.NewV7()).String(),
		UserAgent:             "user-agent-1",
		ClientIP:              "client-ip-1",
		IsBlocked:             false,
		RefreshTokenExpiredAt: time.Now().Add(24 * time.Hour),
	}

	session2 := dto.Session{
		SessionUUID:           uuid.Must(uuid.NewV7()).String(),
		AccountID:             account.ID,
		RefreshToken:          uuid.Must(uuid.NewV7()).String(),
		UserAgent:             "user-agent-2",
		ClientIP:              "client-ip-2",
		IsBlocked:             false,
		RefreshTokenExpiredAt: time.Now().Add(24 * time.Hour),
	}

	_, err := testDB.Auth().CreateSession(ctx, session1)
	require.NoError(t, err)

	_, err = testDB.Auth().CreateSession(ctx, session2)
	require.NoError(t, err)

	// Block only session1 by UUID
	err = testDB.Auth().SetSessionAsBlocked(ctx, session1.SessionUUID)
	require.NoError(t, err)

	// session1 should be blocked
	got1, err := testDB.Auth().GetSessionByUUID(ctx, session1.SessionUUID)
	require.NoError(t, err)
	require.True(t, got1.IsBlocked)

	// session2 should NOT be blocked
	got2, err := testDB.Auth().GetSessionByUUID(ctx, session2.SessionUUID)
	require.NoError(t, err)
	require.False(t, got2.IsBlocked)
}
