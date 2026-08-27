package crypto

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	c := NewCrypto()

	t.Run("Should return an argon2id hash", func(t *testing.T) {
		hash, err := c.HashPassword("123456")
		require.NoError(t, err)
		require.NotEmpty(t, hash)
		require.True(t, strings.HasPrefix(hash, "$argon2id$"), "hash should be argon2id format, got: %s", hash)
	})

	t.Run("Should generate different hashes for same password", func(t *testing.T) {
		hash1, err := c.HashPassword("same-password")
		require.NoError(t, err)

		hash2, err := c.HashPassword("same-password")
		require.NoError(t, err)

		require.NotEqual(t, hash1, hash2, "each hash should use a unique salt")
	})

	t.Run("Should handle empty password", func(t *testing.T) {
		hash, err := c.HashPassword("")
		require.NoError(t, err)
		require.NotEmpty(t, hash)
	})
}

func TestCheckPassword(t *testing.T) {
	c := NewCrypto()

	t.Run("Should verify correct password with argon2id hash", func(t *testing.T) {
		password := "my-secure-password"
		hash, err := c.HashPassword(password)
		require.NoError(t, err)

		err = c.CheckPassword(password, hash)
		require.NoError(t, err)
	})

	t.Run("Should reject wrong password with argon2id hash", func(t *testing.T) {
		hash, err := c.HashPassword("correct-password")
		require.NoError(t, err)

		err = c.CheckPassword("wrong-password", hash)
		require.Error(t, err)
	})

	t.Run("Should reject invalid argon2id format", func(t *testing.T) {
		err := c.CheckPassword("password", "$argon2id$invalid")
		require.Error(t, err)
	})
}

func TestCheckPassword_RoundTrip(t *testing.T) {
	c := NewCrypto()

	passwords := []string{
		"simple",
		"c0mpl3x!P@ssw0rd#2026",
		"com espaços e acentuação",
		"🔐emoji-password",
		strings.Repeat("a", 100),
	}

	for _, pw := range passwords {
		t.Run("Round trip: "+pw[:min(len(pw), 20)], func(t *testing.T) {
			hash, err := c.HashPassword(pw)
			require.NoError(t, err)

			err = c.CheckPassword(pw, hash)
			require.NoError(t, err)

			err = c.CheckPassword(pw+"x", hash)
			require.Error(t, err)
		})
	}
}

func TestHashPasswordCostIsIndependentOfTheInput(t *testing.T) {
	c := NewCrypto()

	// A decoy hash is only useful against a timing oracle while it costs what a
	// real one costs, whatever the input it was minted from.
	decoy, err := c.HashPassword("any-decoy-source")
	require.NoError(t, err)

	genuine, err := c.HashPassword("a-real-password")
	require.NoError(t, err)

	want := fmt.Sprintf("m=%d,t=%d,p=%d", argonMemory, argonIterations, argonParallelism)
	require.Equal(t, want, costParameters(t, decoy))
	require.Equal(t, want, costParameters(t, genuine))
}

func costParameters(t *testing.T, hash string) string {
	t.Helper()

	parts := strings.Split(hash, "$")
	require.Len(t, parts, 6)

	return parts[3]
}
