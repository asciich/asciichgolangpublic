package gitutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_Pull(t *testing.T) {
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

				t.Run("pull new commit from origin", func(t *testing.T) {
					// Create and initialize the origin repository.
					originDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, originDir, &filesoptions.DeleteOptions{})

					originRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, originDir)

					err = originRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					// Clone the origin into a second directory.
					cloneDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, cloneDir, &filesoptions.DeleteOptions{})

					cloneRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, cloneDir)

					err = cloneRepo.CloneRepositoryByPathOrUrl(ctx, originDir)
					require.NoError(t, err)

					// Create a new commit in the origin repository.
					_, err = originRepo.WriteStringToFile(ctx, "hello.txt", "hello world", &filesoptions.WriteOptions{})
					require.NoError(t, err)

					err = originRepo.AddFileByPath(ctx, "hello.txt")
					require.NoError(t, err)

					_, err = originRepo.Commit(ctx, &gitparameteroptions.GitCommitOptions{
						Message: "Add hello.txt",
					})
					require.NoError(t, err)

					// Pull from origin into the clone.
					err = cloneRepo.Pull(ctx)
					require.NoError(t, err)

					// Verify the clone has the new file.
					fileExists, err := cloneRepo.FileByPathExists(ctx, "hello.txt")
					require.NoError(t, err)
					require.True(t, fileExists)
				})

				t.Run("pull when already up to date", func(t *testing.T) {
					// Create and initialize the origin repository.
					originDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, originDir, &filesoptions.DeleteOptions{})

					originRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, originDir)

					err = originRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					// Clone the origin into a second directory.
					cloneDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, cloneDir, &filesoptions.DeleteOptions{})

					cloneRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, cloneDir)

					err = cloneRepo.CloneRepositoryByPathOrUrl(ctx, originDir)
					require.NoError(t, err)

					// Pull when there are no new commits.
					err = cloneRepo.Pull(ctx)
					require.NoError(t, err)
				})
			},
		)
	}
}

func Test_Fetch(t *testing.T) {
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

				t.Run("fetch new commit from origin", func(t *testing.T) {
					originDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, originDir, &filesoptions.DeleteOptions{})

					originRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, originDir)

					err = originRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					cloneDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, cloneDir, &filesoptions.DeleteOptions{})

					cloneRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, cloneDir)

					err = cloneRepo.CloneRepositoryByPathOrUrl(ctx, originDir)
					require.NoError(t, err)

					_, err = originRepo.WriteStringToFile(ctx, "hello.txt", "hello world", &filesoptions.WriteOptions{})
					require.NoError(t, err)

					err = originRepo.AddFileByPath(ctx, "hello.txt")
					require.NoError(t, err)

					_, err = originRepo.Commit(ctx, &gitparameteroptions.GitCommitOptions{
						Message: "Add hello.txt",
					})
					require.NoError(t, err)

					err = cloneRepo.Fetch(ctx)
					require.NoError(t, err)

					fileExists, err := cloneRepo.FileByPathExists(ctx, "hello.txt")
					require.NoError(t, err)
					require.False(t, fileExists)
				})

				t.Run("fetch when already up to date", func(t *testing.T) {
					originDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, originDir, &filesoptions.DeleteOptions{})

					originRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, originDir)

					err = originRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					cloneDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, cloneDir, &filesoptions.DeleteOptions{})

					cloneRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, cloneDir)

					err = cloneRepo.CloneRepositoryByPathOrUrl(ctx, originDir)
					require.NoError(t, err)

					err = cloneRepo.Fetch(ctx)
					require.NoError(t, err)
				})
			},
		)
	}
}
