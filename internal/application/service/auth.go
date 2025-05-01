package service

import (
	"context"
	"time"

	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/internal/application/dto"
	"github.com/diegoclair/go_boilerplate/internal/domain"
	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	"github.com/diegoclair/go_boilerplate/internal/domain/entity"
	"github.com/diegoclair/go_utils/logger"
	"github.com/diegoclair/go_utils/mysqlutils"
	"github.com/diegoclair/go_utils/resterrors"
	"github.com/diegoclair/go_utils/validator"
)

const (
	wrongLogin            string = "Document or password are wrong"
	errDeactivatedAccount string = "Account is deactivated"
)

type authApp struct {
	cache               contract.CacheManager
	crypto              contract.Crypto
	dm                  contract.DataManager
	log                 logger.Logger
	validator           validator.Validator
	accountSvc          contract.AccountApp
	accessTokenDuration time.Duration
}

func newAuthApp(infra domain.Infrastructure, accountSvc contract.AccountApp, accessTokenDuration time.Duration) *authApp {
	return &authApp{
		cache:               infra.CacheManager(),
		crypto:              infra.Crypto(),
		dm:                  infra.DataManager(),
		log:                 infra.Logger(),
		validator:           infra.Validator(),
		accountSvc:          accountSvc,
		accessTokenDuration: accessTokenDuration,
	}
}

func (s *authApp) Login(ctx context.Context, input dto.LoginInput) (account entity.Account, err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	err = input.Validate(ctx, s.validator)
	if err != nil {
		s.log.Errorw(ctx, "error or invalid input", logger.Err(err))
		return account, err
	}

	account, err = s.dm.Account().GetAccountByDocument(ctx, input.CPF)
	if err != nil {
		s.log.Errorw(ctx, "error getting account by document", logger.Err(err))
		return account, resterrors.NewUnauthorizedError(wrongLogin)
	}

	ctx = context.WithValue(ctx, infra.AccountUUIDKey, account.UUID) // set account uuid in context to be used in logs

	if !account.Active {
		s.log.Error(ctx, "account is not active")
		return account, resterrors.NewUnauthorizedError(errDeactivatedAccount)
	}

	s.log.Infow(ctx, "account information used to login",
		logger.Int64("account_id", account.ID),
		logger.String("name", account.Name),
	)

	err = s.crypto.CheckPassword(input.Password, account.Password)
	if err != nil {
		s.log.Error(ctx, "wrong password")
		return account, resterrors.NewUnauthorizedError(wrongLogin)
	}

	return account, nil
}

func (s *authApp) CreateSession(ctx context.Context, session dto.Session) (err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	err = session.Validate(ctx, s.validator)
	if err != nil {
		s.log.Errorw(ctx, "error or invalid input", logger.Err(err))
		return err
	}

	_, err = s.dm.Auth().CreateSession(ctx, session)
	if err != nil {
		s.log.Errorw(ctx, "error creating session", logger.Err(err))
		return err
	}

	return nil
}

func (s *authApp) GetSessionByUUID(ctx context.Context, sessionUUID string) (session dto.Session, err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	session, err = s.dm.Auth().GetSessionByUUID(ctx, sessionUUID)
	if err != nil {
		if mysqlutils.SQLNotFound(err.Error()) {
			return session, resterrors.NewUnauthorizedError("session not found")
		}
		s.log.Errorw(ctx, "error getting session", logger.Err(err))
		return session, err
	}

	return session, nil
}

func (s *authApp) Logout(ctx context.Context, accessToken string) (err error) {
	s.log.Info(ctx, "Process Started")
	defer s.log.Info(ctx, "Process Finished")

	loggedAccountID, err := s.accountSvc.GetLoggedAccountID(ctx)
	if err != nil {
		return err
	}

	// access token will be on cache for 3 minutes after it duration
	err = s.cache.SetStringWithExpiration(ctx, accessToken, "true", s.accessTokenDuration+3*time.Minute)
	if err != nil {
		s.log.Errorw(ctx, "error logging out", logger.Err(err))
		return err
	}

	err = s.dm.Auth().SetSessionAsBlocked(ctx, loggedAccountID)
	if err != nil {
		s.log.Errorw(ctx, "error logging out", logger.Err(err))
		return err
	}

	return nil
}
