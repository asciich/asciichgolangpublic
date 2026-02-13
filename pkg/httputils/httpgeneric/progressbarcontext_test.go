package httpgeneric_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpgeneric"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_WithDownloadProgress(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNBytes(nil, 0)
		require.EqualValues(t, 0, httpgeneric.GetProgressEveryNBytes(ctx))
	})

	t.Run("negative value", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNBytes(getCtx(), -1)
		require.EqualValues(t, 0, httpgeneric.GetProgressEveryNBytes(ctx))
	})

	t.Run("zero", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNBytes(getCtx(), 0)
		require.EqualValues(t, 0, httpgeneric.GetProgressEveryNBytes(ctx))
	})

	t.Run("one", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNBytes(getCtx(), 1)
		require.EqualValues(t, 1, httpgeneric.GetProgressEveryNBytes(ctx))
	})

	t.Run("one", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNBytes(getCtx(), 10)
		require.EqualValues(t, 10, httpgeneric.GetProgressEveryNBytes(ctx))
	})
}

func Test_WithDownloadProgress_kBytes(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNkBytes(nil, 0)
		require.EqualValues(t, 0, httpgeneric.GetProgressEveryNBytes(ctx))
	})

	t.Run("negative value", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNkBytes(getCtx(), -1)
		require.EqualValues(t, 0, httpgeneric.GetProgressEveryNBytes(ctx))
	})

	t.Run("zero", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNkBytes(getCtx(), 0)
		require.EqualValues(t, 0, httpgeneric.GetProgressEveryNBytes(ctx))
	})

	t.Run("one", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNkBytes(getCtx(), 1)
		require.EqualValues(t, 1*1024, httpgeneric.GetProgressEveryNBytes(ctx))
	})

	t.Run("one", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNkBytes(getCtx(), 10)
		require.EqualValues(t, 10*1024, httpgeneric.GetProgressEveryNBytes(ctx))
	})
}

func Test_WithDownloadProgress_MBytes(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNMBytes(nil, 0)
		require.EqualValues(t, 0, httpgeneric.GetProgressEveryNBytes(ctx))
	})

	t.Run("negative value", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNMBytes(getCtx(), -1)
		require.EqualValues(t, 0, httpgeneric.GetProgressEveryNBytes(ctx))
	})

	t.Run("zero", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNMBytes(getCtx(), 0)
		require.EqualValues(t, 0, httpgeneric.GetProgressEveryNBytes(ctx))
	})

	t.Run("one", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNMBytes(getCtx(), 1)
		require.EqualValues(t, 1*1024*1024, httpgeneric.GetProgressEveryNBytes(ctx))
	})

	t.Run("one", func(t *testing.T) {
		ctx := httpgeneric.WithDownloadProgressEveryNMBytes(getCtx(), 10)
		require.EqualValues(t, 10*1024*1024, httpgeneric.GetProgressEveryNBytes(ctx))
	})
}
