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

func Test_ListBranchNames(t *testing.T) {
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

				t.Run("initialized repo with initial commit has default branch", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					branchNames, err := gitRepo.ListBranchNames(ctx)
					require.NoError(t, err)
					require.Len(t, branchNames, 1)
				})

				t.Run("list multiple branches", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					err = gitRepo.CreateBranch(ctx, &parameteroptions.CreateBranchOptions{Name: "feature-a"})
					require.NoError(t, err)

					err = gitRepo.CreateBranch(ctx, &parameteroptions.CreateBranchOptions{Name: "feature-b"})
					require.NoError(t, err)

					branchNames, err := gitRepo.ListBranchNames(ctx)
					require.NoError(t, err)
					require.Len(t, branchNames, 3)
					require.Contains(t, branchNames, "feature-a")
					require.Contains(t, branchNames, "feature-b")
				})
			},
		)
	}
}

func Test_CreateBranch(t *testing.T) {
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

				t.Run("create branch in initialized repo", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					err = gitRepo.CreateBranch(ctx, &parameteroptions.CreateBranchOptions{Name: "test-branch"})
					require.NoError(t, err)

					branchExists, err := gitRepo.BranchByNameExists(ctx, "test-branch")
					require.NoError(t, err)
					require.True(t, branchExists)
				})

				t.Run("nil options returns error", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					err = gitRepo.CreateBranch(ctx, nil)
					require.Error(t, err)
				})
			},
		)
	}
}

func Test_CheckoutBranchByName(t *testing.T) {
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

				t.Run("checkout existing branch", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					err = gitRepo.CreateBranch(ctx, &parameteroptions.CreateBranchOptions{Name: "feature-x"})
					require.NoError(t, err)

					err = gitRepo.CheckoutBranchByName(ctx, "feature-x")
					require.NoError(t, err)

					currentBranch, err := gitRepo.GetCurrentBranchName(ctx)
					require.NoError(t, err)
					require.EqualValues(t, "feature-x", currentBranch)
				})

				t.Run("checkout non existing branch returns error", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					err = gitRepo.CheckoutBranchByName(ctx, "non-existing-branch")
					require.Error(t, err)
				})
			},
		)
	}
}

func Test_GetCurrentBranchName(t *testing.T) {
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

				t.Run("returns default branch after init", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					currentBranch, err := gitRepo.GetCurrentBranchName(ctx)
					require.NoError(t, err)
					require.NotEmpty(t, currentBranch)
				})

				t.Run("returns checked out branch name", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					err = gitRepo.CreateBranch(ctx, &parameteroptions.CreateBranchOptions{Name: "my-branch"})
					require.NoError(t, err)

					err = gitRepo.CheckoutBranchByName(ctx, "my-branch")
					require.NoError(t, err)

					currentBranch, err := gitRepo.GetCurrentBranchName(ctx)
					require.NoError(t, err)
					require.EqualValues(t, "my-branch", currentBranch)
				})
			},
		)
	}
}
func Test_DeleteBranchByName(t *testing.T) {
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

				t.Run("delete existing branch", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					defaultBranch, err := gitRepo.GetCurrentBranchName(ctx)
					require.NoError(t, err)

					err = gitRepo.CreateBranch(ctx, &parameteroptions.CreateBranchOptions{Name: "to-delete"})
					require.NoError(t, err)

					err = gitRepo.CheckoutBranchByName(ctx, defaultBranch)
					require.NoError(t, err)

					branchExists, err := gitRepo.BranchByNameExists(ctx, "to-delete")
					require.NoError(t, err)
					require.True(t, branchExists)

					err = gitRepo.DeleteBranchByName(ctx, "to-delete")
					require.NoError(t, err)

					branchExists, err = gitRepo.BranchByNameExists(ctx, "to-delete")
					require.NoError(t, err)
					require.False(t, branchExists)
				})

				t.Run("delete non existing branch is idempotent", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					err = gitRepo.DeleteBranchByName(ctx, "non-existing-branch")
					require.NoError(t, err)
				})

				t.Run("delete branch twice is idempotent", func(t *testing.T) {
					tempDir, err := tempfiles.CreateTempDir(ctx)
					require.NoError(t, err)
					defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

					gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

					err = gitRepo.Init(ctx, &parameteroptions.CreateRepositoryOptions{
						InitializeWithDefaultAuthor: true,
						InitializeWithEmptyCommit:   true,
					})
					require.NoError(t, err)

					defaultBranch, err := gitRepo.GetCurrentBranchName(ctx)
					require.NoError(t, err)

					err = gitRepo.CreateBranch(ctx, &parameteroptions.CreateBranchOptions{Name: "ephemeral"})
					require.NoError(t, err)

					err = gitRepo.CheckoutBranchByName(ctx, defaultBranch)
					require.NoError(t, err)

					err = gitRepo.DeleteBranchByName(ctx, "ephemeral")
					require.NoError(t, err)

					err = gitRepo.DeleteBranchByName(ctx, "ephemeral")
					require.NoError(t, err)

					branchExists, err := gitRepo.BranchByNameExists(ctx, "ephemeral")
					require.NoError(t, err)
					require.False(t, branchExists)
				})
			},
		)
	}
}
