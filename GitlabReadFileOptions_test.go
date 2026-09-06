package asciichgolangpublic_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	asciichgolangpublic "github.com/asciich/asciichgolangpublic"
)

func TestGitlabReadFileOptions_GetDeepCopy(t *testing.T) {
	t.Run("nil fields", func(t *testing.T) {
		original := &asciichgolangpublic.GitlabReadFileOptions{}
		copy := original.GetDeepCopy()

		require.EqualValues(t, original.Path, copy.Path)
		require.EqualValues(t, original.BranchName, copy.BranchName)

		// Modify copy
		copy.Path = "modified"
		copy.BranchName = "modified"

		// Original should be unchanged
		require.EqualValues(t, "", original.Path)
		require.EqualValues(t, "", original.BranchName)
	})

	t.Run("with values", func(t *testing.T) {
		original := &asciichgolangpublic.GitlabReadFileOptions{
			Path:       "README.md",
			BranchName: "main",
		}
		copy := original.GetDeepCopy()

		require.EqualValues(t, original.Path, copy.Path)
		require.EqualValues(t, original.BranchName, copy.BranchName)

		// Modify copy
		copy.Path = "modified.md"
		copy.BranchName = "develop"

		// Original should be unchanged
		require.EqualValues(t, "README.md", original.Path)
		require.EqualValues(t, "main", original.BranchName)
	})
}
