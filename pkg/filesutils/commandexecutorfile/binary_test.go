package commandexecutorfile_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
)

func TestIsStaticallyLinkedBinary(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()

		isStatic, err := commandexecutorfile.IsStaticallyLinkedBinary(ctx, nil, "/bin/ls")
		require.Error(t, err)
		require.False(t, isStatic)
	})

	t.Run("empty filePath returns error", func(t *testing.T) {
		ctx := getCtx()

		isStatic, err := commandexecutorfile.IsStaticallyLinkedBinary(ctx, commandexecutorexecoo.Exec(), "")
		require.Error(t, err)
		require.False(t, isStatic)
	})

	t.Run("/bin/ls is not statically linked", func(t *testing.T) {
		ctx := getCtx()

		isStatic, err := commandexecutorfile.IsStaticallyLinkedBinary(ctx, commandexecutorexecoo.Exec(), "/bin/ls")
		require.NoError(t, err)
		require.False(t, isStatic)
	})

	t.Run("/bin/sh is not statically linked", func(t *testing.T) {
		ctx := getCtx()

		isStatic, err := commandexecutorfile.IsStaticallyLinkedBinary(ctx, commandexecutorexecoo.Exec(), "/bin/sh")
		require.NoError(t, err)
		require.False(t, isStatic)
	})

	t.Run("nonexistent file returns false", func(t *testing.T) {
		ctx := getCtx()

		isStatic, err := commandexecutorfile.IsStaticallyLinkedBinary(ctx, commandexecutorexecoo.Exec(), "/nonexistent/file")
		// The file command doesn't return an error for nonexistent files, but the output will contain "cannot open"
		// In this case we expect no error but isStatic should be false
		require.False(t, isStatic)
		// Check that the output contains error message about file not existing
		if err != nil {
			require.True(t, strings.Contains(err.Error(), "cannot open") || strings.Contains(err.Error(), "No such file"))
		}
	})
}
