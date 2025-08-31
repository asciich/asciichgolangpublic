package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
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
				ctx := getCtx()

				repo, err := TemporaryGitRepositories().CreateEmptyTemporaryGitRepository(
					ctx,
					&parameteroptions.CreateRepositoryOptions{
						BareRepository:              tt.bareRepository,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				localGitRepo, ok := repo.(*LocalGitRepository)
				require.True(t, ok)

				hash, err := repo.GetCurrentCommitHash(ctx)
				require.NoError(t, err)
				require.EqualValues(t, hash, mustutils.Must(localGitRepo.GetCurrentCommitGoGitHash(ctx)).String())
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
				ctx := getCtx()

				directory, err := tempfilesoo.CreateEmptyTemporaryDirectory(ctx)
				require.NoError(t, err)
				defer directory.Delete(ctx, &filesoptions.DeleteOptions{})

				require.EqualValues(t, "localhost", directory.MustGetHostDescription())

				repo, err := GetLocalGitReposioryFromDirectory(directory)
				require.NoError(t, err)

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
				ctx := getCtx()

				repo, err := TemporaryGitRepositories().CreateEmptyTemporaryGitRepository(
					ctx,
					&parameteroptions.CreateRepositoryOptions{
						BareRepository:              tt.bareRepository,
						InitializeWithEmptyCommit:   false,
						InitializeWithDefaultAuthor: true,
					})
				require.NoError(t, err)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				err = repo.SetGitConfig(
					ctx,
					&gitparameteroptions.GitConfigSetOptions{
						Name:  "a",
						Email: "b@example.net",
					},
				)
				require.NoError(t, err)

				// First commit
				_, err = repo.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message:    "message 1",
						AllowEmpty: true,
					},
				)
				require.NoError(t, err)

				firstCommit, err := repo.GetCurrentCommit(ctx)
				require.NoError(t, err)

				commitMessage, err := firstCommit.GetCommitMessage(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "message 1", commitMessage)

				hasParentCommit, err := firstCommit.HasParentCommit(ctx)
				require.NoError(t, err)
				require.False(t, hasParentCommit)

				firstCommitDirectParents, err := firstCommit.GetParentCommits(ctx, &parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				require.NoError(t, err)
				require.Len(t, firstCommitDirectParents, 0)

				firstCommitAllParents, err := firstCommit.GetParentCommits(ctx, &parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				require.NoError(t, err)
				require.Len(t, firstCommitAllParents, 0)

				// Second commit
				_, err = repo.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message:    "message 2",
						AllowEmpty: true,
					},
				)
				require.NoError(t, err)
				secondCommit, err := repo.GetCurrentCommit(ctx)
				require.NoError(t, err)

				commitMessage, err = secondCommit.GetCommitMessage(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "message 2", commitMessage)

				hasParent, err := secondCommit.HasParentCommit(ctx)
				require.NoError(t, err)
				require.True(t, hasParent)

				secondCommitDirectParents, err := secondCommit.GetParentCommits(ctx, &parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				require.NoError(t, err)
				require.Len(t, secondCommitDirectParents, 1)

				commitMessage, err = secondCommitDirectParents[0].GetCommitMessage(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "message 1", commitMessage)

				secondCommitAllParents, err := secondCommit.GetParentCommits(ctx, &parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				require.NoError(t, err)
				require.Len(t, secondCommitAllParents, 1)

				commitMessage, err = secondCommitAllParents[0].GetCommitMessage(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "message 1", commitMessage)

				// Third Commit
				_, err = repo.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message:    "message 3",
						AllowEmpty: true,
					},
				)
				require.NoError(t, err)
				thirdCommit, err := repo.GetCurrentCommit(ctx)
				require.NoError(t, err)

				commitMessage, err = thirdCommit.GetCommitMessage(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "message 3", commitMessage)

				hasParent, err = thirdCommit.HasParentCommit(ctx)
				require.NoError(t, err)
				require.True(t, hasParent)

				thirdCommitDirectParents, err := thirdCommit.GetParentCommits(ctx, &parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				require.NoError(t, err)
				require.Len(t, thirdCommitDirectParents, 1)

				commitMessage, err = thirdCommitDirectParents[0].GetCommitMessage(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "message 2", commitMessage)

				thirdCommitAllParents, err := thirdCommit.GetParentCommits(ctx, &parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: true})
				require.NoError(t, err)
				require.Len(t, thirdCommitAllParents, 2)

				commitMessage, err = thirdCommitAllParents[0].GetCommitMessage(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "message 2", commitMessage)

				commitMessage, err = thirdCommitAllParents[1].GetCommitMessage(ctx)
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
				ctx := getCtx()

				repo, err := TemporaryGitRepositories().CreateEmptyTemporaryGitRepository(
					ctx,
					&parameteroptions.CreateRepositoryOptions{
						BareRepository:              tt.bareRepository,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					})
				require.NoError(t, err)

				commit, err := repo.GetCurrentCommit(ctx)
				require.NoError(t, err)

				commitMessage, err := commit.GetCommitMessage(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "Initial empty commit during repo initialization", commitMessage)

				authorString, err := commit.GetAuthorString(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "asciichgolangpublic git repo initializer <asciichgolangpublic@example.net>", authorString)

				authorEmail, err := commit.GetAuthorEmail(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "asciichgolangpublic@example.net", authorEmail)

				ageSeconds, err := commit.GetAgeSeconds(ctx)
				require.NoError(t, err)
				require.Greater(t, ageSeconds, 0.)
				require.Less(t, ageSeconds, 2.)

				hasParentCommits, err := commit.HasParentCommit(ctx)
				require.NoError(t, err)
				require.False(t, hasParentCommits)

				parentCommits, err := commit.GetParentCommits(ctx, &parameteroptions.GitCommitGetParentsOptions{IncludeParentsOfParents: false})
				require.NoError(t, err)
				require.Len(t, parentCommits, 0)
			},
		)
	}
}
