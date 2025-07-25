package gnupg

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// Return a temporary file of the given 'implementationName'.
//
// Use defer file.Delete(verbose) to after calling this function to ensure
// the file is deleted after the test is over.
func getFileToTest(implementationName string) (file filesinterfaces.File) {
	const verbose = true

	if implementationName == "localFile" {
		file = files.MustGetLocalFileByPath(
			mustutils.Must(tempfilesoo.CreateEmptyTemporaryFileAndGetPath(verbose)),
		)
	} else if implementationName == "localCommandExecutorFile" {
		file = files.MustGetLocalCommandExecutorFileByPath(
			mustutils.Must(tempfilesoo.CreateEmptyTemporaryFileAndGetPath(verbose)),
		)
	} else {
		logging.LogFatalWithTracef("unknown implementationName='%s'", implementationName)
	}

	return file
}

func TestGnuPg_SignAndValidate(t *testing.T) {
	testutils.SkipIfRunningInGithub(t)

	tests := []struct {
		implementationName string
		contentString      string
		asciiArmor         bool
	}{
		{"localFile", "hello world", true},
		{"localCommandExecutorFile", "hello world", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				toTest := getFileToTest(tt.implementationName)
				defer toTest.Delete(verbose)

				signatureFile, err := toTest.MustGetParentDirectory().GetFileInDirectory(
					toTest.MustGetBaseName() + ".asc",
				)
				require.NoError(t, err)
				defer signatureFile.Delete(verbose)

				require.True(t, toTest.MustExists(verbose))
				require.False(t, signatureFile.MustExists(verbose))

				MustSignFile(
					toTest,
					&GnuPGSignOptions{
						DetachedSign: true,
						AsciiArmor:   tt.asciiArmor,
						Verbose:      verbose,
					},
				)

				require.True(t, toTest.MustExists(verbose))
				require.True(t, signatureFile.MustExists(verbose))

				MustCheckSignatureValid(signatureFile, verbose)
			},
		)
	}
}
