package kubernetesimplementationindependend_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubernetesimplementationindependend"
)

func Test_SanitizeKindName(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		sanitized, err := kubernetesimplementationindependend.SanitizeKindName("")
		require.Error(t, err)
		require.Empty(t, sanitized)
	})

	t.Run("Only one char", func(t *testing.T) {
		sanitized, err := kubernetesimplementationindependend.SanitizeKindName("s")
		require.Error(t, err)
		require.Empty(t, sanitized)
	})

	t.Run("Secret", func(t *testing.T) {
		sanitized, err := kubernetesimplementationindependend.SanitizeKindName("Secret")
		require.NoError(t, err)
		require.EqualValues(t, "Secret", sanitized)
	})

	t.Run("secret", func(t *testing.T) {
		sanitized, err := kubernetesimplementationindependend.SanitizeKindName("secret")
		require.NoError(t, err)
		require.EqualValues(t, "Secret", sanitized)
	})
}
