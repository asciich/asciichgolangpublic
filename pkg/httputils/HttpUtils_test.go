package httputils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/httputils"
)

func Test_WithDownloadProgress(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNBytes(nil, 0)
		require.EqualValues(t, 0, httputils.GetProgressEveryNBytes(ctx))
	})

	t.Run("negative value", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNBytes(getCtx(), -1)
		require.EqualValues(t, 0, httputils.GetProgressEveryNBytes(ctx))
	})

	t.Run("zero", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNBytes(getCtx(), 0)
		require.EqualValues(t, 0, httputils.GetProgressEveryNBytes(ctx))
	})

	t.Run("one", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNBytes(getCtx(), 1)
		require.EqualValues(t, 1, httputils.GetProgressEveryNBytes(ctx))
	})

	t.Run("one", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNBytes(getCtx(), 10)
		require.EqualValues(t, 10, httputils.GetProgressEveryNBytes(ctx))
	})
}

func Test_WithDownloadProgress_kBytes(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNkBytes(nil, 0)
		require.EqualValues(t, 0, httputils.GetProgressEveryNBytes(ctx))
	})

	t.Run("negative value", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNkBytes(getCtx(), -1)
		require.EqualValues(t, 0, httputils.GetProgressEveryNBytes(ctx))
	})

	t.Run("zero", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNkBytes(getCtx(), 0)
		require.EqualValues(t, 0, httputils.GetProgressEveryNBytes(ctx))
	})

	t.Run("one", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNkBytes(getCtx(), 1)
		require.EqualValues(t, 1*1024, httputils.GetProgressEveryNBytes(ctx))
	})

	t.Run("one", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNkBytes(getCtx(), 10)
		require.EqualValues(t, 10*1024, httputils.GetProgressEveryNBytes(ctx))
	})
}

func Test_WithDownloadProgress_MBytes(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNMBytes(nil, 0)
		require.EqualValues(t, 0, httputils.GetProgressEveryNBytes(ctx))
	})

	t.Run("negative value", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNMBytes(getCtx(), -1)
		require.EqualValues(t, 0, httputils.GetProgressEveryNBytes(ctx))
	})

	t.Run("zero", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNMBytes(getCtx(), 0)
		require.EqualValues(t, 0, httputils.GetProgressEveryNBytes(ctx))
	})

	t.Run("one", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNMBytes(getCtx(), 1)
		require.EqualValues(t, 1*1024*1024, httputils.GetProgressEveryNBytes(ctx))
	})

	t.Run("one", func(t *testing.T) {
		ctx := httputils.WithDownloadProgressEveryNMBytes(getCtx(), 10)
		require.EqualValues(t, 10*1024*1024, httputils.GetProgressEveryNBytes(ctx))
	})
}
