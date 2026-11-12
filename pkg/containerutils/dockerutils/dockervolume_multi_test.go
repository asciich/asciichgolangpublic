package dockerutils_test

import (
	"compress/gzip"
	"context"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

const multiTestFilename = "testfile.txt"

// multiTestContent returns the content written into the volume with the given name.
func multiTestContent(volumeName string) string {
	return "Content of " + volumeName + "\n"
}

// ensureAbsentVolumes makes sure the volumes do not exist now and are removed again after the test.
func ensureAbsentVolumes(t *testing.T, ctx context.Context, volumeNames ...string) {
	t.Helper()

	for _, name := range volumeNames {
		require.NoError(t, runCommandToRemoveVolume(ctx, name))
	}

	t.Cleanup(func() {
		for _, name := range volumeNames {
			require.NoError(t, runCommandToRemoveVolume(ctx, name))
		}
	})
}

// ensureVolumesWithContent creates the volumes, each with a testfile containing multiTestContent.
func ensureVolumesWithContent(t *testing.T, ctx context.Context, docker dockerinterfaces.Docker, volumeNames ...string) {
	t.Helper()

	ensureAbsentVolumes(t, ctx, volumeNames...)

	for _, name := range volumeNames {
		require.NoError(t, runCommandToCreateVolume(ctx, name))
		require.NoError(t, writeContentToVolume(ctx, docker, name, multiTestFilename, multiTestContent(name)))
	}
}

func requireVolumeContent(t *testing.T, ctx context.Context, docker dockerinterfaces.Docker, volumeNames ...string) {
	t.Helper()

	for _, name := range volumeNames {
		content, err := readContentFromVolume(ctx, docker, name, multiTestFilename)
		require.NoError(t, err)
		require.EqualValues(t, multiTestContent(name), content)
	}
}

func readManifestOfArchive(t *testing.T, archivePath string) *dockergeneric.Manifest {
	t.Helper()

	f, err := os.Open(archivePath)
	require.NoError(t, err)
	defer f.Close()

	gz, err := gzip.NewReader(f)
	require.NoError(t, err)
	defer gz.Close()

	manifest, err := dockergeneric.ReadManifest(gz)
	require.NoError(t, err)

	return manifest
}

func Test_BackupAndRestoreVolumes(t *testing.T) {
	tests := []struct {
		implementationName string
	}{
		{"commandExecutorDocker"},
		{"nativeDocker"},
	}

	const volumeA = "test-multi-backup-volume-a"
	const volumeB = "test-multi-backup-volume-b"

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				docker := getDockerImplementationByName(tt.implementationName)

				t.Run("round trip keeps the content of all volumes", func(t *testing.T) {
					ensureVolumesWithContent(t, ctx, docker, volumeA, volumeB)
					archivePath := filepath.Join(t.TempDir(), "multi.tar.gz")

					require.NoError(t, docker.BackupVolumes(ctx, []string{volumeA, volumeB}, archivePath))
					require.FileExists(t, archivePath)

					require.NoError(t, runCommandToRemoveVolume(ctx, volumeA))
					require.NoError(t, runCommandToRemoveVolume(ctx, volumeB))

					volumes, err := docker.ListVolumes(ctx)
					require.NoError(t, err)
					require.NotContains(t, volumes, volumeA)
					require.NotContains(t, volumes, volumeB)

					require.NoError(t, docker.RestoreVolumes(ctx, archivePath))

					volumes, err = docker.ListVolumes(ctx)
					require.NoError(t, err)
					require.Contains(t, volumes, volumeA)
					require.Contains(t, volumes, volumeB)

					requireVolumeContent(t, ctx, docker, volumeA, volumeB)
				})

				t.Run("manifest lists the sorted volume names without duplicates", func(t *testing.T) {
					ensureVolumesWithContent(t, ctx, docker, volumeA, volumeB)
					archivePath := filepath.Join(t.TempDir(), "multi.tar.gz")

					require.NoError(t, docker.BackupVolumes(ctx, []string{volumeB, volumeA, volumeB}, archivePath))

					manifest := readManifestOfArchive(t, archivePath)
					require.EqualValues(t, []string{volumeA, volumeB}, manifest.Volumes)
					require.True(t, sort.StringsAreSorted(manifest.Volumes))
				})

				t.Run("backup of a non existing volume fails and leaves no archive and no volume behind", func(t *testing.T) {
					const missingVolume = "test-multi-backup-missing-volume"
					ensureVolumesWithContent(t, ctx, docker, volumeA)
					ensureAbsentVolumes(t, ctx, missingVolume)
					archivePath := filepath.Join(t.TempDir(), "multi.tar.gz")

					err := docker.BackupVolumes(ctx, []string{volumeA, missingVolume}, archivePath)
					require.Error(t, err)
					require.NoFileExists(t, archivePath)

					volumes, err := docker.ListVolumes(ctx)
					require.NoError(t, err)
					require.NotContains(t, volumes, missingVolume)
				})

				t.Run("no volume names returns an error", func(t *testing.T) {
					archivePath := filepath.Join(t.TempDir(), "multi.tar.gz")

					require.Error(t, docker.BackupVolumes(ctx, []string{}, archivePath))
					require.NoFileExists(t, archivePath)
				})

				t.Run("invalid volume name returns an error", func(t *testing.T) {
					archivePath := filepath.Join(t.TempDir(), "multi.tar.gz")

					require.Error(t, docker.BackupVolumes(ctx, []string{"../etc"}, archivePath))
					require.NoFileExists(t, archivePath)
				})

				t.Run("restore of a single volume archive fails with ErrNoManifest", func(t *testing.T) {
					const singleVolume = "test-multi-backup-single-volume"
					ensureVolumesWithContent(t, ctx, docker, singleVolume)
					archivePath := filepath.Join(t.TempDir(), "single.tar.gz")

					// The archive of the single volume backup has no manifest.
					require.NoError(t, docker.BackupVolume(ctx, singleVolume, archivePath))
					require.NoError(t, runCommandToRemoveVolume(ctx, singleVolume))

					err := docker.RestoreVolumes(ctx, archivePath)
					require.ErrorIs(t, err, dockergeneric.ErrNoManifest)

					// Nothing was restored.
					volumes, err := docker.ListVolumes(ctx)
					require.NoError(t, err)
					require.NotContains(t, volumes, singleVolume)
				})

				t.Run("restore of a non existing archive returns an error", func(t *testing.T) {
					archivePath := filepath.Join(t.TempDir(), "does-not-exist.tar.gz")

					require.Error(t, docker.RestoreVolumes(ctx, archivePath))
				})

				t.Run("empty archive path returns an error", func(t *testing.T) {
					require.Error(t, docker.BackupVolumes(ctx, []string{volumeA}, ""))
					require.Error(t, docker.RestoreVolumes(ctx, ""))
				})
			},
		)
	}
}

