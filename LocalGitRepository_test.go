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

				assert.False(repo.MustIsInitialized())

				for i := 0; i < 3; i++ {
					repo.MustInit(&CreateRepositoryOptions{Verbose: verbose})
					assert.True(repo.MustIsInitialized())
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

				repo.CreateFileInDirectory("hello.txt")
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
				clonedRepo2.MustCreateFileInDirectory(fileName)
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

				assert.EqualValues(bareRepo.MustGetLocalPath(), bareRepo.MustGetRootDirectoryPath())
				assert.EqualValues(nonBareRepo.MustGetLocalPath(), nonBareRepo.MustGetRootDirectoryPath())
			},
		)
	}
}

func TestLocalGitRepositoryGetCommit(t *testing.T) {
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
						Verbose:                   verbose,
						BareRepository:            false,
						InitializeWithEmptyCommit: false,
					})
				nonBareRepo := TemporaryDirectories().MustCreateEmptyTemporaryGitRepository(
					&CreateRepositoryOptions{
						Verbose:        verbose,
						BareRepository: false,
					})

				assert.EqualValues(bareRepo.MustGetLocalPath(), bareRepo.MustGetRootDirectoryPath())
				assert.EqualValues(nonBareRepo.MustGetLocalPath(), nonBareRepo.MustGetRootDirectoryPath())
			},
		)
	}
}
