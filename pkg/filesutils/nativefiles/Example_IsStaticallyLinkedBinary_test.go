package nativefiles_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

// This example shows how to use the nativefiles.IsStaticallyLinkedBinary function.
func Test_Example_IsStaticallyLinkedBinary(t *testing.T) {
	// Use a context with verbose output enabled:
	ctx := contextutils.ContextVerbose()

	// Check a dynamically linked system binary:
	isStaticallyLinked, err := nativefiles.IsStaticallyLinkedBinary(ctx, "/bin/ls")
	require.NoError(t, err)
	require.False(t, isStaticallyLinked, "/bin/ls should be dynamically linked")

	// Create a statically linked binary for testing:
	tempDir, err := tempfiles.CreateTempDir(ctx)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

	binaryPath := filepath.Join(tempDir, "static_binary")
	sourcePath := filepath.Join(tempDir, "main.go")

	// Create a minimal Go program:
	sourceCode := `package main
func main() {}
`
	err = nativefiles.WriteString(ctx, sourcePath, sourceCode)
	require.NoError(t, err)

	// Compile with static linking:
	cmd := exec.Command("go", "build", "-ldflags", "-extldflags '-static'", "-o", binaryPath, sourcePath)
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Skipping example: failed to compile static binary: %v, output: %s", err, string(output))
		t.Skip("Cannot create statically linked binary for testing")
	}

	// Check the statically linked binary:
	isStaticallyLinked, err = nativefiles.IsStaticallyLinkedBinary(ctx, binaryPath)
	require.NoError(t, err)
	require.True(t, isStaticallyLinked, "static_binary should be statically linked")
}
