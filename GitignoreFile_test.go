package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestGitignoreFileAddFileToIgnore(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				gitignoreFilePath := tempfiles.MustCreateEmptyTemporaryFileAndGetPath(verbose)

				gitignoreFile, err := GetGitignoreFileByPath(gitignoreFilePath)
				require.NoError(t, err)
				defer gitignoreFile.Delete(verbose)

				for i := 0; i < 3; i++ {
					err := gitignoreFile.AddFileToIgnore("test", "comment", verbose)
					require.NoError(t, err)
					require.EqualValues(t, 2, gitignoreFile.MustGetNumberOfNonEmptyLines())
				}

				for i := 0; i < 3; i++ {
					err := gitignoreFile.AddFileToIgnore("test2", "comment2", verbose)
					require.NoError(t, err)
					require.EqualValues(t, 4, gitignoreFile.MustGetNumberOfNonEmptyLines())
				}
			},
		)
	}
}

func TestGitignoreFileContainsIgnoreOnNonExistingFile(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				nonExitstingFile := tempfiles.MustCreateEmptyTemporaryFile(verbose)
				nonExitstingFile.MustDelete(verbose)

				gitignoreFile, err := GetGitignoreFileByFile(nonExitstingFile)
				require.NoError(t, err)
				defer gitignoreFile.Delete(verbose)

				require.False(t, gitignoreFile.MustExists(verbose))

				containsIgnore, err := gitignoreFile.ContainsIgnore("abc")
				require.Error(t, err)
				require.False(t, containsIgnore)
			},
		)
	}
}

func TestGitignoreFileContainsIgnoreOnEmptyFile(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				emptyFile := tempfiles.MustCreateEmptyTemporaryFile(verbose)

				gitignoreFile, err := GetGitignoreFileByFile(emptyFile)
				require.NoError(t, err)
				defer gitignoreFile.Delete(verbose)

				require.True(t, gitignoreFile.MustExists(verbose))
				containsIgnore, err := gitignoreFile.ContainsIgnore("abc")
				require.NoError(t, err)
				require.False(t, containsIgnore)
			},
		)
	}
}
