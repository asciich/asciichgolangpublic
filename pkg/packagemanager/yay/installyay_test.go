package yay_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/dockerutils/nativedocker"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/yay"
)

func TestInstallYay(t *testing.T) {
	ctx := getCtx()

	const containerName = "test-install-yay"
	require.NoError(t, nativedocker.NewDocker().RemoveContainer(ctx, containerName, &dockeroptions.RemoveOptions{Force: true}))

	container, err := nativedocker.NewDocker().RunContainer(
		ctx, &dockeroptions.DockerRunContainerOptions{
			Name: containerName,
			ImageName: "archlinux",
			Command:   []string{"sleep", "2m"},
		},
	)
	require.NoError(t, err)
	defer container.Remove(ctx, &dockeroptions.RemoveOptions{Force: true})

	// On default arch linux yay is not installed:
	isInstalled, err := yay.IsInstalled(ctx, container)
	require.NoError(t, err)
	require.False(t, isInstalled)

	// Install yay:
	err = yay.InstallYay(ctx, container, &packagemanageroptions.InstallPackageOptions{
		UpdateDatabaseFirst: true,
		UseSudo: false, // Sudo is not installed by default on archlinux and therefore not available in the test container.
	})
	require.NoError(t, err)

	// Now yay is installed and available as command:
	isInstalled, err = yay.IsInstalled(ctx, container)
	require.NoError(t, err)
	require.True(t, isInstalled)
}
