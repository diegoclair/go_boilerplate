package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/diegoclair/go_boilerplate/infra/contract"
)

// The caller turns this into its own error code — the lockout has no vocabulary
// of the application's errors.
type Status struct {
	Locked    bool
	Remaining int
}

// A wait, never a deactivation: the count expires on its own, so the way back in
// costs nothing and no write to the account can be lost.
type Lockout struct {
	store       contract.LoginAttemptStore
	namespace   string
	maxAttempts int
	window      time.Duration
}

func NewLockout(store contract.LoginAttemptStore, namespace string, maxAttempts int, window time.Duration) (*Lockout, error) {
	switch {
	case store == nil:
		return nil, errors.New("lockout: attempt store is required")
	case strings.TrimSpace(namespace) == "":
		return nil, errors.New("lockout: namespace is required")
	case maxAttempts <= 0:
		return nil, errors.New("lockout: max attempts must be positive")
	case window <= 0:
		return nil, errors.New("lockout: window must be positive")
	}

	return &Lockout{store: store, namespace: namespace, maxAttempts: maxAttempts, window: window}, nil
}

func (l *Lockout) Check(ctx context.Context, identity string) (Status, error) {
	count, err := l.store.GetInt(ctx, l.key(identity))
	if err != nil {
		if errors.Is(err, contract.ErrCacheMiss) {
			return Status{Remaining: l.maxAttempts}, nil
		}
		return Status{}, err
	}

	return l.status(count), nil
}

func (l *Lockout) Fail(ctx context.Context, identity string) (Status, error) {
	count, err := l.store.IncrBy(ctx, l.key(identity), 1, l.window)
	if err != nil {
		return Status{}, err
	}

	return l.status(count), nil
}

func (l *Lockout) Reset(ctx context.Context, identity string) error {
	return l.store.Delete(ctx, l.key(identity))
}

func (l *Lockout) status(count int64) Status {
	if count >= int64(l.maxAttempts) {
		return Status{Locked: true}
	}

	return Status{Remaining: l.maxAttempts - int(count)}
}

// Case and padding are the cheapest way to get a fresh counter, so the identity
// is normalised before it becomes a key.
func (l *Lockout) key(identity string) string {
	return l.namespace + ":login_attempt:" + strings.ToLower(strings.TrimSpace(identity))
}
