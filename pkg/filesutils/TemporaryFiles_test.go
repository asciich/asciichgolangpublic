package filesutils_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutortempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfilesoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

// This test suite ensure the different implementations behave in the same way.

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func getTemporaryFileToTest(implementationName string) (fileToTest filesinterfaces.File) {
	ctxSilent := contextutils.WithSilent(getCtx())

	temporayFilePath := mustutils.Must(tempfiles.CreateTemporaryFile(ctxSilent))

	return getFileToTest(implementationName, temporayFilePath)
}

func getFileToTest(implementationName string, path string) (fileToTest filesinterfaces.File) {
	if implementationName == "localFile" {
		return mustutils.Must(files.GetLocalFileByPath(path))
	}

	if implementationName == "localCommandExecutorFile" {
		return mustutils.Must(files.GetLocalCommandExecutorFileByPath(path))
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
				ctx := getCtx()

				sourceFile := getTemporaryFileToTest(tt.implementationName)
				err := sourceFile.WriteString(ctx, tt.content, &filesoptions.WriteOptions{})
				require.NoError(t, err)
				defer sourceFile.Delete(ctx, &filesoptions.DeleteOptions{})

				require.EqualValues(t, tt.content, sourceFile.MustReadAsString())

				tempFile, err := tempfilesoo.CreateTemporaryFileFromFile(ctx, sourceFile)
				require.NoError(t, err)
				defer tempFile.Delete(ctx, &filesoptions.DeleteOptions{})

				require.EqualValues(t, tt.content, tempFile.MustReadAsString())
			},
		)
	}
}

func TestCreateEmptyTemporaryFile(t *testing.T) {
	tests := []struct {
		name           string
		createFunction func(context.Context) (filesinterfaces.File, error)
	}{
		{"tempfilesoo", tempfilesoo.CreateEmptyTemporaryFile},
		{
			"commandexecutor",
			func(ctx context.Context) (filesinterfaces.File, error) {
				return commandexecutortempfilesoo.CreateEmptyTemporaryFile(ctx, commandexecutorexecoo.NewExec())
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				ctx := getCtx()

				file, err := tt.createFunction(ctx)
				require.NoError(t, err)

				defer func() {
					err := file.Delete(ctx, &filesoptions.DeleteOptions{})
					require.NoError(t, err)
				}()

				exists, err := file.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				content, err := file.ReadAsString()
				require.NoError(t, err)
				require.EqualValues(t, "", content)

				localPath, err := file.GetPath()
				require.NoError(t, err)
				require.True(t, strings.HasPrefix(localPath, "/tmp/"))
			},
		)
	}
}

func TestCreateEmptyTemporaryDirectory(t *testing.T) {
	tests := []struct {
		name           string
		createFunction func(context.Context) (filesinterfaces.Directory, error)
	}{
		{"tempfilesoo", tempfilesoo.CreateEmptyTemporaryDirectory},
		{
			"commandexecutor",
			func(ctx context.Context) (filesinterfaces.Directory, error) {
				return commandexecutortempfilesoo.CreateEmptyTemporaryDirectory(ctx, commandexecutorexecoo.NewExec())
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				ctx := getCtx()

				directory, err := tt.createFunction(ctx)
				require.NoError(t, err)

				defer func() {
					err := directory.Delete(ctx, &filesoptions.DeleteOptions{})
					require.NoError(t, err)
				}()

				exists, err := directory.Exists(ctx)
				require.NoError(t, err)
				require.True(t, exists)

				localPath, err := directory.GetPath()
				require.NoError(t, err)
				require.True(t, strings.HasPrefix(localPath, "/tmp/"))
			},
		)
	}
}
