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
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
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

func Test_IsGitRepository(t *testing.T) {
	var tests = []struct {
		implementationName string
	}{
		{"nativegitoo"},
		{"commandexecutorgitoo"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				t.Run("initialized repo returns true", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					_, err = nativegit.InitializeEmptyGitRepository(ctx, tempDir)
					require.NoError(t, err)

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					isGitRepo, err := gitRepo.IsGitRepository(ctx)
					require.NoError(t, err)
					require.True(t, isGitRepo)
				})

				t.Run("empty directory returns false", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					isGitRepo, err := gitRepo.IsGitRepository(ctx)
					require.NoError(t, err)
					require.False(t, isGitRepo)
				})
			},
		)
	}
}
