package aptget_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/aptget"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
)

func getAptGetInContainer(ctx context.Context, t *testing.T) (*aptget.AptGet, containerinterfaces.Container) {
	container, err := nativedocker.NewDocker().GetContainerByName("test-aptget-is-package-update-available")
	require.NoError(t, err)
	container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	err = container.Run(ctx, &dockeroptions.DockerRunContainerOptions{
		// Pinned dated tag (equivalent to the immutable archlinux tag) so the
		// "update available" assumption stays reproducible over time:
		ImageName: "ubuntu:noble-20260810",
		Command:   []string{"sleep", "2m"},
	})
	require.NoError(t, err)

	aptGet, err := aptget.NewAptGet(container)
	require.NoError(t, err)

	return aptGet, container
}

func Test_InstallPackage(t *testing.T) {
	t.Run("empty package name", func(t *testing.T) {
		ctx := getCtx()
		aptGet, container := getAptGetInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		err := aptGet.InstallPackages(ctx, []string{}, &packagemanageroptions.InstallPackageOptions{})
		require.Error(t, err)
	})
	t.Run("already installed", func(t *testing.T) {
		ctx := getCtx()
		aptGet, container := getAptGetInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		// base-files is guaranteed to be already installed in the base image.
		const packageName = "base-files"

		isInstalled, err := aptGet.IsPackageInstalled(ctx, packageName)
		require.NoError(t, err)
		require.True(t, isInstalled)

		// Installing an already installed package must be a no-op and must not fail.
		ctx = contextutils.WithChangeIndicator(ctx)
		err = aptGet.InstallPackages(ctx, []string{packageName}, &packagemanageroptions.InstallPackageOptions{
			UpdateDatabaseFirst: false,
			Force:               false,
		})
		require.NoError(t, err)

		// Idempotence: nothing changed because it was already installed.
		require.False(t, contextutils.IsChanged(ctx))

		// Still installed afterwards.
		isInstalled, err = aptGet.IsPackageInstalled(ctx, packageName)
		require.NoError(t, err)
		require.True(t, isInstalled)
	})

	t.Run("vim", func(t *testing.T) {
		ctx := getCtx()
		aptGet, container := getAptGetInContainer(ctx, t)
		defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

		const packageName = "vim"
		isInstalled, err := aptGet.IsPackageInstalled(ctx, packageName)
		require.NoError(t, err)
		require.False(t, isInstalled)

		ctx = contextutils.WithChangeIndicator(ctx)
		err = aptGet.InstallPackages(ctx, []string{packageName}, &packagemanageroptions.InstallPackageOptions{
			UpdateDatabaseFirst: true,
			Force:               false,
		})
		require.NoError(t, err)
		// Since the package got installed we do expect a change:
		require.True(t, contextutils.IsChanged(ctx))

		isInstalled, err = aptGet.IsPackageInstalled(ctx, packageName)
		require.NoError(t, err)
		require.True(t, isInstalled)
	})
}
