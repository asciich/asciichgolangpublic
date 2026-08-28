package gitutils_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_GetAuthorEmailByCommitHash(t *testing.T) {
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

				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				currentHash, err := gitRepo.GetCurrentCommitHash(ctx)
				require.NoError(t, err)

				authorEmail, err := gitRepo.GetAuthorEmailByCommitHash(currentHash)
				require.NoError(t, err)
				require.NotEmpty(t, authorEmail)
				require.Contains(t, authorEmail, "@")
			},
		)
	}
}

func Test_GetAuthorStringByCommitHash(t *testing.T) {
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

				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				currentHash, err := gitRepo.GetCurrentCommitHash(ctx)
				require.NoError(t, err)

				authorString, err := gitRepo.GetAuthorStringByCommitHash(currentHash)
				require.NoError(t, err)
				require.NotEmpty(t, authorString)
			},
		)
	}
}

func Test_GetDirectoryByPath(t *testing.T) {
	var tests = []struct {
		implementationName string
	}{
		{"nativegitoo"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				subDir, err := gitRepo.GetDirectoryByPath(ctx, "subdir")
				require.NoError(t, err)
				require.NotNil(t, subDir)

				path, err := subDir.GetPath()
				require.NoError(t, err)
				require.Contains(t, path, "subdir")
			},
		)
	}
}

func Test_GetCommitAgeDurationByCommitHash(t *testing.T) {
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

				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				currentHash, err := gitRepo.GetCurrentCommitHash(ctx)
				require.NoError(t, err)

				ageDuration, err := gitRepo.GetCommitAgeDurationByCommitHash(currentHash)
				require.NoError(t, err)
				require.NotNil(t, ageDuration)
				require.Less(t, ageDuration.Seconds(), 10.0)
			},
		)
	}
}

func Test_GetCommitAgeSecondsByCommitHash(t *testing.T) {
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

				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				currentHash, err := gitRepo.GetCurrentCommitHash(ctx)
				require.NoError(t, err)

				ageSeconds, err := gitRepo.GetCommitAgeSecondsByCommitHash(currentHash)
				require.NoError(t, err)
				require.Less(t, ageSeconds, 10.0)
			},
		)
	}
}

func Test_GetCommitMessageByCommitHash(t *testing.T) {
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

				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				currentHash, err := gitRepo.GetCurrentCommitHash(ctx)
				require.NoError(t, err)

				commitMessage, err := gitRepo.GetCommitMessageByCommitHash(currentHash)
				require.NoError(t, err)
				require.NotEmpty(t, commitMessage)
			},
		)
	}
}

func Test_GetCommitParentsByCommitHash(t *testing.T) {
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

				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				_, err = gitRepo.WriteStringToFile(ctx, "test.txt", "content", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				err = gitRepo.AddFileByPath(ctx, "test.txt")
				require.NoError(t, err)

				_, err = gitRepo.Commit(ctx, &gitparameteroptions.GitCommitOptions{
					Message: "Second commit",
				})
				require.NoError(t, err)

				currentHash, err := gitRepo.GetCurrentCommitHash(ctx)
				require.NoError(t, err)

				parents, err := gitRepo.GetCommitParentsByCommitHash(ctx, currentHash, &parameteroptions.GitCommitGetParentsOptions{})
				require.NoError(t, err)
				require.Len(t, parents, 1)

				parentHash, err := parents[0].GetHash(ctx)
				require.NoError(t, err)
				require.NotEqual(t, currentHash, parentHash)
			},
		)
	}

	t.Run("root commit has no parents - nativegitoo", func(t *testing.T) {
		ctx := getCtx()

		tempDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

		gitRepo := getGitRepoByImplementationNameAndPath(t, "nativegitoo", tempDir)

		err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
			InitializeWithEmptyCommit:   true,
			InitializeWithDefaultAuthor: true,
		})
		require.NoError(t, err)

		currentHash, err := gitRepo.GetCurrentCommitHash(ctx)
		require.NoError(t, err)

		parents, err := gitRepo.GetCommitParentsByCommitHash(ctx, currentHash, &parameteroptions.GitCommitGetParentsOptions{})
		require.NoError(t, err)
		require.Len(t, parents, 0)
	})
}

func Test_GetCommitTimeByCommitHash(t *testing.T) {
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

				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				currentHash, err := gitRepo.GetCurrentCommitHash(ctx)
				require.NoError(t, err)

				commitTime, err := gitRepo.GetCommitTimeByCommitHash(currentHash)
				require.NoError(t, err)
				require.NotNil(t, commitTime)
				require.Less(t, time.Since(*commitTime).Seconds(), 10.0)
			},
		)
	}
}

func Test_GetCurrentCommit(t *testing.T) {
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

				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				currentCommit, err := gitRepo.GetCurrentCommit(ctx)
				require.NoError(t, err)
				require.NotNil(t, currentCommit)

				hash, err := currentCommit.GetHash(ctx)
				require.NoError(t, err)
				require.NotEmpty(t, hash)
			},
		)
	}
}

func Test_GetGitStatusOutput(t *testing.T) {
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

				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				statusOutput, err := gitRepo.GetGitStatusOutput(ctx)
				require.NoError(t, err)
				require.IsType(t, "", statusOutput)
			},
		)
	}

	t.Run("status shows uncommitted changes - nativegitoo", func(t *testing.T) {
		ctx := getCtx()

		tempDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

		gitRepo := getGitRepoByImplementationNameAndPath(t, "nativegitoo", tempDir)

		err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
			InitializeWithEmptyCommit:   true,
			InitializeWithDefaultAuthor: true,
		})
		require.NoError(t, err)

		_, err = gitRepo.WriteStringToFile(ctx, "uncommitted.txt", "content", &filesoptions.WriteOptions{})
		require.NoError(t, err)

		statusOutput, err := gitRepo.GetGitStatusOutput(ctx)
		require.NoError(t, err)
		require.Contains(t, statusOutput, "uncommitted.txt")
	})
}

func Test_GetHashByTagName(t *testing.T) {
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

				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				currentHash, err := gitRepo.GetCurrentCommitHash(ctx)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(ctx, &gitparameteroptions.GitRepositoryCreateTagOptions{
					TagName: "v1.0.0",
				})
				require.NoError(t, err)

				hash, err := gitRepo.GetHashByTagName("v1.0.0")
				require.NoError(t, err)
				require.Equal(t, currentHash, hash)
			},
		)
	}
}
