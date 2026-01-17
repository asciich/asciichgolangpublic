package archlinuxutils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxutils/archlinuxutils"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/pacman"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_UpdateArchlinuxKeyringPackage(t *testing.T) {
	ctx := getCtx()

	// Let's take an fixed image version where we know the package is outdated:
	const imageName = "archlinux:base-20250727.0.390543"

	container, err := nativedocker.NewDocker().GetContainerByName("test-update-archlinux-keyring")
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

	// We expect an update of the "archlinux-keyring" package:
	const packageName = "archlinux-keyring"
	updateAvailabe, err := pacman.IsPackageUpdateAvailable(ctx, packageName, &packagemanageroptions.UpdateDatabaseOptions{})
	require.NoError(t, err)
	require.True(t, updateAvailabe)

	err = archlinuxutils.UpdateArchLinuxKeyringPackage(ctx, container, false)
	require.NoError(t, err)

	// The keyring should now be up to date and no further updates available:
	updateAvailabe, err = pacman.IsPackageUpdateAvailable(ctx, packageName, &packagemanageroptions.UpdateDatabaseOptions{})
	require.NoError(t, err)
	require.False(t, updateAvailabe)
}
