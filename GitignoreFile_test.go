package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/tempfiles"
	"github.com/asciich/asciichgolangpublic/testutils"
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
				assert := assert.New(t)

				const verbose bool = true

				gitignoreFilePath := tempfiles.MustCreateEmptyTemporaryFileAndGetPath(verbose)

				gitignoreFile := MustGetGitignoreFileByPath(gitignoreFilePath)
				defer gitignoreFile.Delete(verbose)

				for i := 0; i < 3; i++ {
					gitignoreFile.MustAddFileToIgnore("test", "comment", verbose)

					assert.EqualValues(2, gitignoreFile.MustGetNumberOfNonEmptyLines())
				}

				for i := 0; i < 3; i++ {
					gitignoreFile.MustAddFileToIgnore("test2", "comment2", verbose)

					assert.EqualValues(4, gitignoreFile.MustGetNumberOfNonEmptyLines())
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
				assert := assert.New(t)

				const verbose bool = true

				nonExitstingFile := tempfiles.MustCreateEmptyTemporaryFile(verbose)
				nonExitstingFile.MustDelete(verbose)

				gitignoreFile := MustGetGitignoreFileByFile(nonExitstingFile)
				defer gitignoreFile.Delete(verbose)

				assert.False(gitignoreFile.MustExists(verbose))

				assert.False(gitignoreFile.ContainsIgnore("abc"))
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
				assert := assert.New(t)

				const verbose bool = true

				emptyFile := tempfiles.MustCreateEmptyTemporaryFile(verbose)

				gitignoreFile := MustGetGitignoreFileByFile(emptyFile)
				defer gitignoreFile.Delete(verbose)

				assert.True(gitignoreFile.MustExists(verbose))

				assert.False(gitignoreFile.ContainsIgnore("abc"))
			},
		)
	}
}
