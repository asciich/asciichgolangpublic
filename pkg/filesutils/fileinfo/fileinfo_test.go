package fileinfo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/fileinfo"
)

func Test_NewFileInfo(t *testing.T) {
	t.Run("not nil", func(t *testing.T) {
		f := fileinfo.NewFileInfo()
		require.NotNil(t, f)
	})
}

func Test_GetPath(t *testing.T) {
	t.Run("empty path returns error", func(t *testing.T) {
		f := fileinfo.NewFileInfo()
		_, err := f.GetPath()
		require.Error(t, err)
	})

	t.Run("path set returns path", func(t *testing.T) {
		f := fileinfo.NewFileInfo()
		err := f.SetPath("/tmp/test.txt")
		require.NoError(t, err)

		path, err := f.GetPath()
		require.NoError(t, err)
		require.EqualValues(t, "/tmp/test.txt", path)
	})
}

func Test_SetPath(t *testing.T) {
	t.Run("empty string returns error", func(t *testing.T) {
		f := fileinfo.NewFileInfo()
		err := f.SetPath("")
		require.Error(t, err)
	})

	t.Run("non empty string sets path", func(t *testing.T) {
		f := fileinfo.NewFileInfo()
		err := f.SetPath("/tmp/abc")
		require.NoError(t, err)

		path, err := f.GetPath()
		require.NoError(t, err)
		require.EqualValues(t, "/tmp/abc", path)
	})
}

func Test_GetSizeBytes(t *testing.T) {
	t.Run("default value is zero", func(t *testing.T) {
		f := fileinfo.NewFileInfo()
		sizeBytes, err := f.GetSizeBytes()
		require.NoError(t, err)
		require.EqualValues(t, 0, sizeBytes)
	})

	t.Run("returns set value", func(t *testing.T) {
		f := fileinfo.NewFileInfo()
		err := f.SetSizeBytes(1024)
		require.NoError(t, err)

		sizeBytes, err := f.GetSizeBytes()
		require.NoError(t, err)
		require.EqualValues(t, 1024, sizeBytes)
	})

	t.Run("negative value", func(t *testing.T) {
		f := fileinfo.NewFileInfo()
		err := f.SetSizeBytes(-1)
		require.NoError(t, err)

		sizeBytes, err := f.GetSizeBytes()
		require.NoError(t, err)
		require.EqualValues(t, -1, sizeBytes)
	})
}

func Test_SetSizeBytes(t *testing.T) {
	t.Run("set zero", func(t *testing.T) {
		f := fileinfo.NewFileInfo()
		err := f.SetSizeBytes(0)
		require.NoError(t, err)

		sizeBytes, err := f.GetSizeBytes()
		require.NoError(t, err)
		require.EqualValues(t, 0, sizeBytes)
	})

	t.Run("set positive value", func(t *testing.T) {
		f := fileinfo.NewFileInfo()
		err := f.SetSizeBytes(4096)
		require.NoError(t, err)

		sizeBytes, err := f.GetSizeBytes()
		require.NoError(t, err)
		require.EqualValues(t, 4096, sizeBytes)
	})
}

func Test_GetPathAndSizeBytes(t *testing.T) {
	t.Run("path not set returns error", func(t *testing.T) {
		f := fileinfo.NewFileInfo()
		_, _, err := f.GetPathAndSizeBytes()
		require.Error(t, err)
	})

	t.Run("path and size set returns both", func(t *testing.T) {
		f := fileinfo.NewFileInfo()
		err := f.SetPath("/tmp/test.txt")
		require.NoError(t, err)

		err = f.SetSizeBytes(2048)
		require.NoError(t, err)

		path, sizeBytes, err := f.GetPathAndSizeBytes()
		require.NoError(t, err)
		require.EqualValues(t, "/tmp/test.txt", path)
		require.EqualValues(t, 2048, sizeBytes)
	})

	t.Run("path set size default", func(t *testing.T) {
		f := fileinfo.NewFileInfo()
		err := f.SetPath("/tmp/empty.txt")
		require.NoError(t, err)

		path, sizeBytes, err := f.GetPathAndSizeBytes()
		require.NoError(t, err)
		require.EqualValues(t, "/tmp/empty.txt", path)
		require.EqualValues(t, 0, sizeBytes)
	})
}
