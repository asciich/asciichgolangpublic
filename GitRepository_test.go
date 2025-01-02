package asciichgolangpublic

import (
	"strconv"
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

func TestGitRepository_IsGitRepository(t *testing.T) {
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

				// An non existing directory is not a git repository:
				assert.False(repo.MustIsGitRepository(verbose))

				Directories().MustCreateLocalDirectoryByPath(repo.MustGetPath(), verbose)

				// The directory exists but is empty which is not a git directory:
				assert.False(repo.MustIsGitRepository(verbose))

				for i := 0; i < 2; i++ {
					repo.MustInit(
						&CreateRepositoryOptions{
							Verbose:        verbose,
							BareRepository: tt.bareRepository,
						},
					)
					assert.True(repo.MustIsGitRepository(verbose))
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

func TestGitRepository_Init_fullInOneStep(t *testing.T) {
	tests := []struct {
		implementationName string
		bare               bool
	}{
		{"localGitRepository", false},
		{"localCommandExecutorRepository", false},
		{"localGitRepository", true},
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
				repo.Delete(verbose)

				for i := 0; i < 2; i++ {
					repo.MustInit(
						&CreateRepositoryOptions{
							Verbose:                     verbose,
							InitializeWithDefaultAuthor: true,
							InitializeWithEmptyCommit:   true,
							BareRepository:              tt.bare,
						},
					)
					assert.True(repo.MustExists(verbose))
					assert.True(repo.MustIsInitialized(verbose))
					assert.True(repo.MustHasInitialCommit(verbose))
					assert.EqualValues(
						tt.bare,
						repo.MustIsBareRepository(verbose),
					)
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

func TestGitRepository_HasNoUncommittedChanges(t *testing.T) {
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

				assert.True(repo.MustHasNoUncommittedChanges(verbose))

				repo.MustCreateFileInDirectory(verbose, "hello.txt")
				assert.False(repo.MustHasNoUncommittedChanges(verbose))
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

func TestGitRepository_CheckHasNoUncommittedChanges(t *testing.T) {
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

				assert.Nil(repo.CheckHasNoUncommittedChanges(verbose))

				repo.MustCreateFileInDirectory(verbose, "hello.txt")
				assert.NotNil(repo.CheckHasNoUncommittedChanges(verbose))
			},
		)
	}
}

func TestGitRepository_GetRootDirectoryPath(t *testing.T) {
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
					repo.MustGetRootDirectory(verbose).MustGetPath(),
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

func TestGitRepository_PullAndPush(t *testing.T) {
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
				assert := assert.New(t)

				const verbose bool = true

				upstreamRepo := getGitRepositoryToTest(tt.implementationUpstream)
				defer upstreamRepo.MustDelete(verbose)
				upstreamRepo.MustDelete(verbose)
				upstreamRepo.MustInit(
					&CreateRepositoryOptions{
						BareRepository:              true,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
						Verbose:                     verbose,
					},
				)

				clonedRepo := getGitRepositoryToTest(tt.implementationCloned)
				defer clonedRepo.MustDelete(verbose)
				clonedRepo.MustDelete(verbose)
				clonedRepo.MustCloneRepository(upstreamRepo, verbose)

				clonedRepo.MustSetGitConfig(
					&GitConfigSetOptions{
						Name:  "Test User",
						Email: "user@example.com",
					},
				)

				clonedRepo2 := getGitRepositoryToTest(tt.implementationCloned2)
				defer clonedRepo2.MustDelete(verbose)
				clonedRepo2.MustDelete(verbose)
				clonedRepo2.MustCloneRepository(upstreamRepo, verbose)

				clonedRepo2.MustSetGitConfig(
					&GitConfigSetOptions{
						Name:  "Test User2",
						Email: "user2@example.com",
					},
				)

				assert.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo.MustGetCurrentCommitHash(verbose),
				)

				assert.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo2.MustGetCurrentCommitHash(verbose),
				)

				fileName := "abc.txt"
				clonedRepo2.MustCreateFileInDirectory(verbose, fileName)
				clonedRepo2.MustAddFileByPath(fileName, verbose)
				clonedRepo2.MustCommit(
					&GitCommitOptions{
						Message: "another commit",
						Verbose: verbose,
					},
				)

				assert.NotEqualValues(
					upstreamRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo2.MustGetCurrentCommitHash(verbose),
				)

				assert.NotEqualValues(
					clonedRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo2.MustGetCurrentCommitHash(verbose),
				)

				clonedRepo2.MustPush(verbose)
				assert.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo2.MustGetCurrentCommitHash(verbose),
				)
				assert.NotEqualValues(
					upstreamRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo.MustGetCurrentCommitHash(verbose),
				)

				clonedRepo.MustPull(verbose)
				assert.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo2.MustGetCurrentCommitHash(verbose),
				)
				assert.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo.MustGetCurrentCommitHash(verbose),
				)
			},
		)
	}
}

