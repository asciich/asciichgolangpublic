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

func Test_PushToRemote(t *testing.T) {
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

				t.Run("push to remote when already up to date", func(t *testing.T) {
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

					err = cloneRepo.PushToRemote(ctx, "origin")
					require.NoError(t, err)
				})
			},
		)
	}
}


func Test_PushTagsToRemote(t *testing.T) {
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

				t.Run("push tags to remote when already up to date", func(t *testing.T) {
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

					err = cloneRepo.PushTagsToRemote(ctx, "origin")
					require.NoError(t, err)
				})
			},
		)
	}
}
