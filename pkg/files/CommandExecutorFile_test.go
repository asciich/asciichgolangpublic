package files_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
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
				ctx := getCtx()

				temporaryFilePath := createTempFileAndGetPath()

				var fileToTest filesinterfaces.File
				var err error
				fileToTest, err = files.GetLocalCommandExecutorFileByPath(temporaryFilePath)
				require.NoError(t, err)
				defer fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				exists, err := fileToTest.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				require.EqualValues(t, "", fileToTest.MustReadAsString())

				fileToTest.WriteString(ctx, tt.testContent, &filesoptions.WriteOptions{})

				exists, err = fileToTest.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				require.EqualValues(t, tt.testContent, fileToTest.MustReadAsString())

				err = fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				exists, err = fileToTest.Exists(ctx)
				require.NoError(t, err)
				require.False(t, exists)
			},
		)
	}
}
