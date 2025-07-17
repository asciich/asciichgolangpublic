package gnupg

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tempfiles"
	"github.com/asciich/asciichgolangpublic/testutils"
)

// Return a temporary file of the given 'implementationName'.
//
// Use defer file.Delete(verbose) to after calling this function to ensure
// the file is deleted after the test is over.
func getFileToTest(implementationName string) (file files.File) {
	const verbose = true

	if implementationName == "localFile" {
		file = files.MustGetLocalFileByPath(
			tempfiles.MustCreateEmptyTemporaryFileAndGetPath(verbose),
		)
	} else if implementationName == "localCommandExecutorFile" {
		file = files.MustGetLocalCommandExecutorFileByPath(
			tempfiles.MustCreateEmptyTemporaryFileAndGetPath(verbose),
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
				require := require.New(t)

				const verbose bool = true

				toTest := getFileToTest(tt.implementationName)
				defer toTest.Delete(verbose)

				signatureFile := toTest.MustGetParentDirectory().MustGetFileInDirectory(
					toTest.MustGetBaseName() + ".asc",
				)
				defer signatureFile.Delete(verbose)

				require.True(toTest.MustExists(verbose))
				require.False(signatureFile.MustExists(verbose))

				MustSignFile(
					toTest,
					&GnuPGSignOptions{
						DetachedSign: true,
						AsciiArmor:   tt.asciiArmor,
						Verbose:      verbose,
					},
				)

				require.True(toTest.MustExists(verbose))
				require.True(signatureFile.MustExists(verbose))

				MustCheckSignatureValid(signatureFile, verbose)
			},
		)
	}
}
