package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
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

				commitMessage, err := firstCommit.GetCommitMessage()
				require.NoError(t, err)
				require.EqualValues(t, "message 1", commitMessage)

				hasParentCommit, err := firstCommit.HasParentCommit()
				require.NoError(t, err)
				require.False(t, hasParentCommit)

				firstCommitDirectParents, err := firstCommit.GetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				require.NoError(t, err)
				require.Len(t, firstCommitDirectParents, 0)

				firstCommitAllParents, err := firstCommit.GetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				require.NoError(t, err)
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

				commitMessage, err = secondCommit.GetCommitMessage()
				require.NoError(t, err)
				require.EqualValues(t, "message 2", commitMessage)

				hasParent, err := secondCommit.HasParentCommit()
				require.NoError(t, err)
				require.True(t, hasParent)

				secondCommitDirectParents, err := secondCommit.GetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				require.NoError(t, err)
				require.Len(t, secondCommitDirectParents, 1)

				commitMessage, err = secondCommitDirectParents[0].GetCommitMessage()
				require.NoError(t, err)
				require.EqualValues(t, "message 1", commitMessage)

				secondCommitAllParents, err := secondCommit.GetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				require.NoError(t, err)
				require.Len(t, secondCommitAllParents, 1)

				commitMessage, err = secondCommitAllParents[0].GetCommitMessage()
				require.NoError(t, err)
				require.EqualValues(t, "message 1", commitMessage)

				// Third Commit
				_, err = repo.Commit(
					&gitparameteroptions.GitCommitOptions{
						Message:    "message 3",
						AllowEmpty: true,
					},
				)
				thirdCommit, err := repo.GetCurrentCommit(verbose)
				require.NoError(t, err)

				commitMessage, err = thirdCommit.GetCommitMessage()
				require.NoError(t, err)
				require.EqualValues(t, "message 3", commitMessage)

				hasParent, err = thirdCommit.HasParentCommit()
				require.NoError(t, err)
				require.True(t, hasParent)

				thirdCommitDirectParents, err := thirdCommit.GetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				require.NoError(t, err)
				require.Len(t, thirdCommitDirectParents, 1)

				commitMessage, err = thirdCommitDirectParents[0].GetCommitMessage()
				require.NoError(t, err)
				require.EqualValues(t, "message 2", commitMessage)

				thirdCommitAllParents, err := thirdCommit.GetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				require.NoError(t, err)
				require.Len(t, thirdCommitAllParents, 2)

				commitMessage, err = thirdCommitAllParents[0].GetCommitMessage()
				require.NoError(t, err)
				require.EqualValues(t, "message 2", commitMessage)

				commitMessage, err = thirdCommitAllParents[1].GetCommitMessage()
				require.NoError(t, err)
				require.EqualValues(t, "message 1", commitMessage)
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

				commitMessage, err := commit.GetCommitMessage()
				require.NoError(t, err)
				require.EqualValues(t, "Initial empty commit during repo initialization", commitMessage)

				authorString, err := commit.GetAuthorString()
				require.NoError(t, err)
				require.EqualValues(t, "asciichgolangpublic git repo initializer <asciichgolangpublic@example.net>", authorString)

				authorEmail, err := commit.GetAuthorEmail()
				require.NoError(t, err)
				require.EqualValues(t, "asciichgolangpublic@example.net", authorEmail)

				ageSeconds, err := commit.GetAgeSeconds()
				require.NoError(t, err)
				require.Greater(t, ageSeconds, 0.)
				require.Less(t, ageSeconds, 2.)

				hasParentCommits, err := commit.HasParentCommit()
				require.NoError(t, err)
				require.False(t, hasParentCommits)

				parentCommits, err := commit.GetParentCommits(&parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				require.NoError(t, err)
				require.Len(t, parentCommits, 0)
			},
		)
	}
}
