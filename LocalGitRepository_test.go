package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO move to GitRepository_test.go and run for all implementations


// TODO move to GitRepository_test.go and run for all implementations
// TODO move to GitRepository_test.go and run for all implementations
func TestLocalGitRepositoryGetParentCommits(t *testing.T) {
	tests := []struct {
		bareRepository bool
	}{
		{false},
	}
	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := TemporaryDirectories().MustCreateEmptyTemporaryGitRepository(
					&CreateRepositoryOptions{
						Verbose:                     verbose,
						BareRepository:              tt.bareRepository,
						InitializeWithEmptyCommit:   false,
						InitializeWithDefaultAuthor: true,
					})
				defer repo.Delete(verbose)

				repo.MustSetGitConfig(
					&GitConfigSetOptions{
						Name:    "a",
						Email:   "b@example.net",
						Verbose: true,
					},
				)

				// First commit
				repo.MustCommit(
					&GitCommitOptions{
						Message:    "message 1",
						AllowEmpty: true,
					},
				)

				firstCommit := repo.MustGetCurrentCommit()
				assert.EqualValues("message 1", firstCommit.MustGetCommitMessage())
				assert.False(firstCommit.MustHasParentCommit())

				firstCommitDirectParents := firstCommit.MustGetParentCommits(&GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				assert.Len(firstCommitDirectParents, 0)

				firstCommitAllParents := firstCommit.MustGetParentCommits(&GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				assert.Len(firstCommitAllParents, 0)

				// Second commit
				repo.MustCommit(
					&GitCommitOptions{
						Message:    "message 2",
						AllowEmpty: true,
					},
				)
				secondCommit := repo.MustGetCurrentCommit()
				assert.EqualValues("message 2", secondCommit.MustGetCommitMessage())
				assert.True(secondCommit.MustHasParentCommit())

				secondCommitDirectParents := secondCommit.MustGetParentCommits(&GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				assert.Len(secondCommitDirectParents, 1)
				assert.EqualValues("message 1", secondCommitDirectParents[0].MustGetCommitMessage())

				secondCommitAllParents := secondCommit.MustGetParentCommits(&GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				assert.Len(secondCommitAllParents, 1)
				assert.EqualValues("message 1", secondCommitAllParents[0].MustGetCommitMessage())

				// Third Commit
				repo.MustCommit(
					&GitCommitOptions{
						Message:    "message 3",
						AllowEmpty: true,
					},
				)
				thirdCommit := repo.MustGetCurrentCommit()
				assert.EqualValues("message 3", thirdCommit.MustGetCommitMessage())
				assert.True(thirdCommit.MustHasParentCommit())

				thirdCommitDirectParents := thirdCommit.MustGetParentCommits(&GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				assert.Len(thirdCommitDirectParents, 1)
				assert.EqualValues("message 2", thirdCommitDirectParents[0].MustGetCommitMessage())

				thirdCommitAllParents := thirdCommit.MustGetParentCommits(&GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				assert.Len(thirdCommitAllParents, 2)
				assert.EqualValues("message 2", thirdCommitAllParents[0].MustGetCommitMessage())
				assert.EqualValues("message 1", thirdCommitAllParents[1].MustGetCommitMessage())
			},
		)
	}
}
