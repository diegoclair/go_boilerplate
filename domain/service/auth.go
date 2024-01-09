package service

import (
	"context"

	"log/slog"

	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_utils-lib/v2/resterrors"
)

const (
	wrongLogin string = "Document or password are wrong"
)

type authService struct {
	svc *service
}

func newAuthService(svc *service) contract.AuthService {
	return &authService{
		svc: svc,
	}
}

func (s *authService) Login(ctx context.Context, cpf, secret string) (account entity.Account, err error) {

	s.svc.log.Info(ctx, "Process Started")
	defer s.svc.log.Info(ctx, "Process Finished")

	account, err = s.svc.dm.Account().GetAccountByDocument(ctx, cpf)
	if err != nil {
		s.svc.log.Error(ctx, err.Error())
		return account, resterrors.NewUnauthorizedError(wrongLogin)
	}

	s.svc.log.Infow(ctx, "account information used to login",
		slog.Group("accountInfo",
			slog.Int64("account_id", account.ID),
			slog.String("account_uuid", account.UUID),
			slog.String("name", account.Name),
		))

	err = s.svc.crypto.CheckPassword(secret, account.Password)
	if err != nil {
		s.svc.log.Error(ctx, "wrong password")
		return account, resterrors.NewUnauthorizedError(wrongLogin)
	}

	return account, nil
}

func (s *authService) CreateSession(ctx context.Context, session entity.Session) (err error) {

	s.svc.log.Info(ctx, "Process Started")
	defer s.svc.log.Info(ctx, "Process Finished")

	err = s.svc.dm.Auth().CreateSession(ctx, session)
	if err != nil {
		s.svc.log.Error(ctx, err.Error())
		return err
	}

	return nil
}

func (s *authService) GetSessionByUUID(ctx context.Context, sessionUUID string) (session entity.Session, err error) {

	s.svc.log.Info(ctx, "Process Started")
	defer s.svc.log.Info(ctx, "Process Finished")

	session, err = s.svc.dm.Auth().GetSessionByUUID(ctx, sessionUUID)
	if err != nil {
		s.svc.log.Error(ctx, err.Error())
		return session, err
	}

	return session, nil
}
