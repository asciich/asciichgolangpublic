package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalGitRepositoryInit(t *testing.T) {

	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				emptyDir := TemporaryDirectories().MustCreateEmptyTemporaryDirectory(verbose)
				defer emptyDir.MustDelete(verbose)

				repo := MustGetLocalGitReposioryFromDirectory(emptyDir)

				assert.False(repo.MustIsInitialized(verbose))

				for i := 0; i < 3; i++ {
					repo.MustInit(&CreateRepositoryOptions{Verbose: verbose})
					assert.True(repo.MustIsInitialized(verbose))
				}
			},
		)
	}
}

func TestLocalGitRepositoryHasUncommittedChanges(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := GitRepositories().MustCreateTemporaryInitializedRepository(
					&CreateRepositoryOptions{
						Verbose: verbose,
					},
				)
				defer repo.MustDelete(verbose)

				assert.False(repo.MustHasUncommittedChanges())
				assert.True(repo.MustHasNoUncommittedChanges())

				repo.CreateFileInDirectory(verbose, "hello.txt")
				assert.True(repo.MustHasUncommittedChanges())
				assert.False(repo.MustHasNoUncommittedChanges())
			},
		)
	}
}

func TestLocalGitRepositoryIsBareRepository(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := GitRepositories().MustCreateTemporaryInitializedRepository(
					&CreateRepositoryOptions{
						Verbose:        verbose,
						BareRepository: false,
					},
				)
				defer repo.MustDelete(verbose)

				assert.False(repo.MustIsBareRepository(verbose))

				repo_bare := GitRepositories().MustCreateTemporaryInitializedRepository(
					&CreateRepositoryOptions{
						Verbose:        verbose,
						BareRepository: true,
					},
				)
				defer repo.MustDelete(verbose)

				assert.True(repo_bare.MustIsBareRepository(verbose))
			},
		)
	}
}

func TestLocalGitRepositoryPullAndPush(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				upstreamRepo := GitRepositories().MustCreateTemporaryInitializedRepository(
					&CreateRepositoryOptions{
						Verbose:                   verbose,
						BareRepository:            true,
						InitializeWithEmptyCommit: true,
					},
				)
				defer upstreamRepo.MustDelete(verbose)

				clonedRepo := GitRepositories().MustCloneToTemporaryDirectory(upstreamRepo.MustGetLocalPath(), verbose)
				defer clonedRepo.MustDelete(verbose)
				clonedRepo.MustSetGitConfig(
					&GitConfigSetOptions{
						Name:  "Test User",
						Email: "user@example.com",
					},
				)

				clonedRepo2 := GitRepositories().MustCloneToTemporaryDirectory(upstreamRepo.MustGetLocalPath(), verbose)
				defer clonedRepo2.MustDelete(verbose)
				clonedRepo2.MustSetGitConfig(
					&GitConfigSetOptions{
						Name:  "Test User2",
						Email: "user2@example.com",
					},
				)

				assert.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(),
					clonedRepo.MustGetCurrentCommitHash(),
				)
				assert.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(),
					clonedRepo2.MustGetCurrentCommitHash(),
				)

				fileName := "abc.txt"
				clonedRepo2.MustCreateFileInDirectory(verbose, fileName)
				clonedRepo2.MustAdd(fileName)
				clonedRepo2.MustCommit(
					&GitCommitOptions{
						Message: "another commit",
						Verbose: verbose,
					},
				)

				assert.NotEqualValues(
					upstreamRepo.MustGetCurrentCommitHash(),
					clonedRepo2.MustGetCurrentCommitHash(),
				)

				assert.NotEqualValues(
					clonedRepo.MustGetCurrentCommitHash(),
					clonedRepo2.MustGetCurrentCommitHash(),
				)

				clonedRepo2.MustPush(verbose)
				assert.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(),
					clonedRepo2.MustGetCurrentCommitHash(),
				)
				assert.NotEqualValues(
					upstreamRepo.MustGetCurrentCommitHash(),
					clonedRepo.MustGetCurrentCommitHash(),
				)

				clonedRepo.MustPull(verbose)
				assert.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(),
					clonedRepo2.MustGetCurrentCommitHash(),
				)
				assert.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(),
					clonedRepo.MustGetCurrentCommitHash(),
				)
			},
		)
	}
}

func TestLocalGitRepositoryGetRootDirectory(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}
	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				bareRepo := TemporaryDirectories().MustCreateEmptyTemporaryGitRepository(
					&CreateRepositoryOptions{
						Verbose:        verbose,
						BareRepository: true,
					})
				nonBareRepo := TemporaryDirectories().MustCreateEmptyTemporaryGitRepository(
					&CreateRepositoryOptions{
						Verbose:        verbose,
						BareRepository: false,
					})

				assert.EqualValues(
					bareRepo.MustGetLocalPath(),
					bareRepo.MustGetRootDirectoryPath(verbose),
				)
				assert.EqualValues(
					nonBareRepo.MustGetLocalPath(),
					nonBareRepo.MustGetRootDirectoryPath(verbose),
				)
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
			MustFormatAsTestname(tt),
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

				commit := repo.MustGetCurrentCommit()
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
