package installutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/installutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestInstallFromPath(t *testing.T) {
	tests := []struct {
		content string
		mode    string
	}{
		{"hello_world", "u=rwx,g=,o="},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose = true

				sourceFile, err := tempfilesoo.CreateFromStringAndGetPath(tt.content, verbose)
				require.NoError(t, err)
				defer files.DeleteFileByPath(sourceFile, verbose)

				destFile, err := tempfilesoo.CreateEmptyTemporaryFile(verbose)
				require.NoError(t, err)
				destFilePath, err := destFile.GetPath()
				require.NoError(t, err)
				defer destFile.Delete(verbose)
				err = destFile.Delete(verbose)
				require.NoError(t, err)

				require.NoFileExists(t, destFilePath)

				installutils.MustInstall(
					&installutils.InstallOptions{
						SrcPath:     sourceFile,
						InstallPath: destFilePath,
						Mode:        tt.mode,
						Verbose:     verbose,
					},
				)

				require.FileExists(t, destFilePath)

				accessPermissions, err := destFile.GetAccessPermissionsString()
				require.NoError(t, err)
				require.EqualValues(t, tt.mode, accessPermissions)
			},
		)
	}
}
