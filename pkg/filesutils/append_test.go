package filesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// TestFile_AppendString tests AppendString and AppendBytes methods
func TestFile_AppendString(t *testing.T) {
	tests := []struct {
		implementationName string
		initialContent     string
		appendedContent    string
		expectedContent    string
	}{
		{"nativefilesoo", "", "hello", "hello"},
		{"nativefilesoo", "hello", " world", "hello world"},
		{"nativefilesoo", "line1\n", "line2\n", "line1\nline2\n"},
		{"commandExecutorFileExec", "", "hello", "hello"},
		{"commandExecutorFileExec", "hello", " world", "hello world"},
		{"commandExecutorFileBash", "", "hello", "hello"},
		{"commandExecutorFileBash", "hello", " world", "hello world"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				fileToTest := getTemporaryFileToTest(tt.implementationName)
				defer fileToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				// Write initial content
				if tt.initialContent != "" {
					err := fileToTest.WriteString(ctx, tt.initialContent, &filesoptions.WriteOptions{})
					require.NoError(t, err)
				}

				// Append content
				err := fileToTest.AppendString(ctx, tt.appendedContent)
				require.NoError(t, err)

				// Read and verify
				content, err := fileToTest.ReadAsString(ctx)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedContent, content)
			},
		)
	}
}
