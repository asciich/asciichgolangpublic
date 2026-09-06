package asciichgolangpublic_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	asciichgolangpublic "github.com/asciich/asciichgolangpublic"
)

func TestGitlabCreateTagOptions_GetDeepCopy(t *testing.T) {
	t.Run("nil fields", func(t *testing.T) {
		original := &asciichgolangpublic.GitlabCreateTagOptions{}
		copy := original.GetDeepCopy()

		require.EqualValues(t, original.Name, copy.Name)
		require.EqualValues(t, original.Ref, copy.Ref)

		// Modify copy
		copy.Name = "modified"
		copy.Ref = "modified"

		// Original should be unchanged
		require.EqualValues(t, "", original.Name)
		require.EqualValues(t, "", original.Ref)
	})

	t.Run("with values", func(t *testing.T) {
		original := &asciichgolangpublic.GitlabCreateTagOptions{
			Name: "v1.0.0",
			Ref:  "main",
		}
		copy := original.GetDeepCopy()

		require.EqualValues(t, original.Name, copy.Name)
		require.EqualValues(t, original.Ref, copy.Ref)

		// Modify copy
		copy.Name = "v2.0.0"
		copy.Ref = "develop"

		// Original should be unchanged
		require.EqualValues(t, "v1.0.0", original.Name)
		require.EqualValues(t, "main", original.Ref)
	})
}
