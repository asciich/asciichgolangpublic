package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestGitCommit_CreateTag(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				err := repo.Init(
					ctx,
					&parameteroptions.CreateRepositoryOptions{
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				commitToTag, err := repo.GetCurrentCommit(ctx)
				require.NoError(t, err)
				tagList, err := commitToTag.ListTags(ctx)
				require.NoError(t, err)
				require.EqualValues(t, []GitTag{}, tagList)

				// Add a newer commit to validate the given commit is tagged NOT the latest one
				newerCommit, err := repo.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						AllowEmpty: true,
						Message:    "this commit should not be tagged",
					},
				)
				require.NoError(t, err)

				createdTag, err := commitToTag.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "first_tag",
					},
				)
				require.NoError(t, err)

				commitToTagHash, err := commitToTag.GetHash(ctx)
				require.NoError(t, err)

				createdTagHash, err := createdTag.GetHash(ctx)
				require.NoError(t, err)

				newerCommitHash, err := newerCommit.GetHash(ctx)
				require.NoError(t, err)

				require.EqualValues(t, commitToTagHash, createdTagHash)
				require.NotEqualValues(t, newerCommitHash, createdTagHash)
			},
		)
	}
}

func TestGitCommit_ListTagsNames(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				err := repo.Init(
					ctx, 
					&parameteroptions.CreateRepositoryOptions{
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				currentCommit, err := repo.GetCurrentCommit(ctx)
				require.NoError(t, err)

				tagList, err := currentCommit.ListTags(ctx)
				require.NoError(t, err)
				require.EqualValues(t, []GitTag{}, tagList)

				_, err = currentCommit.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "first_tag",
					},
				)
				require.NoError(t, err)

				tagNames, err := currentCommit.ListTagNames(ctx)
				require.NoError(t, err)

				require.EqualValues(t, "first_tag", tagNames[0])
				require.Len(t, tagNames, 1)

				_, err = currentCommit.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "second_tag",
					},
				)
				require.NoError(t, err)

				tagNames, err = currentCommit.ListTagNames(ctx)
				require.NoError(t, err)

				require.EqualValues(t, "first_tag", tagNames[0])
				require.EqualValues(t, "second_tag", tagNames[1])
				require.Len(t, tagNames, 2)
			},
		)
	}
}

func TestGitCommit_ListVersionTagNames(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				err := repo.Init(
					ctx,
					&parameteroptions.CreateRepositoryOptions{
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				currentCommit, err := repo.GetCurrentCommit(ctx)
				require.NoError(t, err)

				tagList, err := currentCommit.ListTags(ctx)
				require.NoError(t, err)

				require.EqualValues(t, []GitTag{}, tagList)

				_, err = currentCommit.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "first_tag",
					},
				)
				require.NoError(t, err)

				_, err = currentCommit.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
					},
				)
				require.NoError(t, err)

				_, err = currentCommit.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
					},
				)
				require.NoError(t, err)

				_, err = currentCommit.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "another_tag",
					},
				)
				require.NoError(t, err)

				versionTagNames, err := currentCommit.ListVersionTagNames(ctx)
				require.NoError(t, err)

				require.EqualValues(t, "v0.1.2", versionTagNames[0])
				require.EqualValues(t, "v1.0.0", versionTagNames[1])
				require.Len(t, versionTagNames, 2)
			},
		)
	}
}

func TestGitCommit_GetNewestTagVersion(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				err := repo.Init(
					ctx,
					&parameteroptions.CreateRepositoryOptions{
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				currentCommit, err := repo.GetCurrentCommit(ctx)
				require.NoError(t, err)

				tags, err := currentCommit.ListTags(ctx)
				require.NoError(t, err)

				require.EqualValues(t, []GitTag{}, tags)

				_, err = currentCommit.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "first_tag",
					},
				)
				require.NoError(t, err)

				_, err = currentCommit.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
					},
				)
				require.NoError(t, err)

				_, err = currentCommit.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
					},
				)
				require.NoError(t, err)

				_, err = currentCommit.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "another_tag",
					},
				)
				require.NoError(t, err)

				newestTagVersion, err := currentCommit.GetNewestTagVersionString(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "v1.0.0", newestTagVersion)
			},
		)
	}
}

func TestGitCommit_GetNewestTagVersionOrNilIfUnset(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				err := repo.Init(
					ctx, 
					&parameteroptions.CreateRepositoryOptions{
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				currentCommit, err := repo.GetCurrentCommit(ctx)
				require.NoError(t, err)

				tagList, err := currentCommit.ListTags(ctx)
				require.NoError(t, err)
				require.EqualValues(t, []GitTag{}, tagList)

				newestTagVersion, err := currentCommit.GetNewestTagVersionOrNilIfUnset(ctx)
				require.NoError(t, err)
				require.Nil(t, newestTagVersion)

				_, err = currentCommit.CreateTag(
					ctx, 
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "first_tag",
					},
				)
				require.NoError(t, err)

				newestTagVersion, err = currentCommit.GetNewestTagVersionOrNilIfUnset(ctx)
				require.NoError(t, err)
				require.Nil(t, newestTagVersion)

				_, err = currentCommit.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
					},
				)
				require.NoError(t, err)

				newestTagVersion, err = currentCommit.GetNewestTagVersionOrNilIfUnset(ctx)
				require.NoError(t, err)
				require.NotNil(t, newestTagVersion)

				require.EqualValues(t, "v1.0.0", newestTagVersion.String())

				_, err = currentCommit.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
					},
				)
				require.NoError(t, err)

				newestTagVersion, err = currentCommit.GetNewestTagVersionOrNilIfUnset(ctx)
				require.NoError(t, err)
				require.NotNil(t, newestTagVersion)

				require.EqualValues(t, "v1.0.0", newestTagVersion.String())

				_, err = currentCommit.CreateTag(
					ctx, 
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "another_tag",
					},
				)
				require.NoError(t, err)

				newestTagVersion, err = currentCommit.GetNewestTagVersionOrNilIfUnset(ctx)
				require.NoError(t, err)
				require.NotNil(t, newestTagVersion)

				require.EqualValues(t, "v1.0.0", newestTagVersion.String())
			},
		)
	}
}

func TestGitCommit_HasVersionTag(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				err := repo.Init(
					ctx, 
					&parameteroptions.CreateRepositoryOptions{
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				currentCommit, err := repo.GetCurrentCommit(ctx)
				require.NoError(t, err)
				hasVersionTag, err := currentCommit.HasVersionTag(ctx)
				require.NoError(t, err)
				require.False(t, hasVersionTag)

				_, err = currentCommit.CreateTag(
					ctx, 
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "first_tag",
					},
				)
				require.NoError(t, err)
				hasVersionTag, err = currentCommit.HasVersionTag(ctx)
				require.NoError(t, err)
				require.False(t, hasVersionTag)

				_, err = currentCommit.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
					},
				)
				require.NoError(t, err)
				hasVersionTag, err = currentCommit.HasVersionTag(ctx)
				require.NoError(t, err)
				require.True(t, hasVersionTag)

				_, err = currentCommit.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
					},
				)
				require.NoError(t, err)
				hasVersionTag, err = currentCommit.HasVersionTag(ctx)
				require.NoError(t, err)
				require.True(t, hasVersionTag)
			},
		)
	}
}
