package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/testutils"
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
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustInit(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)

				commitToTag := repo.MustGetCurrentCommit(verbose)

				require.EqualValues(
					[]GitTag{},
					commitToTag.MustListTags(verbose),
				)

				// Add a newer commit to validate the given commit is tagged NOT the latest one
				newerCommit := repo.MustCommit(
					&GitCommitOptions{
						AllowEmpty: true,
						Verbose:    verbose,
						Message:    "this commit should not be tagged",
					},
				)

				createdTag := commitToTag.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "first_tag",
						Verbose: verbose,
					},
				)

				require.EqualValues(
					commitToTag.MustGetHash(),
					createdTag.MustGetHash(),
				)

				require.NotEqualValues(
					newerCommit.MustGetHash(),
					createdTag.MustGetHash(),
				)
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
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustInit(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)

				currentCommit := repo.MustGetCurrentCommit(verbose)

				require.EqualValues(
					[]GitTag{},
					currentCommit.MustListTags(verbose),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "first_tag",
						Verbose: verbose,
					},
				)

				require.EqualValues(
					"first_tag",
					currentCommit.MustListTagNames(verbose)[0],
				)
				require.Len(
					currentCommit.MustListTagNames(verbose),
					1,
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "second_tag",
						Verbose: verbose,
					},
				)

				require.EqualValues(
					"first_tag",
					currentCommit.MustListTagNames(verbose)[0],
				)
				require.EqualValues(
					"second_tag",
					currentCommit.MustListTagNames(verbose)[1],
				)
				require.Len(
					currentCommit.MustListTagNames(verbose),
					2,
				)
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
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustInit(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)

				currentCommit := repo.MustGetCurrentCommit(verbose)

				require.EqualValues(
					[]GitTag{},
					currentCommit.MustListTags(verbose),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "first_tag",
						Verbose: verbose,
					},
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "another_tag",
						Verbose: verbose,
					},
				)

				require.EqualValues(
					"v0.1.2",
					currentCommit.MustListVersionTagNames(verbose)[0],
				)
				require.EqualValues(
					"v1.0.0",
					currentCommit.MustListVersionTagNames(verbose)[1],
				)
				require.Len(
					currentCommit.MustListVersionTagNames(verbose),
					2,
				)
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
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustInit(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)

				currentCommit := repo.MustGetCurrentCommit(verbose)

				require.EqualValues(
					[]GitTag{},
					currentCommit.MustListTags(verbose),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "first_tag",
						Verbose: verbose,
					},
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "another_tag",
						Verbose: verbose,
					},
				)

				require.EqualValues(
					"v1.0.0",
					currentCommit.MustGetNewestTagVersion(verbose).MustGetAsString(),
				)
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
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustInit(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)

				currentCommit := repo.MustGetCurrentCommit(verbose)

				require.EqualValues(
					[]GitTag{},
					currentCommit.MustListTags(verbose),
				)
				require.Nil(
					currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "first_tag",
						Verbose: verbose,
					},
				)

				require.Nil(
					currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				require.EqualValues(
					"v1.0.0",
					currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose).MustGetAsString(),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				require.EqualValues(
					"v1.0.0",
					currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose).MustGetAsString(),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "another_tag",
						Verbose: verbose,
					},
				)

				require.EqualValues(
					"v1.0.0",
					currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose).MustGetAsString(),
				)
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
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustInit(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)

				currentCommit := repo.MustGetCurrentCommit(verbose)

				require.False(
					currentCommit.MustHasVersionTag(verbose),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "first_tag",
						Verbose: verbose,
					},
				)

				require.False(
					currentCommit.MustHasVersionTag(verbose),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				require.True(
					currentCommit.MustHasVersionTag(verbose),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				require.True(
					currentCommit.MustHasVersionTag(verbose),
				)
			},
		)
	}
}
