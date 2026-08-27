package service

import (
	"context"
	"errors"
	"time"

	"github.com/diegoclair/apperr"
	"github.com/diegoclair/appvalidator/apperrmap"
	"github.com/diegoclair/go_boilerplate/infra"
	"github.com/diegoclair/go_boilerplate/infra/auth"
	"github.com/diegoclair/go_boilerplate/internal/application/dto"
	"github.com/diegoclair/go_boilerplate/internal/domain"
	"github.com/diegoclair/go_boilerplate/internal/domain/contract"
	"github.com/diegoclair/go_boilerplate/internal/domain/entity"
	"github.com/diegoclair/go_boilerplate/internal/domain/errcodes"
	"github.com/diegoclair/logger"
)

// Prefixes the attempt counters so they never collide with another key family.
const lockoutNamespace = "account"

type authApp struct {
	cache               contract.CacheManager
	dm                  contract.DataManager
	log                 logger.Logger
	validator           apperrmap.Validator
	accountSvc          contract.AccountApp
	accessTokenDuration time.Duration
	signIn              *auth.SignIn[entity.Account]
}

func newAuthApp(infra domain.Infrastructure, accountSvc contract.AccountApp, accessTokenDuration time.Duration) (*authApp, error) {
	lockout, err := auth.NewLockout(infra.CacheManager(), lockoutNamespace, domain.MaxLoginAttempts, domain.LoginAttemptWindow)
	if err != nil {
		return nil, err
	}

	dm := infra.DataManager()

	signIn, err := auth.NewSignIn(auth.SignInDeps[entity.Account]{
		Lockout: lockout,
		Crypto:  infra.Crypto(),
		Logger:  infra.Logger(),
		Errors: auth.SignInErrors{
			WrongCredentials: errcodes.ErrInvalidCredentials,
			Locked:           errcodes.ErrMaxLoginAttempts,
			Deactivated:      errcodes.ErrDeactivatedAccount,
		},
		Find: func(ctx context.Context, document string) (entity.Account, error) {
			return dm.Account().GetAccountByDocument(ctx, document)
		},
		PasswordHash: func(account entity.Account) string { return account.Password },
		IsActive:     func(account entity.Account) bool { return account.Active },
	})
	if err != nil {
		return nil, err
	}

	return &authApp{
		cache:               infra.CacheManager(),
		dm:                  dm,
		log:                 infra.Logger(),
		validator:           infra.Validator(),
		accountSvc:          accountSvc,
		accessTokenDuration: accessTokenDuration,
		signIn:              signIn,
	}, nil
}

func (s *authApp) Login(ctx context.Context, input dto.LoginInput) (account entity.Account, err error) {
	err = input.Validate(ctx, s.validator)
	if err != nil {
		s.log.Error(ctx, "error or invalid input", logger.Err(err))
		return entity.Account{}, err
	}

	account, err = s.signIn.Authenticate(ctx, input.CPF, input.Password)
	if err != nil {
		if errors.Is(err, errcodes.ErrMaxLoginAttempts) {
			s.log.Warn(ctx, "login refused: too many attempts")
		}
		return entity.Account{}, err
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
