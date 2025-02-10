package files

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
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
				require := require.New(t)

				const verbose bool = true

				tempDir, err := os.MkdirTemp("", "tempdir_for_testing")
				require.Nil(t, err)

				var directory Directory = MustGetLocalDirectoryByPath(tempDir)
				defer directory.Delete(verbose)

				require.True(directory.MustExists(verbose))

				for i := 0; i < 2; i++ {
					directory.MustDelete(verbose)
					require.False(directory.MustExists(verbose))
				}

				for i := 0; i < 2; i++ {
					createdDir := Directories().MustCreateLocalDirectoryByPath(directory.MustGetLocalPath(), verbose)
					require.True(directory.MustExists(verbose))
					require.True(createdDir.MustExists(verbose))
				}

			},
		)
	}
}