func TestGitRepository_AddFilesByPath(t *testing.T) {
	tests := []struct {
		implementationUpstream string
		implementationCloned   string
	}{
		{"localGitRepository", "localGitRepository"},
		{"localCommandExecutorRepository", "localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				upstreamRepo := getGitRepositoryToTest(tt.implementationUpstream)
				defer upstreamRepo.MustDelete(verbose)
				upstreamRepo.MustDelete(verbose)
				upstreamRepo.MustInit(
					&CreateRepositoryOptions{
						BareRepository:              true,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
						Verbose:                     verbose,
					},
				)

				clonedRepo := getGitRepositoryToTest(tt.implementationCloned)
				defer clonedRepo.MustDelete(verbose)
				clonedRepo.MustDelete(verbose)
				clonedRepo.MustCloneRepository(upstreamRepo, verbose)

				clonedRepo.MustSetGitConfig(
					&GitConfigSetOptions{
						Name:  "Test User",
						Email: "user@example.com",
					},
				)

				clonedRepo2 := getGitRepositoryToTest(tt.implementationCloned)
				defer clonedRepo2.MustDelete(verbose)
				clonedRepo2.MustDelete(verbose)
				clonedRepo2.MustCloneRepository(upstreamRepo, verbose)

				clonedRepo2.MustSetGitConfig(
					&GitConfigSetOptions{
						Name:  "Test User2",
						Email: "user2@example.com",
					},
				)

				fileName := "abc.txt"
				fileName2 := "abc.txt"
				clonedRepo.MustCreateFileInDirectory(verbose, fileName)
				clonedRepo.MustCreateFileInDirectory(verbose, fileName2)
				clonedRepo.MustAddFilesByPath(
					[]string{fileName, fileName2},
					verbose,
				)
				clonedRepo.MustCommit(
					&GitCommitOptions{
						Message: "another commit",
						Verbose: verbose,
					},
				)
				clonedRepo.MustPush(verbose)
				
				assert.False(clonedRepo2.MustFileByPathExists(fileName, verbose))
				assert.False(clonedRepo2.MustFileByPathExists(fileName2, verbose))
				
				clonedRepo2.MustPull(verbose)
				assert.True(clonedRepo2.MustFileByPathExists(fileName, verbose))
				assert.True(clonedRepo2.MustFileByPathExists(fileName2, verbose))				
			},
		)
	}
}

func TestGitRepository_FileByPathExists(t *testing.T) {
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

				assert.False(repo.MustFileByPathExists("hello.txt", verbose))

				repo.MustCreateFileInDirectory(verbose, "hello.txt")
				assert.True(repo.MustFileByPathExists("hello.txt", verbose))
			},
		)
	}
}

func TestGitRepository_ListFiles(t *testing.T) {
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

				repo.MustCreateFileInDirectory(verbose, "a.txt")
				repo.MustCreateFileInDirectory(verbose, "b.txt")
				repo.MustCreateFileInDirectory(verbose, "c.txt")
				repo.MustCreateFileInDirectory(verbose, "cb.txt")

				files := repo.MustListFiles(
					&ListFileOptions{
						MatchBasenamePattern: []string{"^b\\.txt$"},
						Verbose:              verbose,
					},
				)
				assert.Len(files, 1)
				assert.EqualValues(
					"b.txt",
					files[0].MustGetBaseName(),
				)

				files = repo.MustListFiles(
					&ListFileOptions{
						MatchBasenamePattern: []string{"^.*b\\.txt$"},
						Verbose:              verbose,
					},
				)
				assert.Len(files, 2)
			},
		)
	}
}

func TestGitRepository_CreateTag(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustCommit(
					&GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)

				tagList := gitRepo.MustListTagNames(verbose)
				assert.Len(tagList, 0)

				expectedTags := []string{}
				for i := 0; i < 5; i++ {

					tagNameToAdd := "ExampleTag" + strconv.Itoa(i)

					gitRepo.MustCreateTag(
						&GitRepositoryCreateTagOptions{
							TagName: tagNameToAdd,
							Verbose: verbose,
						},
					)
					expectedTags = append(expectedTags, tagNameToAdd)

					tagList := gitRepo.MustListTagNames(verbose)

					assert.EqualValues(expectedTags, tagList)
				}
			},
		)
	}
}

func TestGitRepository_ListTags(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustCommit(
					&GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "abc",
						Verbose: verbose,
					},
				)

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "abcd",
						Verbose: verbose,
					},
				)

				tags := gitRepo.MustListTags(verbose)

				assert.EqualValues(
					"abc",
					tags[0].MustGetName(),
				)
				assert.EqualValues(
					"abcd",
					tags[1].MustGetName(),
				)
			},
		)
	}
}

func TestGitRepository_GetLatestTagVersionOrNilIfNotFound(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustCommit(
					&GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)

				assert.Nil(gitRepo.MustGetLatestTagVersionOrNilIfNotFound(verbose))

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				assert.NotNil(gitRepo.MustGetLatestTagVersionOrNilIfNotFound(verbose))

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				assert.NotNil(gitRepo.MustGetLatestTagVersionOrNilIfNotFound(verbose))
			},
		)
	}
}

