package auth

import (
	"context"
	"testing"

	"github.com/diegoclair/apperr"
	"github.com/diegoclair/go_boilerplate/infra/contract"
	"github.com/diegoclair/go_boilerplate/infra/mocks"
	"github.com/diegoclair/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testDocument = "01234567890"
	testPassword = "01234567890"
	testKey      = testNamespace + ":login_attempt:" + testDocument
	testHash     = "real-hash"
	testDecoy    = "decoy-hash"
)

var (
	errWrongCredentials = apperr.Define(apperr.KindAuthentication, "TEST_WRONG_CREDENTIALS", "invalid credentials")
	errLocked           = apperr.Define(apperr.KindAuthentication, "TEST_LOCKED", "too many attempts")
	errDeactivated      = apperr.Define(apperr.KindAuthentication, "TEST_DEACTIVATED", "account is deactivated")
)

type testAccount struct {
	id     int64
	hash   string
	active bool
}

type signInMocks struct {
	store  *mocks.MockLoginAttemptStore
	crypto *mocks.MockPasswordHasher
	found  testAccount
	err    error
	lookup int
}

func newTestSignIn(t *testing.T, ctrl *gomock.Controller) (*SignIn[testAccount], *signInMocks) {
	t.Helper()

	m := &signInMocks{
		store:  mocks.NewMockLoginAttemptStore(ctrl),
		crypto: mocks.NewMockPasswordHasher(ctrl),
		found:  testAccount{id: 9, hash: testHash, active: true},
	}
	m.crypto.EXPECT().HashPassword(gomock.Any()).Return(testDecoy, nil).Times(1)

	lockout, err := NewLockout(m.store, testNamespace, testMaxAttempts, testWindow)
	require.NoError(t, err)

	signIn, err := NewSignIn(testDeps(m, lockout))
	require.NoError(t, err)

	return signIn, m
}

func testDeps(m *signInMocks, lockout *Lockout) SignInDeps[testAccount] {
	return SignInDeps[testAccount]{
		Lockout: lockout,
		Crypto:  m.crypto,
		Logger:  logger.NewNoop(),
		Errors: SignInErrors{
			WrongCredentials: errWrongCredentials,
			Locked:           errLocked,
			Deactivated:      errDeactivated,
		},
		Find: func(context.Context, string) (testAccount, error) {
			m.lookup++
			return m.found, m.err
		},
		PasswordHash: func(a testAccount) string { return a.hash },
		IsActive:     func(a testAccount) bool { return a.active },
	}
}

// A dependency accepted as nil trades a refused boot for a flow that is quietly
// missing half of what makes it safe.
func TestNewSignIn_RefusesAnIncompleteFlow(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := &signInMocks{store: mocks.NewMockLoginAttemptStore(ctrl), crypto: mocks.NewMockPasswordHasher(ctrl)}
	lockout, err := NewLockout(m.store, testNamespace, testMaxAttempts, testWindow)
	require.NoError(t, err)

	tests := map[string]func(*SignInDeps[testAccount]){
		"no lockout":      func(d *SignInDeps[testAccount]) { d.Lockout = nil },
		"no crypto":       func(d *SignInDeps[testAccount]) { d.Crypto = nil },
		"no logger":       func(d *SignInDeps[testAccount]) { d.Logger = nil },
		"no refusal code": func(d *SignInDeps[testAccount]) { d.Errors.Locked = nil },
		"no lookup":       func(d *SignInDeps[testAccount]) { d.Find = nil },
		"no password":     func(d *SignInDeps[testAccount]) { d.PasswordHash = nil },
		"no active flag":  func(d *SignInDeps[testAccount]) { d.IsActive = nil },
	}

	for name, breakIt := range tests {
		t.Run(name, func(t *testing.T) {
			deps := testDeps(m, lockout)
			breakIt(&deps)

			got, err := NewSignIn(deps)

			require.Error(t, err)
			require.Nil(t, got)
		})
	}
}

