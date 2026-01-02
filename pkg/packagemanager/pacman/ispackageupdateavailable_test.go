package pacman_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/pacman"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_IsPackageUpdateAvailable(t *testing.T) {
	ctx := getCtx()

	const packageName = "archlinux-keyring"

	// Let's take an fixed image version where we know the package is outdated:
	const imageName = "archlinux:base-20250727.0.390543"

	container, err := nativedocker.NewDocker().GetContainerByName("test-pacman-is-package-update-available")
	require.NoError(t, err)
	container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
	defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	err = container.Run(ctx, &dockeroptions.DockerRunContainerOptions{
		ImageName: imageName,
		Command:   []string{"sleep", "1m"},
	})
	require.NoError(t, err)

	pacman, err := pacman.NewPacman(container)
	require.NoError(t, err)

	// We expect an update available:
	updateAvailabe, err := pacman.IsPackageUpdateAvailalbe(ctx, packageName)
	require.NoError(t, err)
	require.True(t, updateAvailabe)

	// Let's update the package:
	err = pacman.UpdatePackage(ctx, packageName)
	require.NoError(t, err)

	// After the update we expect no update available:
	updateAvailabe, err = pacman.IsPackageUpdateAvailalbe(ctx, packageName)
	require.NoError(t, err)
	require.False(t, updateAvailabe)
}
