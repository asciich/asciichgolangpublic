package dockerutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
)

// Test_Example_PullContainerImage demonstrates how to pull a container image using PullContainerImage.
// This example shows pulling an image and verifying its basic properties.
func Test_Example_PullContainerImage(t *testing.T) {
	// Use a context with verbose output enabled to see progress:
	ctx := contextutils.ContextVerbose()

	// Pull a container image (e.g., alpine:latest):
	image, err := dockerutils.PullContainerImage(ctx, "alpine:latest")
	require.NoError(t, err)
	require.NotNil(t, image)

	// Get the image name to verify it was pulled correctly:
	name, err := image.GetName()
	require.NoError(t, err)
	require.Equal(t, "alpine:latest", name)

	// Check if the image exists in the local Docker daemon:
	exists, err := image.Exists(ctx)
	require.NoError(t, err)
	require.True(t, exists)
}

// Test_Example_PullContainerImage_specificVersion demonstrates pulling a specific version of an image.
// This is recommended for production use to ensure reproducible builds.
func Test_Example_PullContainerImage_specificVersion(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Pull a specific version of an image (e.g., ubuntu:22.04):
	image, err := dockerutils.PullContainerImage(ctx, "ubuntu:22.04")
	require.NoError(t, err)
	require.NotNil(t, image)

	// Verify the image name:
	name, err := image.GetName()
	require.NoError(t, err)
	require.Equal(t, "ubuntu:22.04", name)

	// Check that the image exists:
	exists, err := image.Exists(ctx)
	require.NoError(t, err)
	require.True(t, exists)
}

// Test_Example_PullContainerImage_withErrorHandling demonstrates proper error handling when pulling images.
func Test_Example_PullContainerImage_withErrorHandling(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Attempt to pull a non-existent image to demonstrate error handling:
	_, err := dockerutils.PullContainerImage(ctx, "nonexistent-image:invalid-tag")
	require.Error(t, err)
}

// Test_Example_PullContainerImage_verifyImageExists demonstrates verifying an image after pulling.
func Test_Example_PullContainerImage_verifyImageExists(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	// Pull an image:
	image, err := dockerutils.PullContainerImage(ctx, "busybox:latest")
	require.NoError(t, err)
	require.NotNil(t, image)

	// Get the image name:
	name, err := image.GetName()
	require.NoError(t, err)

	// Verify the image exists:
	exists, err := image.Exists(ctx)
	require.NoError(t, err)
	require.True(t, exists)

	// Verify the image has the expected name:
	require.Contains(t, name, "busybox")
}
