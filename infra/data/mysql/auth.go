package mysql

import (
	"context"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_utils/mysqlutils"
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
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return sessionID, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		session.SessionUUID,
		session.AccountID,
		session.RefreshToken,
		session.UserAgent,
		session.ClientIP,
		session.IsBlocked,
		session.RefreshTokenExpiredAt,
	)
	if err != nil {
		return sessionID, mysqlutils.HandleMySQLError(err)
	}

	sessionID, err = result.LastInsertId()
	if err != nil {
		return sessionID, mysqlutils.HandleMySQLError(err)
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

		WHERE	ts.session_uuid 		= 	?

	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return session, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, sessionUUID)

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
		return session, err
	}

	return session, nil
}
