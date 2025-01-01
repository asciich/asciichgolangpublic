package asciichgolangpublic

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreCommitConfigFile_UpdateDependency(t *testing.T) {
	type TestCase struct {
		testDataDir string
	}

	tests := []TestCase{}

	testDataDirectory := MustGetLocalGitRepositoryByPath(".").MustGetSubDirectory("testdata", "PreCommitConfigFile", "UpdateDependency")
	for _, testDirectory := range testDataDirectory.MustListSubDirectories(&ListDirectoryOptions{Recursive: false}) {
		tests = append(tests, TestCase{testDirectory.MustGetLocalPath()})
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				inputFile := MustGetPreCommitConfigByLocalPath(filepath.Join(tt.testDataDir, "input"))
				preCommitFile := MustGetPreCommitConfigByFile(TemporaryFiles().MustCreateTemporaryFileFromFile(inputFile, verbose))
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
					sourceFiles:         []File{preCommitFile},
					targetVersionString: "v0.10.0",
				}

				preCommitFile.MustUpdateDependency(
					dependency,
					&UpdateDependenciesOptions{
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

				assert.EqualValues(expectedOutputSha, updatedSha)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				gitRepo := getGitRepositoryToTest(tt.implementationName)
				defer gitRepo.Delete(verbose)

				gitRepo.MustWriteStringToFile("# placeholder", verbose, ".pre-commit-config.yaml")

				preCommitConfigFile := MustGetPreCommitConfigInGitRepository(gitRepo)
				assert.True(preCommitConfigFile.MustExists(verbose))
				assert.True(strings.HasSuffix(preCommitConfigFile.MustGetPath(), "/.pre-commit-config.yaml"))
			},
		)
	}
}
