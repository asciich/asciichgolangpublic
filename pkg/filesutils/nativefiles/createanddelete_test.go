package nativefiles_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func Test_CreateAndDeleteFile(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		ctx := getCtx()

		// We use a temporary file path for testing:
		filePath, err := tempfiles.CreateTemporaryFile(ctx)
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, filePath)
		err = nativefiles.Delete(ctx, filePath)
		require.NoError(t, err)

		// Test Delete first to ensure the file is absent in any case for the next test steps.
		// Is done twice to test idempotence.
		for range 2 {
			err = nativefiles.Delete(ctx, filePath)
			require.NoError(t, err)

			require.False(t, nativefiles.Exists(ctx, filePath))
		}

		// Create the file again
		// Is done twice to test idempotence.
		for range 2 {
			err = nativefiles.Create(ctx, filePath)
			require.NoError(t, err)

			require.True(t, nativefiles.Exists(ctx, filePath))
		}

		// Test Delete.
		// Is done twice to test idempotence.
		for range 2 {
			err = nativefiles.Delete(ctx, filePath)
			require.NoError(t, err)

			require.False(t, nativefiles.Exists(ctx, filePath))
		}
	})
}
