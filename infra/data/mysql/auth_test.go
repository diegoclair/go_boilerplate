package mysql

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/diegoclair/go_boilerplate/internal/application/dto"
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

	sessionID, err := testMysql.Auth().CreateSession(ctx, session)
	require.NoError(t, err)
	require.NotZero(t, sessionID)

	session2, err := testMysql.Auth().GetSessionByUUID(ctx, session.SessionUUID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)
	validateTwoSessions(t, session, session2)
}

func TestCreateSessionErrorsWithMock(t *testing.T) {
	testForInsertErrorsWithMock(t, func(db *sql.DB) error {
		_, err := newAuthRepo(db).CreateSession(context.Background(), dto.Session{})
		return err
	})
}

func TestGetSessionErrorsWithMock(t *testing.T) {
	testForSelectErrorsWithMock(t, "session_id", func(db *sql.DB) error {
		_, err := newAuthRepo(db).GetSessionByUUID(context.Background(), "session-uuid")
		return err
	})
}

func TestSetSessionAsBlocked(t *testing.T) {
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

	sessionID, err := testMysql.Auth().CreateSession(ctx, session)
	require.NoError(t, err)
	require.NotZero(t, sessionID)

	err = testMysql.Auth().SetSessionAsBlocked(ctx, session.AccountID)
	require.NoError(t, err)

	session2, err := testMysql.Auth().GetSessionByUUID(ctx, session.SessionUUID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)
	require.True(t, session2.IsBlocked)
}

func TestSetSessionAsBlockedErrorsWithMock(t *testing.T) {
	testForUpdateDeleteErrorsWithMock(t, func(db *sql.DB) error {
		return newAuthRepo(db).SetSessionAsBlocked(context.Background(), 1)
	})
}
