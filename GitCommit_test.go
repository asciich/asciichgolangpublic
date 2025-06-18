package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
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
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				err := repo.Init(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				commitToTag, err := repo.GetCurrentCommit(verbose)
				require.NoError(t, err)
				require.EqualValues(t, []GitTag{}, commitToTag.MustListTags(verbose))

				// Add a newer commit to validate the given commit is tagged NOT the latest one
				newerCommit, err := repo.Commit(
					&gitparameteroptions.GitCommitOptions{
						AllowEmpty: true,
						Verbose:    verbose,
						Message:    "this commit should not be tagged",
					},
				)
				require.NoError(t, err)

				createdTag := commitToTag.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "first_tag",
						Verbose: verbose,
					},
				)

				require.EqualValues(t, commitToTag.MustGetHash(), createdTag.MustGetHash())
				require.NotEqualValues(t, newerCommit.MustGetHash(), createdTag.MustGetHash())
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
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				err := repo.Init(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				currentCommit, err := repo.GetCurrentCommit(verbose)
				require.NoError(t, err)

				require.EqualValues(t, []GitTag{}, currentCommit.MustListTags(verbose))

				currentCommit.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "first_tag",
						Verbose: verbose,
					},
				)

				require.EqualValues(t, "first_tag", currentCommit.MustListTagNames(verbose)[0])
				require.Len(t, currentCommit.MustListTagNames(verbose), 1)

				currentCommit.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "second_tag",
						Verbose: verbose,
					},
				)

				require.EqualValues(t, "first_tag", currentCommit.MustListTagNames(verbose)[0])
				require.EqualValues(t, "second_tag", currentCommit.MustListTagNames(verbose)[1])
				require.Len(t, currentCommit.MustListTagNames(verbose), 2)
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
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				err := repo.Init(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				currentCommit, err := repo.GetCurrentCommit(verbose)
				require.NoError(t, err)

				require.EqualValues(t, []GitTag{}, currentCommit.MustListTags(verbose))

				currentCommit.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "first_tag",
						Verbose: verbose,
					},
				)

				currentCommit.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				currentCommit.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				currentCommit.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "another_tag",
						Verbose: verbose,
					},
				)

				require.EqualValues(t, "v0.1.2", currentCommit.MustListVersionTagNames(verbose)[0])
				require.EqualValues(t, "v1.0.0", currentCommit.MustListVersionTagNames(verbose)[1])
				require.Len(t, currentCommit.MustListVersionTagNames(verbose), 2)
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
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				err := repo.Init(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				currentCommit, err := repo.GetCurrentCommit(verbose)
				require.NoError(t, err)

				require.EqualValues(t, []GitTag{}, currentCommit.MustListTags(verbose))

				currentCommit.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "first_tag",
						Verbose: verbose,
					},
				)

				currentCommit.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				currentCommit.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				currentCommit.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "another_tag",
						Verbose: verbose,
					},
				)

				require.EqualValues(t, "v1.0.0", mustutils.Must(currentCommit.MustGetNewestTagVersion(verbose).GetAsString()))
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
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				err := repo.Init(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				currentCommit, err := repo.GetCurrentCommit(verbose)
				require.NoError(t, err)

				require.EqualValues(t, []GitTag{}, currentCommit.MustListTags(verbose))
				require.Nil(t, currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose))

				currentCommit.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "first_tag",
						Verbose: verbose,
					},
				)

				require.Nil(t, currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose))

				currentCommit.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				require.EqualValues(t, "v1.0.0", mustutils.Must(currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose).GetAsString()))

				currentCommit.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				require.EqualValues(t, "v1.0.0", mustutils.Must(currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose).GetAsString()))

				currentCommit.MustCreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "another_tag",
						Verbose: verbose,
					},
				)

				require.EqualValues(t, "v1.0.0", mustutils.Must(currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose).GetAsString()))
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
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				err := repo.Init(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				currentCommit, err := repo.GetCurrentCommit(verbose)
				require.NoError(t, err)
				require.False(t, currentCommit.MustHasVersionTag(verbose))

				_, err = currentCommit.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "first_tag",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)
				require.False(t, currentCommit.MustHasVersionTag(verbose))

				_, err = currentCommit.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)
				require.True(t, currentCommit.MustHasVersionTag(verbose))

				_, err = currentCommit.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)
				require.True(t, currentCommit.MustHasVersionTag(verbose))
			},
		)
	}
}
