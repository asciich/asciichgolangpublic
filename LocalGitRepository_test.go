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
				require.True(t, ok)

				hash, err := repo.GetCurrentCommitHash(verbose)
				require.NoError(t, err)
				require.EqualValues(t, hash, localGitRepo.MustGetCurrentCommitGoGitHash(verbose).String())
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
				const verbose bool = true

				directory := tempfiles.MustCreateEmptyTemporaryDirectory(verbose)
				defer directory.Delete(verbose)

				require.EqualValues(t, "localhost", directory.MustGetHostDescription())

				repo := MustGetLocalGitReposioryFromDirectory(directory)

				repoPath, err := repo.GetPath()
				require.NoError(t, err)
				require.EqualValues(t, directory.MustGetPath(), repoPath)

				hostDescription, err := repo.GetHostDescription()
				require.NoError(t, err)
				require.EqualValues(t, "localhost", hostDescription)
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
				const verbose bool = true

				repo := TemporaryGitRepositories().MustCreateEmptyTemporaryGitRepository(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						BareRepository:              tt.bareRepository,
						InitializeWithEmptyCommit:   false,
						InitializeWithDefaultAuthor: true,
					})
				defer repo.Delete(verbose)

				err := repo.SetGitConfig(
					&gitparameteroptions.GitConfigSetOptions{
						Name:    "a",
						Email:   "b@example.net",
						Verbose: true,
					},
				)
				require.NoError(t, err)

				// First commit
				_, err = repo.Commit(
					&gitparameteroptions.GitCommitOptions{
						Message:    "message 1",
						AllowEmpty: true,
					},
				)
				require.NoError(t, err)

				firstCommit, err := repo.GetCurrentCommit(verbose)
				require.NoError(t, err)
				require.EqualValues(t, "message 1", firstCommit.MustGetCommitMessage())
				require.False(t, firstCommit.MustHasParentCommit())

				firstCommitDirectParents := firstCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				require.Len(t, firstCommitDirectParents, 0)

				firstCommitAllParents := firstCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				require.Len(t, firstCommitAllParents, 0)

				// Second commit
				_, err = repo.Commit(
					&gitparameteroptions.GitCommitOptions{
						Message:    "message 2",
						AllowEmpty: true,
					},
				)
				secondCommit, err := repo.GetCurrentCommit(verbose)
				require.NoError(t, err)
				require.EqualValues(t, "message 2", secondCommit.MustGetCommitMessage())
				require.True(t, secondCommit.MustHasParentCommit())

				secondCommitDirectParents := secondCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				require.Len(t, secondCommitDirectParents, 1)
				require.EqualValues(t, "message 1", secondCommitDirectParents[0].MustGetCommitMessage())

				secondCommitAllParents := secondCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				require.Len(t, secondCommitAllParents, 1)
				require.EqualValues(t, "message 1", secondCommitAllParents[0].MustGetCommitMessage())

				// Third Commit
				_, err = repo.Commit(
					&gitparameteroptions.GitCommitOptions{
						Message:    "message 3",
						AllowEmpty: true,
					},
				)
				thirdCommit, err := repo.GetCurrentCommit(verbose)
				require.NoError(t, err)
				require.EqualValues(t, "message 3", thirdCommit.MustGetCommitMessage())
				require.True(t, thirdCommit.MustHasParentCommit())

				thirdCommitDirectParents := thirdCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				require.Len(t, thirdCommitDirectParents, 1)
				require.EqualValues(t, "message 2", thirdCommitDirectParents[0].MustGetCommitMessage())

				thirdCommitAllParents := thirdCommit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				require.Len(t, thirdCommitAllParents, 2)
				require.EqualValues(t, "message 2", thirdCommitAllParents[0].MustGetCommitMessage())
				require.EqualValues(t, "message 1", thirdCommitAllParents[1].MustGetCommitMessage())
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
				const verbose bool = true

				repo := TemporaryGitRepositories().MustCreateEmptyTemporaryGitRepository(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:                     verbose,
						BareRepository:              tt.bareRepository,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					})

				commit, err := repo.GetCurrentCommit(verbose)
				require.NoError(t, err)
				require.EqualValues(t, "Initial empty commit during repo initialization", commit.MustGetCommitMessage())
				require.EqualValues(t, "asciichgolangpublic git repo initializer <asciichgolangpublic@example.net>", commit.MustGetAuthorString())
				require.EqualValues(t, "asciichgolangpublic@example.net", commit.MustGetAuthorEmail())
				require.Greater(t, commit.MustGetAgeSeconds(), 0.)
				require.Less(t, commit.MustGetAgeSeconds(), 2.)
				require.False(t, commit.MustHasParentCommit())
				require.Len(t, commit.MustGetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false}), 0)
			},
		)
	}
}
