package asciichgolangpublic

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func getGitRepositoryToTest(implementationName string) (repo GitRepository) {
	const verbose = true

	if implementationName == "localGitRepository" {
		repo = MustGetLocalGitReposioryFromDirectory(
			mustutils.Must(tempfilesoo.CreateEmptyTemporaryDirectory(verbose)),
		)
	} else if implementationName == "localCommandExecutorRepository" {
		repo = MustGetLocalCommandExecutorGitRepositoryByDirectory(
			mustutils.Must(tempfilesoo.CreateEmptyTemporaryDirectory(verbose)),
		)
	} else {
		logging.LogFatalWithTracef("unknown implementationName='%s'", implementationName)
	}

	mustutils.Must0(repo.Init(
		&parameteroptions.CreateRepositoryOptions{
			Verbose:                     verbose,
			InitializeWithEmptyCommit:   true,
			InitializeWithDefaultAuthor: true,
		},
	))

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
				ctx := getCtx()
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				err := repo.Delete(verbose)
				require.NoError(t, err)

				require.False(t, mustutils.Must(repo.Exists(verbose)))
				require.False(t, mustutils.Must(repo.IsInitialized(verbose)))
				require.False(t, mustutils.Must(repo.HasInitialCommit(verbose)))

				for i := 0; i < 2; i++ {
					err = repo.Init(
						&parameteroptions.CreateRepositoryOptions{
							Verbose:        verbose,
							BareRepository: tt.bareRepository,
						},
					)
					require.True(t, mustutils.Must(repo.Exists(verbose)))
					require.True(t, mustutils.Must(repo.IsInitialized(verbose)))
					require.False(t, mustutils.Must(repo.HasInitialCommit(verbose)))
					isBare, err := repo.IsBareRepository(ctx)
					require.NoError(t, err)
					require.EqualValues(t, tt.bareRepository, isBare)
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
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				err := repo.Delete(verbose)
				require.NoError(t, err)

				// An non existing directory is not a git repository:
				isRepo, err := repo.IsGitRepository(verbose)
				require.NoError(t, err)
				require.False(t, isRepo)

				path, err := repo.GetPath()
				require.NoError(t, err)
				_, err = files.Directories().CreateLocalDirectoryByPath(path, verbose)
				require.NoError(t, err)

				// The directory exists but is empty which is not a git directory:
				isRepo, err = repo.IsGitRepository(verbose)
				require.NoError(t, err)
				require.False(t, isRepo)

				for i := 0; i < 2; i++ {
					mustutils.Must0(repo.Init(
						&parameteroptions.CreateRepositoryOptions{
							Verbose:        verbose,
							BareRepository: tt.bareRepository,
						},
					))
					isRepo, err = repo.IsGitRepository(verbose)
					require.NoError(t, err)
					require.True(t, isRepo)
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
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				for i := 0; i < 2; i++ {
					err := repo.Delete(verbose)
					require.NoError(t, err)
					require.False(t, mustutils.Must(repo.Exists(verbose)))
					require.False(t, mustutils.Must(repo.IsInitialized(verbose)))
					require.False(t, mustutils.Must(repo.HasInitialCommit(verbose)))
				}

				for i := 0; i < 2; i++ {
					err := repo.Create(verbose)
					require.NoError(t, err)
					require.True(t, mustutils.Must(repo.Exists(verbose)))
					require.False(t, mustutils.Must(repo.IsInitialized(verbose)))
					require.False(t, mustutils.Must(repo.HasInitialCommit(verbose)))
				}

				for i := 0; i < 2; i++ {
					err := repo.Init(
						&parameteroptions.CreateRepositoryOptions{
							Verbose: verbose,
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(repo.Exists(verbose)))
					require.True(t, mustutils.Must(repo.IsInitialized(verbose)))
					require.False(t, mustutils.Must(repo.HasInitialCommit(verbose)))
				}

				for i := 0; i < 2; i++ {
					err := repo.Init(
						&parameteroptions.CreateRepositoryOptions{
							Verbose:                     verbose,
							InitializeWithEmptyCommit:   true,
							InitializeWithDefaultAuthor: true,
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(repo.Exists(verbose)))
					require.True(t, mustutils.Must(repo.IsInitialized(verbose)))
					require.True(t, mustutils.Must(repo.HasInitialCommit(verbose)))
				}

				for i := 0; i < 2; i++ {
					err := repo.Delete(verbose)
					require.NoError(t, err)
					require.False(t, mustutils.Must(repo.Exists(verbose)))
					require.False(t, mustutils.Must(repo.IsInitialized(verbose)))
					require.False(t, mustutils.Must(repo.HasInitialCommit(verbose)))
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
				ctx := getCtx()
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)
				repo.Delete(verbose)

				for i := 0; i < 2; i++ {
					err := repo.Init(
						&parameteroptions.CreateRepositoryOptions{
							Verbose:                     verbose,
							InitializeWithDefaultAuthor: true,
							InitializeWithEmptyCommit:   true,
							BareRepository:              tt.bare,
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(repo.Exists(verbose)))
					require.True(t, mustutils.Must(repo.IsInitialized(verbose)))
					require.True(t, mustutils.Must(repo.HasInitialCommit(verbose)))
					isBare, err := repo.IsBareRepository(ctx)
					require.NoError(t, err)
					require.EqualValues(t, tt.bare, isBare)
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
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				for i := 0; i < 2; i++ {
					err := repo.Delete(verbose)
					require.NoError(t, err)
					require.False(t, mustutils.Must(repo.Exists(verbose)))
				}

				for i := 0; i < 2; i++ {
					err := repo.CreateAndInit(
						&parameteroptions.CreateRepositoryOptions{
							InitializeWithEmptyCommit:   true,
							InitializeWithDefaultAuthor: true,
							Verbose:                     verbose,
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(repo.Exists(verbose)))
					require.True(t, mustutils.Must(repo.IsInitialized(verbose)))
				}

				for i := 0; i < 2; i++ {
					err := repo.Delete(verbose)
					require.NoError(t, err)
					require.False(t, mustutils.Must(repo.Exists(verbose)))
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
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				require.True(t, mustutils.Must(repo.HasNoUncommittedChanges(verbose)))

				_, err := repo.CreateFileInDirectory(verbose, "hello.txt")
				require.NoError(t, err)
				require.False(t, mustutils.Must(repo.HasNoUncommittedChanges(verbose)))
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
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				require.False(t, mustutils.Must(repo.HasUncommittedChanges(verbose)))

				_, err := repo.CreateFileInDirectory(verbose, "hello.txt")
				require.NoError(t, err)
				require.True(t, mustutils.Must(repo.HasUncommittedChanges(verbose)))
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
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				require.Nil(t, repo.CheckHasNoUncommittedChanges(verbose))

				_, err := repo.CreateFileInDirectory(verbose, "hello.txt")
				require.NoError(t, err)
				require.NotNil(t, repo.CheckHasNoUncommittedChanges(verbose))
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
				ctx := getCtx()
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				err := repo.Delete(verbose)
				require.NoError(t, err)

				err = repo.Init(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:        verbose,
						BareRepository: tt.bareRepository,
					},
				)
				require.NoError(t, err)

				rootDirPath, err := repo.GetRootDirectoryPath(ctx)
				require.NoError(t, err)
				require.EqualValues(t, mustutils.Must(repo.GetPath()), rootDirPath)
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
				const verbose bool = true
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				err := repo.Delete(verbose)
				require.NoError(t, err)

				err = repo.Init(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:        verbose,
						BareRepository: tt.bareRepository,
					},
				)
				require.NoError(t, err)

				rootDir, err := repo.GetRootDirectory(ctx)
				require.NoError(t, err)

				rootDirPath, err := rootDir.GetPath()
				require.NoError(t, err)

				require.EqualValues(t, mustutils.Must(repo.GetPath()), rootDirPath)
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
				const verbose bool = true
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				err := repo.Delete(verbose)
				require.NoError(t, err)

				err = repo.Init(
					&parameteroptions.CreateRepositoryOptions{
						Verbose:        verbose,
						BareRepository: tt.bareRepository,
					},
				)
				require.NoError(t, err)

				expectedRootDirectory, err := repo.GetPath()
				require.NoError(t, err)

				subDir, err := repo.CreateSubDirectory("sub_directory", verbose)
				require.NoError(t, err)

				repoUsingSubDir1, err := GetCommandExecutorGitRepositoryFromDirectory(subDir)
				require.NoError(t, err)

				require.EqualValues(t, expectedRootDirectory, mustutils.Must(repoUsingSubDir1.GetRootDirectoryPath(ctx)))

				repoUsingSubDir2 := MustGetLocalGitReposioryFromDirectory(subDir)

				require.EqualValues(t, expectedRootDirectory, mustutils.Must(repoUsingSubDir2.GetRootDirectoryPath(ctx)))
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
				const verbose bool = true

				upstreamRepo := getGitRepositoryToTest(tt.implementationName)
				defer upstreamRepo.Delete(verbose)
				err := upstreamRepo.Delete(verbose)
				require.NoError(t, err)

				err = upstreamRepo.Init(
					&parameteroptions.CreateRepositoryOptions{
						BareRepository:              true,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
						Verbose:                     verbose,
					},
				)
				require.NoError(t, err)

				clonedRepo := getGitRepositoryToTest(tt.implementationName)
				defer clonedRepo.Delete(verbose)
				err = clonedRepo.Delete(verbose)
				require.NoError(t, err)

				for i := 0; i < 2; i++ {
					err = clonedRepo.CloneRepository(upstreamRepo, verbose)
					require.NoError(t, err)
				}

				require.EqualValues(t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(verbose)),
					mustutils.Must(clonedRepo.GetCurrentCommitHash(verbose)),
				)
			},
		)
	}
}

func TestGitRepository_PullAndPush(t *testing.T) {
	ctx := getCtx()

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
				const verbose bool = true

				upstreamRepo := getGitRepositoryToTest(tt.implementationUpstream)
				defer upstreamRepo.Delete(verbose)
				err := upstreamRepo.Delete(verbose)
				require.NoError(t, err)
				err = upstreamRepo.Init(
					&parameteroptions.CreateRepositoryOptions{
						BareRepository:              true,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
						Verbose:                     verbose,
					},
				)
				require.NoError(t, err)

				clonedRepo := getGitRepositoryToTest(tt.implementationCloned)
				defer clonedRepo.Delete(verbose)
				err = clonedRepo.Delete(verbose)
				require.NoError(t, err)

				err = clonedRepo.CloneRepository(upstreamRepo, verbose)
				require.NoError(t, err)

				err = clonedRepo.SetGitConfig(
					&gitparameteroptions.GitConfigSetOptions{
						Name:  "Test User",
						Email: "user@example.com",
					},
				)
				require.NoError(t, err)

				clonedRepo2 := getGitRepositoryToTest(tt.implementationCloned2)
				defer clonedRepo2.Delete(verbose)
				err = clonedRepo2.Delete(verbose)
				require.NoError(t, err)
				err = clonedRepo2.CloneRepository(upstreamRepo, verbose)
				require.NoError(t, err)

				err = clonedRepo2.SetGitConfig(
					&gitparameteroptions.GitConfigSetOptions{
						Name:  "Test User2",
						Email: "user2@example.com",
					},
				)
				require.NoError(t, err)

				require.EqualValues(
					t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(verbose)),
					mustutils.Must(clonedRepo.GetCurrentCommitHash(verbose)),
				)

				require.EqualValues(
					t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(verbose)),
					mustutils.Must(clonedRepo2.GetCurrentCommitHash(verbose)),
				)

				fileName := "abc.txt"
				_, err = clonedRepo2.CreateFileInDirectory(verbose, fileName)
				require.NoError(t, err)

				err = clonedRepo2.AddFileByPath(fileName, verbose)
				require.NoError(t, err)

				_, err = clonedRepo2.Commit(
					&gitparameteroptions.GitCommitOptions{
						Message: "another commit",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				require.NotEqualValues(t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(verbose)),
					mustutils.Must(clonedRepo2.GetCurrentCommitHash(verbose)),
				)

				require.NotEqualValues(t,
					mustutils.Must(clonedRepo.GetCurrentCommitHash(verbose)),
					mustutils.Must(clonedRepo2.GetCurrentCommitHash(verbose)),
				)

				mustutils.Must0(clonedRepo2.Push(ctx))
				require.EqualValues(t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(verbose)),
					mustutils.Must(clonedRepo2.GetCurrentCommitHash(verbose)),
				)
				require.NotEqualValues(t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(verbose)),
					mustutils.Must(clonedRepo.GetCurrentCommitHash(verbose)),
				)

				err = clonedRepo.Pull(ctx)
				require.NoError(t, err)
				require.EqualValues(t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(verbose)),
					mustutils.Must(clonedRepo2.GetCurrentCommitHash(verbose)),
				)
				require.EqualValues(t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(verbose)),
					mustutils.Must(clonedRepo.GetCurrentCommitHash(verbose)),
				)
			},
		)
	}
}

