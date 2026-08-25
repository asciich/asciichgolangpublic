package nativegit_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/nativegit"
)

func TestInitializeEmptyGitRepository(t *testing.T) {
	t.Run("creates directory if not exists", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "nonexistent", "repo")
		defer os.RemoveAll(dir)

		repo, err := nativegit.InitializeEmptyGitRepository(context.Background(), dir)
		require.NoError(t, err)
		require.NotNil(t, repo)

		_, err = os.Stat(filepath.Join(dir, ".git"))
		require.NoError(t, err)
	})

	t.Run("is idempotent", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "repo")
		defer os.RemoveAll(dir)

		repo1, err := nativegit.InitializeEmptyGitRepository(context.Background(), dir)
		require.NoError(t, err)
		require.NotNil(t, repo1)

		repo2, err := nativegit.InitializeEmptyGitRepository(context.Background(), dir)
		require.NoError(t, err)
		require.NotNil(t, repo2)
	})

	t.Run("returns error for empty path", func(t *testing.T) {
		_, err := nativegit.InitializeEmptyGitRepository(context.Background(), "")
		require.Error(t, err)
	})
}
