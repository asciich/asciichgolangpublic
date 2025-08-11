package installutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
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
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				const verbose = true

				sourceFile, err := tempfilesoo.CreateFromStringAndGetPath(ctx, tt.content)
				require.NoError(t, err)
				defer nativefiles.Delete(ctx, sourceFile, &filesoptions.DeleteOptions{})

				destFile, err := tempfilesoo.CreateEmptyTemporaryFile(ctx)
				require.NoError(t, err)
				destFilePath, err := destFile.GetPath()
				require.NoError(t, err)
				defer destFile.Delete(ctx, &filesoptions.DeleteOptions{})
				err = destFile.Delete(ctx, &filesoptions.DeleteOptions{})
				require.NoError(t, err)

				require.NoFileExists(t, destFilePath)

				err = installutils.Install(
					ctx,
					&installutils.InstallOptions{
						SrcPath:     sourceFile,
						InstallPath: destFilePath,
						Mode:        tt.mode,
					},
				)
				require.NoError(t, err)

				require.FileExists(t, destFilePath)

				accessPermissions, err := destFile.GetAccessPermissionsString()
				require.NoError(t, err)
				require.EqualValues(t, tt.mode, accessPermissions)
			},
		)
	}
}
