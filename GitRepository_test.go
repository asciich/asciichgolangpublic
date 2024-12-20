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
							Verbose:                     verbose,
							InitializeWithEmptyCommit:   true,
							InitializeWithDefaultAuthor: true,
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

func TestGitRepository_HasUncommittedChanges(t *testing.T) {
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

				assert.False(repo.MustHasUncommittedChanges(verbose))

				repo.MustCreateFileInDirectory(verbose, "hello.txt")
				assert.True(repo.MustHasUncommittedChanges(verbose))
			},
		)
	}
}

func TestGitRepository_GetRootDirectory(t *testing.T) {
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

				repo.MustInit(
					&CreateRepositoryOptions{
						Verbose:        verbose,
						BareRepository: tt.bareRepository,
					},
				)

				assert.EqualValues(
					repo.MustGetPath(),
					repo.MustGetRootDirectoryPath(verbose),
				)
			},
		)
	}
}

func TestGitRepository_GetRootDirectory_from_subdirectory(t *testing.T) {
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

				repo.MustInit(
					&CreateRepositoryOptions{
						Verbose:        verbose,
						BareRepository: tt.bareRepository,
					},
				)

				expectedRootDirectory := repo.MustGetPath()

				subDir := repo.MustCreateSubDirectory("sub_directory", verbose)

				repoUsingSubDir1 := MustGetCommandExecutorGitRepositoryFromDirectory(subDir)

				assert.EqualValues(
					expectedRootDirectory,
					repoUsingSubDir1.MustGetRootDirectoryPath(verbose),
				)

				repoUsingSubDir2 := MustGetLocalGitReposioryFromDirectory(subDir)

				assert.EqualValues(
					expectedRootDirectory,
					repoUsingSubDir2.MustGetRootDirectoryPath(verbose),
				)
			},
		)
	}
}

func TestLocalGitRepositoryPullAndPush(t *testing.T) {
	tests := []struct {
		implementationUpstream string
		implementationCloned   string
		implementationCloned2  string
	}{
		{"localGitRepository", "localGitRepository", "localGitRepository"},
		{"localGitRepository", "localGitRepository", "localCommandExecutorRepository"},
		{"localCommandExecutorRepository", "localCommandExecutorRepository", "localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				// assert := assert.New(t)

				const verbose bool = true

				upstreamRepo := getGitRepositoryToTest(tt.implementationUpstream)
				defer upstreamRepo.MustDelete(verbose)

				clonedRepo := getGitRepositoryToTest(tt.implementationCloned)
				defer clonedRepo.Delete(verbose)

				clonedRepo.MustDelete(verbose)
				clonedRepo.MustCloneRepository(upstreamRepo, verbose)

				/*
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
				*/
			},
		)
	}
}
