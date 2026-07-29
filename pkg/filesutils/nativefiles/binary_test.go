package nativefiles_test

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func Test_IsStaticallyLinkedBinary(t *testing.T) {
	ctx := getCtx()

	t.Run("empty string", func(t *testing.T) {
		isStaticallyLinked, err := nativefiles.IsStaticallyLinkedBinary(ctx, "")
		require.Error(t, err)
		require.False(t, isStaticallyLinked)
	})

	t.Run("nonexisting file", func(t *testing.T) {
		isStaticallyLinked, err := nativefiles.IsStaticallyLinkedBinary(ctx, "/this/file/does/not/exist")
		require.Error(t, err)
		require.False(t, isStaticallyLinked)
	})

	t.Run("directory", func(t *testing.T) {
		isStaticallyLinked, err := nativefiles.IsStaticallyLinkedBinary(ctx, "/etc")
		require.Error(t, err)
		require.False(t, isStaticallyLinked)
	})

	t.Run("text file", func(t *testing.T) {
		isStaticallyLinked, err := nativefiles.IsStaticallyLinkedBinary(ctx, "/etc/hosts")
		require.NoError(t, err)
		require.False(t, isStaticallyLinked)
	})

	t.Run("dynamically linked binary", func(t *testing.T) {
		// Most system binaries are dynamically linked
		isStaticallyLinked, err := nativefiles.IsStaticallyLinkedBinary(ctx, "/bin/ls")
		require.NoError(t, err)
		require.False(t, isStaticallyLinked)
	})

	t.Run("statically linked binary", func(t *testing.T) {
		// Create a simple statically linked binary for testing
		tempDir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, tempDir, &filesoptions.DeleteOptions{})

		binaryPath := filepath.Join(tempDir, "static_binary")
		sourcePath := filepath.Join(tempDir, "main.go")

		// Create a minimal Go program
		sourceCode := `package main
func main() {}
`
		err = nativefiles.WriteString(ctx, sourcePath, sourceCode)
		require.NoError(t, err)

		// Compile with static linking
		cmd := exec.Command("go", "build", "-ldflags", "-extldflags '-static'", "-o", binaryPath, sourcePath)
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
		output, err := cmd.CombinedOutput()
		if err != nil {
			// If static compilation fails, skip this test but log the output
			t.Logf("Skipping static binary test: failed to compile static binary: %v, output: %s", err, string(output))
			t.Skip("Cannot create statically linked binary for testing")
		}

		isStaticallyLinked, err := nativefiles.IsStaticallyLinkedBinary(ctx, binaryPath)
		require.NoError(t, err)
		require.True(t, isStaticallyLinked)
	})
}

func Test_IsStaticallyLinkedBinary_WithContextCancellation(t *testing.T) {
	t.Run("cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		isStaticallyLinked, err := nativefiles.IsStaticallyLinkedBinary(ctx, "/bin/ls")
		require.Error(t, err)
		require.False(t, isStaticallyLinked)
	})
}
