package service

import (
	"context"
	"time"

	"log/slog"

	"github.com/diegoclair/go_boilerplate/application/dto"
	"github.com/diegoclair/go_boilerplate/domain/contract"
	"github.com/diegoclair/go_boilerplate/domain/entity"
	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_utils/resterrors"
)

const (
	wrongLogin            string = "Document or password are wrong"
	errDeactivatedAccount string = "Account is deactivated"
)

type authService struct {
	svc        *service
	accountSvc contract.AccountService
}

func newAuthService(svc *service, accountSvc contract.AccountService) contract.AuthService {
	return &authService{
		svc:        svc,
		accountSvc: accountSvc,
	}
}

func (s *authService) Login(ctx context.Context, input dto.LoginInput) (account entity.Account, err error) {
	s.svc.log.Info(ctx, "Process Started")
	defer s.svc.log.Info(ctx, "Process Finished")

	err = input.Validate(ctx, s.svc.validator)
	if err != nil {
		s.svc.log.Errorf(ctx, "error or invalid input: %s", err.Error())
		return account, err
	}

	account, err = s.svc.dm.Account().GetAccountByDocument(ctx, input.CPF)
	if err != nil {
		s.svc.log.Errorf(ctx, "error getting account by document: %s", err.Error())
		return account, resterrors.NewUnauthorizedError(wrongLogin)
	}

	ctx = context.WithValue(ctx, infra.AccountUUIDKey, account.UUID) // set account uuid in context to be used in logs

	if !account.Active {
		s.svc.log.Error(ctx, "account is not active")
		return account, resterrors.NewUnauthorizedError(errDeactivatedAccount)
	}

	s.svc.log.Infow(ctx, "account information used to login",
		slog.Group("accountInfo",
			slog.Int64("account_id", account.ID),
			slog.String("name", account.Name),
		))

	err = s.svc.crypto.CheckPassword(input.Password, account.Password)
	if err != nil {
		s.svc.log.Error(ctx, "wrong password")
		return account, resterrors.NewUnauthorizedError(wrongLogin)
	}

	return account, nil
}

func (s *authService) CreateSession(ctx context.Context, session dto.Session) (err error) {
	s.svc.log.Info(ctx, "Process Started")
	defer s.svc.log.Info(ctx, "Process Finished")

	err = session.Validate(ctx, s.svc.validator)
	if err != nil {
		s.svc.log.Errorf(ctx, "error or invalid input: %s", err.Error())
		return err
	}

	_, err = s.svc.dm.Auth().CreateSession(ctx, session)
	if err != nil {
		s.svc.log.Errorf(ctx, "error creating session: %s", err.Error())
		return err
	}

	return nil
}

func (s *authService) GetSessionByUUID(ctx context.Context, sessionUUID string) (session dto.Session, err error) {
	s.svc.log.Info(ctx, "Process Started")
	defer s.svc.log.Info(ctx, "Process Finished")

	session, err = s.svc.dm.Auth().GetSessionByUUID(ctx, sessionUUID)
	if err != nil {
		s.svc.log.Errorf(ctx, "error getting session: %s", err.Error())
		return session, err
	}

	return session, nil
}

func (s *authService) Logout(ctx context.Context, accessToken string) (err error) {
	s.svc.log.Info(ctx, "Process Started")
	defer s.svc.log.Info(ctx, "Process Finished")

	loggedAccountID, err := s.accountSvc.GetLoggedAccountID(ctx)
	if err != nil {
		return err
	}

	// access token will be on cache for 3 minutes after it duration
	err = s.svc.cache.SetStringWithExpiration(ctx, accessToken, "true", s.svc.cfg.App.Auth.AccessTokenDuration+3*time.Minute)
	if err != nil {
		s.svc.log.Errorf(ctx, "error logging out: %s", err.Error())
		return err
	}

	err = s.svc.dm.Auth().SetSessionAsBlocked(ctx, loggedAccountID)
	if err != nil {
		s.svc.log.Errorf(ctx, "error logging out: %s", err.Error())
		return err
	}

	return nil
}
