package dockerutils_test

import (
	"context"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// The convenience functions BackupVolume and RestoreVolume of the parent
// package delegate to the native implementation.
func Test_BackupAndRestoreVolume_ConvenienceFunctions(t *testing.T) {
	ctx := getCtx()
	docker := getDockerImplementationByName("nativeDocker")

	const volumeName = "test-convenience-volume"
	const testContent = "Hello from convenience volume!\n"
	const testFilename = "testfile.txt"

	// Ensure the volume is absent before the test and cleaned up afterwards:
	err := runCommandToRemoveVolume(ctx, volumeName)
	require.NoError(t, err)
	defer func() {
		err := runCommandToRemoveVolume(ctx, volumeName)
		if err != nil {
			logging.LogErrorByCtxf(ctx, "Failed to cleanup volume '%s': %v", volumeName, err)
		}
	}()

	// Create the volume and write a test file into it:
	err = runCommandToCreateVolume(ctx, volumeName)
	require.NoError(t, err)

	err = writeContentToVolume(ctx, docker, volumeName, testFilename, testContent)
	require.NoError(t, err)

	// UID, GID and access rights have to be preserved by BackupVolume and RestoreVolume:
	const testOwner = "1234:5678"
	const testMode = "741"
	err = setOwnerAndModeInVolume(ctx, volumeName, testFilename, testOwner, testMode)
	require.NoError(t, err)

	ownerAndMode, err := getOwnerAndModeInVolume(ctx, volumeName, testFilename)
	require.NoError(t, err)
	require.EqualValues(t, testOwner+":"+testMode, ownerAndMode)

	// Backup the volume through the convenience function:
	archivePath := filepath.Join(t.TempDir(), "backup.tar.gz")
	err = dockerutils.BackupVolume(ctx, volumeName, archivePath)
	require.NoError(t, err)
	require.FileExists(t, archivePath)

	// Remove the volume and restore it from the archive:
	err = runCommandToRemoveVolume(ctx, volumeName)
	require.NoError(t, err)

	err = dockerutils.RestoreVolume(ctx, volumeName, archivePath)
	require.NoError(t, err)

	// The content has to be restored:
	content, err := readContentFromVolume(ctx, docker, volumeName, testFilename)
	require.NoError(t, err)
	require.EqualValues(t, testContent, content, "Content should be restored from backup")

	// The UID, GID and access rights have to be restored as well:
	ownerAndMode, err = getOwnerAndModeInVolume(ctx, volumeName, testFilename)
	require.NoError(t, err)
	require.EqualValues(t, testOwner+":"+testMode, ownerAndMode, "UID, GID and access rights should be restored from backup")
}

func Test_ListVolumesNames(t *testing.T) {
	ctx := getCtx()

	const volumeName = "test-list-volumes-names-volume"

	t.Run("returns alphabetically sorted names", func(t *testing.T) {
		volumeNames, err := dockerutils.ListVolumesNames(ctx)
		require.NoError(t, err)
		require.NotNil(t, volumeNames)

		// The output has to be sorted to be deterministic:
		require.True(t, sort.StringsAreSorted(volumeNames))
	})

	t.Run("contains a created volume", func(t *testing.T) {
		// Ensure the volume is absent before the test and cleaned up afterwards:
		err := runCommandToRemoveVolume(ctx, volumeName)
		require.NoError(t, err)
		defer func() {
			err := runCommandToRemoveVolume(ctx, volumeName)
			if err != nil {
				logging.LogErrorByCtxf(ctx, "Failed to cleanup volume '%s': %v", volumeName, err)
			}
		}()

		volumeNames, err := dockerutils.ListVolumesNames(ctx)
		require.NoError(t, err)
		require.NotContains(t, volumeNames, volumeName)

		err = runCommandToCreateVolume(ctx, volumeName)
		require.NoError(t, err)

		volumeNames, err = dockerutils.ListVolumesNames(ctx)
		require.NoError(t, err)
		require.Contains(t, volumeNames, volumeName)
	})
}

// The convenience function ListVolumeFileNames of the parent package delegates
// to the native implementation.
func Test_ListVolumeFileNames_ConvenienceFunction(t *testing.T) {
	ctx := getCtx()
	docker := getDockerImplementationByName("nativeDocker")

	const volumeName = "test-convenience-list-volume-files-volume"
	const testContent = "Hello from convenience volume!\n"

	// Ensure the volume is absent before the test and cleaned up afterwards:
	err := runCommandToRemoveVolume(ctx, volumeName)
	require.NoError(t, err)
	defer func() {
		err := runCommandToRemoveVolume(ctx, volumeName)
		if err != nil {
			logging.LogErrorByCtxf(ctx, "Failed to cleanup volume '%s': %v", volumeName, err)
		}
	}()

	// Create the volume and write some files into it, including a hidden one:
	err = runCommandToCreateVolume(ctx, volumeName)
	require.NoError(t, err)

	err = writeContentToVolume(ctx, docker, volumeName, "testfile.txt", testContent)
	require.NoError(t, err)

	err = writeContentToVolume(ctx, docker, volumeName, ".hidden", testContent)
	require.NoError(t, err)

	// List all files inside the volume through the convenience function:
	fileNames, err := dockerutils.ListVolumeFileNames(ctx, volumeName)
	require.NoError(t, err)

	// The output has to be sorted to be deterministic:
	require.True(t, sort.StringsAreSorted(fileNames))

	// Hidden files have to be listed as well:
	require.EqualValues(t, []string{".hidden", "testfile.txt"}, fileNames)
}

func setOwnerAndModeInVolume(ctx context.Context, volumeName string, filename string, owner string, mode string) error {
	cmd := []string{
		"docker", "run", "--rm",
		"-v", volumeName + ":/volume",
		"alpine:latest",
		"/bin/sh", "-c", `chown "$1" "/volume/$2" && chmod "$3" "/volume/$2"`,
		"sh", owner, filename, mode, // $0, $1, $2, $3
	}

	executor := commandexecutorexecoo.Exec()
	_, err := executor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: cmd,
	})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to set owner '%s' and mode '%s' for '%s' in volume '%s': %w", owner, mode, filename, volumeName, err)
	}

	return nil
}

func getOwnerAndModeInVolume(ctx context.Context, volumeName string, filename string) (string, error) {
	cmd := []string{
		"docker", "run", "--rm",
		"-v", volumeName + ":/volume",
		"alpine:latest",
		"stat", "-c", "%u:%g:%a", "/volume/" + filename,
	}

	executor := commandexecutorexecoo.Exec()
	stdout, err := executor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: cmd,
	})
	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to read owner and mode of '%s' in volume '%s': %w", filename, volumeName, err)
	}

	return strings.TrimSpace(stdout), nil
}
