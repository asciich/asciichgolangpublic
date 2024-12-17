package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getGitRepositoryToTest(implementationName string) (repo GitRepository) {
	const verbose = true

	if implementationName == "localGitRepository" {
		repo = MustGetLocalGitReposioryFromDirectory(
			TemporaryDirectories().MustCreateEmptyTemporaryDirectory(verbose),
		)
	} else if implementationName == "localCommandExecutorRepository" {
		repo = MustGetLocalCommandExecutorGitRepositoryByDirectory(
			TemporaryDirectories().MustCreateEmptyTemporaryDirectory(verbose),
		)
	} else {
		LogFatalWithTracef("unknown implementationName='%s'", implementationName)
	}

	repo.MustInit(
		&CreateRepositoryOptions{
			Verbose:                     verbose,
			InitializeWithEmptyCommit:   true,
			InitializeWithDefaultAuthor: true,
		},
	)

	return repo
}

func TestGitRepository_Init_minimal(t *testing.T) {
	tests := []struct {
		implementationName string
		bareRepository     bool
	}{
		{"localGitRepository", false},
		{"localGitRepository", true},
		{"localCommandExecutorRepository", false},
		{"localCommandExecutorRepository", true},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustDelete(verbose)

				assert.False(repo.MustExists(verbose))
				assert.False(repo.MustIsInitialized(verbose))
				assert.False(repo.MustHasInitialCommit(verbose))

				for i := 0; i < 2; i++ {
					repo.MustInit(
						&CreateRepositoryOptions{
							Verbose:        verbose,
							BareRepository: tt.bareRepository,
						},
					)
					assert.True(repo.MustExists(verbose))
					assert.True(repo.MustIsInitialized(verbose))
					assert.False(repo.MustHasInitialCommit(verbose))
					assert.EqualValues(
						tt.bareRepository,
						repo.MustIsBareRepository(verbose),
					)
				}
			},
		)
	}
}

func TestGitRepository_Init(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				for i := 0; i < 2; i++ {
					repo.MustDelete(verbose)
					assert.False(repo.MustExists(verbose))
					assert.False(repo.MustIsInitialized(verbose))
					assert.False(repo.MustHasInitialCommit(verbose))
				}

				for i := 0; i < 2; i++ {
					repo.MustCreate(verbose)
					assert.True(repo.MustExists(verbose))
					assert.False(repo.MustIsInitialized(verbose))
					assert.False(repo.MustHasInitialCommit(verbose))
				}

				for i := 0; i < 2; i++ {
					repo.MustInit(
						&CreateRepositoryOptions{
							Verbose: verbose,
						},
					)
					assert.True(repo.MustExists(verbose))
					assert.True(repo.MustIsInitialized(verbose))
					assert.False(repo.MustHasInitialCommit(verbose))
				}

				for i := 0; i < 2; i++ {
					repo.MustInit(
						&CreateRepositoryOptions{
							Verbose:                   verbose,
							InitializeWithEmptyCommit: true,
						},
					)
					assert.True(repo.MustExists(verbose))
					assert.True(repo.MustIsInitialized(verbose))
					assert.True(repo.MustHasInitialCommit(verbose))
				}

				for i := 0; i < 2; i++ {
					repo.MustDelete(verbose)
					assert.False(repo.MustExists(verbose))
					assert.False(repo.MustIsInitialized(verbose))
					assert.False(repo.MustHasInitialCommit(verbose))
				}
			},
		)
	}
}

func TestGitRepository_CreateAndDeleteRepository(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				for i := 0; i < 2; i++ {
					repo.MustDelete(verbose)
					assert.False(repo.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					repo.MustCreateAndInit(
						&CreateRepositoryOptions{
							InitializeWithEmptyCommit:   true,
							InitializeWithDefaultAuthor: true,
							Verbose:                     verbose,
						},
					)
					assert.True(repo.MustExists(verbose))
					assert.True(repo.MustIsInitialized(verbose))
				}

				for i := 0; i < 2; i++ {
					repo.MustDelete(verbose)
					assert.False(repo.MustExists(verbose))
				}

			},
		)
	}
}
