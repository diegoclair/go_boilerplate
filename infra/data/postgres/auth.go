package postgres

import (
	"context"

	"github.com/diegoclair/go_boilerplate/internal/application/dto"
	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
)

type authRepo struct {
	db dbConn
}

func newAuthRepo(db dbConn) contract.AuthRepo {
	return &authRepo{
		db: db,
	}
}

func (r *authRepo) CreateSession(ctx context.Context, session dto.Session) (sessionID int64, err error) {
	query := `
		INSERT INTO tab_session (
			session_uuid,
			account_id,
			refresh_token,
			user_agent,
			client_ip,
			is_blocked,
			refresh_token_expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING session_id;
	`

	err = r.db.QueryRow(ctx, query,
		session.SessionUUID,
		session.AccountID,
		session.RefreshToken,
		session.UserAgent,
		session.ClientIP,
		session.IsBlocked,
		session.RefreshTokenExpiredAt,
	).Scan(&sessionID)
	if err != nil {
		return sessionID, handleDBError(err)
	}

	return sessionID, nil
}

func (r *authRepo) GetSessionByUUID(ctx context.Context, sessionUUID string) (session dto.Session, err error) {
	query := `
		SELECT
			ts.session_id,
			ts.session_uuid,
			ta.account_id,
			ts.refresh_token,
			ts.user_agent,
			ts.client_ip,
			ts.is_blocked,
			ts.refresh_token_expires_at

		FROM 	tab_session 			ts

		INNER JOIN tab_account ta
			ON ta.account_id = ts.account_id

		WHERE	ts.session_uuid 		= 	$1
	`

	row := r.db.QueryRow(ctx, query, sessionUUID)

	err = row.Scan(
		&session.SessionID,
		&session.SessionUUID,
		&session.AccountID,
		&session.RefreshToken,
		&session.UserAgent,
		&session.ClientIP,
		&session.IsBlocked,
		&session.RefreshTokenExpiredAt,
	)
	if err != nil {
		return session, handleDBError(err)
	}

	return session, nil
}

func (r *authRepo) SetSessionAsBlocked(ctx context.Context, sessionUUID string) (err error) {
	query := `
		UPDATE tab_session
		SET is_blocked = true,
			update_at  = NOW()
		WHERE session_uuid = $1;
	`

	_, err = r.db.Exec(ctx, query, sessionUUID)
	if err != nil {
		return handleDBError(err)
	}

	return nil
}
