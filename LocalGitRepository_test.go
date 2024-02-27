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

				repo := MustGetLocalGitRepositoryByPath(emptyDir.MustGetLocalPath())

				assert.False(repo.MustIsInitialized())

				for i := 0; i < 3; i++ {
					repo.MustInit(verbose)
					assert.True(repo.MustIsInitialized())
				}
			},
		)
	}
}

func TestLocalGitRepositoryHasUncommitedChanges(t *testing.T) {

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

				repo := GitRepositories().MustCreateTemporaryInitializedRepository(verbose)
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

				upstreamRepo := GitRepositories().MustCreateTemporaryInitializedRepository(verbose)
				defer upstreamRepo.MustDelete(verbose)
				upstreamRepo.MustSetGitConfig(
					&GitConfigSetOptions{
						Name:  "Test User",
						Email: "user@example.com",
					},
				)

				upstreamRepo.MustCommit(
					&GitCommitOptions{
						Message:    "inital commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)

				clonedRepo := GitRepositories().MustCloneToTemporaryDirectory(upstreamRepo.MustGetLocalPath(), verbose)
				defer clonedRepo.MustDelete(verbose)
				clonedRepo.MustSetGitConfig(
					&GitConfigSetOptions{
						Name:  "Test User2",
						Email: "user2@example.com",
					},
				)

				assert.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(),
					clonedRepo.MustGetCurrentCommitHash(),
				)

				fileName := "abc.txt"
				upstreamRepo.MustCreateFileInDirectory(fileName)
				upstreamRepo.MustAdd(fileName)
				upstreamRepo.MustCommit(
					&GitCommitOptions{
						Message: "another commit",
						Verbose: verbose,
					},
				)

				assert.NotEqualValues(
					upstreamRepo.MustGetCurrentCommitHash(),
					clonedRepo.MustGetCurrentCommitHash(),
				)

				clonedRepo.MustPull(verbose)
				assert.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(),
					clonedRepo.MustGetCurrentCommitHash(),
				)

				fileName2 := "abc2.txt"
				clonedRepo.MustCreateFileInDirectory(fileName2)
				clonedRepo.MustAdd(fileName2)
				clonedRepo.MustCommit(
					&GitCommitOptions{
						Message: "another commit2",
						Verbose: verbose,
					},
				)

				assert.NotEqualValues(
					upstreamRepo.MustGetCurrentCommitHash(),
					clonedRepo.MustGetCurrentCommitHash(),
				)
			},
		)
	}
}
