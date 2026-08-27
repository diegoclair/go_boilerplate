package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/diegoclair/apperr"
	"github.com/diegoclair/go_boilerplate/infra/contract"
	"github.com/diegoclair/logger"
)

// The input is irrelevant and not a secret; what matters is that the hash carries
// whatever cost the hasher is configured with today, so an unknown identity pays it too.
const decoySource = "go-boilerplate-decoy"

// The codes a refused sign-in answers with, so the flow stays free of any
// application's error vocabulary.
type SignInErrors struct {
	WrongCredentials *apperr.Definition
	Locked           *apperr.Definition
	Deactivated      *apperr.Definition
}

type SignInDeps[T any] struct {
	Lockout *Lockout
	Crypto  contract.PasswordHasher
	Logger  logger.Logger
	Errors  SignInErrors
	// A not-found from Find is answered like a wrong password, never as a failure.
	Find         func(ctx context.Context, identity string) (T, error)
	PasswordHash func(account T) string
	IsActive     func(account T) bool
}

type SignIn[T any] struct {
	lockout      *Lockout
	crypto       contract.PasswordHasher
	log          logger.Logger
	errs         SignInErrors
	decoy        string
	find         func(ctx context.Context, identity string) (T, error)
	passwordHash func(account T) string
	isActive     func(account T) bool
}

func NewSignIn[T any](d SignInDeps[T]) (*SignIn[T], error) {
	switch {
	case d.Lockout == nil:
		return nil, errors.New("sign in: lockout is required")
	case d.Crypto == nil:
		return nil, errors.New("sign in: crypto is required")
	case d.Logger == nil:
		return nil, errors.New("sign in: logger is required")
	case d.Errors.WrongCredentials == nil || d.Errors.Locked == nil || d.Errors.Deactivated == nil:
		return nil, errors.New("sign in: the refusal codes are required")
	case d.Find == nil || d.PasswordHash == nil || d.IsActive == nil:
		return nil, errors.New("sign in: the caller's own account lookups are required")
	}

	decoy, err := d.Crypto.HashPassword(decoySource)
	if err != nil {
		return nil, fmt.Errorf("sign in: mint the decoy hash: %w", err)
	}

	return &SignIn[T]{
		lockout:      d.Lockout,
		crypto:       d.Crypto,
		log:          d.Logger,
		errs:         d.Errors,
		decoy:        decoy,
		find:         d.Find,
		passwordHash: d.PasswordHash,
		isActive:     d.IsActive,
	}, nil
}

// An unknown identity and a wrong password answer the same, so the endpoint never
// becomes a way to ask who has an account here.
func (s *SignIn[T]) Authenticate(ctx context.Context, identity, password string) (T, error) {
	var zero T

	// Answering from the counter alone keeps the ceiling reachable by an identity
	// that has no account, so the answer never depends on the account existing.
	status, err := s.lockout.Check(ctx, identity)
	if err != nil {
		// A counter that cannot be read is a defence that is down, not a reason
		// to refuse everyone who wants to log in.
		s.log.Error(ctx, "error reading the login attempt count", logger.Err(err))
	}
	if status.Locked {
		return zero, s.errs.Locked
	}

	account, err := s.find(ctx, identity)
	if err != nil {
		if !apperr.IsNotFound(err) {
			s.log.Error(ctx, "error loading the account to log in", logger.Err(err))
			return zero, err
		}

		// Without the decoy the unknown identity answers instantly, and the timing
		// tells whoever asks who is registered.
		_ = s.crypto.CheckPassword(password, s.decoy)

		return zero, s.countFailure(ctx, identity)
	}

	if err := s.crypto.CheckPassword(password, s.passwordHash(account)); err != nil {
		return zero, s.countFailure(ctx, identity)
	}

	if err := s.lockout.Reset(ctx, identity); err != nil {
		// The password was right: a counter left standing is worth a log, never
		// a refusal.
		s.log.Error(ctx, "error clearing the login attempt count", logger.Err(err))
	}

	// Read only once the password is right: any earlier and the answer confirms
	// the account to whoever merely guessed the identity.
	if !s.isActive(account) {
		return account, s.errs.Deactivated
	}

	return account, nil
}

func (s *SignIn[T]) countFailure(ctx context.Context, identity string) error {
	status, err := s.lockout.Fail(ctx, identity)
	if err != nil {
		// The count is unknown, so the answer carries no remaining attempts —
		// inventing one here would be a number the next attempt contradicts.
		s.log.Error(ctx, "error counting the failed login", logger.Err(err))
		return s.errs.WrongCredentials
	}

	if status.Locked {
		return s.errs.Locked
	}

	return s.errs.WrongCredentials.WithMeta("remaining_attempts", status.Remaining)
}
