package containerimagehandler_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerimagehandler"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func TestListImageAndTagsInArchvie(t *testing.T) {
	t.Run("alpine", func(t *testing.T) {
		ctx := getCtx()

		archivePath, err := containerimagehandler.DownloadImageAsTeporaryArchive(ctx, "alpine")
		require.NoError(t, err)
		defer nativefiles.Delete(ctx, archivePath, &filesoptions.DeleteOptions{})

		imageNamesAndTags, err := containerimagehandler.ListImageNamesAndTagsInArchive(ctx, archivePath)
		require.NoError(t, err)
		require.EqualValues(t, []string{"alpine:latest"}, imageNamesAndTags)
	})
}
