package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
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
				require := require.New(t)

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
				require.True(ok)

				require.EqualValues(
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
				require := require.New(t)

				const verbose bool = true

				directory := tempfiles.MustCreateEmptyTemporaryDirectory(verbose)
				defer directory.Delete(verbose)

				require.EqualValues(
					"localhost",
					directory.MustGetHostDescription(),
				)

				repo := MustGetLocalGitReposioryFromDirectory(directory)

				require.EqualValues(
					directory.MustGetPath(),
					repo.MustGetPath(),
				)

				require.EqualValues(
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
				require := require.New(t)

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
					&gitparameteroptions.GitConfigSetOptions{
						Name:    "a",
						Email:   "b@example.net",
						Verbose: true,
					},
				)

				// First commit
				repo.MustCommit(
					&gitparameteroptions.GitCommitOptions{
						Message:    "message 1",
						AllowEmpty: true,
					},
				)

				firstCommit := repo.MustGetCurrentCommit(verbose)
				require.EqualValues("message 1", firstCommit.MustGetCommitMessage())
				require.False(firstCommit.MustHasParentCommit())

				firstCommitDirectParents := firstCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				require.Len(firstCommitDirectParents, 0)

				firstCommitAllParents := firstCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				require.Len(firstCommitAllParents, 0)

				// Second commit
				repo.MustCommit(
					&gitparameteroptions.GitCommitOptions{
						Message:    "message 2",
						AllowEmpty: true,
					},
				)
				secondCommit := repo.MustGetCurrentCommit(verbose)
				require.EqualValues("message 2", secondCommit.MustGetCommitMessage())
				require.True(secondCommit.MustHasParentCommit())

				secondCommitDirectParents := secondCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				require.Len(secondCommitDirectParents, 1)
				require.EqualValues("message 1", secondCommitDirectParents[0].MustGetCommitMessage())

				secondCommitAllParents := secondCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				require.Len(secondCommitAllParents, 1)
				require.EqualValues("message 1", secondCommitAllParents[0].MustGetCommitMessage())

				// Third Commit
				repo.MustCommit(
					&gitparameteroptions.GitCommitOptions{
						Message:    "message 3",
						AllowEmpty: true,
					},
				)
				thirdCommit := repo.MustGetCurrentCommit(verbose)
				require.EqualValues("message 3", thirdCommit.MustGetCommitMessage())
				require.True(thirdCommit.MustHasParentCommit())

				thirdCommitDirectParents := thirdCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				require.Len(thirdCommitDirectParents, 1)
				require.EqualValues("message 2", thirdCommitDirectParents[0].MustGetCommitMessage())

				thirdCommitAllParents := thirdCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				require.Len(thirdCommitAllParents, 2)
				require.EqualValues("message 2", thirdCommitAllParents[0].MustGetCommitMessage())
				require.EqualValues("message 1", thirdCommitAllParents[1].MustGetCommitMessage())
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
				require := require.New(t)

				const verbose bool = true

				repo := TemporaryGitRepositories().MustCreateEmptyTemporaryGitRepository(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						BareRepository:              tt.bareRepository,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					})

				commit := repo.MustGetCurrentCommit(verbose)
				require.EqualValues("Initial empty commit during repo initialization", commit.MustGetCommitMessage())
				require.EqualValues("asciichgolangpublic git repo initializer <asciichgolangpublic@example.net>", commit.MustGetAuthorString())
				require.EqualValues("asciichgolangpublic@example.net", commit.MustGetAuthorEmail())
				require.Greater(commit.MustGetAgeSeconds(), 0.)
				require.Less(commit.MustGetAgeSeconds(), 2.)
				require.False(commit.MustHasParentCommit())
				require.Len(commit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false}), 0)
			},
		)
	}
}
