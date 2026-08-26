package gitutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/nativegit"
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