// The archive format must be identical for all implementations.
func Test_BackupAndRestoreVolumes_AcrossImplementations(t *testing.T) {
	tests := []struct {
		backupImplementationName  string
		restoreImplementationName string
	}{
		{"commandExecutorDocker", "nativeDocker"},
		{"nativeDocker", "commandExecutorDocker"},
	}

	const volumeA = "test-multi-cross-volume-a"
	const volumeB = "test-multi-cross-volume-b"

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				ctx := getCtx()
				backupDocker := getDockerImplementationByName(tt.backupImplementationName)
				restoreDocker := getDockerImplementationByName(tt.restoreImplementationName)

				ensureVolumesWithContent(t, ctx, backupDocker, volumeA, volumeB)
				archivePath := filepath.Join(t.TempDir(), "multi.tar.gz")

				require.NoError(t, backupDocker.BackupVolumes(ctx, []string{volumeA, volumeB}, archivePath))

				require.NoError(t, runCommandToRemoveVolume(ctx, volumeA))
				require.NoError(t, runCommandToRemoveVolume(ctx, volumeB))

				require.NoError(t, restoreDocker.RestoreVolumes(ctx, archivePath))

				requireVolumeContent(t, ctx, restoreDocker, volumeA, volumeB)
			},
		)
	}
}

// The convenience functions of the parent package delegate to the native implementation.
func Test_BackupAndRestoreVolumes_ConvenienceFunctions(t *testing.T) {
	ctx := getCtx()
	docker := getDockerImplementationByName("nativeDocker")

	const volumeA = "test-multi-convenience-volume-a"
	const volumeB = "test-multi-convenience-volume-b"

	t.Run("round trip keeps the content of all volumes", func(t *testing.T) {
		ensureVolumesWithContent(t, ctx, docker, volumeA, volumeB)
		archivePath := filepath.Join(t.TempDir(), "multi.tar.gz")

		require.NoError(t, dockerutils.BackupVolumes(ctx, []string{volumeA, volumeB}, archivePath))

		require.NoError(t, runCommandToRemoveVolume(ctx, volumeA))
		require.NoError(t, runCommandToRemoveVolume(ctx, volumeB))

		require.NoError(t, dockerutils.RestoreVolumes(ctx, archivePath))

		requireVolumeContent(t, ctx, docker, volumeA, volumeB)
	})
}
