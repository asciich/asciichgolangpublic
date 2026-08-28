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

func Test_CommitHasParentCommitByCommitHash(t *testing.T) {
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

				// Create and initialize a git repository
				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				// Initialize with an empty commit (first commit)
				err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				// Get the first commit hash
				firstCommitHash, err := gitRepo.GetCurrentCommitHash(ctx)
				require.NoError(t, err)

				// First commit should have no parent
				hasParent, err := gitRepo.CommitHasParentCommitByCommitHash(firstCommitHash)
				require.NoError(t, err)
				require.False(t, hasParent)

				// Create a second commit
				_, err = gitRepo.WriteBytesToFile(ctx, "test.txt", []byte("content"), &filesoptions.WriteOptions{})
				require.NoError(t, err)

				err = gitRepo.AddFileByPath(ctx, "test.txt")
				require.NoError(t, err)

				_, err = gitRepo.Commit(ctx, &gitparameteroptions.GitCommitOptions{
					Message:          "Second commit",
					CommitAllChanges: false,
					AllowEmpty:       false,
				})
				require.NoError(t, err)

				// Get the second commit hash
				secondCommitHash, err := gitRepo.GetCurrentCommitHash(ctx)
				require.NoError(t, err)
				require.NotEqual(t, firstCommitHash, secondCommitHash)

				// Second commit should have a parent
				hasParent, err = gitRepo.CommitHasParentCommitByCommitHash(secondCommitHash)
				require.NoError(t, err)
				require.True(t, hasParent)
			},
		)
	}
}
