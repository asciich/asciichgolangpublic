package files_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
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
				ctx := getCtx()

				tempDir, err := os.MkdirTemp("", "tempdir_for_testing")
				require.NoError(t, err)

				var directory filesinterfaces.Directory
				directory, err = files.GetLocalDirectoryByPath(tempDir)
				require.NoError(t, err)
				defer directory.Delete(ctx, &filesoptions.DeleteOptions{})

				exists, err := directory.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				for i := 0; i < 2; i++ {
					err = directory.Delete(ctx, &filesoptions.DeleteOptions{})
					require.NoError(t, err)

					exists, err = directory.Exists(ctx)
					require.NoError(t, err)
					require.False(t, exists)
				}

				for i := 0; i < 2; i++ {
					localPath, err := directory.GetLocalPath()
					require.NoError(t, err)

					createdDir, err := files.Directories().CreateLocalDirectoryByPath(ctx, localPath, &filesoptions.CreateOptions{})
					require.NoError(t, err)

					dirExists, err := directory.Exists(ctx)
					require.NoError(t, err)
					require.True(t, dirExists)

					createdExists, err := createdDir.Exists(ctx)
					require.NoError(t, err)
					require.True(t, createdExists)
				}

			},
		)
	}
}
