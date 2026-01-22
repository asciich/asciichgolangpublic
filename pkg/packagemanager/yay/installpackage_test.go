package yay_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/yay"
)

func getYayInContainer(ctx context.Context, t *testing.T) (*yay.Yay, containerinterfaces.Container) {
	container, err := nativedocker.NewDocker().GetContainerByName("test-yay-is-package-update-available")
	require.NoError(t, err)
	err = container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})
	require.NoError(t, err)

	err = container.Run(ctx, &dockeroptions.DockerRunContainerOptions{
		ImageName: "archlinux:base-20250727.0.390543",
		Command:   []string{"sleep", "2m"},
	})
	require.NoError(t, err)

	yay, err := yay.NewYay(container)
	require.NoError(t, err)

	err = yay.InstallYay(ctx, &packagemanageroptions.InstallPackageOptions{UpdateDatabaseFirst: true})
	require.NoError(t, err)

	return yay, container
}

func Test_InstallPackageYay(t *testing.T) {
	t.Run("empty package name", func(t *testing.T) {
		ctx := getCtx()
		yay, container := getYayInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		err := yay.InstallPackages(ctx, []string{}, &packagemanageroptions.InstallPackageOptions{})
		require.Error(t, err)
	})

	t.Run("already installed", func(t *testing.T) {
		ctx := getCtx()
		yay, container := getYayInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		// The archlinux-keyring package is already installed but there is an update available in the used container:
		const packageName = "archlinux-keyring"
		isInstalled, err := yay.IsPackageInstalled(ctx, packageName)
		require.NoError(t, err)
		require.True(t, isInstalled)

		isUpdateAvailable, err := yay.IsPackageUpdateAvailable(ctx, packageName, &packagemanageroptions.UpdateDatabaseOptions{})
		require.NoError(t, err)
		require.True(t, isUpdateAvailable)

		ctx = contextutils.WithChangeIndicator(ctx)
		err = yay.InstallPackages(ctx, []string{packageName}, &packagemanageroptions.InstallPackageOptions{
			UpdateDatabaseFirst: false,
			Force:               false,
		})
		require.NoError(t, err)
		// Since already installed we do not expect any change:
		require.False(t, contextutils.IsChanged(ctx))

		// We expect the package still to be installed and not updated:
		isInstalled, err = yay.IsPackageInstalled(ctx, packageName)
		require.NoError(t, err)
		require.True(t, isInstalled)

		isUpdateAvailable, err = yay.IsPackageUpdateAvailable(ctx, packageName, &packagemanageroptions.UpdateDatabaseOptions{})
		require.NoError(t, err)
		require.True(t, isUpdateAvailable)
	})

	t.Run("vim", func(t *testing.T) {
		ctx := getCtx()
		yay, container := getYayInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		const packageName = "vim"
		isInstalled, err := yay.IsPackageInstalled(ctx, packageName)
		require.NoError(t, err)
		require.False(t, isInstalled)

		ctx = contextutils.WithChangeIndicator(ctx)
		err = yay.InstallPackages(ctx, []string{packageName}, &packagemanageroptions.InstallPackageOptions{
			UpdateDatabaseFirst: true,
			Force:               false,
		})
		require.NoError(t, err)
		// Since the package got installed we do expect a change:
		require.True(t, contextutils.IsChanged(ctx))

		isInstalled, err = yay.IsPackageInstalled(ctx, packageName)
		require.NoError(t, err)
		require.True(t, isInstalled)
	})
}
