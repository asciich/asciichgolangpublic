package filesutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// TestDirectory_GetBaseName tests Directory.GetBaseName method
func TestDirectory_GetBaseName(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativedirectoryoo"},
		{"commandExecutorFileExec"},
		{"commandExecutorFileBash"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				dirToTest := getTemporaryDirectoryToTest(tt.implementationName)
				defer dirToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				baseName, err := dirToTest.GetBaseName()
				require.NoError(t, err)
				require.NotEmpty(t, baseName)
			},
		)
	}
}

func TestDirectory_GetDirName(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativedirectoryoo"},
		{"commandExecutorFileExec"},
		{"commandExecutorFileBash"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()

				dirToTest := getTemporaryDirectoryToTest(tt.implementationName)
				defer dirToTest.Delete(ctx, &filesoptions.DeleteOptions{})

				dirName, err := dirToTest.GetDirName()
				require.NoError(t, err)
				require.NotEmpty(t, dirName)
			},
		)
	}
}

// TestDirectory_GetHostDescription tests Directory.GetHostDescription method
func TestDirectory_GetHostDescription(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"nativedirectoryoo"},
		{"commandExecutorFileExec"},
		{"commandExecutorFileBash"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				dirToTest := getTemporaryDirectoryToTest(tt.implementationName)
				defer dirToTest.Delete(getCtx(), &filesoptions.DeleteOptions{})

				hostDesc, err := dirToTest.GetHostDescription()
				require.NoError(t, err)
				require.NotEmpty(t, hostDesc)
			},
		)
	}
}
