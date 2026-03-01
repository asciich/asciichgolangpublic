package containerimagehandler_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerimagehandler"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

// This example shows how to add or delete files from a container image archive.
func Test_Example_AddAndDeleteFilesToImageArchive(t *testing.T) {
	// Enable verbose output:
	ctx := contextutils.ContextVerbose()

	// Download the archive
	archivePath, err := containerimagehandler.DownloadImageAsTeporaryArchive(ctx, "alpine")
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, archivePath, &filesoptions.DeleteOptions{})

	// In the default ubuntu our example file does not exist:
	const examplePath = "/example.txt"
	exists, err := containerimagehandler.FileInArchiveExists(ctx, archivePath, examplePath)
	require.NoError(t, err)
	require.False(t, exists)

	fileList, err := containerimagehandler.ListFilesInArchive(ctx, archivePath)
	require.NoError(t, err)
	require.NotContains(t, fileList, examplePath)

	// Let's add the /example.txt using a temporary file as source:
	srcFile, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "This is the example content to add.")
	require.NoError(t, err)

	err = containerimagehandler.AddFileToArchive(ctx, archivePath, &containeroptions.AddFileToImageOptions{
		SourceFilePath:         srcFile,
		PathInImage:            examplePath,
		NewImageNameAndTag:     "exampleaddanddeletefiles:latest",
		OverwriteSourceArchive: true,
	})
	require.NoError(t, err)

	// Now own file exists in the image:
	exists, err = containerimagehandler.FileInArchiveExists(ctx, archivePath, examplePath)
	require.NoError(t, err)
	require.True(t, exists)

	fileList, err = containerimagehandler.ListFilesInArchive(ctx, archivePath)
	require.NoError(t, err)
	require.Contains(t, fileList, examplePath)

	// Also the content matches:
	content, err := containerimagehandler.ReadFileFromArchiveAsString(ctx, archivePath, examplePath)
	require.NoError(t, err)
	require.EqualValues(t, "This is the example content to add.", content)

	// Let's delete the file again:
	err = containerimagehandler.DeleteFileInArchive(ctx, archivePath, &containeroptions.DeleteFileFromImageOptions{
		PathInImage:            examplePath,
		NewImageNameAndTag:     "exampleaddanddeletefiles:latest",
		OverwriteSourceArchive: true,
	},
	)
	require.NoError(t, err)

	// Now own file exists in the image:
	exists, err = containerimagehandler.FileInArchiveExists(ctx, archivePath, examplePath)
	require.NoError(t, err)
	require.False(t, exists)

	fileList, err = containerimagehandler.ListFilesInArchive(ctx, archivePath)
	require.NoError(t, err)
	require.NotContains(t, fileList, examplePath)
}
