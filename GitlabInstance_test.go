package asciichgolangpublic_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	asciichgolangpublic "github.com/asciich/asciichgolangpublic"
)

func TestGitlabInstance_GetDeepCopy(t *testing.T) {
	t.Run("nil fields", func(t *testing.T) {
		original := asciichgolangpublic.NewGitlabInstance()
		copy := original.GetDeepCopy()

		require.NotNil(t, copy)

		// Modify copy
		err := copy.SetFqdn("modified.example.com")
		require.NoError(t, err)

		// Original should be unchanged
		originalFqdn, err := original.GetFqdn()
		if err == nil {
			require.NotEqual(t, "modified.example.com", originalFqdn)
		}
	})

	t.Run("with FQDN set", func(t *testing.T) {
		original, err := asciichgolangpublic.GetGitlabByFQDN("gitlab.example.com")
		require.NoError(t, err)

		copy := original.GetDeepCopy()

		require.NotNil(t, copy)

		// Verify FQDN is copied
		copyFqdn, err := copy.GetFqdn()
		require.NoError(t, err)
		require.EqualValues(t, "gitlab.example.com", copyFqdn)

		// Modify copy's FQDN
		err = copy.SetFqdn("modified.example.com")
		require.NoError(t, err)

		// Original should be unchanged
		originalFqdn, err := original.GetFqdn()
		require.NoError(t, err)
		require.EqualValues(t, "gitlab.example.com", originalFqdn)

		copyFqdn, err = copy.GetFqdn()
		require.NoError(t, err)
		require.EqualValues(t, "modified.example.com", copyFqdn)
	})

	t.Run("with FQDN set verifies copy works", func(t *testing.T) {
		original, err := asciichgolangpublic.GetGitlabByFQDN("gitlab.example.com")
		require.NoError(t, err)

		copy := original.GetDeepCopy()

		require.NotNil(t, copy)

		// Verify both have the same FQDN initially
		originalFqdn, err := original.GetFqdn()
		require.NoError(t, err)
		copyFqdn, err := copy.GetFqdn()
		require.NoError(t, err)
		require.EqualValues(t, originalFqdn, copyFqdn)

		// Modify copy's FQDN
		err = copy.SetFqdn("modified.example.com")
		require.NoError(t, err)

		// Original should be unchanged
		originalFqdn, err = original.GetFqdn()
		require.NoError(t, err)
		copyFqdn, err = copy.GetFqdn()
		require.NoError(t, err)
		require.EqualValues(t, "gitlab.example.com", originalFqdn)
		require.EqualValues(t, "modified.example.com", copyFqdn)
	})
}
