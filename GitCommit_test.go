package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustInit(
					&CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)

				commitToTag := repo.MustGetCurrentCommit()

				assert.EqualValues(
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

				assert.EqualValues(
					commitToTag.MustGetHash(),
					createdTag.MustGetHash(),
				)

				assert.NotEqualValues(
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustInit(
					&CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)

				currentCommit := repo.MustGetCurrentCommit()

				assert.EqualValues(
					[]GitTag{},
					currentCommit.MustListTags(verbose),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "first_tag",
						Verbose: verbose,
					},
				)

				assert.EqualValues(
					"first_tag",
					currentCommit.MustListTagNames(verbose)[0],
				)
				assert.Len(
					currentCommit.MustListTagNames(verbose),
					1,
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "second_tag",
						Verbose: verbose,
					},
				)

				assert.EqualValues(
					"first_tag",
					currentCommit.MustListTagNames(verbose)[0],
				)
				assert.EqualValues(
					"second_tag",
					currentCommit.MustListTagNames(verbose)[1],
				)
				assert.Len(
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustInit(
					&CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)

				currentCommit := repo.MustGetCurrentCommit()

				assert.EqualValues(
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

				assert.EqualValues(
					"v0.1.2",
					currentCommit.MustListVersionTagNames(verbose)[0],
				)
				assert.EqualValues(
					"v1.0.0",
					currentCommit.MustListVersionTagNames(verbose)[1],
				)
				assert.Len(
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustInit(
					&CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)

				currentCommit := repo.MustGetCurrentCommit()

				assert.EqualValues(
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

				assert.EqualValues(
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustInit(
					&CreateRepositoryOptions{
						Verbose:                     verbose,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)

				currentCommit := repo.MustGetCurrentCommit()

				assert.EqualValues(
					[]GitTag{},
					currentCommit.MustListTags(verbose),
				)
				assert.Nil(
					currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "first_tag",
						Verbose: verbose,
					},
				)

				assert.Nil(
					currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				assert.EqualValues(
					"v1.0.0",
					currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose).MustGetAsString(),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				assert.EqualValues(
					"v1.0.0",
					currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose).MustGetAsString(),
				)

				currentCommit.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "another_tag",
						Verbose: verbose,
					},
				)

				assert.EqualValues(
					"v1.0.0",
					currentCommit.MustGetNewestTagVersionOrNilIfUnset(verbose).MustGetAsString(),
				)
			},
		)
	}
}
