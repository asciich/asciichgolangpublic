package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tempfiles"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestGetCurrentCommitGoGitHash(t *testing.T) {
	tests := []struct {
		bareRepository bool
	}{
		{false},
	}
	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := TemporaryGitRepositories().MustCreateEmptyTemporaryGitRepository(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						BareRepository:              tt.bareRepository,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					})
				defer repo.Delete(verbose)

				localGitRepo, ok := repo.(*LocalGitRepository)
				assert.True(ok)

				assert.EqualValues(
					repo.MustGetCurrentCommitHash(verbose),
					localGitRepo.MustGetCurrentCommitGoGitHash(verbose).String(),
				)
			},
		)
	}
}

func TestLocalGitRepository_GetLocalGitReposioryFromDirectory(t *testing.T) {
	tests := []struct {
		bareRepository bool
	}{
		{false},
	}
	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				directory := tempfiles.MustCreateEmptyTemporaryDirectory(verbose)
				defer directory.Delete(verbose)

				assert.EqualValues(
					"localhost",
					directory.MustGetHostDescription(),
				)

				repo := MustGetLocalGitReposioryFromDirectory(directory)

				assert.EqualValues(
					directory.MustGetPath(),
					repo.MustGetPath(),
				)

				assert.EqualValues(
					"localhost",
					repo.MustGetHostDescription(),
				)
			},
		)
	}
}

// TODO move to GitRepository_test.go and run for all implementations
func TestLocalGitRepositoryGetParentCommits(t *testing.T) {
	tests := []struct {
		bareRepository bool
	}{
		{false},
	}
	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := TemporaryGitRepositories().MustCreateEmptyTemporaryGitRepository(
					&parameteroptions.CreateRepositoryOptions{
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

				firstCommit := repo.MustGetCurrentCommit(verbose)
				assert.EqualValues("message 1", firstCommit.MustGetCommitMessage())
				assert.False(firstCommit.MustHasParentCommit())

				firstCommitDirectParents := firstCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				assert.Len(firstCommitDirectParents, 0)

				firstCommitAllParents := firstCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				assert.Len(firstCommitAllParents, 0)

				// Second commit
				repo.MustCommit(
					&GitCommitOptions{
						Message:    "message 2",
						AllowEmpty: true,
					},
				)
				secondCommit := repo.MustGetCurrentCommit(verbose)
				assert.EqualValues("message 2", secondCommit.MustGetCommitMessage())
				assert.True(secondCommit.MustHasParentCommit())

				secondCommitDirectParents := secondCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				assert.Len(secondCommitDirectParents, 1)
				assert.EqualValues("message 1", secondCommitDirectParents[0].MustGetCommitMessage())

				secondCommitAllParents := secondCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				assert.Len(secondCommitAllParents, 1)
				assert.EqualValues("message 1", secondCommitAllParents[0].MustGetCommitMessage())

				// Third Commit
				repo.MustCommit(
					&GitCommitOptions{
						Message:    "message 3",
						AllowEmpty: true,
					},
				)
				thirdCommit := repo.MustGetCurrentCommit(verbose)
				assert.EqualValues("message 3", thirdCommit.MustGetCommitMessage())
				assert.True(thirdCommit.MustHasParentCommit())

				thirdCommitDirectParents := thirdCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				assert.Len(thirdCommitDirectParents, 1)
				assert.EqualValues("message 2", thirdCommitDirectParents[0].MustGetCommitMessage())

				thirdCommitAllParents := thirdCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				assert.Len(thirdCommitAllParents, 2)
				assert.EqualValues("message 2", thirdCommitAllParents[0].MustGetCommitMessage())
				assert.EqualValues("message 1", thirdCommitAllParents[1].MustGetCommitMessage())
			},
		)
	}
}

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

				repo := TemporaryGitRepositories().MustCreateEmptyTemporaryGitRepository(
					&parameteroptions.CreateRepositoryOptions{
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
				assert.Len(commit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false}), 0)
			},
		)
	}
}
