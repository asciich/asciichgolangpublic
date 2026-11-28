package precommitutils_test

import (
	"context"

	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/commandexecutorgitoo"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/localgitoo"
	precommitutils "github.com/asciich/asciichgolangpublic/pkg/gitutils/precommitutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getGitRepositoryToTest(implementationName string) (repo gitinterfaces.GitRepository) {
	ctx := getCtx()

	if implementationName == "localGitRepository" {
		repo = localgitoo.MustGetLocalGitRepositoryFromDirectory(
			mustutils.Must(tempfilesoo.CreateEmptyTemporaryDirectory(ctx)),
		)
	} else if implementationName == "localCommandExecutorRepository" {
		repo = mustutils.Must(commandexecutorgitoo.CreateLocalTemporaryRepository(ctx, &parameteroptions.CreateRepositoryOptions{}))
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

func TestPreCommitConfigFile_UpdateDependency(t *testing.T) {
	ctx := getCtx()
	type TestCase struct {
		testDataDir string
	}

	tests := []TestCase{}

	testDataDirectory, err := localgitoo.MustGetLocalGitRepositoryByPath(".").GetDirectoryByPath(ctx, "testdata", "PreCommitConfigFile", "UpdateDependency")
	require.NoError(t, err)
	for _, testDirectory := range mustutils.Must(testDataDirectory.ListSubDirectories(ctx, &parameteroptions.ListDirectoryOptions{Recursive: false})) {
		localPath, err := testDirectory.GetLocalPath()
		require.NoError(t, err)
		tests = append(tests, TestCase{localPath})
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				inputFile := precommitutils.MustGetPreCommitConfigByLocalPath(filepath.Join(tt.testDataDir, "input"))
				preCommitFile := precommitutils.MustGetPreCommitConfigByFile(mustutils.Must(tempfilesoo.CreateTemporaryFileFromFile(ctx, inputFile)))
				defer preCommitFile.Delete(ctx, &filesoptions.DeleteOptions{})

				expectedOutput := precommitutils.MustGetPreCommitConfigByLocalPath(filepath.Join(tt.testDataDir, "expected_output"))

				exists, err := expectedOutput.Exists(ctx)
				require.NoError(t, err)
				if !exists {
					if os.Getenv("UPDATE_EXPECTED") == "1" {
						err := expectedOutput.Create(ctx, &filesoptions.CreateOptions{})
						require.NoError(t, err)
					}
				}

				dependency := precommitutils.NewDependencyGitRepository()
				err = dependency.SetUrl("https://gitlab.asciich.ch/gitlab_management/pre-commit")
				require.NoError(t, err)
				err = dependency.SetVersionString("v0.1.0")
				require.NoError(t, err)
				err = dependency.SetSourceFiles([]filesinterfaces.File{preCommitFile})
				require.NoError(t, err)
				err = dependency.SetTargetVersionString("v0.10.0")
				require.NoError(t, err)

				_, err = preCommitFile.UpdateDependency(
					ctx,
					dependency,
					&parameteroptions.UpdateDependenciesOptions{
						Commit: false,
					},
				)
				require.NoError(t, err)

				updatedSha, err := preCommitFile.GetSha256Sum(ctx)
				require.NoError(t, err)
				expectedOutputSha, err := expectedOutput.GetSha256Sum(ctx)
				require.NoError(t, err)

				if expectedOutputSha != updatedSha {
					if os.Getenv("UPDATE_EXPECTED") == "1" {
						err := preCommitFile.CopyToFile(ctx, expectedOutput, &filesoptions.CopyOptions{})
						require.NoError(t, err)
					}
				}

				require.EqualValues(t, expectedOutputSha, updatedSha)
			},
		)
	}
}

func TestPreCommitConfigFile_GetPreCommitConfigInGitRepository(t *testing.T) {
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

				outFile, err := gitRepo.WriteStringToFile(ctx, ".pre-commit-config.yaml", "# placeholder", &filesoptions.WriteOptions{})
				require.NoError(t, err)
				require.NotNil(t, outFile)

				preCommitConfigFile := precommitutils.MustGetPreCommitConfigFileInGitRepository(gitRepo)

				exists, err := preCommitConfigFile.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)
				got, err := preCommitConfigFile.GetPath()
				require.NoError(t, err)
				require.True(t, strings.HasSuffix(got, "/.pre-commit-config.yaml"))
			},
		)
	}
}

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}
