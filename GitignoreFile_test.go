package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
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
				ctx := getCtx()

				gitignoreFilePath, err := tempfilesoo.CreateEmptyTemporaryFileAndGetPath(ctx)
				require.NoError(t, err)

				gitignoreFile, err := GetGitignoreFileByPath(gitignoreFilePath)
				require.NoError(t, err)
				defer gitignoreFile.Delete(ctx, &filesoptions.DeleteOptions{})

				for i := 0; i < 3; i++ {
					err := gitignoreFile.AddFileToIgnore("test", "comment", verbose)
					require.NoError(t, err)

					nEmptyLines, err := gitignoreFile.GetNumberOfNonEmptyLines()
					require.NoError(t, err)
					require.EqualValues(t, 2, nEmptyLines)
				}

				for i := 0; i < 3; i++ {
					err := gitignoreFile.AddFileToIgnore("test2", "comment2", verbose)
					require.NoError(t, err)

					nEmptyLines, err := gitignoreFile.GetNumberOfNonEmptyLines()
					require.NoError(t, err)
					require.EqualValues(t, 4, nEmptyLines)
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
				ctx := getCtx()

				nonExitstingFile, err := tempfilesoo.CreateEmptyTemporaryFile(ctx)
				require.NoError(t, err)
				err = nonExitstingFile.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				gitignoreFile, err := GetGitignoreFileByFile(nonExitstingFile)
				require.NoError(t, err)
				defer gitignoreFile.Delete(ctx, &filesoptions.DeleteOptions{})

				require.False(t, mustutils.Must(gitignoreFile.Exists(ctx)))

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
				ctx := getCtx()

				emptyFile, err := tempfilesoo.CreateEmptyTemporaryFile(ctx)
				require.NoError(t, err)

				gitignoreFile, err := GetGitignoreFileByFile(emptyFile)
				require.NoError(t, err)
				defer gitignoreFile.Delete(ctx,&filesoptions.DeleteOptions{})

				require.True(t, mustutils.Must(gitignoreFile.Exists(ctx)))
				containsIgnore, err := gitignoreFile.ContainsIgnore("abc")
				require.NoError(t, err)
				require.False(t, containsIgnore)
			},
		)
	}
}