func TestGitRepository_GetLatestTagVersion(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustCommit(
					&GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				latestVersion := gitRepo.MustGetLatestTagVersion(verbose)

				assert.EqualValues(
					"v1.0.0",
					latestVersion.MustGetAsString(),
				)
			},
		)
	}
}

func TestGitRepository_GetLatestTagVersionAsString(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustCommit(
					&GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				assert.EqualValues(
					"v1.0.0",
					gitRepo.MustGetLatestTagVersionAsString(verbose),
				)
			},
		)
	}
}

func TestGitRepository_GetCurrentCommitsNewestVersion(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustCommit(
					&GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				assert.EqualValues(
					"v0.1.2",
					gitRepo.MustGetCurrentCommitsNewestVersion(verbose).MustGetAsString(),
				)

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				assert.EqualValues(
					"v1.0.0",
					gitRepo.MustGetCurrentCommitsNewestVersion(verbose).MustGetAsString(),
				)
			},
		)
	}
}

func TestGitRepository_GetCurrentCommitsNewestVersionOrNilIfUnset(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustCommit(
					&GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)

				assert.EqualValues(
					nil,
					gitRepo.MustGetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose),
				)

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				assert.EqualValues(
					"v0.1.2",
					gitRepo.MustGetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose).MustGetAsString(),
				)

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				assert.EqualValues(
					"v1.0.0",
					gitRepo.MustGetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose).MustGetAsString(),
				)
			},
		)
	}
}

func TestGitRepository_IsGolangApplication_emptyRepo(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				assert.False(gitRepo.MustIsGolangApplication(verbose))
			},
		)
	}
}

func TestGitRepository_CheckIsGolangApplication_emptyRepo(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				assert.NotNil(gitRepo.CheckIsGolangApplication(verbose))
			},
		)
	}
}

func TestGitRepository_IsGolangApplication_onlyGoMod(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")

				assert.False(gitRepo.MustIsGolangApplication(verbose))
			},
		)
	}
}

func TestGitRepository_CheckIsGolangApplication_onlyGoMod(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")

				assert.NotNil(gitRepo.CheckIsGolangApplication(verbose))
			},
		)
	}
}

func TestGitRepository_CheckIsGolangApplication_NoMainFunction(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")
				gitRepo.MustWriteStringToFile("package main\nfunc abc() bool {\n\treturn true\n}\n", verbose, "main.go")

				assert.NotNil(gitRepo.CheckIsGolangApplication(verbose))
			},
		)
	}
}

func TestGitRepository_IsGolangApplication_NoMainFunction(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")
				gitRepo.MustWriteStringToFile("package main\nfunc abc() bool {\n\treturn true\n}\n", verbose, "main.go")

				assert.False(gitRepo.MustIsGolangApplication(verbose))
			},
		)
	}
}

func TestGitRepository_IsGolangApplication(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")
				gitRepo.MustWriteStringToFile("package main\nfunc main() {\n\treturn\n}\n", verbose, "main.go")

				assert.True(gitRepo.MustIsGolangApplication(verbose))
			},
		)
	}
}

func TestGitRepository_CheckIsGolangApplication(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")
				gitRepo.MustWriteStringToFile("package main\nfunc main() {\n\treturn\n}\n", verbose, "main.go")

				assert.Nil(gitRepo.CheckIsGolangApplication(verbose))
			},
		)
	}
}

func TestGitRepository_GetFileByPath(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("hello world\n", verbose, "test.txt")

				testTxtFile := gitRepo.MustGetFileByPath("test.txt")

				assert.EqualValues(
					"hello world\n",
					testTxtFile.MustReadAsString(),
				)
			},
		)
	}
}

func TestGitRepository_IsGolangPackage_emptyRepo(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				assert.False(
					gitRepo.MustIsGolangPackage(verbose),
				)
			},
		)
	}
}

func TestGitRepository_IsGolangPackage_onlyGoMod(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")

				assert.True(
					gitRepo.MustIsGolangPackage(verbose),
				)
			},
		)
	}
}

func TestGitRepository_IsGolangPackage_mainFunctionIsNotAPackage(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")
				gitRepo.MustWriteStringToFile("package main\nfunc main() {\n\treturn\n}\n", verbose, "main.go")

				assert.False(
					gitRepo.MustIsGolangPackage(verbose),
				)
			},
		)
	}
}

func TestGitRepository_CheckIsGolangPackage_emptyRepo(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				assert.NotNil(
					gitRepo.CheckIsGolangPackage(verbose),
				)
			},
		)
	}
}

func TestGitRepository_CheckIsGolangPackage_onlyGoMod(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")

				assert.Nil(
					gitRepo.CheckIsGolangPackage(verbose),
				)
			},
		)
	}
}

func TestGitRepository_CheckIsGolangPackage_mainFunctionIsNotAPackage(t *testing.T) {
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

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")
				gitRepo.MustWriteStringToFile("package main\nfunc main() {\n\treturn\n}\n", verbose, "main.go")

				assert.NotNil(
					gitRepo.CheckIsGolangPackage(verbose),
				)
			},
		)
	}
}
