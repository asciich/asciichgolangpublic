package asciichgolangpublic

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/tempfiles"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestPreCommitConfigFile_UpdateDependency(t *testing.T) {
	type TestCase struct {
		testDataDir string
	}

	tests := []TestCase{}

	testDataDirectory := MustGetLocalGitRepositoryByPath(".").MustGetSubDirectory("testdata", "PreCommitConfigFile", "UpdateDependency")
	for _, testDirectory := range mustutils.Must(testDataDirectory.ListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false})) {
		tests = append(tests, TestCase{testDirectory.MustGetLocalPath()})
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				inputFile := MustGetPreCommitConfigByLocalPath(filepath.Join(tt.testDataDir, "input"))
				preCommitFile := MustGetPreCommitConfigByFile(tempfiles.MustCreateTemporaryFileFromFile(inputFile, verbose))
				defer preCommitFile.Delete(verbose)

				expectedOutput := MustGetPreCommitConfigByLocalPath(filepath.Join(tt.testDataDir, "expected_output"))

				if !expectedOutput.MustExists(verbose) {
					if os.Getenv("UPDATE_EXPECTED") == "1" {
						expectedOutput.MustCreate(verbose)
					}
				}

				dependency := &DependencyGitRepository{
					url:                 "https://gitlab.asciich.ch/gitlab_management/pre-commit",
					versionString:       "v0.1.0",
					sourceFiles:         []files.File{preCommitFile},
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

				require.EqualValues(expectedOutputSha, updatedSha)
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
				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				outFile, err := gitRepo.WriteStringToFile("# placeholder", verbose, ".pre-commit-config.yaml")
				require.NoError(t, err)
				require.NotNil(t, outFile)

				preCommitConfigFile := MustGetPreCommitConfigFileInGitRepository(gitRepo)
				require.True(t, preCommitConfigFile.MustExists(verbose))
				require.True(t, strings.HasSuffix(preCommitConfigFile.MustGetPath(), "/.pre-commit-config.yaml"))
			},
		)
	}
}
