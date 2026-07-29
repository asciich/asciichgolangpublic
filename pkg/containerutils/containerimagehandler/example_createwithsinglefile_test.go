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

// This example shows how an image archive with only one file can be created.
//
// The given file is packed in a way the resulting image:
//   - consists of 1 layer.
//   - containing only the added file.
//
// This method can be used to ship a static linked binary with no runtime dependencies as a single file container.
func Test_Example_CreateWithSingleFile(t *testing.T) {
	// Enable verbose output
	ctx := contextutils.ContextVerbose()

	// Use a real statically linked binary (Go compiler) for the test
	// Text files or dynamically linked binaries will fail the static link check
	srcBinaryPath := "/usr/lib/go-1.24/bin/go"
	
	// Verify the binary exists and is statically linked before proceeding
	isStaticallyLinked, err := nativefiles.IsStaticallyLinkedBinary(ctx, srcBinaryPath)
	require.NoError(t, err)
	require.True(t, isStaticallyLinked, "Source binary must be statically linked")

	// Create a temporary file to store the output
	outDir, err := tempfiles.CreateTempDir(ctx)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, outDir, &filesoptions.DeleteOptions{})
	archivePath := filepath.Join(outDir, "example_latest.tar")

	// Create the container image archive:
	err = containerimagehandler.CreateSingleFileArchive(
		ctx,
		archivePath,
		&containeroptions.CreateSingleFileArchiveOptions{
			SourceFilePath:     srcBinaryPath,
			PathInImage:        "/go",
			NewImageNameAndTag: "example:latest",
			Mode:               pointerutils.ToInt64Pointer(0755),
			Architecture:       "amd64",
		},
	)
	require.NoError(t, err)

	// There is only one file in the whole archive:
	fileNames, err := containerimagehandler.ListFilesInArchive(ctx, archivePath)
	require.NoError(t, err)
	require.EqualValues(t, []string{"/go"}, fileNames)

	// Verify we can read the binary from the archive
	content, err := containerimagehandler.ReadFileFromArchiveAsBytes(ctx, archivePath, "/go")
	require.NoError(t, err)
	require.Greater(t, len(content), 0, "Binary content should not be empty")
}
