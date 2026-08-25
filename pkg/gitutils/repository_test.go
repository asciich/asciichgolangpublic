package gitutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/nativegit"
)

func Test_NewGitRepositoryFromDirectory(t *testing.T) {
	t.Run("nativefilesoo", func(t *testing.T) {
		ctx := getCtx()
		tempDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

		_, err = nativegit.InitializeEmptyGitRepository(ctx, tempDir)
		require.NoError(t, err)

		dir, err := nativefilesoo.NewDirectoryByPath(tempDir)
		require.NoError(t, err)

		gitRepo, err := gitutils.NewGitRepositoryFromDirectory(ctx, dir)
		require.NoError(t, err)
		require.NotNil(t, gitRepo)
	})

	t.Run("commandexecutoroo", func(t *testing.T) {
		ctx := getCtx()
		tempDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

		_, err = nativegit.InitializeEmptyGitRepository(ctx, tempDir)
		require.NoError(t, err)

		dir, err := commandexecutorfileoo.NewDirectory(commandexecutorexecoo.Exec(), tempDir)
		require.NoError(t, err)

		gitRepo, err := gitutils.NewGitRepositoryFromDirectory(ctx, dir)
		require.NoError(t, err)
		require.NotNil(t, gitRepo)
	})
}
