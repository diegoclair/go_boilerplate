package service

import (
	"context"
	"time"

	"github.com/diegoclair/apperr"
	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/internal/application/dto"
	"github.com/diegoclair/go_boilerplate/internal/domain"
	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	"github.com/diegoclair/go_boilerplate/internal/domain/entity"
	"github.com/diegoclair/go_boilerplate/internal/domain/errcodes"
	"github.com/diegoclair/logger"
	"github.com/diegoclair/appvalidator/apperrmap"
)

type authApp struct {
	cache               contract.CacheManager
	crypto              contract.Crypto
	dm                  contract.DataManager
	log                 logger.Logger
	validator           apperrmap.Validator
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
	err = input.Validate(ctx, s.validator)
	if err != nil {
		s.log.Error(ctx, "error or invalid input", logger.Err(err))
		return account, err
	}

	account, err = s.dm.Account().GetAccountByDocument(ctx, input.CPF)
	if err != nil {
		s.log.Error(ctx, "error getting account by document", logger.Err(err))
		return account, errcodes.ErrInvalidCredentials
	}

	ctx = context.WithValue(ctx, infra.AccountUUIDKey, account.UUID)
	ctx = logger.WithAttrs(ctx, logger.Attr("account_id", account.ID))

	if !account.Active {
		s.log.Error(ctx, "account is not active")
		return account, errcodes.ErrDeactivatedAccount
	}

	err = s.crypto.CheckPassword(input.Password, account.Password)
	if err != nil {
		s.log.Error(ctx, "wrong password")
		return account, errcodes.ErrInvalidCredentials
	}

	return account, nil
}

func (s *authApp) CreateSession(ctx context.Context, session dto.Session) (err error) {
	err = session.Validate(ctx, s.validator)
	if err != nil {
		s.log.Error(ctx, "error or invalid input", logger.Err(err))
		return err
	}

	_, err = s.dm.Auth().CreateSession(ctx, session)
	if err != nil {
		s.log.Error(ctx, "error creating session", logger.Err(err))
		return err
	}

	return nil
}

func (s *authApp) GetSessionByUUID(ctx context.Context, sessionUUID string) (session dto.Session, err error) {
	ctx = logger.WithAttrs(ctx, logger.Attr("session_uuid", sessionUUID))

	session, err = s.dm.Auth().GetSessionByUUID(ctx, sessionUUID)
	if err != nil {
		if apperr.IsNotFound(err) {
			return session, errcodes.ErrSessionNotFound
		}
		s.log.Error(ctx, "error getting session", logger.Err(err))
		return session, err
	}

	return session, nil
}

func (s *authApp) Logout(ctx context.Context, accessToken string) (err error) {
	sessionUUID, ok := ctx.Value(infra.SessionKey).(string)
	if !ok || sessionUUID == "" {
		s.log.Error(ctx, "session UUID not found in context")
		return errcodes.ErrSessionNotFound
	}

	// access token will be on cache for 3 minutes after it duration
	// this is to avoid the user to login again with the same access token (used in the middleware)
	err = s.cache.Set(ctx, accessToken, "true", s.accessTokenDuration+3*time.Minute)
	if err != nil {
		s.log.Error(ctx, "error logging out", logger.Err(err))
		return err
	}

	err = s.dm.Auth().SetSessionAsBlocked(ctx, sessionUUID)
	if err != nil {
		s.log.Error(ctx, "error logging out", logger.Err(err))
		return err
	}

	return nil
}
