package nativefiles_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
)

func TestListFiles(t *testing.T) {
	t.Run("empty dir allow empty list", func(t *testing.T) {
		ctx := getCtx()

		dir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)

		fileList, err := nativefiles.ListFiles(ctx, dir, &parameteroptions.ListFileOptions{
			AllowEmptyListIfNoFileIsFound: true,
		})
		require.NoError(t, err)
		require.Empty(t, fileList)
	})

	t.Run("empty dir", func(t *testing.T) {
		ctx := getCtx()

		dir, err := tempfiles.CreateTempDir(ctx)
		require.NoError(t, err)

		fileList, err := nativefiles.ListFiles(ctx, dir, &parameteroptions.ListFileOptions{
			AllowEmptyListIfNoFileIsFound: false,
		})
		require.Error(t, err)
		require.Nil(t, fileList)
	})
}
