package files

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
)

// TODO move to File_test.go and run against all implementations.
func TestCommandExecutorFileReadAndWrite(t *testing.T) {
	tests := []struct {
		testContent string
	}{
		{"testcase"},
		{"testcase\n"},
		{"multyLine\nContent"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				temporaryFilePath := createTempFileAndGetPath()

				var fileToTest File = MustGetLocalCommandExecutorFileByPath(temporaryFilePath)
				defer fileToTest.Delete(verbose)

				require.True(fileToTest.MustExists(verbose))
				require.EqualValues(
					"",
					fileToTest.MustReadAsString(),
				)

				fileToTest.WriteString(tt.testContent, verbose)
				require.True(fileToTest.MustExists(verbose))
				require.EqualValues(
					tt.testContent,
					fileToTest.MustReadAsString(),
				)

				fileToTest.MustDelete(verbose)
				require.False(fileToTest.MustExists(verbose))
			},
		)
	}
}