// Without a decoy the unknown-identity path answers instantly, which is the
// oracle this whole flow exists to close — so booting without one is refused.
func TestNewSignIn_RefusesToBootWithoutADecoy(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := &signInMocks{store: mocks.NewMockLoginAttemptStore(ctrl), crypto: mocks.NewMockPasswordHasher(ctrl)}
	m.crypto.EXPECT().HashPassword(gomock.Any()).Return("", assert.AnError)

	lockout, err := NewLockout(m.store, testNamespace, testMaxAttempts, testWindow)
	require.NoError(t, err)

	got, err := NewSignIn(testDeps(m, lockout))

	require.Error(t, err)
	require.Nil(t, got)
}

// At the ceiling the answer comes from the counter alone; a lookup here would pay
// a different cost for an identity that has an account than for one that does not.
func TestAuthenticate_TheCeilingIsAnsweredBeforeAnyLookup(t *testing.T) {
	ctrl := gomock.NewController(t)
	signIn, m := newTestSignIn(t, ctrl)

	m.store.EXPECT().GetInt(gomock.Any(), testKey).Return(int64(testMaxAttempts), nil)

	_, err := signIn.Authenticate(context.Background(), testDocument, testPassword)

	require.ErrorIs(t, err, errLocked)
	assert.Zero(t, m.lookup, "the account may not be looked up once the ceiling is reached")
}

// The two rejections have to be one answer: same code, same remaining count and
// the same hashing cost paid, or the endpoint tells anyone who asks who is registered.
func TestAuthenticate_UnknownIdentityAnswersExactlyLikeAWrongPassword(t *testing.T) {
	const attempts = int64(2)

	unknownCtrl := gomock.NewController(t)
	unknownSignIn, unknown := newTestSignIn(t, unknownCtrl)
	unknown.err = apperr.ErrRecordNotFound
	unknown.store.EXPECT().GetInt(gomock.Any(), testKey).Return(attempts, nil)
	unknown.crypto.EXPECT().CheckPassword(testPassword, testDecoy).Return(assert.AnError)
	unknown.store.EXPECT().IncrBy(gomock.Any(), testKey, int64(1), testWindow).Return(attempts+1, nil)

	_, unknownErr := unknownSignIn.Authenticate(context.Background(), testDocument, testPassword)

	knownCtrl := gomock.NewController(t)
	knownSignIn, known := newTestSignIn(t, knownCtrl)
	known.store.EXPECT().GetInt(gomock.Any(), testKey).Return(attempts, nil)
	known.crypto.EXPECT().CheckPassword(testPassword, testHash).Return(assert.AnError)
	known.store.EXPECT().IncrBy(gomock.Any(), testKey, int64(1), testWindow).Return(attempts+1, nil)

	_, knownErr := knownSignIn.Authenticate(context.Background(), testDocument, testPassword)

	require.ErrorIs(t, unknownErr, errWrongCredentials)
	require.ErrorIs(t, knownErr, errWrongCredentials)
	assert.Equal(t, apperr.GetMeta(knownErr), apperr.GetMeta(unknownErr))
	assert.Equal(t, testMaxAttempts-3, apperr.GetMeta(unknownErr)["remaining_attempts"])
}

func TestAuthenticate_TheCountThatReachesTheCeilingAnswersLocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	signIn, m := newTestSignIn(t, ctrl)

	m.store.EXPECT().GetInt(gomock.Any(), testKey).Return(int64(testMaxAttempts-1), nil)
	m.crypto.EXPECT().CheckPassword(testPassword, testHash).Return(assert.AnError)
	m.store.EXPECT().IncrBy(gomock.Any(), testKey, int64(1), testWindow).Return(int64(testMaxAttempts), nil)

	_, err := signIn.Authenticate(context.Background(), testDocument, testPassword)

	require.ErrorIs(t, err, errLocked)
	assert.Nil(t, apperr.GetMeta(err)["remaining_attempts"])
}

