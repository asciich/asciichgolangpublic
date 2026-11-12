package dockergeneric_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockergeneric"
)

func Test_GetMultiVolumeMountPoint(t *testing.T) {
	t.Run("mount point is below the volumes directory of the archive root", func(t *testing.T) {
		mountPoint, err := dockergeneric.GetMultiVolumeMountPoint("my-volume")
		require.NoError(t, err)
		require.Equal(t, "/archive/volumes/my-volume", mountPoint)
	})

	t.Run("empty volume name returns an error", func(t *testing.T) {
		_, err := dockergeneric.GetMultiVolumeMountPoint("")
		require.Error(t, err)
	})

	t.Run("volume name with path traversal returns an error", func(t *testing.T) {
		_, err := dockergeneric.GetMultiVolumeMountPoint("../etc")
		require.Error(t, err)
	})
}

func Test_SortedUniqueVolumeNames(t *testing.T) {
	t.Run("names are sorted", func(t *testing.T) {
		names, err := dockergeneric.SortedUniqueVolumeNames([]string{"c", "a", "b"})
		require.NoError(t, err)
		require.Equal(t, []string{"a", "b", "c"}, names)
		require.True(t, sort.StringsAreSorted(names))
	})

	t.Run("duplicates are removed", func(t *testing.T) {
		names, err := dockergeneric.SortedUniqueVolumeNames([]string{"b", "a", "b", "a"})
		require.NoError(t, err)
		require.Equal(t, []string{"a", "b"}, names)
	})

	t.Run("input slice is not modified", func(t *testing.T) {
		input := []string{"b", "a"}
		_, err := dockergeneric.SortedUniqueVolumeNames(input)
		require.NoError(t, err)
		require.Equal(t, []string{"b", "a"}, input)
	})

	t.Run("empty list returns an error", func(t *testing.T) {
		_, err := dockergeneric.SortedUniqueVolumeNames([]string{})
		require.Error(t, err)
	})

	t.Run("nil list returns an error", func(t *testing.T) {
		_, err := dockergeneric.SortedUniqueVolumeNames(nil)
		require.Error(t, err)
	})

	t.Run("invalid name returns an error", func(t *testing.T) {
		_, err := dockergeneric.SortedUniqueVolumeNames([]string{"a", "b/c"})
		require.Error(t, err)
	})
}

func Test_ParseManifestOutput(t *testing.T) {
	t.Run("valid manifest is parsed", func(t *testing.T) {
		manifest, err := dockergeneric.ParseManifestOutput(`{"version":1,"volumes":["a","b"]}`)
		require.NoError(t, err)
		require.Equal(t, []string{"a", "b"}, manifest.Volumes)
	})

	t.Run("surrounding whitespace and newlines are ignored", func(t *testing.T) {
		manifest, err := dockergeneric.ParseManifestOutput("\r\n  {\"version\":1,\"volumes\":[\"a\"]}\r\n")
		require.NoError(t, err)
		require.Equal(t, []string{"a"}, manifest.Volumes)
	})

	t.Run("no manifest marker returns ErrNoManifest", func(t *testing.T) {
		_, err := dockergeneric.ParseManifestOutput(dockergeneric.NoManifestMarker + "\n")
		require.ErrorIs(t, err, dockergeneric.ErrNoManifest)
	})

	t.Run("garbage returns an error which is not ErrNoManifest", func(t *testing.T) {
		_, err := dockergeneric.ParseManifestOutput("this is not json")
		require.Error(t, err)
		require.NotErrorIs(t, err, dockergeneric.ErrNoManifest)
	})

	t.Run("empty output returns an error which is not ErrNoManifest", func(t *testing.T) {
		_, err := dockergeneric.ParseManifestOutput("")
		require.Error(t, err)
		require.NotErrorIs(t, err, dockergeneric.ErrNoManifest)
	})
}

func Test_MultiVolumeScripts(t *testing.T) {
	t.Run("backup script writes manifest and archive into the mounted directories", func(t *testing.T) {
		require.Contains(t, dockergeneric.BackupVolumesScript, "/archive/manifest.json")
		require.Contains(t, dockergeneric.BackupVolumesScript, "tar czf")
		require.Contains(t, dockergeneric.BackupVolumesScript, "-C /archive")
		require.Contains(t, dockergeneric.BackupVolumesScript, `mv "/backup/$1" "/backup/$3"`)
	})

	t.Run("backup script removes the temporary file on failure", func(t *testing.T) {
		require.Contains(t, dockergeneric.BackupVolumesScript, "trap")
		require.Contains(t, dockergeneric.BackupVolumesScript, "set -e")
	})

	t.Run("restore script extracts into the archive root", func(t *testing.T) {
		require.Contains(t, dockergeneric.RestoreVolumesScript, "tar xzf")
		require.Contains(t, dockergeneric.RestoreVolumesScript, "-C /archive")
	})

	t.Run("read manifest script prints the no manifest marker", func(t *testing.T) {
		require.Contains(t, dockergeneric.ReadManifestScript, dockergeneric.NoManifestMarker)
		require.Contains(t, dockergeneric.ReadManifestScript, "manifest")
	})
}
