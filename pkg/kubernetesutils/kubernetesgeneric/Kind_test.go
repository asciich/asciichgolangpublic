package kubernetesgeneric_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesgeneric"
)

func Test_SanitizeKindName(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		sanitized, err := kubernetesgeneric.SanitizeKindName("")
		require.Error(t, err)
		require.Empty(t, sanitized)
	})

	t.Run("Only one char", func(t *testing.T) {
		sanitized, err := kubernetesgeneric.SanitizeKindName("s")
		require.Error(t, err)
		require.Empty(t, sanitized)
	})

	t.Run("Secret", func(t *testing.T) {
		sanitized, err := kubernetesgeneric.SanitizeKindName("Secret")
		require.NoError(t, err)
		require.EqualValues(t, "Secret", sanitized)
	})

	t.Run("secret", func(t *testing.T) {
		sanitized, err := kubernetesgeneric.SanitizeKindName("secret")
		require.NoError(t, err)
		require.EqualValues(t, "Secret", sanitized)
	})
}
