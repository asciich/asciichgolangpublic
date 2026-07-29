package containerimagehandler_test

import (
	"os"
	"os/exec"
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

	// Create a temporary directory for our test files
	tempDir, err := tempfiles.CreateTempDir(ctx)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

	// Create a minimal Go source file
	sourcePath := filepath.Join(tempDir, "main.go")
	sourceCode := `package main
func main() {}
`
	err = nativefiles.WriteString(ctx, sourcePath, sourceCode)
	require.NoError(t, err)

	// Build a statically linked binary from the source
	binaryPath := filepath.Join(tempDir, "testbinary")
	cmd := exec.Command("go", "build", "-ldflags", "-extldflags '-static'", "-o", binaryPath, sourcePath)
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	_, err = cmd.CombinedOutput()
	require.NoError(t, err)

	// Verify the binary is statically linked
	isStaticallyLinked, err := nativefiles.IsStaticallyLinkedBinary(ctx, binaryPath)
	require.NoError(t, err)
	require.True(t, isStaticallyLinked, "Built binary must be statically linked")

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
			SourceFilePath:     binaryPath,
			PathInImage:        "/testbinary",
			NewImageNameAndTag: "example:latest",
			Mode:               pointerutils.ToInt64Pointer(0755),
			Architecture:       "amd64",
		},
	)
	require.NoError(t, err)

	// There is only one file in the whole archive:
	fileNames, err := containerimagehandler.ListFilesInArchive(ctx, archivePath)
	require.NoError(t, err)
	require.EqualValues(t, []string{"/testbinary"}, fileNames)

	// Verify we can read the binary from the archive
	content, err := containerimagehandler.ReadFileFromArchiveAsBytes(ctx, archivePath, "/testbinary")
	require.NoError(t, err)
	require.Greater(t, len(content), 0, "Binary content should not be empty")
}
