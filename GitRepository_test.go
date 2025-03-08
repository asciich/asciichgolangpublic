package asciichgolangpublic

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/contextutils"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/tempfiles"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func getGitRepositoryToTest(implementationName string) (repo GitRepository) {
	const verbose = true

	if implementationName == "localGitRepository" {
		repo = MustGetLocalGitReposioryFromDirectory(
			tempfiles.MustCreateEmptyTemporaryDirectory(verbose),
		)
	} else if implementationName == "localCommandExecutorRepository" {
		repo = MustGetLocalCommandExecutorGitRepositoryByDirectory(
			tempfiles.MustCreateEmptyTemporaryDirectory(verbose),
		)
	} else {
		logging.LogFatalWithTracef("unknown implementationName='%s'", implementationName)
	}

	repo.MustInit(
		&parameteroptions.CreateRepositoryOptions{
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustDelete(verbose)

				require.False(repo.MustExists(verbose))
				require.False(repo.MustIsInitialized(verbose))
				require.False(repo.MustHasInitialCommit(verbose))

				for i := 0; i < 2; i++ {
					repo.MustInit(
						&parameteroptions.CreateRepositoryOptions{
							Verbose:        verbose,
							BareRepository: tt.bareRepository,
						},
					)
					require.True(repo.MustExists(verbose))
					require.True(repo.MustIsInitialized(verbose))
					require.False(repo.MustHasInitialCommit(verbose))
					require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustDelete(verbose)

				// An non existing directory is not a git repository:
				require.False(repo.MustIsGitRepository(verbose))

				files.Directories().MustCreateLocalDirectoryByPath(repo.MustGetPath(), verbose)

				// The directory exists but is empty which is not a git directory:
				require.False(repo.MustIsGitRepository(verbose))

				for i := 0; i < 2; i++ {
					repo.MustInit(
						&parameteroptions.CreateRepositoryOptions{
							Verbose:        verbose,
							BareRepository: tt.bareRepository,
						},
					)
					require.True(repo.MustIsGitRepository(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				for i := 0; i < 2; i++ {
					repo.MustDelete(verbose)
					require.False(repo.MustExists(verbose))
					require.False(repo.MustIsInitialized(verbose))
					require.False(repo.MustHasInitialCommit(verbose))
				}

				for i := 0; i < 2; i++ {
					repo.MustCreate(verbose)
					require.True(repo.MustExists(verbose))
					require.False(repo.MustIsInitialized(verbose))
					require.False(repo.MustHasInitialCommit(verbose))
				}

				for i := 0; i < 2; i++ {
					repo.MustInit(
						&parameteroptions.CreateRepositoryOptions{
							Verbose: verbose,
						},
					)
					require.True(repo.MustExists(verbose))
					require.True(repo.MustIsInitialized(verbose))
					require.False(repo.MustHasInitialCommit(verbose))
				}

				for i := 0; i < 2; i++ {
					repo.MustInit(
						&parameteroptions.CreateRepositoryOptions{
							Verbose:                     verbose,
							InitializeWithEmptyCommit:   true,
							InitializeWithDefaultAuthor: true,
						},
					)
					require.True(repo.MustExists(verbose))
					require.True(repo.MustIsInitialized(verbose))
					require.True(repo.MustHasInitialCommit(verbose))
				}

				for i := 0; i < 2; i++ {
					repo.MustDelete(verbose)
					require.False(repo.MustExists(verbose))
					require.False(repo.MustIsInitialized(verbose))
					require.False(repo.MustHasInitialCommit(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)
				repo.Delete(verbose)

				for i := 0; i < 2; i++ {
					repo.MustInit(
						&parameteroptions.CreateRepositoryOptions{
							Verbose:                     verbose,
							InitializeWithDefaultAuthor: true,
							InitializeWithEmptyCommit:   true,
							BareRepository:              tt.bare,
						},
					)
					require.True(repo.MustExists(verbose))
					require.True(repo.MustIsInitialized(verbose))
					require.True(repo.MustHasInitialCommit(verbose))
					require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				for i := 0; i < 2; i++ {
					repo.MustDelete(verbose)
					require.False(repo.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					repo.MustCreateAndInit(
						&parameteroptions.CreateRepositoryOptions{
							InitializeWithEmptyCommit:   true,
							InitializeWithDefaultAuthor: true,
							Verbose:                     verbose,
						},
					)
					require.True(repo.MustExists(verbose))
					require.True(repo.MustIsInitialized(verbose))
				}

				for i := 0; i < 2; i++ {
					repo.MustDelete(verbose)
					require.False(repo.MustExists(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				require.True(repo.MustHasNoUncommittedChanges(verbose))

				repo.MustCreateFileInDirectory(verbose, "hello.txt")
				require.False(repo.MustHasNoUncommittedChanges(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				require.False(repo.MustHasUncommittedChanges(verbose))

				repo.MustCreateFileInDirectory(verbose, "hello.txt")
				require.True(repo.MustHasUncommittedChanges(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				require.Nil(repo.CheckHasNoUncommittedChanges(verbose))

				repo.MustCreateFileInDirectory(verbose, "hello.txt")
				require.NotNil(repo.CheckHasNoUncommittedChanges(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustDelete(verbose)

				repo.MustInit(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:        verbose,
						BareRepository: tt.bareRepository,
					},
				)

				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustDelete(verbose)

				repo.MustInit(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:        verbose,
						BareRepository: tt.bareRepository,
					},
				)

				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				repo.MustDelete(verbose)

				repo.MustInit(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:        verbose,
						BareRepository: tt.bareRepository,
					},
				)

				expectedRootDirectory := repo.MustGetPath()

				subDir := repo.MustCreateSubDirectory("sub_directory", verbose)

				repoUsingSubDir1 := MustGetCommandExecutorGitRepositoryFromDirectory(subDir)

				require.EqualValues(
					expectedRootDirectory,
					repoUsingSubDir1.MustGetRootDirectoryPath(verbose),
				)

				repoUsingSubDir2 := MustGetLocalGitReposioryFromDirectory(subDir)

				require.EqualValues(
					expectedRootDirectory,
					repoUsingSubDir2.MustGetRootDirectoryPath(verbose),
				)
			},
		)
	}
}

func TestGitRepository_CloneRepository_idempotence(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				upstreamRepo := getGitRepositoryToTest(tt.implementationName)
				defer upstreamRepo.MustDelete(verbose)
				upstreamRepo.MustDelete(verbose)
				upstreamRepo.MustInit(
					&parameteroptions.CreateRepositoryOptions{
						BareRepository:              true,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
						Verbose:                     verbose,
					},
				)

				clonedRepo := getGitRepositoryToTest(tt.implementationName)
				defer clonedRepo.MustDelete(verbose)
				clonedRepo.MustDelete(verbose)

				for i := 0; i < 2; i++ {
					clonedRepo.MustCloneRepository(upstreamRepo, verbose)
				}

				require.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo.MustGetCurrentCommitHash(verbose),
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				upstreamRepo := getGitRepositoryToTest(tt.implementationUpstream)
				defer upstreamRepo.MustDelete(verbose)
				upstreamRepo.MustDelete(verbose)
				upstreamRepo.MustInit(
					&parameteroptions.CreateRepositoryOptions{
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

				require.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo.MustGetCurrentCommitHash(verbose),
				)

				require.EqualValues(
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

				require.NotEqualValues(
					upstreamRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo2.MustGetCurrentCommitHash(verbose),
				)

				require.NotEqualValues(
					clonedRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo2.MustGetCurrentCommitHash(verbose),
				)

				clonedRepo2.MustPush(verbose)
				require.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo2.MustGetCurrentCommitHash(verbose),
				)
				require.NotEqualValues(
					upstreamRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo.MustGetCurrentCommitHash(verbose),
				)

				clonedRepo.MustPull(verbose)
				require.EqualValues(
					upstreamRepo.MustGetCurrentCommitHash(verbose),
					clonedRepo2.MustGetCurrentCommitHash(verbose),
				)
				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				upstreamRepo := getGitRepositoryToTest(tt.implementationUpstream)
				defer upstreamRepo.MustDelete(verbose)
				upstreamRepo.MustDelete(verbose)
				upstreamRepo.MustInit(
					&parameteroptions.CreateRepositoryOptions{
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

				require.False(clonedRepo2.MustFileByPathExists(fileName, verbose))
				require.False(clonedRepo2.MustFileByPathExists(fileName2, verbose))

				clonedRepo2.MustPull(verbose)
				require.True(clonedRepo2.MustFileByPathExists(fileName, verbose))
				require.True(clonedRepo2.MustFileByPathExists(fileName2, verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				require.False(repo.MustHasUncommittedChanges(verbose))

				require.False(repo.MustFileByPathExists("hello.txt", verbose))

				repo.MustCreateFileInDirectory(verbose, "hello.txt")
				require.True(repo.MustFileByPathExists("hello.txt", verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				require.False(repo.MustHasUncommittedChanges(verbose))

				repo.MustCreateFileInDirectory(verbose, "a.txt")
				repo.MustCreateFileInDirectory(verbose, "b.txt")
				repo.MustCreateFileInDirectory(verbose, "c.txt")
				repo.MustCreateFileInDirectory(verbose, "cb.txt")

				files := repo.MustListFiles(
					&parameteroptions.ListFileOptions{
						MatchBasenamePattern: []string{"^b\\.txt$"},
						Verbose:              verbose,
					},
				)
				require.Len(files, 1)
				require.EqualValues(
					"b.txt",
					files[0].MustGetBaseName(),
				)

				files = repo.MustListFiles(
					&parameteroptions.ListFileOptions{
						MatchBasenamePattern: []string{"^.*b\\.txt$"},
						Verbose:              verbose,
					},
				)
				require.Len(files, 2)
			},
		)
	}
}

func TestGitRepository_ListFilePaths(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				require.False(repo.MustHasUncommittedChanges(verbose))

				repo.MustCreateFileInDirectory(verbose, "a.txt")
				repo.MustCreateFileInDirectory(verbose, "b.txt")
				repo.MustCreateFileInDirectory(verbose, "c.txt")
				repo.MustCreateFileInDirectory(verbose, "cb.txt")

				files := repo.MustListFilePaths(
					&parameteroptions.ListFileOptions{
						MatchBasenamePattern: []string{"^b\\.txt$"},
						Verbose:              verbose,
						ReturnRelativePaths:  true,
					},
				)
				require.Len(files, 1)
				require.EqualValues(
					"b.txt",
					files[0],
				)

				files = repo.MustListFilePaths(
					&parameteroptions.ListFileOptions{
						MatchBasenamePattern: []string{"^.*b\\.txt$"},
						Verbose:              verbose,
						ReturnRelativePaths:  true,
					},
				)
				require.Len(files, 2)
				require.EqualValues(
					[]string{"b.txt", "cb.txt"},
					files,
				)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

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
				require.Len(tagList, 0)

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

					require.EqualValues(expectedTags, tagList)
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

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

				require.EqualValues(
					"abc",
					tags[0].MustGetName(),
				)
				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

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

				require.Nil(gitRepo.MustGetLatestTagVersionOrNilIfNotFound(verbose))

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				require.NotNil(gitRepo.MustGetLatestTagVersionOrNilIfNotFound(verbose))

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				require.NotNil(gitRepo.MustGetLatestTagVersionOrNilIfNotFound(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

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

				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

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

				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

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

				require.EqualValues(
					"v0.1.2",
					gitRepo.MustGetCurrentCommitsNewestVersion(verbose).MustGetAsString(),
				)

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

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

				require.EqualValues(
					nil,
					gitRepo.MustGetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose),
				)

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)

				require.EqualValues(
					"v0.1.2",
					gitRepo.MustGetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose).MustGetAsString(),
				)

				gitRepo.MustCreateTag(
					&GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)

				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				require.False(gitRepo.MustIsGolangApplication(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				require.NotNil(gitRepo.CheckIsGolangApplication(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")

				require.False(gitRepo.MustIsGolangApplication(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")

				require.NotNil(gitRepo.CheckIsGolangApplication(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")
				gitRepo.MustWriteStringToFile("package main\nfunc abc() bool {\n\treturn true\n}\n", verbose, "main.go")

				require.NotNil(gitRepo.CheckIsGolangApplication(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")
				gitRepo.MustWriteStringToFile("package main\nfunc abc() bool {\n\treturn true\n}\n", verbose, "main.go")

				require.False(gitRepo.MustIsGolangApplication(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")
				gitRepo.MustWriteStringToFile("package main\nfunc main() {\n\treturn\n}\n", verbose, "main.go")

				require.True(gitRepo.MustIsGolangApplication(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")
				gitRepo.MustWriteStringToFile("package main\nfunc main() {\n\treturn\n}\n", verbose, "main.go")

				require.Nil(gitRepo.CheckIsGolangApplication(verbose))
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("hello world\n", verbose, "test.txt")

				testTxtFile := gitRepo.MustGetFileByPath("test.txt")

				require.EqualValues(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				require.False(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")

				require.True(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")
				gitRepo.MustWriteStringToFile("package main\nfunc main() {\n\treturn\n}\n", verbose, "main.go")

				require.False(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				require.NotNil(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")

				require.Nil(
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("module example\n", verbose, "go.mod")
				gitRepo.MustWriteStringToFile("package main\nfunc main() {\n\treturn\n}\n", verbose, "main.go")

				require.NotNil(
					gitRepo.CheckIsGolangPackage(verbose),
				)
			},
		)
	}
}

func TestGitRepository_GetGitRepositoryByDirectory(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				repoRootDirectory := gitRepo.MustGetRootDirectory(verbose)

				gitRepo2 := MustGetGitRepositoryByDirectory(repoRootDirectory)

				require.EqualValues(
					gitRepo.MustGetRootDirectoryPath(verbose),
					gitRepo2.MustGetRootDirectoryPath(verbose),
				)

				require.Nil(
					gitRepo.CheckHasNoUncommittedChanges(verbose),
				)

				require.Nil(
					gitRepo2.CheckHasNoUncommittedChanges(verbose),
				)
			},
		)
	}
}

func TestGitRepository_CreateAndDeleteBranch(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				defaultBranchName := gitRepo.MustGetCurrentBranchName(verbose)

				require.False(gitRepo.MustBranchByNameExists("testbranch", verbose))

				for i := 0; i < 2; i++ {
					gitRepo.MustCreateBranch(
						&parameteroptions.CreateBranchOptions{
							Name:    "testbranch",
							Verbose: verbose,
						},
					)
					require.True(gitRepo.MustBranchByNameExists("testbranch", verbose))
				}

				gitRepo.MustCheckoutBranchByName(defaultBranchName, verbose)

				for i := 0; i < 2; i++ {
					gitRepo.MustDeleteBranchByName("testbranch", verbose)
					require.False(gitRepo.MustBranchByNameExists("testbranch", verbose))
				}
			},
		)
	}
}

func TestGitRepository_CheckoutBranch(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				for _, branchName := range []string{"testbranch1", "testbranch2"} {
					gitRepo.MustCreateBranch(
						&parameteroptions.CreateBranchOptions{
							Name:    branchName,
							Verbose: verbose,
						},
					)
					require.True(gitRepo.MustBranchByNameExists(branchName, verbose))

					gitRepo.MustCheckoutBranchByName(branchName, verbose)
					require.EqualValues(
						branchName,
						gitRepo.MustGetCurrentBranchName(verbose),
					)
				}
			},
		)
	}
}

func TestGitRepository_GetCurrentCommitMessage(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustCommit(
					&GitCommitOptions{
						AllowEmpty: true,
						Message:    "commit message",
						Verbose:    verbose,
					},
				)

				require.EqualValues(
					"commit message",
					gitRepo.MustGetCurrentCommitMessage(verbose),
				)
			},
		)
	}
}

func TestGitRepository_CommitIfUncommittedChanges(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustCommit(
					&GitCommitOptions{
						AllowEmpty: true,
						Message:    "commit before testing",
						Verbose:    verbose,
					},
				)

				require.EqualValues(
					"commit before testing",
					gitRepo.MustGetCurrentCommitMessage(verbose),
				)

				gitRepo.MustCommitIfUncommittedChanges(
					&GitCommitOptions{
						Message: "This should not trigger a commit",
						Verbose: verbose,
					},
				)

				require.EqualValues(
					"commit before testing",
					gitRepo.MustGetCurrentCommitMessage(verbose),
				)

				gitRepo.MustWriteStringToFile("hello", verbose, "world.txt")
				gitRepo.MustAddFileByPath("world.txt", verbose)

				gitRepo.MustCommitIfUncommittedChanges(
					&GitCommitOptions{
						Message: "This should trigger a commit",
						Verbose: verbose,
					},
				)

				require.EqualValues(
					"This should trigger a commit",
					gitRepo.MustGetCurrentCommitMessage(verbose),
				)

				gitRepo.MustWriteStringToFile("world", verbose, "world.txt")
				// world.txt is already known in the git repo.
				// No need to explicitly add world.txt again.

				gitRepo.MustCommitIfUncommittedChanges(
					&GitCommitOptions{
						Message: "This should trigger again a commit",
						Verbose: verbose,
					},
				)

				require.EqualValues(
					"This should trigger again a commit",
					gitRepo.MustGetCurrentCommitMessage(verbose),
				)
			},
		)
	}
}

func TestGitRepository_AddAndRemoveRemote(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				require.False(
					gitRepo.MustRemoteByNameExists("example", verbose),
				)

				for i := 0; i < 2; i++ {
					gitRepo.MustAddRemote(
						&GitRemoteAddOptions{
							RemoteName: "example",
							RemoteUrl:  "https://remote.url.example.com",
							Verbose:    verbose,
						},
					)

					require.True(
						gitRepo.MustRemoteByNameExists("example", verbose),
					)
				}

				for i := 0; i < 2; i++ {
					gitRepo.MustRemoveRemoteByName("example", verbose)
					require.False(
						gitRepo.MustRemoteByNameExists("example", verbose),
					)
				}
			},
		)
	}
}

func TestGitRepository_IsPreCommitRepository(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				require.False(
					gitRepo.MustIsPreCommitRepository(verbose),
				)

				gitRepo.MustCreateSubDirectory("pre_commit_hooks", verbose)

				require.True(
					gitRepo.MustIsPreCommitRepository(verbose),
				)
			},
		)
	}
}

func TestGitRepository_CheckIsPreCommitRepository(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				require.NotNil(
					gitRepo.CheckIsPreCommitRepository(verbose),
				)

				gitRepo.MustCreateSubDirectory("pre_commit_hooks", verbose)

				require.Nil(
					gitRepo.CheckIsPreCommitRepository(verbose),
				)
			},
		)
	}
}

func Test_CheckExists(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"localGitRepository"},
		{"localCommandExecutorRepository"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := contextutils.GetVerbosityContextByBool(true)

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(true)

				require.NoError(t, gitRepo.CheckExists(ctx))

				gitRepo.MustDelete(true)

				require.Error(t, gitRepo.CheckExists(ctx))
			},
		)
	}
}
