package gitutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/nativegit"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_CloneRepositoryByPathOrUrl(t *testing.T) {
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

				// Create and initialize an empty git repository to clone from.
				originDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, originDir, &filesoptions.DeleteOptions{})

				_, err = nativegit.InitializeEmptyGitRepository(ctx, originDir)
				require.NoError(t, err)

				// Create a target directory to clone into.
				cloneDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, cloneDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, cloneDir)

				err = gitRepo.CloneRepositoryByPathOrUrl(ctx, originDir)
				require.NoError(t, err)

				isGitRepo, err := gitRepo.IsGitRepository(ctx)
				require.NoError(t, err)
				require.True(t, isGitRepo)
			},
		)
	}
}

func Test_IsInitialized(t *testing.T) {
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

					isInitialized, err := gitRepo.IsInitialized(ctx)
					require.NoError(t, err)
					require.True(t, isInitialized)
				})

				t.Run("empty directory returns false", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					isInitialized, err := gitRepo.IsInitialized(ctx)
					require.NoError(t, err)
					require.False(t, isInitialized)
				})
			},
		)
	}
}

func Test_CloneRepository(t *testing.T) {
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

				// Create and initialize an empty git repository to clone from.
				originDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, originDir, &filesoptions.DeleteOptions{})

				_, err = nativegit.InitializeEmptyGitRepository(ctx, originDir)
				require.NoError(t, err)

				// Create a target directory to clone into.
				cloneDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, cloneDir, &filesoptions.DeleteOptions{})

				originRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, originDir)
				cloneRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, cloneDir)

				err = cloneRepo.CloneRepository(ctx, originRepo)
				require.NoError(t, err)

				isGitRepo, err := cloneRepo.IsGitRepository(ctx)
				require.NoError(t, err)
				require.True(t, isGitRepo)
			},
		)
	}

	t.Run("nil repository returns error", func(t *testing.T) {
		ctx := getCtx()

		tempDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

		gitRepo := getGitRepoByImplementationNameAndPath(t, "commandexecutorgitoo", tempDir)

		err = gitRepo.CloneRepository(ctx, nil)
		require.Error(t, err)
	})

	t.Run("different host descriptions returns error", func(t *testing.T) {
		ctx := getCtx()

		// Create origin repository
		originDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, originDir, &filesoptions.DeleteOptions{})

		_, err = nativegit.InitializeEmptyGitRepository(ctx, originDir)
		require.NoError(t, err)

		// Create clone target directory
		cloneDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, cloneDir, &filesoptions.DeleteOptions{})

		// Use different implementations for origin and clone (both are localhost, so this should work)
		originRepo := getGitRepoByImplementationNameAndPath(t, "nativegitoo", originDir)
		cloneRepo := getGitRepoByImplementationNameAndPath(t, "commandexecutorgitoo", cloneDir)

		// This should work since both are on localhost
		err = cloneRepo.CloneRepository(ctx, originRepo)
		require.NoError(t, err)

		isGitRepo, err := cloneRepo.IsGitRepository(ctx)
		require.NoError(t, err)
		require.True(t, isGitRepo)
	})
}

func Test_CloneToTemporaryRepository(t *testing.T) {
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

				// Create and initialize a git repository with an initial commit
				originDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, originDir, &filesoptions.DeleteOptions{})

				originRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, originDir)

				err = originRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				// Clone to temporary repository
				clonedRepo, err := originRepo.CloneToTemporaryRepository(ctx)
				require.NoError(t, err)
				require.NotNil(t, clonedRepo)

				// Clean up the cloned repository
				defer func() {
					err := clonedRepo.Delete(ctx, &filesoptions.DeleteOptions{})
					require.NoError(t, err)
				}()

				// Verify the cloned repository is a valid git repository
				isGitRepo, err := clonedRepo.IsGitRepository(ctx)
				require.NoError(t, err)
				require.True(t, isGitRepo)

				// Verify the cloned repository has the origin remote
				remoteConfigs, err := clonedRepo.GetRemoteConfigs(ctx)
				require.NoError(t, err)
				require.Len(t, remoteConfigs, 1)

				remoteName, err := remoteConfigs[0].GetRemoteName()
				require.NoError(t, err)
				require.EqualValues(t, "origin", remoteName)

				remoteUrl, err := remoteConfigs[0].GetUrlFetch()
				require.NoError(t, err)
				require.EqualValues(t, originDir, remoteUrl)
			},
		)
	}
}
