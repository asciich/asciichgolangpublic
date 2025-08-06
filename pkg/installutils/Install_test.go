package installutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/installutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

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
				ctx := getCtx()
				const verbose = true

				sourceFile, err := tempfilesoo.CreateFromStringAndGetPath(ctx, tt.content)
				require.NoError(t, err)
				defer files.DeleteFileByPath(sourceFile, verbose)

				destFile, err := tempfilesoo.CreateEmptyTemporaryFile(ctx)
				require.NoError(t, err)
				destFilePath, err := destFile.GetPath()
				require.NoError(t, err)
				defer destFile.Delete(ctx, &filesoptions.DeleteOptions{})
				err = destFile.Delete(ctx, &filesoptions.DeleteOptions{})
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
