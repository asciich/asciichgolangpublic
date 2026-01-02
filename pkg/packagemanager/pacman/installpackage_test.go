package pacman_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/pacman"
)

func getPacmanInContainer(ctx context.Context, t *testing.T) (*pacman.Pacman, containerinterfaces.Container) {
	container, err := nativedocker.NewDocker().GetContainerByName("test-pacman-is-package-update-available")
	require.NoError(t, err)
	container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	err = container.Run(ctx, &dockeroptions.DockerRunContainerOptions{
		ImageName: "archlinux:base-20250727.0.390543",
		Command:   []string{"sleep", "1m"},
	})
	require.NoError(t, err)

	pacman, err := pacman.NewPacman(container)
	require.NoError(t, err)

	return pacman, container
}

func Test_InstallPackage(t *testing.T) {
	t.Run("empty package name", func(t *testing.T) {
		ctx := getCtx()
		pacman, container := getPacmanInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		err := pacman.InstallPackage(ctx, "", &packagemanageroptions.InstallPackageOptions{})
		require.Error(t, err)
	})

	t.Run("already installed", func(t *testing.T) {
		ctx := getCtx()
		pacman, container := getPacmanInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		// The archlinux-keyring package is already installed but there is an update available in the used container:
		const packageName = "archlinux-keyring"
		isInstalled, err := pacman.IsPackageInstalled(ctx, packageName)
		require.NoError(t, err)
		require.True(t, isInstalled)

		isUpdateAvailable, err := pacman.IsPackageUpdateAvailalbe(ctx, packageName, &packagemanageroptions.UpdateDatabaseOptions{})
		require.NoError(t, err)
		require.True(t, isUpdateAvailable)

		ctx = contextutils.WithChangeIndicator(ctx)
		err = pacman.InstallPackage(ctx, packageName, &packagemanageroptions.InstallPackageOptions{
			UpdateDatabaseFirst: false,
			Force:               false,
		})
		require.NoError(t, err)
		// Since already installed we do not expect any change:
		require.False(t, contextutils.IsChanged(ctx))

		// We expect the package still to be installed and not updated:
		isInstalled, err = pacman.IsPackageInstalled(ctx, packageName)
		require.NoError(t, err)
		require.True(t, isInstalled)

		isUpdateAvailable, err = pacman.IsPackageUpdateAvailalbe(ctx, packageName, &packagemanageroptions.UpdateDatabaseOptions{})
		require.NoError(t, err)
		require.True(t, isUpdateAvailable)
	})

	t.Run("vim", func(t *testing.T) {
		ctx := getCtx()
		pacman, container := getPacmanInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		const packageName = "vim"
		isInstalled, err := pacman.IsPackageInstalled(ctx, packageName)
		require.NoError(t, err)
		require.False(t, isInstalled)

		ctx = contextutils.WithChangeIndicator(ctx)
		err = pacman.InstallPackage(ctx, packageName, &packagemanageroptions.InstallPackageOptions{
			UpdateDatabaseFirst: true,
			Force:               false,
		})
		require.NoError(t, err)
		// Since the package got installed we do expect a change:
		require.True(t, contextutils.IsChanged(ctx))

		isInstalled, err = pacman.IsPackageInstalled(ctx, packageName)
		require.NoError(t, err)
		require.True(t, isInstalled)
	})
}
