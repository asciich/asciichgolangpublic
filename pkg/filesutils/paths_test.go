package filesutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetAbsolutePath_RelativePath(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	absPath, err := GetAbsolutePath("somefile.txt")
	require.NoError(t, err)
	require.Equal(t, filepath.Join(cwd, "somefile.txt"), absPath)
}

func TestGetAbsolutePath_AlreadyAbsolute(t *testing.T) {
	absPath, err := GetAbsolutePath("/tmp/somefile.txt")
	require.NoError(t, err)
	require.Equal(t, "/tmp/somefile.txt", absPath)
}

func TestGetAbsolutePath_DotPath(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	absPath, err := GetAbsolutePath(".")
	require.NoError(t, err)
	require.Equal(t, cwd, absPath)
}

func TestGetAbsolutePath_NestedRelativePath(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	absPath, err := GetAbsolutePath("subdir/file.go")
	require.NoError(t, err)
	require.Equal(t, filepath.Join(cwd, "subdir", "file.go"), absPath)
}

func TestGetAbsolutePath_ParentDirectory(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	absPath, err := GetAbsolutePath("../file.go")
	require.NoError(t, err)
	require.Equal(t, filepath.Join(filepath.Dir(cwd), "file.go"), absPath)
}

func TestGetAbsolutePath_EmptyString(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	absPath, err := GetAbsolutePath("")
	require.NoError(t, err)
	require.Equal(t, cwd, absPath)
}