// A number the next attempt contradicts is worse than no number at all.
func TestAuthenticate_AnUncountableFailureAnswersWithoutARemainingCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	signIn, m := newTestSignIn(t, ctrl)

	m.store.EXPECT().GetInt(gomock.Any(), testKey).Return(int64(0), contract.ErrCacheMiss)
	m.crypto.EXPECT().CheckPassword(testPassword, testHash).Return(assert.AnError)
	m.store.EXPECT().IncrBy(gomock.Any(), testKey, int64(1), testWindow).Return(int64(0), assert.AnError)

	_, err := signIn.Authenticate(context.Background(), testDocument, testPassword)

	require.ErrorIs(t, err, errWrongCredentials)
	assert.Nil(t, apperr.GetMeta(err)["remaining_attempts"])
}

// A counter that cannot be read is a defence that is down, not a reason to refuse
// everyone who wants to log in.
func TestAuthenticate_AnUnreadableCounterDoesNotRefuseTheRightPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	signIn, m := newTestSignIn(t, ctrl)

	m.store.EXPECT().GetInt(gomock.Any(), testKey).Return(int64(0), assert.AnError)
	m.crypto.EXPECT().CheckPassword(testPassword, testHash).Return(nil)
	m.store.EXPECT().Delete(gomock.Any(), testKey).Return(nil)

	account, err := signIn.Authenticate(context.Background(), testDocument, testPassword)

	require.NoError(t, err)
	assert.Equal(t, int64(9), account.id)
}

// The password is right: a counter left standing is worth a log, never a refusal.
func TestAuthenticate_AFailedResetDoesNotRefuseTheRightPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	signIn, m := newTestSignIn(t, ctrl)

	m.store.EXPECT().GetInt(gomock.Any(), testKey).Return(int64(0), contract.ErrCacheMiss)
	m.crypto.EXPECT().CheckPassword(testPassword, testHash).Return(nil)
	m.store.EXPECT().Delete(gomock.Any(), testKey).Return(assert.AnError)

	account, err := signIn.Authenticate(context.Background(), testDocument, testPassword)

	require.NoError(t, err)
	assert.Equal(t, int64(9), account.id)
}

func TestAuthenticate_DeactivatedAccountIsOnlyRevealedByTheRightPassword(t *testing.T) {
	t.Run("wrong password answers wrong credentials", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		signIn, m := newTestSignIn(t, ctrl)
		m.found.active = false

		m.store.EXPECT().GetInt(gomock.Any(), testKey).Return(int64(0), contract.ErrCacheMiss)
		m.crypto.EXPECT().CheckPassword(testPassword, testHash).Return(assert.AnError)
		m.store.EXPECT().IncrBy(gomock.Any(), testKey, int64(1), testWindow).Return(int64(1), nil)

		_, err := signIn.Authenticate(context.Background(), testDocument, testPassword)

		require.ErrorIs(t, err, errWrongCredentials)
		assert.NotErrorIs(t, err, errDeactivated)
	})

	// The counter is cleared first: the password was right, so the next visit must
	// not start closer to the ceiling than a fresh one.
	t.Run("right password answers deactivated after the reset", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		signIn, m := newTestSignIn(t, ctrl)
		m.found.active = false

		m.store.EXPECT().GetInt(gomock.Any(), testKey).Return(int64(0), contract.ErrCacheMiss)
		check := m.crypto.EXPECT().CheckPassword(testPassword, testHash).Return(nil)
		m.store.EXPECT().Delete(gomock.Any(), testKey).Return(nil).After(check)

		_, err := signIn.Authenticate(context.Background(), testDocument, testPassword)

		require.ErrorIs(t, err, errDeactivated)
	})
}

// A store that is down is a system failure, not a refused credential — mapping it
// to one would hand the counter a failure the visitor never made.
func TestAuthenticate_LookupFailureIsASystemError(t *testing.T) {
	ctrl := gomock.NewController(t)
	signIn, m := newTestSignIn(t, ctrl)
	m.err = assert.AnError

	m.store.EXPECT().GetInt(gomock.Any(), testKey).Return(int64(0), contract.ErrCacheMiss)

	_, err := signIn.Authenticate(context.Background(), testDocument, testPassword)

	require.ErrorIs(t, err, assert.AnError)
	assert.NotErrorIs(t, err, errWrongCredentials)
}
