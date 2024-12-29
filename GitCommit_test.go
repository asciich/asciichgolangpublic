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

func TestGitCommit_ListTags(t *testing.T) {
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
					currentCommit.MustListTags(verbose)[0].MustGetName(),
				)
				assert.Len(
					currentCommit.MustListTags(verbose),
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
					currentCommit.MustListTags(verbose)[0].MustGetName(),
				)
				assert.EqualValues(
					"second_tag",
					currentCommit.MustListTags(verbose)[1].MustGetName(),
				)
				assert.Len(
					currentCommit.MustListTags(verbose),
					2,
				)
			},
		)
	}
}
