package hostreconciler_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/asciich/asciichgolangpublic/pkg/hostreconciler"
)

// createConfigDir creates a temporary config directory containing a host.yaml with the given content.
func createConfigDir(t *testing.T, hostYamlContent string) string {
	t.Helper()

	dirPath := t.TempDir()

	// TODO: Replace by the corresponding nativefiles function once verified.
	err := os.WriteFile(filepath.Join(dirPath, hostreconciler.HostConfigFileName), []byte(hostYamlContent), 0644)
	require.NoError(t, err)

	return dirPath
}

func TestDirectorySource_New(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		source, err := hostreconciler.NewDirectorySource("")
		require.Error(t, err)
		require.Nil(t, source)
	})

	t.Run("valid path", func(t *testing.T) {
		source, err := hostreconciler.NewDirectorySource("/tmp/config")
		require.NoError(t, err)
		require.NotNil(t, source)

		dirPath, err := source.GetDirPath()
		require.NoError(t, err)
		require.EqualValues(t, "/tmp/config", dirPath)
	})

	t.Run("implements Source interface", func(t *testing.T) {
		source, err := hostreconciler.NewDirectorySource("/tmp/config")
		require.NoError(t, err)

		var _ hostreconciler.Source = source
	})
}

func TestDirectorySource_GetHostConfigFilePath(t *testing.T) {
	t.Run("host.yaml inside config dir", func(t *testing.T) {
		source, err := hostreconciler.NewDirectorySource("/tmp/config")
		require.NoError(t, err)

		path, err := source.GetHostConfigFilePath()
		require.NoError(t, err)
		require.EqualValues(t, "/tmp/config/host.yaml", path)
	})
}

func TestDirectorySource_GetHostName(t *testing.T) {
	ctx := context.Background()

	t.Run("minimal config", func(t *testing.T) {
		dirPath := createConfigDir(t, "---\nhostname: host.example.com\n")

		source, err := hostreconciler.NewDirectorySource(dirPath)
		require.NoError(t, err)

		hostname, err := source.GetHostName(ctx)
		require.NoError(t, err)
		require.EqualValues(t, "host.example.com", hostname)
	})

	t.Run("minimal config without document separator", func(t *testing.T) {
		dirPath := createConfigDir(t, "hostname: host.example.com\n")

		source, err := hostreconciler.NewDirectorySource(dirPath)
		require.NoError(t, err)

		hostname, err := source.GetHostName(ctx)
		require.NoError(t, err)
		require.EqualValues(t, "host.example.com", hostname)
	})

	t.Run("config dir does not exist", func(t *testing.T) {
		source, err := hostreconciler.NewDirectorySource(filepath.Join(t.TempDir(), "does-not-exist"))
		require.NoError(t, err)

		hostname, err := source.GetHostName(ctx)
		require.Error(t, err)
		require.Empty(t, hostname)
	})

	t.Run("host.yaml missing", func(t *testing.T) {
		source, err := hostreconciler.NewDirectorySource(t.TempDir())
		require.NoError(t, err)

		hostname, err := source.GetHostName(ctx)
		require.Error(t, err)
		require.Empty(t, hostname)
	})

	t.Run("empty file", func(t *testing.T) {
		dirPath := createConfigDir(t, "")

		source, err := hostreconciler.NewDirectorySource(dirPath)
		require.NoError(t, err)

		hostname, err := source.GetHostName(ctx)
		require.Error(t, err)
		require.Empty(t, hostname)
	})

	t.Run("only document separator", func(t *testing.T) {
		dirPath := createConfigDir(t, "---\n")

		source, err := hostreconciler.NewDirectorySource(dirPath)
		require.NoError(t, err)

		hostname, err := source.GetHostName(ctx)
		require.Error(t, err)
		require.Empty(t, hostname)
	})

	t.Run("empty hostname", func(t *testing.T) {
		dirPath := createConfigDir(t, "---\nhostname: \"\"\n")

		source, err := hostreconciler.NewDirectorySource(dirPath)
		require.NoError(t, err)

		hostname, err := source.GetHostName(ctx)
		require.Error(t, err)
		require.Empty(t, hostname)
	})

	t.Run("unknown field is rejected", func(t *testing.T) {
		dirPath := createConfigDir(t, "---\nhostnmae: host.example.com\n")

		source, err := hostreconciler.NewDirectorySource(dirPath)
		require.NoError(t, err)

		hostname, err := source.GetHostName(ctx)
		require.Error(t, err)
		require.Empty(t, hostname)
	})

	t.Run("invalid yaml", func(t *testing.T) {
		dirPath := createConfigDir(t, "---\nhostname: [unclosed\n")

		source, err := hostreconciler.NewDirectorySource(dirPath)
		require.NoError(t, err)

		hostname, err := source.GetHostName(ctx)
		require.Error(t, err)
		require.Empty(t, hostname)
	})

	t.Run("multiple documents are rejected", func(t *testing.T) {
		dirPath := createConfigDir(t, "---\nhostname: a.example.com\n---\nhostname: b.example.com\n")

		source, err := hostreconciler.NewDirectorySource(dirPath)
		require.NoError(t, err)

		hostname, err := source.GetHostName(ctx)
		require.Error(t, err)
		require.Empty(t, hostname)
	})
}
