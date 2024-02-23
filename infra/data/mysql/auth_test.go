package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/stretchr/testify/assert"
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

func TestCreateSessionErrorsWithMock(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setupTest func(dbMock sqlmock.Sqlmock)
	}{
		{
			name: "Should return error when the prepare fails",
			setupTest: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the exec fails",
			setupTest: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectExec().WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			if tt.setupTest != nil {
				tt.setupTest(mock)
			}

			// wantErr will only work with mocked database
			repo := newAuthRepo(db)
			err = repo.CreateSession(ctx, dto.Session{})
			assert.Error(t, err)
		})
	}
}

func TestGetSessionErrorsWithMock(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		setupTest func(dbMock sqlmock.Sqlmock)
	}{
		{
			name: "Should return error when the prepare fails",
			setupTest: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").WillReturnError(assert.AnError)
			},
		},
		{
			name: "Should return error when the query fails",
			setupTest: func(dbMock sqlmock.Sqlmock) {
				dbMock.ExpectPrepare("").ExpectQuery().WillReturnError(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			if tt.setupTest != nil {
				tt.setupTest(mock)
			}

			repo := newAuthRepo(db)
			_, err = repo.GetSessionByUUID(ctx, "session-uuid")
			assert.Error(t, err)
		})
	}
}
