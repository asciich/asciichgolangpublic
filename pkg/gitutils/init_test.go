package gitutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_Init(t *testing.T) {
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

				t.Run("init non bare repository", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{})
					require.NoError(t, err)

					isInitialized, err := gitRepo.IsInitialized(ctx)
					require.NoError(t, err)
					require.True(t, isInitialized)

					isBare, err := gitRepo.IsBareRepository(ctx)
					require.NoError(t, err)
					require.False(t, isBare)
				})

				t.Run("init bare repository", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						BareRepository: true,
					})
					require.NoError(t, err)

					isInitialized, err := gitRepo.IsInitialized(ctx)
					require.NoError(t, err)
					require.True(t, isInitialized)

					isBare, err := gitRepo.IsBareRepository(ctx)
					require.NoError(t, err)
					require.True(t, isBare)
				})

				t.Run("init already initialized repo is idempotent", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{})
					require.NoError(t, err)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{})
					require.NoError(t, err)

					isInitialized, err := gitRepo.IsInitialized(ctx)
					require.NoError(t, err)
					require.True(t, isInitialized)
				})

				t.Run("init with empty commit", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					hasInitialCommit, err := gitRepo.HasInitialCommit(ctx)
					require.NoError(t, err)
					require.True(t, hasInitialCommit)
				})

				t.Run("nil options returns error", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, nil)
					require.Error(t, err)
				})
			},
		)
	}
}
