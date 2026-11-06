package dockerutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

// This example shows how to list all files inside a docker volume.
// A temporary container is started with the volume mounted, lists all files
// including the hidden ones and is deleted afterwards.
// The returned file names are sorted alphabetically to be deterministic.
func Test_Example_ListVolumeFileNames(t *testing.T) {
	// use a context with verbose output enabled:
	ctx := contextutils.ContextVerbose()

	// Get docker on local host
	docker, err := dockerutils.GetDockerOnLocalHost()
	require.NoError(t, err)

	// Define the name of our volume:
	const volumeName = "example-list-volume-files-volume"
	const testContent = "Hello from example volume!\n"

	// Ensure the volume is absent before we start (cleanup from previous failed runs):
	err = runCommandToRemoveVolume(ctx, volumeName)
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

	// List all files inside the volume:
	fileNames, err := dockerutils.ListVolumeFileNames(ctx, volumeName)
	require.NoError(t, err)

	// The names are sorted alphabetically and hidden files are included:
	require.EqualValues(t, []string{".hidden", "testfile.txt"}, fileNames)
}
