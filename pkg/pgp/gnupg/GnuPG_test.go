package gnupg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

// Return a temporary file of the given 'implementationName'.
//
// Use defer file.Delete(verbose) to after calling this function to ensure
// the file is deleted after the test is over.
func getFileToTest(implementationName string) (file filesinterfaces.File) {
	ctx := getCtx()

	if implementationName == "localFile" {
		var err error
		file, err = files.GetLocalFileByPath(
			mustutils.Must(tempfilesoo.CreateEmptyTemporaryFileAndGetPath(ctx)),
		)
		if err != nil {
			panic(err)
		}
	} else if implementationName == "localCommandExecutorFile" {
		file = files.MustGetLocalCommandExecutorFileByPath(
			mustutils.Must(tempfilesoo.CreateEmptyTemporaryFileAndGetPath(ctx)),
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
				ctx := getCtx()

				toTest := getFileToTest(tt.implementationName)
				defer toTest.Delete(ctx, &filesoptions.DeleteOptions{})

				signatureFile, err := mustutils.Must(toTest.GetParentDirectory()).GetFileInDirectory(
					mustutils.Must(toTest.GetBaseName()) + ".asc",
				)
				require.NoError(t, err)
				defer signatureFile.Delete(ctx, &filesoptions.DeleteOptions{})

				require.True(t, mustutils.Must(toTest.Exists(ctx)))
				require.False(t, mustutils.Must(signatureFile.Exists(ctx)))

				MustSignFile(
					toTest,
					&GnuPGSignOptions{
						DetachedSign: true,
						AsciiArmor:   tt.asciiArmor,
						Verbose:      verbose,
					},
				)

				require.True(t, mustutils.Must(toTest.Exists(ctx)))
				require.True(t, mustutils.Must(signatureFile.Exists(ctx)))

				MustCheckSignatureValid(signatureFile, verbose)
			},
		)
	}
}
