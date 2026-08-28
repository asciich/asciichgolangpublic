package gitutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func Test_ListTagNames(t *testing.T) {
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
				ctx := context.Background()

				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				tagNames, err := gitRepo.ListTagNames(ctx)
				require.NoError(t, err)
				require.Len(t, tagNames, 0)

				_, err = gitRepo.CreateTag(ctx, &gitparameteroptions.GitRepositoryCreateTagOptions{
					TagName: "v1.0.0",
				})
				require.NoError(t, err)

				tagNames, err = gitRepo.ListTagNames(ctx)
				require.NoError(t, err)
				require.Len(t, tagNames, 1)
				require.Contains(t, tagNames, "v1.0.0")
			},
		)
	}
}

func Test_ListTags(t *testing.T) {
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
				ctx := context.Background()

				tempDir, err := tempfiles.CreateTempDir(ctx)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

				gitRepo := getGitRepoByImplementationNameAndPath(t, tt.implementationName, tempDir)

				err = gitRepo.CreateAndInit(ctx, &parameteroptions.CreateRepositoryOptions{
					InitializeWithEmptyCommit:   true,
					InitializeWithDefaultAuthor: true,
				})
				require.NoError(t, err)

				tags, err := gitRepo.ListTags(ctx)
				require.NoError(t, err)
				require.Len(t, tags, 0)

				_, err = gitRepo.CreateTag(ctx, &gitparameteroptions.GitRepositoryCreateTagOptions{
					TagName: "v1.0.0",
				})
				require.NoError(t, err)

				tags, err = gitRepo.ListTags(ctx)
				require.NoError(t, err)
				require.Len(t, tags, 1)

				name, err := tags[0].GetName()
				require.NoError(t, err)
				require.Equal(t, "v1.0.0", name)
			},
		)
	}
}

func Test_ListTagsForCommitHash(t *testing.T) {
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
				ctx := context.Background()

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

				tags, err := gitRepo.ListTagsForCommitHash(ctx, currentHash)
				require.NoError(t, err)
				require.Len(t, tags, 0)

				_, err = gitRepo.CreateTag(ctx, &gitparameteroptions.GitRepositoryCreateTagOptions{
					TagName: "v1.0.0",
				})
				require.NoError(t, err)

				tags, err = gitRepo.ListTagsForCommitHash(ctx, currentHash)
				require.NoError(t, err)
				require.Len(t, tags, 1)

				name, err := tags[0].GetName()
				require.NoError(t, err)
				require.Equal(t, "v1.0.0", name)
			},
		)
	}
}

func Test_CreateTag(t *testing.T) {
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
				ctx := context.Background()

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

				createdTag, err := gitRepo.CreateTag(ctx, &gitparameteroptions.GitRepositoryCreateTagOptions{
					TagName: "v1.0.0",
				})
				require.NoError(t, err)
				require.NotNil(t, createdTag)

				name, err := createdTag.GetName()
				require.NoError(t, err)
				require.Equal(t, "v1.0.0", name)

				hash, err := createdTag.GetHash(ctx)
				require.NoError(t, err)
				require.Equal(t, currentHash, hash)
			},
		)
	}
}
