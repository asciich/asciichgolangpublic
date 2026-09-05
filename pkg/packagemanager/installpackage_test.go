package packagemanager_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func getTestContainer(ctx context.Context, t *testing.T, imageName string) containerinterfaces.Container {
	const containerName = "test-packagemanager"

	logging.LogInfoByCtxf(ctx, "Going to start container '%s' using image '%s'.", containerName, imageName)
	container, err := nativedocker.NewDocker().GetContainerByName(containerName)
	require.NoError(t, err)
	container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	err = container.Run(ctx, &dockeroptions.DockerRunContainerOptions{
		ImageName: imageName,
		Command:   []string{"sleep", "2m"},
	})
	require.NoError(t, err)

	return container
}

func getImagesTotest() []string {
	return []string{
		"archlinux:base-20260419.0.517065",
		"ubuntu:noble-20260810",
	}
}

// getAlreadyInstalledPackage returns a package name that is guaranteed to be
// already installed in the given base image.
func getAlreadyInstalledPackage(t *testing.T, image string) string {
	switch image {
	case "archlinux:base-20260419.0.517065":
		// The filesystem package is always installed in Arch Linux:
		return "filesystem"
	case "ubuntu:noble-20260810":
		// The base-files package is always installed in Ubuntu:
		return "base-files"
	default:
		t.Fatalf("no already-installed package known for image '%s'", image)
		return ""
	}
}

func Test_InstallPackages_emptyPackage(t *testing.T) {
	for _, image := range getImagesTotest() {
		t.Run("empty package name", func(t *testing.T) {
			ctx := getCtx()
			container := getTestContainer(ctx, t, image)
			defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

			err := packagemanager.InstallPackages(ctx, container, []string{}, &packagemanageroptions.InstallPackageOptions{})
			require.Error(t, err)
		})
	}
}

func Test_InstallPackages_alreadyInstalled(t *testing.T) {
	for _, image := range getImagesTotest() {
		t.Run("already installed", func(t *testing.T) {
			ctx := getCtx()
			container := getTestContainer(ctx, t, image)
			defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

			packageName := getAlreadyInstalledPackage(t, image)

			isInstalled, err := packagemanager.IsPackageInstalled(ctx, container, packageName)
			require.NoError(t, err)
			require.True(t, isInstalled)

			// Installing an already installed package must be a no-op and must not fail.
			ctx = contextutils.WithChangeIndicator(ctx)
			err = packagemanager.InstallPackages(ctx, container, []string{packageName}, &packagemanageroptions.InstallPackageOptions{
				UpdateDatabaseFirst: false,
				Force:               false,
			})
			require.NoError(t, err)
			// Idempotence: since already installed we do not expect any change:
			require.False(t, contextutils.IsChanged(ctx))

			// We expect the package still to be installed:
			isInstalled, err = packagemanager.IsPackageInstalled(ctx, container, packageName)
			require.NoError(t, err)
			require.True(t, isInstalled)
		})
	}
}

func Test_InstallPackages_vim(t *testing.T) {
	for _, image := range getImagesTotest() {
		t.Run("vim", func(t *testing.T) {
			ctx := getCtx()
			container := getTestContainer(ctx, t, image)
			defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

			const packageName = "vim"
			isInstalled, err := packagemanager.IsPackageInstalled(ctx, container, packageName)
			require.NoError(t, err)
			require.False(t, isInstalled)

			ctx = contextutils.WithChangeIndicator(ctx)
			err = packagemanager.InstallPackages(ctx, container, []string{packageName}, &packagemanageroptions.InstallPackageOptions{
				UpdateDatabaseFirst: true,
				Force:               false,
			})
			require.NoError(t, err)
			// Since the package got installed we do expect a change:
			require.True(t, contextutils.IsChanged(ctx))

			isInstalled, err = packagemanager.IsPackageInstalled(ctx, container, packageName)
			require.NoError(t, err)
			require.True(t, isInstalled)
		})
	}
}
