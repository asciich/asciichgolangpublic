package asciichgolangpublic

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
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
	ctx := getCtx()

	if implementationName == "localGitRepository" {
		repo = MustGetLocalGitReposioryFromDirectory(
			mustutils.Must(tempfilesoo.CreateEmptyTemporaryDirectory(ctx)),
		)
	} else if implementationName == "localCommandExecutorRepository" {
		repo = mustutils.Must(GetLocalCommandExecutorGitRepositoryByDirectory(
			mustutils.Must(tempfilesoo.CreateEmptyTemporaryDirectory(ctx)),
		))
	} else {
		logging.LogFatalWithTracef("unknown implementationName='%s'", implementationName)
	}

	mustutils.Must0(repo.Init(
		ctx,
		&parameteroptions.CreateRepositoryOptions{
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				err := repo.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				require.False(t, mustutils.Must(repo.Exists(ctx)))
				require.False(t, mustutils.Must(repo.IsInitialized(ctx)))
				require.False(t, mustutils.Must(repo.HasInitialCommit(ctx)))

				for i := 0; i < 2; i++ {
					err = repo.Init(
						ctx,
						&parameteroptions.CreateRepositoryOptions{
							BareRepository: tt.bareRepository,
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(repo.Exists(ctx)))
					require.True(t, mustutils.Must(repo.IsInitialized(ctx)))
					require.False(t, mustutils.Must(repo.HasInitialCommit(ctx)))
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				err := repo.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				// An non existing directory is not a git repository:
				isRepo, err := repo.IsGitRepository(ctx)
				require.NoError(t, err)
				require.False(t, isRepo)

				path, err := repo.GetPath()
				require.NoError(t, err)
				_, err = files.Directories().CreateLocalDirectoryByPath(ctx, path, &filesoptions.CreateOptions{})
				require.NoError(t, err)

				// The directory exists but is empty which is not a git directory:
				isRepo, err = repo.IsGitRepository(ctx)
				require.NoError(t, err)
				require.False(t, isRepo)

				for i := 0; i < 2; i++ {
					mustutils.Must0(repo.Init(
						ctx,
						&parameteroptions.CreateRepositoryOptions{
							BareRepository: tt.bareRepository,
						},
					))
					isRepo, err = repo.IsGitRepository(ctx)
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				for i := 0; i < 2; i++ {
					err := repo.Delete(ctx, &filesoptions.DeleteOptions{})
					require.NoError(t, err)
					require.False(t, mustutils.Must(repo.Exists(ctx)))
					require.False(t, mustutils.Must(repo.IsInitialized(ctx)))
					require.False(t, mustutils.Must(repo.HasInitialCommit(ctx)))
				}

				for i := 0; i < 2; i++ {
					err := repo.Create(ctx, &filesoptions.CreateOptions{})
					require.NoError(t, err)
					require.True(t, mustutils.Must(repo.Exists(ctx)))
					require.False(t, mustutils.Must(repo.IsInitialized(ctx)))
					require.False(t, mustutils.Must(repo.HasInitialCommit(ctx)))
				}

				for i := 0; i < 2; i++ {
					err := repo.Init(ctx, &parameteroptions.CreateRepositoryOptions{})
					require.NoError(t, err)
					require.True(t, mustutils.Must(repo.Exists(ctx)))
					require.True(t, mustutils.Must(repo.IsInitialized(ctx)))
					require.False(t, mustutils.Must(repo.HasInitialCommit(ctx)))
				}

				for i := 0; i < 2; i++ {
					err := repo.Init(
						ctx,
						&parameteroptions.CreateRepositoryOptions{
							InitializeWithEmptyCommit:   true,
							InitializeWithDefaultAuthor: true,
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(repo.Exists(ctx)))
					require.True(t, mustutils.Must(repo.IsInitialized(ctx)))
					require.True(t, mustutils.Must(repo.HasInitialCommit(ctx)))
				}

				for i := 0; i < 2; i++ {
					err := repo.Delete(ctx, &filesoptions.DeleteOptions{})
					require.NoError(t, err)
					require.False(t, mustutils.Must(repo.Exists(ctx)))
					require.False(t, mustutils.Must(repo.IsInitialized(ctx)))
					require.False(t, mustutils.Must(repo.HasInitialCommit(ctx)))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})
				repo.Delete(ctx, &filesoptions.DeleteOptions{})

				for i := 0; i < 2; i++ {
					err := repo.Init(
						ctx,
						&parameteroptions.CreateRepositoryOptions{
							InitializeWithDefaultAuthor: true,
							InitializeWithEmptyCommit:   true,
							BareRepository:              tt.bare,
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(repo.Exists(ctx)))
					require.True(t, mustutils.Must(repo.IsInitialized(ctx)))
					require.True(t, mustutils.Must(repo.HasInitialCommit(ctx)))
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				for i := 0; i < 2; i++ {
					err := repo.Delete(ctx, &filesoptions.DeleteOptions{})
					require.NoError(t, err)
					require.False(t, mustutils.Must(repo.Exists(ctx)))
				}

				for i := 0; i < 2; i++ {
					err := repo.CreateAndInit(
						ctx,
						&parameteroptions.CreateRepositoryOptions{
							InitializeWithEmptyCommit:   true,
							InitializeWithDefaultAuthor: true,
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(repo.Exists(ctx)))
					require.True(t, mustutils.Must(repo.IsInitialized(ctx)))
				}

				for i := 0; i < 2; i++ {
					err := repo.Delete(ctx, &filesoptions.DeleteOptions{})
					require.NoError(t, err)
					require.False(t, mustutils.Must(repo.Exists(ctx)))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				require.True(t, mustutils.Must(repo.HasNoUncommittedChanges(ctx)))

				_, err := repo.CreateFileInDirectory(ctx, "hello.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				require.False(t, mustutils.Must(repo.HasNoUncommittedChanges(ctx)))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				require.False(t, mustutils.Must(repo.HasUncommittedChanges(ctx)))

				_, err := repo.CreateFileInDirectory(ctx, "hello.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				require.True(t, mustutils.Must(repo.HasUncommittedChanges(ctx)))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				require.Nil(t, repo.CheckHasNoUncommittedChanges(ctx))

				_, err := repo.CreateFileInDirectory(ctx, "hello.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				require.NotNil(t, repo.CheckHasNoUncommittedChanges(ctx))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				err := repo.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				err = repo.Init(
					ctx,
					&parameteroptions.CreateRepositoryOptions{
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				err := repo.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				err = repo.Init(
					ctx,
					&parameteroptions.CreateRepositoryOptions{
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				err := repo.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				err = repo.Init(
					ctx,
					&parameteroptions.CreateRepositoryOptions{
						BareRepository: tt.bareRepository,
					},
				)
				require.NoError(t, err)

				expectedRootDirectory, err := repo.GetPath()
				require.NoError(t, err)

				subDir, err := repo.CreateSubDirectory(ctx, "sub_directory", &filesoptions.CreateOptions{})
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				upstreamRepo := getGitRepositoryToTest(tt.implementationName)
				defer upstreamRepo.Delete(ctx, &filesoptions.DeleteOptions{})
				err := upstreamRepo.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				err = upstreamRepo.Init(
					ctx,
					&parameteroptions.CreateRepositoryOptions{
						BareRepository:              true,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				clonedRepo := getGitRepositoryToTest(tt.implementationName)
				defer clonedRepo.Delete(ctx, &filesoptions.DeleteOptions{})
				err = clonedRepo.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				for i := 0; i < 2; i++ {
					err = clonedRepo.CloneRepository(ctx, upstreamRepo)
					require.NoError(t, err)
				}

				require.EqualValues(t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(ctx)),
					mustutils.Must(clonedRepo.GetCurrentCommitHash(ctx)),
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				upstreamRepo := getGitRepositoryToTest(tt.implementationUpstream)
				defer upstreamRepo.Delete(ctx, &filesoptions.DeleteOptions{})
				err := upstreamRepo.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)
				err = upstreamRepo.Init(
					ctx,
					&parameteroptions.CreateRepositoryOptions{
						BareRepository:              true,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				clonedRepo := getGitRepositoryToTest(tt.implementationCloned)
				defer clonedRepo.Delete(ctx, &filesoptions.DeleteOptions{})
				err = clonedRepo.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				err = clonedRepo.CloneRepository(ctx, upstreamRepo)
				require.NoError(t, err)

				err = clonedRepo.SetGitConfig(
					ctx,
					&gitparameteroptions.GitConfigSetOptions{
						Name:  "Test User",
						Email: "user@example.com",
					},
				)
				require.NoError(t, err)

				clonedRepo2 := getGitRepositoryToTest(tt.implementationCloned2)
				defer clonedRepo2.Delete(ctx, &filesoptions.DeleteOptions{})
				err = clonedRepo2.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)
				err = clonedRepo2.CloneRepository(ctx, upstreamRepo)
				require.NoError(t, err)

				err = clonedRepo2.SetGitConfig(
					ctx,
					&gitparameteroptions.GitConfigSetOptions{
						Name:  "Test User2",
						Email: "user2@example.com",
					},
				)
				require.NoError(t, err)

				require.EqualValues(
					t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(ctx)),
					mustutils.Must(clonedRepo.GetCurrentCommitHash(ctx)),
				)

				require.EqualValues(
					t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(ctx)),
					mustutils.Must(clonedRepo2.GetCurrentCommitHash(ctx)),
				)

				fileName := "abc.txt"
				_, err = clonedRepo2.CreateFileInDirectory(ctx, fileName, &filesoptions.CreateOptions{})
				require.NoError(t, err)

				err = clonedRepo2.AddFileByPath(ctx, fileName)
				require.NoError(t, err)

				_, err = clonedRepo2.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message: "another commit",
					},
				)
				require.NoError(t, err)

				require.NotEqualValues(t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(ctx)),
					mustutils.Must(clonedRepo2.GetCurrentCommitHash(ctx)),
				)

				require.NotEqualValues(t,
					mustutils.Must(clonedRepo.GetCurrentCommitHash(ctx)),
					mustutils.Must(clonedRepo2.GetCurrentCommitHash(ctx)),
				)

				mustutils.Must0(clonedRepo2.Push(ctx))
				require.EqualValues(t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(ctx)),
					mustutils.Must(clonedRepo2.GetCurrentCommitHash(ctx)),
				)
				require.NotEqualValues(t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(ctx)),
					mustutils.Must(clonedRepo.GetCurrentCommitHash(ctx)),
				)

				err = clonedRepo.Pull(ctx)
				require.NoError(t, err)
				require.EqualValues(t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(ctx)),
					mustutils.Must(clonedRepo2.GetCurrentCommitHash(ctx)),
				)
				require.EqualValues(t,
					mustutils.Must(upstreamRepo.GetCurrentCommitHash(ctx)),
					mustutils.Must(clonedRepo.GetCurrentCommitHash(ctx)),
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				upstreamRepo := getGitRepositoryToTest(tt.implementationUpstream)
				defer upstreamRepo.Delete(ctx, &filesoptions.DeleteOptions{})
				err := upstreamRepo.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				err = upstreamRepo.Init(
					ctx,
					&parameteroptions.CreateRepositoryOptions{
						BareRepository:              true,
						InitializeWithEmptyCommit:   true,
						InitializeWithDefaultAuthor: true,
					},
				)
				require.NoError(t, err)

				clonedRepo := getGitRepositoryToTest(tt.implementationCloned)
				defer clonedRepo.Delete(ctx, &filesoptions.DeleteOptions{})
				err = clonedRepo.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				err = clonedRepo.CloneRepository(ctx, upstreamRepo)
				require.NoError(t, err)

				err = clonedRepo.SetGitConfig(
					ctx,
					&gitparameteroptions.GitConfigSetOptions{
						Name:  "Test User",
						Email: "user@example.com",
					},
				)
				require.NoError(t, err)

				clonedRepo2 := getGitRepositoryToTest(tt.implementationCloned)
				defer clonedRepo2.Delete(ctx, &filesoptions.DeleteOptions{})
				err = clonedRepo2.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				err = clonedRepo2.CloneRepository(ctx, upstreamRepo)
				require.NoError(t, err)

				err = clonedRepo2.SetGitConfig(
					ctx,
					&gitparameteroptions.GitConfigSetOptions{
						Name:  "Test User2",
						Email: "user2@example.com",
					},
				)
				require.NoError(t, err)

				fileName := "abc.txt"
				fileName2 := "abc.txt"
				_, err = clonedRepo.CreateFileInDirectory(ctx, fileName, &filesoptions.CreateOptions{})
				require.NoError(t, err)

				_, err = clonedRepo.CreateFileInDirectory(ctx, fileName2, &filesoptions.CreateOptions{})
				require.NoError(t, err)

				err = clonedRepo.AddFilesByPath(ctx, []string{fileName, fileName2})
				require.NoError(t, err)

				_, err = clonedRepo.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message: "another commit",
					},
				)
				require.NoError(t, err)

				err = clonedRepo.Push(ctx)
				require.NoError(t, err)

				require.False(t, mustutils.Must(clonedRepo2.FileByPathExists(ctx, fileName)))
				require.False(t, mustutils.Must(clonedRepo2.FileByPathExists(ctx, fileName2)))

				err = clonedRepo2.Pull(ctx)
				require.NoError(t, err)
				require.True(t, mustutils.Must(clonedRepo2.FileByPathExists(ctx, fileName)))
				require.True(t, mustutils.Must(clonedRepo2.FileByPathExists(ctx, fileName2)))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				require.False(t, mustutils.Must(repo.HasUncommittedChanges(ctx)))
				require.False(t, mustutils.Must(repo.FileByPathExists(ctx, "hello.txt")))

				_, err := repo.CreateFileInDirectory(ctx, "hello.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)

				require.True(t, mustutils.Must(repo.FileByPathExists(ctx, "hello.txt")))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				require.False(t, mustutils.Must(repo.HasUncommittedChanges(ctx)))

				_, err := repo.CreateFileInDirectory(ctx, "a.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				_, err = repo.CreateFileInDirectory(ctx, "b.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				_, err = repo.CreateFileInDirectory(ctx, "c.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				_, err = repo.CreateFileInDirectory(ctx, "cb.txt", &filesoptions.CreateOptions{})
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const verbose bool = true

				repo := getGitRepositoryToTest(tt.implementationName)
				defer repo.Delete(ctx, &filesoptions.DeleteOptions{})

				require.False(t, mustutils.Must(repo.HasUncommittedChanges(ctx)))

				_, err := repo.CreateFileInDirectory(ctx, "a.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				_, err = repo.CreateFileInDirectory(ctx, "b.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				_, err = repo.CreateFileInDirectory(ctx, "c.txt", &filesoptions.CreateOptions{})
				require.NoError(t, err)
				_, err = repo.CreateFileInDirectory(ctx, "cb.txt", &filesoptions.CreateOptions{})
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
					},
				)
				require.NoError(t, err)

				tagList, err := gitRepo.ListTagNames(ctx)
				require.NoError(t, err)
				require.Len(t, tagList, 0)

				expectedTags := []string{}
				for i := 0; i < 5; i++ {

					tagNameToAdd := "ExampleTag" + strconv.Itoa(i)

					_, err = gitRepo.CreateTag(
						ctx,
						&gitparameteroptions.GitRepositoryCreateTagOptions{
							TagName: tagNameToAdd,
						},
					)
					require.NoError(t, err)
					expectedTags = append(expectedTags, tagNameToAdd)

					tagList, err := gitRepo.ListTagNames(ctx)
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
					},
				)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "abc",
					},
				)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "abcd",
					},
				)
				require.NoError(t, err)

				tags, err := gitRepo.ListTags(ctx)
				require.NoError(t, err)

				require.EqualValues(t, "abc", mustutils.Must(tags[0].GetName()))
				require.EqualValues(t, "abcd", mustutils.Must(tags[1].GetName()))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
					},
				)
				require.NoError(t, err)

				require.Nil(t, mustutils.Must(gitRepo.GetLatestTagVersionOrNilIfNotFound(ctx)))

				_, err = gitRepo.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
					},
				)
				require.NoError(t, err)

				require.NotNil(t, mustutils.Must(gitRepo.GetLatestTagVersionOrNilIfNotFound(ctx)))

				_, err = gitRepo.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
					},
				)
				require.NoError(t, err)

				require.NotNil(t, mustutils.Must(gitRepo.GetLatestTagVersionOrNilIfNotFound(ctx)))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
					},
				)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
					},
				)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
					},
				)
				require.NoError(t, err)

				latestVersion, err := gitRepo.GetLatestTagVersion(ctx)
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
					},
				)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
					},
				)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, "v1.0.0", mustutils.Must(gitRepo.GetLatestTagVersionAsString(ctx)))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
					},
				)
				require.NoError(t, err)

				_, err = gitRepo.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
					},
				)
				require.NoError(t, err)

				newestVersion, err := gitRepo.GetCurrentCommitsNewestVersion(ctx)
				require.NoError(t, err)
				newestVersionString, err := newestVersion.GetAsString()
				require.NoError(t, err)
				require.EqualValues(t, "v0.1.2", newestVersionString)

				_, err = gitRepo.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
					},
				)
				require.NoError(t, err)

				newestVersion, err = gitRepo.GetCurrentCommitsNewestVersion(ctx)
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message:    "initial empty commit",
						AllowEmpty: true,
					},
				)
				require.NoError(t, err)

				require.Nil(t, mustutils.Must(gitRepo.GetCurrentCommitsNewestVersionOrNilIfNotPresent(ctx)))

				_, err = gitRepo.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v0.1.2",
					},
				)
				require.NoError(t, err)

				newestVersion, err := gitRepo.GetCurrentCommitsNewestVersionOrNilIfNotPresent(ctx)
				require.NoError(t, err)
				require.EqualValues(t, "v0.1.2", mustutils.Must(newestVersion.GetAsString()))

				_, err = gitRepo.CreateTag(
					ctx,
					&gitparameteroptions.GitRepositoryCreateTagOptions{
						TagName: "v1.0.0",
					},
				)
				require.NoError(t, err)

				newestVersion, err = gitRepo.GetCurrentCommitsNewestVersionOrNilIfNotPresent(ctx)
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				require.False(t, mustutils.Must(gitRepo.IsGolangApplication(ctx)))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				require.NotNil(t, gitRepo.CheckIsGolangApplication(ctx))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.WriteStringToFile(ctx, "go.mod", "module example\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				require.False(t, mustutils.Must(gitRepo.IsGolangApplication(ctx)))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.WriteStringToFile(ctx, "go.mod", "module example\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				require.NotNil(t, gitRepo.CheckIsGolangApplication(ctx))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.WriteStringToFile(ctx, "go.mod", "module example\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				_, err = gitRepo.WriteStringToFile(ctx, "main.go", "package main\nfunc abc() bool {\n\treturn true\n}\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				require.NotNil(t, gitRepo.CheckIsGolangApplication(ctx))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.WriteStringToFile(ctx, "go.mod", "module example\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)
				_, err = gitRepo.WriteStringToFile(ctx, "main.go", "package main\nfunc abc() bool {\n\treturn true\n}\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				require.False(t, mustutils.Must(gitRepo.IsGolangApplication(ctx)))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.WriteStringToFile(ctx, "go.mod", "module example\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)
				_, err = gitRepo.WriteStringToFile(ctx, "main.go", "package main\nfunc main() {\n\treturn\n}\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				require.True(t, mustutils.Must(gitRepo.IsGolangApplication(ctx)))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.WriteStringToFile(ctx, "go.mod", "module example\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				_, err = gitRepo.WriteStringToFile(ctx, "main.go", "package main\nfunc main() {\n\treturn\n}\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				require.Nil(t, gitRepo.CheckIsGolangApplication(ctx))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.WriteStringToFile(ctx, "test.txt", "hello world\n", &filesoptions.WriteOptions{})
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				require.False(t, mustutils.Must(gitRepo.IsGolangPackage(ctx)))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.WriteStringToFile(ctx, "go.mod", "module example\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				require.True(t, mustutils.Must(gitRepo.IsGolangPackage(ctx)))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.WriteStringToFile(ctx, "go.mod", "module example\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)
				_, err = gitRepo.WriteStringToFile(ctx, "main.go", "package main\nfunc main() {\n\treturn\n}\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				require.False(t, mustutils.Must(gitRepo.IsGolangPackage(ctx)))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				require.Error(t, gitRepo.CheckIsGolangPackage(ctx))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.WriteStringToFile(ctx, "go.mod", "module example\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				require.Nil(t, gitRepo.CheckIsGolangPackage(ctx))
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
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.WriteStringToFile(ctx, "go.mod", "module example\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)
				_, err = gitRepo.WriteStringToFile(ctx, "main.go", "package main\nfunc main() {\n\treturn\n}\n", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				require.NotNil(t, gitRepo.CheckIsGolangPackage(ctx))
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				repoRootDirectory, err := gitRepo.GetRootDirectory(ctx)
				require.NoError(t, err)

				gitRepo2 := MustGetGitRepositoryByDirectory(repoRootDirectory)

				require.EqualValues(
					t,
					mustutils.Must(gitRepo.GetRootDirectoryPath(ctx)),
					mustutils.Must(gitRepo2.GetRootDirectoryPath(ctx)),
				)

				require.Nil(t, gitRepo.CheckHasNoUncommittedChanges(ctx))

				require.Nil(t, gitRepo2.CheckHasNoUncommittedChanges(ctx))
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				defaultBranchName, err := gitRepo.GetCurrentBranchName(ctx)
				require.NoError(t, err)

				require.False(t, mustutils.Must(gitRepo.BranchByNameExists(ctx, "testbranch")))

				for i := 0; i < 2; i++ {
					err = gitRepo.CreateBranch(
						ctx,
						&parameteroptions.CreateBranchOptions{
							Name: "testbranch",
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(gitRepo.BranchByNameExists(ctx, "testbranch")))
				}

				err = gitRepo.CheckoutBranchByName(ctx, defaultBranchName)
				require.NoError(t, err)

				for i := 0; i < 2; i++ {
					err = gitRepo.DeleteBranchByName(ctx, "testbranch")
					require.NoError(t, err)

					require.False(t, mustutils.Must(gitRepo.BranchByNameExists(ctx, "testbranch")))
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				for _, branchName := range []string{"testbranch1", "testbranch2"} {
					err := gitRepo.CreateBranch(
						ctx,
						&parameteroptions.CreateBranchOptions{
							Name: branchName,
						},
					)
					require.NoError(t, err)
					require.True(t, mustutils.Must(gitRepo.BranchByNameExists(ctx, branchName)))

					err = gitRepo.CheckoutBranchByName(ctx, branchName)
					require.NoError(t, err)
					require.EqualValues(t, branchName, mustutils.Must(gitRepo.GetCurrentBranchName(ctx)))
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						AllowEmpty: true,
						Message:    "commit message",
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, "commit message", mustutils.Must(gitRepo.GetCurrentCommitMessage(ctx)))
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				_, err := gitRepo.Commit(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						AllowEmpty: true,
						Message:    "commit before testing",
					},
				)
				require.NoError(t, err)

				require.EqualValues(t,
					"commit before testing",
					mustutils.Must(gitRepo.GetCurrentCommitMessage(ctx)),
				)

				_, err = gitRepo.CommitIfUncommittedChanges(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message: "This should not trigger a commit",
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, "commit before testing", mustutils.Must(gitRepo.GetCurrentCommitMessage(ctx)))

				_, err = gitRepo.WriteStringToFile(ctx, "world.txt", "hello", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				err = gitRepo.AddFileByPath(ctx, "world.txt")
				require.NoError(t, err)

				_, err = gitRepo.CommitIfUncommittedChanges(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message: "This should trigger a commit",
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, "This should trigger a commit", mustutils.Must(gitRepo.GetCurrentCommitMessage(ctx)))

				_, err = gitRepo.WriteStringToFile(ctx, "world.txt", "world", &filesoptions.WriteOptions{})
				require.NoError(t, err)

				// world.txt is already known in the git repo.
				// No need to explicitly add world.txt again.

				_, err = gitRepo.CommitIfUncommittedChanges(
					ctx,
					&gitparameteroptions.GitCommitOptions{
						Message: "This should trigger again a commit",
					},
				)
				require.NoError(t, err)

				require.EqualValues(t, "This should trigger again a commit", mustutils.Must(gitRepo.GetCurrentCommitMessage(ctx)))
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				require.False(t, mustutils.Must(gitRepo.RemoteByNameExists(ctx, "example")))

				for i := 0; i < 2; i++ {
					err := gitRepo.AddRemote(
						ctx,
						&gitparameteroptions.GitRemoteAddOptions{
							RemoteName: "example",
							RemoteUrl:  "https://remote.url.example.com",
							Verbose:    verbose,
						},
					)
					require.NoError(t, err)

					require.True(t, mustutils.Must(gitRepo.RemoteByNameExists(ctx, "example")))
				}

				for i := 0; i < 2; i++ {
					err := gitRepo.RemoveRemoteByName(ctx, "example")
					require.NoError(t, err)
					require.False(t, mustutils.Must(gitRepo.RemoteByNameExists(ctx, "example")))
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				require.False(t, mustutils.Must(gitRepo.IsPreCommitRepository(ctx)))

				_, err := gitRepo.CreateSubDirectory(ctx, "pre_commit_hooks", &filesoptions.CreateOptions{})
				require.NoError(t, err)

				require.True(t, mustutils.Must(gitRepo.IsPreCommitRepository(ctx)))
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				require.NotNil(t, gitRepo.CheckIsPreCommitRepository(ctx))

				_, err := gitRepo.CreateSubDirectory(ctx, "pre_commit_hooks", &filesoptions.CreateOptions{})
				require.NoError(t, err)

				require.Nil(t, gitRepo.CheckIsPreCommitRepository(ctx))
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := contextutils.GetVerbosityContextByBool(true)

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})

				require.NoError(t, gitRepo.CheckExists(ctx))

				err := gitRepo.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				require.Error(t, gitRepo.CheckExists(ctx))
			},
		)
	}
}
