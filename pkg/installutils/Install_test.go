package installutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/files"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/installutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tempfiles"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/testutils"
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

				sourceFile := tempfiles.MustCreateFromStringAndGetPath(tt.content, verbose)
				defer files.MustDeleteFileByPath(sourceFile, verbose)

				destFile := tempfiles.MustCreateEmptyTemporaryFile(verbose)
				destFilePath := destFile.MustGetPath()
				defer destFile.MustDelete(verbose)
				destFile.MustDelete(verbose)

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
				require.EqualValues(
					t,
					tt.mode,
					destFile.MustGetAccessPermissionsString(),
				)
			},
		)
	}
}
