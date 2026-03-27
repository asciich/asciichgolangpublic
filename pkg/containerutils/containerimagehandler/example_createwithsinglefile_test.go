package containerimagehandler_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerimagehandler"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/pointerutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

// This example shows how an image archive wiht only one file can be created.
//
// The given file is packed in a way the resulting image:
//   - consists of 1 layer.
//   - containing only the added file.
//
// This method can be used to ship a static linked binary with no runtime dependencies as a single file container.
func Test_Example_CreateWithSingleFile(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.ContextVerbose()

	// Create an example file
	tempFilePath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, "example context\nwith a new line.\n")
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, tempFilePath, &filesoptions.DeleteOptions{})

	// Create a temporary file to store the ouput
	outDir, err := tempfiles.CreateTempDir(ctx)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, outDir, &filesoptions.DeleteOptions{})
	archivePath := filepath.Join(outDir, "example_latest.tar")

	// Create the container image archive:
	err = containerimagehandler.CreateSingleFileArchive(
		ctx,
		archivePath,
		&containeroptions.CreateSingleFileArchiveOptions{
			SourceFilePath:     tempFilePath,
			PathInImage:        "/testfile.txt",
			NewImageNameAndTag: "example:latest",
			Mode:               pointerutils.ToInt64Pointer(0644),
			Architecture:       "amd64",
		},
	)
	require.NoError(t, err)

	// There is only one file in the whole archive:
	fileNames, err := containerimagehandler.ListFilesInArchive(ctx, archivePath)
	require.NoError(t, err)
	require.EqualValues(t, []string{"/testfile.txt"}, fileNames)

	// Check the content as well.
	content, err := containerimagehandler.ReadFileFromArchiveAsString(ctx, archivePath, "/testfile.txt")
	require.NoError(t, err)
	require.EqualValues(t, "example context\nwith a new line.\n", content)
}
