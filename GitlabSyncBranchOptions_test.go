package asciichgolangpublic_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	asciichgolangpublic "github.com/asciich/asciichgolangpublic"
)

func TestGitlabSyncBranchOptions_GetDeepCopy(t *testing.T) {
	t.Run("nil fields", func(t *testing.T) {
		original := &asciichgolangpublic.GitlabSyncBranchOptions{}
		copy := original.GetDeepCopy()

		require.EqualValues(t, original.TargetBranch, copy.TargetBranch)
		require.EqualValues(t, original.TargetBranchName, copy.TargetBranchName)
		require.EqualValues(t, original.PathsToSync, copy.PathsToSync)

		// Modify copy
		copy.TargetBranchName = "modified"
		copy.PathsToSync = []string{"modified"}

		// Original should be unchanged
		require.EqualValues(t, "", original.TargetBranchName)
		require.Nil(t, original.PathsToSync)
	})

	t.Run("with PathsToSync slice", func(t *testing.T) {
		original := &asciichgolangpublic.GitlabSyncBranchOptions{
			TargetBranchName: "main",
			PathsToSync:      []string{"file1.txt", "file2.txt"},
		}
		copy := original.GetDeepCopy()

		require.EqualValues(t, original.TargetBranchName, copy.TargetBranchName)
		require.EqualValues(t, original.PathsToSync, copy.PathsToSync)

		// Modify copy's slice
		copy.PathsToSync[0] = "modified.txt"
		copy.PathsToSync = append(copy.PathsToSync, "file3.txt")

		// Original should be unchanged (deep copy of slice)
		require.EqualValues(t, []string{"file1.txt", "file2.txt"}, original.PathsToSync)
		require.EqualValues(t, []string{"modified.txt", "file2.txt", "file3.txt"}, copy.PathsToSync)
	})

	t.Run("with TargetBranch object", func(t *testing.T) {
		originalBranch := &asciichgolangpublic.GitlabBranch{}
		err := originalBranch.SetName("main")
		require.NoError(t, err)

		original := &asciichgolangpublic.GitlabSyncBranchOptions{
			TargetBranch: originalBranch,
		}
		copy := original.GetDeepCopy()

		require.NotNil(t, copy.TargetBranch)
		require.EqualValues(t, original.TargetBranch, copy.TargetBranch)

		// Modify copy's branch
		err = copy.TargetBranch.SetName("develop")
		require.NoError(t, err)

		// Original should be unchanged (deep copy of nested object)
		originalName, err := original.TargetBranch.GetName()
		require.NoError(t, err)
		require.EqualValues(t, "main", originalName)

		copyName, err := copy.TargetBranch.GetName()
		require.NoError(t, err)
		require.EqualValues(t, "develop", copyName)
	})

	t.Run("with all fields set", func(t *testing.T) {
		originalBranch := &asciichgolangpublic.GitlabBranch{}
		err := originalBranch.SetName("main")
		require.NoError(t, err)

		original := &asciichgolangpublic.GitlabSyncBranchOptions{
			TargetBranch:     originalBranch,
			TargetBranchName: "main",
			PathsToSync:      []string{"file1.txt", "file2.txt"},
		}
		copy := original.GetDeepCopy()

		// Verify all fields are copied
		require.EqualValues(t, original.TargetBranchName, copy.TargetBranchName)
		require.EqualValues(t, original.PathsToSync, copy.PathsToSync)
		require.NotNil(t, copy.TargetBranch)

		// Modify copy
		copy.TargetBranchName = "develop"
		copy.PathsToSync[0] = "modified.txt"
		err = copy.TargetBranch.SetName("develop")
		require.NoError(t, err)

		// Original should be unchanged
		require.EqualValues(t, "main", original.TargetBranchName)
		require.EqualValues(t, []string{"file1.txt", "file2.txt"}, original.PathsToSync)
		originalName, err := original.TargetBranch.GetName()
		require.NoError(t, err)
		require.EqualValues(t, "main", originalName)
	})
}
