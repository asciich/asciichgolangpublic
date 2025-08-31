package asciichgolangpublic

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestPreCommitConfigFile_UpdateDependency(t *testing.T) {
	type TestCase struct {
		testDataDir string
	}

	tests := []TestCase{}

	testDataDirectory := MustGetLocalGitRepositoryByPath(".").MustGetSubDirectory("testdata", "PreCommitConfigFile", "UpdateDependency")
	for _, testDirectory := range mustutils.Must(testDataDirectory.ListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false})) {
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

				inputFile := MustGetPreCommitConfigByLocalPath(filepath.Join(tt.testDataDir, "input"))
				preCommitFile := MustGetPreCommitConfigByFile(mustutils.Must(tempfilesoo.CreateTemporaryFileFromFile(ctx, inputFile)))
				defer preCommitFile.Delete(ctx, &filesoptions.DeleteOptions{})

				expectedOutput := MustGetPreCommitConfigByLocalPath(filepath.Join(tt.testDataDir, "expected_output"))

				exists, err := expectedOutput.Exists(ctx)
				require.NoError(t, err)
				if !exists {
					if os.Getenv("UPDATE_EXPECTED") == "1" {
						err := expectedOutput.Create(ctx, &filesoptions.CreateOptions{})
						require.NoError(t, err)
					}
				}

				dependency := &DependencyGitRepository{
					url:                 "https://gitlab.asciich.ch/gitlab_management/pre-commit",
					versionString:       "v0.1.0",
					sourceFiles:         []filesinterfaces.File{preCommitFile},
					targetVersionString: "v0.10.0",
				}

				preCommitFile.MustUpdateDependency(
					dependency,
					&parameteroptions.UpdateDependenciesOptions{
						Commit:  false,
						Verbose: verbose,
					},
				)

				updatedSha := preCommitFile.MustGetSha256Sum()
				expectedOutputSha := expectedOutput.MustGetSha256Sum()

				if expectedOutputSha != updatedSha {
					if os.Getenv("UPDATE_EXPECTED") == "1" {
						preCommitFile.MustCopyToFile(expectedOutput, verbose)
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

				preCommitConfigFile := MustGetPreCommitConfigFileInGitRepository(gitRepo)

				exists, err := preCommitConfigFile.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)
				require.True(t, strings.HasSuffix(preCommitConfigFile.MustGetPath(), "/.pre-commit-config.yaml"))
			},
		)
	}
}
