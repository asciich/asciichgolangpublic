package files_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// TODO test against all Directory implementations.
func TestDirectoriesCreateLocalDirectoryByPath(t *testing.T) {

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

				tempDir, err := os.MkdirTemp("", "tempdir_for_testing")
				require.Nil(t, err)

				var directory files.Directory = files.MustGetLocalDirectoryByPath(tempDir)
				defer directory.Delete(verbose)

				exists, err := directory.Exists(verbose)
				require.NoError(t, err)
				require.True(t, exists)

				for i := 0; i < 2; i++ {
					err = directory.Delete(verbose)
					require.NoError(t, err)

					exists, err = directory.Exists(verbose)
					require.NoError(t, err)
					require.False(t, exists)
				}

				for i := 0; i < 2; i++ {
					localPath, err := directory.GetLocalPath()
					require.NoError(t, err)
					createdDir := files.Directories().MustCreateLocalDirectoryByPath(localPath, verbose)
					
					dirExists, err:= directory.Exists(verbose)
					require.NoError(t, err)
					require.True(t, dirExists)

					createdExists, err := createdDir.Exists(verbose)
					require.NoError(t, err)
					require.True(t, createdExists)
				}

			},
		)
	}
}
