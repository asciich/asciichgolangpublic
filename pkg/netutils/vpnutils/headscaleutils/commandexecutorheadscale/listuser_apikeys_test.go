package commandexecutorheadscale_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/commandexecutorheadscale"
)

func TestListUsersEntry_ToUser(t *testing.T) {
	const userJSON = `{"id": 1, "name": "reto", "created_at": {"seconds": 1775153792}}`

	t.Run("username comes from name and display name stays empty", func(t *testing.T) {
		entry := commandexecutorheadscale.ListUsersEntry{}
		require.NoError(t, json.Unmarshal([]byte(userJSON), &entry))

		user, err := entry.ToUser()
		require.NoError(t, err)

		username, err := user.GetUsername()
		require.NoError(t, err)
		require.Equal(t, "reto", username)

		name, err := user.GetName()
		require.NoError(t, err)
		require.Empty(t, name)

		created, err := user.GetCreated()
		require.NoError(t, err)
		require.EqualValues(t, 1775153792, created.Unix())
	})

	t.Run("missing created time is rejected", func(t *testing.T) {
		entry := commandexecutorheadscale.ListUsersEntry{ID: 1, Name: "reto"}
		_, err := entry.ToUser()
		require.Error(t, err)
	})

	t.Run("invalid id is rejected", func(t *testing.T) {
		entry := commandexecutorheadscale.ListUsersEntry{}
		_, err := entry.ToUser()
		require.Error(t, err)
	})
}

func TestListApiKeyEntry_ToApiKey(t *testing.T) {
	t.Run("missing expiration maps to nil", func(t *testing.T) {
		entry := commandexecutorheadscale.ListApiKeyEntry{}
		require.NoError(t, json.Unmarshal(
			[]byte(`{"id": 2, "prefix": "abcdef", "created_at": {"seconds": 1775153792}}`),
			&entry,
		))

		key, err := entry.ToApiKey()
		require.NoError(t, err)

		expiration, err := key.GetExpiration()
		require.NoError(t, err)
		require.Nil(t, expiration)
	})

	t.Run("empty prefix is rejected", func(t *testing.T) {
		entry := commandexecutorheadscale.ListApiKeyEntry{ID: 2}
		entry.CreatedAt.Seconds = 1775153792

		_, err := entry.ToApiKey()
		require.Error(t, err)
	})
}
