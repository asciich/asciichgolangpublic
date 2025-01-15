package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestLocalGitRepositoryCreateEmptyTemporaryGitRepository(t *testing.T) {
	tests := []struct {
		bareRepository bool
	}{
		{true},
		{false},
	}
	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := TemporaryDirectories().MustCreateEmptyTemporaryGitRepository(
					&CreateRepositoryOptions{
						Verbose:                     verbose,
						BareRepository:              tt.bareRepository,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					})

				commit := repo.MustGetCurrentCommit(verbose)
				assert.EqualValues("Initial empty commit during repo initialization", commit.MustGetCommitMessage())
				assert.EqualValues("asciichgolangpublic git repo initializer <asciichgolangpublic@example.net>", commit.MustGetAuthorString())
				assert.EqualValues("asciichgolangpublic@example.net", commit.MustGetAuthorEmail())
				assert.Greater(commit.MustGetAgeSeconds(), 0.)
				assert.Less(commit.MustGetAgeSeconds(), 1.)
				assert.False(commit.MustHasParentCommit())
				assert.Len(commit.MustGetParentCommits(&GitCommitGetParentsOptions{IncludeParentsOfParents: false}), 0)
			},
		)
	}
}
