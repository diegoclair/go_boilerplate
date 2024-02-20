package mysql

import (
	"context"

	"github.com/diegoclair/go_boilerplate/application/contract"
	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_utils/mysqlutils"
)

type authRepo struct {
	db dbConnection
}

func newAuthRepo(db dbConnection) contract.AuthRepo {
	return &authRepo{
		db: db,
	}
}

func (r *authRepo) CreateSession(ctx context.Context, session dto.Session) (err error) {
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

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		session.SessionUUID,
		session.AccountID,
		session.RefreshToken,
		session.UserAgent,
		session.ClientIP,
		session.IsBlocked,
		session.RefreshTokenExpiredAt,
	)
	if err != nil {
		return mysqlutils.HandleMySQLError(err)
	}

	return nil
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

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return session, mysqlutils.HandleMySQLError(err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(sessionUUID)
	if err != nil {
		return session, mysqlutils.HandleMySQLError(err)
	}
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
