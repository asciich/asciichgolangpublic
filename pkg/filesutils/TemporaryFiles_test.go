package filesutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

// This test suite ensure the different implementations behave in the same way.
func getFileToTest(implementationName string) (fileToTest filesinterfaces.File) {
	ctxSilent := contextutils.WithSilent(getCtx())

	temporayFile := mustutils.Must(tempfilesoo.CreateEmptyTemporaryFileAndGetPath(ctxSilent))

	if implementationName == "localFile" {
		return mustutils.Must(files.GetLocalFileByPath(temporayFile))
	}

	if implementationName == "localCommandExecutorFile" {
		return files.MustGetLocalCommandExecutorFileByPath(temporayFile)
	}

	logging.LogFatalWithTracef("Unknown implementation name '%s'", implementationName)
	return nil
}

func TestTemporaryFilesCreateFromFile(t *testing.T) {
	tests := []struct {
		implementationName string
		content            string
	}{
		{"localFile", "testcase"},
		{"localCommandExecutorFile", "testcase"},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true
				ctx := getCtx()

				sourceFile := getFileToTest(tt.implementationName)
				err := sourceFile.WriteString(tt.content, verbose)
				require.NoError(t, err)
				defer sourceFile.Delete(verbose)

				require.EqualValues(t, tt.content, sourceFile.MustReadAsString())

				tempFile, err := tempfilesoo.CreateTemporaryFileFromFile(ctx, sourceFile)
				require.NoError(t, err)
				defer tempFile.Delete(verbose)

				require.EqualValues(t, tt.content, tempFile.MustReadAsString())
			},
		)
	}
}
