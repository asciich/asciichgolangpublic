package filesutils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefilesoo"
)

// A Simple example how to get a local file as objectoriented file
func Test_ExampleNewFileByPath(t *testing.T) {
	// Enable verbose output:
	ctx := contextutils.ContextVerbose()

	// Create a temporary file using plain golang:
	tmpFile, err := os.CreateTemp("", "demo-*.txt")
	require.NoError(t, err)
	filePath := tmpFile.Name()
	defer nativefiles.Delete(ctx, filePath, &filesoptions.DeleteOptions{})

	// Get the local file:
	file, err := nativefilesoo.NewFileByPath(filePath)
	require.NoError(t,err)

	// Now we can work with this file in an object oriented way.
	// We read the size for example:
	size, err := file.GetSizeBytes(ctx)
	require.NoError(t, err)
	require.EqualValues(t, int64(0), size)
}
