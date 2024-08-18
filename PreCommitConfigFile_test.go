package asciichgolangpublic

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreCommitConfigFileUpdateDependency(t *testing.T) {
	type TestCase struct {
		testDataDir string
	}

	tests := []TestCase{}

	testDataDirectory := MustGetLocalGitRepositoryByPath(".").MustGetSubDirectory("testdata", "PreCommitConfigFile", "UpdateDependency")
	for _, testDirectory := range testDataDirectory.MustGetSubDirectories(&ListDirectoryOptions{Recursive: false}) {
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

				if !expectedOutput.MustExists() {
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