func TestGitRepository_AddFilesByPath(t *testing.T) {
	ctx := getCtx()

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
				const verbose bool = true

				upstreamRepo := getGitRepositoryToTest(tt.implementationUpstream)
				defer upstreamRepo.Delete(verbose)
				err := upstreamRepo.Delete(verbose)
				require.NoError(t, err)

				err = upstreamRepo.Init(
					&parameteroptions.CreateRepositoryOptions{
						BareRepository:              true,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
						Verbose:                     verbose,
					},
				)
				require.NoError(t, err)

				clonedRepo := getGitRepositoryToTest(tt.implementationCloned)
				defer clonedRepo.Delete(verbose)
				err = clonedRepo.Delete(verbose)
				require.NoError(t, err)

				err = clonedRepo.CloneRepository(upstreamRepo, verbose)
				require.NoError(t, err)

				err = clonedRepo.SetGitConfig(
					&gitparameteroptions.GitConfigSetOptions{
						Name:  "Test User",
						Email: "user@example.com",
					},
				)
				require.NoError(t, err)

				clonedRepo2 := getGitRepositoryToTest(tt.implementationCloned)
				defer clonedRepo2.Delete(verbose)
				err = clonedRepo2.Delete(verbose)
				require.NoError(t, err)

				err = clonedRepo2.CloneRepository(upstreamRepo, verbose)
				require.NoError(t, err)

				err = clonedRepo2.SetGitConfig(
					&gitparameteroptions.GitConfigSetOptions{
						Name:  "Test User2",
						Email: "user2@example.com",
					},
				)
				require.NoError(t, err)

				fileName := "abc.txt"
				fileName2 := "abc.txt"
				_, err = clonedRepo.CreateFileInDirectory(verbose, fileName)
				require.NoError(t, err)

				_, err = clonedRepo.CreateFileInDirectory(verbose, fileName2)
				require.NoError(t, err)

				err = clonedRepo.AddFilesByPath(
					[]string{fileName, fileName2},
					verbose,
				)
				require.NoError(t, err)

				_, err = clonedRepo.Commit(
					&gitparameteroptions.GitCommitOptions{
						Message: "another commit",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				err = clonedRepo.Push(ctx)
				require.NoError(t, err)

				require.False(t, mustutils.Must(clonedRepo2.FileByPathExists(fileName, verbose)))
				require.False(t, mustutils.Must(clonedRepo2.FileByPathExists(fileName2, verbose)))

				err = clonedRepo2.Pull(ctx)
				require.NoError(t, err)
				require.True(t, mustutils.Must(clonedRepo2.FileByPathExists(fileName, verbose)))
				require.True(t, mustutils.Must(clonedRepo2.FileByPathExists(fileName2, verbose)))
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
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				require.False(t, mustutils.Must(repo.HasUncommittedChanges(verbose)))
				require.False(t, mustutils.Must(repo.FileByPathExists("hello.txt", verbose)))

				_, err := repo.CreateFileInDirectory(verbose, "hello.txt")
				require.NoError(t, err)

				require.True(t, mustutils.Must(repo.FileByPathExists("hello.txt", verbose)))
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
				const verbose bool = true
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				require.False(t, mustutils.Must(repo.HasUncommittedChanges(verbose)))

				_, err := repo.CreateFileInDirectory(verbose, "a.txt")
				require.NoError(t, err)
				_, err = repo.CreateFileInDirectory(verbose, "b.txt")
				require.NoError(t, err)
				_, err = repo.CreateFileInDirectory(verbose, "c.txt")
				require.NoError(t, err)
				_, err = repo.CreateFileInDirectory(verbose, "cb.txt")
				require.NoError(t, err)

				files, err := repo.ListFiles(
					ctx,
					&parameteroptions.ListFileOptions{
						MatchBasenamePattern: []string{"^b\\.txt$"},
					},
				)
				require.NoError(t, err)
				require.Len(t, files, 1)

				baseName, err := files[0].GetBaseName()
				require.NoError(t, err)
				require.EqualValues(t, "b.txt", baseName)

				files, err = repo.ListFiles(
					ctx,
					&parameteroptions.ListFileOptions{
						MatchBasenamePattern: []string{"^.*b\\.txt$"},
					},
				)
				require.NoError(t, err)
				require.Len(t, files, 2)
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
				ctx := getCtx()
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(verbose)

				require.False(t, mustutils.Must(repo.HasUncommittedChanges(verbose)))

				_, err := repo.CreateFileInDirectory(verbose, "a.txt")
				require.NoError(t, err)
				_, err = repo.CreateFileInDirectory(verbose, "b.txt")
				require.NoError(t, err)
				_, err = repo.CreateFileInDirectory(verbose, "c.txt")
				require.NoError(t, err)
				_, err = repo.CreateFileInDirectory(verbose, "cb.txt")
				require.NoError(t, err)

				files, err := repo.ListFilePaths(
					ctx,
					&parameteroptions.ListFileOptions{
						MatchBasenamePattern: []string{"^b\\.txt$"},
						ReturnRelativePaths:  true,
					},
				)
				require.NoError(t, err)
				require.Len(t, files, 1)
				require.EqualValues(t, "b.txt", files[0])

				files, err = repo.ListFilePaths(
					ctx,
					&parameteroptions.ListFileOptions{
						MatchBasenamePattern: []string{"^.*b\\.txt$"},
						ReturnRelativePaths:  true,
					},
				)
				require.NoError(t, err)
				require.Len(t, files, 2)
				require.EqualValues(t, []string{"b.txt", "cb.txt"}, files)
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.Commit(
					&gitparameteroptions.GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)
				require.NoError(t, err)

				tagList, err := gitRepo.ListTagNames(verbose)
				require.NoError(t, err)
				require.Len(t, tagList, 0)

				expectedTags := []string{}
				for i := 0; i < 5; i++ {

					tagNameToAdd := "ExampleTag" + strconv.Itoa(i)

					_, err = gitRepo.CreateTag(
						&gitparameteroptions.GitRepositoryCreateTagOptions{
							TagName: tagNameToAdd,
							Verbose: verbose,
						},
					)
					require.NoError(t, err)
					expectedTags = append(expectedTags, tagNameToAdd)

					tagList, err := gitRepo.ListTagNames(verbose)
					require.NoError(t, err)

					require.EqualValues(t, expectedTags, tagList)
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.Commit(
					&gitparameteroptions.GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "abc",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "abcd",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				tags, err := gitRepo.ListTags(verbose)
				require.NoError(t, err)

				require.EqualValues(t, "abc", tags[0].MustGetName())
				require.EqualValues(t, "abcd", tags[1].MustGetName())
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.Commit(
					&gitparameteroptions.GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)
				require.NoError(t, err)

				require.Nil(t, mustutils.Must(gitRepo.GetLatestTagVersionOrNilIfNotFound(verbose)))

				_, err = gitRepo.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				require.NotNil(t, mustutils.Must(gitRepo.GetLatestTagVersionOrNilIfNotFound(verbose)))

				_, err = gitRepo.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				require.NotNil(t, mustutils.Must(gitRepo.GetLatestTagVersionOrNilIfNotFound(verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.Commit(
					&gitparameteroptions.GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				latestVersion, err := gitRepo.GetLatestTagVersion(verbose)
				require.NoError(t, err)

				require.EqualValues(t, "v1.0.0", mustutils.Must(latestVersion.GetAsString()))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.Commit(
					&gitparameteroptions.GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, "v1.0.0", mustutils.Must(gitRepo.GetLatestTagVersionAsString(verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.Commit(
					&gitparameteroptions.GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				newestVersion, err := gitRepo.GetCurrentCommitsNewestVersion(verbose)
				require.NoError(t, err)
				newestVersionString, err := newestVersion.GetAsString()
				require.NoError(t, err)
				require.EqualValues(t, "v0.1.2", newestVersionString)

				_, err = gitRepo.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				newestVersion, err = gitRepo.GetCurrentCommitsNewestVersion(verbose)
				require.NoError(t, err)
				newestVersionString, err = newestVersion.GetAsString()
				require.NoError(t, err)
				require.EqualValues(t, "v1.0.0", newestVersionString)
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.Commit(
					&gitparameteroptions.GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
						Verbose:    verbose,
					},
				)
				require.NoError(t, err)

				require.Nil(t, mustutils.Must(gitRepo.GetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose)))

				_, err = gitRepo.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				newestVersion, err := gitRepo.GetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose)
				require.NoError(t, err)
				require.EqualValues(t, "v0.1.2", mustutils.Must(newestVersion.GetAsString()))

				_, err = gitRepo.CreateTag(
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				newestVersion, err = gitRepo.GetCurrentCommitsNewestVersionOrNilIfNotPresent(verbose)
				require.NoError(t, err)
				require.EqualValues(t, "v1.0.0", mustutils.Must(newestVersion.GetAsString()))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				require.False(t, mustutils.Must(gitRepo.IsGolangApplication(verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				require.NotNil(t, gitRepo.CheckIsGolangApplication(verbose))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.WriteStringToFile("module example\n", verbose, "go.mod")
				require.NoError(t, err)

				require.False(t, mustutils.Must(gitRepo.IsGolangApplication(verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.WriteStringToFile("module example\n", verbose, "go.mod")
				require.NoError(t, err)

				require.NotNil(t, gitRepo.CheckIsGolangApplication(verbose))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.WriteStringToFile("module example\n", verbose, "go.mod")
				require.NoError(t, err)

				_, err = gitRepo.WriteStringToFile("package main\nfunc abc() bool {\n\treturn true\n}\n", verbose, "main.go")
				require.NoError(t, err)

				require.NotNil(t, gitRepo.CheckIsGolangApplication(verbose))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.WriteStringToFile("module example\n", verbose, "go.mod")
				require.NoError(t, err)
				_, err = gitRepo.WriteStringToFile("package main\nfunc abc() bool {\n\treturn true\n}\n", verbose, "main.go")
				require.NoError(t, err)

				require.False(t, mustutils.Must(gitRepo.IsGolangApplication(verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.WriteStringToFile("module example\n", verbose, "go.mod")
				require.NoError(t, err)
				_, err = gitRepo.WriteStringToFile("package main\nfunc main() {\n\treturn\n}\n", verbose, "main.go")
				require.NoError(t, err)

				require.True(t, mustutils.Must(gitRepo.IsGolangApplication(verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.WriteStringToFile("module example\n", verbose, "go.mod")
				require.NoError(t, err)

				_, err = gitRepo.WriteStringToFile("package main\nfunc main() {\n\treturn\n}\n", verbose, "main.go")
				require.NoError(t, err)

				require.Nil(t, gitRepo.CheckIsGolangApplication(verbose))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.WriteStringToFile("hello world\n", verbose, "test.txt")
				require.NoError(t, err)

				testTxtFile, err := gitRepo.GetFileByPath("test.txt")
				require.NoError(t, err)

				require.EqualValues(t, "hello world\n", testTxtFile.MustReadAsString())
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				require.False(t, mustutils.Must(gitRepo.IsGolangPackage(verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.WriteStringToFile("module example\n", verbose, "go.mod")
				require.NoError(t, err)

				require.True(t, mustutils.Must(gitRepo.IsGolangPackage(verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.WriteStringToFile("module example\n", verbose, "go.mod")
				require.NoError(t, err)
				_, err = gitRepo.WriteStringToFile("package main\nfunc main() {\n\treturn\n}\n", verbose, "main.go")
				require.NoError(t, err)

				require.False(t, mustutils.Must(gitRepo.IsGolangPackage(verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.WriteStringToFile("module example\n", verbose, "go.mod")
				require.NoError(t, err)

				require.Nil(t, gitRepo.CheckIsGolangPackage(verbose))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.WriteStringToFile("module example\n", verbose, "go.mod")
				require.NoError(t, err)
				_, err = gitRepo.WriteStringToFile("package main\nfunc main() {\n\treturn\n}\n", verbose, "main.go")
				require.NoError(t, err)

				require.NotNil(t, gitRepo.CheckIsGolangPackage(verbose))
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
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				repoRootDirectory, err := gitRepo.GetRootDirectory(ctx)
				require.NoError(t, err)

				gitRepo2 := MustGetGitRepositoryByDirectory(repoRootDirectory)

				require.EqualValues(
					t,
					mustutils.Must(gitRepo.GetRootDirectoryPath(ctx)),
					mustutils.Must(gitRepo2.GetRootDirectoryPath(ctx)),
				)

				require.Nil(t, gitRepo.CheckHasNoUncommittedChanges(verbose))

				require.Nil(t, gitRepo2.CheckHasNoUncommittedChanges(verbose))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				defaultBranchName, err := gitRepo.GetCurrentBranchName(verbose)
				require.NoError(t, err)

				require.False(t, mustutils.Must(gitRepo.BranchByNameExists("testbranch", verbose)))

				for i := 0; i < 2; i++ {
					err = gitRepo.CreateBranch(
						&parameteroptions.CreateBranchOptions{
							Name:    "testbranch",
							Verbose: verbose,
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(gitRepo.BranchByNameExists("testbranch", verbose)))
				}

				err = gitRepo.CheckoutBranchByName(defaultBranchName, verbose)
				require.NoError(t, err)

				for i := 0; i < 2; i++ {
					err = gitRepo.DeleteBranchByName("testbranch", verbose)
					require.NoError(t, err)

					require.False(t, mustutils.Must(gitRepo.BranchByNameExists("testbranch", verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				for _, branchName := range []string{"testbranch1", "testbranch2"} {
					err := gitRepo.CreateBranch(
						&parameteroptions.CreateBranchOptions{
							Name:    branchName,
							Verbose: verbose,
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(gitRepo.BranchByNameExists(branchName, verbose)))

					err = gitRepo.CheckoutBranchByName(branchName, verbose)
					require.NoError(t, err)
					require.EqualValues(t, branchName, mustutils.Must(gitRepo.GetCurrentBranchName(verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.Commit(
					&gitparameteroptions.GitCommitOptions{
						AllowEmpty: true,
						Message:    "commit message",
						Verbose:    verbose,
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, "commit message", mustutils.Must(gitRepo.GetCurrentCommitMessage(verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				_, err := gitRepo.Commit(
					&gitparameteroptions.GitCommitOptions{
						AllowEmpty: true,
						Message:    "commit before testing",
						Verbose:    verbose,
					},
				)
				require.NoError(t, err)

				require.EqualValues(t,
					"commit before testing",
					mustutils.Must(gitRepo.GetCurrentCommitMessage(verbose)),
				)

				_, err = gitRepo.CommitIfUncommittedChanges(
					&gitparameteroptions.GitCommitOptions{
						Message: "This should not trigger a commit",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, "commit before testing", mustutils.Must(gitRepo.GetCurrentCommitMessage(verbose)))

				_, err = gitRepo.WriteStringToFile("hello", verbose, "world.txt")
				require.NoError(t, err)

				err = gitRepo.AddFileByPath("world.txt", verbose)
				require.NoError(t, err)

				_, err = gitRepo.CommitIfUncommittedChanges(
					&gitparameteroptions.GitCommitOptions{
						Message: "This should trigger a commit",
						Verbose: verbose,
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, "This should trigger a commit", mustutils.Must(gitRepo.GetCurrentCommitMessage(verbose)))

				_, err = gitRepo.WriteStringToFile("world", verbose, "world.txt")
				require.NoError(t, err)

				// world.txt is already known in the git repo.
				// No need to explicitly add world.txt again.

				_, err = gitRepo.CommitIfUncommittedChanges(
					&gitparameteroptions.GitCommitOptions{
						Message: "This should trigger again a commit",
						Verbose: verbose,
					},
				)

				require.EqualValues(t, "This should trigger again a commit", mustutils.Must(gitRepo.GetCurrentCommitMessage(verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				require.False(t, mustutils.Must(gitRepo.RemoteByNameExists("example", verbose)))

				for i := 0; i < 2; i++ {
					err := gitRepo.AddRemote(
						&gitparameteroptions.GitRemoteAddOptions{
							RemoteName: "example",
							RemoteUrl:  "https://remote.url.example.com",
							Verbose:    verbose,
						},
					)
					require.NoError(t, err)

					require.True(t, mustutils.Must(gitRepo.RemoteByNameExists("example", verbose)))
				}

				for i := 0; i < 2; i++ {
					err := gitRepo.RemoveRemoteByName("example", verbose)
					require.NoError(t, err)
					require.False(t, mustutils.Must(gitRepo.RemoteByNameExists("example", verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				require.False(t, mustutils.Must(gitRepo.IsPreCommitRepository(verbose)))

				_, err := gitRepo.CreateSubDirectory("pre_commit_hooks", verbose)
				require.NoError(t, err)

				require.True(t, mustutils.Must(gitRepo.IsPreCommitRepository(verbose)))
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				require.NotNil(t, gitRepo.CheckIsPreCommitRepository(verbose))

				_, err := gitRepo.CreateSubDirectory("pre_commit_hooks", verbose)
				require.NoError(t, err)

				require.Nil(t, gitRepo.CheckIsPreCommitRepository(verbose))
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

				err := gitRepo.Delete(true)
				require.NoError(t, err)

				require.Error(t, gitRepo.CheckExists(ctx))
			},
		)
	}
}
