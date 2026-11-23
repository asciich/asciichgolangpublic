package headscalegeneric_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalegeneric"
)

func TestUser(t *testing.T) {
	created := time.Date(2026, 4, 2, 18, 16, 32, 0, time.UTC)

	t.Run("empty name and email are allowed", func(t *testing.T) {
		user := headscalegeneric.NewUser()
		require.NoError(t, user.SetID(1))
		require.NoError(t, user.SetUsername("reto"))
		require.NoError(t, user.SetCreated(created))

		name, err := user.GetName()
		require.NoError(t, err)
		require.Empty(t, name)

		email, err := user.GetEmail()
		require.NoError(t, err)
		require.Empty(t, email)

		username, err := user.GetUsername()
		require.NoError(t, err)
		require.Equal(t, "reto", username)
	})

	t.Run("unset required fields return an error", func(t *testing.T) {
		user := headscalegeneric.NewUser()

		_, err := user.GetID()
		require.Error(t, err)

		_, err = user.GetUsername()
		require.Error(t, err)

		_, err = user.GetCreated()
		require.Error(t, err)
	})

	t.Run("invalid values", func(t *testing.T) {
		user := headscalegeneric.NewUser()
		require.Error(t, user.SetID(0))
		require.Error(t, user.SetUsername(""))
		require.Error(t, user.SetCreated(time.Time{}))
	})

	t.Run("deep copy is independent", func(t *testing.T) {
		user := headscalegeneric.NewUser()
		require.NoError(t, user.SetUsername("reto"))

		cp, err := user.GetDeepCopy()
		require.NoError(t, err)
		require.NotSame(t, user, cp)
		require.Equal(t, user, cp)

		require.NoError(t, cp.SetUsername("other"))

		username, err := user.GetUsername()
		require.NoError(t, err)
		require.Equal(t, "reto", username)
	})

	t.Run("nil user", func(t *testing.T) {
		var user *headscalegeneric.User
		_, err := user.GetDeepCopy()
		require.Error(t, err)
	})
}

func TestApiKey(t *testing.T) {
	created := time.Date(2026, 4, 2, 18, 16, 32, 0, time.UTC)

	t.Run("zero time expiration maps to nil", func(t *testing.T) {
		key := headscalegeneric.NewApiKey()

		zero := time.Time{}
		require.NoError(t, key.SetExpiration(&zero))

		expiration, err := key.GetExpiration()
		require.NoError(t, err)
		require.Nil(t, expiration)
	})

	t.Run("fields are returned", func(t *testing.T) {
		key := headscalegeneric.NewApiKey()
		expiration := created.Add(24 * time.Hour)

		require.NoError(t, key.SetID(3))
		require.NoError(t, key.SetPrefix("abcdef"))
		require.NoError(t, key.SetExpiration(&expiration))
		require.NoError(t, key.SetCreated(created))

		prefix, err := key.GetPrefix()
		require.NoError(t, err)
		require.Equal(t, "abcdef", prefix)

		got, err := key.GetExpiration()
		require.NoError(t, err)
		require.True(t, got.Equal(expiration))
	})

	t.Run("unset required fields return an error", func(t *testing.T) {
		key := headscalegeneric.NewApiKey()

		_, err := key.GetID()
		require.Error(t, err)

		_, err = key.GetPrefix()
		require.Error(t, err)

		_, err = key.GetCreated()
		require.Error(t, err)
	})

	t.Run("deep copy copies the expiration", func(t *testing.T) {
		key := headscalegeneric.NewApiKey()
		expiration := created.Add(24 * time.Hour)
		require.NoError(t, key.SetExpiration(&expiration))

		cp, err := key.GetDeepCopy()
		require.NoError(t, err)

		newExpiration := created.Add(48 * time.Hour)
		require.NoError(t, cp.SetExpiration(&newExpiration))

		got, err := key.GetExpiration()
		require.NoError(t, err)
		require.True(t, got.Equal(expiration))
	})

	t.Run("nil api key", func(t *testing.T) {
		var key *headscalegeneric.ApiKey
		_, err := key.GetDeepCopy()
		require.Error(t, err)
	})
}
