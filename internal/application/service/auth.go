package service

import (
	"context"
	"fmt"
	"time"

	"github.com/diegoclair/apperr"
	"github.com/diegoclair/appvalidator/apperrmap"
	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/internal/application/dto"
	"github.com/diegoclair/go_boilerplate/internal/domain"
	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	"github.com/diegoclair/go_boilerplate/internal/domain/entity"
	"github.com/diegoclair/go_boilerplate/internal/domain/errcodes"
	"github.com/diegoclair/logger"
)

// The input is irrelevant and not a secret; what matters is that the hash carries
// whatever cost the crypto is configured with today.
const decoySource = "go-boilerplate-decoy"

type authApp struct {
	cache               contract.CacheManager
	crypto              contract.Crypto
	dm                  contract.DataManager
	log                 logger.Logger
	validator           apperrmap.Validator
	accountSvc          contract.AccountApp
	accessTokenDuration time.Duration
	decoyPassword       string
}

func newAuthApp(infra domain.Infrastructure, accountSvc contract.AccountApp, accessTokenDuration time.Duration) (*authApp, error) {
	cryptoClient := infra.Crypto()

	decoyPassword, err := cryptoClient.HashPassword(decoySource)
	if err != nil {
		return nil, fmt.Errorf("auth service: mint the decoy hash: %w", err)
	}

	return &authApp{
		cache:               infra.CacheManager(),
		crypto:              cryptoClient,
		dm:                  infra.DataManager(),
		log:                 infra.Logger(),
		validator:           infra.Validator(),
		accountSvc:          accountSvc,
		accessTokenDuration: accessTokenDuration,
		decoyPassword:       decoyPassword,
	}, nil
}

func (s *authApp) Login(ctx context.Context, input dto.LoginInput) (account entity.Account, err error) {
	err = input.Validate(ctx, s.validator)
	if err != nil {
		s.log.Error(ctx, "error or invalid input", logger.Err(err))
		return entity.Account{}, err
	}

	account, err = s.dm.Account().GetAccountByDocument(ctx, input.CPF)
	if err != nil {
		if !apperr.IsNotFound(err) {
			s.log.Error(ctx, "error getting account by document", logger.Err(err))
			return entity.Account{}, err
		}

		// Paying the hashing cost here keeps the answer for a document nobody
		// owns as slow as the answer for one that exists.
		_ = s.crypto.CheckPassword(input.Password, s.decoyPassword)

		s.log.Error(ctx, "account not found")
		return entity.Account{}, errcodes.ErrInvalidCredentials
	}

	ctx = context.WithValue(ctx, infra.AccountUUIDKey, account.UUID)
	ctx = logger.WithAttrs(ctx, logger.Attr("account_id", account.ID))

	err = s.crypto.CheckPassword(input.Password, account.Password)
	if err != nil {
		s.log.Error(ctx, "wrong password")
		return entity.Account{}, errcodes.ErrInvalidCredentials
	}

	// Only whoever already proved the password gets told the account exists but is off.
	if !account.Active {
		s.log.Error(ctx, "account is not active")
		return entity.Account{}, errcodes.ErrDeactivatedAccount
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

	// Outliving the token itself is what stops the middleware from honouring a
	// logged-out token until it expires on its own.
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
