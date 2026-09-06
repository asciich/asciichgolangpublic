package asciichgolangpublic_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	asciichgolangpublic "github.com/asciich/asciichgolangpublic"
)

func TestGitlabWriteFileOptions_GetDeepCopy(t *testing.T) {
	t.Run("nil fields", func(t *testing.T) {
		original := &asciichgolangpublic.GitlabWriteFileOptions{}
		copy := original.GetDeepCopy()

		require.EqualValues(t, original.Path, copy.Path)
		require.EqualValues(t, original.BranchName, copy.BranchName)
		require.EqualValues(t, original.CommitMessage, copy.CommitMessage)
		require.EqualValues(t, original.Content, copy.Content)

		// Modify copy
		copy.Path = "modified"
		copy.BranchName = "modified"
		copy.CommitMessage = "modified"
		copy.Content = []byte("modified")

		// Original should be unchanged
		require.EqualValues(t, "", original.Path)
		require.EqualValues(t, "", original.BranchName)
		require.EqualValues(t, "", original.CommitMessage)
		require.Nil(t, original.Content)
	})

	t.Run("with Content byte slice", func(t *testing.T) {
		original := &asciichgolangpublic.GitlabWriteFileOptions{
			Path:          "README.md",
			BranchName:    "main",
			CommitMessage: "Add README",
			Content:       []byte("content here"),
		}
		copy := original.GetDeepCopy()

		require.EqualValues(t, original.Path, copy.Path)
		require.EqualValues(t, original.BranchName, copy.BranchName)
		require.EqualValues(t, original.CommitMessage, copy.CommitMessage)
		require.EqualValues(t, original.Content, copy.Content)

		// Modify copy's byte slice
		copy.Content[0] = 'X'
		copy.Content = append(copy.Content, byte('!'))

		// Original should be unchanged (deep copy of byte slice)
		require.EqualValues(t, []byte("content here"), original.Content)
		require.NotEqualValues(t, original.Content, copy.Content)
	})

	t.Run("with empty Content slice", func(t *testing.T) {
		original := &asciichgolangpublic.GitlabWriteFileOptions{
			Path:          "test.txt",
			BranchName:    "main",
			CommitMessage: "Test",
			Content:       []byte{},
		}
		copy := original.GetDeepCopy()

		require.EqualValues(t, original.Content, copy.Content)

		// Modify copy
		copy.Content = append(copy.Content, byte('a'))

		// Original should be unchanged
		require.EqualValues(t, []byte{}, original.Content)
		require.EqualValues(t, []byte{'a'}, copy.Content)
	})
}
